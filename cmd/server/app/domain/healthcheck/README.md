# Contributor Health Check API

## Overview

Endpoint health check untuk contributor memungkinkan client untuk langsung mengecek status contributor tanpa melalui KV Store. Ini memberikan data **real-time** dan lebih **scalable**.

## Endpoint

### `GET /v1/health`

Returns the current health status of the contributor.

**Response:**

```json
{
  "status": "online",
  "active_requests": 5,
  "is_busy": true,
  "available_models": [
    "Qwen3-8B-GGUF/Qwen3-8B-Q8_0.gguf",
    "Llama-3-8B-GGUF/Llama-3-8B-Q8_0.gguf"
  ]
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `status` | string | Contributor status: `"online"` or `"offline"` |
| `active_requests` | integer | Current number of active requests being processed |
| `is_busy` | boolean | `true` if `active_requests > 0` |
| `available_models` | array | List of model IDs available on this contributor |

**HTTP Status Codes:**

- `200 OK` - Contributor is online and responding
- `503 Service Unavailable` - Contributor is offline or not responding

## Client Usage Example

### JavaScript/TypeScript

```typescript
interface ContributorHealth {
  status: 'online' | 'offline';
  active_requests: number;
  is_busy: boolean;
  available_models: string[];
}

async function checkContributorHealth(endpointUrl: string): Promise<ContributorHealth> {
  const response = await fetch(`${endpointUrl}/v1/health`, {
    method: 'GET',
    signal: AbortSignal.timeout(5000), // 5 second timeout
  });
  
  if (!response.ok) {
    throw new Error(`Contributor ${endpointUrl} is not responding`);
  }
  
  return await response.json();
}

async function selectBestContributor(contributors: string[]): Promise<string | null> {
  // Filter online and not busy contributors
  const available = await Promise.all(
    contributors.map(async (url) => {
      try {
        const health = await checkContributorHealth(url);
        return { url, health };
      } catch {
        return null; // Contributor offline
      }
    })
  );
  
  // Filter out offline and busy contributors
  const healthy = available
    .filter((c): c is { url: string; health: ContributorHealth } => 
      c !== null && c.health.status === 'online' && !c.health.is_busy
    );
  
  if (healthy.length === 0) {
    return null; // No available contributors
  }
  
  // Select one with least active requests
  healthy.sort((a, b) => a.health.active_requests - b.health.active_requests);
  
  return healthy[0].url;
}
```

### Go

```go
type ContributorHealth struct {
    Status          string   `json:"status"`
    ActiveRequests  int64    `json:"active_requests"`
    IsBusy          bool     `json:"is_busy"`
    AvailableModels []string `json:"available_models"`
}

func CheckContributorHealth(ctx context.Context, endpointURL string) (*ContributorHealth, error) {
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpointURL+"/v1/health", nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("contributor returned status: %d", resp.StatusCode)
    }
    
    var health ContributorHealth
    if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
        return nil, err
    }
    
    return &health, nil
}

func SelectBestContributor(ctx context.Context, contributors []string) (string, error) {
    type result struct {
        url    string
        health *ContributorHealth
        err    error
    }
    
    // Check all contributors concurrently
    results := make(chan result, len(contributors))
    for _, url := range contributors {
        go func(url string) {
            health, err := CheckContributorHealth(ctx, url)
            results <- result{url: url, health: health, err: err}
        }(url)
    }
    
    // Collect healthy contributors
    var available []result
    for i := 0; i < len(contributors); i++ {
        r := <-results
        if r.err == nil && r.health.Status == "online" && !r.health.IsBusy {
            available = append(available, r)
        }
    }
    
    if len(available) == 0 {
        return "", fmt.Errorf("no available contributors")
    }
    
    // Select one with least active requests
    sort.Slice(available, func(i, j int) bool {
        return available[i].health.ActiveRequests < available[j].health.ActiveRequests
    })
    
    return available[0].url, nil
}
```

### cURL

```bash
# Check single contributor
curl http://localhost:8080/v1/health

# Response:
# {
#   "status": "online",
#   "active_requests": 0,
#   "is_busy": false,
#   "available_models": ["Qwen3-8B-GGUF/Qwen3-8B-Q8_0.gguf"]
# }

# Check multiple contributors (bash)
for url in "http://contributor1:8080" "http://contributor2:8080"; do
  echo "Checking $url..."
  curl -s "$url/v1/health" | jq -c "{url: \"$url\", status: .status, is_busy: .is_busy, active: .active_requests}"
done
```

## Architecture Comparison

### Before (Push-based with KV Store)

```
┌─────────────┐      Heartbeat       ┌─────────────┐
│ Contributor │ ───────────────────► │  KV Store   │
│             │    (every 30s)       │             │
└─────────────┘                      └──────┬──────┘
                                           │
                                           │ Read (stale)
                                           │
                                     ┌─────▼──────┐
                                     │   Client   │
                                     └────────────┘
```

**Issues:**
- Stale data (up to 30s old)
- KV Store bottleneck
- Rate limiting concerns

### After (Pull-based Direct Check)

```
┌─────────────┐
│ Contributor │
│  /v1/health │
└──────┬──────┘
       │
       │ Direct HTTP
       │ (real-time)
       │
┌──────▼──────┐
│   Client    │
└─────────────┘
```

**Benefits:**
- ✅ Real-time data
- ✅ No KV Store dependency
- ✅ No rate limiting
- ✅ Scalable (client controls check frequency)
- ✅ Accurate busy status

## Best Practices

### 1. Timeout Handling

Always use timeouts when checking contributor health:

```typescript
const timeout = 5000; // 5 seconds
const response = await fetch(url, { 
  signal: AbortSignal.timeout(timeout) 
});
```

### 2. Concurrent Checks

Check multiple contributors concurrently for faster selection:

```go
// Check all contributors in parallel
for _, url := range contributors {
    go checkHealth(url)
}
```

### 3. Caching

Cache health status briefly to avoid excessive checks:

```typescript
const CACHE_TTL = 10000; // 10 seconds
const cache = new Map<string, {health: Health, timestamp: number}>();

function getCachedHealth(url: string): Health | null {
  const cached = cache.get(url);
  if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
    return cached.health;
  }
  return null;
}
```

### 4. Fallback Strategy

Have a fallback when no contributors are available:

```typescript
const contributor = await selectBestContributor(urls);
if (!contributor) {
  // Fallback: queue request or return error
  await queueRequest(request);
}
```

## Migration Guide

### From KV Store to Direct Health Check

**Before:**

```typescript
// Get contributors from KV
const contributors = await kv.getOnlineContributors();

// Filter by LastSeen (stale!)
const available = contributors.filter(c => 
  Date.now() - new Date(c.LastSeen).getTime() < 120000
);
```

**After:**

```typescript
// Direct health check (real-time)
const available = await Promise.all(
  contributors.map(async (c) => {
    const health = await fetch(`${c.endpoint_url}/v1/health`);
    return health.status === 'online' && !health.is_busy;
  })
);
```

## Monitoring

Track these metrics on client side:

- Health check success rate
- Average response time
- Contributor availability %
- Busy ratio (busy checks / total checks)

Example:

```typescript
const metrics = {
  totalChecks: 0,
  successfulChecks: 0,
  contributorsFound: 0,
  contributorsBusy: 0,
};

async function checkWithMetrics(url: string) {
  metrics.totalChecks++;
  try {
    const health = await checkContributorHealth(url);
    metrics.successfulChecks++;
    if (health.is_busy) metrics.contributorsBusy++;
    if (health.status === 'online') metrics.contributorsFound++;
    return health;
  } catch {
    return null;
  }
}
```
