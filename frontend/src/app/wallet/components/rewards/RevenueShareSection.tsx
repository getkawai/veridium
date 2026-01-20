import { Card, Button, Statistic, Row, Col, App, Skeleton, Tag, Table, Progress, Modal } from 'antd';
import { useState, useEffect, useCallback } from 'react';
import { DeAIService, JarvisService } from '@@/github.com/kawai-network/veridium/internal/services';
import { Flexbox } from 'react-layout-kit';
import { TrendingUp, DollarSign, Gift, Info, Percent, Users, ExternalLink, PieChart } from 'lucide-react';
import { useUserStore } from '@/store/user';
import { Browser } from '@wailsio/runtime';
import { TokenUSDT } from '@web3icons/react';
import type { NetworkInfo, ClaimableReward, RevenueShareStatsResponse } from '@@/github.com/kawai-network/veridium/internal/services/models';

interface RevenueShareRecord {
  period_id: number;
  index: number;
  amount: string;
  formatted: string;
  proof: string[];
  claim_status: string;
  claim_tx_hash?: string;
  created_at: string;
  kawai_balance: string;
  total_supply: string;
  share_percentage: string;
}

interface RevenueShareStats {
  total_earned: string;
  total_claimable: string;
  pending_claims: number;
  current_kawai_balance: string;
  current_share_percentage: string;
  estimated_weekly_usdt: string;
  unclaimed_records: RevenueShareRecord[];
}

interface RevenueShareSectionProps {
  currentNetwork: NetworkInfo | null;
  theme: any;
  styles: any;
  onRefresh?: (refreshFn: () => void) => void;
}

export const RevenueShareSection = ({ currentNetwork, theme, styles, onRefresh }: RevenueShareSectionProps) => {
  const { message } = App.useApp();
  const userAddress = useUserStore((s) => s.walletAddress);
  
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<RevenueShareStats | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [claimLoading, setClaimLoading] = useState<Set<number>>(new Set());
  const [confirmModal, setConfirmModal] = useState<RevenueShareRecord | null>(null);
  const [gasEstimate, setGasEstimate] = useState<string | null>(null);
  const [estimateLoading, setEstimateLoading] = useState(false);

  const [blockchainError, setBlockchainError] = useState(false);

  const loadRevenueShareStats = useCallback(async (address: string, showMessage = false) => {
    if (!address) {
      setError('No wallet connected');
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);
    setBlockchainError(false);
    
    try {
      // Get claimable rewards (USDT type only)
      const rewards = await DeAIService.GetClaimableRewards();
      
      if (!rewards) {
        setError('Failed to load revenue share data. Please try again.');
        setLoading(false);
        return;
      }

      // Filter for USDT rewards only (revenue distribution)
      const usdtRewards = (rewards.unclaimed_proofs || [])
        .filter((p): p is ClaimableReward => p !== null && p.reward_type === 'usdt');

      // Get real blockchain data
      let revenueStats: RevenueShareStatsResponse | null = null;
      try {
        const stats = await DeAIService.GetRevenueShareStats();
        if (stats) {
          revenueStats = stats as RevenueShareStatsResponse;
        }
      } catch (e) {
        console.warn('Failed to get revenue share stats from blockchain:', e);
        setBlockchainError(true);
        if (showMessage) {
          message.warning('Unable to fetch KAWAI balance from blockchain');
        }
      }

      const revenueShareStats: RevenueShareStats = {
        total_earned: rewards.total_usdt_claimable_formatted || '0.00',
        total_claimable: rewards.total_usdt_claimable_formatted || '0.00',
        pending_claims: usdtRewards.length,
        current_kawai_balance: revenueStats?.kawai_balance_formatted || '0',
        current_share_percentage: revenueStats?.share_percentage || '0',
        estimated_weekly_usdt: '0.00', // TODO: Calculate from historical platform revenue
        unclaimed_records: usdtRewards.map((r) => ({
          period_id: r.period_id,
          index: r.index,
          amount: r.amount,
          formatted: r.formatted,
          proof: r.proof,
          claim_status: r.claim_status,
          claim_tx_hash: r.claim_tx_hash,
          created_at: r.created_at,
          kawai_balance: revenueStats?.kawai_balance || '0',
          total_supply: revenueStats?.total_supply || '0',
          share_percentage: revenueStats?.share_percentage || '0',
        })),
      };

      setStats(revenueShareStats);

      if (showMessage) {
        message.success('Revenue share data refreshed');
      }
    } catch (e: any) {
      console.error('Failed to load revenue share stats:', e);
      setError(e.message || 'Failed to load revenue share data');
      
      if (showMessage) {
        message.error('Failed to refresh revenue share data');
      }
    } finally {
      setLoading(false);
    }
  }, [message]);

  useEffect(() => {
    if (userAddress) {
      loadRevenueShareStats(userAddress);
    }
  }, [userAddress, loadRevenueShareStats]);

  // Expose refresh function to parent
  useEffect(() => {
    if (onRefresh && userAddress) {
      onRefresh(() => loadRevenueShareStats(userAddress, true));
    }
  }, [onRefresh, loadRevenueShareStats, userAddress]);

  const estimateGas = async () => {
    setEstimateLoading(true);
    try {
      if (currentNetwork) {
        const est = await JarvisService.EstimateGas(currentNetwork.id);
        if (est) {
          const totalGwei = est.maxGasPriceGwei * 150000; // USDT claim ~150k gas
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

  const openConfirmModal = (record: RevenueShareRecord) => {
    setConfirmModal(record);
    setGasEstimate(null);
    estimateGas();
  };

  const handleClaimRevenue = async (record: RevenueShareRecord) => {
    if (claimLoading.has(record.period_id)) return;
    setClaimLoading(prev => new Set(prev).add(record.period_id));

    try {
      const result = await DeAIService.ClaimUSDTReward(
        record.period_id,
        record.index,
        record.amount,
        record.proof
      );

      if (result?.tx_hash) {
        const explorerUrl = currentNetwork?.explorerURL || 'https://testnet.monadexplorer.com';
        message.success(
          <span>
            Revenue claimed! Tx: {result.tx_hash.substring(0, 10)}...
            <a
              onClick={() => Browser.OpenURL(`${explorerUrl}/tx/${result.tx_hash}`)}
              style={{ marginLeft: 8, cursor: 'pointer' }}
            >
              View <ExternalLink size={12} style={{ verticalAlign: 'middle' }} />
            </a>
          </span>
        );
        setTimeout(() => userAddress && loadRevenueShareStats(userAddress, true), 3000);
      }
    } catch (e: any) {
      console.error('Revenue claim failed:', e);
      message.error(e.message || 'Claim failed');
    } finally {
      setClaimLoading(prev => {
        const next = new Set(prev);
        next.delete(record.period_id);
        return next;
      });
      setConfirmModal(null);
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

  if (error) {
    return (
      <Flexbox style={{ width: '100%' }} gap={20}>
        <div className={styles.placeholderCard}>
          <TrendingUp size={48} color={theme.colorTextQuaternary} style={{ marginBottom: 16 }} />
          <h3 style={{ margin: '0 0 8px', color: theme.colorError }}>Error Loading Revenue Data</h3>
          <p style={{ color: theme.colorTextSecondary, margin: '0 0 16px' }}>{error}</p>
          <Button onClick={() => userAddress && loadRevenueShareStats(userAddress, true)} icon={<Info size={16} />}>
            Retry
          </Button>
        </div>
      </Flexbox>
    );
  }

  const sharePercentage = parseFloat(stats?.current_share_percentage || '0');

  return (
    <Flexbox style={{ width: '100%' }} gap={20}>
      {/* Blockchain Error Warning */}
      {blockchainError && (
        <Card
          style={{
            background: 'linear-gradient(135deg, #ef444410, #dc262610)',
            border: '1px solid #ef444460',
          }}
        >
          <Flexbox horizontal gap={12} align="center">
            <Info size={20} color="#ef4444" />
            <Flexbox gap={4} flex={1}>
              <span style={{ fontWeight: 600, color: theme.colorError }}>
                Blockchain Data Unavailable
              </span>
              <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
                Unable to fetch your KAWAI balance from the blockchain. Showing 0 as fallback. 
                Your actual holdings are safe - this is just a display issue.
              </span>
            </Flexbox>
          </Flexbox>
        </Card>
      )}

      {/* Phase 1 Notice */}
      <Card
        style={{
          background: 'linear-gradient(135deg, #f59e0b10, #d9770610)',
          border: '1px solid #f59e0b60',
        }}
      >
        <Flexbox horizontal gap={12} align="center">
          <Info size={20} color="#f59e0b" />
          <Flexbox gap={4} flex={1}>
            <span style={{ fontWeight: 600, color: theme.colorWarning }}>
              Phase 1: Revenue Sharing Preview
            </span>
            <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
              Full revenue sharing activates in Phase 2 after 1B KAWAI supply is reached. 
              Your KAWAI balance and share percentage are shown in real-time from the blockchain. 
              Estimated earnings will be calculated once platform revenue data is available.
            </span>
          </Flexbox>
        </Flexbox>
      </Card>

      {/* Hero Banner */}
      <Card
        style={{
          background: 'linear-gradient(135deg, #26a17b20, #1a906520)',
          border: '1px solid #26a17b40',
        }}
      >
        <Flexbox gap={16}>
          <Flexbox horizontal align="center" gap={12}>
            <div
              style={{
                padding: 16,
                borderRadius: 12,
                background: 'linear-gradient(135deg, #26a17b, #1a9065)',
              }}
            >
              <PieChart size={32} color="white" />
            </div>
            <Flexbox gap={4}>
              <h2 style={{ margin: 0, fontSize: 24, fontWeight: 700 }}>Hold-to-Earn Revenue Share</h2>
              <span style={{ fontSize: 14, color: theme.colorTextSecondary }}>
                Earn 100% of platform profit (USDT) proportional to your KAWAI holdings
              </span>
            </Flexbox>
          </Flexbox>
        </Flexbox>
      </Card>

      {/* Summary Cards */}
      <Row gutter={16}>
        <Col xs={24} sm={8}>
          <Card size="small" style={{ background: 'linear-gradient(135deg, #26a17b20, #1a906520)', border: '1px solid #26a17b40' }}>
            {loading ? (
              <Skeleton active paragraph={{ rows: 1 }} />
            ) : (
              <Statistic
                title="Total Earned"
                value={stats?.total_earned || '0.00'}
                prefix={<DollarSign size={20} color="#26a17b" />}
                suffix="USDT"
                precision={2}
                valueStyle={{ color: '#26a17b', fontWeight: 700 }}
              />
            )}
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card size="small" style={{ background: 'linear-gradient(135deg, #22c55e20, #16a34a20)', border: '1px solid #22c55e40' }}>
            {loading ? (
              <Skeleton active paragraph={{ rows: 1 }} />
            ) : (
              <Statistic
                title="Claimable Now"
                value={stats?.total_claimable || '0.00'}
                prefix={<Gift size={20} color="#22c55e" />}
                suffix="USDT"
                precision={2}
                valueStyle={{ color: '#22c55e', fontWeight: 700 }}
              />
            )}
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card size="small" style={{ background: 'linear-gradient(135deg, #3b82f620, #2563eb20)', border: '1px solid #3b82f640' }}>
            {loading ? (
              <Skeleton active paragraph={{ rows: 1 }} />
            ) : (
              <Statistic
                title="Est. Weekly"
                value={stats?.estimated_weekly_usdt || '0.00'}
                prefix={<TrendingUp size={20} color="#3b82f6" />}
                suffix="USDT"
                precision={2}
                valueStyle={{ color: '#3b82f6', fontWeight: 700 }}
              />
            )}
          </Card>
        </Col>
      </Row>

      {/* Your Share Card */}
      <Card
        title={
          <Flexbox horizontal align="center" gap={8}>
            <Percent size={16} />
            <span>Your Share of Revenue</span>
          </Flexbox>
        }
        size="small"
      >
        {loading ? (
          <Skeleton active paragraph={{ rows: 3 }} />
        ) : (
          <Flexbox gap={16}>
            <Flexbox horizontal justify="space-between" align="center">
              <Flexbox gap={4}>
                <span style={{ fontSize: 14, color: theme.colorTextSecondary }}>Your KAWAI Balance</span>
                <span style={{ fontSize: 24, fontWeight: 700 }}>
                  {parseFloat(stats?.current_kawai_balance || '0').toLocaleString()} KAWAI
                </span>
              </Flexbox>
              <Flexbox gap={4} align="flex-end">
                <span style={{ fontSize: 14, color: theme.colorTextSecondary }}>Your Share</span>
                <span style={{ fontSize: 32, fontWeight: 700, color: theme.colorPrimary }}>
                  {sharePercentage.toFixed(4)}%
                </span>
              </Flexbox>
            </Flexbox>

            <Progress
              percent={Math.min(sharePercentage, 100)} // Actual percentage (no scaling)
              strokeColor={{
                '0%': '#26a17b',
                '100%': '#1a9065',
              }}
              status="active"
              showInfo={false}
            />

            <div
              style={{
                padding: 12,
                borderRadius: 8,
                background: theme.colorInfoBg,
                border: `1px solid ${theme.colorInfoBorder}`,
              }}
            >
              <Flexbox horizontal gap={8} align="center">
                <Info size={16} color={theme.colorInfo} />
                <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                  Your weekly USDT = (Your KAWAI / Total Supply) × Weekly Net Profit
                </span>
              </Flexbox>
            </div>
          </Flexbox>
        )}
      </Card>

      {/* Claimable Revenue */}
      <Card
        title={
          <Flexbox horizontal align="center" gap={8}>
            <Gift size={16} />
            <span>Claimable Revenue</span>
            {stats && stats.pending_claims > 0 && (
              <Tag color="green">{stats.pending_claims} available</Tag>
            )}
          </Flexbox>
        }
        size="small"
      >
        {loading ? (
          <Skeleton active paragraph={{ rows: 5 }} />
        ) : !stats?.unclaimed_records || stats.unclaimed_records.length === 0 ? (
          <Flexbox align="center" gap={16} style={{ padding: '40px 20px' }}>
            <TokenUSDT size={64} variant="branded" />
            <span style={{ fontSize: 16, color: theme.colorTextSecondary }}>
              No claimable revenue yet. Keep holding KAWAI to earn weekly USDT!
            </span>
            <div
              style={{
                padding: 16,
                borderRadius: 8,
                background: theme.colorWarningBg,
                border: `1px solid ${theme.colorWarningBorder}`,
                textAlign: 'center',
                maxWidth: 500,
              }}
            >
              <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
                💡 <strong>Revenue sharing starts in Phase 2</strong> after 1B KAWAI supply is reached.
                Currently in Phase 1 (Mining Era).
              </span>
            </div>
          </Flexbox>
        ) : (
          <Table
            dataSource={stats.unclaimed_records}
            rowKey="period_id"
            size="small"
            pagination={false}
            columns={[
              {
                title: 'Period',
                dataIndex: 'period_id',
                key: 'period',
                render: (period: number) => <Tag color="blue">Week #{period}</Tag>,
              },
              {
                title: 'Amount',
                dataIndex: 'formatted',
                key: 'amount',
                render: (formatted: string) => (
                  <span style={{ fontWeight: 600, color: theme.colorSuccess }}>
                    ${formatted} USDT
                  </span>
                ),
              },
              {
                title: 'Your Share',
                dataIndex: 'share_percentage',
                key: 'share',
                render: (share: string) => (
                  <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                    {parseFloat(share).toFixed(4)}%
                  </span>
                ),
              },
              {
                title: 'Date',
                dataIndex: 'created_at',
                key: 'date',
                render: (date: string) => (
                  <span style={{ fontSize: 12 }}>{formatDate(date)}</span>
                ),
              },
              {
                title: 'Status',
                dataIndex: 'claim_status',
                key: 'status',
                render: (status: string) => (
                  <Tag color={status === 'unclaimed' ? 'green' : 'default'}>
                    {status}
                  </Tag>
                ),
              },
              {
                title: 'Action',
                key: 'action',
                render: (_: any, record: RevenueShareRecord) => (
                  <Button
                    type="primary"
                    size="small"
                    icon={<Gift size={14} />}
                    loading={claimLoading.has(record.period_id)}
                    onClick={() => openConfirmModal(record)}
                    style={{
                      background: 'linear-gradient(135deg, #26a17b, #1a9065)',
                      border: 'none',
                    }}
                  >
                    Claim
                  </Button>
                ),
              },
            ]}
          />
        )}
      </Card>

      {/* How It Works */}
      <Card
        title={
          <Flexbox horizontal align="center" gap={8}>
            <Info size={16} />
            <span>How Revenue Sharing Works</span>
          </Flexbox>
        }
        size="small"
        style={{ background: theme.colorFillTertiary, border: 'none' }}
      >
        <Flexbox gap={16}>
          <ul style={{ margin: 0, paddingLeft: 20, fontSize: 13, color: theme.colorTextSecondary }}>
            <li style={{ marginBottom: 8 }}>
              <strong>100% of platform profit</strong> (USDT) is distributed to KAWAI holders every week
            </li>
            <li style={{ marginBottom: 8 }}>
              Your share is <strong>proportional to your KAWAI holdings</strong> (no lock/stake required)
            </li>
            <li style={{ marginBottom: 8 }}>
              Formula: <code>Your USDT = (Your KAWAI / Total Supply) × Weekly Net Profit</code>
            </li>
            <li style={{ marginBottom: 8 }}>
              Weekly settlements generate <strong>Merkle proofs</strong> for gas-efficient claiming
            </li>
            <li style={{ marginBottom: 8 }}>
              <strong>Phase 1 (Current):</strong> Mining Era - KAWAI tokens minted as rewards
            </li>
            <li style={{ marginBottom: 8 }}>
              <strong>Phase 2 (Future):</strong> USDT Era - Revenue sharing begins after 1B supply reached
            </li>
            <li>
              Estimated yield: <strong>Variable, based on platform revenue and holdings</strong>
            </li>
          </ul>

          <div
            style={{
              padding: 16,
              borderRadius: 8,
              background: 'linear-gradient(135deg, #667eea10, #764ba210)',
              border: '1px solid #667eea40',
            }}
          >
            <Flexbox horizontal gap={12} align="center">
              <Users size={24} color="#667eea" />
              <Flexbox gap={4} flex={1}>
                <span style={{ fontWeight: 600, fontSize: 14 }}>
                  Real Yield from Real Revenue
                </span>
                <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                  Unlike most crypto projects, KAWAI derives value from external revenue (AI service payments).
                  This is sustainable, not a ponzi scheme.
                </span>
              </Flexbox>
            </Flexbox>
          </div>
        </Flexbox>
      </Card>

      {/* Confirmation Modal */}
      <Modal
        title={
          <Flexbox horizontal align="center" gap={8}>
            <Info size={18} />
            <span>Confirm Revenue Claim</span>
          </Flexbox>
        }
        open={!!confirmModal}
        onOk={() => confirmModal && handleClaimRevenue(confirmModal)}
        onCancel={() => setConfirmModal(null)}
        okText="Confirm & Claim"
        confirmLoading={confirmModal ? claimLoading.has(confirmModal.period_id) : false}
      >
        {confirmModal && (
          <Flexbox gap={16} padding={8}>
            <div style={{ padding: 16, background: theme.colorFillSecondary, borderRadius: 8 }}>
              <Flexbox horizontal justify="space-between" align="center">
                <span>Revenue Amount:</span>
                <span style={{ fontWeight: 700, fontSize: 18 }}>
                  ${confirmModal.formatted} USDT
                </span>
              </Flexbox>
              <Flexbox horizontal justify="space-between" style={{ marginTop: 8 }}>
                <span style={{ color: theme.colorTextSecondary }}>Period:</span>
                <span>Week #{confirmModal.period_id}</span>
              </Flexbox>
              <Flexbox horizontal justify="space-between" style={{ marginTop: 8 }}>
                <span style={{ color: theme.colorTextSecondary }}>Your Share:</span>
                <span>{parseFloat(confirmModal.share_percentage).toFixed(4)}%</span>
              </Flexbox>
            </div>

            <Flexbox gap={4}>
              <span style={{ fontWeight: 600 }}>Estimated Gas Fee:</span>
              {estimateLoading ? (
                <span>Calculating...</span>
              ) : (
                <span style={{ fontFamily: 'monospace' }}>
                  {gasEstimate || 'Unknown'}
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
