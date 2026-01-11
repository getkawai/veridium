import { ChatMessageError } from "@/types";

/* eslint-disable typescript-sort-keys/interface */

export interface IMessageService {
  updateMessageError(id: string, error: ChatMessageError): Promise<any>;
  updateMessageRAG(id: string, value: any): Promise<any>;
  removeMessagesByAssistant(
    assistantId: string,
    topicId?: string,
  ): Promise<any>;
}
