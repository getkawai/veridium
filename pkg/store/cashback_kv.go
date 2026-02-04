package store

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
)

// Cashback KV operations - using dedicated cashback namespace

// StoreCashbackData stores cashback data in the cashback namespace
func (s *KVStore) StoreCashbackData(ctx context.Context, key string, data []byte) error {
	err := s.client.SetValue(ctx, s.cashbackNamespaceID, key, data)
	if err != nil {
		return fmt.Errorf("failed to write cashback data to KV: %w", err)
	}
	return nil
}

// StoreCashbackDataWithTTL stores cashback data with an expiration time (seconds)
func (s *KVStore) StoreCashbackDataWithTTL(ctx context.Context, key string, data []byte, ttl int) error {
	err := s.client.SetValueWithTTL(ctx, s.cashbackNamespaceID, key, data, ttl)
	if err != nil {
		return fmt.Errorf("failed to write cashback data with TTL to KV: %w", err)
	}
	return nil
}

// GetCashbackData retrieves cashback data from the cashback namespace
func (s *KVStore) GetCashbackData(ctx context.Context, key string) ([]byte, error) {
	value, err := s.client.GetValue(ctx, s.cashbackNamespaceID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get cashback data from KV: %w", err)
	}
	return value, nil
}

// DeleteCashbackData deletes cashback data from the cashback namespace
func (s *KVStore) DeleteCashbackData(ctx context.Context, key string) error {
	err := s.client.DeleteValue(ctx, s.cashbackNamespaceID, key)
	if err != nil {
		return fmt.Errorf("failed to delete cashback data from KV: %w", err)
	}
	return nil
}

// GetSettledCashbackPeriods retrieves the list of settled cashback periods
// This is used to optimize GetClaimableCashbackRecords by only checking relevant periods
func (s *KVStore) GetSettledCashbackPeriods(ctx context.Context) ([]uint64, error) {
	const key = "cashback_settled_periods"
	data, err := s.GetCashbackData(ctx, key)
	if err != nil {
		// If key doesn't exist, return empty list
		return []uint64{}, nil
	}

	var periods []uint64
	if err := json.Unmarshal(data, &periods); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settled periods: %w", err)
	}

	return periods, nil
}

// AddSettledCashbackPeriod adds a period to the list of settled cashback periods
func (s *KVStore) AddSettledCashbackPeriod(ctx context.Context, period uint64) error {
	const key = "cashback_settled_periods"

	// Get existing periods
	periods, err := s.GetSettledCashbackPeriods(ctx)
	if err != nil {
		return fmt.Errorf("failed to get settled periods: %w", err)
	}

	// Check if period already exists
	for _, p := range periods {
		if p == period {
			// Already exists, no need to add
			return nil
		}
	}

	// Add new period and sort
	periods = append(periods, period)
	sort.Slice(periods, func(i, j int) bool {
		return periods[i] < periods[j]
	})

	// Store updated list
	data, err := json.Marshal(periods)
	if err != nil {
		return fmt.Errorf("failed to marshal settled periods: %w", err)
	}

	return s.StoreCashbackData(ctx, key, data)
}
