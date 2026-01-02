package image

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/internal/topic"
)

// Service handles high-level image generation operations with database persistence
type Service struct {
	*StableDiffusion // Embedded Engine
	DB               *database.Service
	TopicService     *topic.TopicService
}

// NewService creates a new image generation service
func NewService(db *database.Service, engine *StableDiffusion) *Service {
	return &Service{
		StableDiffusion: engine,
		DB:              db,
	}
}

// SetTopicService sets the topic service
func (s *Service) SetTopicService(ts *topic.TopicService) {
	s.TopicService = ts
}

// CreateImage handles frontend CreateImageRequest and generates images asynchronously
func (s *Service) CreateImage(req CreateImageRequest) error {
	ctx := context.Background()

	log.Printf("[CreateImage] Starting image generation for topic: %s", req.GenerationTopicId)

	// 1. Get first available model if not specified
	modelPath := s.GetFirstAvailableModel()
	if modelPath == "" {
		log.Printf("[CreateImage] ERROR: No SD model found")
		return fmt.Errorf("no SD model found")
	}
	log.Printf("[CreateImage] Using model: %s", modelPath)

	// 2. Resolve UserID and Topic
	log.Printf("[CreateImage] Checking database service...")
	if s.DB == nil {
		log.Printf("[CreateImage] ERROR: Database service not initialized")
		return fmt.Errorf("database service not initialized")
	}

	const DefaultUserID = "DEFAULT_LOBE_CHAT_USER"
	topic, err := s.DB.Queries().GetGenerationTopic(ctx, req.GenerationTopicId)
	log.Printf("[CreateImage] Fetching generation topic from DB...")
	if err != nil {
		log.Printf("[CreateImage] ERROR: Failed to find generation topic: %v", err)
		return fmt.Errorf("failed to find generation topic: %w", err)
	}
	log.Printf("[CreateImage] Topic found: %s", topic.ID)

	// 3. Create Generation Batch
	now := time.Now().UnixMilli()
	batchID := uuid.New().String()

	configBytes, _ := json.Marshal(req.Params)

	generationBatch := db.CreateGenerationBatchParams{
		ID:                batchID,
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
	if _, err := s.DB.Queries().CreateGenerationBatch(ctx, generationBatch); err != nil {
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
		_, err := s.DB.Queries().CreateAsyncTask(ctx, db.CreateAsyncTaskParams{
			ID:        taskID,
			Type:      sql.NullString{String: "image_generation", Valid: true},
			Status:    sql.NullString{String: "pending", Valid: true},
			Error:     sql.NullString{},
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
		_, err = s.DB.Queries().CreateGeneration(ctx, db.CreateGenerationParams{
			ID:                genID,
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
	go s.generateImagesInBackground(batchID, req.ImageNum, opts)

	return nil
}

// generateImagesInBackground generates images in parallel and updates database records
func (s *Service) generateImagesInBackground(batchID string, imageNum int, opts GenerationOptions) {
	ctx := context.Background()
	outputDir := filepath.Join(paths.FileBase(), "uploads")

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Printf("[Background] ERROR: Failed to create output directory: %v", err)
		return
	}

	log.Printf("[Background] Starting background generation for batch %s", batchID)

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

	// Get all generations for this batch to update them
	generations, err := s.DB.Queries().ListGenerations(ctx, batchID)
	if err != nil {
		log.Printf("[Background] ERROR: Failed to list generations: %v", err)
		return
	}

	// Trigger topic title update if TopicService is available
	if s.TopicService != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("❌ [PANIC] Topic title update panic recovered: %v", r)
				}
			}()

			batch, err := s.DB.Queries().GetGenerationBatch(ctx, batchID)
			if err != nil {
				log.Printf("[Background] Warning: Failed to fetch batch for title update: %v", err)
				return
			}

			if batch.GenerationTopicID != "" {
				log.Printf("[Background] Triggering topic title update for topic %s", batch.GenerationTopicID)
				// Use Background context with timeout - must outlive HTTP request
				bgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				if err := s.TopicService.UpdateGenerationTopicTitleFromPrompt(bgCtx, batch.GenerationTopicID, opts.Prompt); err != nil {
					log.Printf("[Background] Warning: Failed to update topic title: %v", err)
				}
			}
		}()
	}

	if len(generations) != imageNum {
		log.Printf("[Background] WARNING: Expected %d generations, got %d", imageNum, len(generations))
	}

	// Track successful image paths for fallback
	var successfulImagePath string
	var failedIndices []int
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Generate images in parallel
	for i := 0; i < imageNum; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done() // ✅ MUST be first defer to ensure it always executes
			defer func() {
				if r := recover(); r != nil {
					log.Printf("❌ [PANIC] Image generation panic recovered for index %d: %v", index, r)
					mu.Lock()
					failedIndices = append(failedIndices, index)
					mu.Unlock()
				}
			}()

			if index >= len(generations) {
				log.Printf("[Background] ERROR: No generation record for index %d", index)
				return
			}

			generation := generations[index]
			startTime := time.Now()

			// Update task status to "processing"
			if generation.AsyncTaskID.Valid {
				_, err := s.DB.Queries().UpdateAsyncTask(ctx, db.UpdateAsyncTaskParams{
					ID:        generation.AsyncTaskID.String,
					Status:    sql.NullString{String: "processing", Valid: true},
					Error:     sql.NullString{},
					Duration:  sql.NullInt64{},
					UpdatedAt: time.Now().UnixMilli(),
				})
				if err != nil {
					log.Printf("[Background] ERROR: Failed to update task status to processing: %v", err)
				}
			}

			log.Printf("[Background] Generating image %d/%d", index+1, imageNum)
			fileName := fmt.Sprintf("gen_%s_%d.png", batchID, index)
			outputPath := filepath.Join(outputDir, fileName)

			// Create local copy of opts
			localOpts := opts
			localOpts.OutputPath = outputPath

			// Rotate through models with retry strategy
			startModelIndex := index % len(availableModels)
			var remoteErr error

			// Try all available models until one succeeds
			for attempt := 0; attempt < len(availableModels); attempt++ {
				modelIndex := (startModelIndex + attempt) % len(availableModels)
				localOpts.Model = availableModels[modelIndex]
				log.Printf("[Background] Image %d: Attempt %d/%d using model: %s", index, attempt+1, len(availableModels), localOpts.Model)

				// Try remote generation
				remoteErr = s.generateImageRemote(localOpts)
				if remoteErr == nil {
					log.Printf("[Background] Image %d: Success with model %s", index, localOpts.Model)
					break
				}

				log.Printf("[Background] Image %d: Failed with model %s: %v", index, localOpts.Model, remoteErr)
			}

			mu.Lock()
			if remoteErr != nil {
				log.Printf("[Background] Remote generation failed for image %d: %v", index, remoteErr)
				failedIndices = append(failedIndices, index)
			} else {
				log.Printf("[Background] Image %d generated successfully", index)
				if successfulImagePath == "" {
					successfulImagePath = outputPath
				}
			}
			mu.Unlock()

			// Update database record
			if remoteErr != nil {
				// Update task with error
				if generation.AsyncTaskID.Valid {
					_, err := s.DB.Queries().UpdateAsyncTask(ctx, db.UpdateAsyncTaskParams{
						ID:        generation.AsyncTaskID.String,
						Status:    sql.NullString{String: "error", Valid: true},
						Error:     sql.NullString{String: remoteErr.Error(), Valid: true},
						Duration:  sql.NullInt64{Int64: time.Since(startTime).Milliseconds(), Valid: true},
						UpdatedAt: time.Now().UnixMilli(),
					})
					if err != nil {
						log.Printf("[Background] ERROR: Failed to update task with error: %v", err)
					}
				}
			} else {
				// Success - update with file data
				if err := s.updateGenerationWithFile(ctx, generation.ID, outputPath, fileName, opts, generation.AsyncTaskID.String, startTime); err != nil {
					log.Printf("[Background] ERROR: Failed to update generation with file: %v", err)
					remoteErr = err
				}
			}
		}(i)
	}

	// Wait for all generations to complete
	wg.Wait()
	log.Printf("[Background] All parallel generations completed for batch %s", batchID)

	// Handle failed images by copying successful one
	if len(failedIndices) > 0 && successfulImagePath != "" {
		log.Printf("[Background] Copying successful image to %d failed slot(s)", len(failedIndices))
		for _, idx := range failedIndices {
			if idx >= len(generations) {
				continue
			}

			generation := generations[idx]
			failedFileName := fmt.Sprintf("gen_%s_%d.png", batchID, idx)
			failedOutputPath := filepath.Join(outputDir, failedFileName)

			// Copy successful image
			if err := copyFile(successfulImagePath, failedOutputPath); err != nil {
				log.Printf("[Background] ERROR: Failed to copy image to slot %d: %v", idx, err)
				continue
			}

			// Update generation with copied file
			s.updateGenerationWithFile(ctx, generation.ID, failedOutputPath, failedFileName, opts, generation.AsyncTaskID.String, time.Now())
		}
	}
}

// updateGenerationWithFile updates a generation record with file data after successful generation
func (s *Service) updateGenerationWithFile(ctx context.Context, generationID, outputPath, fileName string, opts GenerationOptions, taskID string, startTime time.Time) error {
	// Calculate file info
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		log.Printf("[Background] ERROR: Failed to stat file: %v", err)
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// Calculate hash
	f, err := os.Open(outputPath)
	if err != nil {
		log.Printf("[Background] ERROR: Failed to open file for hashing: %v", err)
		return fmt.Errorf("failed to open file for hashing: %w", err)
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		f.Close()
		log.Printf("[Background] ERROR: Failed to hash file: %v", err)
		return fmt.Errorf("failed to hash file: %w", err)
	}
	f.Close()
	fileHash := hex.EncodeToString(hash.Sum(nil))
	fileUrl := "/files/uploads/" + fileName

	// Create GlobalFile if not exists
	_, err = s.DB.Queries().GetGlobalFile(ctx, fileHash)
	if err != nil && err == sql.ErrNoRows {
		_, err = s.DB.Queries().CreateGlobalFile(ctx, db.CreateGlobalFileParams{
			HashID:   fileHash,
			FileType: "image/png",
			Size:     fileInfo.Size(),
			Url:      fileUrl,
			Creator:  sql.NullString{String: "default", Valid: true},
		})
		if err != nil {
			log.Printf("[Background] Warning: failed to create global file: %v", err)
		}
	}

	// Create File record
	savedFile, err := s.DB.Queries().CreateFile(ctx, db.CreateFileParams{
		// UserID removed
		FileType: "image/png",
		FileHash: sql.NullString{String: fileHash, Valid: true},
		Name:     fileName,
		Size:     fileInfo.Size(),
		Url:      fileUrl,
		Source:   sql.NullString{String: "ImageGeneration", Valid: true},
		Metadata: sql.NullString{String: "{}", Valid: true},
	})
	if err != nil {
		log.Printf("[Background] ERROR: Failed to create file record: %v", err)
		return fmt.Errorf("failed to create file record: %w", err)
	}

	// Create asset JSON
	assetMap := map[string]interface{}{
		"type":         "image",
		"url":          fileUrl,
		"width":        opts.Width,
		"height":       opts.Height,
		"thumbnailUrl": fileUrl,
		"originalUrl":  fileUrl,
	}
	assetBytes, _ := json.Marshal(assetMap)

	// Update generation record
	_, err = s.DB.Queries().UpdateGeneration(ctx, db.UpdateGenerationParams{
		ID:          generationID,
		AsyncTaskID: sql.NullString{String: taskID, Valid: true},
		FileID:      sql.NullString{String: savedFile.ID, Valid: true},
		Asset:       sql.NullString{String: string(assetBytes), Valid: true},
		UpdatedAt:   time.Now().UnixMilli(),
	})
	if err != nil {
		log.Printf("[Background] ERROR: Failed to update generation: %v", err)
		return fmt.Errorf("failed to update generation: %w", err)
	}

	// Update async_task to success
	_, err = s.DB.Queries().UpdateAsyncTask(ctx, db.UpdateAsyncTaskParams{
		ID:        taskID,
		Status:    sql.NullString{String: "success", Valid: true},
		Error:     sql.NullString{},
		Duration:  sql.NullInt64{Int64: time.Since(startTime).Milliseconds(), Valid: true},
		UpdatedAt: time.Now().UnixMilli(),
	})
	if err != nil {
		log.Printf("[Background] ERROR: Failed to update task to success: %v", err)
	}

	log.Printf("[Background] Successfully updated generation %s with file %s", generationID, fileName)
	return nil
}
