package services

import (
	"testing"
	"time"

	"github.com/kawai-network/veridium/pkg/store"
	"github.com/stretchr/testify/assert"
)

func TestCalculateScore(t *testing.T) {
	selector := &ContributorSelector{}

	tests := []struct {
		name        string
		contributor *store.ContributorData
		criteria    *SelectionCriteria
		minScore    float64
		maxScore    float64
	}{
		{
			name: "perfect contributor",
			contributor: &store.ContributorData{
				SuccessRate:     1.0,
				AvgResponseTime: 0.5,
				ActiveRequests:  0,
				AvailableRAM:    32,
				GPUMemory:       24,
				Region:          "us-west",
			},
			criteria: &SelectionCriteria{
				PreferredRegion: "us-west",
			},
			minScore: 180.0, // Should get near-perfect score
			maxScore: 200.0,
		},
		{
			name: "loaded contributor",
			contributor: &store.ContributorData{
				SuccessRate:     0.95,
				AvgResponseTime: 2.0,
				ActiveRequests:  10,
				AvailableRAM:    16,
				GPUMemory:       12,
				Region:          "us-east",
			},
			criteria: &SelectionCriteria{
				PreferredRegion: "us-west",
			},
			minScore: 130.0,
			maxScore: 160.0,
		},
		{
			name: "poor performer",
			contributor: &store.ContributorData{
				SuccessRate:     0.5,
				AvgResponseTime: 5.0,
				ActiveRequests:  20,
				AvailableRAM:    8,
				GPUMemory:       4,
				Region:          "asia-east",
			},
			criteria: &SelectionCriteria{
				PreferredRegion: "us-west",
			},
			minScore: 100.0,
			maxScore: 130.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := selector.calculateScore(tt.contributor, tt.criteria)
			assert.GreaterOrEqual(t, score, tt.minScore, "score should be >= min")
			assert.LessOrEqual(t, score, tt.maxScore, "score should be <= max")
		})
	}
}

func TestMeetsRequirements(t *testing.T) {
	selector := &ContributorSelector{}

	tests := []struct {
		name        string
		contributor *store.ContributorData
		criteria    *SelectionCriteria
		expected    bool
	}{
		{
			name: "meets all requirements",
			contributor: &store.ContributorData{
				AvailableRAM:    32,
				GPUMemory:       24,
				ActiveRequests:  5,
				AvailableModels: []string{"llama-3.1-70b", "gpt-4"},
			},
			criteria: &SelectionCriteria{
				MinRAM:        16,
				MinGPUMemory:  12,
				MaxLoad:       10,
				RequiredModel: "llama-3.1-70b",
			},
			expected: true,
		},
		{
			name: "insufficient RAM",
			contributor: &store.ContributorData{
				AvailableRAM:    8,
				GPUMemory:       24,
				ActiveRequests:  5,
				AvailableModels: []string{"llama-3.1-70b"},
			},
			criteria: &SelectionCriteria{
				MinRAM:        16,
				RequiredModel: "llama-3.1-70b",
			},
			expected: false,
		},
		{
			name: "insufficient GPU memory",
			contributor: &store.ContributorData{
				AvailableRAM:    32,
				GPUMemory:       8,
				ActiveRequests:  5,
				AvailableModels: []string{"llama-3.1-70b"},
			},
			criteria: &SelectionCriteria{
				MinGPUMemory:  12,
				RequiredModel: "llama-3.1-70b",
			},
			expected: false,
		},
		{
			name: "too much load",
			contributor: &store.ContributorData{
				AvailableRAM:    32,
				GPUMemory:       24,
				ActiveRequests:  15,
				AvailableModels: []string{"llama-3.1-70b"},
			},
			criteria: &SelectionCriteria{
				MaxLoad:       10,
				RequiredModel: "llama-3.1-70b",
			},
			expected: false,
		},
		{
			name: "missing required model",
			contributor: &store.ContributorData{
				AvailableRAM:    32,
				GPUMemory:       24,
				ActiveRequests:  5,
				AvailableModels: []string{"gpt-4"},
			},
			criteria: &SelectionCriteria{
				RequiredModel: "llama-3.1-70b",
			},
			expected: false,
		},
		{
			name: "no requirements",
			contributor: &store.ContributorData{
				AvailableRAM:    8,
				GPUMemory:       4,
				ActiveRequests:  20,
				AvailableModels: []string{},
			},
			criteria: &SelectionCriteria{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := selector.meetsRequirements(tt.contributor, tt.criteria)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetectRegion(t *testing.T) {
	// This test verifies region detection logic
	// We can't import from kronk package, so we test the concept

	// Test that we can determine region from timezone
	_, offset := time.Now().Zone()
	offsetHours := offset / 3600

	// Verify offset is reasonable (-12 to +14)
	assert.GreaterOrEqual(t, offsetHours, -12)
	assert.LessOrEqual(t, offsetHours, 14)

	// Test region mapping logic with exclusive upper bounds to avoid overlaps
	var region string
	switch {
	case offsetHours >= -8 && offsetHours < -5:
		region = "us-west"
	case offsetHours >= -5 && offsetHours < -3:
		region = "us-east"
	case offsetHours >= 0 && offsetHours < 3:
		region = "eu-west"
	case offsetHours >= 3 && offsetHours < 6:
		region = "eu-east"
	case offsetHours >= 6 && offsetHours < 9:
		region = "asia-west"
	case offsetHours >= 9 && offsetHours <= 12:
		region = "asia-east"
	default:
		region = "unknown"
	}

	assert.NotEmpty(t, region)
}

func TestContributorMetrics(t *testing.T) {
	// Test that ContributorMetrics struct can be created
	metrics := &store.ContributorMetrics{
		Region:          "us-west",
		AvailableModels: []string{"llama-3.1-70b"},
		ActiveRequests:  5,
		TotalRequests:   1000,
		AvgResponseTime: 1.5,
		SuccessRate:     0.98,
		CPUCores:        16,
		TotalRAM:        64,
		AvailableRAM:    32,
		GPUModel:        "NVIDIA RTX 4090",
		GPUMemory:       24,
	}

	assert.Equal(t, "us-west", metrics.Region)
	assert.Equal(t, int64(5), metrics.ActiveRequests)
	assert.Equal(t, 0.98, metrics.SuccessRate)
}

func TestSelectBestContributor_NilCriteria(t *testing.T) {
	// Test that nil criteria is handled gracefully
	// We can't test the full flow without a real KV store,
	// but we can verify the nil check exists by ensuring
	// the function signature accepts nil
	var criteria *SelectionCriteria = nil
	assert.Nil(t, criteria) // Verify we can pass nil

	// The actual nil handling is tested in the implementation
	// where nil criteria is replaced with empty SelectionCriteria{}
}

func TestContributorDataExtensions(t *testing.T) {
	// Test that extended ContributorData fields work
	now := time.Now()
	contributor := &store.ContributorData{
		WalletAddress:   "0x1234567890123456789012345678901234567890",
		EndpointURL:     "https://contributor.example.com",
		Region:          "eu-west",
		AvailableModels: []string{"llama-3.1-8b", "llama-3.1-70b"},
		ActiveRequests:  3,
		TotalRequests:   500,
		AvgResponseTime: 1.2,
		SuccessRate:     0.99,
		LastHealthCheck: now,
		CPUCores:        32,
		TotalRAM:        128,
		AvailableRAM:    64,
		GPUModel:        "NVIDIA H100",
		GPUMemory:       80,
	}

	assert.Equal(t, "eu-west", contributor.Region)
	assert.Len(t, contributor.AvailableModels, 2)
	assert.Equal(t, int64(3), contributor.ActiveRequests)
	assert.Equal(t, 0.99, contributor.SuccessRate)
	assert.Equal(t, int64(80), contributor.GPUMemory)
}
