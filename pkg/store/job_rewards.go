package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/kawai-network/y/types"
)

// JobRewardRecord is defined in pkg/types to avoid circular dependency
type JobRewardRecord = types.JobRewardRecord

// SaveJobReward stores a job reward record in KV
// Key format: job_rewards:{contributor}:{timestamp_unix}
func (s *KVStore) SaveJobReward(ctx context.Context, record *JobRewardRecord) error {
	key := fmt.Sprintf("job_rewards:%s:%d", record.ContributorAddress, record.Timestamp.Unix())

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal job reward: %w", err)
	}

	if err := s.client.SetValue(ctx, s.contributorsNamespaceID, key, data); err != nil {
		return fmt.Errorf("failed to save job reward: %w", err)
	}

	return nil
}

// GetJobRewardsSinceLastSettlement retrieves all unsettled job rewards for a contributor
func (s *KVStore) GetJobRewardsSinceLastSettlement(ctx context.Context, contributorAddress string, rewardType types.RewardType) ([]*JobRewardRecord, error) {
	// List all keys for this contributor
	prefix := fmt.Sprintf("job_rewards:%s:", contributorAddress)

	keys, err := s.client.ListKeysSimple(ctx, s.contributorsNamespaceID, prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to list job rewards: %w", err)
	}

	var records []*JobRewardRecord

	for _, key := range keys {
		// Get the record
		value, err := s.client.GetValue(ctx, s.contributorsNamespaceID, key)
		if err != nil {
			slog.Warn("Failed to get job reward", "key", key, "error", err)
			continue
		}

		var record JobRewardRecord
		if err := json.Unmarshal(value, &record); err != nil {
			slog.Warn("Failed to unmarshal job reward", "key", key, "error", err)
			continue
		}

		// Filter by reward type and settlement status
		if record.RewardType == rewardType && !record.IsSettled {
			records = append(records, &record)
		}
	}

	return records, nil
}

// MarkJobRewardsAsSettled marks job rewards as settled for a specific period
func (s *KVStore) MarkJobRewardsAsSettled(ctx context.Context, contributorAddress string, periodID int64) error {
	prefix := fmt.Sprintf("job_rewards:%s:", contributorAddress)

	keys, err := s.client.ListKeysSimple(ctx, s.contributorsNamespaceID, prefix)
	if err != nil {
		return fmt.Errorf("failed to list job rewards: %w", err)
	}

	for _, key := range keys {
		// Get the record
		value, err := s.client.GetValue(ctx, s.contributorsNamespaceID, key)
		if err != nil {
			continue
		}

		var record JobRewardRecord
		if err := json.Unmarshal(value, &record); err != nil {
			continue
		}

		// Mark as settled
		if !record.IsSettled {
			record.IsSettled = true
			record.SettledPeriodID = periodID

			data, err := json.Marshal(record)
			if err != nil {
				continue
			}

			if err := s.client.SetValue(ctx, s.contributorsNamespaceID, key, data); err != nil {
				slog.Warn("Failed to mark job reward as settled", "key", key, "error", err)
			}
		}
	}

	return nil
}

// GetAllUnsettledJobRewards retrieves all unsettled job rewards across all contributors
// Used for settlement generation
func (s *KVStore) GetAllUnsettledJobRewards(ctx context.Context, rewardType types.RewardType) (map[string][]*JobRewardRecord, error) {
	// List all keys in contributors namespace with job_rewards prefix (with pagination)
	keys, err := s.client.ListAllKeys(ctx, s.contributorsNamespaceID, "job_rewards:")
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	result := make(map[string][]*JobRewardRecord)

	// Scan all job_rewards keys
	for _, key := range keys {
		// Get the record
		value, err := s.client.GetValue(ctx, s.contributorsNamespaceID, key)
		if err != nil {
			slog.Warn("Failed to get job reward", "key", key, "error", err)
			continue
		}

		var record JobRewardRecord
		if err := json.Unmarshal(value, &record); err != nil {
			slog.Warn("Failed to unmarshal job reward", "key", key, "error", err)
			continue
		}

		// Filter by reward type and settlement status
		if record.RewardType != rewardType || record.IsSettled {
			continue
		}

		// Group by contributor address
		result[record.ContributorAddress] = append(result[record.ContributorAddress], &record)
	}

	return result, nil
}
