import { Card, Button, Statistic, Row, Col, App, Skeleton, Divider } from 'antd';
import { useState, useEffect, useCallback } from 'react';
import { ReferralService } from '@@/github.com/kawai-network/veridium/internal/services';
import { Flexbox } from 'react-layout-kit';
import { Users, DollarSign, Gift, Share2, Copy, CheckCircle, Info, Zap, TrendingUp } from 'lucide-react';
import { useUserStore } from '@/store/user';
import { copyText } from '@/utils/clipboard';
import type { NetworkInfo } from '@@/github.com/kawai-network/veridium/internal/services/models';

interface ReferralRewardsSectionProps {
  currentNetwork: NetworkInfo | null;
  theme: any;
  styles: any;
  onRefresh?: (refreshFn: () => void) => void;
}

export const ReferralRewardsSection = ({ currentNetwork, theme, styles, onRefresh }: ReferralRewardsSectionProps) => {
  const { message } = App.useApp();
  const userAddress = useUserStore((s) => s.walletAddress);
  
  const [loading, setLoading] = useState(true);
  const [stats, setStats] = useState<{
    code: string;
    total_referrals: number;
    total_earnings_usdt: number;
    total_earnings_kawai: string;
  } | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const loadReferralStats = useCallback(async (address: string, showMessage = false) => {
    if (!address) {
      setError('No wallet connected');
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);
    
    try {
      const result = await ReferralService.GetReferralStats(address);

      if (!result) {
        setError('Failed to load referral data. Please try again.');
        return;
      }

      setStats(result);

      if (showMessage) {
        message.success('Referral data refreshed');
      }
    } catch (e: any) {
      console.error('Failed to load referral stats:', e);
      setError(e.message || 'Failed to load referral data');
      
      if (showMessage) {
        message.error('Failed to refresh referral data');
      }
    } finally {
      setLoading(false);
    }
  }, [message]);

  useEffect(() => {
    if (userAddress) {
      loadReferralStats(userAddress);
    }
  }, [userAddress, loadReferralStats]);

  // Expose refresh function to parent
  useEffect(() => {
    if (onRefresh && userAddress) {
      const refreshFn = () => loadReferralStats(userAddress, true);
      onRefresh(refreshFn);
    }
  }, [onRefresh, userAddress, loadReferralStats]);

  const handleCopyCode = async () => {
    if (!stats?.code) return;
    
    try {
      await copyText(stats.code);
      setCopied(true);
      message.success('Referral code copied!');
      setTimeout(() => setCopied(false), 2000);
    } catch (e) {
      message.error('Failed to copy code');
    }
  };

  const handleShare = async () => {
    if (!stats?.code) return;

    const shareText = `Join Kawai DeAI Network and get 10 USDT + 200 KAWAI FREE! Use my code: ${stats.code}\n\nDecentralized AI • No credit card • Instant access\n\n`;
    const shareUrl = `https://getkawai.com?ref=${stats.code}`;
    const fullShareText = shareText + shareUrl;

    try {
      await copyText(fullShareText);
      message.success('Referral link copied to clipboard!');
    } catch (err) {
      message.error('Failed to copy referral link');
    }
  };

  const formatKawaiAmount = (rawAmount: string) => {
    if (!rawAmount || rawAmount === '0') return '0.0000';
    try {
      const amount = BigInt(rawAmount);
      const divisor = BigInt(10 ** 18);
      const wholePart = amount / divisor;
      const fractionalPart = amount % divisor;
      
      const fractionalStr = fractionalPart.toString().padStart(18, '0');
      const trimmedFractional = fractionalStr.slice(0, 4);
      
      return `${wholePart}.${trimmedFractional}`;
    } catch (e) {
      console.error('Failed to format KAWAI amount:', e);
      return '0.0000';
    }
  };

  // Check if error is "no referral code" (expected for new users)
  const isNewUser = error?.includes('no referral code') || error?.includes('key not found');

  if (error && !isNewUser) {
    // Only show error UI for unexpected errors
    return (
      <Flexbox style={{ width: '100%' }} gap={20}>
        <div className={styles.placeholderCard}>
          <Users size={48} color={theme.colorTextQuaternary} style={{ marginBottom: 16 }} />
          <h3 style={{ margin: '0 0 8px', color: theme.colorError }}>Error Loading Referral Data</h3>
          <p style={{ color: theme.colorTextSecondary, margin: '0 0 16px' }}>{error}</p>
          <Button onClick={() => userAddress && loadReferralStats(userAddress, true)} icon={<Info size={16} />}>
            Retry
          </Button>
        </div>
      </Flexbox>
    );
  }

  // For new users without referral code, show welcome/create UI
  if (isNewUser) {
    return (
      <Flexbox style={{ width: '100%' }} gap={20}>
        <div className={styles.placeholderCard}>
          <Users size={48} color={theme.colorPrimary} style={{ marginBottom: 16 }} />
          <h3 style={{ margin: '0 0 8px', color: theme.colorText }}>Create Your Referral Code</h3>
          <p style={{ color: theme.colorTextSecondary, margin: '0 0 16px', maxWidth: 400, textAlign: 'center' }}>
            Start earning rewards by referring friends! Generate your unique referral code and earn 5 USDT + 100 KAWAI for each successful referral.
          </p>
          <Button 
            type="primary" 
            size="large"
            icon={<Gift size={16} />}
            onClick={async () => {
              try {
                setLoading(true);
                setError(null);
                // Call CreateReferralCode service
                const code = await ReferralService.CreateReferralCode(userAddress!);
                message.success(`Referral code created: ${code}`);
                // Reload stats to show the new code
                await loadReferralStats(userAddress!, false);
              } catch (e: any) {
                console.error('Failed to create referral code:', e);
                message.error(e.message || 'Failed to create referral code');
                setError(e.message || 'Failed to create referral code');
              } finally {
                setLoading(false);
              }
            }}
            loading={loading}
          >
            Generate Referral Code
          </Button>
        </div>
      </Flexbox>
    );
  }

  return (
    <Flexbox style={{ width: '100%' }} gap={20}>
      {/* Summary Cards */}
      <Row gutter={16}>
        <Col xs={24} sm={8}>
          <Card size="small" style={{ background: 'linear-gradient(135deg, #3b82f620, #2563eb20)', border: '1px solid #3b82f640' }}>
            {loading ? (
              <Skeleton active paragraph={{ rows: 1 }} />
            ) : (
              <Statistic
                title="Total Referrals"
                value={stats?.total_referrals || 0}
                prefix={<Users size={20} color="#3b82f6" />}
                valueStyle={{ color: '#3b82f6', fontWeight: 700 }}
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
                title="Earned (USDT)"
                value={stats?.total_earnings_usdt || 0}
                prefix={<DollarSign size={20} color="#22c55e" />}
                suffix="USDT"
                precision={2}
                valueStyle={{ color: '#22c55e', fontWeight: 700 }}
              />
            )}
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card size="small" style={{ background: 'linear-gradient(135deg, #667eea20, #764ba220)', border: '1px solid #667eea40' }}>
            {loading ? (
              <Skeleton active paragraph={{ rows: 1 }} />
            ) : (
              <Statistic
                title="Earned (KAWAI)"
                value={formatKawaiAmount(stats?.total_earnings_kawai || '0')}
                prefix={<Gift size={20} color="#667eea" />}
                valueStyle={{ color: '#667eea', fontWeight: 700 }}
              />
            )}
          </Card>
        </Col>
      </Row>

      {/* Referral Code Card */}
      <Card
        title={
          <Flexbox horizontal align="center" gap={8}>
            <Share2 size={16} />
            <span>Your Referral Code</span>
          </Flexbox>
        }
        size="small"
      >
        {loading ? (
          <Skeleton active paragraph={{ rows: 2 }} />
        ) : (
          <Flexbox gap={16}>
            <div
              style={{
                padding: '24px',
                background: theme.colorBgLayout,
                borderRadius: 12,
                border: `2px dashed ${theme.colorPrimary}`,
                textAlign: 'center',
              }}
            >
              <span
                style={{
                  fontSize: 32,
                  fontWeight: 700,
                  letterSpacing: 4,
                  color: theme.colorPrimary,
                  fontFamily: 'Monaco, Courier New, monospace',
                  display: 'block',
                  marginBottom: 16,
                }}
              >
                {stats?.code || 'LOADING'}
              </span>
              <Flexbox horizontal justify="center" gap={8}>
                <Button
                  icon={copied ? <CheckCircle size={16} /> : <Copy size={16} />}
                  onClick={handleCopyCode}
                  type={copied ? 'default' : 'primary'}
                >
                  {copied ? 'Copied!' : 'Copy Code'}
                </Button>
                <Button
                  icon={<Share2 size={16} />}
                  onClick={handleShare}
                  type="default"
                >
                  Copy Share Link
                </Button>
              </Flexbox>
            </div>

            <div
              style={{
                padding: '12px 16px',
                background: theme.colorInfoBg,
                borderRadius: 8,
                textAlign: 'center',
              }}
            >
              <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                Share this code with friends. They get <strong>10 USDT + 200 KAWAI</strong>, you get <strong>5 USDT + 100 KAWAI</strong> per referral + <strong>5% lifetime mining commission</strong>!
              </span>
            </div>
          </Flexbox>
        )}
      </Card>

      {/* Mining Commission Explainer */}
      <Card
        title={
          <Flexbox horizontal align="center" gap={8}>
            <Zap size={16} />
            <span>Lifetime Mining Commission</span>
          </Flexbox>
        }
        size="small"
        style={{ 
          background: 'linear-gradient(135deg, #f59e0b10, #d9770610)',
          border: '1px solid #f59e0b40'
        }}
      >
        <Flexbox gap={16}>
          <Flexbox horizontal gap={12} align="center">
            <div
              style={{
                padding: 12,
                borderRadius: 8,
                background: 'linear-gradient(135deg, #f59e0b20, #d9770620)',
                border: '1px solid #f59e0b60',
              }}
            >
              <TrendingUp size={24} color="#f59e0b" />
            </div>
            <Flexbox gap={4} flex={1}>
              <span style={{ fontWeight: 600, fontSize: 15 }}>
                Earn 5% of all mining rewards from your referrals — forever!
              </span>
              <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
                Every time someone you referred uses AI, you automatically earn 5% of the mining rewards in KAWAI tokens.
              </span>
            </Flexbox>
          </Flexbox>

          <Divider style={{ margin: 0 }} />

          <Flexbox gap={12}>
            <span style={{ fontSize: 13, fontWeight: 600, color: theme.colorText }}>
              How it works:
            </span>
            <ul style={{ margin: 0, paddingLeft: 20, fontSize: 13, color: theme.colorTextSecondary }}>
              <li style={{ marginBottom: 8 }}>
                Your referral uses AI (e.g., generates 1000 tokens)
              </li>
              <li style={{ marginBottom: 8 }}>
                Mining reward is distributed: <strong>85% contributor</strong>, <strong>5% developer</strong>, <strong>5% user</strong>, <strong>5% you (affiliator)</strong>
              </li>
              <li style={{ marginBottom: 8 }}>
                You receive 5% of the mining rewards automatically — no action needed
              </li>
              <li>
                This continues <strong>for the lifetime of your referral's usage</strong> — passive income!
              </li>
            </ul>
          </Flexbox>

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
                <strong>Note:</strong> Mining commission is tracked on-chain and will be claimable weekly via Merkle proof (same as mining rewards).
              </span>
            </Flexbox>
          </div>
        </Flexbox>
      </Card>

      {/* Referral Benefits */}
      <Card
        title={
          <Flexbox horizontal align="center" gap={8}>
            <Gift size={16} />
            <span>Referral Benefits</span>
          </Flexbox>
        }
        size="small"
      >
        <Flexbox gap={12}>
          <div
            style={{
              padding: 16,
              borderRadius: 8,
              background: 'linear-gradient(135deg, #22c55e10, #16a34a10)',
              border: '1px solid #22c55e40',
            }}
          >
            <Flexbox horizontal justify="space-between" align="center">
              <Flexbox gap={4}>
                <span style={{ fontWeight: 600, fontSize: 14 }}>Your Friend Gets</span>
                <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                  When they sign up with your code
                </span>
              </Flexbox>
              <Flexbox align="flex-end" gap={4}>
                <span style={{ fontSize: 20, fontWeight: 700, color: '#22c55e' }}>
                  10 USDT
                </span>
                <span style={{ fontSize: 16, fontWeight: 600, color: '#667eea' }}>
                  + 200 KAWAI
                </span>
              </Flexbox>
            </Flexbox>
          </div>

          <div
            style={{
              padding: 16,
              borderRadius: 8,
              background: 'linear-gradient(135deg, #3b82f610, #2563eb10)',
              border: '1px solid #3b82f640',
            }}
          >
            <Flexbox horizontal justify="space-between" align="center">
              <Flexbox gap={4}>
                <span style={{ fontWeight: 600, fontSize: 14 }}>You Get</span>
                <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                  For each successful referral
                </span>
              </Flexbox>
              <Flexbox align="flex-end" gap={4}>
                <span style={{ fontSize: 20, fontWeight: 700, color: '#22c55e' }}>
                  5 USDT
                </span>
                <span style={{ fontSize: 16, fontWeight: 600, color: '#667eea' }}>
                  + 100 KAWAI
                </span>
              </Flexbox>
            </Flexbox>
          </div>
        </Flexbox>
      </Card>

      {/* How It Works */}
      <Card
        title={
          <Flexbox horizontal align="center" gap={8}>
            <Info size={16} />
            <span>How Referrals Work</span>
          </Flexbox>
        }
        size="small"
        style={{ background: theme.colorFillTertiary, border: 'none' }}
      >
        <Flexbox gap={12}>
          <Flexbox horizontal gap={12} align="flex-start">
            <div
              style={{
                width: 24,
                height: 24,
                borderRadius: '50%',
                background: theme.colorPrimary,
                color: 'white',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontWeight: 600,
                fontSize: 12,
                flexShrink: 0,
              }}
            >
              1
            </div>
            <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
              Share your referral code with friends via social media, email, or direct message
            </span>
          </Flexbox>

          <Flexbox horizontal gap={12} align="flex-start">
            <div
              style={{
                width: 24,
                height: 24,
                borderRadius: '50%',
                background: theme.colorPrimary,
                color: 'white',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontWeight: 600,
                fontSize: 12,
                flexShrink: 0,
              }}
            >
              2
            </div>
            <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
              They sign up and enter your code during wallet setup
            </span>
          </Flexbox>

          <Flexbox horizontal gap={12} align="flex-start">
            <div
              style={{
                width: 24,
                height: 24,
                borderRadius: '50%',
                background: theme.colorPrimary,
                color: 'white',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontWeight: 600,
                fontSize: 12,
                flexShrink: 0,
              }}
            >
              3
            </div>
            <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
              They instantly receive <strong>10 USDT + 200 KAWAI</strong> (double the normal bonus!)
            </span>
          </Flexbox>

          <Flexbox horizontal gap={12} align="flex-start">
            <div
              style={{
                width: 24,
                height: 24,
                borderRadius: '50%',
                background: theme.colorPrimary,
                color: 'white',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontWeight: 600,
                fontSize: 12,
                flexShrink: 0,
              }}
            >
              4
            </div>
            <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
              You receive <strong>5 USDT + 100 KAWAI</strong> as a thank you for spreading the word
            </span>
          </Flexbox>

          <Flexbox horizontal gap={12} align="flex-start">
            <div
              style={{
                width: 24,
                height: 24,
                borderRadius: '50%',
                background: theme.colorPrimary,
                color: 'white',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontWeight: 600,
                fontSize: 12,
                flexShrink: 0,
              }}
            >
              5
            </div>
            <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
              <strong>Bonus:</strong> Earn 5% of all mining rewards every time your referral uses AI — forever!
            </span>
          </Flexbox>

          <Flexbox horizontal gap={12} align="flex-start">
            <div
              style={{
                width: 24,
                height: 24,
                borderRadius: '50%',
                background: theme.colorPrimary,
                color: 'white',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontWeight: 600,
                fontSize: 12,
                flexShrink: 0,
              }}
            >
              6
            </div>
            <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
              <strong>Unlimited referrals = Unlimited earnings!</strong> No caps, no limits.
            </span>
          </Flexbox>
        </Flexbox>
      </Card>

      {/* Stats Note */}
      {!loading && stats && stats.total_referrals === 0 && (
        <div
          style={{
            padding: '16px',
            background: theme.colorWarningBg,
            borderRadius: 8,
            border: `1px solid ${theme.colorWarningBorder}`,
            textAlign: 'center',
          }}
        >
          <span style={{ fontSize: 13, color: theme.colorTextSecondary }}>
            🎯 <strong>Start earning today!</strong> Share your referral code and watch your rewards grow.
          </span>
        </div>
      )}
    </Flexbox>
  );
};

