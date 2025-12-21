// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

/**
 * @title OTCMarket
 * @dev Simple OTC Market for P2P trading of DeAI Tokens vs USDT.
 */
contract OTCMarket is ReentrancyGuard {
    using SafeERC20 for IERC20;

    struct Order {
        uint256 id;
        address seller;
        uint256 tokenAmount;
        uint256 priceInUSDT;
        bool isActive;
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
    event OrderCancelled(uint256 indexed orderId, address indexed seller);
    event OrderFulfilled(
        uint256 indexed orderId,
        address indexed buyer,
        address indexed seller,
        uint256 amount,
        uint256 price
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

        // Lock tokens in escrow
        tokenDeAI.safeTransferFrom(msg.sender, address(this), _amount);

        uint256 orderId = orders.length;
        orders.push(
            Order({
                id: orderId,
                seller: msg.sender,
                tokenAmount: _amount,
                priceInUSDT: _priceInUSDT,
                isActive: true
            })
        );

        emit OrderCreated(orderId, msg.sender, _amount, _priceInUSDT);
    }

    /**
     * @notice Buy a specific order.
     * @param _orderId ID of the order to buy.
     */
    function buyOrder(uint256 _orderId) external nonReentrant {
        require(_orderId < orders.length, "Invalid Order ID");
        Order storage order = orders[_orderId];
        require(order.isActive, "Order not active");

        // Mark as inactive immediately to prevent re-entrancy
        order.isActive = false;

        // Calculate fee (if any)
        uint256 feeAmount = (order.priceInUSDT * FEE_BPS) / 10000;
        uint256 sellerAmount = order.priceInUSDT - feeAmount;

        // Transfer USDT from Buyer -> Selller (+ Fee)
        if (feeAmount > 0 && feeRecipient != address(0)) {
            usdt.safeTransferFrom(msg.sender, feeRecipient, feeAmount);
        }
        usdt.safeTransferFrom(msg.sender, order.seller, sellerAmount);

        // Transfer DeAI Token from Escrow -> Buyer
        tokenDeAI.safeTransfer(msg.sender, order.tokenAmount);

        emit OrderFulfilled(
            _orderId,
            msg.sender,
            order.seller,
            order.tokenAmount,
            order.priceInUSDT
        );
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

        order.isActive = false;

        // Return tokens to seller
        tokenDeAI.safeTransfer(msg.sender, order.tokenAmount);

        emit OrderCancelled(_orderId, msg.sender);
    }

    function getDocumentsCount() external view returns (uint256) {
        return orders.length;
    }
}
