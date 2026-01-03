'use client';

import { Button, Text, Icon, CopyButton } from '@lobehub/ui';
import { Card, Flex, Statistic, Divider, Empty, List, Tag } from 'antd';
import { createStyles } from 'antd-style';
import { Users, DollarSign, Gift, Share2 } from 'lucide-react';
import { memo } from 'react';

const useStyles = createStyles(({ css, token }) => ({
  card: css`
    background: ${token.colorBgContainer};
    border-radius: ${token.borderRadiusLG}px;
    border: 1px solid ${token.colorBorder};
  `,
  referralCode: css`
    padding: 16px;
    background: ${token.colorBgLayout};
    border-radius: 8px;
    text-align: center;
    border: 2px dashed ${token.colorPrimary};
  `,
  codeText: css`
    font-size: 32px;
    font-weight: 700;
    letter-spacing: 4px;
    color: ${token.colorPrimary};
    font-family: 'Monaco', 'Courier New', monospace;
  `,
  statCard: css`
    background: ${token.colorBgElevated};
    padding: 16px;
    border-radius: 8px;
    border: 1px solid ${token.colorBorder};
  `,
}));

interface ReferralStats {
  referralCode: string;
  totalReferrals: number;
  totalEarnings: number;
  pendingReferrals: number;
  recentReferrals: Array<{
    address: string;
    timestamp: number;
    earned: number;
    status: 'pending' | 'completed';
  }>;
}

interface ReferralDashboardProps {
  stats: ReferralStats;
  onShare?: () => void;
}

export const ReferralDashboard = memo<ReferralDashboardProps>(({ stats, onShare }) => {
  const { styles } = useStyles();

  const shareText = `Join Kawai DeAI Network and get 10 USDT + 200 KAWAI FREE! Use my code: ${stats.referralCode}\n\nDecentralized AI • No credit card • Instant access\n\n`;

  const handleShare = async () => {
    if (navigator.share) {
      try {
        await navigator.share({
          title: 'Join Kawai DeAI Network',
          text: shareText,
          url: `https://kawai.network?ref=${stats.referralCode}`,
        });
      } catch (err) {
        // User cancelled share
      }
    } else {
      // Fallback: copy to clipboard
      try {
        await navigator.clipboard.writeText(shareText + `https://kawai.network?ref=${stats.referralCode}`);
        onShare?.();
      } catch (clipboardErr) {
        console.error('Failed to copy to clipboard:', clipboardErr);
        // Show error message to user
        alert('Failed to copy referral link. Please try again.');
      }
    }
  };

  return (
    <Flex vertical gap="large">
      {/* Referral Code Card */}
      <Card className={styles.card}>
        <Flex vertical gap="middle">
          <Text strong style={{ fontSize: 16 }}>Your Referral Code</Text>
          <div className={styles.referralCode}>
            <Text className={styles.codeText}>{stats.referralCode}</Text>
            <Flex justify="center" gap="small" style={{ marginTop: 12 }}>
              <CopyButton content={stats.referralCode} />
              <Button
                type="primary"
                icon={<Share2 size={16} />}
                onClick={handleShare}
              >
                Share
              </Button>
            </Flex>
          </div>
          <Text type="secondary" style={{ fontSize: 12, textAlign: 'center' }}>
            Share this code with friends. They get 10 USDT + 200 KAWAI, you get 5 USDT + 100 KAWAI per referral!
          </Text>
        </Flex>
      </Card>

      {/* Stats Grid */}
      <Flex gap="middle" wrap="wrap">
        <div className={styles.statCard} style={{ flex: 1, minWidth: 150 }}>
          <Statistic
            title="Total Referrals"
            value={stats.totalReferrals}
            prefix={<Icon icon={Users} size={20} />}
            valueStyle={{ color: '#3f8600' }}
          />
        </div>
        <div className={styles.statCard} style={{ flex: 1, minWidth: 150 }}>
          <Statistic
            title="Total Earned"
            value={stats.totalEarnings}
            prefix={<Icon icon={DollarSign} size={20} />}
            suffix="USDT"
            precision={2}
            valueStyle={{ color: '#cf1322' }}
          />
        </div>
        <div className={styles.statCard} style={{ flex: 1, minWidth: 150 }}>
          <Statistic
            title="Pending"
            value={stats.pendingReferrals}
            prefix={<Icon icon={Gift} size={20} />}
            valueStyle={{ color: '#faad14' }}
          />
        </div>
      </Flex>

      {/* Recent Referrals */}
      <Card className={styles.card}>
        <Text strong style={{ fontSize: 16, display: 'block', marginBottom: 16 }}>
          Recent Referrals
        </Text>
        {stats.recentReferrals.length === 0 ? (
          <Empty
            description="No referrals yet"
            image={Empty.PRESENTED_IMAGE_SIMPLE}
          >
            <Text type="secondary">Share your code to start earning!</Text>
          </Empty>
        ) : (
          <List
            dataSource={stats.recentReferrals}
            renderItem={(item) => (
              <List.Item>
                <Flex justify="space-between" style={{ width: '100%' }} align="center">
                  <div>
                    <Text strong style={{ display: 'block' }}>
                      {item.address.length > 14 
                        ? `${item.address.slice(0, 8)}...${item.address.slice(-6)}`
                        : item.address
                      }
                    </Text>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                      {new Date(item.timestamp).toLocaleDateString()}
                    </Text>
                  </div>
                  <Flex align="center" gap="small">
                    <Text strong style={{ color: '#52c41a' }}>
                      +{item.earned} USDT
                    </Text>
                    <Tag color={item.status === 'completed' ? 'success' : 'warning'}>
                      {item.status}
                    </Tag>
                  </Flex>
                </Flex>
              </List.Item>
            )}
          />
        )}
      </Card>

      <Divider />

      {/* How It Works */}
      <Card className={styles.card}>
        <Text strong style={{ fontSize: 16, display: 'block', marginBottom: 16 }}>
          How Referrals Work
        </Text>
        <Flex vertical gap="small">
          <Flex gap="small">
            <Text strong style={{ color: '#1890ff' }}>1.</Text>
            <Text>Share your referral code with friends</Text>
          </Flex>
          <Flex gap="small">
            <Text strong style={{ color: '#1890ff' }}>2.</Text>
            <Text>They sign up and enter your code during onboarding</Text>
          </Flex>
          <Flex gap="small">
            <Text strong style={{ color: '#1890ff' }}>3.</Text>
            <Text>They get 10 USDT + 200 KAWAI bonus (instead of 5 USDT + 100 KAWAI)</Text>
          </Flex>
          <Flex gap="small">
            <Text strong style={{ color: '#1890ff' }}>4.</Text>
            <Text>You get 5 USDT + 100 KAWAI for each successful referral</Text>
          </Flex>
          <Flex gap="small">
            <Text strong style={{ color: '#1890ff' }}>5.</Text>
            <Text>Unlimited referrals = Unlimited earnings!</Text>
          </Flex>
        </Flex>
      </Card>
    </Flex>
  );
});

ReferralDashboard.displayName = 'ReferralDashboard';

