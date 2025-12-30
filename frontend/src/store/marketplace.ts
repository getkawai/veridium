import { create } from 'zustand';
import { subscribeWithSelector } from 'zustand/middleware';
import { MarketplaceService } from '@@/github.com/kawai-network/veridium/internal/services';
import type { Order, MarketStats, OrderHistoryEntry, TradeHistoryEntry } from '@@/github.com/kawai-network/veridium/internal/services/models';

interface MarketplaceState {
  // Data
  activeOrders: Order[];
  marketStats: MarketStats | null;
  userOrders: Order[];
  orderHistory: OrderHistoryEntry[];
  tradeHistory: TradeHistoryEntry[];
  
  // UI State - granular loading states
  loading: {
    activeOrders: boolean;
    marketStats: boolean;
    userOrders: boolean;
    orderHistory: boolean;
    tradeHistory: boolean;
  };
  refreshing: boolean;
  error: string | null;
  
  // Actions
  loadMarketplaceData: (address: string) => Promise<void>;
  refreshData: (address: string) => Promise<void>;
  createSellOrder: (tokenAmount: string, usdtPrice: string) => Promise<boolean>;
  buyOrder: (orderID: string) => Promise<boolean>;
  buyPartialOrder: (orderID: string, tokenAmount: string) => Promise<boolean>;
  cancelOrder: (orderID: string) => Promise<boolean>;
  
  // Real-time updates
  updateMarketStats: (stats: MarketStats) => void;
  addOrder: (order: Order) => void;
  updateOrderStatus: (orderID: string, status: string) => void;
  updateOrderPartialFill: (orderID: string, remainingAmount: string) => void;
  handleTradeCompleted: (trade: any, userAddress: string) => void;
  
  // Reset
  reset: () => void;
}

const initialState = {
  activeOrders: [],
  marketStats: null,
  userOrders: [],
  orderHistory: [],
  tradeHistory: [],
  loading: {
    activeOrders: false,
    marketStats: false,
    userOrders: false,
    orderHistory: false,
    tradeHistory: false,
  },
  refreshing: false,
  error: null,
};

export const useMarketplaceStore = create<MarketplaceState>()(
  subscribeWithSelector((set, get) => ({
    ...initialState,

    loadMarketplaceData: async (address: string) => {
      if (!address) return;
      
      set({ 
        loading: { 
          activeOrders: true, 
          marketStats: true, 
          userOrders: true, 
          orderHistory: true, 
          tradeHistory: true,
        }, 
        error: null 
      });
      
      try {
        // Load data sequentially to show progressive loading
        const activeOrdersPromise = MarketplaceService.GetActiveOrders('price_asc', {}).then(data => {
          set(state => ({ 
            activeOrders: data || [], 
            loading: { ...state.loading, activeOrders: false } 
          }));
          return data;
        });

        const marketStatsPromise = MarketplaceService.GetMarketStats().then(data => {
          set(state => ({ 
            marketStats: data || null, 
            loading: { ...state.loading, marketStats: false } 
          }));
          return data;
        });

        const userOrdersPromise = MarketplaceService.GetUserOrders(address).then(data => {
          set(state => ({ 
            userOrders: data || [], 
            loading: { ...state.loading, userOrders: false } 
          }));
          return data;
        });

        const orderHistoryPromise = MarketplaceService.GetOrderHistory(address).then(result => {
          const data = result?.orders || [];
          set(state => ({ 
            orderHistory: data, 
            loading: { ...state.loading, orderHistory: false } 
          }));
          return data;
        });

        const tradeHistoryPromise = MarketplaceService.GetTradeHistory(address).then(data => {
          set(state => ({ 
            tradeHistory: data || [], 
            loading: { ...state.loading, tradeHistory: false } 
          }));
          return data;
        });

        // Wait for all to complete
        await Promise.all([
          activeOrdersPromise,
          marketStatsPromise,
          userOrdersPromise,
          orderHistoryPromise,
          tradeHistoryPromise,
        ]);

        set({ error: null });
      } catch (error) {
        console.error('Failed to load marketplace data:', error);
        set({ 
          loading: {
            activeOrders: false,
            marketStats: false,
            userOrders: false,
            orderHistory: false,
            tradeHistory: false,
          }, 
          error: error instanceof Error ? error.message : 'Failed to load marketplace data' 
        });
      }
    },

    refreshData: async (address: string) => {
      if (!address) return;
      
      set({ refreshing: true });
      
      try {
        const [activeOrders, marketStats, userOrders] = await Promise.all([
          MarketplaceService.GetActiveOrders('price_asc', {}),
          MarketplaceService.GetMarketStats(),
          MarketplaceService.GetUserOrders(address),
        ]);

        set({
          activeOrders: activeOrders || [],
          marketStats: marketStats || null,
          userOrders: userOrders || [],
          refreshing: false,
        });
      } catch (error) {
        console.error('Failed to refresh marketplace data:', error);
        set({ refreshing: false });
      }
    },

    createSellOrder: async (tokenAmount: string, usdtPrice: string) => {
      try {
        const result = await MarketplaceService.CreateSellOrder(tokenAmount, usdtPrice);
        if (result?.success) {
          return true;
        } else {
          set({ error: result?.error || 'Failed to create order' });
          return false;
        }
      } catch (error) {
        console.error('Failed to create order:', error);
        set({ error: error instanceof Error ? error.message : 'Failed to create order' });
        return false;
      }
    },

    buyOrder: async (orderID: string) => {
      try {
        const result = await MarketplaceService.BuyOrder(orderID);
        if (result?.success) {
          return true;
        } else {
          set({ error: result?.error || 'Failed to execute trade' });
          return false;
        }
      } catch (error) {
        console.error('Failed to execute trade:', error);
        set({ error: error instanceof Error ? error.message : 'Failed to execute trade' });
        return false;
      }
    },

    buyPartialOrder: async (orderID: string, tokenAmount: string) => {
      try {
        const result = await MarketplaceService.BuyPartialOrder(orderID, tokenAmount);
        if (result?.success) {
          return true;
        } else {
          set({ error: result?.error || 'Failed to execute partial trade' });
          return false;
        }
      } catch (error) {
        console.error('Failed to execute partial trade:', error);
        set({ error: error instanceof Error ? error.message : 'Failed to execute partial trade' });
        return false;
      }
    },

    cancelOrder: async (orderID: string) => {
      try {
        await MarketplaceService.CancelOrder(orderID);
        return true;
      } catch (error) {
        console.error('Failed to cancel order:', error);
        set({ error: error instanceof Error ? error.message : 'Failed to cancel order' });
        return false;
      }
    },

    updateMarketStats: (stats: MarketStats) => {
      set({ marketStats: stats });
    },

    addOrder: (order: Order) => {
      set(state => ({
        activeOrders: [order, ...state.activeOrders],
      }));
    },

    updateOrderStatus: (orderID: string, status: string) => {
      set(state => ({
        activeOrders: state.activeOrders.map(order => 
          order.id === orderID ? { ...order, status } : order
        ),
        userOrders: state.userOrders.map(order => 
          order.id === orderID ? { ...order, status } : order
        ),
      }));
    },

    updateOrderPartialFill: (orderID: string, remainingAmount: string) => {
      set(state => ({
        activeOrders: state.activeOrders.map(order => 
          order.id === orderID ? { ...order, remainingAmount, status: 'active' } : order
        ),
        userOrders: state.userOrders.map(order => 
          order.id === orderID ? { ...order, remainingAmount, status: 'active' } : order
        ),
      }));
    },

    handleTradeCompleted: (trade: any, userAddress: string) => {
      // Remove filled orders from active orders
      set(state => ({
        activeOrders: state.activeOrders.filter(order => 
          order.id !== trade.orderID || order.status !== 'filled'
        ),
      }));
    },

    reset: () => {
      set(initialState);
    },
  }))
);