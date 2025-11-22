
/* eslint-disable typescript-sort-keys/interface */

export interface IMessageService {
  removeMessagesByAssistant(assistantId: string, topicId?: string): Promise<any>;
}
