/**
 * Feature Flags Configuration
 * 
 * This file controls feature rollout for gradual migration from frontend to backend.
 * 
 * IMPORTANT: Always start with flags set to FALSE. Only enable after thorough testing.
 * Each flag can be toggled independently for easy rollback.
 */

export const FEATURE_FLAGS = {
  /**
   * Phase 1: Use backend for chat operations
   * 
   * When enabled, sendMessage() will call AgentChatService.Chat() instead of
   * the original frontend chat flow.
   * 
   * Status: 🧪 EXPERIMENTAL - Testing Phase
   * Default: false (use original frontend logic)
   * 
   * To enable: Set to true and test thoroughly
   * To rollback: Set back to false (instant rollback, no code changes needed)
   */
  USE_BACKEND_CHAT: false, // DISABLED: Backend has issues, use original frontend

  /**
   * Phase 3: Use backend tools engine
   * 
   * When enabled, tools will be executed via backend ToolsEngineBridge
   * instead of frontend tool orchestration.
   * 
   * Status: ⏳ TODO - Depends on Phase 1-2
   * Default: false
   */
  USE_BACKEND_TOOLS: false,

  /**
   * Phase 3: Use backend RAG
   * 
   * When enabled, knowledge base queries will go through backend RAGWorkflow
   * instead of frontend RAG implementation.
   * 
   * Status: ⏳ TODO - Depends on Phase 1-2
   * Default: false
   */
  USE_BACKEND_RAG: false,

  /**
   * Phase 5: Use backend streaming
   * 
   * When enabled, chat responses will stream via Wails events
   * instead of synchronous responses.
   * 
   * Status: ⏳ TODO - Depends on Wails v3 streaming API
   * Default: false
   */
  USE_BACKEND_STREAMING: false,

  /**
   * Development: Enable verbose logging for migration debugging
   * 
   * When enabled, logs all backend calls, fallbacks, and state changes
   * for easier debugging during migration.
   * 
   * Default: false (enable in development only)
   */
  DEBUG_MIGRATION: false,
} as const;

/**
 * Check if a feature flag is enabled
 * 
 * @param feature - The feature flag to check
 * @returns true if the feature is enabled, false otherwise
 */
export function isFeatureEnabled(feature: keyof typeof FEATURE_FLAGS): boolean {
  return FEATURE_FLAGS[feature];
}

/**
 * Get all feature flags status
 * 
 * Useful for debugging and monitoring which features are enabled.
 * 
 * @returns Object with all feature flags and their current values
 */
export function getFeatureFlagsStatus() {
  return {
    ...FEATURE_FLAGS,
    timestamp: new Date().toISOString(),
  };
}

/**
 * Log a migration event (only if DEBUG_MIGRATION is enabled)
 * 
 * @param event - Event name
 * @param data - Additional data to log
 */
export function logMigrationEvent(event: string, data?: any) {
  if (FEATURE_FLAGS.DEBUG_MIGRATION) {
    console.log(`[Migration] ${event}`, data || '');
  }
}

/**
 * Check if backend is available
 * 
 * This checks if Wails bindings are properly loaded.
 * Useful for development mode detection.
 * 
 * @returns true if backend bindings are available
 */
export function isBackendAvailable(): boolean {
  try {
    // Check if Wails runtime is available
    return typeof window !== 'undefined' && 'go' in window;
  } catch {
    return false;
  }
}

// Export types
export type FeatureFlag = keyof typeof FEATURE_FLAGS;
export type FeatureFlagsStatus = typeof FEATURE_FLAGS;

