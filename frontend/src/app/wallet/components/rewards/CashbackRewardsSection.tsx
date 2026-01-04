import { Card, Table, Progress, Tag, Statistic, Row, Col, Button, Empty, Skeleton, Tooltip, App } from 'antd';
import { Gift, Clock, Award, Info, DollarSign, Coins } from 'lucide-react';
import { Flexbox } from 'react-layout-kit';
import { useState, useEffect, useCallback } from 'react';
import { createStyles } from 'antd-style';
import { CashbackService } from '@@/github.com/kawai-network/veridium/internal/services';
import { useUserStore } from '@/store/user';
import type { CashbackStatsResponse } from '@@/github.com/kawai-network/veridium/internal/services/models';
import type { NetworkInfo } from '@@/github.com/kawai-network/veridium/internal/services/models';

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
  { level: 1, min: 100, max: 499, rate: 2, cap: 10000, label: 'Silver' },
  { level: 2, min: 500, max: 999, rate: 3, cap: 15000, label: 'Gold' },
  { level: 3, min: 1000, max: 4999, rate: 4, cap: 18000, label: 'Platinum' },
  { level: 4, min: 5000, max: Infinity, rate: 5, cap: 20000, label: 'Diamond' },
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

  const loadCashbackStats = useCallback(async (address: string, showMessage = false) => {
    if (!address) {
      setError('No wallet connected');
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);
    
    try {
      const [statsResult, periodResult] = await Promise.all([
        CashbackService.GetCashbackStats(address),
        CashbackService.GetCurrentPeriod(),
      ]);

      if (!statsResult) {
        setError('Failed to load cashback data. Please try again.');
        return;
      }

      setStats(statsResult);
      setCurrentPeriod(periodResult);

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

  const calculateTierProgress = () => {
    if (!stats) return { percent: 0, current: 0, next: 0 };
    
    const currentTier = CASHBACK_TIERS[stats.currentTier || 0];
    const nextTier = CASHBACK_TIERS[Math.min((stats.currentTier || 0) + 1, 4)];
    const totalDeposits = parseFloat(stats.totalDeposits || '0');
    
    if (currentTier.level === 4) {
      return { percent: 100, current: totalDeposits, next: totalDeposits };
    }
    
    const progress = ((totalDeposits - currentTier.min) / (nextTier.min - currentTier.min)) * 100;
    return {
      percent: Math.min(Math.max(progress, 0), 100),
      current: totalDeposits,
      next: nextTier.min,
    };
  };

  if (error) {
    return (
      <Flexbox style={{ width: '100%' }} gap={20}>
        <div className={propStyles.placeholderCard}>
          <Award size={48} color={theme.colorTextQuaternary} style={{ marginBottom: 16}} />
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
  const currentTier = CASHBACK_TIERS[stats?.currentTier || 0];
  const nextTier = CASHBACK_TIERS[Math.min((stats?.currentTier || 0) + 1, 4)];

  return (
    <Flexbox style={{ width: '100%' }} gap={20}>
      {/* Coming Soon Banner */}
      <div
        style={{
          padding: '12px 16px',
          background: theme.colorInfoBg,
          borderRadius: 8,
          border: `1px solid ${theme.colorInfoBorder}`,
          display: 'flex',
          alignItems: 'center',
          gap: 12,
        }}
      >
        <Info size={16} color={theme.colorInfo} />
        <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
          🚧 <strong>Cashback claiming coming soon!</strong> Backend Merkle proof generation is in progress. You can view your stats and tier progress now.
        </span>
      </div>

      {/* Summary Cards */}
      <Row gutter={16}>
        <Col xs={24} sm={8}>
          <Card size="small" style={{ background: 'linear-gradient(135deg, #667eea20, #764ba220)', border: '1px solid #667eea40' }}>
            {loading ? (
              <Skeleton active paragraph={{ rows: 1 }} />
            ) : (
              <Statistic
                title="Total Cashback Earned"
                value={stats?.totalCashbackEarned || '0'}
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
                value={stats?.claimableAmount || '0'}
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
                title="Pending (This Period)"
                value={stats?.pendingAmount || '0'}
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
                    <Tag color="blue" size="small">Current</Tag>
                  )}
                </div>
              ))}
            </Flexbox>
          </Flexbox>
        )}
      </Card>

      {/* Deposit History */}
      <Card title="Deposit History" size="small">
        {loading ? (
          <Skeleton active paragraph={{ rows: 5 }} />
        ) : !stats?.depositHistory || stats.depositHistory.length === 0 ? (
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={
              <Flexbox gap={8} align="center">
                <span style={{ color: theme.colorTextSecondary }}>
                  No deposits yet. Make your first deposit to start earning cashback!
                </span>
                <span style={{ fontSize: 16, fontWeight: 600, color: theme.colorSuccess }}>
                  🎁 First deposit gets 5% bonus!
                </span>
              </Flexbox>
            }
          >
            <Button 
              type="primary" 
              icon={<DollarSign size={16} />}
              onClick={onOpenDepositModal}
            >
              Make First Deposit
            </Button>
          </Empty>
        ) : (
          <Table
            dataSource={stats.depositHistory}
            rowKey="txHash"
            pagination={{ pageSize: 10, showSizeChanger: false }}
            size="small"
            columns={[
              {
                title: 'Date',
                dataIndex: 'timestamp',
                key: 'date',
                render: (timestamp: string) => new Date(timestamp).toLocaleDateString(),
              },
              {
                title: 'Deposit Amount',
                dataIndex: 'amount',
                key: 'amount',
                render: (amount: string) => (
                  <span style={{ fontWeight: 600 }}>
                    ${parseFloat(amount).toFixed(2)} USDT
                  </span>
                ),
              },
              {
                title: 'Rate',
                dataIndex: 'cashbackRate',
                key: 'rate',
                render: (rate: number, record: any) => (
                  <Flexbox horizontal align="center" gap={4}>
                    <span>{rate}%</span>
                    {rate === 5 && (
                      <Tooltip title="First deposit bonus">
                        <Tag color="gold" size="small">First Deposit</Tag>
                      </Tooltip>
                    )}
                  </Flexbox>
                ),
              },
              {
                title: 'Cashback Earned',
                dataIndex: 'cashbackAmount',
                key: 'cashback',
                render: (amount: string) => (
                  <span style={{ color: '#667eea', fontWeight: 600 }}>
                    {parseFloat(amount).toFixed(2)} KAWAI
                  </span>
                ),
              },
              {
                title: 'Status',
                dataIndex: 'claimed',
                key: 'status',
                render: (claimed: boolean) => (
                  <Tag color={claimed ? 'success' : 'warning'}>
                    {claimed ? 'Claimed' : 'Pending'}
                  </Tag>
                ),
              },
              {
                title: 'Action',
                key: 'action',
                render: (record: any) => (
                  record.claimed ? (
                    <Button size="small" type="text" disabled>
                      Claimed
                    </Button>
                  ) : (
                    <Tooltip title="Claiming coming soon! Backend integration in progress.">
                      <Button size="small" type="primary" icon={<Gift size={14} />} disabled>
                        Claim
                      </Button>
                    </Tooltip>
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
            <li>Earn <strong>1-5% KAWAI cashback</strong> on every USDT deposit</li>
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

