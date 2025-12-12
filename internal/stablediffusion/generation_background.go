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
	"os"
	"path/filepath"
	"sync"
	"time"

	db "github.com/kawai-network/veridium/internal/database/generated"
)

// generateImagesInBackground generates images in parallel and updates database records
func (sdrm *StableDiffusion) generateImagesInBackground(batchID, userID string, imageNum int, opts GenerationOptions) {
	ctx := context.Background()
	outputDir := "files/uploads"

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
	generations, err := sdrm.DB.Queries().ListGenerations(ctx, db.ListGenerationsParams{
		GenerationBatchID: batchID,
		UserID:            userID,
	})
	if err != nil {
		log.Printf("[Background] ERROR: Failed to list generations: %v", err)
		return
	}

	// Trigger topic title update if TopicService is available
	if sdrm.TopicService != nil {
		go func() {
			batch, err := sdrm.DB.Queries().GetGenerationBatch(ctx, db.GetGenerationBatchParams{
				ID:     batchID,
				UserID: userID,
			})
			if err != nil {
				log.Printf("[Background] Warning: Failed to fetch batch for title update: %v", err)
				return
			}

			if batch.GenerationTopicID != "" {
				log.Printf("[Background] Triggering topic title update for topic %s", batch.GenerationTopicID)
				if err := sdrm.TopicService.UpdateGenerationTopicTitleFromPrompt(context.Background(), batch.GenerationTopicID, userID, opts.Prompt); err != nil {
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
			defer wg.Done()

			if index >= len(generations) {
				log.Printf("[Background] ERROR: No generation record for index %d", index)
				return
			}

			generation := generations[index]
			startTime := time.Now()

			// Update task status to "processing"
			if generation.AsyncTaskID.Valid {
				_, err := sdrm.DB.Queries().UpdateAsyncTask(ctx, db.UpdateAsyncTaskParams{
					ID:        generation.AsyncTaskID.String,
					UserID:    userID,
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
				remoteErr = sdrm.generateImageRemote(localOpts)
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
					_, err := sdrm.DB.Queries().UpdateAsyncTask(ctx, db.UpdateAsyncTaskParams{
						ID:        generation.AsyncTaskID.String,
						UserID:    userID,
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
				sdrm.updateGenerationWithFile(ctx, generation.ID, userID, outputPath, fileName, opts, generation.AsyncTaskID.String, startTime)
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
			sdrm.updateGenerationWithFile(ctx, generation.ID, userID, failedOutputPath, failedFileName, opts, generation.AsyncTaskID.String, time.Now())
		}
	}
}

// updateGenerationWithFile updates a generation record with file data after successful generation
func (sdrm *StableDiffusion) updateGenerationWithFile(ctx context.Context, generationID, userID, outputPath, fileName string, opts GenerationOptions, taskID string, startTime time.Time) {
	// Calculate file info
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		log.Printf("[Background] ERROR: Failed to stat file: %v", err)
		return
	}

	// Calculate hash
	f, err := os.Open(outputPath)
	if err != nil {
		log.Printf("[Background] ERROR: Failed to open file for hashing: %v", err)
		return
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		f.Close()
		log.Printf("[Background] ERROR: Failed to hash file: %v", err)
		return
	}
	f.Close()
	fileHash := hex.EncodeToString(hash.Sum(nil))
	fileUrl := fmt.Sprintf("/files/uploads/%s", fileName)

	// Create GlobalFile if not exists
	_, err = sdrm.DB.Queries().GetGlobalFile(ctx, fileHash)
	if err != nil && err == sql.ErrNoRows {
		_, err = sdrm.DB.Queries().CreateGlobalFile(ctx, db.CreateGlobalFileParams{
			HashID:   fileHash,
			FileType: "image/png",
			Size:     fileInfo.Size(),
			Url:      fileUrl,
			Creator:  userID,
		})
		if err != nil {
			log.Printf("[Background] Warning: failed to create global file: %v", err)
		}
	}

	// Create File record
	savedFile, err := sdrm.DB.Queries().CreateFile(ctx, db.CreateFileParams{
		UserID:   userID,
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
		return
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
	_, err = sdrm.DB.Queries().UpdateGeneration(ctx, db.UpdateGenerationParams{
		ID:          generationID,
		UserID:      userID,
		AsyncTaskID: sql.NullString{String: taskID, Valid: true},
		FileID:      sql.NullString{String: savedFile.ID, Valid: true},
		Asset:       sql.NullString{String: string(assetBytes), Valid: true},
		UpdatedAt:   time.Now().UnixMilli(),
	})
	if err != nil {
		log.Printf("[Background] ERROR: Failed to update generation: %v", err)
		return
	}

	// Update async_task to success
	_, err = sdrm.DB.Queries().UpdateAsyncTask(ctx, db.UpdateAsyncTaskParams{
		ID:        taskID,
		UserID:    userID,
		Status:    sql.NullString{String: "success", Valid: true},
		Error:     sql.NullString{},
		Duration:  sql.NullInt64{Int64: time.Since(startTime).Milliseconds(), Valid: true},
		UpdatedAt: time.Now().UnixMilli(),
	})
	if err != nil {
		log.Printf("[Background] ERROR: Failed to update task to success: %v", err)
	}

	log.Printf("[Background] Successfully updated generation %s with file %s", generationID, fileName)
}
