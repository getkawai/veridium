// Pipeline Example
// This example demonstrates how to use the pipeline feature
package main

import (
	"fmt"
	"log"

	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
	"github.com/kawai-network/veridium/pkg/stablediffusion/pipeline"
)

func main() {
	fmt.Println("Stable Diffusion Go - Pipeline Example")
	fmt.Println("=======================================")

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

	// Example 1: Simple Text-to-Image Pipeline
	fmt.Println("\n1. Simple Text-to-Image Pipeline...")
	p1 := pipeline.New().
		Add(&pipeline.TextToImageStep{
			Params: stablediffusion.ImgGenParams{
				Prompt:      "a beautiful mountain landscape",
				Width:       512,
				Height:      512,
				SampleSteps: 10,
				CfgScale:    1.0,
				Seed:        42,
			},
			OutputName: "generated.png",
		}).
		Add(&pipeline.SaveStep{
			Destination: "outputs/pipeline_simple.png",
		})

	ctx1 := pipeline.NewContext("./working")
	ctx1.SD = sd

	if err := p1.Execute(ctx1); err != nil {
		log.Fatalf("Pipeline execution failed: %v", err)
	}
	fmt.Println("✓ Pipeline executed successfully!")

	// Example 2: Text-to-Image with Upscaling Pipeline
	fmt.Println("\n2. Text-to-Image with Upscaling Pipeline...")
	p2 := pipeline.New().
		Add(&pipeline.TextToImageStep{
			Params: stablediffusion.ImgGenParams{
				Prompt:      "a futuristic city at night",
				Width:       512,
				Height:      512,
				SampleSteps: 10,
				CfgScale:    1.0,
				Seed:        123,
			},
			OutputName: "generated.png",
		}).
		Add(&pipeline.UpscaleStep{
			Params: stablediffusion.UpscalerParams{
				EsrganPath: "models/RealESRGAN_x4plus.pth",
				NThreads:   4,
			},
			Factor:     4,
			OutputName: "upscaled.png",
		}).
		Add(&pipeline.SaveStep{
			Destination: "outputs/pipeline_upscaled.png",
		})

	ctx2 := pipeline.NewContext("./working")
	ctx2.SD = sd

	if err := p2.Execute(ctx2); err != nil {
		log.Printf("Note: Upscale step may fail without ESRGAN model: %v", err)
	} else {
		fmt.Println("✓ Pipeline executed successfully!")
	}

	// Example 3: Image-to-Image Pipeline
	fmt.Println("\n3. Image-to-Image Pipeline...")
	p3 := pipeline.New().
		Add(&pipeline.LoadStep{
			Source: "outputs/pipeline_simple.png",
		}).
		Add(&pipeline.ImageToImageStep{
			Params: stablediffusion.ImgGenParams{
				Prompt:      "a beautiful mountain landscape at sunset",
				SampleSteps: 15,
				CfgScale:    5.0,
				Strength:    0.75,
				Seed:        456,
			},
			OutputName: "img2img_output.png",
		}).
		Add(&pipeline.SaveStep{
			Destination: "outputs/pipeline_img2img.png",
		})

	ctx3 := pipeline.NewContext("./working")
	ctx3.SD = sd

	if err := p3.Execute(ctx3); err != nil {
		log.Printf("Pipeline execution failed: %v", err)
	} else {
		fmt.Println("✓ Pipeline executed successfully!")
	}

	// Example 4: Using Preset Pipelines
	fmt.Println("\n4. Using Preset Pipelines...")

	// Simple txt2img using preset
	p4 := pipeline.Txt2ImgPipeline(
		sd,
		stablediffusion.ImgGenParams{
			Prompt:      "an abstract painting with vibrant colors",
			Width:       512,
			Height:      512,
			SampleSteps: 10,
			CfgScale:    1.0,
			Seed:        789,
		},
		"outputs/pipeline_preset.png",
	)

	ctx4 := pipeline.NewContext("./working")
	ctx4.SD = sd

	if err := p4.Execute(ctx4); err != nil {
		log.Fatalf("Preset pipeline execution failed: %v", err)
	}
	fmt.Println("✓ Preset pipeline executed successfully!")

	// Example 5: Pipeline with Error Handling
	fmt.Println("\n5. Pipeline with Error Handling...")
	p5 := pipeline.NewWithConfig(&pipeline.Config{
		StopOnError:      true,
		KeepIntermediate: true,
	}).
		Add(&pipeline.TextToImageStep{
			Params: stablediffusion.ImgGenParams{
				Prompt:      "a serene lake with mountains",
				Width:       512,
				Height:      512,
				SampleSteps: 10,
				CfgScale:    1.0,
			},
			OutputName: "step1.png",
		}).
		OnError(func(err error) error {
			fmt.Printf("Error occurred: %v\n", err)
			return err // Return nil to continue, or err to stop
		})

	ctx5 := pipeline.NewContext("./working")
	ctx5.SD = sd

	if err := p5.Execute(ctx5); err != nil {
		log.Printf("Pipeline execution failed: %v", err)
	} else {
		fmt.Println("✓ Pipeline executed successfully!")
	}

	// Example 6: Pipeline Builder Pattern
	fmt.Println("\n6. Using Pipeline Builder...")
	p6 := pipeline.NewBuilder().
		StopOnError(true).
		KeepIntermediate(false).
		Then(&pipeline.TextToImageStep{
			Params: stablediffusion.ImgGenParams{
				Prompt:      "a magical forest with glowing mushrooms",
				Width:       512,
				Height:      512,
				SampleSteps: 10,
				CfgScale:    1.0,
				Seed:        999,
			},
			OutputName: "generated.png",
		}).
		Then(&pipeline.SaveStep{
			Destination: "outputs/pipeline_builder.png",
		}).
		Build()

	ctx6 := pipeline.NewContext("./working")
	ctx6.SD = sd

	if err := p6.Execute(ctx6); err != nil {
		log.Fatalf("Builder pipeline execution failed: %v", err)
	}
	fmt.Println("✓ Builder pipeline executed successfully!")

	// Example 7: Complex Multi-Step Pipeline
	fmt.Println("\n7. Complex Multi-Step Pipeline...")
	fmt.Println("   (Load -> Generate -> Upscale -> Save)")

	p7 := pipeline.New().
		Add(&pipeline.TextToImageStep{
			Params: stablediffusion.ImgGenParams{
				Prompt:      "a steampunk airship flying over clouds",
				Width:       512,
				Height:      512,
				SampleSteps: 20,
				CfgScale:    5.0,
				Seed:        111,
			},
			OutputName: "step1_generate.png",
		}).
		Add(&pipeline.SaveStep{
			Destination: "outputs/pipeline_complex_step1.png",
		}).
		Add(&pipeline.SaveStep{
			Destination: "outputs/pipeline_complex_final.png",
		})

	ctx7 := pipeline.NewContext("./working")
	ctx7.SD = sd

	if err := p7.Execute(ctx7); err != nil {
		log.Printf("Complex pipeline execution failed: %v", err)
	} else {
		fmt.Println("✓ Complex pipeline executed successfully!")
	}

	fmt.Println("\n✅ All pipeline examples completed!")
}
