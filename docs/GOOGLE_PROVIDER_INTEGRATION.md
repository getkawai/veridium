# Integrasi Google Provider ke Veridium

Dokumen ini menjelaskan cara mengintegrasikan Google Provider (Gemini & Vertex AI) ke dalam sistem Veridium.

## Overview

Google Provider telah diimplementasikan di `pkg/fantasy/providers/google/` dan siap untuk diintegrasikan ke dalam aplikasi Veridium. Provider ini mendukung:

- Gemini API (dengan API key)
- Vertex AI (dengan Google Cloud credentials)
- Text generation dan streaming
- Tool calling
- Object generation
- Reasoning/thinking mode
- Multi-modal (text + image)

## Struktur File

```
pkg/fantasy/providers/google/
├── google.go              # Implementasi utama provider
├── auth.go                # Authentication helpers
├── error.go               # Error handling
├── provider_options.go    # Provider-specific options
├── slice.go               # Utility functions
├── google_test.go         # Unit tests
└── README.md              # Dokumentasi provider
```

## Cara Integrasi

### 1. Import Provider di Context

Edit file `internal/app/context.go` untuk menambahkan import:

```go
import (
    // ... existing imports ...
    googleprovider "github.com/kawai-network/veridium/pkg/fantasy/providers/google"
)
```

### 2. Tambahkan Konfigurasi Environment

Tambahkan ke `.env`:

```bash
# Google Gemini API
GEMINI_API_KEY=your_gemini_api_key_here

# Atau untuk Vertex AI
GOOGLE_CLOUD_PROJECT=your-project-id
GOOGLE_CLOUD_LOCATION=us-central1
```

### 3. Inisialisasi Provider

Ada beberapa cara untuk mengintegrasikan provider:

#### Option A: Sebagai Provider Standalone

```go
func initGoogleProvider(ctx context.Context) (fantasy.LanguageModel, error) {
    apiKey := os.Getenv("GEMINI_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("GEMINI_API_KEY not set")
    }

    provider, err := googleprovider.New(
        googleprovider.WithGeminiAPIKey(apiKey),
    )
    if err != nil {
        return nil, err
    }

    return provider.LanguageModel(ctx, "gemini-1.5-flash")
}
```

#### Option B: Dalam Pooled Provider Chain

Tambahkan ke fungsi `initChatModel` atau `initSummaryModel` di `context.go`:

```go
// Tambahkan Google Gemini ke chain
if apiKey := os.Getenv("GEMINI_API_KEY"); apiKey != "" {
    log.Printf("🔍 %s: Initializing Google Gemini...", taskName)
    if provider, err := googleprovider.New(
        googleprovider.WithGeminiAPIKey(apiKey),
    ); err == nil {
        if geminiModel, err := provider.LanguageModel(bgCtx, "gemini-1.5-flash"); err == nil {
            chain = append(chain, geminiModel)
            log.Printf("%s: Google Gemini (gemini-1.5-flash)", taskName)
        }
    }
}
```

#### Option C: Dalam Pooled Provider dengan API Key Pool

```go
pooledProviders = append(pooledProviders, pooled.ProviderConfig{
    Name:   "google-gemini",
    Weight: 1.0,
    APIKeys: []string{
        constant.GetRandomGeminiApiKey(), // Implementasi di constant package
    },
    CreateClient: func(apiKey string) (fantasy.LanguageModel, error) {
        provider, err := googleprovider.New(
            googleprovider.WithGeminiAPIKey(apiKey),
        )
        if err != nil {
            return nil, err
        }
        return provider.LanguageModel(bgCtx, "gemini-1.5-flash")
    },
})
```

### 4. Tambahkan Helper di Constant Package

Edit `internal/constant/api_keys.go`:

```go
var geminiApiKeys = []string{
    // Tambahkan API keys dari environment atau config
}

func GetRandomGeminiApiKey() string {
    if len(geminiApiKeys) == 0 {
        if key := os.Getenv("GEMINI_API_KEY"); key != "" {
            return key
        }
        return ""
    }
    return geminiApiKeys[rand.Intn(len(geminiApiKeys))]
}

func init() {
    // Load dari environment
    if keys := os.Getenv("GEMINI_API_KEYS"); keys != "" {
        geminiApiKeys = strings.Split(keys, ",")
    }
}
```

## Contoh Penggunaan dalam Veridium

### 1. Chat Service

```go
// Di internal/services/chat_service.go
func (s *ChatService) GenerateResponse(ctx context.Context, message string) (string, error) {
    response, err := s.model.Generate(ctx, fantasy.Call{
        Prompt: fantasy.Prompt{
            {
                Role: fantasy.MessageRoleUser,
                Content: []fantasy.MessagePart{
                    fantasy.TextPart{Text: message},
                },
            },
        },
        Temperature: fantasy.Opt(0.7),
    })
    if err != nil {
        return "", err
    }
    return response.Content.Text(), nil
}
```

### 2. RAG dengan Google

```go
// Di internal/services/rag_processor.go
func (r *RAGProcessor) ProcessWithGoogle(ctx context.Context, query string, docs []string) (string, error) {
    // Gabungkan dokumen sebagai context
    context := strings.Join(docs, "\n\n")
    
    response, err := r.googleModel.Generate(ctx, fantasy.Call{
        Prompt: fantasy.Prompt{
            {
                Role: fantasy.MessageRoleSystem,
                Content: []fantasy.MessagePart{
                    fantasy.TextPart{
                        Text: fmt.Sprintf("Use the following context to answer questions:\n\n%s", context),
                    },
                },
            },
            {
                Role: fantasy.MessageRoleUser,
                Content: []fantasy.MessagePart{
                    fantasy.TextPart{Text: query},
                },
            },
        },
    })
    if err != nil {
        return "", err
    }
    return response.Content.Text(), nil
}
```

### 3. Structured Output untuk Knowledge Base

```go
// Di internal/services/knowledge_base.go
type ExtractedInfo struct {
    Title   string   `json:"title"`
    Summary string   `json:"summary"`
    Tags    []string `json:"tags"`
}

func (kb *KnowledgeBaseService) ExtractInfo(ctx context.Context, content string) (*ExtractedInfo, error) {
    objResponse, err := kb.googleModel.GenerateObject(ctx, fantasy.ObjectCall{
        Prompt: fantasy.Prompt{
            {
                Role: fantasy.MessageRoleUser,
                Content: []fantasy.MessagePart{
                    fantasy.TextPart{
                        Text: fmt.Sprintf("Extract title, summary, and tags from:\n\n%s", content),
                    },
                },
            },
        },
        Schema: schema.FromStruct(ExtractedInfo{}),
    })
    if err != nil {
        return nil, err
    }
    return objResponse.Object.(*ExtractedInfo), nil
}
```

## Testing

### Unit Tests

```bash
# Test provider
go test ./pkg/fantasy/providers/google/...

# Test dengan coverage
go test -cover ./pkg/fantasy/providers/google/...
```

### Integration Tests

```bash
# Set API key
export GEMINI_API_KEY=your_key_here

# Run example
go run examples/google_provider_example.go
```

## Model Selection Guide

### Untuk Chat (Real-time)
- `gemini-1.5-flash` - Cepat, murah, cocok untuk chat
- `gemini-1.5-flash-8b` - Lebih cepat, lebih murah

### Untuk Summary/Analysis
- `gemini-1.5-pro` - Lebih akurat, context window besar
- `gemini-2.0-flash-exp` - Experimental, fitur terbaru

### Untuk Reasoning Tasks
- `gemini-2.0-flash-thinking-exp` - Dengan thinking mode
- Gunakan dengan `ThinkingConfig`

## Cost Optimization

1. **Gunakan Flash untuk Chat**: `gemini-1.5-flash` lebih murah
2. **Context Caching**: Untuk dokumen besar yang sering digunakan
3. **Pooled Provider**: Rotasi API keys untuk rate limiting
4. **Fallback Chain**: Gemini → Pollinations → Local

## Monitoring

Tambahkan logging untuk tracking:

```go
log.Printf("Google Provider: Model=%s, Tokens=%d, Cost=$%.4f",
    response.Model,
    response.Usage.TotalTokens,
    calculateCost(response.Usage),
)
```

## Troubleshooting

### Error: "API key not valid"
- Periksa `GEMINI_API_KEY` di `.env`
- Pastikan API key aktif di Google AI Studio

### Error: "Quota exceeded"
- Gunakan pooled provider dengan multiple keys
- Tambahkan fallback provider

### Error: "Model not found"
- Periksa nama model (case-sensitive)
- Beberapa model hanya tersedia di Vertex AI

## Next Steps

1. ✅ Provider sudah diimplementasikan
2. ⏳ Tambahkan ke `context.go` initialization
3. ⏳ Tambahkan API key management di `constant` package
4. ⏳ Integrasikan ke chat service
5. ⏳ Tambahkan monitoring dan logging
6. ⏳ Setup context caching untuk RAG

## Resources

- [Google AI Studio](https://aistudio.google.com/)
- [Gemini API Docs](https://ai.google.dev/docs)
- [Vertex AI Docs](https://cloud.google.com/vertex-ai/docs)
- [Provider README](../pkg/fantasy/providers/google/README.md)
