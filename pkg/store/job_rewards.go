package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

// JobRewardRecord stores detailed reward split for a single job
// This is used to generate 9-field Merkle leaves for MiningRewardDistributor
type JobRewardRecord struct {
	Timestamp          time.Time `json:"timestamp"`
	ContributorAddress string    `json:"contributor_address"`
	UserAddress        string    `json:"user_address"`
	ReferrerAddress    string    `json:"referrer_address"`  // Empty if non-referral
	DeveloperAddress   string    `json:"developer_address"` // From GetRandomTreasuryAddress()

	ContributorAmount string `json:"contributor_amount"` // Contributor reward amount (85% or 90% of total)
	DeveloperAmount   string `json:"developer_amount"`   // Developer reward amount (5% of total)
	UserAmount        string `json:"user_amount"`        // User reward amount (5% of total)
	AffiliatorAmount  string `json:"affiliator_amount"`  // Affiliator reward amount (5% of total or 0)

	TokenUsage  int64  `json:"token_usage"`
	RewardType  string `json:"reward_type"` // "kawai" or "usdt"
	HasReferrer bool   `json:"has_referrer"`

	// For tracking settlement
	SettledPeriodID int64 `json:"settled_period_id,omitempty"`
	IsSettled       bool  `json:"is_settled,omitempty"`
}

// SaveJobReward stores a job reward record in KV
// Key format: job_rewards:{contributor}:{timestamp_unix}
func (s *KVStore) SaveJobReward(ctx context.Context, record *JobRewardRecord) error {
	key := fmt.Sprintf("job_rewards:%s:%d", record.ContributorAddress, record.Timestamp.Unix())

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal job reward: %w", err)
	}

	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.contributorsNamespaceID,
		Key:         key,
		Value:       data,
	})

	if err != nil {
		return fmt.Errorf("failed to save job reward: %w", err)
	}

	return nil
}

// GetJobRewardsSinceLastSettlement retrieves all unsettled job rewards for a contributor
func (s *KVStore) GetJobRewardsSinceLastSettlement(ctx context.Context, contributorAddress string, rewardType string) ([]*JobRewardRecord, error) {
	// List all keys for this contributor
	prefix := fmt.Sprintf("job_rewards:%s:", contributorAddress)

	resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.ListWorkersKVsParams{
		NamespaceID: s.contributorsNamespaceID,
		Prefix:      prefix,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list job rewards: %w", err)
	}

	var records []*JobRewardRecord

	for _, key := range resp.Result {
		// Get the record
		value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
			NamespaceID: s.contributorsNamespaceID,
			Key:         key.Name,
		})
		if err != nil {
			slog.Warn("Failed to get job reward", "key", key.Name, "error", err)
			continue
		}

		var record JobRewardRecord
		if err := json.Unmarshal(value, &record); err != nil {
			slog.Warn("Failed to unmarshal job reward", "key", key.Name, "error", err)
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

	resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.ListWorkersKVsParams{
		NamespaceID: s.contributorsNamespaceID,
		Prefix:      prefix,
	})
	if err != nil {
		return fmt.Errorf("failed to list job rewards: %w", err)
	}

	for _, key := range resp.Result {
		// Get the record
		value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
			NamespaceID: s.contributorsNamespaceID,
			Key:         key.Name,
		})
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

			_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
				NamespaceID: s.contributorsNamespaceID,
				Key:         key.Name,
				Value:       data,
			})
			if err != nil {
				slog.Warn("Failed to mark job reward as settled", "key", key.Name, "error", err)
			}
		}
	}

	return nil
}

// GetAllUnsettledJobRewards retrieves all unsettled job rewards across all contributors
// Used for settlement generation
func (s *KVStore) GetAllUnsettledJobRewards(ctx context.Context, rewardType string) (map[string][]*JobRewardRecord, error) {
	// List all keys in contributors namespace
	resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.ListWorkersKVsParams{
		NamespaceID: s.contributorsNamespaceID,
		Limit:       1000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	result := make(map[string][]*JobRewardRecord)

	// Scan all job_rewards keys
	for _, key := range resp.Result {
		// Only process job_rewards keys
		if !strings.HasPrefix(key.Name, "job_rewards:") {
			continue
		}

		// Get the record
		value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
			NamespaceID: s.contributorsNamespaceID,
			Key:         key.Name,
		})
		if err != nil {
			slog.Warn("Failed to get job reward", "key", key.Name, "error", err)
			continue
		}

		var record JobRewardRecord
		if err := json.Unmarshal(value, &record); err != nil {
			slog.Warn("Failed to unmarshal job reward", "key", key.Name, "error", err)
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
