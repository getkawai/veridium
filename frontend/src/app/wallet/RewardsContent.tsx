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
  Clock,
  Info
} from 'lucide-react';
import { ActionIcon } from '@lobehub/ui';
import { Browser } from '@wailsio/runtime';
import { Flexbox } from 'react-layout-kit';
import { TokenUSDT } from '@web3icons/react';
import type { RewardsContentProps, ClaimableReward, ClaimableRewardsResponse } from './types';

const RewardsContent = ({ styles, theme, currentNetwork, transactions }: RewardsContentProps) => {
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

  const loadRewards = useCallback(async (showMessage = false) => {
    setLoading(true);
    setError(null);
    try {
      const result = await DeAIService.GetClaimableRewards();

      // Fix Bug #1: Missing Nil Check
      if (!result) {
        // If wallet is locked or service fails silently returning nil
        setError('No wallet connected or failed to load rewards. Please unlock your wallet.');
        setRewards(null);
        return;
      }

      setRewards(result);
      
      // Show success message on manual refresh
      if (showMessage) {
        message.success('Rewards refreshed successfully');
      }
    } catch (e: any) {
      console.error('Failed to load rewards:', e);
      setError(e.message || 'Failed to load rewards');
      
      // Show error message on manual refresh
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

  // Derived state
  const validUnclaimed = useMemo(() => {
    return (rewards?.unclaimed_proofs || []).filter((p): p is ClaimableReward => p !== null);
  }, [rewards]);

  const validPending = useMemo(() => {
    return (rewards?.pending_proofs || []).filter((p): p is ClaimableReward => p !== null);
  }, [rewards]);

  const paginatedUnclaimed = useMemo(() => {
    return validUnclaimed.slice((page - 1) * pageSize, page * pageSize);
  }, [validUnclaimed, page]);

  const estimateGas = async (proof: ClaimableReward) => {
    setEstimateLoading(true);
    try {
      // Use generic network gas estimate since specific claim estimate is not available
      if (currentNetwork) {
        const est = await JarvisService.EstimateGas(currentNetwork.id);
        if (est) {
          // Estimate: 300,000 gas * gasPrice (in Gwei)
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
    setGasEstimate(null); // Reset
    estimateGas(proof);
  };

  const handleClaim = async (proof: ClaimableReward) => {
    const proofKey = `${proof.period_id}-${proof.reward_type}`;
    if (claimLoading.has(proofKey)) return;

    setClaimLoading(prev => new Set(prev).add(proofKey));

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
        const explorerUrl = currentNetwork?.explorerURL || 'https://testnet.monadexplorer.com'; // Fix explorerURL typo

        message.success(
          <span>
            Claim submitted! Tx: {result.tx_hash.substring(0, 10)}...
            <a
              onClick={() => Browser.OpenURL(`${explorerUrl}/tx/${result.tx_hash}`)}
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
      message.error(e.message || 'Claim failed');
      setTimeout(() => loadRewards(), 1000);
    } finally {
      setClaimLoading(prev => {
        const next = new Set(prev);
        next.delete(proofKey);
        return next;
      });
      setConfirmModal(null); // Close modal if open
    }
  };

  const handleClaimAll = async () => {
    if (validUnclaimed.length === 0) return;
    setIsClaimAll(true);

    // Process sequentially to avoid nonce issues
    for (const proof of validUnclaimed) {
      // Skip if already claiming
      const proofKey = `${proof.period_id}-${proof.reward_type}`;
      if (claimLoading.has(proofKey)) continue;

      await handleClaim(proof);
      // Small delay between claims
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
          <Button onClick={() => loadRewards(true)} icon={<Repeat2 size={16} />}>Retry</Button>
        </div>
      </Flexbox>
    );
  }

  const hasUnclaimedRewards = validUnclaimed.length > 0;
  const hasPendingRewards = validPending.length > 0;

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
        <ActionIcon 
          icon={Repeat2} 
          onClick={() => loadRewards(true)} 
          title="Refresh rewards" 
          loading={loading}
        />
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
            {paginatedUnclaimed.map((proof: ClaimableReward, idx: number) => {
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
                          const explorerUrl = currentNetwork?.explorerURL || 'https://testnet.monadexplorer.com';
                          Browser.OpenURL(`${explorerUrl}/tx/${proof.claim_tx_hash}`);
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

      {/* Recent Activity / History */}
      <Card title="Recent Activity" size="small">
        <Table
          dataSource={transactions.slice(0, 5)}
          rowKey="txHash"
          pagination={false}
          size="small"
          columns={[
            {
              title: 'Type',
              dataIndex: 'txType',
              key: 'type',
              render: (t: string) => <Tag>{t || 'TX'}</Tag>
            },
            {
              title: 'Hash',
              dataIndex: 'txHash',
              render: (h: string) => <span style={{ fontFamily: 'monospace' }}>{h.substring(0, 10)}...</span>
            },
            {
              title: 'Date',
              dataIndex: 'createdAt',
              render: (d: number | string) => <span style={{ fontSize: 12 }}>{formatDate(d.toString())}</span>
            },
            {
              title: 'Status',
              dataIndex: 'status',
              key: 'status',
              render: (s: string) => (
                <Flexbox horizontal align="center" gap={4}>
                  {s === 'confirmed' ? <CheckCircle size={14} color="#22c55e" /> : <Clock size={14} color="#f59e0b" />}
                  <span style={{ fontSize: 12, textTransform: 'capitalize' }}>{s}</span>
                </Flexbox>
              )
            }
          ]}
        />
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

