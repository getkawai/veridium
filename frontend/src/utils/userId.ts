import { getUserStoreState } from '@/store/user';

const DEFAULT_USER_ID = 'DEFAULT_LOBE_CHAT_USER';

export const getResolvedUserId = (): string => {
  try {
    const walletAddress = getUserStoreState().walletAddress?.trim();
    const result = walletAddress || DEFAULT_USER_ID;
    console.log('[getResolvedUserId]', { walletAddress, result, usingFallback: !walletAddress });
    return result;
  } catch (error) {
    console.error('[getResolvedUserId] Error:', error);
    return DEFAULT_USER_ID;
  }
};

export const getRequiredUserId = (): string => {
  const userId = getResolvedUserId();
  if (!userId) {
    throw new Error('Wallet address is required');
  }
  return userId;
};
