import {
  ChatMessageError,
  ChatMessagePluginError,
  UpdateMessageRAGParams,
} from '@/types';

/* eslint-disable typescript-sort-keys/interface */

export interface IMessageService {
  updateMessageError(id: string, error: ChatMessageError): Promise<any>;
  updateMessagePluginError(id: string, value: ChatMessagePluginError | null): Promise<any>;
  updateMessageRAG(id: string, value: UpdateMessageRAGParams): Promise<void>;
  removeMessagesByAssistant(assistantId: string, topicId?: string): Promise<any>;
}
