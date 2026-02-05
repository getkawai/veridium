package remote

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	cfv6 "github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/ai"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/kawai-network/veridium/internal/constant"
)

// CloudflareGenerator implements Generator using Cloudflare Workers AI
type CloudflareGenerator struct {
	accountID string
	client    *cfv6.Client
}

// NewCloudflareGenerator creates a new Cloudflare generator
func NewCloudflareGenerator() *CloudflareGenerator {
	creds := constant.GetRandomCloudflareApiKey()
	parts := strings.Split(creds, ":")

	var accountID, apiKey string
	if len(parts) == 2 {
		accountID = parts[0]
		apiKey = parts[1]
	}

	client := cfv6.NewClient(
		option.WithAPIToken(apiKey),
	)

	return &CloudflareGenerator{
		accountID: accountID,
		client:    client,
	}
}

// Ensure interface compliance
var _ Generator = (*CloudflareGenerator)(nil)

func (g *CloudflareGenerator) IsAvailable() bool {
	return g.client != nil && g.accountID != ""
}

// GetAvailableModels returns list of available models
func (g *CloudflareGenerator) GetAvailableModels() []string {
	return []string{
		"@cf/black-forest-labs/flux-1-schnell",
	}
}

// Generate generates an image using Cloudflare Workers AI
func (g *CloudflareGenerator) Generate(ctx context.Context, opts GenerationOptions) error {
	if !g.IsAvailable() {
		return fmt.Errorf("Cloudflare generator is not available")
	}

	model := opts.Model
	if model == "" {
		model = "@cf/black-forest-labs/flux-1-schnell"
	}

	// Prepare inputs using strong types
	params := ai.AIRunParamsBodyTextToImage{
		Prompt: cfv6.F(opts.Prompt),
	}

	if opts.Width != 0 {
		params.Width = cfv6.F(int64(opts.Width))
	}
	if opts.Height != 0 {
		params.Height = cfv6.F(int64(opts.Height))
	}
	if opts.Steps != 0 {
		params.NumSteps = cfv6.F(int64(opts.Steps))
	}
	if opts.Seed != nil {
		params.Seed = cfv6.F(*opts.Seed)
	}

	// Execute Run using AI service
	// We use the typed params. body is the union.
	// We cast our struct to the union type implicitly or explicitly if needed?
	// The struct implements the interface.

	resp, err := g.client.AI.Run(ctx, model, ai.AIRunParams{
		AccountID: cfv6.F(g.accountID),
		Body:      params,
	})

	if err != nil {
		return fmt.Errorf("failed to run workers AI: %w", err)
	}

	// Response handling
	// Use reflection to inspect or try to cast if it's a known wrapper
	// Since we expect binary, but the SDK returns a Union...
	// If the SDK parses it as JSON, we might get garbage or error if it's really binary.
	// HOWEVER, Cloudflare AI usually returns binary for image models.
	// If the SDK assumes JSON, it might fail to unmarshal.
	// If it fails to unmarshal, 'err' above should have triggered?

	// Let's assume for now that if it succeeds, 'resp' contains something useful.
	// If it's *AIRunResponseUnion (ptr to interface), we dereference.

	var actualResp interface{}
	if resp != nil {
		actualResp = *resp
	}

	switch v := actualResp.(type) {
	case []byte:
		if err := os.WriteFile(opts.OutputPath, v, 0644); err != nil {
			return err
		}
	case io.ReadCloser:
		defer v.Close()
		f, err := os.Create(opts.OutputPath)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(f, v); err != nil {
			return err
		}
	case io.Reader:
		// Stream to file
		f, err := os.Create(opts.OutputPath)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := io.Copy(f, v); err != nil {
			return err
		}
	default:
		// Fallback: If the SDK failed to give us the binary (likely because it tried to parse JSON),
		// we might need to use raw request. But let's see what we got.
		return fmt.Errorf("unexpected response type from Cloudflare AI: %T", v)
	}

	log.Printf("[CloudflareGen] Image saved to %s", opts.OutputPath)
	return nil
}
