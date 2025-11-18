/**
 * Backend-powered Context Engineering
 * This replaces the frontend context engineering with backend processing
 */

import { OpenAIChatMessage, UIChatMessage } from '@/types';
import { processMessagesBackend } from '@/services/contextEngineBackend';

interface ContextEngineeringContext {
  enableHistoryCount?: boolean;
  historyCount?: number;
  historySummary?: string;
  inputTemplate?: string;
  isWelcomeQuestion?: boolean;
  messages: UIChatMessage[];
  model: string;
  provider: string;
  sessionId?: string;
  systemRole?: string;
  tools?: string[];
}

/**
 * Context engineering using backend Go service
 * This provides the same interface as the frontend version but uses backend processing
 */
export const contextEngineeringBackend = async ({
  messages = [],
  tools,
  model,
  provider,
  systemRole,
  inputTemplate,
  enableHistoryCount,
  historyCount,
  historySummary,
  sessionId,
  isWelcomeQuestion,
}: ContextEngineeringContext): Promise<OpenAIChatMessage[]> => {
  // Call backend service
  const result = await processMessagesBackend({
    messages,
    tools,
    model,
    provider,
    systemRole,
    inputTemplate,
    enableHistoryCount,
    historyCount,
    historySummary,
    sessionId,
    isWelcomeQuestion,
  });

  // Backend returns the same format as frontend
  return result as OpenAIChatMessage[];
};

