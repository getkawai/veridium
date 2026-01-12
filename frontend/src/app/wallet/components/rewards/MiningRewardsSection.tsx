import { Card, Button, Empty, Tag, Spin, App, Skeleton, Modal, Pagination, Table } from 'antd';
import { useState, useEffect, useCallback, useMemo } from 'react';
import { DeAIService, JarvisService } from '@@/github.com/kawai-network/veridium/internal/services';
import {
  History,
  Gift,
  ExternalLink,
  Repeat2,
  Coins,
  CheckCircle,
  Info
} from 'lucide-react';
import { Browser, Dialogs } from '@wailsio/runtime';
import { Flexbox } from 'react-layout-kit';
import { TokenUSDT } from '@web3icons/react';
import type { 
  NetworkInfo, 
  ClaimableReward, 
  ClaimableRewardsResponse 
} from '@@/github.com/kawai-network/veridium/internal/services/models';

interface MiningRewardsSectionProps {
  currentNetwork: NetworkInfo | null;
  theme: any;
  styles: any;
  onRefresh?: (refreshFn: () => void) => void;
}

export const MiningRewardsSection = ({ currentNetwork, theme, styles, onRefresh }: MiningRewardsSectionProps) => {
  const { message } = App.useApp();
  const [loading, setLoading] = useState(true);
  const [claimLoading, setClaimLoading] = useState<Set<string>>(new Set());
  const [rewards, setRewards] = useState<ClaimableRewardsResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Pagination & Modals
  const [page, setPage] = useState(1);
  const pageSize = 5;
  const [confirmModal, setConfirmModal] = useState<ClaimableReward | null>(null);
  const [gasEstimate, setGasEstimate] = useState<string | null>(null);
  const [estimateLoading, setEstimateLoading] = useState(false);
  const [isClaimAll, setIsClaimAll] = useState(false);

  // Helper functions
  const isValidTxHash = (hash: string): boolean => {
    return /^0x[a-fA-F0-9]{64}$/.test(hash);
  };

  const getExplorerUrl = (txHash: string): string => {
    if (!isValidTxHash(txHash)) {
      console.warn('Invalid transaction hash:', txHash);
      return '#'; // Return placeholder for invalid hashes
    }
    
    const baseUrl = currentNetwork?.explorerURL || 'https://testnet.monadexplorer.com';
    const cleanUrl = baseUrl.replace(/\/$/, '');
    return `${cleanUrl}/tx/${txHash}`;
  };

  const validateNetwork = (network: NetworkInfo | null): boolean => {
    if (!network) return false;
    if (!network.explorerURL) {
      console.warn('Network missing explorer URL:', network);
      return false;
    }
    return true;
  };

  const loadRewards = useCallback(async (showMessage = false) => {
    setLoading(true);
    setError(null);
    try {
      const result = await DeAIService.GetClaimableRewards();

      if (!result) {
        setError('No wallet connected or failed to load rewards. Please unlock your wallet.');
        setRewards(null);
        return;
      }

      setRewards(result);
      
      if (showMessage) {
        message.success('Rewards refreshed successfully');
      }
    } catch (e: any) {
      console.error('Failed to load rewards:', e);
      setError(e.message || 'Failed to load rewards');
      
      if (showMessage) {
        message.error('Failed to refresh rewards');
      }
    } finally {
      setLoading(false);
    }
  }, [message]);

  useEffect(() => {
    loadRewards();
  }, [loadRewards]);

  // Validate network configuration
  useEffect(() => {
    if (currentNetwork && !validateNetwork(currentNetwork)) {
      console.warn('Current network configuration issues detected');
    }
  }, [currentNetwork]);

  // Expose refresh function to parent
  useEffect(() => {
    if (onRefresh) {
      onRefresh(() => loadRewards(true));
    }
  }, [onRefresh, loadRewards]);

  const validUnclaimed = useMemo(() => {
    return (rewards?.unclaimed_proofs || []).filter((p): p is ClaimableReward => p !== null);
  }, [rewards]);

  const validPending = useMemo(() => {
    return (rewards?.pending_proofs || []).filter((p): p is ClaimableReward => p !== null);
  }, [rewards]);

  // Auto-refresh when there are pending claims
  useEffect(() => {
    if (validPending.length > 0) {
      const interval = setInterval(() => {
        loadRewards(false); // Silent refresh
      }, 10000); // Every 10 seconds

      return () => clearInterval(interval);
    }
  }, [validPending.length, loadRewards]);

  const paginatedUnclaimed = useMemo(() => {
    return validUnclaimed.slice((page - 1) * pageSize, page * pageSize);
  }, [validUnclaimed, page]);

  const estimateGas = async (proof: ClaimableReward) => {
    setEstimateLoading(true);
    try {
      if (currentNetwork) {
        const est = await JarvisService.EstimateGas(currentNetwork.id);
        if (est) {
          const totalGwei = est.maxGasPriceGwei * 300000;
          const totalEth = totalGwei / 1e9;
          setGasEstimate(`${totalEth.toFixed(6)} ${currentNetwork.nativeTokenSymbol}`);
        } else {
          setGasEstimate('Unknown');
        }
      } else {
        setGasEstimate('Unknown');
      }
    } catch (e) {
      console.error('Gas estimation failed:', e);
      setGasEstimate('Unknown');
    } finally {
      setEstimateLoading(false);
    }
  };

  const openConfirmModal = (proof: ClaimableReward) => {
    setConfirmModal(proof);
    setGasEstimate(null);
    estimateGas(proof);
  };

  const handleClaim = async (proof: ClaimableReward) => {
    const proofKey = `${proof.period_id}-${proof.reward_type}`;
    if (claimLoading.has(proofKey)) return;

    setClaimLoading(prev => new Set(prev).add(proofKey));

    try {
      let result: any; // Explicitly type as any for now
      if (proof.reward_type === 'kawai') {
        // Mining rewards use ClaimMiningReward with 9-field format
        result = await DeAIService.ClaimMiningReward(
          proof.period_id,
          proof.contributor_amount || proof.amount, // Contributor amount
          proof.developer_amount || "0",            // Developer amount
          proof.user_amount || "0",                 // User amount
          proof.affiliator_amount || "0",           // Affiliator amount
          proof.developer_address || "0x0000000000000000000000000000000000000000", // Developer address
          proof.user_address || "0x0000000000000000000000000000000000000000",      // User address
          proof.affiliator_address || "0x0000000000000000000000000000000000000000", // Affiliator address
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
        const explorerUrl = getExplorerUrl(result.tx_hash);

        message.success(
          <span>
            Claim submitted! Tx: {result.tx_hash.substring(0, 10)}...
            <a
              onClick={() => Browser.OpenURL(explorerUrl)}
              style={{ marginLeft: 8, cursor: 'pointer' }}
            >
              View <ExternalLink size={12} style={{ verticalAlign: 'middle' }} />
            </a>
          </span>
        );
        setTimeout(() => loadRewards(), 3000);
      }
    } catch (e: any) {
      console.error('Claim failed:', e);
      
      // Always show error dialog (better visibility than toast)
      const errorMsg = e.message || 'Claim failed';
      const isInsufficientFunds = errorMsg.toLowerCase().includes('insufficient funds');
      
      if (isInsufficientFunds) {
        await Dialogs.Error({
          Title: 'Insufficient Funds for Gas',
          Message: 'You need MON tokens to pay for gas fees.\n\n' +
                   'After receiving MON, try claiming again.'
        });
      } else {
        await Dialogs.Error({
          Title: 'Claim Transaction Failed',
          Message: errorMsg
        });
      }
      
      setTimeout(() => loadRewards(), 1000);
    } finally {
      setClaimLoading(prev => {
        const next = new Set(prev);
        next.delete(proofKey);
        return next;
      });
      setConfirmModal(null);
    }
  };

  const handleClaimAll = async () => {
    if (validUnclaimed.length === 0) return;
    setIsClaimAll(true);

    for (const proof of validUnclaimed) {
      const proofKey = `${proof.period_id}-${proof.reward_type}`;
      if (claimLoading.has(proofKey)) continue;

      await handleClaim(proof);
      await new Promise(r => setTimeout(r, 2000));
    }

    setIsClaimAll(false);
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

  const formatAccumulating = (rawAmount: string, decimals: number) => {
    if (!rawAmount || rawAmount === '0') return '0';
    try {
      const amount = BigInt(rawAmount);
      const divisor = BigInt(10 ** decimals);
      const wholePart = amount / divisor;
      const fractionalPart = amount % divisor;
      
      const fractionalStr = fractionalPart.toString().padStart(decimals, '0');
      const precision = decimals === 18 ? 4 : 2;
      const trimmedFractional = fractionalStr.slice(0, precision);
      
      return `${wholePart}.${trimmedFractional}`;
    } catch (e) {
      console.error('Failed to format accumulating amount:', e);
      return '0';
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
      <Flexbox style={{ width: '100%' }} gap={20}>
        <div className={styles.placeholderCard}>
          <Gift size={48} color={theme.colorTextQuaternary} style={{ marginBottom: 16 }} />
          <h3 style={{ margin: '0 0 8px', color: theme.colorError }}>Error Loading Rewards</h3>
          <p style={{ color: theme.colorTextSecondary, margin: '0 0 16px' }}>{error}</p>
          <Button onClick={() => loadRewards(true)} icon={<Repeat2 size={16} />}>Retry</Button>
        </div>
      </Flexbox>
    );
  }

  const hasUnclaimedRewards = validUnclaimed.length > 0;
  const hasPendingRewards = validPending.length > 0;

  return (
    <Flexbox style={{ width: '100%' }} gap={20}>
      {/* Summary Cards */}
      <Flexbox horizontal gap={16} style={{ flexWrap: 'wrap' }}>
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
                  Accumulating: {formatAccumulating(rewards?.current_kawai_accumulating || '0', 18)} KAWAI
                </span>
              </>
            )}
          </Flexbox>
        </Card>

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
                  Accumulating: ${formatAccumulating(rewards?.current_usdt_accumulating || '0', 6)} USDT
                </span>
              </>
            )}
          </Flexbox>
        </Card>
      </Flexbox>

      {/* Unclaimed Rewards */}
      <Card
        title={
          <Flexbox horizontal align="center" gap={8} justify="space-between" width="100%">
            <Flexbox horizontal align="center" gap={8}>
              <Gift size={16} />
              <span>Unclaimed Rewards</span>
              {hasUnclaimedRewards && (
                <Tag color="green" style={{ marginLeft: 8 }}>
                  {validUnclaimed.length} available
                </Tag>
              )}
            </Flexbox>
            {hasUnclaimedRewards && (
              <Button
                size="small"
                type="primary"
                ghost
                onClick={handleClaimAll}
                loading={isClaimAll}
              >
                Claim All
              </Button>
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
            {paginatedUnclaimed.map((proof: ClaimableReward) => {
              const proofKey = `${proof.period_id}-${proof.reward_type}`;
              const isLoading = claimLoading.has(proofKey);
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
                      onClick={() => openConfirmModal(proof)}
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
            <Flexbox justify="center" style={{ marginTop: 16 }}>
              <Pagination
                current={page}
                total={validUnclaimed.length}
                pageSize={pageSize}
                onChange={setPage}
                size="small"
                hideOnSinglePage
              />
            </Flexbox>
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
              <Tag color="orange">{validPending.length} pending</Tag>
            </Flexbox>
          }
          size="small"
        >
          <Flexbox gap={8}>
            {validPending.map((proof: ClaimableReward) => {
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
                        onClick={() => {
                          const explorerUrl = getExplorerUrl(proof.claim_tx_hash);
                          Browser.OpenURL(explorerUrl);
                        }}
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

      {/* Recent Activity */}
      <Card title="Recent Mining Activity" size="small">
        {loading ? (
          <Flexbox gap={12}>
            <Skeleton active paragraph={{ rows: 1, width: ['100%', '80%', '60%', '90%', '70%'] }} title={false} />
            <Skeleton active paragraph={{ rows: 1, width: ['90%', '70%', '80%', '60%', '85%'] }} title={false} />
            <Skeleton active paragraph={{ rows: 1, width: ['80%', '90%', '70%', '85%', '75%'] }} title={false} />
          </Flexbox>
        ) : (
          (() => {
            // Create recent activity from confirmed claims
            const getRecentActivity = () => {
              // Use confirmed_proofs directly from backend
              const confirmedClaims = (rewards?.confirmed_proofs || [])
                .filter((p): p is ClaimableReward => p !== null)
                .sort((a, b) => {
                  const aTime = a.claimed_at ? new Date(a.claimed_at).getTime() : new Date(a.created_at).getTime();
                  const bTime = b.claimed_at ? new Date(b.claimed_at).getTime() : new Date(b.created_at).getTime();
                  return bTime - aTime;
                })
                .slice(0, 10) // Show last 10 claims
                .map(proof => ({
                  key: `${proof.period_id}-${proof.reward_type}-${proof.claim_tx_hash}`,
                  txHash: proof.claim_tx_hash!,
                  txType: 'Mining Claim',
                  createdAt: proof.claimed_at || proof.created_at,
                  status: 'confirmed',
                  amount: proof.formatted,
                  rewardType: proof.reward_type
                }));

              return confirmedClaims;
            };

            const recentActivity = getRecentActivity();

          return (
            <Table
              dataSource={recentActivity}
              rowKey="key"
              pagination={false}
              size="small"
              locale={{
                emptyText: recentActivity.length === 0 ? (
                  <Empty
                    image={Empty.PRESENTED_IMAGE_SIMPLE}
                    description={
                      <span style={{ color: theme.colorTextSecondary }}>
                        No mining claims yet.
                        <br />
                        {validUnclaimed.length > 0 
                          ? "Claim your rewards above to see activity here."
                          : "Keep contributing to earn mining rewards!"
                        }
                      </span>
                    }
                  />
                ) : undefined
              }}
              columns={[
                {
                  title: 'Type',
                  dataIndex: 'txType',
                  key: 'type',
                  render: (t: string) => <Tag color="blue">{t}</Tag>
                },
                {
                  title: 'Amount',
                  key: 'amount',
                  render: (record: any) => (
                    <Flexbox horizontal align="center" gap={4}>
                      <Coins size={12} color="#667eea" />
                      <span style={{ fontSize: 12, fontWeight: 500 }}>
                        {record.amount} {record.rewardType === 'kawai' ? 'KAWAI' : 'USDT'}
                      </span>
                    </Flexbox>
                  )
                },
                {
                  title: 'Hash',
                  dataIndex: 'txHash',
                  render: (h: string) => (
                    <a
                      onClick={() => {
                        const explorerUrl = getExplorerUrl(h);
                        Browser.OpenURL(explorerUrl);
                      }}
                      style={{ cursor: 'pointer', fontFamily: 'monospace', fontSize: 12 }}
                    >
                      {h.substring(0, 10)}... <ExternalLink size={10} style={{ verticalAlign: 'middle' }} />
                    </a>
                  )
                },
                {
                  title: 'Date',
                  dataIndex: 'createdAt',
                  render: (d: string) => <span style={{ fontSize: 12 }}>{formatDate(d)}</span>
                },
                {
                  title: 'Status',
                  dataIndex: 'status',
                  key: 'status',
                  render: (status: string) => (
                    <Flexbox horizontal align="center" gap={4}>
                      {status === 'confirmed' ? (
                        <>
                          <CheckCircle size={14} color="#22c55e" />
                          <span style={{ fontSize: 12, textTransform: 'capitalize' }}>Confirmed</span>
                        </>
                      ) : (
                        <>
                          <span style={{ 
                            width: 14, 
                            height: 14, 
                            borderRadius: '50%', 
                            backgroundColor: getStatusColor(status),
                            display: 'inline-block'
                          }} />
                          <span style={{ fontSize: 12, textTransform: 'capitalize' }}>{status}</span>
                        </>
                      )}
                    </Flexbox>
                  )
                }
              ]}
            />
          );
        })()
        )}
      </Card>

      {/* Info Card */}
      <Card size="small" style={{ background: theme.colorFillTertiary, border: 'none' }}>
        <Flexbox gap={8}>
          <span style={{ fontSize: 12, fontWeight: 600, color: theme.colorTextSecondary }}>
            How Mining Rewards Work
          </span>
          <ul style={{ margin: 0, paddingLeft: 16, fontSize: 12, color: theme.colorTextTertiary }}>
            <li>Rewards accumulate off-chain as you contribute AI compute resources</li>
            <li>Weekly settlements generate Merkle proofs for gas-efficient claiming</li>
            <li>KAWAI tokens are minted during Phase 1 (Mining Era)</li>
            <li>USDT payments begin in Phase 2 after KAWAI max supply is reached</li>
          </ul>
        </Flexbox>
      </Card>

      {/* Confirmation Modal */}
      <Modal
        title={
          <Flexbox horizontal align="center" gap={8}>
            <Info size={18} />
            <span>Confirm Claim</span>
          </Flexbox>
        }
        open={!!confirmModal}
        onOk={() => confirmModal && handleClaim(confirmModal)}
        onCancel={() => setConfirmModal(null)}
        okText="Confirm & Claim"
        confirmLoading={confirmModal ? claimLoading.has(`${confirmModal.period_id}-${confirmModal.reward_type}`) : false}
      >
        {confirmModal && (
          <Flexbox gap={16} padding={8}>
            <div style={{ padding: 16, background: theme.colorFillSecondary, borderRadius: 8 }}>
              <Flexbox horizontal justify="space-between" align="center">
                <span>Reward Amount:</span>
                <span style={{ fontWeight: 700, fontSize: 18 }}>
                  {confirmModal.formatted} {confirmModal.reward_type === 'kawai' ? 'KAWAI' : 'USDT'}
                </span>
              </Flexbox>
              <Flexbox horizontal justify="space-between" style={{ marginTop: 8 }}>
                <span style={{ color: theme.colorTextSecondary }}>Period:</span>
                <span>#{confirmModal.period_id}</span>
              </Flexbox>
            </div>

            <Flexbox gap={4}>
              <span style={{ fontWeight: 600 }}>Estimated Gas Fee:</span>
              {estimateLoading ? (
                <Spin size="small" />
              ) : (
                <span style={{ fontFamily: 'monospace' }}>
                  {gasEstimate || 'Calculating...'}
                </span>
              )}
              <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                Network: {currentNetwork?.name}
              </span>
            </Flexbox>
          </Flexbox>
        )}
      </Modal>
    </Flexbox>
  );
};

