// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title OTCMarket
 * @notice P2P OTC Market for trading KAWAI tokens vs USDC with partial fill support
 * @dev Secure P2P trading system with ReentrancyGuard protection
 * 
 * Website: https://getkawai.com
 * Docs: https://getkawai.com/docs
 * 
 * Features:
 * - Create sell orders (tokens locked in contract)
 * - Buy full or partial amounts
 * - Cancel orders (refund remaining tokens)
 * - Zero fees (configurable for future)
 * - Efficient order querying (pagination support)
 */
contract OTCMarket is ReentrancyGuard {
    using SafeERC20 for IERC20;

    struct Order {
        uint256 id;
        address seller;
        uint256 tokenAmount;      // Original token amount
        uint256 priceInUSDT;      // Original total price in USDT
        uint256 remainingAmount;  // ✅ NEW: Remaining tokens available
        bool isActive;            // true if remainingAmount > 0
    }

    IERC20 public immutable tokenDeAI;
    IERC20 public immutable usdt;

    // Fee in basis points (1 = 0.01%). Example: 100 = 1%
    uint256 public constant FEE_BPS = 0;
    address public feeRecipient;

    Order[] public orders;

    event OrderCreated(
        uint256 indexed orderId,
        address indexed seller,
        uint256 amount,
        uint256 price
    );
    
    event OrderCancelled(
        uint256 indexed orderId, 
        address indexed seller
    );
    
    event OrderFulfilled(
        uint256 indexed orderId,
        address indexed buyer,
        address indexed seller,
        uint256 amount,
        uint256 price
    );
    
    // ✅ NEW: Event for partial fills
    event OrderPartiallyFilled(
        uint256 indexed orderId,
        address indexed buyer,
        address indexed seller,
        uint256 amountFilled,
        uint256 remainingAmount,
        uint256 pricePaid
    );

    constructor(address _tokenDeAI, address _usdt, address _feeRecipient) {
        require(_tokenDeAI != address(0), "Invalid token address");
        require(_usdt != address(0), "Invalid USDT address");
        tokenDeAI = IERC20(_tokenDeAI);
        usdt = IERC20(_usdt);
        feeRecipient = _feeRecipient;
    }

    /**
     * @notice Create a sell order.
     * @param _amount Amount of DeAI tokens to sell.
     * @param _priceInUSDT Total price in USDT (not per token).
     */
    function createOrder(
        uint256 _amount,
        uint256 _priceInUSDT
    ) external nonReentrant {
        require(_amount > 0, "Amount must be > 0");
        require(_priceInUSDT > 0, "Price must be > 0");

        // Lock tokens in contract
        tokenDeAI.safeTransferFrom(msg.sender, address(this), _amount);

        uint256 orderId = orders.length;
        orders.push(
            Order({
                id: orderId,
                seller: msg.sender,
                tokenAmount: _amount,
                priceInUSDT: _priceInUSDT,
                remainingAmount: _amount,  // ✅ NEW: Initialize remaining amount
                isActive: true
            })
        );

        emit OrderCreated(orderId, msg.sender, _amount, _priceInUSDT);
    }

    /**
     * @notice Buy a specific order (full amount).
     * @param _orderId ID of the order to buy.
     */
    function buyOrder(uint256 _orderId) external nonReentrant {
        require(_orderId < orders.length, "Invalid Order ID");
        Order storage order = orders[_orderId];
        require(order.isActive, "Order not active");

        // Buy the full remaining amount
        _executeTrade(_orderId, order.remainingAmount);
    }

    /**
     * @notice Buy a partial amount from an order.
     * @param _orderId ID of the order to buy from.
     * @param _amount Amount of tokens to buy (must be <= remainingAmount).
     */
    function buyOrderPartial(
        uint256 _orderId, 
        uint256 _amount
    ) external nonReentrant {
        require(_orderId < orders.length, "Invalid Order ID");
        Order storage order = orders[_orderId];
        require(order.isActive, "Order not active");
        require(_amount > 0, "Amount must be > 0");
        require(_amount <= order.remainingAmount, "Amount exceeds remaining");

        _executeTrade(_orderId, _amount);
    }

    /**
     * @notice Internal function to execute a trade (full or partial).
     * @param _orderId ID of the order.
     * @param _amount Amount of tokens to trade.
     */
    function _executeTrade(uint256 _orderId, uint256 _amount) internal {
        Order storage order = orders[_orderId];

        // Calculate proportional USDT price
        uint256 proportionalPrice = (order.priceInUSDT * _amount) / order.tokenAmount;
        
        // Prevent zero-price trades due to integer truncation
        require(proportionalPrice > 0, "Amount too small, price rounds to zero");

        // Update remaining amount
        order.remainingAmount -= _amount;

        // Mark as inactive if fully filled
        if (order.remainingAmount == 0) {
            order.isActive = false;
        }

        // Calculate fee (if any)
        uint256 feeAmount = (proportionalPrice * FEE_BPS) / 10000;
        uint256 sellerAmount = proportionalPrice - feeAmount;

        // Transfer USDT from Buyer -> Seller (+ Fee)
        if (feeAmount > 0 && feeRecipient != address(0)) {
            usdt.safeTransferFrom(msg.sender, feeRecipient, feeAmount);
        }
        usdt.safeTransferFrom(msg.sender, order.seller, sellerAmount);

        // Transfer DeAI Token from contract -> Buyer
        tokenDeAI.safeTransfer(msg.sender, _amount);

        // Emit appropriate event
        if (order.remainingAmount == 0) {
            // Fully filled
            emit OrderFulfilled(
                _orderId,
                msg.sender,
                order.seller,
                _amount,
                proportionalPrice
            );
        } else {
            // Partially filled
            emit OrderPartiallyFilled(
                _orderId,
                msg.sender,
                order.seller,
                _amount,
                order.remainingAmount,
                proportionalPrice
            );
        }
    }

    /**
     * @notice Cancel your own order.
     * @param _orderId ID of the order to cancel.
     */
    function cancelOrder(uint256 _orderId) external nonReentrant {
        require(_orderId < orders.length, "Invalid Order ID");
        Order storage order = orders[_orderId];
        require(order.seller == msg.sender, "Not your order");
        require(order.isActive, "Order already sold/cancelled");

        uint256 refundAmount = order.remainingAmount;
        order.isActive = false;
        order.remainingAmount = 0;

        // Return remaining tokens to seller
        tokenDeAI.safeTransfer(msg.sender, refundAmount);

        emit OrderCancelled(_orderId, msg.sender);
    }

    /**
     * @notice Get the total number of orders.
     */
    function getOrdersCount() external view returns (uint256) {
        return orders.length;
    }

    /**
     * @notice Get order details by ID.
     * @param _orderId ID of the order.
     */
    function getOrder(uint256 _orderId) external view returns (
        uint256 id,
        address seller,
        uint256 tokenAmount,
        uint256 priceInUSDT,
        uint256 remainingAmount,
        bool isActive
    ) {
        require(_orderId < orders.length, "Invalid Order ID");
        Order storage order = orders[_orderId];
        return (
            order.id,
            order.seller,
            order.tokenAmount,
            order.priceInUSDT,
            order.remainingAmount,
            order.isActive
        );
    }

    /**
     * @notice Get multiple orders at once (for efficient querying).
     * @param _orderIds Array of order IDs to fetch.
     */
    function getOrders(uint256[] calldata _orderIds) external view returns (Order[] memory) {
        Order[] memory result = new Order[](_orderIds.length);
        for (uint256 i = 0; i < _orderIds.length; i++) {
            require(_orderIds[i] < orders.length, "Invalid Order ID");
            result[i] = orders[_orderIds[i]];
        }
        return result;
    }

    /**
     * @notice Get all active orders (paginated to avoid gas limit).
     * @param _offset Starting index.
     * @param _limit Maximum number of orders to return.
     */
    function getActiveOrders(
        uint256 _offset, 
        uint256 _limit
    ) external view returns (Order[] memory) {
        require(_limit <= 100, "Limit too high");
        
        // Count active orders
        uint256 activeCount = 0;
        for (uint256 i = 0; i < orders.length; i++) {
            if (orders[i].isActive) {
                activeCount++;
            }
        }
        
        // Calculate result size
        uint256 resultSize = activeCount > _offset ? activeCount - _offset : 0;
        if (resultSize > _limit) {
            resultSize = _limit;
        }
        
        // Fill result array
        Order[] memory result = new Order[](resultSize);
        uint256 resultIndex = 0;
        uint256 currentOffset = 0;
        
        for (uint256 i = 0; i < orders.length && resultIndex < resultSize; i++) {
            if (orders[i].isActive) {
                if (currentOffset >= _offset) {
                    result[resultIndex] = orders[i];
                    resultIndex++;
                }
                currentOffset++;
            }
        }
        
        return result;
    }

    /**
     * @notice Get orders by seller (for user's order history).
     * @param _seller Address of the seller.
     * @param _offset Starting index.
     * @param _limit Maximum number of orders to return.
     */
    function getOrdersBySeller(
        address _seller,
        uint256 _offset,
        uint256 _limit
    ) external view returns (Order[] memory) {
        require(_limit <= 100, "Limit too high");
        
        // Count seller's orders
        uint256 sellerOrderCount = 0;
        for (uint256 i = 0; i < orders.length; i++) {
            if (orders[i].seller == _seller) {
                sellerOrderCount++;
            }
        }
        
        // Calculate result size
        uint256 resultSize = sellerOrderCount > _offset ? sellerOrderCount - _offset : 0;
        if (resultSize > _limit) {
            resultSize = _limit;
        }
        
        // Fill result array
        Order[] memory result = new Order[](resultSize);
        uint256 resultIndex = 0;
        uint256 currentOffset = 0;
        
        for (uint256 i = 0; i < orders.length && resultIndex < resultSize; i++) {
            if (orders[i].seller == _seller) {
                if (currentOffset >= _offset) {
                    result[resultIndex] = orders[i];
                    resultIndex++;
                }
                currentOffset++;
            }
        }
        
        return result;
    }
}
