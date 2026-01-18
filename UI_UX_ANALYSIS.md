# Analysis of UI/UX and Potential Bugs - Veridium Frontend

## 1. UI/UX Analysis

### Theme Inconsistency (Light Mode)
- **Issue:** The `balanceCard` style in `src/app/wallet/wallet.tsx` (lines 52-97) uses hardcoded dark colors (`#1a1a2e`, `#16213e`) and gradients.
- **Impact:** If the application is switched to Light Mode, this card will remain dark, potentially clashing with the rest of the light theme.
- **Recommendation:** Use theme tokens (e.g., `token.colorBgContainer`, `token.colorPrimary`) instead of hardcoded hex values, or define specific light/dark overrides.

### "Smart Deposit" Experience
- **Issue:** The deposit flow waits for a fixed 2 seconds before attempting to sync with the backend.
- **Impact:** Blockchain transactions often take longer than 2 seconds. If the sync happens before the transaction is indexed, it fails, leaving the user with a "Deposit successful but sync failed" message.
- **Recommendation:** Implement a polling mechanism that checks for the transaction status or the balance update for a longer duration (e.g., up to 30 seconds) before failing. Provide a visible "Retry Sync" button.

### Send Form - Custom Token
- **Issue:** The `SendForm` component has an internal state logic for `customTokenAddress`, but the UI (Select dropdown) only offers "Native", "USDT", and "KAWAI".
- **Impact:** Users cannot easily send custom tokens via the UI.
- **Recommendation:** Add a "Custom Token" option to the dropdown that, when selected, reveals an input field for the token contract address.

### Network Filtering
- **Issue:** The `NetworkSwitcher` filters networks based on `backendConfig`. If the config is restricted (e.g., to "testnet"), the user cannot verify or see other networks even if they exist in the wallet.
- **Impact:** unexpected restrictions for power users who might want to switch networks for testing or other purposes.

## 2. Potential Bugs & Technical Issues

### Hardcoded Configuration
- **Issue:** `DEFAULT_CHAIN_ID` is hardcoded to `10143` (Monad Testnet) in `src/app/wallet/wallet.tsx`.
- **Impact:** If the default network changes, this code requires manual updates.
- **Recommendation:** Move this to a central configuration file (`src/config/network.ts` or similar).

### Error Handling in Transfers
- **Issue:** The `handleSend` function relies on string matching (e.g., "insufficient funds") to provide user-friendly error messages.
- **Impact:** If the underlying library or node error messages change (e.g., with a different RPC provider), the error matching might fail, showing generic errors.
- **Recommendation:** Use error codes if available from the `JarvisService` or the web3 provider.

### SetupForm Mnemonic Verification
- **Issue:** In `SetupForm` (Wallet Creation), the user is asked to save the mnemonic, but there is no step to verify they have actually saved it (e.g., "Enter word #4").
- **Impact:** A user might click "I have written it down" without actually doing so, leading to potential loss of funds if they forget their password.

## 3. Code Quality Observations

- **File Size:** `src/app/wallet/wallet.tsx` is over 1100 lines long. It contains multiple component definitions (`MenuContent`, `NetworkSwitcher`, `SendForm`, `AddTokenModal`, etc.).
- **Recommendation:** Refactor these sub-components into their own files (e.g., `src/app/wallet/components/NetworkSwitcher.tsx`) to improve readability and maintainability.

## 4. Specific Code Locations
- **`src/app/wallet/wallet.tsx`**: Main logic and potential bugs.
- **`src/app/wallet/HomeContent.tsx`**: `balanceCard` rendering.
