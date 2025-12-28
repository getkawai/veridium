# Requirements Document

## Introduction

The P2P Marketplace is an internal trading platform that enables Contributors to sell their earned KAWAI tokens to Investors using USDT, without requiring initial liquidity pools. This marketplace uses smart contract escrow to ensure secure atomic swaps and allows natural price discovery through supply and demand forces.

## Glossary

- **Contributor**: Users who earn KAWAI tokens by providing AI compute services
- **Investor**: Users who want to buy KAWAI tokens using USDT
- **Order**: A sell listing created by Contributors offering KAWAI tokens at a specific USDT price
- **Escrow_System**: Smart contract that holds tokens/USDT during trades to ensure atomic swaps
- **OTC_Market**: The deployed smart contract handling peer-to-peer trading operations
- **Order_Book**: The collection of all active buy/sell orders in the marketplace
- **Atomic_Swap**: A trade execution where both parties' assets are exchanged simultaneously or not at all

## Requirements

### Requirement 1: Order Creation and Management

**User Story:** As a Contributor, I want to create sell orders for my KAWAI tokens, so that I can convert them to USDT when I need liquidity.

#### Acceptance Criteria

1. WHEN a Contributor creates a sell order, THE Order_Management_System SHALL validate the user has sufficient KAWAI token balance
2. WHEN creating an order, THE Order_Management_System SHALL require token amount, USDT price, and seller wallet address
3. WHEN an order is created, THE Escrow_System SHALL lock the specified KAWAI tokens from the seller's wallet
4. WHEN an order is successfully created, THE Order_Management_System SHALL assign a unique order ID and store it in the Order_Book
5. WHEN a Contributor cancels their order, THE Escrow_System SHALL return the locked KAWAI tokens to their wallet

### Requirement 2: Order Discovery and Browsing

**User Story:** As an Investor, I want to browse available KAWAI token orders, so that I can find tokens at prices I'm willing to pay.

#### Acceptance Criteria

1. WHEN an Investor visits the marketplace, THE Order_Display_System SHALL show all active sell orders sorted by price
2. WHEN displaying orders, THE Order_Display_System SHALL show token amount, USDT price, price per token, and seller address
3. WHEN filtering orders, THE Order_Display_System SHALL allow sorting by price, amount, and creation date
4. WHEN an order is partially filled, THE Order_Display_System SHALL update the remaining available amount in real-time
5. WHEN an order is completed or cancelled, THE Order_Display_System SHALL remove it from the active listings

### Requirement 3: Secure Trade Execution

**User Story:** As an Investor, I want to buy KAWAI tokens safely, so that I can be confident the trade will complete atomically without risk of loss.

#### Acceptance Criteria

1. WHEN an Investor initiates a buy order, THE Trade_Execution_System SHALL validate the buyer has sufficient USDT balance
2. WHEN executing a trade, THE Escrow_System SHALL simultaneously transfer KAWAI tokens to buyer and USDT to seller
3. IF either transfer fails during execution, THEN THE Escrow_System SHALL revert both transfers to maintain atomic swap property
4. WHEN a trade completes successfully, THE Trade_Execution_System SHALL emit a trade completion event
5. WHEN a partial buy occurs, THE Order_Management_System SHALL update the remaining order amount and keep it active

### Requirement 4: Order History and Tracking

**User Story:** As a marketplace user, I want to view my trading history, so that I can track my past transactions and current order status.

#### Acceptance Criteria

1. WHEN a user requests their order history, THE History_System SHALL return all their past orders with status and timestamps
2. WHEN displaying trade history, THE History_System SHALL show order details, execution price, and completion status
3. WHEN an order status changes, THE Notification_System SHALL update the user interface in real-time
4. WHEN a user has active orders, THE Order_Tracking_System SHALL display current status and remaining amounts
5. WHEN viewing order details, THE History_System SHALL show creation time, last update, and transaction hashes

### Requirement 5: Price Discovery and Market Data

**User Story:** As a marketplace participant, I want to see market pricing information, so that I can make informed trading decisions.

#### Acceptance Criteria

1. WHEN users view the marketplace, THE Market_Data_System SHALL display current lowest ask price and highest recent trade price
2. WHEN calculating market metrics, THE Analytics_System SHALL compute 24-hour trading volume and price ranges
3. WHEN showing price trends, THE Market_Data_System SHALL display recent trade history with timestamps and prices
4. WHEN orders are created at different prices, THE Price_Discovery_System SHALL enable natural price formation through supply and demand
5. WHEN displaying market depth, THE Order_Book_System SHALL show distribution of orders across different price levels

### Requirement 6: Smart Contract Integration

**User Story:** As a system administrator, I want the marketplace to integrate securely with deployed smart contracts, so that all trades are executed on-chain with proper validation.

#### Acceptance Criteria

1. WHEN interfacing with blockchain, THE Blockchain_Client SHALL connect to the deployed OTC_Market contract on Monad testnet
2. WHEN creating orders, THE Contract_Interface SHALL call the smart contract's createOrder function with proper parameters
3. WHEN executing trades, THE Contract_Interface SHALL invoke the buyOrder function ensuring atomic execution
4. WHEN cancelling orders, THE Contract_Interface SHALL call cancelOrder and verify token return to seller
5. WHEN contract events are emitted, THE Event_Listener SHALL capture and process OrderCreated, OrderFilled, and OrderCancelled events

### Requirement 7: Wails Service Integration and Data Management

**User Story:** As a frontend developer, I want Wails service bindings for marketplace operations, so that I can build a responsive desktop user interface.

#### Acceptance Criteria

1. WHEN the frontend calls marketplace service methods, THE Marketplace_Service SHALL validate input data and create orders via smart contract
2. WHEN the frontend requests order listings, THE Marketplace_Service SHALL return structured data with filtering and sorting capabilities
3. WHEN processing buy requests through service calls, THE Trade_Service SHALL execute trades and return transaction status and details
4. WHEN storing order data, THE Data_Storage_System SHALL persist order information in the multi-namespace KV store
5. WHEN service method errors occur, THE Error_Handling_System SHALL return appropriate Go error types with descriptive messages

### Requirement 8: Wallet Integration and Authorization

**User Story:** As a marketplace user, I want to use my connected wallet for authentication, so that only I can manage my orders and execute trades within the desktop application.

#### Acceptance Criteria

1. WHEN a user's wallet is connected in the desktop app, THE Wallet_Integration_System SHALL use the active wallet address for marketplace operations
2. WHEN creating orders through Wails services, THE Authorization_System SHALL ensure only the connected wallet owner can create orders for their tokens
3. WHEN cancelling orders via service calls, THE Authorization_System SHALL verify the requester is the original order creator
4. WHEN accessing user-specific data through services, THE Access_Control_System SHALL filter results based on connected wallet address
5. WHEN Wails service methods require wallet operations, THE Wallet_Service SHALL handle transaction signing and blockchain interactions