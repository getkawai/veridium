import { SendMessageServerParams, SendMessageServerResponse, StructureOutputParams } from '@/types';

class AiChatService {
  sendMessageInServer = async (
    params: SendMessageServerParams,
    abortController: AbortController,
  ): Promise<SendMessageServerResponse> => {
    // Mock delay to simulate network request
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Check if operation was aborted
    if (abortController?.signal.aborted) {
      throw new Error('Operation aborted');
    }

    // Generate mock IDs
    const userMessageId = `msg-${Date.now()}-user`;
    const assistantMessageId = `msg-${Date.now()}-assistant`;
    const topicId = params.topicId || `topic-${Date.now()}`;

    // Mock user message
    const userMessage = {
      id: userMessageId,
      content: params.newUserMessage.content,
      role: 'user' as const,
      createdAt: Date.now(),
      updatedAt: Date.now(),
      meta: {
        avatar: 'user',
        title: 'User',
      },
      topicId,
      sessionId: params.sessionId,
      threadId: params.threadId,
    };

    // Mock assistant response
    const assistantMessage = {
      id: assistantMessageId,
      content: `Thank you for your message: "${params.newUserMessage.content}". This is a mock response from the AI assistant using the ${params.newAssistantMessage.model} model.`,
      role: 'assistant' as const,
      createdAt: Date.now(),
      updatedAt: Date.now(),
      meta: {
        avatar: 'assistant',
        title: 'AI Assistant',
      },
      topicId,
      sessionId: params.sessionId,
      threadId: params.threadId,
      usage: {
        completionTokens: 150,
        promptTokens: 50,
        totalTokens: 200,
      },
      performance: {
        tps: 25,
        ttft: 800,
        duration: 1200,
      },
    };

    return {
      userMessageId,
      assistantMessageId,
      topicId,
      isCreateNewTopic: !params.topicId,
      messages: [userMessage, assistantMessage],
      topics: params.topicId ? undefined : [{
        id: topicId,
        title: params.newTopic?.title || 'New Conversation',
        createdAt: Date.now(),
        updatedAt: Date.now(),
        favorite: false,
        sessionId: params.sessionId,
      }],
    };
  };

  generateJSON = async (
    params: Omit<StructureOutputParams, 'keyVaultsPayload'>,
    abortController: AbortController,
  ) => {
    // Mock delay to simulate network request
    await new Promise(resolve => setTimeout(resolve, 800));

    // Check if operation was aborted
    if (abortController?.signal.aborted) {
      throw new Error('Operation aborted');
    }

    // Mock structured output based on schema
    let mockData: any = {};

    if (params.schema?.schema?.properties) {
      const properties = params.schema.schema.properties;

      // Generate mock data for each property
      Object.keys(properties).forEach(key => {
        const prop = properties[key];
        if (prop.type === 'string') {
          mockData[key] = `Mock ${key} value`;
        } else if (prop.type === 'number') {
          mockData[key] = Math.floor(Math.random() * 100);
        } else if (prop.type === 'boolean') {
          mockData[key] = Math.random() > 0.5;
        } else if (prop.type === 'array') {
          mockData[key] = [`mock item 1`, `mock item 2`];
        } else {
          mockData[key] = `Mock ${key} data`;
        }
      });
    } else {
      // Default mock structure
      mockData = {
        name: 'Mock Generated Object',
        description: 'This is a mock response from the structured output service',
        value: Math.floor(Math.random() * 1000),
        status: 'completed',
        timestamp: new Date().toISOString(),
      };
    }

    return {
      data: mockData,
      usage: {
        completionTokens: 120,
        promptTokens: 80,
        totalTokens: 200,
      },
      performance: {
        tps: 30,
        ttft: 600,
        duration: 1000,
      },
    };
  };

  // sendGroupMessageInServer = async (params: SendMessageServerParams) => {
  //   return lambdaClient.aiChat.sendGroupMessageInServer.mutate(cleanObject(params));
  // };
}

export const aiChatService = new AiChatService();
