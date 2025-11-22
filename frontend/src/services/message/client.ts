
import { INBOX_SESSION_ID } from '@/const/session';
import { clientDB } from '@/database/client/db';
import { MessageModel } from '@/database/models/message';
import { BaseClientService } from '@/services/baseClientService';

import { IMessageService } from './type';

export class ClientService extends BaseClientService implements IMessageService {
  private get messageModel(): MessageModel {
    return new MessageModel(clientDB as any, this.userId);
  }

  updateMessageError: IMessageService['updateMessageError'] = async (id, error) => {
    return this.messageModel.update(id, { error });
  };


  updateMessagePluginError: IMessageService['updateMessagePluginError'] = async (id, value) => {
    return this.messageModel.updateMessagePlugin(id, { error: value });
  };

  updateMessageRAG: IMessageService['updateMessageRAG'] = async (id, value) => {
    console.log(id, value);
    throw new Error('not implemented');
  };

  removeMessagesByAssistant: IMessageService['removeMessagesByAssistant'] = async (
    sessionId,
    topicId,
  ) => {
    return this.messageModel.deleteMessagesBySession(this.toDbSessionId(sessionId), topicId);
  };


  private toDbSessionId = (sessionId: string | undefined) => {
    return sessionId === INBOX_SESSION_ID ? undefined : sessionId;
  };
}
