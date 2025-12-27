'use client';

import { useAutoAnimate } from '@formkit/auto-animate/react';
import { useSize } from 'ahooks';
import { ActionIcon, Avatar, Text } from '@lobehub/ui';
import { Popover, App } from 'antd';
import { createStyles, useTheme } from 'antd-style';
import { Plus, CheckCircle2, Copy } from 'lucide-react';
import { memo, useEffect, useRef, useState } from 'react';
import { Flexbox } from 'react-layout-kit';

import { genAvatar } from '@/utils/avatar';

import { WalletService } from '@@/github.com/kawai-network/veridium/internal/services';
import type { WalletInfo } from '@@/github.com/kawai-network/veridium/internal/services';

const useStyles = createStyles(({ css, token }) => ({
  accountItem: css`
    cursor: pointer;
    transition: all 0.2s ease;
    border-radius: 8px;
    width: 100%;
    
    &:hover {
      background: ${token.colorFillTertiary};
    }
  `,
  activeItem: css`
    background: ${token.colorFillSecondary};
  `,
  avatar: css`
    flex: none;
    pointer-events: none;
  `,
  header: css`
    margin-bottom: 12px;
    padding-inline: 4px;
    width: 100%;
  `
}));

interface AccountListProps {
  activeAddress: string;
  onAccountSwitch: (address: string) => void;
  onAddAccount: () => void;
}

const AccountList = memo<AccountListProps>(({ activeAddress, onAccountSwitch, onAddAccount }) => {
  const { styles, cx } = useStyles();
  const theme = useTheme();
  const { message } = App.useApp();
  const [parent] = useAutoAnimate();
  const [wallets, setWallets] = useState<WalletInfo[]>([]);

  const ref = useRef(null);
  const size = useSize(ref);
  const width = size?.width || 80;
  const showMoreInfo = Boolean(width > 120);

  const handleCopy = (e: React.MouseEvent, address: string) => {
    e.stopPropagation();
    navigator.clipboard.writeText(address);
    message.success('Address copied!');
  };

  const fetchWallets = async () => {
    try {
      const status = await WalletService.GetStatus();
      setWallets(status.wallets || []);
    } catch (e) {
      console.error('Failed to fetch wallets', e);
    }
  };

  useEffect(() => {
    fetchWallets();
  }, [activeAddress]);

  return (
    <Flexbox
      align="center"
      gap={12}
      paddingInline={8}
      ref={ref}
      style={{
        maxHeight: '100%',
        overflowY: 'auto',
        paddingTop: 16,
      }}
      width={'100%'}
    >
      {/* Header / Add Button */}
      {showMoreInfo ? (
        <Flexbox align="center" horizontal justify="space-between" className={styles.header}>
          <Text fontSize={14} weight={500}>Accounts {wallets.length ? wallets.length : ''}</Text>
          <ActionIcon
            icon={Plus}
            onClick={onAddAccount}
            size={'small'}
            title="Add Account"
            tooltipProps={{ placement: 'left' }}
          />
        </Flexbox>
      ) : (
        <ActionIcon
          icon={Plus}
          onClick={onAddAccount}
          size={{ blockSize: 48, size: 20 }}
          title="Add Account"
          tooltipProps={{ placement: 'left' }}
          variant={'filled'}
          style={{ marginBottom: 12 }}
        />
      )}

      {/* List */}
      <Flexbox align="center" gap={12} ref={parent} width={'100%'}>
        {wallets.map((wallet, index) => {
          const isActive = wallet.address === activeAddress;

          const avatarInfo = genAvatar(wallet.address);

          const tooltipContent = (
            <Flexbox align={'center'} flex={1} gap={16} horizontal justify={'space-between'} style={{ overflow: 'hidden' }}>
              <Flexbox flex={1} style={{ overflow: 'hidden' }}>
                <Text ellipsis fontSize={14} weight={500}>
                  {wallet.description || 'Unnamed Account'}
                </Text>
                <Text ellipsis fontSize={12} type={'secondary'} style={{ fontFamily: 'monospace' }}>
                  {wallet.address.substring(0, 6)}...{wallet.address.substring(wallet.address.length - 4)}
                </Text>
              </Flexbox>
              <Flexbox horizontal gap={4} align="center">
                <ActionIcon
                  icon={Copy}
                  onClick={(e) => handleCopy(e, wallet.address)}
                  size="small"
                  title="Copy Address"
                />
                {isActive && <CheckCircle2 size={14} color={theme.colorSuccess} />}
              </Flexbox>
            </Flexbox>
          );

          const item = (
            <Flexbox
              align={'center'}
              gap={12}
              horizontal
              justify={'center'}
              onClick={() => !isActive && onAccountSwitch(wallet.address)}
              className={cx(styles.accountItem, isActive && styles.activeItem)}
              style={{
                cursor: 'pointer',
                padding: showMoreInfo ? '8px' : '0',
                paddingBottom: index === wallets.length - 1 ? 4 : 0,
              }}
              width={'100%'}
            >
              <Avatar
                avatar={avatarInfo.emoji}
                background={isActive ? theme.colorFillSecondary : undefined}
                bordered={isActive}
                shape="square"
                size={48}
                className={styles.avatar}
                style={{ fontSize: 24, background: !isActive ? avatarInfo.background : undefined }}
              />
              {showMoreInfo && tooltipContent}
            </Flexbox>
          );

          return (
            <Popover
              key={wallet.address}
              arrow={false}
              content={tooltipContent}
              placement={'left'}
              styles={{ body: { width: 200 } }}
              trigger={showMoreInfo ? [] : ['hover']}
            >
              {item}
            </Popover>
          );
        })}
      </Flexbox>
    </Flexbox>
  );
});

AccountList.displayName = 'AccountList';

export default AccountList;
