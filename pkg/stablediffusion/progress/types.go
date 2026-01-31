// Package progress provides enhanced progress tracking for stable-diffusion-go
package progress

import (
	"fmt"
	"time"
	"sync/atomic"
)

// GenerationStage represents the current stage of generation
type GenerationStage int

const (
	// StageInitializing - Initial setup
	StageInitializing GenerationStage = iota
	// StageLoadingModel - Loading the model
	StageLoadingModel
	// StageEncodingPrompt - Encoding the prompt
	StageEncodingPrompt
	// StageGenerating - Main generation loop
	StageGenerating
	// StageDecoding - Decoding the latent
	StageDecoding
	// StageSaving - Saving the output
	StageSaving
	// StageComplete - Generation complete
	StageComplete
	// StageError - Error occurred
	StageError
	// StageCancelled - Generation was cancelled
	StageCancelled
)

// String returns the string representation of the stage
func (s GenerationStage) String() string {
	switch s {
	case StageInitializing:
		return "Initializing"
	case StageLoadingModel:
		return "Loading Model"
	case StageEncodingPrompt:
		return "Encoding Prompt"
	case StageGenerating:
		return "Generating"
	case StageDecoding:
		return "Decoding"
	case StageSaving:
		return "Saving"
	case StageComplete:
		return "Complete"
	case StageError:
		return "Error"
	case StageCancelled:
		return "Cancelled"
	default:
		return "Unknown"
	}
}

// GenerationProgress represents detailed progress information
type GenerationProgress struct {
	Stage         GenerationStage
	Step          int
	TotalSteps    int
	Percentage    float64
	TimeElapsed   time.Duration
	TimeRemaining time.Duration
	Message       string
	Error         error
	Timestamp     time.Time
}

// IsComplete returns true if the generation is complete (success, error, or cancelled)
func (p *GenerationProgress) IsComplete() bool {
	return p.Stage == StageComplete || p.Stage == StageError || p.Stage == StageCancelled
}

// IsSuccessful returns true if the generation completed successfully
func (p *GenerationProgress) IsSuccessful() bool {
	return p.Stage == StageComplete && p.Error == nil
}

// FormatPercentage returns the percentage as a formatted string
func (p *GenerationProgress) FormatPercentage(decimals int) string {
	format := fmt.Sprintf("%%.%df%%%%", decimals)
	return fmt.Sprintf(format, p.Percentage)
}

// ProgressCallback is the callback function type for progress updates
type ProgressCallback func(progress GenerationProgress)

// CancelToken provides cancellation support for generation
type CancelToken struct {
	cancelled atomic.Bool
}

// NewCancelToken creates a new cancel token
func NewCancelToken() *CancelToken {
	return &CancelToken{}
}

// Cancel cancels the operation
func (c *CancelToken) Cancel() {
	c.cancelled.Store(true)
}

// IsCancelled returns true if the operation has been cancelled
func (c *CancelToken) IsCancelled() bool {
	return c.cancelled.Load()
}

// GenerationOptions contains options for generation with progress tracking
type GenerationOptions struct {
	ProgressCallback ProgressCallback
	CancelToken      *CancelToken
	UpdateInterval   time.Duration // Minimum interval between progress updates
}

// DefaultGenerationOptions returns default generation options
func DefaultGenerationOptions() *GenerationOptions {
	return &GenerationOptions{
		UpdateInterval: 100 * time.Millisecond,
	}
}

// ProgressTracker tracks progress for a generation operation
type ProgressTracker struct {
	startTime      time.Time
	stage          GenerationStage
	step           int
	totalSteps     int
	callback       ProgressCallback
	updateInterval time.Duration
	lastUpdate     time.Time
	cancelToken    *CancelToken
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(totalSteps int, callback ProgressCallback, cancelToken *CancelToken) *ProgressTracker {
	return &ProgressTracker{
		startTime:      time.Now(),
		stage:          StageInitializing,
		totalSteps:     totalSteps,
		callback:       callback,
		updateInterval: 100 * time.Millisecond,
		lastUpdate:     time.Now(),
		cancelToken:    cancelToken,
	}
}

// SetStage updates the current stage
func (pt *ProgressTracker) SetStage(stage GenerationStage) {
	pt.stage = stage
	pt.update()
}

// SetStep updates the current step
func (pt *ProgressTracker) SetStep(step int) {
	pt.step = step
	pt.update()
}

// SetMessage sets a custom message
func (pt *ProgressTracker) SetMessage(message string) {
	pt.updateWithMessage(message)
}

// SetError sets an error state
func (pt *ProgressTracker) SetError(err error) {
	pt.stage = StageError
	pt.updateWithError(err)
}

// Complete marks the generation as complete
func (pt *ProgressTracker) Complete() {
	pt.stage = StageComplete
	pt.step = pt.totalSteps
	pt.update()
}

// IsCancelled checks if the operation should be cancelled
func (pt *ProgressTracker) IsCancelled() bool {
	if pt.cancelToken != nil && pt.cancelToken.IsCancelled() {
		pt.stage = StageCancelled
		return true
	}
	return false
}

// update triggers a progress update
func (pt *ProgressTracker) update() {
	pt.updateWithMessage("")
}

// updateWithMessage triggers a progress update with a message
func (pt *ProgressTracker) updateWithMessage(message string) {
	if pt.callback == nil {
		return
	}

	// Check update interval
	if time.Since(pt.lastUpdate) < pt.updateInterval {
		return
	}

	elapsed := time.Since(pt.startTime)
	percentage := pt.calculatePercentage()
	remaining := pt.calculateTimeRemaining(elapsed, percentage)

	progress := GenerationProgress{
		Stage:         pt.stage,
		Step:          pt.step,
		TotalSteps:    pt.totalSteps,
		Percentage:    percentage,
		TimeElapsed:   elapsed,
		TimeRemaining: remaining,
		Message:       message,
		Timestamp:     time.Now(),
	}

	pt.lastUpdate = time.Now()
	pt.callback(progress)
}

// updateWithError triggers a progress update with an error
func (pt *ProgressTracker) updateWithError(err error) {
	if pt.callback == nil {
		return
	}

	elapsed := time.Since(pt.startTime)

	progress := GenerationProgress{
		Stage:       pt.stage,
		Step:        pt.step,
		TotalSteps:  pt.totalSteps,
		Percentage:  pt.calculatePercentage(),
		TimeElapsed: elapsed,
		Error:       err,
		Timestamp:   time.Now(),
	}

	pt.callback(progress)
}

// calculatePercentage calculates the completion percentage
func (pt *ProgressTracker) calculatePercentage() float64 {
	if pt.totalSteps <= 0 {
		return 0
	}

	// Stage-based weighting
	stageWeights := map[GenerationStage]float64{
		StageInitializing:     5,
		StageLoadingModel:     15,
		StageEncodingPrompt:   5,
		StageGenerating:       60,
		StageDecoding:         10,
		StageSaving:           5,
		StageComplete:         100,
		StageError:            0,
		StageCancelled:        0,
	}

	basePercentage := stageWeights[pt.stage]

	// Add progress within the generating stage
	if pt.stage == StageGenerating && pt.totalSteps > 0 {
		stepProgress := (float64(pt.step) / float64(pt.totalSteps)) * stageWeights[StageGenerating]
		basePercentage = stageWeights[StageLoadingModel] + stageWeights[StageEncodingPrompt] + stepProgress
	}

	return min(basePercentage, 100.0)
}

// calculateTimeRemaining estimates the remaining time
func (pt *ProgressTracker) calculateTimeRemaining(elapsed time.Duration, percentage float64) time.Duration {
	if percentage <= 0 {
		return 0
	}

	// Estimate total time based on current progress
	totalEstimated := time.Duration(float64(elapsed) * 100.0 / percentage)
	remaining := totalEstimated - elapsed

	if remaining < 0 {
		return 0
	}

	return remaining
}

// min returns the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}