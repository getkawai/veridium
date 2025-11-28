/**
 * Tools Engineering - Unified tools processing
 * Context-engine removed - tools are now handled by backend
 */
// import { ToolsEngine } from '@/context-engine';
// import type { PluginEnableChecker } from '@/context-engine';
import { ChatCompletionTool, WorkingModel } from '@/types';
import { LobeChatPluginManifest } from '@/chat-plugin-sdk';

import { getToolStoreState } from '@/store/tool';
import { pluginSelectors } from '@/store/tool/selectors';
import { WebBrowsingManifest } from '@/tools/web-browsing';

import { getSearchConfig } from '../getSearchConfig';
import { isCanUseFC } from '../isCanUseFC';
import { shouldEnableTool } from '../toolFilters';

/**
 * Plugin enable checker type
 */
export type PluginEnableChecker = (context: { pluginId: string }) => boolean;

/**
 * Tools engine configuration options
 */
export interface ToolsEngineConfig {
  /** Additional manifests to include beyond the standard ones */
  additionalManifests?: LobeChatPluginManifest[];
  /** Default tool IDs that will always be added to the end of the tools list */
  defaultToolIds?: string[];
  /** Custom enable checker for plugins */
  enableChecker?: PluginEnableChecker;
}

/**
 * Simple tools engine implementation (context-engine removed)
 */
class SimpleToolsEngine {
  private manifests: LobeChatPluginManifest[];
  private enableChecker?: PluginEnableChecker;
  private functionCallChecker: (model: string, provider: string) => boolean;
  private defaultToolIds?: string[];

  constructor(config: {
    manifestSchemas: LobeChatPluginManifest[];
    enableChecker?: PluginEnableChecker;
    functionCallChecker: (model: string, provider: string) => boolean;
    defaultToolIds?: string[];
  }) {
    this.manifests = config.manifestSchemas;
    this.enableChecker = config.enableChecker;
    this.functionCallChecker = config.functionCallChecker;
    this.defaultToolIds = config.defaultToolIds;
  }

  generateTools(params: {
    model: string;
    provider: string;
    toolIds: string[];
  }): ChatCompletionTool[] {
    const { model, provider, toolIds } = params;

    if (!this.functionCallChecker(model, provider)) {
      return [];
    }

    const allToolIds = [...toolIds, ...(this.defaultToolIds || [])];
    const tools: ChatCompletionTool[] = [];

    for (const manifest of this.manifests) {
      if (!allToolIds.includes(manifest.identifier)) continue;
      if (this.enableChecker && !this.enableChecker({ pluginId: manifest.identifier })) continue;

      for (const api of manifest.api || []) {
        tools.push({
          type: 'function',
          function: {
            name: `${manifest.identifier}____${api.name}`,
            description: api.description || '',
            parameters: api.parameters as any,
          },
        });
      }
    }

    return tools;
  }
}

/**
 * Initialize ToolsEngine with current manifest schemas and configurable options
 */
export const createToolsEngine = (config: ToolsEngineConfig = {}): SimpleToolsEngine => {
  const { enableChecker, additionalManifests = [], defaultToolIds } = config;

  const toolStoreState = getToolStoreState();

  // Get all available plugin manifests
  const pluginManifests = pluginSelectors.installedPluginManifestList(toolStoreState);

  // Get all builtin tool manifests
  const builtinManifests = toolStoreState.builtinTools.map(
    (tool) => tool.manifest as LobeChatPluginManifest,
  );

  // Combine all manifests
  const allManifests = [...pluginManifests, ...builtinManifests, ...additionalManifests];

  return new SimpleToolsEngine({
    defaultToolIds,
    enableChecker,
    functionCallChecker: isCanUseFC,
    manifestSchemas: allManifests,
  });
};

export const createChatToolsEngine = (workingModel: WorkingModel) =>
  createToolsEngine({
    // Add WebBrowsingManifest as default tool
    defaultToolIds: [WebBrowsingManifest.identifier],
    // Create search-aware enableChecker for this request
    enableChecker: ({ pluginId }) => {
      // Check platform-specific constraints (e.g., LocalSystem desktop-only)
      if (!shouldEnableTool(pluginId)) {
        return false;
      }

      // For WebBrowsingManifest, apply search logic
      if (pluginId === WebBrowsingManifest.identifier) {
        const searchConfig = getSearchConfig(workingModel.model, workingModel.provider);
        return searchConfig.useApplicationBuiltinSearchTool;
      }

      // For all other plugins, enable by default
      return true;
    },
  });

/**
 * Provides the same functionality using ToolsEngine with enhanced capabilities
 *
 * @param toolIds - Array of tool IDs to generate tools for
 * @param model - Model name for function calling compatibility check (optional)
 * @param provider - Provider name for function calling compatibility check (optional)
 * @returns Array of ChatCompletionTool objects
 */
export const getEnabledTools = (
  toolIds: string[] = [],
  model: string,
  provider: string,
): ChatCompletionTool[] => {
  const toolsEngine = createToolsEngine();

  return (
    toolsEngine.generateTools({
      model: model, // Use provided model or fallback
      provider: provider, // Use provided provider or fallback
      toolIds,
    }) || []
  );
};
