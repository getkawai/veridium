package image

import (
	"context"
	"database/sql"
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

// CreateImage handles frontend CreateImageRequest and generates images asynchronously
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
	log.Printf("[CreateImage] Checking database service...")
	if sdrm.DB == nil {
		log.Printf("[CreateImage] ERROR: Database service not initialized")
		return fmt.Errorf("database service not initialized")
	}

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

	// 4. Create placeholder async_task and generation records
	log.Printf("[CreateImage] Creating placeholder records for %d images...", req.ImageNum)

	// Convert params to GenerationOptions for background worker
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

	// Create placeholder records for each image
	for i := 0; i < req.ImageNum; i++ {
		// Create async_task with pending status
		taskID := uuid.New().String()
		_, err := sdrm.DB.Queries().CreateAsyncTask(ctx, db.CreateAsyncTaskParams{
			ID:        taskID,
			Type:      sql.NullString{String: "image_generation", Valid: true},
			Status:    sql.NullString{String: "pending", Valid: true},
			Error:     sql.NullString{},
			UserID:    topic.UserID,
			Duration:  sql.NullInt64{},
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			log.Printf("[CreateImage] ERROR: Failed to create async task: %v", err)
			return fmt.Errorf("failed to create async task: %w", err)
		}

		// Create generation record with pending task
		genID := uuid.New().String()
		_, err = sdrm.DB.Queries().CreateGeneration(ctx, db.CreateGenerationParams{
			ID:                genID,
			UserID:            topic.UserID,
			GenerationBatchID: batchID,
			AsyncTaskID:       sql.NullString{String: taskID, Valid: true},
			FileID:            sql.NullString{}, // Will be filled when image completes
			Seed:              sql.NullInt64{Int64: 0, Valid: false},
			Asset:             sql.NullString{}, // Will be filled when image completes
			CreatedAt:         now,
			UpdatedAt:         now,
		})
		if err != nil {
			log.Printf("[CreateImage] ERROR: Failed to create generation record: %v", err)
			return fmt.Errorf("failed to create generation record: %w", err)
		}

		log.Printf("[CreateImage] Created placeholder for image %d/%d (task: %s, gen: %s)", i+1, req.ImageNum, taskID, genID)
	}

	log.Printf("[CreateImage] All placeholder records created, returning to frontend")

	// 5. Launch background goroutine to generate images
	go sdrm.generateImagesInBackground(batchID, topic.UserID, req.ImageNum, opts)

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
