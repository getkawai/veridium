# Contributor Discovery Implementation

## Overview
Implemented a comprehensive contributor discovery and load balancing system that enables desktop users to automatically find and connect to the best available contributor nodes in the DePIN network.

## Problem Solved
Previously, contributors could register and send heartbeats, but desktop users had no way to discover or select contributors. This created a critical gap in the network architecture.

## Solution Architecture

### 1. Enhanced Contributor Metadata (`pkg/store/contributor.go`)
Extended `ContributorData` struct with discovery fields:
- **Geographic**: `Region` (e.g., "us-west", "eu-central")
- **Capacity**: `AvailableModels`, `CPUCores`, `TotalRAM`, `AvailableRAM`, `GPUModel`, `GPUMemory`
- **Performance**: `ActiveRequests`, `TotalRequests`, `AvgResponseTime`, `SuccessRate`
- **Health**: `LastHealthCheck` timestamp

### 2. Metrics Update System (`pkg/store/contributor.go`)
Added `UpdateContributorMetrics()` method:
- Updates all discovery metadata atomically
- Called every 30s by contributor heartbeat
- Stores parsed hardware specs for efficient filtering

### 3. Enhanced Heartbeat (`cmd/server/api/services/kronk/kronk.go`)
Contributors now send rich metrics every 30s:
- Auto-detect region based on timezone
- List available models from cache
- Calculate success rate and avg response time
- Report current load and hardware capacity

Helper functions:
- `detectRegion()`: Maps UTC offset to geographic regions
- `getAvailableModels()`: Extracts model IDs from cache

### 4. Contributor Selector Service (`internal/services/contributor_selector.go`)
Smart selection algorithm with scoring system:
- **Success Rate** (30 points): Higher = better
- **Response Time** (25 points): Lower = better  
- **Current Load** (20 points): Lower = better
- **Hardware Capacity** (15 points): More RAM/GPU = better
- **Region Match** (10 points): Bonus for preferred region

Filtering criteria:
- Minimum RAM requirement
- Minimum GPU memory requirement
- Maximum concurrent load
- Required model availability
- Region preference

### 5. Desktop Integration (`internal/services/deai_service.go`)
Added three Wails-exposed methods:

#### `GetAvailableContributors() ([]*ContributorInfo, error)`
Returns all online contributors with scores, sorted by quality.

#### `SelectBestContributor(preferredRegion, requiredModel, minRAM, minGPUMemory) (*ContributorInfo, error)`
Selects optimal contributor based on criteria.

#### `GetContributorStats() (map[string]interface{}, error)`
Returns network statistics:
- Total online/active contributors
- Total requests served
- Average success rate
- Contributors by region
- Contributors by model

## Data Flow

### Contributor Side (Every 30s)
1. Detect region from timezone
2. Get available models from cache
3. Calculate performance metrics
4. Update KV store via `UpdateContributorMetrics()`

### Desktop Side (On-Demand)
1. User initiates AI request
2. Desktop calls `SelectBestContributor()` with criteria
3. Selector queries KV for online contributors
4. Filters by requirements (RAM, GPU, model)
5. Scores remaining candidates
6. Returns best contributor endpoint
7. Desktop connects to contributor URL

## Key Features

### Automatic Region Detection
Uses timezone offset to map contributors to regions:
- `us-west`, `us-east`
- `eu-west`, `eu-east`
- `asia-west`, `asia-east`

### Load Balancing
Exponential penalty for high load:
- 0 requests = 20 points
- 10 requests = 7.4 points
- 20 requests = 2.7 points

### Health Checks
Contributors with stale health checks (>2 minutes) are filtered out automatically.

### Hardware-Aware Selection
Considers actual hardware capacity:
- RAM: Normalized to 32GB max
- GPU: Normalized to 24GB VRAM max

## Usage Example

### Frontend (TypeScript)
```typescript
// Get best contributor for Llama 3.1 70B
const contributor = await SelectBestContributor(
  "us-west",      // preferred region
  "llama-3.1-70b", // required model
  32,             // min 32GB RAM
  24              // min 24GB GPU VRAM
);

// Connect to contributor
const response = await fetch(`${contributor.endpoint_url}/v1/chat/completions`, {
  method: 'POST',
  body: JSON.stringify(chatRequest)
});
```

### View All Contributors
```typescript
const contributors = await GetAvailableContributors();
console.log(`Found ${contributors.length} contributors`);
contributors.forEach(c => {
  console.log(`${c.wallet_address}: ${c.region}, ${c.score} points`);
});
```

### Network Stats
```typescript
const stats = await GetContributorStats();
console.log(`Online: ${stats.total_online}`);
console.log(`Regions: ${JSON.stringify(stats.regions)}`);
console.log(`Models: ${JSON.stringify(stats.models)}`);
```

## Files Modified

1. `pkg/store/contributor.go`
   - Extended `ContributorData` struct (17 new fields)
   - Added `UpdateContributorMetrics()` method
   - Added `ContributorMetrics` struct

2. `cmd/server/api/services/kronk/kronk.go`
   - Enhanced heartbeat goroutine with metrics
   - Added `detectRegion()` helper
   - Added `getAvailableModels()` helper

3. `internal/services/deai_service.go`
   - Added `contributorSelector` field
   - Updated `NewDeAIService()` constructor
   - Added `GetAvailableContributors()` method
   - Added `SelectBestContributor()` method
   - Added `GetContributorStats()` method
   - Added `ContributorInfo` struct

4. `internal/services/contributor_selector.go` (NEW)
   - Core selection logic
   - Scoring algorithm
   - Filtering system
   - Statistics aggregation

## Testing

Build verification:
```bash
go build -o /dev/null ./cmd/server
# ✓ Builds successfully
```

## Next Steps

1. **Frontend Integration**: Update desktop app to use new discovery methods
2. **Fallback Logic**: Handle case when no contributors available
3. **Metrics Dashboard**: Build UI to visualize contributor network
4. **Performance Monitoring**: Track selection accuracy and response times
5. **Geographic Optimization**: Add latency-based region detection
6. **Model Preloading**: Suggest contributors to preload popular models

## Benefits

- **Zero Configuration**: Desktop users automatically discover contributors
- **Optimal Performance**: Smart selection based on multiple factors
- **Load Distribution**: Prevents overloading single contributors
- **Geographic Optimization**: Prefers nearby contributors
- **Hardware Matching**: Routes heavy models to powerful nodes
- **Fault Tolerance**: Automatically excludes unhealthy contributors
- **Scalability**: Supports unlimited contributors via KV store

## Technical Decisions

### Why Cloudflare KV?
- Already in use for contributor registration
- No additional infrastructure needed
- Global edge network for low latency
- Simple key-value model fits use case

### Why 30s Heartbeat?
- Balance between freshness and KV write costs
- Matches existing heartbeat interval
- Sufficient for detecting offline nodes (2min timeout)

### Why Score-Based Selection?
- Flexible: Easy to adjust weights
- Transparent: Clear why contributor was chosen
- Extensible: Can add more factors later
- Fair: Considers multiple dimensions

## Monitoring

Key metrics to track:
- Contributor discovery success rate
- Average selection time
- Load distribution across contributors
- Regional coverage
- Model availability

## Security Considerations

- Contributors self-report metrics (trust model)
- Desktop validates contributor endpoints
- Future: Add reputation system based on actual performance
- Future: Implement contributor verification via challenges
