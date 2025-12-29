import { Card, Button, Empty, Tag, Spin, App, Skeleton } from 'antd';
import { useState, useEffect, useCallback } from 'react';
import { DeAIService } from '@@/github.com/kawai-network/veridium/internal/services';
import {
  History,
  Gift,
  ExternalLink,
  Repeat2,
  Coins,
} from 'lucide-react';
import { ActionIcon } from '@lobehub/ui';
import { Browser } from '@wailsio/runtime';
import { Flexbox } from 'react-layout-kit';
import { TokenUSDT } from '@web3icons/react';
import type { RewardsContentProps } from './types';

const RewardsContent = ({ styles, theme }: RewardsContentProps) => {
  const { message } = App.useApp();
  const [loading, setLoading] = useState(true);
  const [claimLoading, setClaimLoading] = useState<string | null>(null);
  const [rewards, setRewards] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);

  const loadRewards = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await DeAIService.GetClaimableRewards();
      setRewards(result);
    } catch (e: any) {
      console.error('Failed to load rewards:', e);
      setError(e.message || 'Failed to load rewards');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadRewards();
  }, [loadRewards]);

  const handleClaim = async (proof: any) => {
    const proofKey = `${proof.period_id}-${proof.reward_type}`;
    setClaimLoading(proofKey);
    try {
      let result;
      if (proof.reward_type === 'kawai') {
        result = await DeAIService.ClaimKawaiReward(
          proof.period_id,
          proof.index,
          proof.amount,
          proof.proof
        );
      } else {
        result = await DeAIService.ClaimUSDTReward(
          proof.period_id,
          proof.index,
          proof.amount,
          proof.proof
        );
      }

      if (result?.tx_hash) {
        message.success(
          <span>
            Claim submitted! Tx: {result.tx_hash.substring(0, 10)}...
            <a
              onClick={() => Browser.OpenURL(`https://testnet.monadexplorer.com/tx/${result.tx_hash}`)}
              style={{ marginLeft: 8, cursor: 'pointer' }}
            >
              View <ExternalLink size={12} style={{ verticalAlign: 'middle' }} />
            </a>
          </span>
        );
        // Refresh rewards after claim
        setTimeout(() => loadRewards(), 3000);
      }
    } catch (e: any) {
      console.error('Claim failed:', e);
      message.error(e.message || 'Claim failed');
    } finally {
      setClaimLoading(null);
    }
  };

  const formatDate = (dateStr: string) => {
    if (!dateStr) return '-';
    try {
      return new Date(dateStr).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
      });
    } catch {
      return '-';
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'unclaimed':
        return 'green';
      case 'pending':
        return 'orange';
      case 'confirmed':
        return 'blue';
      case 'failed':
        return 'red';
      default:
        return 'default';
    }
  };

  if (error) {
    return (
      <Flexbox style={{ maxWidth: 800 }} gap={20}>
        <div>
          <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>Rewards</h2>
          <span style={{ color: theme.colorTextSecondary, fontSize: 13 }}>Claim your KAWAI & USDT rewards</span>
        </div>
        <div className={styles.placeholderCard}>
          <Gift size={48} color={theme.colorTextQuaternary} style={{ marginBottom: 16 }} />
          <h3 style={{ margin: '0 0 8px', color: theme.colorError }}>Error Loading Rewards</h3>
          <p style={{ color: theme.colorTextSecondary, margin: '0 0 16px' }}>{error}</p>
          <Button onClick={loadRewards} icon={<Repeat2 size={16} />}>Retry</Button>
        </div>
      </Flexbox>
    );
  }

  const hasUnclaimedRewards = rewards?.unclaimed_proofs?.length > 0;
  const hasPendingRewards = rewards?.pending_proofs?.length > 0;

  return (
    <Flexbox style={{ maxWidth: 800 }} gap={20}>
      {/* Header */}
      <Flexbox horizontal justify="space-between" align="center">
        <div>
          <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>Rewards</h2>
          <span style={{ color: theme.colorTextSecondary, fontSize: 13 }}>
            Claim your KAWAI & USDT rewards from AI compute contributions
          </span>
        </div>
        <ActionIcon icon={Repeat2} onClick={loadRewards} title="Refresh" loading={loading} />
      </Flexbox>

      {/* Summary Cards */}
      <Flexbox horizontal gap={16} style={{ flexWrap: 'wrap' }}>
        {/* KAWAI Claimable */}
        <Card
          size="small"
          style={{
            flex: '1 1 200px',
            background: 'linear-gradient(135deg, #667eea20, #764ba220)',
            border: '1px solid #667eea40',
          }}
        >
          <Flexbox gap={8}>
            {loading ? (
              <Skeleton active paragraph={{ rows: 1, width: ['100%'] }} title={false} />
            ) : (
              <>
                <Flexbox horizontal align="center" gap={8}>
                  <Coins size={20} color="#667eea" />
                  <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>KAWAI Claimable</span>
                </Flexbox>
                <span style={{ fontSize: 24, fontWeight: 700 }}>
                  {rewards?.total_kawai_claimable_formatted || '0.0000'}
                </span>
                <span style={{ fontSize: 11, color: theme.colorTextTertiary }}>
                  Accumulating: {rewards?.current_kawai_accumulating || '0'}
                </span>
              </>
            )}
          </Flexbox>
        </Card>

        {/* USDT Claimable */}
        <Card
          size="small"
          style={{
            flex: '1 1 200px',
            background: 'linear-gradient(135deg, #26a17b20, #1a906520)',
            border: '1px solid #26a17b40',
          }}
        >
          <Flexbox gap={8}>
            {loading ? (
              <Skeleton active paragraph={{ rows: 1, width: ['100%'] }} title={false} />
            ) : (
              <>
                <Flexbox horizontal align="center" gap={8}>
                  <TokenUSDT size={20} variant="branded" />
                  <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>USDT Claimable</span>
                </Flexbox>
                <span style={{ fontSize: 24, fontWeight: 700 }}>
                  ${rewards?.total_usdt_claimable_formatted || '0.00'}
                </span>
                <span style={{ fontSize: 11, color: theme.colorTextTertiary }}>
                  Accumulating: ${rewards?.current_usdt_accumulating || '0'}
                </span>
              </>
            )}
          </Flexbox>
        </Card>
      </Flexbox>

      {/* Unclaimed Rewards */}
      <Card
        title={
          <Flexbox horizontal align="center" gap={8}>
            <Gift size={16} />
            <span>Unclaimed Rewards</span>
            {hasUnclaimedRewards && (
              <Tag color="green" style={{ marginLeft: 8 }}>
                {rewards.unclaimed_proofs.length} available
              </Tag>
            )}
          </Flexbox>
        }
        size="small"
      >
        {loading ? (
          <Skeleton active paragraph={{ rows: 3 }} />
        ) : !hasUnclaimedRewards ? (
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={
              <span style={{ color: theme.colorTextSecondary }}>
                No unclaimed rewards available.
                <br />
                Keep contributing to earn more!
              </span>
            }
          />
        ) : (
          <Flexbox gap={12}>
            {rewards.unclaimed_proofs.map((proof: any, idx: number) => {
              const proofKey = `${proof.period_id}-${proof.reward_type}`;
              const isLoading = claimLoading === proofKey;
              const isKawai = proof.reward_type === 'kawai';

              return (
                <div
                  key={proofKey}
                  style={{
                    padding: 16,
                    borderRadius: 12,
                    border: `1px solid ${theme.colorBorderSecondary}`,
                    background: theme.colorBgContainer,
                  }}
                >
                  <Flexbox horizontal justify="space-between" align="center">
                    <Flexbox gap={4}>
                      <Flexbox horizontal align="center" gap={8}>
                        {isKawai ? (
                          <Coins size={18} color="#667eea" />
                        ) : (
                          <TokenUSDT size={18} variant="branded" />
                        )}
                        <span style={{ fontWeight: 600, fontSize: 16 }}>
                          {proof.formatted || proof.amount} {isKawai ? 'KAWAI' : 'USDT'}
                        </span>
                        <Tag color={getStatusColor(proof.claim_status)}>
                          {proof.claim_status}
                        </Tag>
                      </Flexbox>
                      <span style={{ fontSize: 11, color: theme.colorTextTertiary }}>
                        Period: {formatDate(proof.created_at)} • Index #{proof.index}
                      </span>
                    </Flexbox>
                    <Button
                      type="primary"
                      onClick={() => handleClaim(proof)}
                      loading={isLoading}
                      icon={<Gift size={14} />}
                      style={{
                        background: isKawai
                          ? 'linear-gradient(135deg, #667eea, #764ba2)'
                          : 'linear-gradient(135deg, #26a17b, #1a9065)',
                        border: 'none',
                      }}
                    >
                      Claim
                    </Button>
                  </Flexbox>
                </div>
              );
            })}
          </Flexbox>
        )}
      </Card>

      {/* Pending Claims */}
      {hasPendingRewards && (
        <Card
          title={
            <Flexbox horizontal align="center" gap={8}>
              <History size={16} />
              <span>Pending Claims</span>
              <Tag color="orange">{rewards.pending_proofs.length} pending</Tag>
            </Flexbox>
          }
          size="small"
        >
          <Flexbox gap={8}>
            {rewards.pending_proofs.map((proof: any) => {
              const proofKey = `${proof.period_id}-${proof.reward_type}`;
              const isKawai = proof.reward_type === 'kawai';

              return (
                <Flexbox
                  key={proofKey}
                  horizontal
                  justify="space-between"
                  align="center"
                  style={{
                    padding: '8px 12px',
                    borderRadius: 8,
                    background: theme.colorFillTertiary,
                  }}
                >
                  <Flexbox horizontal align="center" gap={8}>
                    {isKawai ? (
                      <Coins size={16} color="#667eea" />
                    ) : (
                      <TokenUSDT size={16} variant="branded" />
                    )}
                    <span style={{ fontWeight: 500 }}>
                      {proof.formatted} {isKawai ? 'KAWAI' : 'USDT'}
                    </span>
                  </Flexbox>
                  <Flexbox horizontal align="center" gap={8}>
                    <Spin size="small" />
                    <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                      Confirming...
                    </span>
                    {proof.claim_tx_hash && (
                      <a
                        onClick={() => Browser.OpenURL(`https://testnet.monadexplorer.com/tx/${proof.claim_tx_hash}`)}
                        style={{ cursor: 'pointer' }}
                      >
                        <ExternalLink size={14} />
                      </a>
                    )}
                  </Flexbox>
                </Flexbox>
              );
            })}
          </Flexbox>
        </Card>
      )}

      {/* Info Card */}
      <Card size="small" style={{ background: theme.colorFillTertiary, border: 'none' }}>
        <Flexbox gap={8}>
          <span style={{ fontSize: 12, fontWeight: 600, color: theme.colorTextSecondary }}>
            How Rewards Work
          </span>
          <ul style={{ margin: 0, paddingLeft: 16, fontSize: 12, color: theme.colorTextTertiary }}>
            <li>Rewards accumulate off-chain as you contribute AI compute resources</li>
            <li>Weekly settlements generate Merkle proofs for gas-efficient claiming</li>
            <li>KAWAI tokens are minted during Phase 1 (Mining Era)</li>
            <li>USDT payments begin in Phase 2 after KAWAI max supply is reached</li>
          </ul>
        </Flexbox>
      </Card>
    </Flexbox>
  );
};

export default RewardsContent;

