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
  // Convert UIChatMessage[] (camelCase) to backend Message[] (PascalCase)
  const backendMessages = messages.map((msg) => ({
    ID: msg.id || '',
    Role: msg.role || '',
    Content: msg.content,
    CreatedAt: msg.createdAt || 0,
    UpdatedAt: msg.updatedAt || 0,
    Meta: msg.meta || {},
  }));

  // Call backend service
  const result = await processMessagesBackend({
    messages: backendMessages as any, // Type cast because of case mismatch
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

  // Convert backend response (PascalCase) back to frontend format (camelCase)
  return result.map((msg: any) => ({
    id: msg.ID || msg.id,
    role: msg.Role || msg.role,
    content: msg.Content || msg.content,
    createdAt: msg.CreatedAt || msg.createdAt,
    updatedAt: msg.UpdatedAt || msg.updatedAt,
    meta: msg.Meta || msg.meta,
  })) as OpenAIChatMessage[];
};

