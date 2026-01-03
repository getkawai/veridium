'use client';

import { Button, Text, Icon } from '@lobehub/ui';
import { Flex, Input, message } from 'antd';
import { createStyles } from 'antd-style';
import { Gift, ChevronDown, ChevronUp } from 'lucide-react';
import { memo, useState } from 'react';

const useStyles = createStyles(({ css, token, prefixCls }) => ({
  banner: css`
    padding: 16px;
    background: linear-gradient(135deg, ${token.colorPrimaryBg} 0%, ${token.colorPrimaryBgHover} 100%);
    border: 1px solid ${token.colorPrimaryBorder};
    border-radius: 8px;
    margin-bottom: 16px;
  `,
  bonusAmount: css`
    font-size: 24px;
    font-weight: 700;
    color: ${token.colorPrimary};
  `,
  comparisonText: css`
    font-size: 12px;
    color: ${token.colorTextSecondary};
    text-decoration: line-through;
  `,
  expandButton: css`
    padding: 0;
    height: auto;
    font-size: 12px;
  `,
}));

interface ReferralBannerProps {
  onReferralApplied?: (code: string) => void;
}

export const ReferralBanner = memo<ReferralBannerProps>(({ onReferralApplied }) => {
  const { styles, theme } = useStyles();
  const [isExpanded, setIsExpanded] = useState(false);
  const [referralCode, setReferralCode] = useState('');
  const [isApplied, setIsApplied] = useState(false);

  const handleApplyCode = () => {
    if (!referralCode.trim()) {
      message.error('Please enter a referral code');
      return;
    }

    // Validate code format (6 alphanumeric characters)
    if (!/^[A-Z0-9]{6}$/.test(referralCode.toUpperCase())) {
      message.error('Invalid referral code format');
      return;
    }

    setIsApplied(true);
    message.success('Referral code applied! You\'ll get 10 USDT + 200 KAWAI bonus 🎉');
    onReferralApplied?.(referralCode.toUpperCase());
  };

  return (
    <div className={styles.banner}>
      <Flex vertical gap="small">
        {/* Main Banner */}
        <Flex align="center" justify="space-between">
          <Flex align="center" gap="small">
            <Icon icon={Gift} size={24} style={{ color: theme.colorPrimary }} />
            <div>
              <Text strong style={{ display: 'block' }}>
                {isApplied ? '🎉 Bonus Upgraded!' : 'Have a Referral Code?'}
              </Text>
              <Text type="secondary" style={{ fontSize: 12 }}>
                {isApplied ? 'You\'ll receive 10 USDT + 200 KAWAI' : 'Get extra USDT + KAWAI bonus'}
              </Text>
            </div>
          </Flex>
          {!isApplied && (
            <Button
              type="text"
              size="small"
              className={styles.expandButton}
              onClick={() => setIsExpanded(!isExpanded)}
              icon={isExpanded ? <ChevronUp size={14} /> : <ChevronDown size={14} />}
            >
              {isExpanded ? 'Hide' : 'Enter Code'}
            </Button>
          )}
        </Flex>

        {/* Expanded Input */}
        {isExpanded && !isApplied && (
          <Flex gap="small" style={{ marginTop: 8 }}>
            <Input
              placeholder="Enter 6-digit code"
              value={referralCode}
              onChange={(e) => setReferralCode(e.target.value.toUpperCase())}
              maxLength={6}
              style={{ textTransform: 'uppercase' }}
            />
            <Button type="primary" onClick={handleApplyCode}>
              Apply
            </Button>
          </Flex>
        )}

        {/* Bonus Display */}
        {isApplied && (
          <Flex vertical gap="small" style={{ marginTop: 8 }}>
            <Flex align="center" gap="small">
              <Text className={styles.bonusAmount}>10 USDT</Text>
              <Text className={styles.comparisonText}>5 USDT</Text>
              <Text type="success" style={{ fontSize: 12, fontWeight: 600 }}>
                +100%
              </Text>
            </Flex>
            <Flex align="center" gap="small">
              <Text className={styles.bonusAmount}>200 KAWAI</Text>
              <Text className={styles.comparisonText}>100 KAWAI</Text>
              <Text type="success" style={{ fontSize: 12, fontWeight: 600 }}>
                +100%
              </Text>
            </Flex>
          </Flex>
        )}
      </Flex>
    </div>
  );
});

ReferralBanner.displayName = 'ReferralBanner';

