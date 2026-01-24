// SPDX-License-Identifier: MIT
pragma solidity ^0.8.33;

import {Test, console} from "forge-std/Test.sol";
import {OTCMarket} from "../contracts/OTCMarket.sol";
import {KawaiToken} from "../contracts/KawaiToken.sol";
import {MockStablecoin} from "../contracts/MockStablecoin.sol";

contract OTCMarketTest is Test {
    OTCMarket public market;
    KawaiToken public kawai;
    MockStablecoin public usdt;

    address public admin = address(1);
    address public seller = address(2);
    address public buyer = address(3);
    address public feeRecipient = address(4);

    uint256 constant KAWAI_AMOUNT = 1000 ether;
    uint256 constant USDT_PRICE = 500 ether;

    // Events (must be declared in test contract for expectEmit)
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
    
    event OrderPartiallyFilled(
        uint256 indexed orderId,
        address indexed buyer,
        address indexed seller,
        uint256 amountFilled,
        uint256 remainingAmount,
        uint256 pricePaid
    );

    function setUp() public {
        // Deploy contracts
        vm.startPrank(admin);
        kawai = new KawaiToken(admin, admin);
        usdt = new MockStablecoin();
        market = new OTCMarket(address(kawai), address(usdt), feeRecipient);
        vm.stopPrank();

        // Setup seller
        vm.startPrank(admin);
        kawai.mint(seller, KAWAI_AMOUNT * 10);
        vm.stopPrank();

        vm.startPrank(seller);
        kawai.approve(address(market), type(uint256).max);
        vm.stopPrank();

        // Setup buyer
        vm.startPrank(admin);
        usdt.mint(buyer, USDT_PRICE * 10);
        vm.stopPrank();

        vm.startPrank(buyer);
        usdt.approve(address(market), type(uint256).max);
        vm.stopPrank();
    }

    // ========================================
    // Test: Order Creation
    // ========================================

    function testCreateOrder() public {
        vm.startPrank(seller);
        
        uint256 sellerBalanceBefore = kawai.balanceOf(seller);
        
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);
        
        // Check seller balance decreased
        assertEq(kawai.balanceOf(seller), sellerBalanceBefore - KAWAI_AMOUNT);
        
        // Check order created correctly
        (
            uint256 id,
            address orderSeller,
            uint256 tokenAmount,
            uint256 priceInUSDT,
            uint256 remainingAmount,
            bool isActive
        ) = market.getOrder(0);
        
        assertEq(id, 0);
        assertEq(orderSeller, seller);
        assertEq(tokenAmount, KAWAI_AMOUNT);
        assertEq(priceInUSDT, USDT_PRICE);
        assertEq(remainingAmount, KAWAI_AMOUNT);  // ✅ Initially full
        assertTrue(isActive);
        
        vm.stopPrank();
    }

    function testCreateOrderRevertsIfAmountZero() public {
        vm.startPrank(seller);
        vm.expectRevert("Amount must be > 0");
        market.createOrder(0, USDT_PRICE);
        vm.stopPrank();
    }

    function testCreateOrderRevertsIfPriceZero() public {
        vm.startPrank(seller);
        vm.expectRevert("Price must be > 0");
        market.createOrder(KAWAI_AMOUNT, 0);
        vm.stopPrank();
    }

    // ========================================
    // Test: Full Order Buy
    // ========================================

    function testBuyOrderFull() public {
        // Create order
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);

        // Buy full order
        uint256 buyerKawaiBefore = kawai.balanceOf(buyer);
        uint256 sellerUSDTBefore = usdt.balanceOf(seller);

        vm.prank(buyer);
        market.buyOrder(0);

        // Check balances
        assertEq(kawai.balanceOf(buyer), buyerKawaiBefore + KAWAI_AMOUNT);
        assertEq(usdt.balanceOf(seller), sellerUSDTBefore + USDT_PRICE);

        // Check order is now inactive
        (, , , , uint256 remainingAmount, bool isActive) = market.getOrder(0);
        assertEq(remainingAmount, 0);
        assertFalse(isActive);
    }

    // ========================================
    // Test: Partial Order Buy
    // ========================================

    function testBuyOrderPartial() public {
        // Create order
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);

        // Buy 30% of order
        uint256 partialAmount = KAWAI_AMOUNT * 30 / 100;
        uint256 partialPrice = USDT_PRICE * 30 / 100;

        uint256 buyerKawaiBefore = kawai.balanceOf(buyer);
        uint256 sellerUSDTBefore = usdt.balanceOf(seller);

        vm.prank(buyer);
        market.buyOrderPartial(0, partialAmount);

        // Check balances
        assertEq(kawai.balanceOf(buyer), buyerKawaiBefore + partialAmount);
        assertEq(usdt.balanceOf(seller), sellerUSDTBefore + partialPrice);

        // Check order still active with remaining amount
        (, , , , uint256 remainingAmount, bool isActive) = market.getOrder(0);
        assertEq(remainingAmount, KAWAI_AMOUNT - partialAmount);
        assertTrue(isActive);
    }

    function testBuyOrderPartialMultipleTimes() public {
        // Create order
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);

        // Buy 30% first
        uint256 firstAmount = KAWAI_AMOUNT * 30 / 100;
        vm.prank(buyer);
        market.buyOrderPartial(0, firstAmount);

        // Check remaining
        (, , , , uint256 remaining1, bool active1) = market.getOrder(0);
        assertEq(remaining1, KAWAI_AMOUNT - firstAmount);
        assertTrue(active1);

        // Buy another 40%
        uint256 secondAmount = KAWAI_AMOUNT * 40 / 100;
        vm.prank(buyer);
        market.buyOrderPartial(0, secondAmount);

        // Check remaining
        (, , , , uint256 remaining2, bool active2) = market.getOrder(0);
        assertEq(remaining2, KAWAI_AMOUNT - firstAmount - secondAmount);
        assertTrue(active2);

        // Buy remaining 30%
        uint256 thirdAmount = KAWAI_AMOUNT * 30 / 100;
        vm.prank(buyer);
        market.buyOrderPartial(0, thirdAmount);

        // Check order fully filled
        (, , , , uint256 remaining3, bool active3) = market.getOrder(0);
        assertEq(remaining3, 0);
        assertFalse(active3);
    }

    function testBuyOrderPartialRevertsIfExceedsRemaining() public {
        // Create order
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);

        // Try to buy more than available
        vm.prank(buyer);
        vm.expectRevert("Amount exceeds remaining");
        market.buyOrderPartial(0, KAWAI_AMOUNT + 1);
    }

    function testBuyOrderPartialRevertsIfAmountZero() public {
        // Create order
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);

        // Try to buy zero
        vm.prank(buyer);
        vm.expectRevert("Amount must be > 0");
        market.buyOrderPartial(0, 0);
    }

    // ========================================
    // Test: Order Cancellation
    // ========================================

    function testCancelOrder() public {
        // Create order
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);

        uint256 sellerBalanceBefore = kawai.balanceOf(seller);

        // Cancel order
        vm.prank(seller);
        market.cancelOrder(0);

        // Check tokens returned
        assertEq(kawai.balanceOf(seller), sellerBalanceBefore + KAWAI_AMOUNT);

        // Check order inactive
        (, , , , uint256 remainingAmount, bool isActive) = market.getOrder(0);
        assertEq(remainingAmount, 0);
        assertFalse(isActive);
    }

    function testCancelOrderPartiallyFilled() public {
        // Create order
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);

        // Buy 40%
        uint256 partialAmount = KAWAI_AMOUNT * 40 / 100;
        vm.prank(buyer);
        market.buyOrderPartial(0, partialAmount);

        uint256 sellerBalanceBefore = kawai.balanceOf(seller);
        uint256 expectedRefund = KAWAI_AMOUNT - partialAmount;

        // Cancel remaining
        vm.prank(seller);
        market.cancelOrder(0);

        // Check only remaining tokens returned
        assertEq(kawai.balanceOf(seller), sellerBalanceBefore + expectedRefund);

        // Check order inactive
        (, , , , uint256 remainingAmount, bool isActive) = market.getOrder(0);
        assertEq(remainingAmount, 0);
        assertFalse(isActive);
    }

    function testCancelOrderRevertsIfNotSeller() public {
        // Create order
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);

        // Try to cancel as buyer
        vm.prank(buyer);
        vm.expectRevert("Not your order");
        market.cancelOrder(0);
    }

    // ========================================
    // Test: View Functions
    // ========================================

    function testGetOrder() public {
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);

        (
            uint256 id,
            address orderSeller,
            uint256 tokenAmount,
            uint256 priceInUSDT,
            uint256 remainingAmount,
            bool isActive
        ) = market.getOrder(0);

        assertEq(id, 0);
        assertEq(orderSeller, seller);
        assertEq(tokenAmount, KAWAI_AMOUNT);
        assertEq(priceInUSDT, USDT_PRICE);
        assertEq(remainingAmount, KAWAI_AMOUNT);
        assertTrue(isActive);
    }

    function testGetActiveOrders() public {
        // Create 3 orders
        vm.startPrank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);
        market.createOrder(KAWAI_AMOUNT * 2, USDT_PRICE * 2);
        market.createOrder(KAWAI_AMOUNT * 3, USDT_PRICE * 3);
        vm.stopPrank();

        // Get active orders
        OTCMarket.Order[] memory activeOrders = market.getActiveOrders(0, 10);
        
        assertEq(activeOrders.length, 3);
        assertEq(activeOrders[0].tokenAmount, KAWAI_AMOUNT);
        assertEq(activeOrders[1].tokenAmount, KAWAI_AMOUNT * 2);
        assertEq(activeOrders[2].tokenAmount, KAWAI_AMOUNT * 3);
    }

    function testGetActiveOrdersAfterPartialFill() public {
        // Create 2 orders
        vm.startPrank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);
        market.createOrder(KAWAI_AMOUNT * 2, USDT_PRICE * 2);
        vm.stopPrank();

        // Fully fill first order
        vm.prank(buyer);
        market.buyOrder(0);

        // Get active orders (should only return second order)
        OTCMarket.Order[] memory activeOrders = market.getActiveOrders(0, 10);
        
        assertEq(activeOrders.length, 1);
        assertEq(activeOrders[0].id, 1);
        assertEq(activeOrders[0].tokenAmount, KAWAI_AMOUNT * 2);
    }

    function testGetOrdersBySeller() public {
        address seller2 = address(5);
        
        // Setup seller2
        vm.prank(admin);
        kawai.mint(seller2, KAWAI_AMOUNT * 10);
        
        vm.prank(seller2);
        kawai.approve(address(market), type(uint256).max);

        // Create orders from different sellers
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);
        
        vm.prank(seller2);
        market.createOrder(KAWAI_AMOUNT * 2, USDT_PRICE * 2);
        
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT * 3, USDT_PRICE * 3);

        // Get orders by seller
        OTCMarket.Order[] memory sellerOrders = market.getOrdersBySeller(seller, 0, 10);
        
        assertEq(sellerOrders.length, 2);
        assertEq(sellerOrders[0].tokenAmount, KAWAI_AMOUNT);
        assertEq(sellerOrders[1].tokenAmount, KAWAI_AMOUNT * 3);
    }

    // ========================================
    // Test: Events
    // ========================================

    function testOrderCreatedEvent() public {
        vm.expectEmit(true, true, false, true);
        emit OrderCreated(0, seller, KAWAI_AMOUNT, USDT_PRICE);
        
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);
    }

    function testOrderPartiallyFilledEvent() public {
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);

        uint256 partialAmount = KAWAI_AMOUNT * 30 / 100;
        uint256 partialPrice = USDT_PRICE * 30 / 100;
        uint256 remaining = KAWAI_AMOUNT - partialAmount;

        vm.expectEmit(true, true, true, true);
        emit OrderPartiallyFilled(
            0,
            buyer,
            seller,
            partialAmount,
            remaining,
            partialPrice
        );
        
        vm.prank(buyer);
        market.buyOrderPartial(0, partialAmount);
    }

    function testOrderFulfilledEvent() public {
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);

        vm.expectEmit(true, true, true, true);
        emit OrderFulfilled(
            0,
            buyer,
            seller,
            KAWAI_AMOUNT,
            USDT_PRICE
        );
        
        vm.prank(buyer);
        market.buyOrder(0);
    }

    function testOrderCancelledEvent() public {
        vm.prank(seller);
        market.createOrder(KAWAI_AMOUNT, USDT_PRICE);

        vm.expectEmit(true, true, false, true);
        emit OrderCancelled(0, seller);
        
        vm.prank(seller);
        market.cancelOrder(0);
    }
}

