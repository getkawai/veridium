// Context engine removed - using backend for context processing
// import { INBOX_GUIDE_SYSTEMROLE, INBOX_SESSION_ID, isDesktop, isServerMode } from '@/const';
// import {
//   ContextEngine,
//   HistorySummaryProvider,
//   HistoryTruncateProcessor,
//   InboxGuideProvider,
//   InputTemplateProcessor,
//   MessageCleanupProcessor,
//   MessageContentProcessor,
//   PlaceholderVariablesProcessor,
//   SystemRoleInjector,
//   ToolCallProcessor,
//   ToolMessageReorder,
//   ToolNameResolver,
//   ToolSystemRoleProvider,
// } from '@/context-engine';
// import { historySummaryPrompt } from '@/prompts';
import { OpenAIChatMessage, UIChatMessage } from '@/types';
// import { VARIABLE_GENERATORS } from '@/utils/client';

// import { isCanUseFC } from '@/helpers/isCanUseFC';
// import { getToolStoreState } from '@/store/tool';
// import { toolSelectors } from '@/store/tool/selectors';

// import { isCanUseVideo, isCanUseVision } from './helper';

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

// Context engineering is now handled by the backend
// This function is kept for backward compatibility but just passes messages through
export const contextEngineering = async ({
  messages = [],
}: ContextEngineeringContext): Promise<OpenAIChatMessage[]> => {
  // Simply return messages as-is, backend handles context processing
  return messages.map((msg) => ({
    role: msg.role,
    content: msg.content,
  })) as OpenAIChatMessage[];
};
