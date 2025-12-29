import { Card, Table, Button, Modal, Form, Input, Tag, Statistic, Row, Col, Tabs, Empty, message, Skeleton } from 'antd';
import { ShoppingCart, TrendingUp, TrendingDown, Plus, History, Eye, RefreshCw } from 'lucide-react';
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
    handleTradeCompleted,
  } = useMarketplaceStore();

  const [createOrderModal, setCreateOrderModal] = useState(false);
  const [form] = Form.useForm();

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
    const unsubscribeTradeCompleted = Events.On('marketplace:trade_completed', handleTradeCompletedEvent);
    const unsubscribeUserOrderCreated = Events.On(`marketplace:user:${walletAddress}:order_created`, handleOrderCreated);
    const unsubscribeUserOrderStatus = Events.On(`marketplace:user:${walletAddress}:order_status_update`, handleOrderStatusUpdate);
    const unsubscribeUserTradeCompleted = Events.On(`marketplace:user:${walletAddress}:trade_completed`, handleTradeCompletedEvent);

    return () => {
      unsubscribeMarketData();
      unsubscribeOrderCreated();
      unsubscribeOrderStatus();
      unsubscribeTradeCompleted();
      unsubscribeUserOrderCreated();
      unsubscribeUserOrderStatus();
      unsubscribeUserTradeCompleted();
    };
  }, [walletAddress, updateMarketStats, addOrder, updateOrderStatus, handleTradeCompleted, refreshData]);

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

  // Buy order
  const handleBuyOrderClick = async (orderID: string, partial?: boolean, amount?: string) => {
    const success = partial && amount 
      ? await buyPartialOrder(orderID, amount)
      : await buyOrder(orderID);
      
    if (success) {
      message.success('Trade executed successfully!');
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
      title: 'Amount (KAWAI)',
      dataIndex: 'remainingAmount',
      key: 'remainingAmount',
      render: (amount: string) => parseFloat(amount).toFixed(2),
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
        <Button
          type="primary"
          size="small"
          onClick={() => handleBuyOrderClick(record.id)}
          disabled={record.seller === walletAddress}
        >
          {record.seller === walletAddress ? 'Your Order' : 'Buy'}
        </Button>
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
      title: 'Amount (KAWAI)',
      dataIndex: 'tokenAmount',
      key: 'tokenAmount',
      render: (amount: string) => parseFloat(amount).toFixed(2),
    },
    {
      title: 'Remaining',
      dataIndex: 'remainingAmount',
      key: 'remainingAmount',
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
    </Flexbox>
  );
};

export default OTCContent;

