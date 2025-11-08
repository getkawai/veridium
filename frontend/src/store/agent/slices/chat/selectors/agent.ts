import {
  DEFAULT_AGENT_CONFIG,
  DEFAULT_MODEL,
  DEFAULT_PROVIDER,
  INBOX_SESSION_ID,
} from '@/const';
import { KnowledgeItem, KnowledgeType, LobeAgentConfig } from '@/types';

import { DEFAULT_OPENING_QUESTIONS } from '@/features/AgentSetting/store/selectors';
import { filterToolIds } from '@/helpers/toolFilters';
import { AgentStoreState } from '@/store/agent/initialState';
import { merge } from '@/utils/merge';

const isInboxSession = (s: AgentStoreState) => s.activeId === INBOX_SESSION_ID;

// ==========   Config   ============== //

const inboxAgentConfig = (s: AgentStoreState) =>
  merge(DEFAULT_AGENT_CONFIG, s.agentMap[INBOX_SESSION_ID]);
const inboxAgentModel = (s: AgentStoreState) => inboxAgentConfig(s).model;

const getAgentConfigById =
  (id: string) =>
  (s: AgentStoreState): LobeAgentConfig =>
    merge(s.defaultAgentConfig, s.agentMap[id]);

const getAgentConfigByAgentId =
  (agentId: string) =>
  (s: AgentStoreState): LobeAgentConfig => {
    // Find the session that contains this agent
    const sessionId = Object.keys(s.agentMap).find((sessionKey) => {
      const agentConfig = s.agentMap[sessionKey];
      return agentConfig?.id === agentId;
    });

    if (sessionId) {
      return merge(s.defaultAgentConfig, s.agentMap[sessionId]);
    }

    // Fallback to default config if agent not found
    return s.defaultAgentConfig;
  };

export const currentAgentConfig = (s: AgentStoreState): LobeAgentConfig =>
  getAgentConfigById(s.activeId)(s);

const currentAgentSystemRole = (s: AgentStoreState) => {
  return currentAgentConfig(s).systemRole;
};

const currentAgentModel = (s: AgentStoreState): string => {
  const config = currentAgentConfig(s);
  const model = config?.model;
  
  // Handle NullString from database (Go type with {String: string, Valid: boolean})
  if (model && typeof model === 'object' && 'Valid' in model && 'String' in model) {
    return (model as any).Valid ? (model as any).String : DEFAULT_MODEL;
  }

  return model || DEFAULT_MODEL;
};

const currentAgentModelProvider = (s: AgentStoreState) => {
  const config = currentAgentConfig(s);

  return config?.provider || DEFAULT_PROVIDER;
};

const currentAgentPlugins = (s: AgentStoreState) => {
  const config = currentAgentConfig(s);

  return Array.isArray(config?.plugins) ? config.plugins : [];
};

/**
 * Get displayable agent plugins by filtering out platform-specific tools
 * that shouldn't be shown in the current environment
 */
const displayableAgentPlugins = (s: AgentStoreState) => {
  const plugins = currentAgentPlugins(s);
  return filterToolIds(plugins);
};

const currentAgentKnowledgeBases = (s: AgentStoreState) => {
  const config = currentAgentConfig(s);

  return Array.isArray(config?.knowledgeBases) ? config.knowledgeBases : [];
};

const currentAgentFiles = (s: AgentStoreState) => {
  const config = currentAgentConfig(s);

  return Array.isArray(config?.files) ? config.files : [];
};

const currentEnabledKnowledge = (s: AgentStoreState) => {
  const knowledgeBases = currentAgentKnowledgeBases(s);
  const files = currentAgentFiles(s);

  return [
    ...files
      .filter((f) => f.enabled)
      .map((f) => ({ fileType: f.type, id: f.id, name: f.name, type: KnowledgeType.File })),
    ...knowledgeBases
      .filter((k) => k.enabled)
      .map((k) => ({ id: k.id, name: k.name, type: KnowledgeType.KnowledgeBase })),
  ] as KnowledgeItem[];
};

const hasSystemRole = (s: AgentStoreState) => {
  const config = currentAgentConfig(s);

  return !!config.systemRole;
};

const hasKnowledgeBases = (s: AgentStoreState) => {
  const knowledgeBases = currentAgentKnowledgeBases(s);

  return knowledgeBases.length > 0;
};

const hasFiles = (s: AgentStoreState) => {
  const files = currentAgentFiles(s);

  return files.length > 0;
};

const hasKnowledge = (s: AgentStoreState) => hasKnowledgeBases(s) || hasFiles(s);
const hasEnabledKnowledge = (s: AgentStoreState) => currentEnabledKnowledge(s).length > 0;
const currentKnowledgeIds = (s: AgentStoreState) => {
  return {
    fileIds: currentAgentFiles(s)
      .filter((item) => item.enabled)
      .map((f) => f.id),
    knowledgeBaseIds: currentAgentKnowledgeBases(s)
      .filter((item) => item.enabled)
      .map((k) => k.id),
  };
};

const isAgentConfigLoading = (s: AgentStoreState) => {
  // During centralized loading phase, show loading if batch hasn't completed
  if (!s.isAllAgentConfigsLoaded) return true;
  
  // After batch loading, check if specific session config is loaded
  return !s.agentConfigInitMap[s.activeId];
};

const openingQuestions = (s: AgentStoreState) => {
  const questions = currentAgentConfig(s).openingQuestions;
  return Array.isArray(questions) ? questions : DEFAULT_OPENING_QUESTIONS;
};
const openingMessage = (s: AgentStoreState) => currentAgentConfig(s).openingMessage || '';

// TTS voice selector - returns the configured TTS voice for the current agent
const currentAgentTTSVoice = (lang?: string) => (s: AgentStoreState) => {
  // TODO: Implement proper TTS voice selection based on agent config and language
  // For now, return undefined to disable TTS voice comparison
  // This was previously connected to lobe-chat backend, needs implementation for go backend
  return undefined;
};

export const agentSelectors = {
  currentAgentConfig,
  currentAgentFiles,
  currentAgentKnowledgeBases,
  currentAgentModel,
  currentAgentModelProvider,
  currentAgentPlugins,
  currentAgentSystemRole,
  currentAgentTTSVoice,
  currentEnabledKnowledge,
  currentKnowledgeIds,
  displayableAgentPlugins,
  getAgentConfigByAgentId,
  getAgentConfigById,
  hasEnabledKnowledge,
  hasKnowledge,
  hasSystemRole,
  inboxAgentConfig,
  inboxAgentModel,
  isAgentConfigLoading,
  isInboxSession,
  openingMessage,
  openingQuestions,
};
