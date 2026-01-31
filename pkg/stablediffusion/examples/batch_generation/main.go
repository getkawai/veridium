// Batch Generation Example
// This example demonstrates how to use the batch generation feature
// to generate multiple images with different seeds
package main

import (
	"fmt"
	"log"
	"time"

	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
	"github.com/kawai-network/veridium/pkg/stablediffusion/batch"
)

func main() {
	fmt.Println("Stable Diffusion Go - Batch Generation Example")
	fmt.Println("===============================================")

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

	// Create batch generator
	batchGen := batch.NewGenerator(sd)

	// Example 1: Generate multiple images with different seeds
	fmt.Println("\n1. Generating 5 images with different seeds...")
	params1, err := batch.NewBuilder().
		WithBaseParams(stablediffusion.ImgGenParams{
			Prompt:      "a beautiful sunset over mountains",
			Width:       512,
			Height:      512,
			SampleSteps: 10,
			CfgScale:    1.0,
		}).
		AddSeedVariations("a beautiful sunset over mountains", 5).
		WithOutputPattern("outputs/sunset_%03d.png").
		WithParallelism(2).
		Build()

	if err != nil {
		log.Fatalf("Failed to build batch params: %v", err)
	}

	start := time.Now()
	result1, err := batchGen.GenerateWithCallback(params1, func(progress batch.BatchProgress) {
		fmt.Printf("\rProgress: %d/%d (%.1f%%) - ETA: %s",
			progress.Completed,
			progress.Total,
			progress.Percentage,
			formatDuration(progress.TimeRemaining))
	})
	if err != nil {
		log.Fatalf("Batch generation failed: %v", err)
	}

	fmt.Printf("\n✓ Generated %d images in %s\n", result1.SuccessCount, time.Since(start))
	for _, path := range result1.GetImages() {
		fmt.Printf("  - %s\n", path)
	}

	// Example 2: Generate with different resolutions
	fmt.Println("\n2. Generating images at different resolutions...")
	resolutions := [][2]int32{
		{512, 512},
		{768, 512},
		{512, 768},
	}

	params2, _ := batch.NewBuilder().
		WithBaseParams(stablediffusion.ImgGenParams{
			Prompt:      "a serene lake with mountains in background",
			SampleSteps: 10,
			CfgScale:    1.0,
		}).
		AddResolutionVariations("a serene lake with mountains in background", resolutions).
		WithOutputPattern("outputs/lake_%dx%d.png").
		Build()

	result2, err := batchGen.Generate(params2)
	if err != nil {
		log.Fatalf("Batch generation failed: %v", err)
	}

	fmt.Printf("✓ Generated %d images at different resolutions\n", result2.SuccessCount)

	// Example 3: Generate with different CFG scales
	fmt.Println("\n3. Generating images with different CFG scales...")
	cfgScales := []float32{1.0, 3.0, 5.0, 7.0}

	params3, _ := batch.NewBuilder().
		WithBaseParams(stablediffusion.ImgGenParams{
			Prompt:      "a cute cat playing with a ball",
			Width:       512,
			Height:      512,
			SampleSteps: 10,
			Seed:        42,
		}).
		AddCfgScaleVariations("a cute cat playing with a ball", cfgScales).
		WithOutputPattern("outputs/cat_cfg_%.1f.png").
		Build()

	result3, err := batchGen.Generate(params3)
	if err != nil {
		log.Fatalf("Batch generation failed: %v", err)
	}

	fmt.Printf("✓ Generated %d images with different CFG scales\n", result3.SuccessCount)

	// Example 4: Using presets
	fmt.Println("\n4. Using QuickTest preset...")
	presetBuilder := batch.NewPresetBuilder()
	params4, _ := presetBuilder.QuickTest("a futuristic city at night").
		WithOutputPattern("outputs/city_%03d.png").
		Build()

	result4, err := batchGen.Generate(params4)
	if err != nil {
		log.Fatalf("Batch generation failed: %v", err)
	}

	fmt.Printf("✓ QuickTest generated %d images\n", result4.SuccessCount)

	fmt.Println("\n✅ All batch generation examples completed!")
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "--:--"
	}
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
