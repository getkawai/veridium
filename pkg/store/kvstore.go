package store

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kawai-network/veridium/internal/constant"
	"golang.org/x/time/rate"
)

type Store interface {
	// Contributor operations
	SaveContributor(ctx context.Context, data *ContributorData) error
	GetContributor(ctx context.Context, address string) (*ContributorData, error)
	ListContributors(ctx context.Context) ([]*ContributorData, error)
	ListActiveContributors(ctx context.Context) ([]*ContributorData, error)
	ListContributorsWithBalance(ctx context.Context, rewardType string) ([]*ContributorData, error)
	GetOnlineContributors(ctx context.Context) ([]*ContributorData, error)
	UpdateHeartbeat(ctx context.Context, address string) error
	DeductSettledRewards(ctx context.Context, address string, rewardType string, amountToDeduct string) error
	RegisterContributor(ctx context.Context, address, endpointURL, hardwareSpecs string) (*ContributorData, error)
	SoftDeleteContributor(ctx context.Context, address string) error
	RestoreContributor(ctx context.Context, address string) error
	RecordJobReward(ctx context.Context, contributorAddress string, tokenUsage int64) error

	// Merkle proof operations (deprecated - use period-specific methods)
	SaveMerkleProof(ctx context.Context, address string, data *MerkleProofData) error
	GetMerkleProof(ctx context.Context, address string) (*MerkleProofData, error)

	// Period-specific Merkle proof operations
	SaveMerkleProofForPeriod(ctx context.Context, address string, periodID int64, data *MerkleProofData) error
	GetMerkleProofForPeriod(ctx context.Context, address string, periodID int64) (*MerkleProofData, error)
	ListMerkleProofs(ctx context.Context, address string) ([]*MerkleProofData, error)
	DeleteMerkleProof(ctx context.Context, address string, periodID int64) error

	// Claim status operations
	MarkClaimPending(ctx context.Context, address string, periodID int64, txHash string) error
	ConfirmClaim(ctx context.Context, address string, periodID int64) error
	MarkClaimFailed(ctx context.Context, address string, periodID int64, reason string) error
	RetryFailedClaim(ctx context.Context, address string, periodID int64) error
	GetPendingClaims(ctx context.Context) ([]*MerkleProofData, error)

	// Settlement operations
	GetSettlementSnapshots(ctx context.Context, rewardType string) ([]*SettlementSnapshot, error)
	PerformSettlement(ctx context.Context, periodID int64, merkleRoot string, rewardType string, proofs map[string]*MerkleProofData) (*SettlementPeriod, error)
	PerformSettlementWithConfig(ctx context.Context, periodID int64, merkleRoot string, rewardType string, proofs map[string]*MerkleProofData, config *SettlementConfig) (*SettlementPeriod, error)
	PerformSettlementParallel(ctx context.Context, periodID int64, merkleRoot string, rewardType string, proofs map[string]*MerkleProofData, workers int) (*SettlementPeriod, error)
	ResumeSettlement(ctx context.Context, periodID int64, proofs map[string]*MerkleProofData, config *SettlementConfig) (*SettlementPeriod, error)
	GetClaimableRewards(ctx context.Context, address string) (map[string]interface{}, error)

	// Settlement period operations
	SaveSettlementPeriod(ctx context.Context, period *SettlementPeriod) error
	GetSettlementPeriod(ctx context.Context, periodID int64) (*SettlementPeriod, error)
	ListSettlementPeriods(ctx context.Context) ([]*SettlementPeriod, error)

	// Admin operations
	EnsureAdminExists(ctx context.Context, adminAddress string) error

	// Marketplace operations
	StoreMarketplaceData(ctx context.Context, key string, data []byte) error
	StoreMarketplaceDataWithTTL(ctx context.Context, key string, data []byte, ttl int) error
	GetMarketplaceData(ctx context.Context, key string) ([]byte, error)
	DeleteMarketplaceData(ctx context.Context, key string) error
}

// SupplyQuerier defines interface for fetching token supply and max supply
type SupplyQuerier interface {
	GetTotalSupply(ctx context.Context) (*big.Int, error)
	GetMaxSupply(ctx context.Context) (*big.Int, error)
}

// KVStore implements Store interface with multiple namespaces
type KVStore struct {
	client    *cloudflare.API
	accountID string

	// Separate namespace IDs for different data types
	contributorsNamespaceID   string
	proofsNamespaceID         string
	settlementsNamespaceID    string
	authzNamespaceID          string // Reverse index: address -> apikey
	p2pMarketplaceNamespaceID string

	// ✅ Rate limiter for KV operations
	rateLimiter *rate.Limiter

	// Optional: Querier for Halving Logic
	supplyQuerier SupplyQuerier
}

// NewMultiNamespaceKVStore creates a new KVStore with separate namespaces
func NewMultiNamespaceKVStore() (*KVStore, error) {
	api, err := cloudflare.NewWithAPIToken(constant.GetCfApiToken())
	if err != nil {
		return nil, fmt.Errorf("failed to create cloudflare client: %w", err)
	}

	return &KVStore{
		client:                    api,
		accountID:                 constant.GetCfAccountId(),
		contributorsNamespaceID:   constant.GetCfKvContributorsNamespaceId(),
		proofsNamespaceID:         constant.GetCfKvProofsNamespaceId(),
		settlementsNamespaceID:    constant.GetCfKvSettlementsNamespaceId(),
		authzNamespaceID:          constant.GetCfKvAuthzNamespaceId(),
		p2pMarketplaceNamespaceID: constant.GetCfKvP2pMarketplaceNamespaceId(),
		rateLimiter:               rate.NewLimiter(rate.Limit(100), 200), // ✅ 100 ops/sec, burst 200
	}, nil
}

// SetSupplyQuerier injects the supply querier dependency
func (s *KVStore) SetSupplyQuerier(sq SupplyQuerier) {
	s.supplyQuerier = sq
}
