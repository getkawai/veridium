// Package pipeline provides pipeline/chaining operations for stable-diffusion-go
package pipeline

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
	"github.com/kawai-network/veridium/pkg/stablediffusion/progress"
)

// Step represents a single step in a pipeline
type Step interface {
	Name() string
	Execute(ctx *Context) error
	Validate() error
}

// Context holds the state during pipeline execution
type Context struct {
	InputPath        string
	OutputPath       string
	WorkingDir       string
	Metadata         map[string]interface{}
	CancelToken      *progress.CancelToken
	ProgressTracker  *progress.ProgressTracker
	SD               *stablediffusion.StableDiffusion
	Upscaler         *stablediffusion.Upscaler
	KeepIntermediate bool
}

// NewContext creates a new pipeline context
func NewContext(workingDir string) *Context {
	return &Context{
		WorkingDir:       workingDir,
		Metadata:         make(map[string]interface{}),
		KeepIntermediate: false,
	}
}

// SetInput sets the input path
func (c *Context) SetInput(path string) {
	c.InputPath = path
}

// SetOutput sets the output path
func (c *Context) SetOutput(path string) {
	c.OutputPath = path
}

// GetIntermediatePath generates a path for intermediate files
func (c *Context) GetIntermediatePath(name string) string {
	return filepath.Join(c.WorkingDir, "intermediate", name)
}

// Pipeline represents a chain of operations
type Pipeline struct {
	steps        []Step
	config       *Config
	errorHandler func(error) error
}

// Config contains pipeline configuration
type Config struct {
	StopOnError      bool
	KeepIntermediate bool
	TempDir          string
	Timeout          time.Duration
}

// DefaultConfig returns default pipeline configuration
func DefaultConfig() *Config {
	return &Config{
		StopOnError:      true,
		KeepIntermediate: false,
		TempDir:          os.TempDir(),
		Timeout:          30 * time.Minute,
	}
}

// New creates a new pipeline with default config
func New() *Pipeline {
	return &Pipeline{
		steps:  make([]Step, 0),
		config: DefaultConfig(),
	}
}

// NewWithConfig creates a new pipeline with custom config
func NewWithConfig(config *Config) *Pipeline {
	return &Pipeline{
		steps:  make([]Step, 0),
		config: config,
	}
}

// Add adds a step to the pipeline
func (p *Pipeline) Add(step Step) *Pipeline {
	p.steps = append(p.steps, step)
	return p
}

// AddAll adds multiple steps to the pipeline
func (p *Pipeline) AddAll(steps ...Step) *Pipeline {
	p.steps = append(p.steps, steps...)
	return p
}

// OnError sets an error handler
func (p *Pipeline) OnError(handler func(error) error) *Pipeline {
	p.errorHandler = handler
	return p
}

// Validate validates all steps in the pipeline
func (p *Pipeline) Validate() error {
	for i, step := range p.steps {
		if err := step.Validate(); err != nil {
			return fmt.Errorf("step %d (%s) validation failed: %w", i, step.Name(), err)
		}
	}
	return nil
}

// Execute runs the pipeline
func (p *Pipeline) Execute(ctx *Context) error {
	if err := p.Validate(); err != nil {
		return err
	}

	// Create working directory if needed
	if ctx.WorkingDir == "" {
		ctx.WorkingDir = p.config.TempDir
	}
	if err := os.MkdirAll(ctx.WorkingDir, 0755); err != nil {
		return fmt.Errorf("failed to create working directory: %w", err)
	}

	// Create intermediate directory
	intermediateDir := filepath.Join(ctx.WorkingDir, "intermediate")
	if err := os.MkdirAll(intermediateDir, 0755); err != nil {
		return fmt.Errorf("failed to create intermediate directory: %w", err)
	}

	// Execute steps
	intermediateFiles := make([]string, 0)

	for i, step := range p.steps {
		// Check cancellation
		if ctx.CancelToken != nil && ctx.CancelToken.IsCancelled() {
			return fmt.Errorf("pipeline cancelled at step %d", i)
		}

		// Execute step
		if err := step.Execute(ctx); err != nil {
			handledErr := err
			if p.errorHandler != nil {
				handledErr = p.errorHandler(err)
			}

			if handledErr != nil && p.config.StopOnError {
				// Cleanup intermediate files
				if !p.config.KeepIntermediate {
					for _, f := range intermediateFiles {
						_ = os.Remove(f)
					}
				}
				return fmt.Errorf("step %d (%s) failed: %w", i, step.Name(), handledErr)
			}
		}

		// Track intermediate files
		if ctx.InputPath != "" && ctx.InputPath != intermediateFiles[len(intermediateFiles)-1] {
			intermediateFiles = append(intermediateFiles, ctx.InputPath)
		}
	}

	// Cleanup intermediate files
	if !p.config.KeepIntermediate {
		for _, f := range intermediateFiles {
			if f != ctx.OutputPath { // Don't delete final output
				_ = os.Remove(f)
			}
		}
		_ = os.Remove(intermediateDir)
	}

	return nil
}

// ExecuteWithContext runs the pipeline with a Go context
func (p *Pipeline) ExecuteWithContext(goCtx context.Context, pipelineCtx *Context) error {
	// Create a cancel token that respects the Go context
	cancelToken := progress.NewCancelToken()
	pipelineCtx.CancelToken = cancelToken

	// Watch for context cancellation
	done := make(chan struct{})
	go func() {
		select {
		case <-goCtx.Done():
			cancelToken.Cancel()
		case <-done:
		}
	}()
	defer close(done)

	return p.Execute(pipelineCtx)
}

// StepCount returns the number of steps in the pipeline
func (p *Pipeline) StepCount() int {
	return len(p.steps)
}

// PipelineResult contains the result of pipeline execution
type PipelineResult struct {
	Success        bool
	OutputPath     string
	StepsCompleted int
	Duration       time.Duration
	Errors         []error
}

// Builder provides a fluent API for building pipelines
type Builder struct {
	pipeline *Pipeline
}

// NewBuilder creates a new pipeline builder
func NewBuilder() *Builder {
	return &Builder{
		pipeline: New(),
	}
}

// WithConfig sets the pipeline configuration
func (b *Builder) WithConfig(config *Config) *Builder {
	b.pipeline.config = config
	return b
}

// Then adds a step to the pipeline
func (b *Builder) Then(step Step) *Builder {
	b.pipeline.Add(step)
	return b
}

// StopOnError sets whether to stop on error
func (b *Builder) StopOnError(stop bool) *Builder {
	b.pipeline.config.StopOnError = stop
	return b
}

// KeepIntermediate sets whether to keep intermediate files
func (b *Builder) KeepIntermediate(keep bool) *Builder {
	b.pipeline.config.KeepIntermediate = keep
	return b
}

// WithTimeout sets the pipeline timeout
func (b *Builder) WithTimeout(timeout time.Duration) *Builder {
	b.pipeline.config.Timeout = timeout
	return b
}

// Build returns the configured pipeline
func (b *Builder) Build() *Pipeline {
	return b.pipeline
}

// Common step names
const (
	StepTextToImage  = "TextToImage"
	StepImageToImage = "ImageToImage"
	StepUpscale      = "Upscale"
	StepPreprocess   = "Preprocess"
	StepConvert      = "Convert"
	StepSave         = "Save"
	StepLoad         = "Load"
)
