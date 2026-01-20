package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kawai-network/veridium/pkg/grab"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: test-resume-download <url> <destination>")
		fmt.Println("\nExample:")
		fmt.Println("  test-resume-download https://huggingface.co/bartowski/Qwen_Qwen3-VL-8B-Instruct-GGUF/resolve/main/Qwen_Qwen3-VL-8B-Instruct-Q4_K_M.gguf /tmp/test-model.gguf")
		fmt.Println("\nTest scenario:")
		fmt.Println("  1. Run this command")
		fmt.Println("  2. Press Ctrl+C after a few seconds (partial download)")
		fmt.Println("  3. Run the same command again")
		fmt.Println("  4. Check if it resumes from where it left off")
		os.Exit(1)
	}

	url := os.Args[1]
	dest := os.Args[2]

	// Check if file already exists
	if info, err := os.Stat(dest); err == nil {
		fmt.Printf("📦 Found existing file: %s (%.2f MB)\n", dest, float64(info.Size())/(1024*1024))
		fmt.Printf("🔄 Will attempt to resume download...\n\n")
	} else {
		fmt.Printf("📥 Starting fresh download...\n\n")
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n\n⚠️  Interrupt received, canceling download...")
		cancel()
	}()

	// Create download request
	req, err := grab.NewRequest(dest, url)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req = req.WithContext(ctx)

	// Start download
	fmt.Printf("🚀 Downloading: %s\n", url)
	fmt.Printf("📁 Destination: %s\n\n", dest)

	client := grab.NewClient()
	resp := client.Do(req)

	// Progress monitoring
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	lastPrinted := time.Now()

	for {
		select {
		case <-ticker.C:
			if resp.IsComplete() {
				goto done
			}

			// Print progress every 500ms
			if time.Since(lastPrinted) >= 500*time.Millisecond {
				progress := resp.Progress() * 100
				speed := resp.BytesPerSecond() / (1024 * 1024) // MB/s
				downloaded := float64(resp.BytesComplete()) / (1024 * 1024)
				total := float64(resp.Size()) / (1024 * 1024)

				if resp.DidResume {
					fmt.Printf("\r🔄 RESUMED | %.1f%% | %.1f/%.1f MB | %.2f MB/s | ETA: %s",
						progress, downloaded, total, speed, resp.ETA().Sub(time.Now()).Round(time.Second))
				} else {
					fmt.Printf("\r📥 DOWNLOADING | %.1f%% | %.1f/%.1f MB | %.2f MB/s | ETA: %s",
						progress, downloaded, total, speed, resp.ETA().Sub(time.Now()).Round(time.Second))
				}
				lastPrinted = time.Now()
			}

		case <-resp.Done:
			goto done
		}
	}

done:
	fmt.Println() // New line after progress

	if err := resp.Err(); err != nil {
		if err == context.Canceled {
			fmt.Println("\n❌ Download canceled by user")
			if info, err := os.Stat(dest); err == nil {
				fmt.Printf("📊 Partial file saved: %.2f MB\n", float64(info.Size())/(1024*1024))
				fmt.Println("💡 Run the same command again to resume!")
			}
			os.Exit(0)
		}
		log.Fatalf("Download failed: %v", err)
	}

	// Success
	fmt.Println("\n✅ Download completed successfully!")
	if info, err := os.Stat(dest); err == nil {
		fmt.Printf("📊 Final size: %.2f MB\n", float64(info.Size())/(1024*1024))
	}

	if resp.DidResume {
		fmt.Printf("🔄 This download was RESUMED from %.2f MB\n", float64(resp.BytesComplete()-resp.Size())/(1024*1024))
	}

	fmt.Printf("⏱️  Total time: %s\n", resp.Duration().Round(time.Second))
	fmt.Printf("📈 Average speed: %.2f MB/s\n", resp.BytesPerSecond()/(1024*1024))
}
