import { DEFAULT_FEATURE_FLAGS, mapFeatureFlagsEnvToState } from './schema';

export const getServerFeatureFlagsValue = () => DEFAULT_FEATURE_FLAGS;

/**
 * Get feature flags from EdgeConfig with fallback to environment variables
 * @param userId - Optional user ID for user-specific feature flag evaluation
 */
export const getServerFeatureFlagsFromEdgeConfig = async (userId?: string) => {
  return DEFAULT_FEATURE_FLAGS;
};

export const serverFeatureFlags = (userId?: string) => {
  const serverConfig = getServerFeatureFlagsValue();

  return mapFeatureFlagsEnvToState(serverConfig, userId);
};

/**
 * Get server feature flags from EdgeConfig and map them to state with user ID
 * @param userId - Optional user ID for user-specific feature flag evaluation
 */
export const getServerFeatureFlagsStateFromEdgeConfig = (userId?: string) => {
  return mapFeatureFlagsEnvToState(DEFAULT_FEATURE_FLAGS, userId);
};

export * from './schema';
