// Configuration Example
// This example demonstrates how to use the configuration management feature
package main

import (
	"fmt"
	"log"

	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
	"github.com/kawai-network/veridium/pkg/stablediffusion/config"
)

func main() {
	fmt.Println("Stable Diffusion Go - Configuration Example")
	fmt.Println("============================================")

	// Example 1: Load configuration from file
	fmt.Println("\n1. Loading configuration from file...")
	cfg, err := config.LoadConfig("../../configs/example.yaml")
	if err != nil {
		log.Printf("Note: Could not load config file: %v", err)
		log.Println("Using default configuration instead...")
		cfg = config.DefaultConfig()
	}

	fmt.Printf("Config Version: %s\n", cfg.Version)
	fmt.Printf("Output Directory: %s\n", cfg.Output.OutputDir)
	fmt.Printf("Default Prompt: %s\n", cfg.Generation.DefaultPrompt)

	// Example 2: Create configuration programmatically
	fmt.Println("\n2. Creating configuration programmatically...")
	customCfg := &config.Config{
		Version: config.Version,
		Models: config.ModelConfig{
			DiffusionModelPath: "models/my_model.gguf",
			VAEPath:            "models/my_vae.safetensors",
		},
		Generation: config.GenerationConfig{
			DefaultPrompt:         "beautiful landscape",
			DefaultNegativePrompt: "blurry",
			DefaultWidth:          768,
			DefaultHeight:         512,
			DefaultSampleSteps:    25,
			DefaultCfgScale:       7.0,
		},
		Output: config.OutputConfig{
			OutputDir:     "./my_outputs",
			NamingPattern: "image_{seed}.png",
			AutoTimestamp: false,
		},
		Performance: config.PerformanceConfig{
			NThreads:           4,
			OffloadParamsToCPU: true,
			DiffusionFlashAttn: true,
		},
	}

	// Validate the configuration
	if err := customCfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	fmt.Println("Custom configuration created and validated successfully!")

	// Example 3: Convert config to ContextParams
	fmt.Println("\n3. Converting configuration to ContextParams...")
	ctxParams, err := customCfg.ToContextParams()
	if err != nil {
		log.Fatalf("Failed to convert config: %v", err)
	}

	fmt.Printf("Diffusion Model: %s\n", ctxParams.DiffusionModelPath)
	fmt.Printf("VAE Path: %s\n", ctxParams.VAEPath)
	fmt.Printf("Threads: %d\n", ctxParams.NThreads)
	fmt.Printf("Flash Attention: %v\n", ctxParams.DiffusionFlashAttn)

	// Example 4: Convert config to ImgGenParams
	fmt.Println("\n4. Converting configuration to ImgGenParams...")
	imgParams := customCfg.ToImgGenParams(nil)

	fmt.Printf("Width: %d\n", imgParams.Width)
	fmt.Printf("Height: %d\n", imgParams.Height)
	fmt.Printf("Sample Steps: %d\n", imgParams.SampleSteps)
	fmt.Printf("CFG Scale: %.2f\n", imgParams.CfgScale)

	// Example 5: Override specific parameters
	fmt.Println("\n5. Overriding specific parameters...")
	overrides := &stablediffusion.ImgGenParams{
		Prompt: "a beautiful sunset over mountains",
		Width:  1024,
		Height: 1024,
		Seed:   12345,
	}
	mergedParams := customCfg.ToImgGenParams(overrides)

	fmt.Printf("Merged Prompt: %s\n", mergedParams.Prompt)
	fmt.Printf("Merged Width: %d\n", mergedParams.Width)
	fmt.Printf("Merged Height: %d\n", mergedParams.Height)
	fmt.Printf("Merged Seed: %d\n", mergedParams.Seed)
	// Note: Other values come from config defaults
	fmt.Printf("Merged CFG Scale (from config): %.2f\n", mergedParams.CfgScale)

	// Example 6: Generate output path
	fmt.Println("\n6. Generating output path...")
	outputPath := customCfg.GenerateOutputPath(12345, "png")
	fmt.Printf("Generated Output Path: %s\n", outputPath)

	// Example 7: Save configuration
	fmt.Println("\n7. Saving configuration to file...")
	if err := config.SaveConfig(customCfg, "./custom_config.yaml"); err != nil {
		log.Printf("Failed to save config: %v", err)
	} else {
		fmt.Println("Configuration saved to ./custom_config.yaml")
	}

	// Example 8: Using ConfigManager
	fmt.Println("\n8. Using ConfigManager...")
	cfgManager, err := config.NewConfigManager("./managed_config.yaml", true)
	if err != nil {
		log.Printf("Failed to create config manager: %v", err)
	} else {
		// Update configuration
		err = cfgManager.Update(func(c *config.Config) {
			c.Generation.DefaultWidth = 1024
			c.Generation.DefaultHeight = 1024
		})
		if err != nil {
			log.Printf("Failed to update config: %v", err)
		} else {
			fmt.Println("Configuration updated and auto-saved!")
			fmt.Printf("New default size: %dx%d\n",
				cfgManager.Get().Generation.DefaultWidth,
				cfgManager.Get().Generation.DefaultHeight)
		}
	}

	// Example 9: Effective prompts
	fmt.Println("\n9. Effective prompts...")
	effectivePrompt := cfg.GetEffectivePrompt("sunset colors")
	fmt.Printf("Effective Prompt: %s\n", effectivePrompt)

	effectiveNegPrompt := cfg.GetEffectiveNegativePrompt("watermark")
	fmt.Printf("Effective Negative Prompt: %s\n", effectiveNegPrompt)

	fmt.Println("\n✅ Configuration example completed!")
}
