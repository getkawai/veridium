import { Card, Modal, Table, Tag, Button, Empty, Popover, Spin, Tooltip } from 'antd';
import { memo, useState } from 'react';
import { JarvisService } from '@@/github.com/kawai-network/veridium/internal/services';
import {
  History,
  Plus,
  Send,
  Eye,
  EyeOff,
  Gift,
  Repeat2,
  Fuel,
  ArrowDownToLine,
  Coins,
} from 'lucide-react';
import { Icon } from '@lobehub/ui';
import { Flexbox } from 'react-layout-kit';
import { useTheme } from 'antd-style';
import { NetworkIcon } from './NetworkIcons';
import { TokenUSDT } from '@web3icons/react';
import type { HomeContentProps } from './types';

// Transaction Link with analysis popup
const TransactionLink = memo<{ txHash: string; networkId?: number }>(({ txHash, networkId }) => {
  const theme = useTheme();
  const [analyzing, setAnalyzing] = useState(false);
  const [analysis, setAnalysis] = useState<any>(null);

  const handleAnalyze = async () => {
    if (!networkId || analyzing) return;
    setAnalyzing(true);
    try {
      const result = await JarvisService.AnalyzeTransaction(txHash, networkId);
      setAnalysis(result);
    } catch (e) {
      console.error('Failed to analyze transaction', e);
    } finally {
      setAnalyzing(false);
    }
  };

  const shortHash = `${txHash.substring(0, 6)}...${txHash.substring(txHash.length - 4)}`;

  return (
    <Popover
      trigger="click"
      onOpenChange={(open) => open && handleAnalyze()}
      content={
        <div style={{ width: 300, maxHeight: 400, overflowY: 'auto' }}>
          {analyzing ? (
            <Flexbox align="center" justify="center" style={{ padding: 20 }}>
              <Spin size="small" />
              <span style={{ marginLeft: 8 }}>Analyzing...</span>
            </Flexbox>
          ) : analysis ? (
            <Flexbox gap={12}>
              <div style={{ fontSize: 11, color: theme.colorTextTertiary }}>TRANSACTION ANALYSIS</div>

              <Flexbox gap={8}>
                <Flexbox horizontal justify="space-between">
                  <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Status</span>
                  <Tag color={analysis.status === 'done' ? 'green' : analysis.status === 'reverted' ? 'red' : 'orange'}>
                    {analysis.status}
                  </Tag>
                </Flexbox>

                <Flexbox horizontal justify="space-between">
                  <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Type</span>
                  <span style={{ fontSize: 12, fontWeight: 600 }}>{analysis.txType || 'Unknown'}</span>
                </Flexbox>

                {analysis.method && (
                  <Flexbox horizontal justify="space-between">
                    <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Method</span>
                    <Tag color="blue" style={{ fontFamily: 'monospace' }}>{analysis.method}</Tag>
                  </Flexbox>
                )}

                {analysis.value && analysis.value !== '0' && (
                  <Flexbox horizontal justify="space-between">
                    <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Value</span>
                    <span style={{ fontSize: 12 }}>{analysis.value}</span>
                  </Flexbox>
                )}

                {analysis.gasUsed && (
                  <Flexbox horizontal justify="space-between">
                    <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Gas Used</span>
                    <span style={{ fontSize: 12 }}>{parseInt(analysis.gasUsed).toLocaleString()}</span>
                  </Flexbox>
                )}

                {analysis.gasCost && (
                  <Flexbox horizontal justify="space-between">
                    <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Gas Cost</span>
                    <span style={{ fontSize: 12 }}>{analysis.gasCost}</span>
                  </Flexbox>
                )}

                {analysis.blockNumber > 0 && (
                  <Flexbox horizontal justify="space-between">
                    <span style={{ color: theme.colorTextSecondary, fontSize: 12 }}>Block</span>
                    <span style={{ fontSize: 12, fontFamily: 'monospace' }}>#{analysis.blockNumber.toLocaleString()}</span>
                  </Flexbox>
                )}
              </Flexbox>

              {/* Decoded Parameters */}
              {analysis.params && analysis.params.length > 0 && (
                <>
                  <div style={{ fontSize: 11, color: theme.colorTextTertiary, marginTop: 8 }}>PARAMETERS</div>
                  <Flexbox gap={4}>
                    {analysis.params.map((param: any, i: number) => (
                      <div key={i} style={{
                        padding: '4px 8px',
                        background: theme.colorFillTertiary,
                        borderRadius: 4,
                        fontSize: 11
                      }}>
                        <span style={{ color: theme.colorTextSecondary }}>{param.name}</span>
                        <span style={{ color: theme.colorTextTertiary }}> ({param.type})</span>
                        <div style={{ fontFamily: 'monospace', wordBreak: 'break-all', marginTop: 2 }}>
                          {param.value?.substring(0, 50)}{param.value?.length > 50 ? '...' : ''}
                        </div>
                      </div>
                    ))}
                  </Flexbox>
                </>
              )}

              {/* Event Logs */}
              {analysis.logs && analysis.logs.length > 0 && (
                <>
                  <div style={{ fontSize: 11, color: theme.colorTextTertiary, marginTop: 8 }}>EVENTS ({analysis.logs.length})</div>
                  <Flexbox gap={4}>
                    {analysis.logs.slice(0, 3).map((log: any, i: number) => (
                      <Tag key={i} color="purple">{log.name || 'Unknown Event'}</Tag>
                    ))}
                    {analysis.logs.length > 3 && (
                      <span style={{ fontSize: 11, color: theme.colorTextTertiary }}>
                        +{analysis.logs.length - 3} more
                      </span>
                    )}
                  </Flexbox>
                </>
              )}

              {analysis.error && (
                <div style={{ color: theme.colorError, fontSize: 12, marginTop: 8 }}>
                  Error: {analysis.error}
                </div>
              )}
            </Flexbox>
          ) : (
            <div style={{ padding: 16, textAlign: 'center', color: theme.colorTextSecondary }}>
              Click to analyze transaction
            </div>
          )}
        </div>
      }
    >
      <span
        style={{
          fontFamily: 'monospace',
          fontSize: 11,
          cursor: 'pointer',
          color: theme.colorPrimary,
          textDecoration: 'underline'
        }}
      >
        {shortHash}
      </span>
    </Popover>
  );
});

const HomeContent = ({
  address,
  balance,
  nativeBalance,
  kawaiBalance,
  balanceVisible,
  setBalanceVisible,
  setModalType,
  transactions,
  styles,
  theme,
  currentNetwork,
  gasEstimate,
  currentBlock,
  balancesLoading
}: HomeContentProps) => {
  const [showAllTx, setShowAllTx] = useState(false);

  return (
    <Flexbox style={{ maxWidth: 900, width: '100%' }} gap={20}>
      {/* Balance Card */}
      <Card className={styles.balanceCard}>
        <div className={styles.eyeButton} onClick={() => setBalanceVisible(!balanceVisible)}>
          {balanceVisible ? <Eye size={16} /> : <EyeOff size={16} />}
        </div>
        <Flexbox horizontal justify="space-between" align="center">
          <Flexbox style={{ flexDirection: 'column' }} gap={4}>
            <span style={{ fontSize: 11, color: theme.colorTextSecondary, textTransform: 'uppercase', letterSpacing: '0.5px' }}>Total Balance</span>
            <div className={styles.statValue}>
              {balancesLoading ? (
                <Spin size="small" />
              ) : (
                <>
                  {balanceVisible ? `$${balance}` : '••••••'}
                  <span style={{ fontSize: 16, color: theme.colorTextTertiary, marginLeft: 6, fontWeight: 500 }}>USDT</span>
                </>
              )}
            </div>
            {/* Native token balance */}
            {currentNetwork && (
              <div style={{ fontSize: 13, color: theme.colorTextSecondary, marginTop: 4 }}>
                {balanceVisible ? nativeBalance : '••••'} {currentNetwork.nativeTokenSymbol}
              </div>
            )}
          </Flexbox>
          {/* Network & Gas Info */}
          <Flexbox gap={8} align="flex-end">
            {gasEstimate && (
              <Tooltip title={`Max Tip: ${gasEstimate.maxTipGwei.toFixed(2)} Gwei`}>
                <Flexbox horizontal align="center" gap={4} style={{
                  padding: '4px 8px',
                  background: 'rgba(255,255,255,0.1)',
                  borderRadius: 8,
                  fontSize: 11
                }}>
                  <Fuel size={12} />
                  <span>{gasEstimate.maxGasPriceGwei.toFixed(1)} Gwei</span>
                </Flexbox>
              </Tooltip>
            )}
            {currentBlock > 0 && (
              <div style={{ fontSize: 10, color: theme.colorTextTertiary }}>
                Block #{currentBlock.toLocaleString()}
              </div>
            )}
          </Flexbox>
        </Flexbox>
      </Card>

      {/* Quick Actions */}
      <Flexbox horizontal gap={12} style={{ marginTop: 4 }}>
        {[
          { label: 'Deposit', icon: Plus, color: '#10b981', action: () => setModalType('deposit') },
          { label: 'Send', icon: Send, color: '#06b6d4', action: () => setModalType('send') },
          { label: 'Receive', icon: ArrowDownToLine, color: '#22c55e', action: () => setModalType('receive') },
          { label: 'Swap', icon: Repeat2, color: '#eab308', action: () => setModalType('swap') },
        ].map((item) => (
          <div key={item.label} className={styles.actionButton} onClick={item.action}>
            <div className={styles.actionCircle} style={{ background: `${item.color}20`, color: item.color }}>
              <item.icon size={24} />
            </div>
            <span style={{ fontWeight: 600, fontSize: 13 }}>{item.label}</span>
          </div>
        ))}
      </Flexbox>

      {/* Token List */}
      <Card
        title={<Flexbox horizontal align="center" gap={8}><Coins size={16} /> Tokens</Flexbox>}
        size="small"
        extra={
          <Button
            type="text"
            icon={<Plus size={14} />}
            size="small"
            onClick={() => setModalType('addToken')}
          >
            Add Token
          </Button>
        }
      >
        <Flexbox gap={8}>
          {/* Native Token */}
          {currentNetwork && (
            <div className={styles.tokenRow}>
              <Flexbox horizontal align="center" gap={12} style={{ flex: 1 }}>
                <div style={{
                  width: 36,
                  height: 36,
                  borderRadius: '50%',
                  background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  color: '#fff',
                  fontWeight: 800,
                  fontSize: 14,
                }}>
                  {currentNetwork && (
                    <NetworkIcon
                      name={currentNetwork.icon || 'ethereum'}
                      size={24}
                      variant="mono"
                      color="#fff"
                    />
                  )}
                </div>
                <div>
                  <div style={{ fontWeight: 600 }}>{currentNetwork.nativeTokenSymbol}</div>
                  <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>Native Token</div>
                </div>
              </Flexbox>
              <div style={{ textAlign: 'right', minWidth: 70 }}>
                <div style={{ fontWeight: 700, fontSize: 14, color: '#fff' }}>{balanceVisible ? nativeBalance : '••••'}</div>
              </div>
            </div>
          )}

          {/* USDT */}
          <div className={styles.tokenRow}>
            <Flexbox horizontal align="center" gap={12} style={{ flex: 1 }}>
              <div style={{
                width: 36,
                height: 36,
                borderRadius: '50%',
                background: '#26a17b',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: '#fff',
                fontWeight: 800,
                fontSize: 16,
                fontFamily: 'Arial, sans-serif',
                textShadow: '0 1px 2px rgba(0,0,0,0.2)'
              }}>
                <TokenUSDT size={36} variant="branded" />
              </div>
              <div>
                <div style={{ fontWeight: 600 }}>USDT</div>
                <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>Tether USD</div>
              </div>
            </Flexbox>
            <Flexbox horizontal align="center" gap={16}>
              <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>$1.00</span>
              <div style={{ textAlign: 'right', minWidth: 70 }}>
                <div style={{ fontWeight: 700, fontSize: 14, color: '#fff' }}>{balanceVisible ? balance : '••••'}</div>
                <div style={{ fontSize: 11, color: theme.colorTextTertiary }}>${balanceVisible ? balance : '••••'}</div>
              </div>
            </Flexbox>
          </div>

          {/* KAWAI */}
          <div className={styles.tokenRow}>
            <Flexbox horizontal align="center" gap={12} style={{ flex: 1 }}>
              <div style={{
                width: 36,
                height: 36,
                borderRadius: '50%',
                background: 'linear-gradient(135deg, #ff9a9e 0%, #fecfef 99%, #fecfef 100%)',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: '#fff',
                fontWeight: 800,
                fontSize: 14,
                fontFamily: 'Arial, sans-serif',
                boxShadow: '0 2px 8px rgba(255, 154, 158, 0.3)'
              }}>
                <Icon icon={Gift} size={20} color="#fff" />
              </div>
              <div>
                <div style={{ fontWeight: 600 }}>KAWAI</div>
                <div style={{ fontSize: 12, color: theme.colorTextSecondary }}>Kawai Token</div>
              </div>
            </Flexbox>
            <Flexbox horizontal align="center" gap={16}>
              <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>-</span>
              <div style={{ textAlign: 'right', minWidth: 70 }}>
                <div style={{ fontWeight: 700, fontSize: 14, color: '#fff' }}>{balanceVisible ? kawaiBalance : '••••'}</div>
                <div style={{ fontSize: 11, color: theme.colorTextTertiary }}>
                  {balanceVisible ? `${kawaiBalance} KAWAI` : '••••'}
                </div>
              </div>
            </Flexbox>
          </div>

        </Flexbox>
      </Card>

      {/* Activity */}
      <Card
        title={<Flexbox horizontal align="center" gap={8}><History size={16} /> Recent Activity</Flexbox>}
        size="small"
        extra={transactions.length > 5 && (
          <Button type="link" size="small" onClick={() => setShowAllTx(true)}>
            View All ({transactions.length})
          </Button>
        )}
      >
        {transactions.length > 0 ? (
          <Table
            dataSource={transactions.slice(0, 5)}
            rowKey="id"
            pagination={false}
            size="small"
            columns={[
              { title: 'Type', dataIndex: 'txType', key: 'txType', render: (type) => <Tag color={type === 'DEPOSIT' ? 'green' : 'blue'}>{type}</Tag> },
              { title: 'Amount', dataIndex: 'amount', key: 'amount', render: (amount, record: any) => <span style={{ color: record.txType === 'DEPOSIT' ? theme.colorSuccess : theme.colorText, fontWeight: 600 }}>{record.txType === 'DEPOSIT' ? '+' : '-'}{amount} USDT</span> },
              { title: 'Date', dataIndex: 'createdAt', key: 'createdAt', render: (date) => new Date(date).toLocaleDateString() },
              {
                title: 'TX',
                dataIndex: 'txHash',
                key: 'txHash',
                render: (txHash) => txHash ? (
                  <TransactionLink txHash={txHash} networkId={currentNetwork?.id} />
                ) : '-'
              },
            ]}
          />
        ) : (
          <Flexbox align="center" gap={16} style={{ padding: '24px 0' }}>
            <Empty description={false} image={Empty.PRESENTED_IMAGE_SIMPLE} />
            <span style={{ color: theme.colorTextSecondary }}>No transactions yet</span>
            <Button
              type="primary"
              size="small"
              onClick={() => window.open('https://testnet.monad.xyz/faucet', '_blank')}
            >
              Get Test Tokens (Faucet)
            </Button>
          </Flexbox>
        )}
      </Card>

      {/* Full Transaction History Modal */}
      <Modal
        title={<Flexbox horizontal align="center" gap={8}><History size={18} /> Transaction History</Flexbox>}
        open={showAllTx}
        onCancel={() => setShowAllTx(false)}
        footer={null}
        width={700}
      >
        <Table
          dataSource={transactions}
          rowKey="id"
          pagination={{ pageSize: 10, showSizeChanger: false, showTotal: (total) => `${total} transactions` }}
          size="small"
          columns={[
            { title: 'Type', dataIndex: 'txType', key: 'txType', width: 100, render: (type) => <Tag color={type === 'DEPOSIT' ? 'green' : 'blue'}>{type}</Tag> },
            { title: 'Amount', dataIndex: 'amount', key: 'amount', render: (amount, record: any) => <span style={{ color: record.txType === 'DEPOSIT' ? theme.colorSuccess : theme.colorText, fontWeight: 600 }}>{record.txType === 'DEPOSIT' ? '+' : '-'}{amount} USDT</span> },
            { title: 'Date', dataIndex: 'createdAt', key: 'createdAt', render: (date) => new Date(date).toLocaleString() },
            {
              title: 'TX Hash',
              dataIndex: 'txHash',
              key: 'txHash',
              render: (txHash) => txHash ? (
                <TransactionLink txHash={txHash} networkId={currentNetwork?.id} />
              ) : '-'
            },
          ]}
        />
      </Modal>
    </Flexbox>
  );
};

export default HomeContent;

