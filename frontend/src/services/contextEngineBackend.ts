/**
 * Backend Context Engine Service
 * Wrapper for calling Go backend context engineering service via Wails bindings
 */

import { UIChatMessage } from '@/types';

// Import generated Wails bindings
import * as ContextEngineService from '@@/github.com/kawai-network/veridium';
import type {
  ContextEngineeringRequest as BackendRequest,
  ContextEngineeringResponse as BackendResponse,
} from '@@/github.com/kawai-network/veridium';

// Re-export types for convenience
export type { BackendRequest as ContextEngineeringRequest, BackendResponse as ContextEngineeringResponse };

/**
 * Call backend context engineering service
 * This replaces the frontend contextEngineering function with a backend call
 */
export async function processMessagesBackend(
  request: BackendRequest,
): Promise<UIChatMessage[]> {
  try {
    // Call Go backend service via Wails bindings
    const response = await ContextEngineService.ContextEngineService.ProcessMessages(request);

    if (response.error) {
      console.error('Backend context engineering error:', response.error);
      throw new Error(response.error);
    }

    return response.messages as UIChatMessage[];
  } catch (error) {
    console.error('Failed to call backend context engineering:', error);
    throw error;
  }
}

/**
 * Get backend engine statistics
 */
export async function getEngineStats(): Promise<Record<string, any>> {
  try {
    return await ContextEngineService.ContextEngineService.GetEngineStats();
  } catch (error) {
    console.error('Failed to get engine stats:', error);
    throw error;
  }
}

/**
 * Validate backend engine configuration
 */
export async function validateConfig(config: any): Promise<{ valid: boolean; errors: string[] }> {
  try {
    const result = await ContextEngineService.ContextEngineService.ValidateConfig(JSON.stringify(config));
    return result as { valid: boolean; errors: string[] };
  } catch (error) {
    console.error('Failed to validate config:', error);
    throw error;
  }
}

