# Free Trial System Implementation Plan

## Overview
Implement a "Free Credits" system where users can claim a one-time trial balance (e.g., 5 USDT) to test the platform. The system uses the existing USDT balance structure but adds a mechanism to prevent abuse (one claim per wallet).

## Architecture
- **Model:** Direct USDT Balance (Virtual).
- **Storage:** Cloudflare KV (via `pkg/store`).
- **Logic:**
  - `HasClaimedTrial(address)`: Checks if a `trial:{address}` key exists.
  - `ClaimFreeTrial(address)`:
    1. Checks eligibility.
    2. Adds 5,000,000 micro USDT (5 USDT) to user's balance.
    3. Sets `trial:{address}` to preventing double claiming.

## Files to Modify

### 1. `pkg/store/balance.go`
- Add `HasClaimedTrial(ctx, address) (bool, error)`
- Add `ClaimFreeTrial(ctx, address) error`
- Define `formatted_key` for trial tracking.

### 2. `pkg/gateway/handler.go`
- Add `HandleClaimTrial` method.
- Validate user authentication/wallet signature if necessary (or just address for now if that's how it works).

### 3. `pkg/gateway/server.go`
- Register the new route: `POST /v1/user/claim-trial`

### 4. `internal/constant/economics.go` (Create/Modify)
- Define `FREE_TRIAL_AMOUNT_MICRO_USDT = 5000000`

## Step-by-Step Implementation

1.  **Define Constants**: Add the constant for the trial amount.
2.  **Store Implementation**: Implement the KV logic in `pkg/store/balance.go`.
3.  **Handler Implementation**: Add the HTTP handler in `pkg/gateway/handler.go`.
4.  **Route Registration**: Wire it up in `pkg/gateway/server.go`.
5.  **Testing**: Verify the flow.

## Future Improvements (Not in this iteration)
- Social verification (Twitter/Discord).
- Referral system.
