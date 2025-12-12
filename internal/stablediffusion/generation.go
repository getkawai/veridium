package stablediffusion

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	db "github.com/kawai-network/veridium/internal/database/generated"
)

// RuntimeImageGenParams matches frontend model-bank/standard-parameters RuntimeImageGenParams
type RuntimeImageGenParams struct {
	Prompt      string   `json:"prompt"`
	ImageUrl    *string  `json:"imageUrl,omitempty"`
	ImageUrls   []string `json:"imageUrls,omitempty"`
	Width       *int     `json:"width,omitempty"`
	Height      *int     `json:"height,omitempty"`
	Size        string   `json:"size,omitempty"`
	AspectRatio string   `json:"aspectRatio,omitempty"`
	Cfg         *float64 `json:"cfg,omitempty"`
	Strength    *float64 `json:"strength,omitempty"`
	Steps       *int     `json:"steps,omitempty"`
	Quality     string   `json:"quality,omitempty"`
	Seed        *int64   `json:"seed,omitempty"`
	SamplerName string   `json:"samplerName,omitempty"`
	Scheduler   string   `json:"scheduler,omitempty"`
}

// CreateImageRequest matches frontend createImage action parameters
type CreateImageRequest struct {
	GenerationTopicId string                `json:"generationTopicId"`
	Provider          string                `json:"provider"`
	Model             string                `json:"model"`
	ImageNum          int                   `json:"imageNum"`
	Params            RuntimeImageGenParams `json:"params"`
}

// GenerationOptions defines internal options for SD binary execution
type GenerationOptions struct {
	Prompt         string
	NegativePrompt string
	ModelPath      string
	OutputPath     string
	ImageUrl       *string
	ImageUrls      []string
	Width          int
	Height         int
	Size           string
	AspectRatio    string
	Steps          int
	Cfg            float64
	Strength       float64
	Seed           *int64
	Quality        string
	SamplerName    string
	Scheduler      string
	OutputFormat   string
	Model          string // Model name for remote API
}

// CreateImage handles frontend CreateImageRequest and generates images
func (sdrm *StableDiffusion) CreateImage(req CreateImageRequest) error {
	ctx := context.Background()

	log.Printf("[CreateImage] Starting image generation for topic: %s", req.GenerationTopicId)

	// 1. Get first available model if not specified
	modelPath := sdrm.GetFirstAvailableModel()
	if modelPath == "" {
		log.Printf("[CreateImage] ERROR: No SD model found")
		return fmt.Errorf("no SD model found")
	}
	log.Printf("[CreateImage] Using model: %s", modelPath)

	// 2. Resolve UserID and Topic
	// We need the UserID from the topic to associate the generation
	// Note: We use the Queries from the injected DB service
	log.Printf("[CreateImage] Checking database service...")
	if sdrm.DB == nil {
		log.Printf("[CreateImage] ERROR: Database service not initialized")
		return fmt.Errorf("database service not initialized")
	}

	// Use DefaultUserID since we operate in single-user mode for now
	const DefaultUserID = "DEFAULT_LOBE_CHAT_USER"
	userID := DefaultUserID
	topic, err := sdrm.DB.Queries().GetGenerationTopic(ctx, db.GetGenerationTopicParams{
		ID:     req.GenerationTopicId,
		UserID: userID,
	})
	log.Printf("[CreateImage] Fetching generation topic from DB...")
	if err != nil {
		log.Printf("[CreateImage] ERROR: Failed to find generation topic: %v", err)
		return fmt.Errorf("failed to find generation topic: %w", err)
	}
	log.Printf("[CreateImage] Topic found for user: %s", topic.UserID)

	// 3. Create Generation Batch
	now := time.Now().UnixMilli()
	batchID := uuid.New().String()

	configBytes, _ := json.Marshal(req.Params)

	generationBatch := db.CreateGenerationBatchParams{
		ID:                batchID,
		UserID:            topic.UserID,
		GenerationTopicID: req.GenerationTopicId,
		Provider:          req.Provider,
		Model:             req.Model,
		Prompt:            req.Params.Prompt,
		Config:            sql.NullString{String: string(configBytes), Valid: true},
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// Handle optional dimensions for batch record
	if req.Params.Width != nil {
		generationBatch.Width = sql.NullInt64{Int64: int64(*req.Params.Width), Valid: true}
	}
	if req.Params.Height != nil {
		generationBatch.Height = sql.NullInt64{Int64: int64(*req.Params.Height), Valid: true}
	}
	if req.Params.AspectRatio != "" {
		generationBatch.Ratio = sql.NullString{String: req.Params.AspectRatio, Valid: true}
	}

	log.Printf("[CreateImage] Creating generation batch in DB...")
	if _, err := sdrm.DB.Queries().CreateGenerationBatch(ctx, generationBatch); err != nil {
		log.Printf("[CreateImage] ERROR: Failed to create generation batch: %v", err)
		return fmt.Errorf("failed to create generation batch: %w", err)
	}
	log.Printf("[CreateImage] Batch created with ID: %s", batchID)

	// 4. Prepare Output Directory
	// Use "files/uploads" as requested, ensuring it exists
	outputDir := "files/uploads"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 5. Convert RuntimeImageGenParams to GenerationOptions
	opts := GenerationOptions{
		Prompt:      req.Params.Prompt,
		ModelPath:   modelPath,
		ImageUrl:    req.Params.ImageUrl,
		ImageUrls:   req.Params.ImageUrls,
		Size:        req.Params.Size,
		AspectRatio: req.Params.AspectRatio,
		Quality:     req.Params.Quality,
		SamplerName: req.Params.SamplerName,
		Scheduler:   req.Params.Scheduler,
		Seed:        req.Params.Seed,
	}

	// Handle optional numeric params
	if req.Params.Width != nil {
		opts.Width = *req.Params.Width
	}
	if req.Params.Height != nil {
		opts.Height = *req.Params.Height
	}
	if req.Params.Steps != nil {
		opts.Steps = *req.Params.Steps
	}
	if req.Params.Cfg != nil {
		opts.Cfg = *req.Params.Cfg
	}
	if req.Params.Strength != nil {
		opts.Strength = *req.Params.Strength
	}

	// 6. Generate multiple images in parallel
	log.Printf("[CreateImage] Generating %d image(s) in parallel...", req.ImageNum)

	// Available models for variation
	availableModels := []string{
		"flux",
		"stable-diffusion",
		"kontext",
		"turbo",
		"nanobanana",
		"seedream",
		"nanobanana-pro",
		"seedream-pro",
		"gptimage",
		"zimage",
		"veo",
		"seedance",
		"seedance-pro",
	}

	// Track successful image paths for fallback
	var successfulImagePath string
	var failedIndices []int
	var mu sync.Mutex // Mutex for thread-safe access to shared variables
	var wg sync.WaitGroup

	// Launch parallel generation
	for i := 0; i < req.ImageNum; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			log.Printf("[CreateImage] Generating image %d/%d", index+1, req.ImageNum)
			// Generate ID for the file
			fileName := fmt.Sprintf("gen_%s_%d.png", batchID, index)
			outputPath := filepath.Join(outputDir, fileName)

			// Create a copy of opts for this goroutine
			localOpts := opts
			localOpts.OutputPath = outputPath

			// Rotate through available models for variation
			modelIndex := index % len(availableModels)
			localOpts.Model = availableModels[modelIndex]
			log.Printf("[CreateImage] Using model: %s for image %d", localOpts.Model, index)

			// Try remote generation
			log.Printf("[CreateImage] Attempting remote generation for image %d...", index)
			remoteErr := sdrm.generateImageRemote(localOpts)

			// Thread-safe update of shared variables
			mu.Lock()
			if remoteErr != nil {
				log.Printf("[CreateImage] Remote generation failed for image %d: %v", index, remoteErr)
				failedIndices = append(failedIndices, index)
			} else {
				log.Printf("[CreateImage] Image %d generated successfully (remote)", index)
				// Store the first successful image path for fallback
				if successfulImagePath == "" {
					successfulImagePath = outputPath
				}
			}
			mu.Unlock()
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	log.Printf("[CreateImage] All parallel generations completed")

	// Handle failed images: copy from successful image if available
	if len(failedIndices) > 0 {
		if successfulImagePath != "" {
			log.Printf("[CreateImage] Copying successful image to %d failed slot(s)", len(failedIndices))
			for _, idx := range failedIndices {
				failedFileName := fmt.Sprintf("gen_%s_%d.png", batchID, idx)
				failedOutputPath := filepath.Join(outputDir, failedFileName)

				// Copy successful image to failed slot
				if err := copyFile(successfulImagePath, failedOutputPath); err != nil {
					log.Printf("[CreateImage] ERROR: Failed to copy image to slot %d: %v", idx, err)
					return fmt.Errorf("failed to copy image for slot %d: %w", idx, err)
				}
				log.Printf("[CreateImage] Copied successful image to slot %d", idx)
			}
		} else {
			// All images failed, return error
			return fmt.Errorf("all remote generations failed and no successful image available")
		}
	}

	// 7. Post-Processing: Database Records for all images
	for i := 0; i < req.ImageNum; i++ {
		fileName := fmt.Sprintf("gen_%s_%d.png", batchID, i)
		outputPath := filepath.Join(outputDir, fileName)

		// 7a. Calculate File Info
		fileInfo, err := os.Stat(outputPath)
		if err != nil {
			return fmt.Errorf("failed to stat generated file: %w", err)
		}

		// Calculate SHA256 Hash
		f, err := os.Open(outputPath)
		if err != nil {
			return fmt.Errorf("failed to open generated file for hashing: %w", err)
		}
		hash := sha256.New()
		if _, err := io.Copy(hash, f); err != nil {
			f.Close()
			return fmt.Errorf("failed to hash generated file: %w", err)
		}
		f.Close()
		fileHash := hex.EncodeToString(hash.Sum(nil))

		fileUrl := fmt.Sprintf("/files/uploads/%s", fileName)

		// 7b. Create GlobalFile (if not exists)
		// Try to get it first
		_, err = sdrm.DB.Queries().GetGlobalFile(ctx, fileHash)
		if err != nil && err == sql.ErrNoRows {
			// Create GlobalFile
			_, err = sdrm.DB.Queries().CreateGlobalFile(ctx, db.CreateGlobalFileParams{
				HashID:   fileHash,
				FileType: "image/png",
				Size:     fileInfo.Size(),
				Url:      fileUrl,
				Creator:  topic.UserID, // Original creator
			})
			if err != nil {
				// Ignore error if it's unique constraint (race condition)
				// But log it? For now assume it's fine.
				fmt.Printf("Warning: failed to create global file: %v\n", err)
			}
		}

		// 7c. Create File Record
		// Source: "ImageGeneration" (enum emulation)
		savedFile, err := sdrm.DB.Queries().CreateFile(ctx, db.CreateFileParams{
			UserID:   topic.UserID,
			FileType: "image/png",
			FileHash: sql.NullString{String: fileHash, Valid: true},
			Name:     fileName,
			Size:     fileInfo.Size(),
			Url:      fileUrl,
			Source:   sql.NullString{String: "ImageGeneration", Valid: true},
			Metadata: sql.NullString{String: "{}", Valid: true}, // Empty metadata default
		})
		if err != nil {
			return fmt.Errorf("failed to create file record: %w", err)
		}

		// 7d. Create Generation Record
		genID := uuid.New().String()
		_, err = sdrm.DB.Queries().CreateGeneration(ctx, db.CreateGenerationParams{
			ID:                genID,
			UserID:            topic.UserID,
			GenerationBatchID: batchID,
			FileID:            sql.NullString{String: savedFile.ID, Valid: true},
			Seed:              sql.NullInt64{Int64: 0, Valid: false}, // Seed tracking could be improved
			CreatedAt:         now,
			UpdatedAt:         now,
			Asset: func() sql.NullString {
				assetMap := map[string]interface{}{
					"type":         "image",
					"url":          fileUrl,
					"width":        opts.Width,
					"height":       opts.Height,
					"thumbnailUrl": fileUrl, // Use same URL for now
					"originalUrl":  fileUrl,
				}
				bytes, _ := json.Marshal(assetMap)
				return sql.NullString{String: string(bytes), Valid: true}
			}(),
		})
		if err != nil {
			return fmt.Errorf("failed to create generation record: %w", err)
		}
	}

	return nil
}

// GetFirstAvailableModel returns the first available SD model
func (sdrm *StableDiffusion) GetFirstAvailableModel() string {
	modelsPath := sdrm.GetModelsPath()
	files, err := os.ReadDir(modelsPath)
	if err != nil {
		return ""
	}

	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			// Check for all supported model formats including GGUF
			if strings.HasSuffix(name, ".ckpt") ||
				strings.HasSuffix(name, ".safetensors") ||
				strings.HasSuffix(name, ".pt") ||
				strings.HasSuffix(name, ".bin") ||
				strings.HasSuffix(name, ".gguf") {
				return filepath.Join(modelsPath, name)
			}
		}
	}
	return ""
}

// GetOutputsPath returns the path for generated images
// Note: This matches the default but we override it in CreateImage to use files/uploads
func (sdrm *StableDiffusion) GetOutputsPath() string {
	homeDir, _ := os.UserHomeDir()
	return homeDir + "/.stable-diffusion/outputs"
}

// CreateImageWithOptions generates an image using GenerationOptions directly
// Used by internal services like image_designer
func (sdrm *StableDiffusion) CreateImageWithOptions(opts GenerationOptions) error {
	return sdrm.createImageInternal(opts)
}

// generateImageRemote generates an image using Pollinations AI remote API
func (sdrm *StableDiffusion) generateImageRemote(opts GenerationOptions) error {
	// Build the API URL
	baseURL := "https://image.pollinations.ai/prompt/"
	encodedPrompt := url.PathEscape(opts.Prompt)
	apiURL := baseURL + encodedPrompt

	// Build query parameters
	params := url.Values{}
	if opts.Width > 0 {
		params.Add("width", strconv.Itoa(opts.Width))
	}
	if opts.Height > 0 {
		params.Add("height", strconv.Itoa(opts.Height))
	}
	// Use specified model or default to flux
	model := opts.Model
	if model == "" {
		model = "flux"
	}
	params.Add("model", model)
	params.Add("enhance", "false")
	params.Add("nologo", "true")

	if len(params) > 0 {
		apiURL += "$" + params.Encode()
	}

	log.Printf("[Remote] Requesting image from: %s", apiURL)

	// Make HTTP request
	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("remote API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote API returned status %d", resp.StatusCode)
	}

	// Create output file
	outFile, err := os.Create(opts.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Download image
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}

	log.Printf("[Remote] Image downloaded successfully to: %s", opts.OutputPath)
	return nil
}

// createImageInternal executes the Stable Diffusion binary to generate an image
func (sdrm *StableDiffusion) createImageInternal(opts GenerationOptions) error {
	binaryPath := sdrm.getBinaryPath()

	// Check if binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("stable diffusion binary not found at %s", binaryPath)
	}

	// Default values
	if opts.Width == 0 {
		opts.Width = 1024
	}
	if opts.Height == 0 {
		opts.Height = 1024
	}
	if opts.Steps == 0 {
		opts.Steps = 20
	}
	if opts.Cfg == 0 {
		opts.Cfg = 7.0
	}

	// Prepare arguments
	args := []string{
		"-m", opts.ModelPath,
		"-p", opts.Prompt,
		"-o", opts.OutputPath,
		"--width", strconv.Itoa(opts.Width),
		"--height", strconv.Itoa(opts.Height),
		"--steps", strconv.Itoa(opts.Steps),
		"--cfg-scale", strconv.FormatFloat(opts.Cfg, 'f', -1, 64),
	}

	// Add seed if specified
	if opts.Seed != nil {
		args = append(args, "--seed", strconv.FormatInt(*opts.Seed, 10))
	}

	if opts.NegativePrompt != "" {
		args = append(args, "-n", opts.NegativePrompt)
	}

	// Add sampler if specified
	if opts.SamplerName != "" {
		args = append(args, "--sampling-method", opts.SamplerName)
	}

	// Add scheduler if specified
	if opts.Scheduler != "" {
		args = append(args, "--schedule", opts.Scheduler)
	}

	// Add strength for img2img
	if opts.Strength > 0 {
		args = append(args, "--strength", strconv.FormatFloat(opts.Strength, 'f', -1, 64))
	}

	// Add input image for img2img
	if opts.ImageUrl != nil && *opts.ImageUrl != "" {
		args = append(args, "-i", *opts.ImageUrl)
	}

	// Log the command for debugging
	log.Printf("[SD] Executing: %s %v", binaryPath, args)

	// Execute command via the injected executor
	if err := sdrm.Executor.Run(sdrm.ctx, binaryPath, args...); err != nil {
		return fmt.Errorf("generation failed: %w", err)
	}

	// Verify output exists
	if _, err := os.Stat(opts.OutputPath); err != nil {
		return fmt.Errorf("output file was not created: %w", err)
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
