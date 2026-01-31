// Progress Tracking Example
// This example demonstrates how to use the progress tracking feature
package main

import (
	"fmt"
	"log"

	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
	"github.com/kawai-network/veridium/pkg/stablediffusion/progress"
)

func main() {
	fmt.Println("Stable Diffusion Go - Progress Tracking Example")
	fmt.Println("================================================")

	// Initialize Stable Diffusion
	sd, err := stablediffusion.NewStableDiffusion(&stablediffusion.ContextParams{
		DiffusionModelPath: "D:\\hf-mirror\\Z-Image-Turbo-GGUF\\z_image_turbo-Q4_K_M.gguf",
		LLMPath:            "D:\\hf-mirror\\Z-Image-Turbo-GGUF\\Qwen3-4B-Instruct-2507-Q4_K_M.gguf",
		VAEPath:            "D:\\hf-mirror\\Z-Image-Turbo-GGUF\\diffusion_pytorch_model.safetensors",
		DiffusionFlashAttn: true,
		OffloadParamsToCPU: true,
	})
	if err != nil {
		log.Fatalf("Failed to create Stable Diffusion instance: %v", err)
	}
	defer sd.Free()

	// Example 1: Simple progress bar
	fmt.Println("\n1. Generating image with simple progress bar...")
	options1 := &progress.GenerationOptions{
		ProgressCallback: progress.SimpleProgressBar(30),
	}

	wrapper1 := progress.NewWrapper(sd, options1)
	err = wrapper1.GenerateImage(&stablediffusion.ImgGenParams{
		Prompt:      "a serene mountain landscape at sunrise",
		Width:       512,
		Height:      512,
		SampleSteps: 20,
		CfgScale:    5.0,
		Seed:        42,
	}, "outputs/progress_simple.png")

	if err != nil {
		log.Fatalf("Generation failed: %v", err)
	}
	fmt.Println("✓ Image generated successfully!")

	// Example 2: Custom progress callback
	fmt.Println("\n2. Generating image with custom progress callback...")
	options2 := &progress.GenerationOptions{
		ProgressCallback: func(p progress.GenerationProgress) {
			fmt.Printf("\r[%s] Step %d/%d | %.1f%% | Elapsed: %s | Remaining: %s | %s",
				p.Stage,
				p.Step,
				p.TotalSteps,
				p.Percentage,
				formatDuration(p.TimeElapsed),
				formatDuration(p.TimeRemaining),
				p.Message,
			)
			if p.IsComplete() {
				fmt.Println()
			}
		},
	}

	wrapper2 := progress.NewWrapper(sd, options2)
	err = wrapper2.GenerateImage(&stablediffusion.ImgGenParams{
		Prompt:      "a futuristic city with flying cars",
		Width:       512,
		Height:      512,
		SampleSteps: 25,
		CfgScale:    7.0,
		Seed:        123,
	}, "outputs/progress_custom.png")

	if err != nil {
		log.Fatalf("Generation failed: %v", err)
	}
	fmt.Println("✓ Image generated successfully!")

	// Example 3: Progress with cancellation support
	fmt.Println("\n3. Generating image with cancellation support...")
	cancelToken := progress.NewCancelToken()

	options3 := &progress.GenerationOptions{
		ProgressCallback: func(p progress.GenerationProgress) {
			fmt.Printf("\r[%s] %.1f%% - %s", p.Stage, p.Percentage, p.Message)
			if p.IsComplete() {
				fmt.Println()
			}
		},
		CancelToken: cancelToken,
	}

	wrapper3 := progress.NewWrapper(sd, options3)

	// In a real application, you might cancel based on user input or timeout
	// For this example, we'll just let it complete
	err = wrapper3.GenerateImage(&stablediffusion.ImgGenParams{
		Prompt:      "an abstract painting with vibrant colors",
		Width:       512,
		Height:      512,
		SampleSteps: 15,
		CfgScale:    5.0,
		Seed:        456,
	}, "outputs/progress_cancel.png")

	if err != nil {
		log.Fatalf("Generation failed: %v", err)
	}
	fmt.Println("✓ Image generated successfully!")

	// Example 4: Using ProgressTracker directly
	fmt.Println("\n4. Using ProgressTracker directly...")
	tracker := progress.NewProgressTracker(20, func(p progress.GenerationProgress) {
		fmt.Printf("\rProgress: %.1f%% [%s] %s", p.Percentage, p.Stage, p.Message)
		if p.IsComplete() {
			fmt.Println()
		}
	}, nil)

	tracker.SetStage(progress.StageInitializing)
	tracker.SetMessage("Starting generation...")

	tracker.SetStage(progress.StageLoadingModel)
	tracker.SetMessage("Loading model...")

	tracker.SetStage(progress.StageEncodingPrompt)
	tracker.SetMessage("Encoding prompt...")

	tracker.SetStage(progress.StageGenerating)
	for i := 0; i <= 20; i++ {
		tracker.SetStep(i)
		tracker.SetMessage(fmt.Sprintf("Generating step %d/20", i))
	}

	tracker.SetStage(progress.StageDecoding)
	tracker.SetMessage("Decoding latent...")

	tracker.SetStage(progress.StageSaving)
	tracker.SetMessage("Saving image...")

	tracker.Complete()
	fmt.Println("✓ Progress tracking demo completed!")

	fmt.Println("\n✅ All progress tracking examples completed!")
}

func formatDuration(d interface{}) string {
	// Helper function to format duration for display
	return "00:00"
}
