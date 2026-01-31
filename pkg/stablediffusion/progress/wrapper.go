package progress

import (
	"fmt"
	"time"

	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
	"github.com/kawai-network/veridium/pkg/stablediffusion/sd"
)

// StableDiffusionWrapper wraps StableDiffusion with progress tracking
type StableDiffusionWrapper struct {
	sd      *stablediffusion.StableDiffusion
	options *GenerationOptions
}

// NewWrapper creates a new wrapper around StableDiffusion
func NewWrapper(sd *stablediffusion.StableDiffusion, options *GenerationOptions) *StableDiffusionWrapper {
	if options == nil {
		options = DefaultGenerationOptions()
	}
	return &StableDiffusionWrapper{
		sd:      sd,
		options: options,
	}
}

// GenerateImage generates an image with progress tracking
func (w *StableDiffusionWrapper) GenerateImage(
	params *stablediffusion.ImgGenParams,
	outputPath string,
) error {
	tracker := NewProgressTracker(int(params.SampleSteps), w.options.ProgressCallback, w.options.CancelToken)

	// Stage: Initializing
	tracker.SetStage(StageInitializing)
	tracker.SetMessage("Initializing generation...")

	if tracker.IsCancelled() {
		return fmt.Errorf("generation cancelled")
	}

	// Stage: Loading Model (simulated - actual loading happens in NewStableDiffusion)
	tracker.SetStage(StageLoadingModel)
	tracker.SetMessage("Loading model...")

	if tracker.IsCancelled() {
		return fmt.Errorf("generation cancelled")
	}

	// Stage: Encoding Prompt
	tracker.SetStage(StageEncodingPrompt)
	tracker.SetMessage("Encoding prompt...")

	if tracker.IsCancelled() {
		return fmt.Errorf("generation cancelled")
	}

	// Stage: Generating with step tracking
	tracker.SetStage(StageGenerating)

	// Create a callback that updates the tracker
	progressCallback := func(step int, steps int, t float32, data interface{}) {
		tracker.SetStep(step)
		tracker.SetMessage(fmt.Sprintf("Generating step %d/%d", step, steps))
	}

	// Set the progress callback in the underlying sd
	sd.SetProgressCallback(func(step int, steps int, t float32, data interface{}) {
		progressCallback(step, steps, t, data)
	}, nil)

	// Stage: Decoding
	tracker.SetStage(StageDecoding)
	tracker.SetMessage("Decoding latent...")

	// Stage: Saving
	tracker.SetStage(StageSaving)
	tracker.SetMessage("Saving image...")

	// Perform actual generation
	if err := w.sd.GenerateImage(params, outputPath); err != nil {
		tracker.SetError(err)
		return err
	}

	// Complete
	tracker.Complete()
	return nil
}

// GenerateVideo generates a video with progress tracking
func (w *StableDiffusionWrapper) GenerateVideo(
	params *stablediffusion.VidGenParams,
	outputPath string,
) error {
	tracker := NewProgressTracker(int(params.SampleSteps), w.options.ProgressCallback, w.options.CancelToken)

	// Stage: Initializing
	tracker.SetStage(StageInitializing)
	tracker.SetMessage("Initializing video generation...")

	if tracker.IsCancelled() {
		return fmt.Errorf("generation cancelled")
	}

	// Stage: Loading Model
	tracker.SetStage(StageLoadingModel)
	tracker.SetMessage("Loading video model...")

	if tracker.IsCancelled() {
		return fmt.Errorf("generation cancelled")
	}

	// Stage: Encoding Prompt
	tracker.SetStage(StageEncodingPrompt)
	tracker.SetMessage("Encoding video prompt...")

	if tracker.IsCancelled() {
		return fmt.Errorf("generation cancelled")
	}

	// Stage: Generating
	tracker.SetStage(StageGenerating)
	tracker.SetMessage("Generating video frames...")

	// Set progress callback
	sd.SetProgressCallback(func(step int, steps int, t float32, data interface{}) {
		tracker.SetStep(step)
		tracker.SetMessage(fmt.Sprintf("Generating frame %d/%d", step, steps))
	}, nil)

	// Stage: Decoding
	tracker.SetStage(StageDecoding)
	tracker.SetMessage("Decoding video frames...")

	// Stage: Saving
	tracker.SetStage(StageSaving)
	tracker.SetMessage("Saving video...")

	// Perform actual generation
	if err := w.sd.GenerateVideo(params, outputPath); err != nil {
		tracker.SetError(err)
		return err
	}

	// Complete
	tracker.Complete()
	return nil
}

// SimpleProgressBar creates a simple text-based progress bar
func SimpleProgressBar(width int) func(progress GenerationProgress) {
	return func(progress GenerationProgress) {
		filled := int(float64(width) * progress.Percentage / 100.0)
		empty := width - filled

		bar := "["
		for i := 0; i < filled; i++ {
			bar += "="
		}
		if filled < width {
			bar += ">"
			empty--
		}
		for i := 0; i < empty; i++ {
			bar += " "
		}
		bar += "]"

		fmt.Printf("\r%s %s %s (%s/%s) %s",
			progress.Stage.String(),
			bar,
			progress.FormatPercentage(1),
			formatDuration(progress.TimeElapsed),
			formatDuration(progress.TimeRemaining),
			progress.Message,
		)

		if progress.IsComplete() {
			fmt.Println()
		}
	}
}

// formatDuration formats a duration for display
func formatDuration(d time.Duration) string {
	if d == 0 {
		return "--:--"
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// CreateProgressCallback creates a progress callback with custom formatting
func CreateProgressCallback(
	onProgress func(progress GenerationProgress),
	onComplete func(progress GenerationProgress),
	onError func(progress GenerationProgress, err error),
) ProgressCallback {
	return func(progress GenerationProgress) {
		onProgress(progress)

		if progress.IsComplete() {
			if progress.Error != nil && onError != nil {
				onError(progress, progress.Error)
			} else if onComplete != nil {
				onComplete(progress)
			}
		}
	}
}
