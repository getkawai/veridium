# Mining Rewards System - Requirements Document

## Executive Summary

The mining rewards system currently has critical issues affecting user experience:
1. **Broken Explorer Links**: Claims show "Confirming..." with non-functional Etherscan links instead of Monad explorer
2. **Empty Recent Activity**: The Recent Activity section doesn't display completed claims properly
3. **Type Safety Issues**: Multiple TypeScript errors causing potential runtime issues

This document outlines the requirements to fix these issues and improve the overall mining rewards experience.

## Problem Statement

### Issue 1: Incorrect Explorer Links
**Current State**: When users claim rewards, the pending claims show "Confirming..." status with links to `https://api.etherscan.io/v2/tx/[hash]` which return 404 errors.

**Expected State**: Links should point to Monad explorer (`https://testnet.monadexplorer.com/tx/[hash]`) and work correctly.

### Issue 2: Recent Activity Not Showing Claims
**Current State**: The Recent Activity section is empty even after successful claims.

**Expected State**: Recent Activity should display confirmed mining claims with proper transaction details.

### Issue 3: TypeScript Errors
**Current State**: Multiple type errors in the MiningRewardsSection component affecting code quality and potential runtime stability.

**Expected State**: Clean TypeScript compilation with proper type safety.

## Functional Requirements

### FR-1: Correct Explorer Integration
- **FR-1.1**: All transaction links must point to Monad explorer
- **FR-1.2**: Explorer URLs must be constructed using `currentNetwork.explorerURL` or fallback to `https://testnet.monadexplorer.com`
- **FR-1.3**: Transaction hashes must be validated before creating links
- **FR-1.4**: Links must open in external browser correctly

### FR-2: Recent Activity Display
- **FR-2.1**: Recent Activity must show confirmed mining claims from the last 30 days
- **FR-2.2**: Each entry must display: Type (Mining Claim), Amount, Transaction Hash, Date, Status (Confirmed)
- **FR-2.3**: Transaction hashes must be clickable and open Monad explorer
- **FR-2.4**: Empty state must show informative message when no claims exist
- **FR-2.5**: Data must refresh automatically after successful claims

### FR-3: Claim Status Tracking
- **FR-3.1**: Claims must progress through states: Unclaimed → Pending → Confirmed
- **FR-3.2**: Pending claims must show transaction hash immediately after submission
- **FR-3.3**: Status must update automatically without manual refresh
- **FR-3.4**: Error states must be clearly communicated

## Technical Requirements

### TR-1: Type Safety
- **TR-1.1**: Fix all TypeScript compilation errors
- **TR-1.2**: Add proper type guards for claim status validation
- **TR-1.3**: Remove unused imports and variables
- **TR-1.4**: Add null checks for optional fields

### TR-2: Data Management
- **TR-2.1**: Implement proper data source for confirmed claims
- **TR-2.2**: Store claim completion data in KV store
- **TR-2.3**: Add automatic status polling for pending transactions
- **TR-2.4**: Optimize API calls to reduce load times

### TR-3: Network Configuration
- **TR-3.1**: Validate network configuration for correct explorer URLs
- **TR-3.2**: Implement fallback explorer URL mechanism
- **TR-3.3**: Test explorer links across different environments
- **TR-3.4**: Add error handling for network configuration issues

## User Experience Requirements

### UX-1: Visual Feedback
- **UX-1.1**: Loading states must be clear and informative
- **UX-1.2**: Error messages must be actionable and user-friendly
- **UX-1.3**: Success confirmations must include transaction links
- **UX-1.4**: Status changes must be visually distinct

### UX-2: Information Architecture
- **UX-2.1**: Recent Activity must be easily scannable
- **UX-2.2**: Transaction details must be accessible but not overwhelming
- **UX-2.3**: Navigation to external explorer must be intuitive
- **UX-2.4**: Empty states must guide users on next actions

## Acceptance Criteria

### AC-1: Explorer Links Work Correctly
- [ ] All pending claims show working Monad explorer links
- [ ] Links open to correct transaction pages
- [ ] No 404 errors when clicking transaction links
- [ ] Fallback URL works when network config is missing

### AC-2: Recent Activity Functions Properly
- [ ] Confirmed claims appear in Recent Activity within 30 seconds
- [ ] All required columns display correct data
- [ ] Transaction hashes are clickable and functional
- [ ] Empty state shows when no claims exist
- [ ] Data refreshes after new claims

### AC-3: Type Safety Achieved
- [ ] Zero TypeScript compilation errors
- [ ] No runtime type-related errors
- [ ] Proper null/undefined handling
- [ ] Clean code without unused variables

### AC-4: User Experience Improved
- [ ] Clear status progression for claims
- [ ] Informative loading and error states
- [ ] Intuitive navigation to transaction details
- [ ] Responsive design works on all screen sizes

## Implementation Priority

### P0 (Critical - Fix Immediately)
1. Fix TypeScript compilation errors
2. Correct explorer URL configuration
3. Fix Recent Activity data source

### P1 (High - Next Sprint)
1. Implement proper claim status tracking
2. Add automatic data refresh
3. Improve error handling and messaging

### P2 (Medium - Future Sprint)
1. Add advanced filtering for Recent Activity
2. Implement real-time status updates
3. Add claim history export functionality

## Success Metrics

### Technical Metrics
- Zero TypeScript compilation errors
- 100% functional explorer links
- < 2 second load time for Recent Activity
- 99% uptime for claim status updates

### User Metrics
- Reduced support tickets about broken links
- Increased engagement with Recent Activity
- Improved claim success rate
- Positive user feedback scores

## Dependencies

### Internal Dependencies
- DeAIService API methods
- KV store for claim data persistence
- Network configuration management
- Wails runtime services

### External Dependencies
- Monad Explorer API availability
- Blockchain RPC endpoint stability
- Browser compatibility for external links

## Risk Mitigation

### Technical Risks
- **Risk**: Network configuration changes break existing functionality
- **Mitigation**: Implement fallback URLs and comprehensive testing

- **Risk**: Data migration causes temporary data loss
- **Mitigation**: Implement gradual migration with rollback capability

### User Experience Risks
- **Risk**: UI changes confuse existing users
- **Mitigation**: Maintain familiar interface while fixing underlying issues

- **Risk**: Performance degradation from additional API calls
- **Mitigation**: Implement caching and optimize data fetching

## Definition of Done

A requirement is considered complete when:
1. All acceptance criteria are met
2. Code passes TypeScript compilation
3. Unit tests cover new functionality
4. Integration tests verify end-to-end flow
5. User acceptance testing confirms improved experience
6. Documentation is updated
7. Code review is completed and approved