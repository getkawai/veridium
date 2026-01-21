import { Card, Progress, Tag, Statistic, Row, Col, Button, Empty, Skeleton, Tooltip, App, Table } from 'antd';
import { Gift, Clock, Award, Info, DollarSign, Coins, ExternalLink } from 'lucide-react';
import { Flexbox } from 'react-layout-kit';
import { useState, useEffect, useCallback } from 'react';
import { createStyles } from 'antd-style';
import { CashbackService, DeAIService } from '@@/github.com/kawai-network/veridium/internal/services';
import { useUserStore } from '@/store/user';
import type { CashbackStatsResponse, ClaimableCashbackRecord } from '@@/github.com/kawai-network/veridium/internal/services/models';
import type { NetworkInfo } from '@@/github.com/kawai-network/veridium/internal/services/models';
import { Browser } from '@wailsio/runtime';

const useStyles = createStyles(({ css, token }) => ({
  tierCard: css`
    background: ${token.colorBgContainer};
    border: 1px solid ${token.colorBorder};
    border-radius: ${token.borderRadiusLG}px;
    padding: 24px;
  `,
  tierBadge: css`
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    background: linear-gradient(135deg, #667eea, #764ba2);
    color: white;
    border-radius: 20px;
    font-weight: 600;
    font-size: 16px;
  `,
  tierItem: css`
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    padding: 12px;
    border-radius: 8px;
    transition: all 0.3s;
    
    &.active {
      background: ${token.colorPrimaryBg};
      border: 2px solid ${token.colorPrimary};
    }
    
    &.locked {
      opacity: 0.5;
    }
  `,
  infoCard: css`
    background: ${token.colorInfoBg};
    border: 1px solid ${token.colorInfoBorder};
    border-radius: ${token.borderRadiusLG}px;
    padding: 16px;
  `,
}));

const CASHBACK_TIERS = [
  { level: 0, min: 0, max: 99, rate: 1, cap: 5000, label: 'Bronze' },
  { level: 1, min: 100, max: 499, rate: 1.25, cap: 10000, label: 'Silver' },
  { level: 2, min: 500, max: 999, rate: 1.5, cap: 15000, label: 'Gold' },
  { level: 3, min: 1000, max: 4999, rate: 1.75, cap: 20000, label: 'Platinum' },
  { level: 4, min: 5000, max: Infinity, rate: 2, cap: 20000, label: 'Diamond' },
];

interface CashbackRewardsSectionProps {
  currentNetwork: NetworkInfo | null;
  theme: any;
  styles: any;
  onOpenDepositModal?: () => void;
  onRefresh?: (refreshFn: () => void) => void;
}

export const CashbackRewardsSection = ({ currentNetwork, theme, styles: propStyles, onOpenDepositModal, onRefresh }: CashbackRewardsSectionProps) => {
  const { styles } = useStyles();
  const { message } = App.useApp();
  const userAddress = useUserStore((s) => s.walletAddress);

  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<CashbackStatsResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [currentPeriod, setCurrentPeriod] = useState<number>(0);
  const [claimableRecords, setClaimableRecords] = useState<ClaimableCashbackRecord[]>([]);
  const [claimLoading, setClaimLoading] = useState<Set<number>>(new Set());

  const loadCashbackStats = useCallback(async (address: string, showMessage = false) => {
    if (!address) {
      setError('No wallet connected');
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const [statsResult, periodResult, recordsResult] = await Promise.all([
        CashbackService.GetCashbackStats(address),
        CashbackService.GetCurrentPeriod(),
        CashbackService.GetClaimableCashback(address),
      ]);

      if (!statsResult) {
        setError('Failed to load cashback data. Please try again.');
        return;
      }

      // Convert wei to KAWAI for display
      const convertedStats = {
        ...statsResult,
        total_cashback: (BigInt(statsResult.total_cashback || '0') / BigInt(10 ** 18)).toString(),
        pending_cashback: (BigInt(statsResult.pending_cashback || '0') / BigInt(10 ** 18)).toString(),
        claimed_cashback: (BigInt(statsResult.claimed_cashback || '0') / BigInt(10 ** 18)).toString(),
        // Convert USDT wei (6 decimals) to USDT
        total_deposit_amount_usdt: Number(BigInt(statsResult.total_deposit_amount || '0') / BigInt(10 ** 6)),
      };

      setStats(convertedStats as any);
      setCurrentPeriod(periodResult);
      setClaimableRecords(recordsResult || []);

      if (showMessage) {
        message.success('Cashback data refreshed');
      }
    } catch (e: any) {
      console.error('Failed to load cashback stats:', e);
      setError(e.message || 'Failed to load cashback data');

      if (showMessage) {
        message.error('Failed to refresh cashback data');
      }
    } finally {
      setLoading(false);
    }
  }, [message]);

  useEffect(() => {
    if (userAddress) {
      loadCashbackStats(userAddress);
    }
  }, [userAddress, loadCashbackStats]);

  // Expose refresh function to parent
  useEffect(() => {
    if (onRefresh && userAddress) {
      onRefresh(() => loadCashbackStats(userAddress, true));
    }
  }, [onRefresh, userAddress, loadCashbackStats]);

  const getCurrentTierLevel = (totalDepositUSDT: number): number => {
    // Determine tier based on total USDT deposited
    for (let i = CASHBACK_TIERS.length - 1; i >= 0; i--) {
      if (totalDepositUSDT >= CASHBACK_TIERS[i].min) {
        return CASHBACK_TIERS[i].level;
      }
    }
    return 0; // Default to Bronze
  };

  const calculateTierProgress = () => {
    if (!stats) return { percent: 0, current: 0, next: 0 };

    const totalDepositUSDT = (stats as any).total_deposit_amount_usdt || 0;
    const currentTierLevel = getCurrentTierLevel(totalDepositUSDT);
    const currentTier = CASHBACK_TIERS[currentTierLevel];
    const nextTier = CASHBACK_TIERS[Math.min(currentTierLevel + 1, 4)];

    if (currentTier.level === 4) {
      return { percent: 100, current: totalDepositUSDT, next: totalDepositUSDT };
    }

    const progress = ((totalDepositUSDT - currentTier.min) / (nextTier.min - currentTier.min)) * 100;
    return {
      percent: Math.min(Math.max(progress, 0), 100),
      current: totalDepositUSDT,
      next: nextTier.min,
    };
  };

  if (error) {
    return (
      <Flexbox style={{ width: '100%' }} gap={20}>
        <div className={propStyles.placeholderCard}>
          <Award size={48} color={theme.colorTextQuaternary} style={{ marginBottom: 16 }} />
          <h3 style={{ margin: '0 0 8px', color: theme.colorError }}>Error Loading Cashback Data</h3>
          <p style={{ color: theme.colorTextSecondary, margin: '0 0 16px' }}>{error}</p>
          <Button onClick={() => userAddress && loadCashbackStats(userAddress, true)} icon={<Info size={16} />}>
            Retry
          </Button>
        </div>
      </Flexbox>
    );
  }

  const tierProgress = calculateTierProgress();
  const totalDepositUSDT = (stats as any)?.total_deposit_amount_usdt || 0;
  const currentTierLevel = getCurrentTierLevel(totalDepositUSDT);
  const currentTier = CASHBACK_TIERS[currentTierLevel];
  const nextTier = CASHBACK_TIERS[Math.min(currentTierLevel + 1, 4)];

  const handleClaimCashback = async (record: ClaimableCashbackRecord) => {
    if (claimLoading.has(record.period)) return;
    setClaimLoading(prev => new Set(prev).add(record.period));
    try {
      // Empty proof is valid for single-leaf Merkle trees
      const result = await DeAIService.ClaimCashbackReward(record.period, record.amount, record.proof || []);
      if (result?.tx_hash) {
        const explorerUrl = currentNetwork?.explorerURL;
        message.success(
          <span>
            Claim confirmed! Tx: {result.tx_hash.substring(0, 10)}...
            {explorerUrl && (
              <a
                onClick={() => Browser.OpenURL(`${explorerUrl}/tx/${result.tx_hash}`)}
                style={{ marginLeft: 8, cursor: 'pointer' }}
              >
                View <ExternalLink size={12} style={{ verticalAlign: 'middle' }} />
              </a>
            )}
          </span>
        );
        setTimeout(() => userAddress && loadCashbackStats(userAddress, true), 3000);
      }
    } catch (e: any) {
      console.error('Cashback claim failed:', e);
      message.error(e.message || 'Claim failed');
    } finally {
      setClaimLoading(prev => { const next = new Set(prev); next.delete(record.period); return next; });
    }
  };

  return (
    <Flexbox style={{ width: '100%' }} gap={20}>
      {/* Summary Cards */}
      <Row gutter={16}>
        <Col xs={24} sm={8}>
          <Card size="small" style={{ background: 'linear-gradient(135deg, #667eea20, #764ba220)', border: '1px solid #667eea40' }}>
            {loading ? (
              <Skeleton active paragraph={{ rows: 1 }} />
            ) : (
              <Statistic
                title="Total Cashback Earned"
                value={stats?.total_cashback || '0'}
                suffix="KAWAI"
                prefix={<Coins size={20} color="#667eea" />}
                valueStyle={{ color: '#667eea', fontWeight: 700 }}
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
                value={stats?.pending_cashback || '0'}
                suffix="KAWAI"
                prefix={<Gift size={20} color="#22c55e" />}
                valueStyle={{ color: '#22c55e', fontWeight: 700 }}
              />
            )}
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card size="small" style={{ background: 'linear-gradient(135deg, #f59e0b20, #d9770620)', border: '1px solid #f59e0b40' }}>
            {loading ? (
              <Skeleton active paragraph={{ rows: 1 }} />
            ) : (
              <Statistic
                title="Claimed"
                value={stats?.claimed_cashback || '0'}
                suffix="KAWAI"
                prefix={<Clock size={20} color="#f59e0b" />}
                valueStyle={{ color: '#f59e0b', fontWeight: 700 }}
              />
            )}
          </Card>
        </Col>
      </Row>

      {/* Tier Progress Section */}
      <Card title="Your Cashback Tier" size="small" className={styles.tierCard}>
        {loading ? (
          <Skeleton active paragraph={{ rows: 3 }} />
        ) : (
          <Flexbox gap={20}>
            <Flexbox horizontal justify="space-between" align="center">
              <div className={styles.tierBadge}>
                <Award size={20} />
                <span>Tier {currentTier.level}: {currentTier.label}</span>
                <Tag color="gold">{currentTier.rate}% Cashback</Tag>
              </div>
              <Tooltip title="Cap per deposit">
                <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>
                  Max: {currentTier.cap.toLocaleString()} KAWAI per deposit
                </span>
              </Tooltip>
            </Flexbox>

            {currentTier.level < 4 && (
              <Flexbox gap={8}>
                <Flexbox horizontal justify="space-between" align="center">
                  <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
                    Progress to Tier {nextTier.level} ({nextTier.label})
                  </span>
                  <span style={{ fontSize: 13, fontWeight: 600 }}>
                    ${tierProgress.current.toFixed(2)} / ${tierProgress.next.toFixed(2)}
                  </span>
                </Flexbox>
                <Progress
                  percent={tierProgress.percent}
                  strokeColor={{
                    '0%': '#667eea',
                    '100%': '#764ba2',
                  }}
                  status="active"
                />
                <span style={{ fontSize: 12, color: theme.colorTextTertiary }}>
                  Deposit ${(tierProgress.next - tierProgress.current).toFixed(2)} more to unlock {nextTier.rate}% cashback
                </span>
              </Flexbox>
            )}

            {currentTier.level === 4 && (
              <div style={{ textAlign: 'center', padding: '20px 0' }}>
                <span style={{ fontSize: 16, color: theme.colorSuccess, fontWeight: 600 }}>
                  🎉 Congratulations! You've reached the highest tier!
                </span>
              </div>
            )}

            {/* All Tiers Visualization */}
            <Flexbox horizontal justify="space-between" gap={8} style={{ marginTop: 16 }}>
              {CASHBACK_TIERS.map((tier) => (
                <div
                  key={tier.level}
                  className={`${styles.tierItem} ${tier.level === currentTier.level ? 'active' : ''} ${tier.level > currentTier.level ? 'locked' : ''}`}
                  style={{ flex: 1 }}
                >
                  <Award size={24} color={tier.level === currentTier.level ? theme.colorPrimary : theme.colorTextTertiary} />
                  <span style={{ fontSize: 11, fontWeight: 600 }}>{tier.label}</span>
                  <span style={{ fontSize: 10, color: theme.colorTextSecondary }}>{tier.rate}%</span>
                  {tier.level === currentTier.level && (
                    <Tag color="blue">Current</Tag>
                  )}
                </div>
              ))}
            </Flexbox>
          </Flexbox>
        )}
      </Card>

      {/* Claimable Cashback History */}
      <Card title="Claimable Cashback" size="small">
        {loading ? (
          <Skeleton active paragraph={{ rows: 5 }} />
        ) : claimableRecords.length === 0 ? (
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={
              <Flexbox gap={8} align="center">
                <span style={{ color: theme.colorTextSecondary }}>
                  No claimable cashback yet. Make your first deposit to start earning!
                </span>
                <span style={{ fontSize: 16, fontWeight: 600, color: theme.colorSuccess }}>
                  🎁 First deposit gets 5% bonus!
                </span>
                <Button
                  type="primary"
                  icon={<DollarSign size={16} />}
                  onClick={onOpenDepositModal}
                  style={{ marginTop: 12 }}
                >
                  Make First Deposit
                </Button>
              </Flexbox>
            }
          />
        ) : (
          <Table
            dataSource={claimableRecords}
            rowKey="period"
            size="small"
            pagination={false}
            columns={[
              {
                title: 'Period',
                dataIndex: 'period',
                key: 'period',
                render: (period: number) => <Tag color="blue">#{period}</Tag>,
              },
              {
                title: 'Amount',
                dataIndex: 'amount',
                key: 'amount',
                render: (amount: string) => {
                  const kawai = (BigInt(amount) / BigInt(10 ** 18)).toString();
                  return <span style={{ fontWeight: 600, color: theme.colorSuccess }}>{kawai} KAWAI</span>;
                },
              },
              {
                title: 'Status',
                dataIndex: 'claimed',
                key: 'claimed',
                render: (claimed: boolean) => (
                  claimed ? (
                    <Tag color="default">Claimed</Tag>
                  ) : (
                    <Tag color="green">Ready to Claim</Tag>
                  )
                ),
              },
              {
                title: 'Action',
                key: 'action',
                render: (_: any, record: ClaimableCashbackRecord) => (
                  record.claimed ? (
                    <span style={{ color: theme.colorTextTertiary, fontSize: 12 }}>Already claimed</span>
                  ) : (
                    <Button
                      type="primary"
                      size="small"
                      icon={<Gift size={14} />}
                      loading={claimLoading.has(record.period)}
                      onClick={() => handleClaimCashback(record)}
                    >
                      Claim
                    </Button>
                  )
                ),
              },
            ]}
          />
        )}
      </Card>

      {/* How It Works */}
      <Card title="How Cashback Works" size="small" className={styles.infoCard}>
        <Flexbox gap={8}>
          <ul style={{ margin: 0, paddingLeft: 20, fontSize: 13, color: theme.colorTextSecondary }}>
            <li>Earn <strong>1-2% KAWAI cashback</strong> on every {currentNetwork?.stablecoinSymbol || 'USDT'} deposit</li>
            <li>First deposit always gets <strong>5% bonus</strong> regardless of amount</li>
            <li>Higher total deposits unlock better cashback rates</li>
            <li>Cashback caps per deposit: 5K-20K KAWAI (prevents abuse)</li>
            <li>Claims available <strong>weekly</strong> after settlement period</li>
            <li><strong>200M KAWAI</strong> allocated for cashback program (~3 year runway)</li>
            <li>Unlimited deposits = Unlimited cashback earnings!</li>
          </ul>
          <Flexbox horizontal gap={8} style={{ marginTop: 12, padding: '12px', background: theme.colorWarningBg, borderRadius: 8 }}>
            <Info size={16} color={theme.colorWarning} />
            <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
              Current Period: <strong>#{currentPeriod}</strong> • Settlement: Weekly (Every Monday)
            </span>
          </Flexbox>
        </Flexbox>
      </Card>
    </Flexbox>
  );
};

