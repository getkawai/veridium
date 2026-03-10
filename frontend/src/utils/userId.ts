import { getUserStoreState } from '@/store/user';

export const getResolvedUserId = (): string => {
  try {
    return (getUserStoreState().walletAddress || '').trim();
  } catch {
    return '';
  }
};

export const getRequiredUserId = (): string => {
  const userId = getResolvedUserId();
  if (!userId) {
    throw new Error('Wallet address is required');
  }
  return userId;
};
