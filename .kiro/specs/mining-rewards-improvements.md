# Mining Rewards System Improvements

## Overview
This spec addresses critical issues in the mining rewards system, specifically around claim confirmation, transaction tracking, and recent activity display.

## Current Issues

### 1. Broken Explorer Links
- **Problem**: Claims showing "Confirming..." with broken Etherscan links instead of Monad explorer
- **Root Cause**: Potential network configuration or transaction hash format issues
- **Impact**: Users cannot verify their transactions

### 2. Recent Activity Display Issues
- **Problem**: Recent Activity section not properly displaying completed claims
- **Root Cause**: Logic filtering confirmed claims from pending_proofs has type issues and incorrect data source
- **Impact**: Users lose visibility into their claim history

### 3. Type Safety Issues
- **Problem**: Multiple TypeScript errors in MiningRewardsSection component
- **Root Cause**: Improper type handling for claim status and transaction hashes
- **Impact**: Runtime errors and poor developer experience

## User Stories

### US-1: Proper Transaction Tracking
**As a** user who has submitted a mining claim  
**I want** to see the correct transaction status with working explorer links  
**So that** I can verify my transaction on the blockchain

**Acceptance Criteria:**
- [ ] Pending claims show "Confirming..." status with correct Monad explorer links
- [ ] Explorer links open to the correct transaction page
- [ ] Transaction hashes are properly formatted and validated
- [ ] Network configuration correctly points to Monad testnet

### US-2: Comprehensive Recent Activity
**As a** user who has successfully claimed rewards  
**I want** to see my recent claim history in the Recent Activity section  
**So that** I can track my mining reward claims over time

**Acceptance Criteria:**
- [ ] Recent Activity displays confirmed claims from the last 30 days
- [ ] Each entry shows: Type, Amount, Transaction Hash, Date, Status
- [ ] Transaction hashes are clickable and open Monad explorer
- [ ] Empty state shows helpful message when no claims exist
- [ ] Data refreshes automatically after successful claims

### US-3: Improved Claim Status Flow
**As a** user claiming mining rewards  
**I want** clear status updates throughout the claim process  
**So that** I understand what's happening with my transaction

**Acceptance Criteria:**
- [ ] Clear status progression: Unclaimed → Pending → Confirmed
- [ ] Pending claims show transaction hash immediately after submission
- [ ] Status updates automatically without manual refresh
- [ ] Error states are clearly communicated with actionable messages
- [ ] Gas estimation works correctly for all claim types

## Technical Requirements

### TR-1: Fix Type Safety Issues
- Fix TypeScript errors in MiningRewardsSection component
- Properly type claim status and transaction hash fields
- Remove unused imports and variables
- Add proper null checks for optional fields

### TR-2: Improve Data Flow
- Separate confirmed claims data source from pending_proofs
- Implement proper claim history storage in KV store
- Add automatic status polling for pending transactions
- Optimize data fetching to reduce API calls

### TR-3: Network Configuration
- Ensure currentNetwork.explorerURL points to correct Monad explorer
- Validate transaction hash format before creating explorer links
- Add fallback explorer URL if network config is missing
- Test explorer links across different environments

### TR-4: Enhanced Recent Activity
- Create dedicated endpoint for recent claim history
- Implement pagination for large claim histories
- Add filtering options (by token type, date range)
- Store claim completion timestamps for accurate sorting

## Implementation Plan

### Phase 1: Fix Critical Issues (Priority: High)
1. **Fix TypeScript Errors**
   - Resolve type issues in MiningRewardsSection component
   - Add proper type guards for claim status validation
   - Fix boolean type assignment errors

2. **Fix Explorer Links**
   - Validate network configuration for Monad explorer
   - Test transaction hash format and explorer URL construction
   - Add error handling for invalid transaction hashes

### Phase 2: Improve Recent Activity (Priority: High)
1. **Data Source Improvements**
   - Create proper confirmed claims data structure
   - Implement claim history storage in backend
   - Add automatic claim status updates

2. **UI Enhancements**
   - Improve Recent Activity table with proper data
   - Add loading states and error handling
   - Implement auto-refresh after successful claims

### Phase 3: Enhanced Features (Priority: Medium)
1. **Status Tracking**
   - Add real-time transaction status polling
   - Implement push notifications for claim confirmations
   - Add transaction receipt validation

2. **User Experience**
   - Add claim history export functionality
   - Implement advanced filtering and search
   - Add claim analytics and statistics

## Testing Strategy

### Unit Tests
- [ ] Test claim status validation logic
- [ ] Test transaction hash formatting
- [ ] Test explorer URL generation
- [ ] Test data filtering and pagination

### Integration Tests
- [ ] Test complete claim flow from submission to confirmation
- [ ] Test network switching and explorer URL updates
- [ ] Test error handling for failed transactions
- [ ] Test data persistence and retrieval

### User Acceptance Tests
- [ ] Test claim submission with various reward types
- [ ] Test transaction tracking across different networks
- [ ] Test Recent Activity display with multiple claims
- [ ] Test error scenarios and recovery flows

## Success Metrics

### Technical Metrics
- [ ] Zero TypeScript compilation errors
- [ ] 100% working explorer links
- [ ] < 2 second load time for Recent Activity
- [ ] 99% uptime for claim status updates

### User Experience Metrics
- [ ] Reduced support tickets about "broken links"
- [ ] Increased user engagement with Recent Activity section
- [ ] Improved claim success rate (fewer failed transactions)
- [ ] Positive user feedback on transaction visibility

## Dependencies

### Backend Services
- DeAIService.GetClaimableRewards()
- DeAIService.ClaimMiningReward()
- DeAIService.ClaimUSDTReward()
- JarvisService.EstimateGas()

### External Services
- Monad Explorer API
- Blockchain RPC endpoints
- Transaction status polling services

### Frontend Components
- Antd Table, Modal, Pagination components
- Browser.OpenURL for external links
- Wails runtime services

## Risk Assessment

### High Risk
- **Network Configuration Changes**: Could break existing functionality
- **Data Migration**: Moving from pending_proofs to dedicated claim history

### Medium Risk
- **Type System Changes**: Could introduce new compilation errors
- **UI Changes**: Could affect user workflow

### Low Risk
- **Explorer URL Updates**: Easy to revert if issues occur
- **Loading State Improvements**: Non-breaking enhancements

## Rollback Plan

1. **Immediate Rollback**: Revert to previous component version if critical errors
2. **Partial Rollback**: Disable new features while keeping bug fixes
3. **Data Rollback**: Restore previous data structure if migration fails
4. **Configuration Rollback**: Revert network settings if explorer links break

## Future Considerations

### Phase 4: Advanced Features
- Real-time WebSocket updates for claim status
- Mobile-responsive claim history interface
- Integration with wallet transaction history
- Advanced analytics and reporting dashboard

### Performance Optimizations
- Implement claim data caching
- Add virtual scrolling for large claim histories
- Optimize API calls with request batching
- Add offline support for claim history viewing