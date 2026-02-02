package services

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/kawai-network/veridium/pkg/store"
)

// ContributorSelector handles selection of best contributor for inference requests
type ContributorSelector struct {
	kv store.Store
}

// NewContributorSelector creates a new contributor selector
func NewContributorSelector(kv store.Store) *ContributorSelector {
	return &ContributorSelector{
		kv: kv,
	}
}

// SelectionCriteria defines criteria for contributor selection
type SelectionCriteria struct {
	PreferredRegion string // Preferred geographic region
	RequiredModel   string // Required model ID
	MinRAM          int64  // Minimum RAM in GB
	MinGPUMemory    int64  // Minimum GPU VRAM in GB
	MaxLoad         int64  // Maximum active requests
}

// ContributorScore represents a contributor with its selection score
type ContributorScore struct {
	Contributor *store.ContributorData
	Score       float64
}

// SelectBestContributor selects the best available contributor based on criteria
// Returns nil if no suitable contributor is found
func (s *ContributorSelector) SelectBestContributor(ctx context.Context, criteria *SelectionCriteria) (*store.ContributorData, error) {
	// Handle nil criteria with default values
	if criteria == nil {
		criteria = &SelectionCriteria{}
	}

	// Get all online contributors
	contributors, err := s.kv.GetOnlineContributors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get online contributors: %w", err)
	}

	if len(contributors) == 0 {
		return nil, fmt.Errorf("no online contributors available")
	}

	// Filter and score contributors
	candidates := make([]*ContributorScore, 0)
	for _, c := range contributors {
		// Skip if not active
		if !c.IsActive {
			continue
		}

		// Skip if last health check is too old (>2 minutes)
		if time.Since(c.LastHealthCheck) > 2*time.Minute {
			continue
		}

		// Apply filters
		if !s.meetsRequirements(c, criteria) {
			continue
		}

		// Calculate score
		score := s.calculateScore(c, criteria)
		candidates = append(candidates, &ContributorScore{
			Contributor: c,
			Score:       score,
		})
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no contributors meet the requirements")
	}

	// Sort by score (highest first)
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	// Return best candidate
	return candidates[0].Contributor, nil
}

// meetsRequirements checks if contributor meets minimum requirements
func (s *ContributorSelector) meetsRequirements(c *store.ContributorData, criteria *SelectionCriteria) bool {
	// Check RAM requirement
	if criteria.MinRAM > 0 && c.AvailableRAM < criteria.MinRAM {
		return false
	}

	// Check GPU memory requirement
	if criteria.MinGPUMemory > 0 && c.GPUMemory < criteria.MinGPUMemory {
		return false
	}

	// Check load limit
	if criteria.MaxLoad > 0 && c.ActiveRequests >= criteria.MaxLoad {
		return false
	}

	// Check model availability
	if criteria.RequiredModel != "" {
		hasModel := false
		for _, model := range c.AvailableModels {
			if model == criteria.RequiredModel {
				hasModel = true
				break
			}
		}
		if !hasModel {
			return false
		}
	}

	return true
}

// calculateScore calculates selection score for a contributor
// Higher score = better candidate
func (s *ContributorSelector) calculateScore(c *store.ContributorData, criteria *SelectionCriteria) float64 {
	score := 100.0

	// Factor 1: Success Rate (0-30 points)
	// Higher success rate = better
	score += c.SuccessRate * 30.0

	// Factor 2: Response Time (0-25 points)
	// Lower response time = better
	// Assume ideal response time is 1s, penalize slower responses
	if c.AvgResponseTime > 0 {
		responseScore := 25.0 * (1.0 / (1.0 + c.AvgResponseTime))
		score += responseScore
	} else {
		// No data yet, give neutral score
		score += 12.5
	}

	// Factor 3: Current Load (0-20 points)
	// Lower load = better
	loadScore := 20.0
	if c.ActiveRequests > 0 {
		// Penalize based on active requests (exponential penalty)
		loadScore = 20.0 * math.Exp(-float64(c.ActiveRequests)/10.0)
	}
	score += loadScore

	// Factor 4: Hardware Capacity (0-15 points)
	// More RAM and GPU memory = better
	ramScore := math.Min(float64(c.AvailableRAM)/32.0, 1.0) * 7.5 // Max at 32GB
	gpuScore := math.Min(float64(c.GPUMemory)/24.0, 1.0) * 7.5    // Max at 24GB
	score += ramScore + gpuScore

	// Factor 5: Region Match (0-10 points)
	// Matching region = bonus
	if criteria.PreferredRegion != "" && c.Region == criteria.PreferredRegion {
		score += 10.0
	}

	return score
}

// GetAvailableContributors returns all online contributors with their scores
func (s *ContributorSelector) GetAvailableContributors(ctx context.Context) ([]*ContributorScore, error) {
	contributors, err := s.kv.GetOnlineContributors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get online contributors: %w", err)
	}

	scores := make([]*ContributorScore, 0)
	for _, c := range contributors {
		if !c.IsActive {
			continue
		}

		// Calculate score with no specific criteria
		score := s.calculateScore(c, &SelectionCriteria{})
		scores = append(scores, &ContributorScore{
			Contributor: c,
			Score:       score,
		})
	}

	// Sort by score
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	return scores, nil
}

// GetContributorStats returns statistics about available contributors
func (s *ContributorSelector) GetContributorStats(ctx context.Context) (map[string]interface{}, error) {
	contributors, err := s.kv.GetOnlineContributors(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get online contributors: %w", err)
	}

	stats := map[string]interface{}{
		"total_online":     len(contributors),
		"total_active":     0,
		"total_requests":   int64(0),
		"avg_success_rate": 0.0,
		"regions":          make(map[string]int),
		"models":           make(map[string]int),
	}

	activeCount := 0
	totalRequests := int64(0)
	totalSuccessRate := 0.0
	regions := make(map[string]int)
	models := make(map[string]int)

	for _, c := range contributors {
		if !c.IsActive {
			continue
		}

		activeCount++
		totalRequests += c.TotalRequests
		totalSuccessRate += c.SuccessRate

		if c.Region != "" {
			regions[c.Region]++
		}

		for _, model := range c.AvailableModels {
			models[model]++
		}
	}

	stats["total_active"] = activeCount
	stats["total_requests"] = totalRequests
	if activeCount > 0 {
		stats["avg_success_rate"] = totalSuccessRate / float64(activeCount)
	}
	stats["regions"] = regions
	stats["models"] = models

	return stats, nil
}
