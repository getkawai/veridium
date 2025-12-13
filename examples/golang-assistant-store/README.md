# Go Assistant Store

Implementasi Go untuk mengambil data assistants dari NPM registry `@lobehub/agents-index`.

## 🚀 Features

- ✅ Fetch agent index dengan multi-locale support
- ✅ Fetch detail agent individual
- ✅ In-memory caching dengan expiration
- ✅ Fallback ke default locale jika locale tidak tersedia
- ✅ Search agents by query
- ✅ Filter by category
- ✅ Whitelist/Blacklist filtering
- ✅ Get categories dengan count

## 📦 Installation

```bash
cd examples/golang-assistant-store
go mod download
```

## 🎯 Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    // Initialize store dengan default URL
    store := NewAssistantStore("")
    
    // Atau dengan custom URL
    // store := NewAssistantStore("https://your-custom-cdn.com/agents")
    
    // Fetch all agents untuk locale tertentu
    agents, err := store.GetAgentIndex("en-US")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d agents\n", len(agents))
    
    // Fetch detail untuk agent tertentu
    detail, err := store.GetAgent("web-development", "en-US")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Agent: %s\n", detail.Meta.Title)
    fmt.Printf("System Role: %s\n", detail.SystemRole)
}
```

### Search Agents

```go
// Search agents by query
results, err := store.SearchAgents("en-US", "web development")
if err != nil {
    log.Fatal(err)
}

for _, agent := range results {
    fmt.Printf("- %s: %s\n", agent.Meta.Title, agent.Meta.Description)
}
```

### Filter by Category

```go
// Get agents in specific category
agents, err := store.FilterByCategory("en-US", "development")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d development agents\n", len(agents))
```

### Get Categories

```go
// Get all categories with counts
categories, err := store.GetCategories("en-US")
if err != nil {
    log.Fatal(err)
}

for category, count := range categories {
    fmt.Printf("%s: %d agents\n", category, count)
}
```

### Whitelist/Blacklist Filtering

```go
// Filter dengan whitelist
filter := &FilterOptions{
    Whitelist: []string{"web-development", "api-design"},
}

agents, err := store.GetAgentIndexWithFilter("en-US", filter)
if err != nil {
    log.Fatal(err)
}

// Filter dengan blacklist
filter = &FilterOptions{
    Blacklist: []string{"spam-agent", "deprecated-agent"},
}

agents, err = store.GetAgentIndexWithFilter("en-US", filter)
```

## 🏃 Run Example

```bash
go run .
```

Output:
```
Found 150 agents

1. Web Development Assistant (web-development)
   Category: development
   Author: lobehub
   Tags: [web html css javascript]

2. API Design Expert (api-design)
   Category: development
   Author: lobehub
   Tags: [api rest graphql]

...

Fetching detail for: web-development
System Role: You are an expert web developer with deep knowledge of HTML, CSS, JavaScript...
```

## 🔧 Configuration

### Custom Base URL

```go
// Gunakan custom CDN atau self-hosted registry
store := NewAssistantStore("https://your-cdn.com/agents")
```

### Environment Variable

```go
import "os"

baseURL := os.Getenv("AGENTS_INDEX_URL")
if baseURL == "" {
    baseURL = DefaultAgentsIndexURL
}

store := NewAssistantStore(baseURL)
```

## 📊 Data Structure

### AgentIndexItem

```go
type AgentIndexItem struct {
    Identifier string    `json:"identifier"`
    Category   string    `json:"category"`
    Author     string    `json:"author"`
    Meta       AgentMeta `json:"meta"`
    CreatedAt  string    `json:"createdAt,omitempty"`
    Homepage   string    `json:"homepage,omitempty"`
}
```

### AgentDetail

```go
type AgentDetail struct {
    Identifier string                 `json:"identifier"`
    Author     string                 `json:"author"`
    SystemRole string                 `json:"systemRole"`
    Meta       AgentMeta              `json:"meta"`
    Config     map[string]interface{} `json:"config,omitempty"`
    Plugins    []string               `json:"plugins,omitempty"`
}
```

### AgentMeta

```go
type AgentMeta struct {
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Tags        []string `json:"tags"`
    Avatar      string   `json:"avatar"`
}
```

## 🎨 Architecture

```
┌─────────────────────────────────────────────┐
│           AssistantStore                    │
├─────────────────────────────────────────────┤
│ - baseURL: string                           │
│ - httpClient: *http.Client                  │
│ - cache: *Cache                             │
├─────────────────────────────────────────────┤
│ + GetAgentIndex(locale) []AgentIndexItem    │
│ + GetAgent(id, locale) *AgentDetail         │
│ + SearchAgents(locale, query) []Agent       │
│ + FilterByCategory(locale, cat) []Agent     │
│ + GetCategories(locale) map[string]int      │
└─────────────────────────────────────────────┘
                    │
                    ├─── Cache (in-memory)
                    │    - Auto expiration
                    │    - Thread-safe
                    │
                    └─── HTTP Client
                         - Timeout: 30s
                         - Fallback to default locale
```

## 🔄 Caching Strategy

- **Index List:** Cache selama 1 jam (3600 detik)
- **Agent Detail:** Cache selama 24 jam (86400 detik)
- **Auto Cleanup:** Expired items dihapus setiap 5 menit
- **Thread-Safe:** Menggunakan `sync.RWMutex`

## 🌐 Supported Locales

- `en-US` (English)
- `zh-CN` (Chinese Simplified)
- `zh-TW` (Chinese Traditional)
- `ja-JP` (Japanese)
- `ko-KR` (Korean)
- `fr-FR` (French)
- `de-DE` (German)
- `es-ES` (Spanish)
- `pt-BR` (Portuguese)
- `ru-RU` (Russian)
- `ar` (Arabic)
- Dan lainnya...

## 🔗 URLs

- **Default NPM Registry:** `https://registry.npmmirror.com/@lobehub/agents-index/v1/files/public`
- **GitHub Repository:** https://github.com/lobehub/lobe-chat-agents
- **NPM Package:** [@lobehub/agents-index](https://www.npmjs.com/package/@lobehub/agents-index)

## 📝 Notes

- Data diambil dari NPM registry mirror, bukan langsung dari GitHub
- Fallback otomatis ke `en-US` jika locale tidak tersedia
- Cache membantu mengurangi network requests
- Thread-safe untuk concurrent access

## 🤝 Contributing

Silakan submit PR atau issue di repository utama LobeChat!

