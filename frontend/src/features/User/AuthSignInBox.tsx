'use client';

import { DOCUMENTS_REFER_URL, PRIVACY_URL, TERMS_URL } from '@/const';
import { Button, Text, Icon, CopyButton } from '@lobehub/ui';
import { LobeHub } from '@lobehub/ui/brand';
import { Col, Flex, Row, Input, Space, Divider, message } from 'antd';
import { createStyles } from 'antd-style';
import { Key, PlusCircle, Unlock, HardDrive, ArrowRight } from 'lucide-react';
import { memo, useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

import BrandWatermark from '@/components/BrandWatermark';
import { useUserStore } from '@/store/user';

const useStyles = createStyles(({ css, token }) => ({
  container: css`
    min-width: 400px;
    border: 1px solid ${token.colorBorder};
    border-radius: ${token.borderRadiusLG}px;
    background: ${token.colorBgContainer};
    overflow: hidden;
  `,
  contentCard: css`
    padding: 32px;
  `,
  footer: css`
    padding: 16px;
    border-top: 1px solid ${token.colorBorder};
    background: ${token.colorBgElevated};
    color: ${token.colorTextDescription};
  `,
  methodItem: css`
    padding: 16px;
    border: 1px solid ${token.colorBorder};
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    &:hover {
      background: ${token.colorBgTextHover};
      border-color: ${token.colorPrimary};
    }
  `,
  title: css`
    margin-bottom: 8px;
    font-size: 20px;
    font-weight: 600;
    display: block;
  `,
  description: css`
    color: ${token.colorTextSecondary};
    margin-bottom: 24px;
    text-align: center;
    display: block;
  `,
}));

type Step = 'welcome' | 'setup' | 'mnemonic' | 'import' | 'unlock';

export default memo(() => {
  const { styles } = useStyles();
  const { t } = useTranslation('auth');

  const [step, setStep] = useState<Step>('welcome');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [mnemonic, setMnemonic] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const {
    hasWallet,
    refreshWalletStatus,
    unlockWallet,
    setupWallet,
    generateMnemonic,
    isWalletLoaded
  } = useUserStore(s => ({
    hasWallet: s.hasWallet,
    refreshWalletStatus: s.refreshWalletStatus,
    unlockWallet: s.unlockWallet,
    setupWallet: s.setupWallet,
    generateMnemonic: s.generateMnemonic,
    isWalletLoaded: s.isWalletLoaded
  }));

  useEffect(() => {
    refreshWalletStatus();
  }, []);

  useEffect(() => {
    if (isWalletLoaded) {
      if (hasWallet) {
        setStep('unlock');
      } else {
        setStep('welcome');
      }
    }
  }, [isWalletLoaded, hasWallet]);

  const handleUnlock = async () => {
    setIsLoading(true);
    const success = await unlockWallet(password);
    setIsLoading(false);
    if (!success) {
      message.error(t('unlockFailed', { defaultValue: 'Invalid password' }));
    }
  };

  const handleSetup = async (method: 'generate' | 'import') => {
    if (password !== confirmPassword) {
      message.error(t('passwordMismatch', { defaultValue: 'Passwords do not match' }));
      return;
    }
    if (password.length < 8) {
      message.error(t('passwordTooShort', { defaultValue: 'Password must be at least 8 characters' }));
      return;
    }

    if (method === 'generate') {
      const phrase = await generateMnemonic();
      setMnemonic(phrase);
      setStep('mnemonic');
    } else {
      setStep('import');
    }
  };

  const handleFinishSetup = async () => {
    setIsLoading(true);
    try {
      await setupWallet(password, mnemonic);
      message.success(t('setupSuccess', { defaultValue: 'Wallet setup complete!' }));
    } catch (err) {
      message.error(t('setupFailed', { defaultValue: 'Failed to setup wallet' }));
    } finally {
      setIsLoading(false);
    }
  };

  const footerBtns = [
    { href: DOCUMENTS_REFER_URL, id: 0, label: t('footerPageLink__help') },
    { href: PRIVACY_URL, id: 1, label: t('footerPageLink__privacy') },
    { href: TERMS_URL, id: 2, label: t('footerPageLink__terms') },
  ];

  const renderContent = () => {
    switch (step) {
      case 'welcome':
        return (
          <Flex vertical align="center">
            <LobeHub size={80} />
            <Text className={styles.title} style={{ marginTop: 24 }}>Welcome To OnChain Wallet</Text>
            <Text className={styles.description}>Your Gateway to Decentralized World</Text>
            <Button type="primary" size="large" onClick={() => setStep('setup')} block>
              Setup wallet
            </Button>
          </Flex>
        );

      case 'setup':
        return (
          <Flex vertical gap="large">
            <div style={{ textAlign: 'center' }}>
              <Text className={styles.title}>Setup password</Text>
              <Text as="p" type="secondary">Input your wallet password</Text>
            </div>
            <Input.Password
              placeholder="Enter password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              size="large"
            />
            <Input.Password
              placeholder="Confirm password"
              value={confirmPassword}
              onChange={e => setConfirmPassword(e.target.value)}
              size="large"
            />
            <Divider plain>Choose Method</Divider>
            <Space direction="vertical" style={{ width: '100%' }}>
              <Flex align="center" gap="small" className={styles.methodItem} onClick={() => handleSetup('import')}>
                <Icon icon={PlusCircle} size={24} />
                <div style={{ flex: 1 }}>
                  <Text strong style={{ display: 'block' }}>Use Existing Mnemonic</Text>
                  <Text type="secondary" style={{ fontSize: 12 }}>Import an existing 12-24 word recovery phrase</Text>
                </div>
                <Icon icon={ArrowRight} />
              </Flex>
              <Flex align="center" gap="small" className={styles.methodItem} onClick={() => handleSetup('generate')}>
                <Icon icon={Key} size={24} />
                <div style={{ flex: 1 }}>
                  <Text strong style={{ display: 'block' }}>Generate New Mnemonic</Text>
                  <Text type="secondary" style={{ fontSize: 12 }}>Create a new wallet with a 12-24 word phrase</Text>
                </div>
                <Icon icon={ArrowRight} />
              </Flex>
              <Flex align="center" gap="small" className={styles.methodItem}>
                <Icon icon={HardDrive} size={24} />
                <div style={{ flex: 1 }}>
                  <Text strong style={{ display: 'block' }}>Restore backup</Text>
                  <Text type="secondary" style={{ fontSize: 12 }}>Recover wallet from backup file</Text>
                </div>
                <Icon icon={ArrowRight} />
              </Flex>
            </Space>
          </Flex>
        );

      case 'mnemonic':
        return (
          <Flex vertical gap="large">
            <div style={{ textAlign: 'center' }}>
              <Text className={styles.title}>Backup Mnemonic</Text>
              <Text as="p" type="secondary">Write down these 12 words in order</Text>
            </div>
            <div style={{ position: 'relative', padding: '24px 16px', background: 'rgba(0,0,0,0.05)', borderRadius: 8, wordBreak: 'break-word', textAlign: 'center' }}>
              <Text code strong style={{ fontSize: 16 }}>{mnemonic}</Text>
              <CopyButton
                content={mnemonic}
                size="small"
                style={{ position: 'absolute', top: 4, right: 4 }}
              />
            </div>
            <Button type="primary" block size="large" onClick={handleFinishSetup} loading={isLoading}>
              I have written it down
            </Button>
            <Button type="link" onClick={() => setStep('setup')}>Back</Button>
          </Flex>
        );

      case 'import':
        return (
          <Flex vertical gap="large">
            <div style={{ textAlign: 'center' }}>
              <Text className={styles.title}>Import Mnemonic</Text>
              <Text as="p" type="secondary">Enter your 12-24 word phrase</Text>
            </div>
            <Input.TextArea
              rows={4}
              placeholder="word1 word2 ..."
              value={mnemonic}
              onChange={e => setMnemonic(e.target.value)}
            />
            <Button type="primary" block size="large" onClick={handleFinishSetup} loading={isLoading}>
              Import Wallet
            </Button>
            <Button type="link" onClick={() => setStep('setup')}>Back</Button>
          </Flex>
        );

      case 'unlock':
        return (
          <Flex vertical align="center" gap="large">
            <LobeHub size={80} />
            <div style={{ textAlign: 'center' }}>
              <Text className={styles.title}>Welcome Back</Text>
              <Text as="p" type="secondary">Unlock your wallet with password</Text>
            </div>
            <Input.Password
              placeholder="Enter password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              size="large"
              onPressEnter={handleUnlock}
              style={{ width: '100%' }}
            />
            <Button
              type="primary"
              size="large"
              icon={<Unlock size={16} />}
              onClick={handleUnlock}
              loading={isLoading}
              block
            >
              Unlock
            </Button>
          </Flex>
        );
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.contentCard}>
        {renderContent()}
      </div>
      <div className={styles.footer}>
        <Row align="middle">
          <Col span={12}>
            <BrandWatermark />
          </Col>
          <Col span={12}>
            <Flex justify="end" gap="small">
              {footerBtns.map(btn => (
                <Button key={btn.id} type="text" size="small" style={{ color: 'inherit' }}>
                  {btn.label}
                </Button>
              ))}
            </Flex>
          </Col>
        </Row>
      </div>
    </div>
  );
});
