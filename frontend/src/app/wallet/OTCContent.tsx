import { Card, Table, Button, Modal, Form, Input, Tag, Statistic, Row, Col, Tabs, Empty, message, Skeleton, Progress, InputNumber, Space } from 'antd';
import { ShoppingCart, TrendingUp, TrendingDown, Plus, History, Eye, RefreshCw, Percent } from 'lucide-react';
import { Flexbox } from 'react-layout-kit';
import { useState, useEffect } from 'react';
import { Events } from '@wailsio/runtime';
import { useUserStore } from '../../store/user';
import { useMarketplaceStore } from '../../store/marketplace';
import type { Order, TradeHistoryEntry } from '@@/github.com/kawai-network/veridium/internal/services/models';
import type { OTCContentProps } from './types';

const OTCContent = ({ styles, theme }: OTCContentProps) => {
  const { walletAddress } = useUserStore();
  const {
    activeOrders,
    marketStats,
    userOrders,
    orderHistory,
    tradeHistory,
    loading,
    refreshing,
    error,
    loadMarketplaceData,
    refreshData,
    createSellOrder,
    buyOrder,
    buyPartialOrder,
    cancelOrder,
    updateMarketStats,
    addOrder,
    updateOrderStatus,
    updateOrderPartialFill,
    handleTradeCompleted,
  } = useMarketplaceStore();

  const [createOrderModal, setCreateOrderModal] = useState(false);
  const [partialBuyModal, setPartialBuyModal] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
  const [form] = Form.useForm();
  const [partialBuyForm] = Form.useForm();

  // Real-time event handlers
  useEffect(() => {
    if (!walletAddress) return;

    const handleMarketDataUpdate = (ev: any) => {
      const data = ev.data;
      if (data) {
        updateMarketStats(data);
      }
    };

    const handleOrderCreated = (ev: any) => {
      const data = ev.data;
      if (data) {
        addOrder(data);
      }
    };

    const handleOrderStatusUpdate = (ev: any) => {
      const data = ev.data;
      if (data) {
        updateOrderStatus(data.orderID, data.status);
      }
    };

    const handleOrderPartiallyFilled = (ev: any) => {
      const data = ev.data;
      if (data) {
        // Update order with new remaining amount
        updateOrderPartialFill(data.orderID, data.remainingAmount);
        // Show notification
        if (data.seller === walletAddress) {
          message.info(`Your order was partially filled! ${data.amountFilled} KAWAI sold.`);
        }
      }
    };

    const handleTradeCompletedEvent = (ev: any) => {
      const data = ev.data;
      if (data && (data.buyer === walletAddress || data.seller === walletAddress)) {
        handleTradeCompleted(data, walletAddress);
        // Refresh data to get updated order book and user data
        refreshData(walletAddress);
      }
    };

    // Subscribe to marketplace events
    const unsubscribeMarketData = Events.On('marketplace:market_data_update', handleMarketDataUpdate);
    const unsubscribeOrderCreated = Events.On('marketplace:order_created', handleOrderCreated);
    const unsubscribeOrderStatus = Events.On('marketplace:order_status_update', handleOrderStatusUpdate);
    const unsubscribeOrderPartiallyFilled = Events.On('marketplace:order_partially_filled', handleOrderPartiallyFilled);
    const unsubscribeTradeCompleted = Events.On('marketplace:trade_completed', handleTradeCompletedEvent);
    const unsubscribeUserOrderCreated = Events.On(`marketplace:user:${walletAddress}:order_created`, handleOrderCreated);
    const unsubscribeUserOrderStatus = Events.On(`marketplace:user:${walletAddress}:order_status_update`, handleOrderStatusUpdate);
    const unsubscribeUserOrderPartiallyFilled = Events.On(`marketplace:user:${walletAddress}:order_partially_filled`, handleOrderPartiallyFilled);
    const unsubscribeUserTradeCompleted = Events.On(`marketplace:user:${walletAddress}:trade_completed`, handleTradeCompletedEvent);

    return () => {
      unsubscribeMarketData();
      unsubscribeOrderCreated();
      unsubscribeOrderStatus();
      unsubscribeOrderPartiallyFilled();
      unsubscribeTradeCompleted();
      unsubscribeUserOrderCreated();
      unsubscribeUserOrderStatus();
      unsubscribeUserOrderPartiallyFilled();
      unsubscribeUserTradeCompleted();
    };
  }, [walletAddress, updateMarketStats, addOrder, updateOrderStatus, updateOrderPartialFill, handleTradeCompleted, refreshData]);

  // Initial data load
  useEffect(() => {
    if (walletAddress) {
      loadMarketplaceData(walletAddress);
    }
  }, [walletAddress, loadMarketplaceData]);

  // Show error messages
  useEffect(() => {
    if (error) {
      message.error(error);
    }
  }, [error]);

  // Create sell order
  const handleCreateOrder = async (values: { tokenAmount: string; usdtPrice: string }) => {
    const success = await createSellOrder(values.tokenAmount, values.usdtPrice);
    if (success) {
      message.success('Order created successfully!');
      setCreateOrderModal(false);
      form.resetFields();
      if (walletAddress) {
        refreshData(walletAddress);
      }
    }
  };

  // Buy order (full)
  const handleBuyOrderClick = async (orderID: string) => {
    const success = await buyOrder(orderID);
      
    if (success) {
      message.success('Trade executed successfully!');
      if (walletAddress) {
        refreshData(walletAddress);
      }
    }
  };

  // Open partial buy modal
  const handlePartialBuyClick = (order: Order) => {
    setSelectedOrder(order);
    setPartialBuyModal(true);
    partialBuyForm.setFieldsValue({
      amount: parseFloat(order.remainingAmount) / 2, // Default to 50%
    });
  };

  // Execute partial buy
  const handlePartialBuySubmit = async (values: { amount: number }) => {
    if (!selectedOrder) return;
    
    const success = await buyPartialOrder(selectedOrder.id, values.amount.toString());
    
    if (success) {
      message.success(`Partial buy executed! Bought ${values.amount} KAWAI`);
      setPartialBuyModal(false);
      setSelectedOrder(null);
      partialBuyForm.resetFields();
      if (walletAddress) {
        refreshData(walletAddress);
      }
    }
  };

  // Cancel order
  const handleCancelOrderClick = async (orderID: string) => {
    const success = await cancelOrder(orderID);
    if (success) {
      message.success('Order cancelled successfully!');
      if (walletAddress) {
        refreshData(walletAddress);
      }
    }
  };

  const orderBookColumns = [
    {
      title: 'Price (USDT)',
      dataIndex: 'pricePerToken',
      key: 'pricePerToken',
      render: (price: string) => (
        <span style={{ fontWeight: 600, color: theme.colorSuccess }}>
          ${parseFloat(price).toFixed(4)}
        </span>
      ),
    },
    {
      title: 'Available / Total',
      key: 'amount',
      render: (record: Order) => {
        const remaining = parseFloat(record.remainingAmount);
        const total = parseFloat(record.tokenAmount);
        const filledPercent = ((total - remaining) / total) * 100;
        
        return (
          <Flexbox gap={4}>
            <div style={{ fontSize: 13 }}>
              <span style={{ fontWeight: 600 }}>{remaining.toFixed(2)}</span>
              <span style={{ color: theme.colorTextSecondary }}> / {total.toFixed(2)} KAWAI</span>
            </div>
            {filledPercent > 0 && (
              <Progress 
                percent={filledPercent} 
                size="small" 
                showInfo={false}
                strokeColor={theme.colorWarning}
              />
            )}
          </Flexbox>
        );
      },
    },
    {
      title: 'Total (USDT)',
      key: 'total',
      render: (record: Order) => {
        const total = parseFloat(record.remainingAmount) * parseFloat(record.pricePerToken);
        return `$${total.toFixed(2)}`;
      },
    },
    {
      title: 'Action',
      key: 'action',
      render: (record: Order) => (
        <Space size="small">
          <Button
            type="primary"
            size="small"
            onClick={() => handleBuyOrderClick(record.id)}
            disabled={record.seller === walletAddress}
          >
            {record.seller === walletAddress ? 'Your Order' : 'Buy All'}
          </Button>
          {record.seller !== walletAddress && (
            <Button
              size="small"
              onClick={() => handlePartialBuyClick(record)}
              icon={<Percent size={14} />}
            >
              Partial
            </Button>
          )}
        </Space>
      ),
    },
  ];

  const userOrderColumns = [
    {
      title: 'Price (USDT)',
      dataIndex: 'pricePerToken',
      key: 'pricePerToken',
      render: (price: string) => `$${parseFloat(price).toFixed(4)}`,
    },
    {
      title: 'Order Progress',
      key: 'progress',
      render: (record: Order) => {
        const remaining = parseFloat(record.remainingAmount);
        const total = parseFloat(record.tokenAmount);
        const filled = total - remaining;
        const filledPercent = (filled / total) * 100;
        
        return (
          <Flexbox gap={4}>
            <div style={{ fontSize: 12 }}>
              <span style={{ color: theme.colorTextSecondary }}>Filled: </span>
              <span style={{ fontWeight: 600 }}>{filled.toFixed(2)}</span>
              <span style={{ color: theme.colorTextSecondary }}> / {total.toFixed(2)} KAWAI</span>
            </div>
            <Progress 
              percent={filledPercent} 
              size="small"
              status={record.status === 'filled' ? 'success' : 'active'}
              strokeColor={record.status === 'filled' ? theme.colorSuccess : theme.colorPrimary}
            />
            <div style={{ fontSize: 11, color: theme.colorTextTertiary }}>
              Remaining: {remaining.toFixed(2)} KAWAI
            </div>
          </Flexbox>
        );
      },
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' ? 'green' : status === 'filled' ? 'blue' : 'red'}>
          {status.toUpperCase()}
        </Tag>
      ),
    },
    {
      title: 'Action',
      key: 'action',
      render: (record: Order) => (
        <Flexbox horizontal gap={8}>
          <Button
            size="small"
            icon={<Eye size={14} />}
          >
            Details
          </Button>
          {record.status === 'active' && (
            <Button
              size="small"
              danger
              onClick={() => handleCancelOrderClick(record.id)}
            >
              Cancel
            </Button>
          )}
        </Flexbox>
      ),
    },
  ];

  return (
    <Flexbox style={{ maxWidth: 1200, width: '100%' }} gap={20}>
      {/* Header */}
      <Flexbox horizontal justify="space-between" align="center">
        <div>
          <h2 style={{ margin: 0, fontSize: 20, fontWeight: 600 }}>OTC Market</h2>
          <span style={{ color: theme.colorTextSecondary, fontSize: 13 }}>P2P trading for KAWAI tokens</span>
        </div>
        <Flexbox horizontal gap={12}>
          <Button
            icon={<RefreshCw size={16} />}
            onClick={() => walletAddress && refreshData(walletAddress)}
            loading={refreshing}
          >
            Refresh
          </Button>
          <Button
            type="primary"
            icon={<Plus size={16} />}
            onClick={() => setCreateOrderModal(true)}
          >
            Create Sell Order
          </Button>
        </Flexbox>
      </Flexbox>

      {/* Market Stats */}
      <Card size="small">
        {loading.marketStats ? (
          <Row gutter={24}>
            <Col span={6}>
              <Skeleton active paragraph={{ rows: 2, width: ['60%', '80%'] }} />
            </Col>
            <Col span={6}>
              <Skeleton active paragraph={{ rows: 2, width: ['60%', '80%'] }} />
            </Col>
            <Col span={6}>
              <Skeleton active paragraph={{ rows: 2, width: ['60%', '80%'] }} />
            </Col>
            <Col span={6}>
              <Skeleton active paragraph={{ rows: 2, width: ['60%', '80%'] }} />
            </Col>
          </Row>
        ) : marketStats ? (
          <Row gutter={24}>
            <Col span={6}>
              <Statistic
                title="Lowest Ask"
                value={parseFloat(marketStats.lowestAskPrice || '0')}
                precision={4}
                prefix="$"
                valueStyle={{ color: theme.colorSuccess }}
              />
            </Col>
            <Col span={6}>
              <Statistic
                title="Highest Bid"
                value={parseFloat(marketStats.highestBidPrice || '0')}
                precision={4}
                prefix="$"
                valueStyle={{ color: theme.colorError }}
              />
            </Col>
            <Col span={6}>
              <Statistic
                title="24h Volume"
                value={parseFloat(marketStats.volume24h || '0')}
                precision={2}
                prefix="$"
              />
            </Col>
            <Col span={6}>
              <Statistic
                title="24h Change"
                value={parseFloat(marketStats.priceChange24h || '0')}
                precision={2}
                suffix="%"
                valueStyle={{ 
                  color: parseFloat(marketStats.priceChange24h || '0') >= 0 
                    ? theme.colorSuccess 
                    : theme.colorError 
                }}
                prefix={parseFloat(marketStats.priceChange24h || '0') >= 0 ? <TrendingUp size={16} /> : <TrendingDown size={16} />}
              />
            </Col>
          </Row>
        ) : (
          <Empty description="No market data available" />
        )}
      </Card>

      {/* Main Content */}
      <Row gutter={24}>
        {/* Order Book */}
        <Col span={14}>
          <Card
            title={
              <Flexbox horizontal align="center" gap={8}>
                <ShoppingCart size={16} />
                Order Book ({activeOrders.length} orders)
              </Flexbox>
            }
            size="small"
          >
            {loading.activeOrders ? (
              <Skeleton active paragraph={{ rows: 8 }} />
            ) : activeOrders.length > 0 ? (
              <Table
                dataSource={activeOrders}
                columns={orderBookColumns}
                rowKey="id"
                pagination={{ pageSize: 10, showSizeChanger: false }}
                size="small"
              />
            ) : (
              <Empty description="No active orders" />
            )}
          </Card>
        </Col>

        {/* Recent Trades */}
        <Col span={10}>
          <Card
            title={
              <Flexbox horizontal align="center" gap={8}>
                <History size={16} />
                Recent Trades
              </Flexbox>
            }
            size="small"
          >
            {loading.marketStats ? (
              <Skeleton active paragraph={{ rows: 6 }} />
            ) : marketStats?.recentTrades && marketStats.recentTrades.length > 0 ? (
              <div style={{ maxHeight: 300, overflowY: 'auto' }}>
                {marketStats.recentTrades.map((trade, index) => (
                  <div
                    key={trade.id || index}
                    style={{
                      padding: '8px 0',
                      borderBottom: index < marketStats!.recentTrades.length - 1 ? `1px solid ${theme.colorBorderSecondary}` : 'none',
                    }}
                  >
                    <Flexbox horizontal justify="space-between">
                      <span style={{ fontSize: 12, color: theme.colorTextSecondary }}>
                        {parseFloat(trade.tokenAmount).toFixed(2)} KAWAI
                      </span>
                      <span style={{ fontSize: 12, fontWeight: 600, color: theme.colorSuccess }}>
                        ${parseFloat(trade.price).toFixed(4)}
                      </span>
                    </Flexbox>
                  </div>
                ))}
              </div>
            ) : (
              <Empty description="No recent trades" />
            )}
          </Card>
        </Col>
      </Row>

      {/* User Orders & History */}
      <Card size="small">
        <Tabs
          items={[
            {
              key: 'orders',
              label: `My Orders (${userOrders.length})`,
              children: loading.userOrders ? (
                <Skeleton active paragraph={{ rows: 5 }} />
              ) : userOrders.length > 0 ? (
                <Table
                  dataSource={userOrders}
                  columns={userOrderColumns}
                  rowKey="id"
                  pagination={{ pageSize: 5, showSizeChanger: false }}
                  size="small"
                />
              ) : (
                <Empty description="No orders yet" />
              ),
            },
            {
              key: 'history',
              label: `Order History (${orderHistory.length})`,
              children: loading.orderHistory ? (
                <Skeleton active paragraph={{ rows: 5 }} />
              ) : orderHistory.length > 0 ? (
                <Table
                  dataSource={orderHistory}
                  columns={[
                    {
                      title: 'Price (USDT)',
                      dataIndex: 'pricePerToken',
                      key: 'pricePerToken',
                      render: (price: string) => `$${parseFloat(price).toFixed(4)}`,
                    },
                    {
                      title: 'Amount (KAWAI)',
                      dataIndex: 'tokenAmount',
                      key: 'tokenAmount',
                      render: (amount: string) => parseFloat(amount).toFixed(2),
                    },
                    {
                      title: 'Filled',
                      dataIndex: 'filledAmount',
                      key: 'filledAmount',
                      render: (amount: string) => parseFloat(amount).toFixed(2),
                    },
                    {
                      title: 'Status',
                      dataIndex: 'status',
                      key: 'status',
                      render: (status: string) => (
                        <Tag color={status === 'active' ? 'green' : status === 'filled' ? 'blue' : 'red'}>
                          {status.toUpperCase()}
                        </Tag>
                      ),
                    },
                    {
                      title: 'Trades',
                      dataIndex: 'tradeCount',
                      key: 'tradeCount',
                    },
                    {
                      title: 'Date',
                      dataIndex: 'createdAt',
                      key: 'createdAt',
                      render: (date: any) => new Date(date).toLocaleDateString(),
                    },
                  ]}
                  rowKey="id"
                  pagination={{ pageSize: 5, showSizeChanger: false }}
                  size="small"
                />
              ) : (
                <Empty description="No order history" />
              ),
            },
            {
              key: 'trades',
              label: `Trade History (${tradeHistory.length})`,
              children: loading.tradeHistory ? (
                <Skeleton active paragraph={{ rows: 5 }} />
              ) : tradeHistory.length > 0 ? (
                <Table
                  dataSource={tradeHistory}
                  columns={[
                    {
                      title: 'Type',
                      key: 'type',
                      render: (record: TradeHistoryEntry) => (
                        <Tag color={record.buyer === walletAddress ? 'blue' : 'green'}>
                          {record.buyer === walletAddress ? 'BUY' : 'SELL'}
                        </Tag>
                      ),
                    },
                    {
                      title: 'Price (USDT)',
                      dataIndex: 'price',
                      key: 'price',
                      render: (price: string) => `$${parseFloat(price).toFixed(4)}`,
                    },
                    {
                      title: 'Amount (KAWAI)',
                      dataIndex: 'tokenAmount',
                      key: 'tokenAmount',
                      render: (amount: string) => parseFloat(amount).toFixed(2),
                    },
                    {
                      title: 'Total (USDT)',
                      dataIndex: 'usdtAmount',
                      key: 'usdtAmount',
                      render: (amount: string) => `$${parseFloat(amount).toFixed(2)}`,
                    },
                    {
                      title: 'Date',
                      dataIndex: 'timestamp',
                      key: 'timestamp',
                      render: (date: any) => new Date(date).toLocaleDateString(),
                    },
                  ]}
                  rowKey="id"
                  pagination={{ pageSize: 5, showSizeChanger: false }}
                  size="small"
                />
              ) : (
                <Empty description="No trade history" />
              ),
            },
          ]}
        />
      </Card>

      {/* Create Order Modal */}
      <Modal
        title="Create Sell Order"
        open={createOrderModal}
        onCancel={() => setCreateOrderModal(false)}
        footer={null}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleCreateOrder}
        >
          <Form.Item
            label="KAWAI Token Amount"
            name="tokenAmount"
            rules={[
              { required: true, message: 'Please enter token amount' },
              { pattern: /^\d+(\.\d+)?$/, message: 'Please enter a valid number' },
            ]}
          >
            <Input placeholder="Enter amount of KAWAI tokens to sell" />
          </Form.Item>
          
          <Form.Item
            label="USDT Price per Token"
            name="usdtPrice"
            rules={[
              { required: true, message: 'Please enter USDT price' },
              { pattern: /^\d+(\.\d+)?$/, message: 'Please enter a valid price' },
            ]}
          >
            <Input placeholder="Enter price per KAWAI token in USDT" />
          </Form.Item>

          <Flexbox horizontal justify="flex-end" gap={12}>
            <Button onClick={() => setCreateOrderModal(false)}>
              Cancel
            </Button>
            <Button type="primary" htmlType="submit">
              Create Order
            </Button>
          </Flexbox>
        </Form>
      </Modal>

      {/* Partial Buy Modal */}
      <Modal
        title="Partial Buy Order"
        open={partialBuyModal}
        onCancel={() => {
          setPartialBuyModal(false);
          setSelectedOrder(null);
          partialBuyForm.resetFields();
        }}
        footer={null}
        width={500}
      >
        {selectedOrder && (
          <Flexbox gap={16}>
            {/* Order Info */}
            <Card size="small" style={{ backgroundColor: theme.colorBgLayout }}>
              <Flexbox gap={8}>
                <Flexbox horizontal justify="space-between">
                  <span style={{ color: theme.colorTextSecondary }}>Price per Token:</span>
                  <span style={{ fontWeight: 600, color: theme.colorSuccess }}>
                    ${parseFloat(selectedOrder.pricePerToken).toFixed(4)}
                  </span>
                </Flexbox>
                <Flexbox horizontal justify="space-between">
                  <span style={{ color: theme.colorTextSecondary }}>Available:</span>
                  <span style={{ fontWeight: 600 }}>
                    {parseFloat(selectedOrder.remainingAmount).toFixed(2)} KAWAI
                  </span>
                </Flexbox>
                <Flexbox horizontal justify="space-between">
                  <span style={{ color: theme.colorTextSecondary }}>Total Available:</span>
                  <span style={{ fontWeight: 600 }}>
                    ${(parseFloat(selectedOrder.remainingAmount) * parseFloat(selectedOrder.pricePerToken)).toFixed(2)}
                  </span>
                </Flexbox>
              </Flexbox>
            </Card>

            {/* Partial Buy Form */}
            <Form
              form={partialBuyForm}
              layout="vertical"
              onFinish={handlePartialBuySubmit}
            >
              <Form.Item
                label="Amount to Buy (KAWAI)"
                name="amount"
                rules={[
                  { required: true, message: 'Please enter amount' },
                  {
                    validator: (_, value) => {
                      const remaining = parseFloat(selectedOrder.remainingAmount);
                      if (value <= 0) {
                        return Promise.reject('Amount must be greater than 0');
                      }
                      if (value > remaining) {
                        return Promise.reject(`Amount cannot exceed ${remaining.toFixed(2)} KAWAI`);
                      }
                      return Promise.resolve();
                    },
                  },
                ]}
              >
                <InputNumber
                  style={{ width: '100%' }}
                  placeholder="Enter amount"
                  min={0}
                  max={parseFloat(selectedOrder.remainingAmount)}
                  step={0.01}
                  precision={2}
                  onChange={(value) => {
                    if (value) {
                      const total = value * parseFloat(selectedOrder.pricePerToken);
                      const percent = (value / parseFloat(selectedOrder.remainingAmount)) * 100;
                      partialBuyForm.setFieldsValue({ 
                        calculatedTotal: total.toFixed(2),
                        calculatedPercent: percent.toFixed(1)
                      });
                    }
                  }}
                />
              </Form.Item>

              {/* Quick Select Buttons */}
              <Flexbox horizontal gap={8} style={{ marginBottom: 16 }}>
                <Button
                  size="small"
                  onClick={() => {
                    const amount = parseFloat(selectedOrder.remainingAmount) * 0.25;
                    partialBuyForm.setFieldsValue({ amount });
                    partialBuyForm.validateFields(['amount']);
                  }}
                >
                  25%
                </Button>
                <Button
                  size="small"
                  onClick={() => {
                    const amount = parseFloat(selectedOrder.remainingAmount) * 0.5;
                    partialBuyForm.setFieldsValue({ amount });
                    partialBuyForm.validateFields(['amount']);
                  }}
                >
                  50%
                </Button>
                <Button
                  size="small"
                  onClick={() => {
                    const amount = parseFloat(selectedOrder.remainingAmount) * 0.75;
                    partialBuyForm.setFieldsValue({ amount });
                    partialBuyForm.validateFields(['amount']);
                  }}
                >
                  75%
                </Button>
                <Button
                  size="small"
                  onClick={() => {
                    const amount = parseFloat(selectedOrder.remainingAmount);
                    partialBuyForm.setFieldsValue({ amount });
                    partialBuyForm.validateFields(['amount']);
                  }}
                >
                  100%
                </Button>
              </Flexbox>

              {/* Calculated Total */}
              <Form.Item dependencies={['amount']}>
                {() => {
                  const amount = partialBuyForm.getFieldValue('amount');
                  if (amount && amount > 0) {
                    const total = amount * parseFloat(selectedOrder.pricePerToken);
                    const percent = (amount / parseFloat(selectedOrder.remainingAmount)) * 100;
                    return (
                      <Card size="small" style={{ backgroundColor: theme.colorInfoBg }}>
                        <Flexbox gap={8}>
                          <Flexbox horizontal justify="space-between">
                            <span style={{ color: theme.colorTextSecondary }}>You will pay:</span>
                            <span style={{ fontWeight: 600, fontSize: 16, color: theme.colorPrimary }}>
                              ${total.toFixed(2)} USDT
                            </span>
                          </Flexbox>
                          <Flexbox horizontal justify="space-between">
                            <span style={{ color: theme.colorTextSecondary }}>You will receive:</span>
                            <span style={{ fontWeight: 600 }}>
                              {amount.toFixed(2)} KAWAI
                            </span>
                          </Flexbox>
                          <Progress 
                            percent={percent} 
                            size="small"
                            format={(percent) => `${percent?.toFixed(1)}% of order`}
                          />
                        </Flexbox>
                      </Card>
                    );
                  }
                  return null;
                }}
              </Form.Item>

              <Flexbox horizontal justify="flex-end" gap={12}>
                <Button onClick={() => {
                  setPartialBuyModal(false);
                  setSelectedOrder(null);
                  partialBuyForm.resetFields();
                }}>
                  Cancel
                </Button>
                <Button type="primary" htmlType="submit">
                  Buy Partial
                </Button>
              </Flexbox>
            </Form>
          </Flexbox>
        )}
      </Modal>
    </Flexbox>
  );
};

export default OTCContent;

