package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kawai-network/veridium/pkg/types"
)

const (
	PeriodCounterKey = "settlement:period_counter"
)

// PeriodCounter tracks the last settlement period ID
type PeriodCounter struct {
	LastPeriodID int64            `json:"last_period_id"`
	RewardType   types.RewardType `json:"reward_type"`
}

// GetNextPeriodID returns the next sequential period ID and increments counter
func (s *KVStore) GetNextPeriodID(ctx context.Context, rewardType types.RewardType) (int64, error) {
	key := fmt.Sprintf("%s:%s", PeriodCounterKey, rewardType)

	// Try to get current counter
	value, err := s.client.GetValue(ctx, s.settlementsNamespaceID, key)

	var counter PeriodCounter
	if err != nil {
		// First time - start from period 1
		counter = PeriodCounter{
			LastPeriodID: 0,
			RewardType:   rewardType,
		}
	} else {
		if err := json.Unmarshal(value, &counter); err != nil {
			return 0, fmt.Errorf("failed to unmarshal counter: %w", err)
		}
	}

	// Increment and save
	nextPeriod := counter.LastPeriodID + 1
	counter.LastPeriodID = nextPeriod

	data, err := json.Marshal(counter)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal counter: %w", err)
	}

	err = s.client.SetValue(ctx, s.settlementsNamespaceID, key, data)
	if err != nil {
		return 0, fmt.Errorf("failed to save counter: %w", err)
	}

	return nextPeriod, nil
}

// GetCurrentPeriodID returns the last used period ID without incrementing
func (s *KVStore) GetCurrentPeriodID(ctx context.Context, rewardType types.RewardType) (int64, error) {
	key := fmt.Sprintf("%s:%s", PeriodCounterKey, rewardType)

	value, err := s.client.GetValue(ctx, s.settlementsNamespaceID, key)

	if err != nil {
		// No counter yet - return 0
		return 0, nil
	}

	var counter PeriodCounter
	if err := json.Unmarshal(value, &counter); err != nil {
		return 0, fmt.Errorf("failed to unmarshal counter: %w", err)
	}

	return counter.LastPeriodID, nil
}

// ResetPeriodCounter resets counter to 0 (for testing/migration)
func (s *KVStore) ResetPeriodCounter(ctx context.Context, rewardType types.RewardType) error {
	key := fmt.Sprintf("%s:%s", PeriodCounterKey, rewardType)

	counter := PeriodCounter{
		LastPeriodID: 0,
		RewardType:   rewardType,
	}

	data, err := json.Marshal(counter)
	if err != nil {
		return fmt.Errorf("failed to marshal counter: %w", err)
	}

	err = s.client.SetValue(ctx, s.settlementsNamespaceID, key, data)

	return err
}
