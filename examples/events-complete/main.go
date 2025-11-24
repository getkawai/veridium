package main

import (
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed assets
var assets embed.FS

// Typed event data structures
type UserData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type NotificationData struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

type ProgressData struct {
	Task    string  `json:"task"`
	Percent float64 `json:"percent"`
	Status  string  `json:"status"`
}

// Event void type untuk events tanpa data
type Void application.Void

func init() {
	// Register events dengan type safety
	application.RegisterEvent[UserData]("user-updated")
	application.RegisterEvent[NotificationData]("notification")
	application.RegisterEvent[ProgressData]("task-progress")
	application.RegisterEvent[Void]("app-shutdown")
}

func main() {
	app := application.New(application.Options{
		Name:        "Event System Demo",
		Description: "Demonstrates complete event system functionality",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	setupEventHandlers(app)

	win := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Event System Demo",
		Name:  "main",
	})

	setupWindowEventHandlers(win, app)

	startBackgroundTasks(app)

	err := app.Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func setupEventHandlers(app *application.App) {
	// Hook system - intercept events before dispatch
	app.Event.RegisterHook("notification", func(e *application.CustomEvent) {
		data := e.Data.(NotificationData)

		// Cancel high priority notifications if not in business hours
		if data.Priority > 7 && !isBusinessHours() {
			fmt.Printf("Cancelling high priority notification outside business hours: %s\n", data.Message)
			e.Cancel()
			return
		}

		fmt.Printf("Hook processed notification: %s\n", data.Message)
	})

	// Regular event listeners
	app.Event.On("user-updated", func(e *application.CustomEvent) {
		userData := e.Data.(UserData)
		app.Logger.Info("User updated", "user", userData)

		// Emit related event
		app.Event.Emit("user-log", fmt.Sprintf("User %s was updated", userData.Name))
	})

	app.Event.On("notification", func(e *application.CustomEvent) {
		data := e.Data.(NotificationData)
		app.Logger.Info("Notification received", "type", data.Type, "priority", data.Priority)
	})

	app.Event.On("task-progress", func(e *application.CustomEvent) {
		data := e.Data.(ProgressData)
		app.Logger.Info("Task progress", "task", data.Task, "progress", data.Percent)
	})

	// Once listener - only fires once
	app.Event.On("app-startup", func(e *application.CustomEvent) {
		app.Logger.Info("App startup event received!")
	})

	// Multiple times listener
	callCount := 0
	app.Event.OnMultiple("limited-event", func(e *application.CustomEvent) {
		callCount++
		app.Logger.Info("Limited event", "count", callCount, "data", e.Data)
	}, 3) // Will only fire 3 times

	// Application events
	app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(event *application.ApplicationEvent) {
		app.Logger.Info("Application started!")

		// Check system theme
		if event.Context().IsDarkMode() {
			app.Logger.Info("System is in dark mode")
		}

		// Emit startup event
		app.Event.Emit("app-startup", "App is ready!")
	})

	app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(event *application.ApplicationEvent) {
		if event.Context().IsDarkMode() {
			app.Logger.Info("Switched to dark mode")
			app.Event.Emit("theme-changed", "dark")
		} else {
			app.Logger.Info("Switched to light mode")
			app.Event.Emit("theme-changed", "light")
		}
	})
}

func setupWindowEventHandlers(win *application.WebviewWindow, app *application.App) {
	// Window closing dengan hook - can prevent closing
	win.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) {
		app.Logger.Info("Window closing...")

		// Check if there are unsaved changes
		if hasUnsavedChanges() {
			app.Logger.Info("Has unsaved changes - cancelling close")
			e.Cancel()
			return
		}

		app.Logger.Info("Window can close safely")
	})

	// Window focus events
	win.OnWindowEvent(events.Common.WindowFocus, func(e *application.WindowEvent) {
		app.Logger.Info("Window gained focus")
		app.Event.Emit("window-focused", win.Name())
	})

	win.OnWindowEvent(events.Common.WindowLostFocus, func(e *application.WindowEvent) {
		app.Logger.Info("Window lost focus")
	})

	// Drag and drop
	win.OnWindowEvent(events.Common.WindowFilesDropped, func(e *application.WindowEvent) {
		files := e.Context().DroppedFiles()
		app.Logger.Info("Files dropped", "count", len(files))

		for _, file := range files {
			app.Event.Emit("file-processed", fmt.Sprintf("Processed: %s", file))
		}
	})

	// Custom window events
	win.OnWindowEvent(events.Common.WindowDidMove, func(e *application.WindowEvent) {
		app.Logger.Info("Window moved")
	})

	win.OnWindowEvent(events.Common.WindowDidResize, func(e *application.WindowEvent) {
		app.Logger.Info("Window resized")
	})
}

func startBackgroundTasks(app *application.App) {
	// Background task - emit periodic progress
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		tasks := []string{"Loading data", "Processing files", "Updating UI", "Saving state"}
		taskIndex := 0

		for {
			select {
			case <-ticker.C:
				progress := ProgressData{
					Task:    tasks[taskIndex%len(tasks)],
					Percent: float64(taskIndex%10) * 10.0,
					Status:  "in-progress",
				}

				app.Event.Emit("task-progress", progress)

				if taskIndex%10 == 0 {
					// Emit user update every 20 seconds
					user := UserData{
						ID:   taskIndex / 10,
						Name: fmt.Sprintf("User%d", taskIndex/10),
						Age:  25 + (taskIndex / 10),
					}
					app.Event.Emit("user-updated", user)
				}

				if taskIndex%15 == 0 {
					// Emit notification
					notification := NotificationData{
						Type:     "system",
						Message:  "Background task completed",
						Priority: 5,
					}
					app.Event.Emit("notification", notification)
				}

				taskIndex++

			case <-app.Context().Done():
				// App is shutting down
				app.Event.Emit("app-shutdown", Void{})
				return
			}
		}
	}()

	// Limited event test
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		callIndex := 0
		for {
			select {
			case <-ticker.C:
				callIndex++
				app.Logger.Info("Emitting limited event", "index", callIndex)
				app.Event.Emit("limited-event", fmt.Sprintf("Call #%d", callIndex))

				if callIndex >= 5 {
					return // Stop after 5 calls (only 3 will be processed)
				}

			case <-app.Context().Done():
				return
			}
		}
	}()
}

// Helper functions
func isBusinessHours() bool {
	now := time.Now()
	hour := now.Hour()
	return hour >= 9 && hour <= 17 // 9 AM to 5 PM
}

func hasUnsavedChanges() bool {
	// Simulate checking for unsaved changes
	return time.Now().Second()%3 == 0 // Random for demo
}
