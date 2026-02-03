/* eslint-disable sort-keys-fix/sort-keys-fix, typescript-sort-keys/interface */
// import { ToolNameResolver } from '@/context-engine';
import {
  ChatMessageError,
  ChatToolPayload,
  CreateMessageParams,
  MessageToolCall,
  ToolsCallingContext,
} from '@/types';
import { LobeChatPluginManifest } from '@/chat-plugin-sdk';
import isEqual from 'fast-deep-equal';
import { StateCreator } from 'zustand/vanilla';

// import { mcpService } from '@/services/mcp';
import { ChatStore } from '@/store/chat/store';
import { useToolStore } from '@/store/tool';
import { pluginSelectors } from '@/store/tool/selectors';
import { builtinTools } from '@/tools';
import { merge } from '@/utils/merge';
import { safeParseJSON } from '@/utils/safeParseJSON';
import { setNamespace } from '@/utils/storeDebug';

import { chatSelectors } from '../message/selectors';


const n = setNamespace('plugin');

export interface ChatPluginAction {
  createAssistantMessageByPlugin: (content: string, parentId: string) => Promise<void>;
  fillPluginMessageContent: (
    id: string,
    content: string,
    triggerAiMessage?: boolean,
  ) => Promise<void>;

  invokeBuiltinTool: (id: string, payload: ChatToolPayload) => Promise<void>;
  invokeDefaultTypePlugin: (id: string, payload: any) => Promise<string | undefined>;
  invokeMarkdownTypePlugin: (id: string, payload: ChatToolPayload) => Promise<void>;
  invokeMCPTypePlugin: (id: string, payload: ChatToolPayload) => Promise<string | undefined>;

  invokeStandaloneTypePlugin: (id: string, payload: ChatToolPayload) => Promise<void>;

  reInvokeToolMessage: (id: string) => Promise<void>;
  triggerAIMessage: (params: {
    parentId?: string;
    traceId?: string;
    threadId?: string;
    inPortalThread?: boolean;
    inSearchWorkflow?: boolean;
  }) => Promise<void>;
  summaryPluginContent: (id: string) => Promise<void>;

  /**
   * @deprecated V1 method
   */
  triggerToolCalls: (
    id: string,
    params?: { threadId?: string; inPortalThread?: boolean; inSearchWorkflow?: boolean },
  ) => Promise<void>;
  updatePluginState: (id: string, value: any) => Promise<void>;
  updatePluginArguments: <T = any>(id: string, value: T, replace?: boolean) => Promise<void>;

  internal_addToolToAssistantMessage: (id: string, tool: ChatToolPayload) => Promise<void>;
  internal_removeToolToAssistantMessage: (id: string, tool_call_id?: string) => Promise<void>;
  /**
   * use the optimistic update value to update the message tools to database
   */
  internal_refreshToUpdateMessageTools: (id: string) => Promise<void>;

  internal_callPluginApi: (id: string, payload: ChatToolPayload) => Promise<string | undefined>;
  internal_invokeDifferentTypePlugin: (id: string, payload: ChatToolPayload) => Promise<any>;
  internal_togglePluginApiCalling: (
    loading: boolean,
    id?: string,
    action?: string,
  ) => AbortController | undefined;
  internal_transformToolCalls: (toolCalls: MessageToolCall[]) => ChatToolPayload[];
  internal_updatePluginError: (id: string, error: ChatMessageError) => Promise<void>;
  internal_constructToolsCallingContext: (id: string) => ToolsCallingContext | undefined;
}

export const chatPlugin: StateCreator<
  ChatStore,
  [['zustand/devtools', never]],
  [],
  ChatPluginAction
> = (set, get) => ({
  createAssistantMessageByPlugin: async (content, parentId) => {
    // Dummy implementation for UI focus
    console.log('Creating assistant message by plugin:', { content, parentId });

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 300));

    // Mock message creation success
    console.log('Assistant message created successfully');
  },

  fillPluginMessageContent: async (id, content, triggerAiMessage) => {
    const { triggerAIMessage, internal_updateMessageContent } = get();

    await internal_updateMessageContent(id, content);

    if (triggerAiMessage) await triggerAIMessage({ parentId: id });
  },
  invokeBuiltinTool: async (id, payload) => {
    const {
      internal_togglePluginApiCalling,
      internal_updateMessageContent,
    } = get();

    // Dummy implementation for UI focus
    console.log('Invoking builtin tool:', { id, payload });

    const params = JSON.parse(payload.arguments);
    internal_togglePluginApiCalling(true, id, n('invokeBuiltinTool/start') as string);

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 800));

    // Mock transformed data
    const mockData = JSON.stringify({
      result: `Mock execution result for ${payload.apiName}`,
      params,
      timestamp: Date.now(),
    });

    internal_togglePluginApiCalling(false, id, n('invokeBuiltinTool/end') as string);

    await internal_updateMessageContent(id, mockData);

    // Mock tool API call
    const mockContent = {
      success: true,
      data: `Mock ${payload.apiName} execution completed`,
    };

    // @ts-ignore
    const { [payload.apiName]: action } = get();
    if (!action) return;

    return await action(id, mockContent);
  },

  invokeDefaultTypePlugin: async (id, payload) => {
    const { internal_callPluginApi } = get();

    const data = await internal_callPluginApi(id, payload);

    if (!data) return;

    return data;
  },

  invokeMarkdownTypePlugin: async (id, payload) => {
    const { internal_callPluginApi } = get();

    await internal_callPluginApi(id, payload);
  },

  invokeStandaloneTypePlugin: async (id, payload) => {
    // Dummy implementation for UI focus
    console.log('Invoking standalone plugin:', { id, payload });

    // Mock validation result
    const mockResult = { valid: true, errors: null };

    if (!mockResult.valid) {
      // Mock error handling
      console.error('Plugin settings invalid:', mockResult.errors);
      return;
    }

    // Simulate successful plugin execution
    await new Promise(resolve => setTimeout(resolve, 500));
    console.log('Standalone plugin executed successfully');
  },

  reInvokeToolMessage: async (id) => {
    const message = chatSelectors.getMessageById(id)(get());
    if (!message || message.role !== 'tool' || !message.plugin) return;

    // if there is error content, then clear the error
    if (!!message.pluginError) {
      get().internal_updateMessagePluginError(id, null);
    }

    const payload: ChatToolPayload = { ...message.plugin, id: message.tool_call_id! };

    await get().internal_invokeDifferentTypePlugin(id, payload);
  },

  triggerAIMessage: async ({ parentId, traceId, threadId, inPortalThread, inSearchWorkflow }) => {
    // TODO: implement this
    // const { activeId, activeTopicId } = get();

    // // Use backendAgentChat to trigger AI response
    // // We send an empty message or just the context to trigger the agent
    // // In the new backend architecture, sending a request with the current session/thread context
    // // should be enough for the agent to pick up the conversation (including recent tool outputs).

    // try {
    //   await backendAgentChat.sendMessage({
    //     session_id: activeId,
    //     topic_id: activeTopicId,
    //     thread_id: threadId,
    //     // We don't send a user message here, just triggering the agent
    //     // The backend should handle "continue" or "reply to tool output" logic
    //     message: undefined,
    //   });

    //   // Refresh messages to show the new AI response
    //   await get().refreshMessages();

    // } catch (error) {
    //   console.error('[Plugin] triggerAIMessage failed:', error);
    // }
  },

  summaryPluginContent: async (id) => {
    // const message = chatSelectors.getMessageById(id)(get());
    // if (!message || message.role !== 'tool') return;

    // // Use backendLibraryChat for stateless summarization
    // const summaryPrompt = `Please summarize the following content:\n\n${message.content}`;

    // try {
    //   // 1. Generate summary using stateless chat completion
    //   const response = await backendLibraryChat.chatCompletion({
    //     messages: [
    //       { role: 'user', content: summaryPrompt }
    //     ],
    //     temperature: 0.5,
    //   });

    //   const summary = response.choices?.[0]?.message?.content;

    //   if (!summary) return;

    //   // 2. Add the summary as an assistant message to the UI
    //   // We use internal_createMessage to add it to the store and DB
    //   // This makes it look like the agent responded, but without the overhead of the full agent loop
    //   const { activeId, activeTopicId, activeThreadId } = get();

    //   await get().internal_createMessage({
    //     role: 'assistant',
    //     content: summary,
    //     sessionId: activeId,
    //     topicId: activeTopicId,
    //     threadId: activeThreadId,
    //     parentId: id, // Link to the tool message
    //   });

    // } catch (error) {
    //   console.error('[Plugin] summaryPluginContent failed:', error);
    // }
  },

  triggerToolCalls: async (assistantId, { threadId, inPortalThread, inSearchWorkflow } = {}) => {
    const message = chatSelectors.getMessageById(assistantId)(get());
    if (!message || !message.tools) return;

    let shouldCreateMessage = false;
    let latestToolId = '';
    const messagePools = message.tools.map(async (payload) => {
      const toolMessage: CreateMessageParams = {
        content: '',
        parentId: assistantId,
        plugin: payload,
        role: 'tool',
        sessionId: get().activeId,
        tool_call_id: payload.id,
        threadId,
        topicId: get().activeTopicId, // if there is activeTopicId，then add it to topicId
        groupId: message.groupId, // Propagate groupId from parent message for group chat
      };

      const id = await get().internal_createMessage(toolMessage);
      if (!id) return;

      // trigger the plugin call
      const data = await get().internal_invokeDifferentTypePlugin(id, payload);

      if (data && !['markdown', 'standalone'].includes(payload.type)) {
        shouldCreateMessage = true;
        latestToolId = id;
      }
    });

    await Promise.all(messagePools);

    await get().internal_toggleMessageInToolsCalling(false, assistantId);

    // only default type tool calls should trigger AI message
    if (!shouldCreateMessage) return;

    const traceId = chatSelectors.getTraceIdByMessageId(latestToolId)(get());

    await get().triggerAIMessage({ traceId, threadId, inPortalThread, inSearchWorkflow });
  },
  updatePluginState: async (id, value) => {
    // Dummy implementation for UI focus
    console.log('Updating plugin state:', { id, value });

    // optimistic update
    get().internal_dispatchMessage({ id, type: 'updateMessage', value: { pluginState: value } });

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 200));
    console.log('Plugin state updated successfully');
  },

  updatePluginArguments: async (id, value, replace = false) => {
    // Dummy implementation for UI focus
    console.log('Updating plugin arguments:', { id, value, replace });

    const toolMessage = chatSelectors.getMessageById(id)(get());
    if (!toolMessage || !toolMessage?.tool_call_id) return;

    const prevArguments = toolMessage?.plugin?.arguments;
    const prevJson = safeParseJSON(prevArguments || '');
    const nextValue = replace ? (value as any) : merge(prevJson || {}, value);
    if (isEqual(prevJson, nextValue)) return;

    // optimistic update
    get().internal_dispatchMessage({
      id,
      type: 'updateMessagePlugin',
      value: { arguments: JSON.stringify(nextValue) },
    });

    // Mock assistant message update
    const assistantMessage = chatSelectors.getMessageById(toolMessage?.parentId || '')(get());
    if (assistantMessage) {
      get().internal_dispatchMessage({
        id: assistantMessage.id,
        type: 'updateMessageTools',
        tool_call_id: toolMessage?.tool_call_id,
        value: { arguments: JSON.stringify(nextValue) },
      });
    }

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 300));
    console.log('Plugin arguments updated successfully');
  },

  internal_addToolToAssistantMessage: async (id, tool) => {
    const assistantMessage = chatSelectors.getMessageById(id)(get());
    if (!assistantMessage) return;

    const { internal_dispatchMessage, internal_refreshToUpdateMessageTools } = get();
    internal_dispatchMessage({
      type: 'addMessageTool',
      value: tool,
      id: assistantMessage.id,
    });

    await internal_refreshToUpdateMessageTools(id);
  },

  internal_removeToolToAssistantMessage: async (id, tool_call_id) => {
    const message = chatSelectors.getMessageById(id)(get());
    if (!message || !tool_call_id) return;

    const { internal_dispatchMessage, internal_refreshToUpdateMessageTools } = get();

    // optimistic update
    internal_dispatchMessage({ type: 'deleteMessageTool', tool_call_id, id: message.id });

    // update the message tools
    await internal_refreshToUpdateMessageTools(id);
  },
  internal_refreshToUpdateMessageTools: async (id) => {
    // Dummy implementation for UI focus
    console.log('Refreshing message tools:', id);

    const message = chatSelectors.getMessageById(id)(get());
    if (!message || !message.tools) return;

    const { internal_toggleMessageLoading } = get();

    internal_toggleMessageLoading(true, id);

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 400));

    internal_toggleMessageLoading(false, id);
    console.log('Message tools refreshed successfully');
  },

  internal_callPluginApi: async (id, payload) => {
    const { internal_updateMessageContent, internal_togglePluginApiCalling } = get();

    // Dummy implementation for UI focus
    console.log('Calling plugin API:', { id, payload });

    const abortController = internal_togglePluginApiCalling(
      true,
      id,
      n('fetchPlugin/start') as string,
    );

    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 1200));

    // Mock API response
    const mockResponse = {
      text: `Mock plugin response for ${payload.apiName}: ${JSON.stringify(payload.arguments)}`,
      traceId: `mock-trace-${Date.now()}`,
    };

    internal_togglePluginApiCalling(false, id, n('fetchPlugin/end') as string);

    await internal_updateMessageContent(id, mockResponse.text);

    return mockResponse.text;
  },

  internal_invokeDifferentTypePlugin: async (id, payload) => {
    switch (payload.type) {
      case 'standalone': {
        return await get().invokeStandaloneTypePlugin(id, payload);
      }

      case 'markdown': {
        return await get().invokeMarkdownTypePlugin(id, payload);
      }

      case 'builtin': {
        return await get().invokeBuiltinTool(id, payload);
      }

      // @ts-ignore
      case 'mcp': {
        return await get().invokeMCPTypePlugin(id, payload);
      }

      default: {
        return await get().invokeDefaultTypePlugin(id, payload);
      }
    }
  },
  invokeMCPTypePlugin: async (id, payload) => {
    const { internal_updateMessageContent, internal_togglePluginApiCalling } = get();

    // Dummy implementation for UI focus
    console.log('Invoking MCP plugin:', { id, payload });

    let data: string = '';

    const abortController = internal_togglePluginApiCalling(
      true,
      id,
      n('fetchPlugin/start') as string,
    );

    // Simulate MCP call delay
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Mock MCP response
    data = `Mock MCP response for ${payload.apiName}: ${JSON.stringify(payload.arguments)}`;

    internal_togglePluginApiCalling(false, id, n('fetchPlugin/end') as string);

    await internal_updateMessageContent(id, data);

    return data;
  },

  internal_togglePluginApiCalling: (loading, id, action) => {
    return get().internal_toggleLoadingArrays('pluginApiLoadingIds', loading, id, action);
  },

  internal_transformToolCalls: (toolCalls) => {
    // const toolNameResolver = new ToolNameResolver();
    // Temporary simple tool name resolver
    const toolNameResolver = {
      resolve: (calls: any[], manifests: Record<string, LobeChatPluginManifest>) => {
        return calls.map((call) => {
          const fnName = call.function?.name || '';
          const [identifier, apiName] = fnName.split('____');
          const manifest = manifests[identifier];
          const api = manifest?.api?.find((a) => a.name === apiName);
          return {
            apiName: apiName || fnName,
            arguments: call.function?.arguments || '{}',
            id: call.id,
            identifier: identifier || fnName,
            type: api ? (manifest?.type || 'default') : 'default',
          };
        });
      },
    };

    // Build manifests map from tool store
    const toolStoreState = useToolStore.getState();
    const manifests: Record<string, LobeChatPluginManifest> = {};

    // Get all installed plugins
    const installedPlugins = pluginSelectors.installedPlugins(toolStoreState);
    for (const plugin of installedPlugins) {
      if (plugin.manifest) {
        manifests[plugin.identifier] = plugin.manifest as LobeChatPluginManifest;
      }
    }

    // Get all builtin tools
    for (const tool of builtinTools) {
      if (tool.manifest) {
        manifests[tool.identifier] = tool.manifest as LobeChatPluginManifest;
      }
    }

    return toolNameResolver.resolve(toolCalls, manifests);
  },
  internal_updatePluginError: async (id, error) => {
    // Dummy implementation for UI focus
    console.log('Updating plugin error:', { id, error });

    get().internal_dispatchMessage({ id, type: 'updateMessage', value: { error } });

    // Simulate async operation
    await new Promise(resolve => setTimeout(resolve, 200));
    console.log('Plugin error updated successfully');
  },

  internal_constructToolsCallingContext: (id: string) => {
    const message = chatSelectors.getMessageById(id)(get());
    if (!message) return;

    return {
      topicId: message.topicId,
    };
  },
});
