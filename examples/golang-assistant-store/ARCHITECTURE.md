# Architecture Documentation

## 📐 System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Application                       │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                       AssistantStore                             │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │  Public API Methods                                         │ │
│  │  - GetAgentIndex(locale)                                   │ │
│  │  - GetAgent(identifier, locale)                            │ │
│  │  - SearchAgents(locale, query)                             │ │
│  │  - FilterByCategory(locale, category)                      │ │
│  │  - GetCategories(locale)                                   │ │
│  │  - GetAgentIndexWithFilter(locale, filter)                 │ │
│  └────────────────────────────────────────────────────────────┘ │
│                             │                                    │
│  ┌──────────────────────────┼──────────────────────────┐        │
│  │                          │                           │        │
│  ▼                          ▼                           ▼        │
│ ┌──────────┐        ┌──────────┐              ┌──────────┐      │
│ │  Cache   │        │  HTTP    │              │  URL     │      │
│ │  Layer   │        │  Client  │              │  Builder │      │
│ └──────────┘        └──────────┘              └──────────┘      │
└─────────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    NPM Registry Mirror                           │
│  https://registry.npmmirror.com/@lobehub/agents-index/...       │
│                                                                  │
│  ├── index.en-US.json     (Agent list - English)                │
│  ├── index.zh-CN.json     (Agent list - Chinese)                │
│  ├── web-dev.en-US.json   (Agent detail - English)              │
│  └── web-dev.zh-CN.json   (Agent detail - Chinese)              │
└─────────────────────────────────────────────────────────────────┘
```

## 🏗️ Component Design

### 1. AssistantStore

**Responsibility:** Main facade for fetching and managing assistant data.

**Key Components:**
- `baseURL`: Configuration for NPM registry endpoint
- `httpClient`: HTTP client with timeout configuration
- `cache`: In-memory cache for performance optimization

**Design Patterns:**
- **Facade Pattern**: Provides simple interface to complex subsystems
- **Strategy Pattern**: URL building strategy based on locale
- **Template Method**: Common fetch logic with locale fallback

### 2. Cache

**Responsibility:** In-memory caching with automatic expiration.

**Features:**
- Thread-safe operations using `sync.RWMutex`
- Automatic cleanup of expired items
- Configurable TTL per item

**Design Patterns:**
- **Singleton-like**: One cache instance per store
- **Observer Pattern**: Background goroutine for cleanup

### 3. HTTP Client

**Configuration:**
- Timeout: 30 seconds
- Automatic retry: Fallback to default locale on 404
- Force cache: Leverages HTTP caching headers

## 🔄 Data Flow

### Fetching Agent Index

```
User Request
    │
    ▼
GetAgentIndex(locale)
    │
    ├─→ Check Cache ──→ [HIT] Return cached data
    │                       
    └─→ [MISS]
         │
         ▼
    Build URL: index.{locale}.json
         │
         ▼
    HTTP GET Request
         │
         ├─→ [200 OK] ──→ Parse JSON ──→ Cache ──→ Return
         │
         ├─→ [404] ──→ Retry with default locale
         │                │
         │                └─→ [200 OK] ──→ Parse JSON ──→ Cache ──→ Return
         │
         └─→ [Error] ──→ Return error
```

### Fetching Agent Detail

```
User Request
    │
    ▼
GetAgent(identifier, locale)
    │
    ├─→ Check Cache ──→ [HIT] Return cached data
    │
    └─→ [MISS]
         │
         ▼
    Build URL: {identifier}.{locale}.json
         │
         ▼
    HTTP GET Request
         │
         ├─→ [200 OK] ──→ Parse JSON ──→ Cache ──→ Return
         │
         ├─→ [404] ──→ Retry with default locale
         │                │
         │                └─→ [200 OK] ──→ Parse JSON ──→ Cache ──→ Return
         │
         └─→ [Error] ──→ Return error
```

## 🎯 Design Decisions

### 1. Why In-Memory Cache?

**Pros:**
- Fast access (microseconds)
- No external dependencies
- Simple implementation
- Suitable for read-heavy workloads

**Cons:**
- Lost on restart
- Not shared across instances
- Memory usage

**Alternatives Considered:**
- Redis: Overkill for simple use case
- File cache: Slower, more complex
- No cache: Too many network requests

### 2. Why Fallback to Default Locale?

**Rationale:**
- Not all agents have translations in all locales
- Better UX to show English than error
- Matches TypeScript implementation behavior

### 3. Why Separate Index and Detail Endpoints?

**Benefits:**
- Faster initial load (index is smaller)
- Bandwidth optimization
- Better caching strategy (different TTL)
- Lazy loading of details

## 🔐 Security Considerations

### 1. SSRF Protection

**Current State:** Basic HTTP client without SSRF protection

**Recommendations:**
```go
// Add IP address validation
func isPrivateIP(ip string) bool {
    // Check for private IP ranges
    // 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
}

// Validate URL before fetch
func (s *AssistantStore) validateURL(url string) error {
    // Parse URL
    // Check scheme (only http/https)
    // Resolve hostname to IP
    // Check if IP is private
}
```

### 2. JSON Parsing Safety

**Current State:** Standard `json.Unmarshal` with struct validation

**Protections:**
- Struct tags define expected types
- Unknown fields are ignored
- No eval or code execution

### 3. Rate Limiting

**Recommendation:**
```go
type RateLimiter struct {
    requests int
    window   time.Duration
    mu       sync.Mutex
}

func (r *RateLimiter) Allow() bool {
    // Implement token bucket or sliding window
}
```

## 📊 Performance Characteristics

### Time Complexity

| Operation | Without Cache | With Cache |
|-----------|--------------|------------|
| GetAgentIndex | O(n) network | O(1) memory |
| GetAgent | O(1) network | O(1) memory |
| SearchAgents | O(n) + network | O(n) memory |
| FilterByCategory | O(n) + network | O(n) memory |

### Space Complexity

| Component | Space |
|-----------|-------|
| Cache (index) | O(n × m) where n=agents, m=locales |
| Cache (details) | O(k) where k=fetched agents |
| HTTP Client | O(1) |

### Network Requests

| Operation | Requests | Cache Impact |
|-----------|----------|--------------|
| First GetAgentIndex | 1-2 (with fallback) | Cached for 1 hour |
| Subsequent GetAgentIndex | 0 | From cache |
| First GetAgent | 1-2 (with fallback) | Cached for 24 hours |
| SearchAgents | 0-2 | Uses cached index |

## 🔧 Configuration Options

### Environment Variables

```bash
# Custom base URL
export AGENTS_INDEX_URL="https://your-cdn.com/agents"

# HTTP timeout (future enhancement)
export HTTP_TIMEOUT="30s"

# Cache TTL (future enhancement)
export CACHE_TTL_INDEX="3600s"
export CACHE_TTL_DETAIL="86400s"
```

### Code Configuration

```go
// Custom HTTP client
store := NewAssistantStore("")
store.httpClient = &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns: 100,
        IdleConnTimeout: 90 * time.Second,
    },
}

// Custom cache
store.cache = NewCacheWithConfig(CacheConfig{
    DefaultTTL: 1 * time.Hour,
    CleanupInterval: 10 * time.Minute,
})
```

## 🚀 Future Enhancements

### 1. Distributed Cache

```go
type CacheBackend interface {
    Get(key string) (interface{}, bool)
    Set(key string, value interface{}, ttl time.Duration)
    Delete(key string)
}

// Redis implementation
type RedisCache struct {
    client *redis.Client
}

// Memcached implementation
type MemcachedCache struct {
    client *memcache.Client
}
```

### 2. Metrics & Observability

```go
type Metrics struct {
    CacheHits   prometheus.Counter
    CacheMisses prometheus.Counter
    HTTPRequests prometheus.Counter
    Latency     prometheus.Histogram
}
```

### 3. Retry Logic

```go
type RetryConfig struct {
    MaxRetries int
    Backoff    time.Duration
}

func (s *AssistantStore) fetchWithRetry(url string, config RetryConfig) (*http.Response, error) {
    // Implement exponential backoff
}
```

### 4. Circuit Breaker

```go
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    state       State // Open, HalfOpen, Closed
}
```

## 📚 References

- [Go HTTP Client Best Practices](https://pkg.go.dev/net/http)
- [Caching Strategies](https://aws.amazon.com/caching/best-practices/)
- [NPM Registry API](https://github.com/npm/registry/blob/master/docs/REGISTRY-API.md)
- [LobeChat Agents Repository](https://github.com/lobehub/lobe-chat-agents)

