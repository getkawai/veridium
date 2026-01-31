package batch

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kawai-network/veridium/pkg/stablediffusion"
)

// Builder provides a fluent API for building batch generation jobs
type Builder struct {
	params BatchImgGenParams
}

// NewBuilder creates a new batch builder
func NewBuilder() *Builder {
	return &Builder{
		params: BatchImgGenParams{
			BaseParams:      stablediffusion.ImgGenParams{},
			Variations:      []VariationParams{},
			OutputPattern:   "batch_%03d.png",
			Parallelism:     1,
			ContinueOnError: true,
		},
	}
}

// WithBaseParams sets the base parameters for all variations
func (b *Builder) WithBaseParams(params stablediffusion.ImgGenParams) *Builder {
	b.params.BaseParams = params
	return b
}

// WithOutputPattern sets the output filename pattern
func (b *Builder) WithOutputPattern(pattern string) *Builder {
	b.params.OutputPattern = pattern
	return b
}

// WithParallelism sets the number of parallel workers
func (b *Builder) WithParallelism(n int) *Builder {
	b.params.Parallelism = n
	return b
}

// WithContinueOnError sets whether to continue on errors
func (b *Builder) WithContinueOnError(continueOnError bool) *Builder {
	b.params.ContinueOnError = continueOnError
	return b
}

// AddVariation adds a single variation
func (b *Builder) AddVariation(v VariationParams) *Builder {
	b.params.Variations = append(b.params.Variations, v)
	return b
}

// AddPromptVariation adds variations with different prompts but same other params
func (b *Builder) AddPromptVariation(prompts ...string) *Builder {
	for _, prompt := range prompts {
		b.params.Variations = append(b.params.Variations, VariationParams{
			Prompt: prompt,
		})
	}
	return b
}

// AddSeedVariations adds variations with different seeds
func (b *Builder) AddSeedVariations(basePrompt string, count int) *Builder {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < count; i++ {
		b.params.Variations = append(b.params.Variations, VariationParams{
			Prompt: basePrompt,
			Seed:   rng.Int63(),
		})
	}
	return b
}

// AddSeedRange adds variations with seeds in a specific range
func (b *Builder) AddSeedRange(basePrompt string, startSeed int64, count int) *Builder {
	for i := 0; i < count; i++ {
		b.params.Variations = append(b.params.Variations, VariationParams{
			Prompt: basePrompt,
			Seed:   startSeed + int64(i),
		})
	}
	return b
}

// AddResolutionVariations adds variations with different resolutions
func (b *Builder) AddResolutionVariations(prompt string, resolutions [][2]int32) *Builder {
	for _, res := range resolutions {
		b.params.Variations = append(b.params.Variations, VariationParams{
			Prompt: prompt,
			Width:  res[0],
			Height: res[1],
		})
	}
	return b
}

// AddCfgScaleVariations adds variations with different CFG scales
func (b *Builder) AddCfgScaleVariations(prompt string, scales []float32) *Builder {
	for _, scale := range scales {
		b.params.Variations = append(b.params.Variations, VariationParams{
			Prompt:   prompt,
			CfgScale: scale,
		})
	}
	return b
}

// AddStepVariations adds variations with different step counts
func (b *Builder) AddStepVariations(prompt string, steps []int32) *Builder {
	for _, step := range steps {
		b.params.Variations = append(b.params.Variations, VariationParams{
			Prompt:      prompt,
			SampleSteps: step,
		})
	}
	return b
}

// AddSamplerVariations adds variations with different samplers
func (b *Builder) AddSamplerVariations(prompt string, samplers []string) *Builder {
	for _, sampler := range samplers {
		b.params.Variations = append(b.params.Variations, VariationParams{
			Prompt:       prompt,
			SampleMethod: sampler,
		})
	}
	return b
}

// AddGridSearch performs grid search over multiple parameters
func (b *Builder) AddGridSearch(
	basePrompt string,
	seeds []int64,
	cfgScales []float32,
	steps []int32,
) *Builder {
	for _, seed := range seeds {
		for _, cfg := range cfgScales {
			for _, step := range steps {
				b.params.Variations = append(b.params.Variations, VariationParams{
					Prompt:      basePrompt,
					Seed:        seed,
					CfgScale:    cfg,
					SampleSteps: step,
				})
			}
		}
	}
	return b
}

// AddImg2ImgVariations adds variations for img2img with different strengths
func (b *Builder) AddImg2ImgVariations(
	prompt string,
	initImagePath string,
	strengths []float32,
) *Builder {
	for _, strength := range strengths {
		b.params.Variations = append(b.params.Variations, VariationParams{
			Prompt:        prompt,
			InitImagePath: initImagePath,
			Strength:      strength,
		})
	}
	return b
}

// Count returns the number of variations
func (b *Builder) Count() int {
	return len(b.params.Variations)
}

// Build returns the configured BatchImgGenParams
func (b *Builder) Build() (*BatchImgGenParams, error) {
	if len(b.params.Variations) == 0 {
		return nil, fmt.Errorf("no variations added to batch")
	}
	return &b.params, nil
}

// MustBuild returns the configured BatchImgGenParams or panics
func (b *Builder) MustBuild() *BatchImgGenParams {
	params, err := b.Build()
	if err != nil {
		panic(err)
	}
	return params
}

// PresetBuilder provides common batch presets
type PresetBuilder struct {
	builder *Builder
}

// NewPresetBuilder creates a new preset builder
func NewPresetBuilder() *PresetBuilder {
	return &PresetBuilder{
		builder: NewBuilder(),
	}
}

// XYPreset creates an XY plot (grid of CFG vs Steps)
func (p *PresetBuilder) XYPreset(
	prompt string,
	cfgScales []float32,
	steps []int32,
) *Builder {
	return p.builder.AddGridSearch(prompt, []int64{42}, cfgScales, steps)
}

// SeedExplorer creates variations exploring different seeds
func (p *PresetBuilder) SeedExplorer(prompt string, count int) *Builder {
	return p.builder.AddSeedVariations(prompt, count).
		WithOutputPattern("seed_explore_%03d_seed{seed}.png")
}

// ResolutionTest creates variations at different resolutions
func (p *PresetBuilder) ResolutionTest(prompt string) *Builder {
	resolutions := [][2]int32{
		{512, 512},
		{768, 512},
		{512, 768},
		{1024, 512},
		{512, 1024},
		{1024, 1024},
	}
	return p.builder.AddResolutionVariations(prompt, resolutions).
		WithOutputPattern("res_test_%03d_{width}x{height}.png")
}

// SamplerComparison creates variations with different samplers
func (p *PresetBuilder) SamplerComparison(prompt string) *Builder {
	samplers := []string{
		"euler",
		"euler_a",
		"heun",
		"dpm2",
		"dpm++2m",
		"dpm++2s_a",
	}
	return p.builder.AddSamplerVariations(prompt, samplers).
		WithOutputPattern("sampler_compare_%s.png")
}

// QuickTest creates a small batch for quick testing
func (p *PresetBuilder) QuickTest(prompt string) *Builder {
	return p.builder.AddSeedVariations(prompt, 4).
		WithParallelism(2).
		WithOutputPattern("quick_test_%03d.png")
}
