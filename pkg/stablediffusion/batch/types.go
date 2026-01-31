// Package batch provides batch image generation capabilities for stable-diffusion-go
package batch

import (
	"fmt"
	"sync"
	"time"

	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
	"github.com/kawai-network/veridium/pkg/stablediffusion/progress"
)

// VariationParams represents parameters for a single variation in a batch
type VariationParams struct {
	Prompt           string
	NegativePrompt   string
	Seed             int64
	CfgScale         float32
	ImageCfgScale    float32
	SampleSteps      int32
	Width            int32
	Height           int32
	SampleMethod     string
	Scheduler        string
	InitImagePath    string
	MaskImagePath    string
	ControlImagePath string
	Strength         float32
}

// BatchImgGenParams contains parameters for batch image generation
type BatchImgGenParams struct {
	BaseParams      stablediffusion.ImgGenParams
	Variations      []VariationParams
	OutputPattern   string // Pattern for output filenames, e.g., "output_%03d.png"
	Parallelism     int    // Number of parallel workers (default: 1)
	ContinueOnError bool   // Whether to continue on individual errors
}

// BatchResult contains the results of a batch generation
type BatchResult struct {
	Images       []string
	Errors       []error
	StartTime    time.Time
	EndTime      time.Time
	Duration     time.Duration
	SuccessCount int
	FailCount    int
	mu           sync.RWMutex
}

// AddSuccess adds a successful result
func (r *BatchResult) AddSuccess(path string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Images = append(r.Images, path)
	r.SuccessCount++
}

// AddError adds an error result
func (r *BatchResult) AddError(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Errors = append(r.Errors, err)
	r.FailCount++
}

// GetImages returns a copy of the generated image paths
func (r *BatchResult) GetImages() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]string, len(r.Images))
	copy(result, r.Images)
	return result
}

// GetErrors returns a copy of the errors
func (r *BatchResult) GetErrors() []error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]error, len(r.Errors))
	copy(result, r.Errors)
	return result
}

// IsSuccessful returns true if all generations were successful
func (r *BatchResult) IsSuccessful() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.FailCount == 0
}

// BatchProgress represents progress for batch generation
type BatchProgress struct {
	Completed     int
	Total         int
	CurrentItem   int
	CurrentStage  progress.GenerationStage
	Percentage    float64
	TimeElapsed   time.Duration
	TimeRemaining time.Duration
	Message       string
}

// BatchCallback is called during batch generation
type BatchCallback func(progress BatchProgress)

// Generator handles batch image generation
type Generator struct {
	sd *stablediffusion.StableDiffusion
}

// NewGenerator creates a new batch generator
func NewGenerator(sd *stablediffusion.StableDiffusion) *Generator {
	return &Generator{sd: sd}
}

// Generate performs batch image generation
func (g *Generator) Generate(params *BatchImgGenParams) (*BatchResult, error) {
	return g.GenerateWithCallback(params, nil)
}

// GenerateWithCallback performs batch generation with progress callback
func (g *Generator) GenerateWithCallback(
	params *BatchImgGenParams,
	callback BatchCallback,
) (*BatchResult, error) {
	result := &BatchResult{
		Images:    make([]string, 0),
		Errors:    make([]error, 0),
		StartTime: time.Now(),
	}

	// Determine number of variations
	total := len(params.Variations)
	if total == 0 {
		return nil, fmt.Errorf("no variations specified")
	}

	// Set default parallelism
	parallelism := params.Parallelism
	if parallelism <= 0 {
		parallelism = 1
	}
	if parallelism > total {
		parallelism = total
	}

	// Create work queue
	type workItem struct {
		index     int
		varParams VariationParams
	}

	workQueue := make(chan workItem, total)
	for i, v := range params.Variations {
		workQueue <- workItem{index: i, varParams: v}
	}
	close(workQueue)

	// Create result channel
	type resultItem struct {
		index int
		path  string
		err   error
	}
	resultChan := make(chan resultItem, total)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < parallelism; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for work := range workQueue {
				// Merge base params with variation
				imgParams := g.mergeParams(&params.BaseParams, &work.varParams)

				// Generate output path
				outputPath := g.generateOutputPath(params.OutputPattern, work.index, work.varParams.Seed)

				// Generate image
				err := g.sd.GenerateImage(imgParams, outputPath)

				resultChan <- resultItem{
					index: work.index,
					path:  outputPath,
					err:   err,
				}
			}
		}(i)
	}

	// Wait for completion in a goroutine
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	completed := 0
	startTime := time.Now()

	for res := range resultChan {
		completed++

		if res.err != nil {
			result.AddError(fmt.Errorf("item %d: %w", res.index, res.err))
			if !params.ContinueOnError {
				// Signal remaining workers to stop (in a real implementation)
				break
			}
		} else {
			result.AddSuccess(res.path)
		}

		// Report progress
		if callback != nil {
			elapsed := time.Since(startTime)
			percentage := float64(completed) / float64(total) * 100.0

			// Estimate remaining time
			var remaining time.Duration
			if completed > 0 {
				avgTime := elapsed / time.Duration(completed)
				remaining = avgTime * time.Duration(total-completed)
			}

			callback(BatchProgress{
				Completed:     completed,
				Total:         total,
				CurrentItem:   res.index,
				Percentage:    percentage,
				TimeElapsed:   elapsed,
				TimeRemaining: remaining,
				Message:       fmt.Sprintf("Generated %d/%d images", completed, total),
			})
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

// GenerateSequential performs sequential batch generation (no parallelism)
func (g *Generator) GenerateSequential(
	params *BatchImgGenParams,
	callback BatchCallback,
) (*BatchResult, error) {
	result := &BatchResult{
		Images:    make([]string, 0),
		Errors:    make([]error, 0),
		StartTime: time.Now(),
	}

	total := len(params.Variations)
	if total == 0 {
		return nil, fmt.Errorf("no variations specified")
	}

	for i, varParams := range params.Variations {
		// Merge base params with variation
		imgParams := g.mergeParams(&params.BaseParams, &varParams)

		// Generate output path
		outputPath := g.generateOutputPath(params.OutputPattern, i, varParams.Seed)

		// Report progress before generation
		if callback != nil {
			callback(BatchProgress{
				Completed:   i,
				Total:       total,
				CurrentItem: i,
				Percentage:  float64(i) / float64(total) * 100.0,
				TimeElapsed: time.Since(result.StartTime),
				Message:     fmt.Sprintf("Generating image %d/%d", i+1, total),
			})
		}

		// Generate image
		err := g.sd.GenerateImage(imgParams, outputPath)

		if err != nil {
			result.AddError(fmt.Errorf("item %d: %w", i, err))
			if !params.ContinueOnError {
				break
			}
		} else {
			result.AddSuccess(outputPath)
		}

		// Report progress after generation
		if callback != nil {
			elapsed := time.Since(result.StartTime)
			completed := i + 1

			var remaining time.Duration
			if completed > 0 {
				avgTime := elapsed / time.Duration(completed)
				remaining = avgTime * time.Duration(total-completed)
			}

			callback(BatchProgress{
				Completed:     completed,
				Total:         total,
				CurrentItem:   i,
				Percentage:    float64(completed) / float64(total) * 100.0,
				TimeElapsed:   elapsed,
				TimeRemaining: remaining,
				Message:       fmt.Sprintf("Generated %d/%d images", completed, total),
			})
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

// mergeParams merges base params with variation params
func (g *Generator) mergeParams(base *stablediffusion.ImgGenParams, variation *VariationParams) *stablediffusion.ImgGenParams {
	result := *base

	if variation.Prompt != "" {
		result.Prompt = variation.Prompt
	}
	if variation.NegativePrompt != "" {
		result.NegativePrompt = variation.NegativePrompt
	}
	if variation.Seed != 0 {
		result.Seed = variation.Seed
	}
	if variation.CfgScale > 0 {
		result.CfgScale = variation.CfgScale
	}
	if variation.ImageCfgScale > 0 {
		result.ImageCfgScale = variation.ImageCfgScale
	}
	if variation.SampleSteps > 0 {
		result.SampleSteps = variation.SampleSteps
	}
	if variation.Width > 0 {
		result.Width = variation.Width
	}
	if variation.Height > 0 {
		result.Height = variation.Height
	}
	if variation.SampleMethod != "" {
		result.SampleMethod = variation.SampleMethod
	}
	if variation.Scheduler != "" {
		result.Scheduler = variation.Scheduler
	}
	if variation.InitImagePath != "" {
		result.InitImagePath = variation.InitImagePath
	}
	if variation.MaskImagePath != "" {
		result.MaskImagePath = variation.MaskImagePath
	}
	if variation.ControlImagePath != "" {
		result.ControlImagePath = variation.ControlImagePath
	}
	if variation.Strength > 0 {
		result.Strength = variation.Strength
	}

	return &result
}

// generateOutputPath generates an output path based on pattern
func (g *Generator) generateOutputPath(pattern string, index int, seed int64) string {
	if pattern == "" {
		pattern = "batch_%03d.png"
	}

	// Replace %d or %03d with the index
	path := fmt.Sprintf(pattern, index)

	// Replace {seed} with the seed value
	path = fmt.Sprintf(path, seed)

	return path
}
