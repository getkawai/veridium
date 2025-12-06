package builtin

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/internal/stablediffusion"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
	"github.com/kawai-network/veridium/types"
)

// ============================================================================
// Response Types (matching frontend expected format)
// ============================================================================

// DallEImageItem matches frontend DallEImageItem interface
type DallEImageItem struct {
	Prompt     string `json:"prompt"`
	PreviewUrl string `json:"previewUrl,omitempty"`
	ImageId    string `json:"imageId,omitempty"`
	Quality    string `json:"quality"` // "standard" | "hd"
	Size       string `json:"size"`    // "1792x1024" | "1024x1024" | "1024x1792"
	Style      string `json:"style"`   // "vivid" | "natural"
}

// ImageDesignerService provides image generation capabilities using Stable Diffusion
type ImageDesignerService struct {
	sdManager   *stablediffusion.StableDiffusionReleaseManager
	outputDir   string
	initialized bool
}

// NewImageDesignerService creates a new image designer service
func NewImageDesignerService() *ImageDesignerService {
	homeDir, _ := os.UserHomeDir()
	outputDir := filepath.Join(homeDir, ".stable-diffusion", "outputs")
	os.MkdirAll(outputDir, 0755)

	return &ImageDesignerService{
		sdManager:   stablediffusion.NewStableDiffusionReleaseManager(),
		outputDir:   outputDir,
		initialized: false,
	}
}

// IsAvailable checks if Stable Diffusion is installed and ready
func (s *ImageDesignerService) IsAvailable() bool {
	if !s.sdManager.IsStableDiffusionInstalled() {
		return false
	}

	// Check if at least one model is installed
	models, err := s.sdManager.CheckInstalledModels()
	if err != nil || len(models) == 0 {
		return false
	}

	return true
}

// GetFirstAvailableModel returns the first available SD model
func (s *ImageDesignerService) GetFirstAvailableModel() string {
	modelsPath := s.sdManager.GetModelsPath()
	files, err := os.ReadDir(modelsPath)
	if err != nil {
		return ""
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if strings.HasSuffix(name, ".ckpt") ||
				strings.HasSuffix(name, ".safetensors") ||
				strings.HasSuffix(name, ".pt") {
				return filepath.Join(modelsPath, name)
			}
		}
	}
	return ""
}

// Text2Image generates images from text prompts using Stable Diffusion
func (s *ImageDesignerService) Text2Image(prompts []string, quality, size, style string, seeds []int) ([]DallEImageItem, error) {
	// Default values
	if quality == "" {
		quality = "standard"
	}
	if size == "" {
		size = "1024x1024"
	}
	if style == "" {
		style = "vivid"
	}

	// Check if SD is available
	if !s.IsAvailable() {
		log.Printf("⚠️  Stable Diffusion not available, using placeholder images")
		return s.generatePlaceholders(prompts, quality, size, style, seeds)
	}

	// Get SD binary and model paths
	sdBinary := s.sdManager.GetBinaryPath()
	modelPath := s.GetFirstAvailableModel()

	if modelPath == "" {
		log.Printf("⚠️  No SD model found, using placeholder images")
		return s.generatePlaceholders(prompts, quality, size, style, seeds)
	}

	log.Printf("🎨 Using Stable Diffusion: %s", filepath.Base(modelPath))

	results := make([]DallEImageItem, 0, len(prompts))

	for i, prompt := range prompts {
		// Generate unique output filename
		imageId := uuid.New().String()
		outputPath := filepath.Join(s.outputDir, fmt.Sprintf("%s.png", imageId))

		// Determine seed
		seed := time.Now().UnixNano() + int64(i)
		if len(seeds) > i && seeds[i] > 0 {
			seed = int64(seeds[i])
		}

		// Parse size
		width, height := parseDallESize(size)

		// Determine steps based on quality
		steps := 20
		if quality == "hd" {
			steps = 30
		}

		// Build SD command
		// sd -m model.safetensors -p "prompt" -o output.png --width 1024 --height 1024 --steps 20 --seed 12345
		args := []string{
			"-m", modelPath,
			"-p", prompt,
			"-o", outputPath,
			"--width", strconv.Itoa(width),
			"--height", strconv.Itoa(height),
			"--steps", strconv.Itoa(steps),
			"--seed", strconv.FormatInt(seed, 10),
		}

		// Add negative prompt for better quality
		negativePrompt := "ugly, blurry, low quality, distorted"
		if style == "natural" {
			negativePrompt = "cartoon, anime, illustration, " + negativePrompt
		}
		args = append(args, "-n", negativePrompt)

		log.Printf("🖼️  Generating image %d/%d: %s", i+1, len(prompts), truncateString(prompt, 50))

		// Execute SD
		cmd := exec.Command(sdBinary, args...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			log.Printf("⚠️  SD generation failed: %v, stderr: %s", err, stderr.String())
			// Fallback to placeholder
			result := s.generateSinglePlaceholder(prompt, quality, size, style, i)
			results = append(results, result)
			continue
		}

		// Check if output was created
		if _, err := os.Stat(outputPath); err != nil {
			log.Printf("⚠️  Output image not found: %s", outputPath)
			result := s.generateSinglePlaceholder(prompt, quality, size, style, i)
			results = append(results, result)
			continue
		}

		// Read image and convert to data URL for preview
		imageData, err := os.ReadFile(outputPath)
		if err != nil {
			log.Printf("⚠️  Failed to read output image: %v", err)
			result := s.generateSinglePlaceholder(prompt, quality, size, style, i)
			results = append(results, result)
			continue
		}

		// Create data URL
		previewUrl := fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(imageData))

		results = append(results, DallEImageItem{
			Prompt:     prompt,
			PreviewUrl: previewUrl,
			ImageId:    imageId,
			Quality:    quality,
			Size:       size,
			Style:      style,
		})

		log.Printf("✅ Generated image: %s", imageId)
	}

	return results, nil
}

// generatePlaceholders generates placeholder images when SD is not available
func (s *ImageDesignerService) generatePlaceholders(prompts []string, quality, size, style string, seeds []int) ([]DallEImageItem, error) {
	results := make([]DallEImageItem, 0, len(prompts))

	for i, prompt := range prompts {
		result := s.generateSinglePlaceholder(prompt, quality, size, style, i)
		if len(seeds) > i {
			// Use seed in placeholder URL for consistency
			width, height := parseDallESize(size)
			result.PreviewUrl = fmt.Sprintf("https://picsum.photos/seed/%d/%d/%d", seeds[i], width, height)
		}
		results = append(results, result)
	}

	return results, nil
}

// generateSinglePlaceholder generates a single placeholder image
func (s *ImageDesignerService) generateSinglePlaceholder(prompt, quality, size, style string, index int) DallEImageItem {
	width, height := parseDallESize(size)
	placeholderUrl := fmt.Sprintf("https://picsum.photos/seed/%d/%d/%d", index, width, height)

	log.Printf("🎨 Generated placeholder for: %s", truncateString(prompt, 50))

	return DallEImageItem{
		Prompt:     prompt,
		PreviewUrl: placeholderUrl,
		Quality:    quality,
		Size:       size,
		Style:      style,
	}
}

// parseDallESize parses DALL-E size string to width and height
func parseDallESize(size string) (int, int) {
	switch size {
	case "1792x1024":
		return 1792, 1024
	case "1024x1792":
		return 1024, 1792
	default:
		return 1024, 1024
	}
}

// truncateString truncates a string to max length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// ============================================================================
// Tool Registration
// ============================================================================

// RegisterImageDesigner registers the lobe-image-designer tool (DALL-E compatible)
func RegisterImageDesigner(registry *tools.ToolRegistry) error {
	service := NewImageDesignerService()

	tool := &types.Tool{
		Type: fantasy.ToolTypeFunction,
		Definition: types.ToolDefinition{
			Name:        "lobe-image-designer__text2image",
			Description: "Create images from text prompts using AI image generation. Generate up to 4 diverse images based on the description.",
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"prompts": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"minItems":    1,
						"maxItems":    4,
						"description": "Array of detailed image descriptions. Create diverse prompts if user doesn't specify exact number.",
					},
					"quality": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"standard", "hd"},
						"default":     "standard",
						"description": "Image quality. 'hd' creates images with finer details.",
					},
					"size": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"1792x1024", "1024x1024", "1024x1792"},
						"default":     "1024x1024",
						"description": "Image resolution. Use 1024x1024 (square) as default, 1792x1024 for wide, 1024x1792 for tall/portrait.",
					},
					"style": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"vivid", "natural"},
						"default":     "vivid",
						"description": "Image style. 'vivid' for hyper-real/dramatic, 'natural' for more realistic.",
					},
					"seeds": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "integer"},
						"description": "Optional seeds for reproducible generation when modifying previous images.",
					},
				},
				"required": []string{"prompts"},
			},
		},
		Executor: func(ctx context.Context, args map[string]string) (string, error) {
			// Parse prompts
			promptsStr := args["prompts"]
			if promptsStr == "" {
				return "", fmt.Errorf("prompts is required")
			}

			var prompts []string
			if err := json.Unmarshal([]byte(promptsStr), &prompts); err != nil {
				return "", fmt.Errorf("failed to parse prompts: %w", err)
			}

			if len(prompts) == 0 {
				return "", fmt.Errorf("at least one prompt is required")
			}
			if len(prompts) > 4 {
				prompts = prompts[:4] // Limit to 4
			}

			// Parse optional parameters
			quality := args["quality"]
			size := args["size"]
			style := args["style"]

			var seeds []int
			if seedsStr := args["seeds"]; seedsStr != "" {
				json.Unmarshal([]byte(seedsStr), &seeds)
			}

			// Generate images
			results, err := service.Text2Image(prompts, quality, size, style, seeds)
			if err != nil {
				return "", err
			}

			resultJSON, err := json.Marshal(results)
			if err != nil {
				return "", fmt.Errorf("failed to marshal results: %w", err)
			}

			log.Printf("🖼️  Generated %d images", len(results))
			return string(resultJSON), nil
		},
		Enabled: true,
	}

	return registry.Register(tool)
}
