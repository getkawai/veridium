'use client';

import { DOCUMENTS_REFER_URL, PRIVACY_URL, TERMS_URL } from '@/const';
import { Button, Text, Icon, CopyButton } from '@lobehub/ui';
import { LobeHub } from '@lobehub/ui/brand';
import { Col, Flex, Row, Input, Space, Divider, message, Modal, Select, Progress } from 'antd';
import { createStyles } from 'antd-style';
import { Key, PlusCircle, Unlock, HardDrive, ArrowRight, Download, FileUp, Trash2, Wallet, AlertTriangle } from 'lucide-react';
import { memo, useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import type { WalletInfo } from '@@/github.com/kawai-network/veridium/internal/services/models';
import { Service as LocalFsService } from '@@/github.com/kawai-network/veridium/pkg/localfs';
import { Dialogs, Browser } from '@wailsio/runtime';

import BrandWatermark from '@/components/BrandWatermark';
import { useUserStore } from '@/store/user';
import { ReferralBanner } from '@/features/Referral/ReferralBanner';

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
  walletItem: css`
    padding: 12px 16px;
    border: 1px solid ${token.colorBorder};
    border-radius: 8px;
    transition: all 0.2s;
    &:hover {
      background: ${token.colorBgTextHover};
    }
  `,
  activeWallet: css`
    border-color: ${token.colorPrimary};
    background: ${token.colorPrimaryBg};
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
  backupWarning: css`
    padding: 16px;
    background: ${token.colorWarningBg};
    border: 1px solid ${token.colorWarningBorder};
    border-radius: 8px;
    margin-bottom: 16px;
  `,
}));

type Step = 'welcome' | 'setup' | 'mnemonic' | 'import' | 'unlock' | 'manage' | 'importKeystore' | 'addWallet';

export default memo(() => {
  const { styles } = useStyles();
  const { t } = useTranslation('clerk');

  const [step, setStep] = useState<Step>('welcome');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [mnemonic, setMnemonic] = useState('');
  const [description, setDescription] = useState('');
  const [keystoreJSON, setKeystoreJSON] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [showBackupReminder, setShowBackupReminder] = useState(false);
  const [selectedWallet, setSelectedWallet] = useState<string>('');
  const [referralCode, setReferralCode] = useState<string>('');
  const [hasReferral, setHasReferral] = useState(false);

  const {
    hasWallet,
    wallets,
    refreshWalletStatus,
    unlockWallet,
    setupWallet,
    generateMnemonic,
    isWalletLoaded,
    createWallet,
    switchWallet,
    deleteWallet,
    exportKeystore,
    importKeystore,
  } = useUserStore(s => ({
    hasWallet: s.hasWallet,
    wallets: s.wallets,
    refreshWalletStatus: s.refreshWalletStatus,
    unlockWallet: s.unlockWallet,
    setupWallet: s.setupWallet,
    generateMnemonic: s.generateMnemonic,
    isWalletLoaded: s.isWalletLoaded,
    createWallet: s.createWallet,
    switchWallet: s.switchWallet,
    deleteWallet: s.deleteWallet,
    exportKeystore: s.exportKeystore,
    importKeystore: s.importKeystore,
  }));

  useEffect(() => {
    refreshWalletStatus();
    
    // Check for referral code in URL
    const urlParams = new URLSearchParams(window.location.search);
    const refCode = urlParams.get('ref');
    if (refCode) {
      setReferralCode(refCode.toUpperCase());
      setHasReferral(true);
    }
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

  const resetForm = () => {
    setPassword('');
    setConfirmPassword('');
    setMnemonic('');
    setDescription('');
    setKeystoreJSON('');
  };

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
    // Validate mnemonic
    const words = mnemonic.trim().split(/\s+/);
    if (words.length !== 12 && words.length !== 24) {
      message.error(t('invalidMnemonic', { defaultValue: 'Mnemonic must be 12 or 24 words' }));
      return;
    }

    // Basic word validation (lowercase letters only)
    const invalidWords = words.filter(word => !/^[a-z]+$/.test(word));
    if (invalidWords.length > 0) {
      message.error(t('invalidMnemonicWords', { defaultValue: 'Mnemonic contains invalid words' }));
      return;
    }

    if (!description.trim()) {
      message.error('Please enter a wallet name');
      return;
    }

    setIsLoading(true);
    try {
      if (hasWallet) {
        await createWallet(password, mnemonic.trim(), description);
      } else {
        await setupWallet(password, mnemonic.trim(), description);
      }
      message.success(t('setupSuccess', { defaultValue: 'Wallet setup complete!' }));
      setShowBackupReminder(true);
      resetForm();
    } catch (err) {
      message.error(t('setupFailed', { defaultValue: 'Failed to setup wallet' }));
    } finally {
      setIsLoading(false);
    }
  };

  const handleImportKeystore = async () => {
    if (!keystoreJSON || !password) {
      message.error('Please provide keystore JSON and password');
      return;
    }
    setIsLoading(true);
    try {
      await importKeystore(keystoreJSON, password, description);
      message.success('Keystore imported successfully!');
      resetForm();
      setStep('unlock');
    } catch (err: any) {
      message.error(err?.message || 'Failed to import keystore');
    } finally {
      setIsLoading(false);
    }
  };

  const handleExportKeystore = async (address: string) => {
    try {
      const json = await exportKeystore(address);
      const blob = new Blob([json], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `keystore-${address}.json`;
      a.click();
      URL.revokeObjectURL(url);
      message.success('Keystore exported!');
    } catch (err) {
      message.error('Failed to export keystore');
    }
  };

  const handleDeleteWallet = async (address: string) => {
    Modal.confirm({
      title: 'Delete Wallet',
      content: 'Are you sure you want to delete this wallet? This action cannot be undone.',
      okText: 'Delete',
      okType: 'danger',
      onOk: async () => {
        const success = await deleteWallet(address);
        if (success) {
          message.success('Wallet deleted');
        } else {
          message.error('Failed to delete wallet');
        }
      }
    });
  };

  const handleSwitchWallet = async () => {
    if (!selectedWallet || !password) {
      message.error('Please select a wallet and enter password');
      return;
    }
    setIsLoading(true);
    const success = await switchWallet(selectedWallet, password);
    setIsLoading(false);
    if (success) {
      message.success('Wallet switched');
      setPassword('');
      setSelectedWallet('');
    } else {
      message.error('Invalid password');
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
            <Text className={styles.title} style={{ marginTop: 24 }}>Welcome To Kawai DeAI Network</Text>
            <Text className={styles.description}>
              {hasReferral 
                ? '🎉 Get 10 USDT + 200 KAWAI FREE!' 
                : 'Get 5 USDT + 100 KAWAI FREE'}
            </Text>
            <Text type="secondary" style={{ fontSize: 12, marginBottom: 16, textAlign: 'center' }}>
              No credit card • No email • Instant access
            </Text>
            <Button type="primary" size="large" onClick={() => setStep('setup')} block>
              Setup wallet & Claim Bonus
            </Button>
          </Flex>
        );

      case 'setup':
        return (
          <Flex vertical gap="large">
            {/* Referral Banner */}
            {!hasReferral && (
              <ReferralBanner 
                onReferralApplied={(code) => {
                  setReferralCode(code);
                  setHasReferral(true);
                }}
              />
            )}
            
            {/* Bonus Display */}
            {hasReferral && (
              <div style={{ 
                padding: '12px 16px', 
                background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
                borderRadius: 8,
                textAlign: 'center'
              }}>
                <Text strong style={{ color: 'white', display: 'block', fontSize: 14 }}>
                  🎉 Referral Applied: {referralCode}
                </Text>
                <Text style={{ color: 'rgba(255,255,255,0.9)', fontSize: 12 }}>
                  You'll receive 10 USDT + 200 KAWAI (instead of 5 USDT + 100 KAWAI)
                </Text>
              </div>
            )}

            <div style={{ textAlign: 'center' }}>
              <Text className={styles.title}>Setup password</Text>
              <Text as="p" type="secondary">Create a secure password for your wallet</Text>
            </div>
            <Input.Password
              placeholder="Enter password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              size="large"
              autoFocus
            />
            {password && (
              <div style={{ marginTop: -8, marginBottom: 8 }}>
                <Progress
                  percent={(() => {
                    let score = 0;
                    if (password.length >= 8) score += 25;
                    if (password.length >= 12) score += 15;
                    if (/[A-Z]/.test(password)) score += 20;
                    if (/[0-9]/.test(password)) score += 20;
                    if (/[^A-Za-z0-9]/.test(password)) score += 20;
                    return Math.min(100, score);
                  })()}
                  strokeColor={(() => {
                    let score = 0;
                    if (password.length >= 8) score += 25;
                    if (password.length >= 12) score += 15;
                    if (/[A-Z]/.test(password)) score += 20;
                    if (/[0-9]/.test(password)) score += 20;
                    if (/[^A-Za-z0-9]/.test(password)) score += 20;
                    const s = Math.min(100, score);
                    if (s < 40) return '#ff4d4f';
                    if (s < 60) return '#faad14';
                    if (s < 80) return '#52c41a';
                    return '#1890ff';
                  })()}
                  showInfo={false}
                  size="small"
                />
                <Text type="secondary" style={{ fontSize: 11 }}>
                  {(() => {
                    let score = 0;
                    if (password.length >= 8) score += 25;
                    if (password.length >= 12) score += 15;
                    if (/[A-Z]/.test(password)) score += 20;
                    if (/[0-9]/.test(password)) score += 20;
                    if (/[^A-Za-z0-9]/.test(password)) score += 20;
                    const s = Math.min(100, score);
                    if (s < 40) return 'Weak - Add more characters';
                    if (s < 60) return 'Fair - Add uppercase or numbers';
                    if (s < 80) return 'Good - Add special characters';
                    return 'Strong password';
                  })()}
                </Text>
              </div>
            )}
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
              <Flex align="center" gap="small" className={styles.methodItem} onClick={() => setStep('importKeystore')}>
                <Icon icon={HardDrive} size={24} />
                <div style={{ flex: 1 }}>
                  <Text strong style={{ display: 'block' }}>Import Keystore</Text>
                  <Text type="secondary" style={{ fontSize: 12 }}>Import from MetaMask or other wallets</Text>
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
            <div className={styles.backupWarning}>
              <Flex gap="small" align="start">
                <Icon icon={AlertTriangle} size={20} style={{ color: '#faad14', flexShrink: 0 }} />
                <Text type="warning" style={{ fontSize: 13 }}>
                  Never share your mnemonic! Anyone with these words can access your funds.
                </Text>
              </Flex>
            </div>
            <div style={{ position: 'relative', padding: '24px 16px', background: 'rgba(0,0,0,0.05)', borderRadius: 8, wordBreak: 'break-word', textAlign: 'center' }}>
              <Text code strong style={{ fontSize: 16 }}>{mnemonic}</Text>
              <CopyButton
                content={mnemonic}
                size="small"
                style={{ position: 'absolute', top: 4, right: 4 }}
              />
            </div>
            <Input
              placeholder="Wallet name (e.g. My Main Wallet)"
              value={description}
              onChange={e => setDescription(e.target.value)}
            />
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
            <Input
              placeholder="Wallet name (e.g. My Main Wallet)"
              value={description}
              onChange={e => setDescription(e.target.value)}
            />
            <Button type="primary" block size="large" onClick={handleFinishSetup} loading={isLoading}>
              Import Wallet
            </Button>
            <Button type="link" onClick={() => setStep('setup')}>Back</Button>
          </Flex>
        );

      case 'importKeystore':
        return (
          <Flex vertical gap="large">
            <div style={{ textAlign: 'center' }}>
              <Text className={styles.title}>Import Keystore</Text>
              <Text as="p" type="secondary">Paste your keystore JSON or upload file</Text>
            </div>
            <Button
              block
              size="large"
              icon={<FileUp size={18} />}
              onClick={async () => {
                try {
                  const result = await Dialogs.OpenFile({
                    CanChooseFiles: true,
                    CanChooseDirectories: false,
                    AllowsMultipleSelection: false,
                    Filters: [
                      {
                        DisplayName: 'Keystore JSON',
                        Pattern: '*.json',
                      },
                    ],
                    Title: 'Select Keystore File',
                  });

                  if (!result) return;

                  const filePath = Array.isArray(result) ? result[0] : result;
                  const fileResult = await LocalFsService.ReadFile({ path: filePath });
                  if (!fileResult || !fileResult.content) {
                    throw new Error('Failed to read file');
                  }
                  setKeystoreJSON(fileResult.content);
                  message.success('Keystore file loaded');
                } catch (error) {
                  console.error('Failed to load keystore:', error);
                  message.error('Failed to load keystore file');
                }
              }}
            >
              Select Keystore File
            </Button>
            <Input.TextArea
              rows={4}
              placeholder="Or paste keystore JSON here..."
              value={keystoreJSON}
              onChange={e => setKeystoreJSON(e.target.value)}
            />
            <Input.Password
              placeholder="Keystore password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              size="large"
            />
            <Input
              placeholder="Wallet description (optional)"
              value={description}
              onChange={e => setDescription(e.target.value)}
            />
            <Button type="primary" block size="large" onClick={handleImportKeystore} loading={isLoading}>
              Import Keystore
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
            {wallets.length > 1 && (
              <Select
                style={{ width: '100%' }}
                placeholder="Select wallet"
                value={selectedWallet || wallets.find(w => w.isActive)?.address}
                onChange={(val) => setSelectedWallet(val)}
                options={wallets.map(w => ({
                  label: `${w.description || 'Wallet'} (${w.address.slice(0, 8)}...)`,
                  value: w.address,
                }))}
              />
            )}
            <Input.Password
              placeholder="Enter password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              size="large"
              onPressEnter={selectedWallet ? handleSwitchWallet : handleUnlock}
              style={{ width: '100%' }}
              autoFocus
            />
            <Button
              type="primary"
              size="large"
              icon={<Unlock size={16} />}
              onClick={selectedWallet && selectedWallet !== wallets.find(w => w.isActive)?.address ? handleSwitchWallet : handleUnlock}
              loading={isLoading}
              block
            >
              Unlock
            </Button>
            {wallets.length > 0 && (
              <Flex gap="small">
                <Button type="link" onClick={() => setStep('manage')}>
                  <Icon icon={Wallet} size={14} style={{ marginRight: 4 }} /> Manage Wallets
                </Button>
                <Button type="link" onClick={() => { resetForm(); setStep('setup'); }}>
                  <Icon icon={PlusCircle} size={14} style={{ marginRight: 4 }} /> Add Wallet
                </Button>
              </Flex>
            )}
          </Flex>
        );

      case 'manage':
        return (
          <Flex vertical gap="large">
            <div style={{ textAlign: 'center' }}>
              <Text className={styles.title}>Manage Wallets</Text>
              <Text as="p" type="secondary">View, export, or delete your wallets</Text>
            </div>
            <Space direction="vertical" style={{ width: '100%' }}>
              {wallets.map((wallet: WalletInfo) => (
                <Flex
                  key={wallet.address}
                  className={`${styles.walletItem} ${wallet.isActive ? styles.activeWallet : ''}`}
                  justify="space-between"
                  align="center"
                >
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <Text strong style={{ display: 'block' }}>{wallet.description || 'Wallet'}</Text>
                    <Text type="secondary" style={{ fontSize: 12 }}>{wallet.address}</Text>
                    {wallet.isActive && <Text type="success" style={{ fontSize: 11, marginLeft: 8 }}>● Active</Text>}
                  </div>
                  <Flex gap="small">
                    <Button
                      size="small"
                      icon={<Download size={14} />}
                      onClick={() => handleExportKeystore(wallet.address)}
                    />
                    {!wallet.isActive && (
                      <Button
                        size="small"
                        danger
                        icon={<Trash2 size={14} />}
                        onClick={() => handleDeleteWallet(wallet.address)}
                      />
                    )}
                  </Flex>
                </Flex>
              ))}
            </Space>
            <Divider />
            <Button type="dashed" block onClick={() => { resetForm(); setStep('setup'); }}>
              <Icon icon={PlusCircle} size={14} style={{ marginRight: 4 }} /> Add New Wallet
            </Button>
            <Button type="link" onClick={() => setStep('unlock')}>Back to Unlock</Button>
          </Flex>
        );
    }
  };

  return (
    <>
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
                  <Button onClick={() => Browser.OpenURL(btn.href)} key={btn.id} type="text" size="small" style={{ color: 'inherit' }}>
                    {btn.label}
                  </Button>
                ))}
              </Flex>
            </Col>
          </Row>
        </div>
      </div>

      {/* Backup Reminder Modal */}
      <Modal
        open={showBackupReminder}
        title="⚠️ Backup Reminder"
        onOk={() => setShowBackupReminder(false)}
        onCancel={() => setShowBackupReminder(false)}
        okText="I understand"
        cancelButtonProps={{ style: { display: 'none' } }}
      >
        <div className={styles.backupWarning}>
          <Space direction="vertical">
            <Text strong>Please backup your mnemonic phrase!</Text>
            <Text type="secondary">
              - Store it in a secure location<br />
              - Never share it with anyone<br />
              - Consider using a hardware wallet for large amounts
            </Text>
          </Space>
        </div>
      </Modal>
    </>
  );
});
