import { Tabs, Button } from 'antd';
import { useState, useRef, useCallback } from 'react';
import { Coins, Award, Users, RefreshCw, TrendingUp } from 'lucide-react';
import { Flexbox } from 'react-layout-kit';
import { MiningRewardsSection } from './components/rewards/MiningRewardsSection';
import { CashbackRewardsSection } from './components/rewards/CashbackRewardsSection';
import { ReferralRewardsSection } from './components/rewards/ReferralRewardsSection';
import { RevenueShareSection } from './components/rewards/RevenueShareSection';
import type { RewardsContentProps } from './types';

const RewardsContent = ({ styles, theme, currentNetwork, transactions, setModalType }: RewardsContentProps) => {
  const [activeTab, setActiveTab] = useState<'mining' | 'cashback' | 'referral' | 'revenue'>('mining');
  const [refreshKey, setRefreshKey] = useState(0);
  
  // Refs to hold refresh callbacks from each section
  const miningRefreshRef = useRef<(() => void) | null>(null);
  const cashbackRefreshRef = useRef<(() => void) | null>(null);
  const referralRefreshRef = useRef<(() => void) | null>(null);
  const revenueRefreshRef = useRef<(() => void) | null>(null);

  const handleRefresh = useCallback(() => {
    // Call the active tab's refresh callback
    switch (activeTab) {
      case 'mining':
        miningRefreshRef.current?.();
        break;
      case 'cashback':
        cashbackRefreshRef.current?.();
        break;
      case 'referral':
        referralRefreshRef.current?.();
        break;
      case 'revenue':
        revenueRefreshRef.current?.();
        break;
    }
  }, [activeTab]);

  const handleOpenDepositModal = useCallback(() => {
    setModalType?.('deposit');
  }, [setModalType]);

  return (
    <Flexbox style={{ maxWidth: 1200, width: '100%' }} gap={20}>
      {/* Header */}
      <Flexbox horizontal justify="space-between" align="center">
        <div>
          <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>Rewards</h2>
          <span style={{ color: theme.colorTextSecondary, fontSize: 13 }}>
            Claim your KAWAI rewards from mining, deposits, referrals, and revenue sharing
          </span>
        </div>
        <Button
          icon={<RefreshCw size={16} />}
          onClick={handleRefresh}
          size="small"
        >
          Refresh
        </Button>
      </Flexbox>

      {/* Tabs */}
      <Tabs
        activeKey={activeTab}
        onChange={(key) => setActiveTab(key as 'mining' | 'cashback' | 'referral' | 'revenue')}
        size="large"
        items={[
          {
            key: 'mining',
            label: (
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <Coins size={16} />
                Mining Rewards
              </span>
            ),
            children: (
              <MiningRewardsSection
                key={`mining-${refreshKey}`}
                currentNetwork={currentNetwork}
                theme={theme}
                styles={styles}
                onRefresh={(refreshFn) => { miningRefreshRef.current = refreshFn; }}
              />
            ),
          },
          {
            key: 'cashback',
            label: (
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <Award size={16} />
                Deposit Cashback
              </span>
            ),
            children: (
              <CashbackRewardsSection
                key={`cashback-${refreshKey}`}
                currentNetwork={currentNetwork}
                theme={theme}
                styles={styles}
                onOpenDepositModal={handleOpenDepositModal}
                onRefresh={(refreshFn) => { cashbackRefreshRef.current = refreshFn; }}
              />
            ),
          },
          {
            key: 'referral',
            label: (
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <Users size={16} />
                Referral Rewards
              </span>
            ),
            children: (
              <ReferralRewardsSection
                key={`referral-${refreshKey}`}
                theme={theme}
                styles={styles}
                onRefresh={(refreshFn) => { referralRefreshRef.current = refreshFn; }}
              />
            ),
          },
          {
            key: 'revenue',
            label: (
              <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <TrendingUp size={16} />
                Revenue Share
              </span>
            ),
            children: (
              <RevenueShareSection
                key={`revenue-${refreshKey}`}
                currentNetwork={currentNetwork}
                theme={theme}
                styles={styles}
                onRefresh={(refreshFn) => { revenueRefreshRef.current = refreshFn; }}
              />
            ),
          },
        ]}
      />
    </Flexbox>
  );
};

export default RewardsContent;
