# Quick Start Guide

## 🚀 5-Minute Quick Start

### Step 1: Clone and Navigate

```bash
cd examples/golang-assistant-store
```

### Step 2: Run the Example

```bash
go run .
```

**Expected Output:**
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
```

## 📖 Common Use Cases

### Use Case 1: List All Agents

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    store := NewAssistantStore("")
    
    agents, err := store.GetAgentIndex("en-US")
    if err != nil {
        log.Fatal(err)
    }
    
    for _, agent := range agents {
        fmt.Printf("- %s: %s\n", agent.Meta.Title, agent.Meta.Description)
    }
}
```

### Use Case 2: Search for Specific Agents

```go
func main() {
    store := NewAssistantStore("")
    
    // Search for web development related agents
    results, err := store.SearchAgents("en-US", "web development")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d agents matching 'web development'\n", len(results))
    for _, agent := range results {
        fmt.Printf("- %s\n", agent.Meta.Title)
    }
}
```

### Use Case 3: Get Agent Details

```go
func main() {
    store := NewAssistantStore("")
    
    // Get detailed information about a specific agent
    detail, err := store.GetAgent("web-development", "en-US")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Title: %s\n", detail.Meta.Title)
    fmt.Printf("Description: %s\n", detail.Meta.Description)
    fmt.Printf("System Role:\n%s\n", detail.SystemRole)
    fmt.Printf("Tags: %v\n", detail.Meta.Tags)
    fmt.Printf("Plugins: %v\n", detail.Plugins)
}
```

### Use Case 4: Browse by Category

```go
func main() {
    store := NewAssistantStore("")
    
    // Get all categories
    categories, err := store.GetCategories("en-US")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Available Categories:")
    for category, count := range categories {
        fmt.Printf("- %s: %d agents\n", category, count)
    }
    
    // Get agents in specific category
    devAgents, err := store.FilterByCategory("en-US", "development")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("\nDevelopment Agents (%d):\n", len(devAgents))
    for _, agent := range devAgents {
        fmt.Printf("- %s\n", agent.Meta.Title)
    }
}
```

### Use Case 5: Multi-Language Support

```go
func main() {
    store := NewAssistantStore("")
    
    // Fetch same agent in different languages
    locales := []string{"en-US", "zh-CN", "ja-JP"}
    
    for _, locale := range locales {
        detail, err := store.GetAgent("web-development", locale)
        if err != nil {
            log.Printf("Failed to fetch %s: %v", locale, err)
            continue
        }
        
        fmt.Printf("\n[%s] %s\n", locale, detail.Meta.Title)
        fmt.Printf("Description: %s\n", detail.Meta.Description)
    }
}
```

### Use Case 6: Custom Filtering

```go
func main() {
    store := NewAssistantStore("")
    
    // Only show specific agents (whitelist)
    filter := &FilterOptions{
        Whitelist: []string{
            "web-development",
            "api-design",
            "code-review",
            "database-expert",
        },
    }
    
    agents, err := store.GetAgentIndexWithFilter("en-US", filter)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Showing %d whitelisted agents:\n", len(agents))
    for _, agent := range agents {
        fmt.Printf("- %s\n", agent.Meta.Title)
    }
}
```

### Use Case 7: Building a CLI Tool

```go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"
)

func main() {
    // Define CLI flags
    locale := flag.String("locale", "en-US", "Locale for agents")
    search := flag.String("search", "", "Search query")
    category := flag.String("category", "", "Filter by category")
    detail := flag.String("detail", "", "Get agent detail by identifier")
    
    flag.Parse()
    
    store := NewAssistantStore("")
    
    // Handle detail request
    if *detail != "" {
        showAgentDetail(store, *detail, *locale)
        return
    }
    
    // Get agent list
    var agents []AgentIndexItem
    var err error
    
    if *search != "" {
        agents, err = store.SearchAgents(*locale, *search)
    } else if *category != "" {
        agents, err = store.FilterByCategory(*locale, *category)
    } else {
        agents, err = store.GetAgentIndex(*locale)
    }
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Display results
    fmt.Printf("Found %d agents\n\n", len(agents))
    for i, agent := range agents {
        fmt.Printf("%d. %s (%s)\n", i+1, agent.Meta.Title, agent.Identifier)
        fmt.Printf("   Category: %s | Author: %s\n", agent.Category, agent.Author)
        fmt.Printf("   %s\n\n", agent.Meta.Description)
    }
}

func showAgentDetail(store *AssistantStore, identifier, locale string) {
    detail, err := store.GetAgent(identifier, locale)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("=== %s ===\n\n", detail.Meta.Title)
    fmt.Printf("Identifier: %s\n", detail.Identifier)
    fmt.Printf("Author: %s\n", detail.Author)
    fmt.Printf("Description: %s\n\n", detail.Meta.Description)
    fmt.Printf("Tags: %v\n", detail.Meta.Tags)
    fmt.Printf("Plugins: %v\n\n", detail.Plugins)
    fmt.Printf("System Role:\n%s\n", detail.SystemRole)
}
```

**Usage:**
```bash
# List all agents
go run . --locale en-US

# Search agents
go run . --search "web development"

# Filter by category
go run . --category development

# Get agent detail
go run . --detail web-development

# Chinese locale
go run . --locale zh-CN --search "网页开发"
```

### Use Case 8: Web API Server

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
)

var store *AssistantStore

func main() {
    store = NewAssistantStore("")
    
    http.HandleFunc("/api/agents", handleGetAgents)
    http.HandleFunc("/api/agents/search", handleSearchAgents)
    http.HandleFunc("/api/agents/detail", handleGetAgentDetail)
    http.HandleFunc("/api/categories", handleGetCategories)
    
    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGetAgents(w http.ResponseWriter, r *http.Request) {
    locale := r.URL.Query().Get("locale")
    if locale == "" {
        locale = "en-US"
    }
    
    agents, err := store.GetAgentIndex(locale)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(agents)
}

func handleSearchAgents(w http.ResponseWriter, r *http.Request) {
    locale := r.URL.Query().Get("locale")
    query := r.URL.Query().Get("q")
    
    if locale == "" {
        locale = "en-US"
    }
    
    agents, err := store.SearchAgents(locale, query)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(agents)
}

func handleGetAgentDetail(w http.ResponseWriter, r *http.Request) {
    locale := r.URL.Query().Get("locale")
    identifier := r.URL.Query().Get("id")
    
    if locale == "" {
        locale = "en-US"
    }
    
    if identifier == "" {
        http.Error(w, "identifier required", http.StatusBadRequest)
        return
    }
    
    detail, err := store.GetAgent(identifier, locale)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(detail)
}

func handleGetCategories(w http.ResponseWriter, r *http.Request) {
    locale := r.URL.Query().Get("locale")
    if locale == "" {
        locale = "en-US"
    }
    
    categories, err := store.GetCategories(locale)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(categories)
}
```

**API Endpoints:**
```bash
# Get all agents
curl http://localhost:8080/api/agents?locale=en-US

# Search agents
curl http://localhost:8080/api/agents/search?locale=en-US&q=web

# Get agent detail
curl http://localhost:8080/api/agents/detail?locale=en-US&id=web-development

# Get categories
curl http://localhost:8080/api/categories?locale=en-US
```

## 🛠️ Development Commands

```bash
# Run the application
make run

# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make fmt

# Run linter
make vet

# Build binary
make build

# Run all checks
make all
```

## 🐛 Troubleshooting

### Issue: "Failed to fetch agent index"

**Solution:**
```bash
# Check network connectivity
curl https://registry.npmmirror.com/@lobehub/agents-index/v1/files/public/index.en-US.json

# Try with custom URL
export AGENTS_INDEX_URL="https://cdn.jsdelivr.net/npm/@lobehub/agents-index/public"
go run .
```

### Issue: "Cache not working"

**Solution:**
```go
// Clear cache manually
store.cache.Clear()

// Or disable cache for debugging
store.cache = nil // Will cause panic - for debugging only
```

### Issue: "Timeout errors"

**Solution:**
```go
// Increase timeout
store.httpClient.Timeout = 60 * time.Second
```

## 📚 Next Steps

1. Read [ARCHITECTURE.md](./ARCHITECTURE.md) for detailed design
2. Check [README.md](./README.md) for full API documentation
3. Explore [example_advanced.go](./example_advanced.go) for more examples
4. Run tests: `make test`

## 🤝 Contributing

Feel free to submit issues or pull requests to improve this implementation!

