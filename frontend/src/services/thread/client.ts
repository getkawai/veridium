import { INBOX_SESSION_ID } from '@/const/session';
import { clientDB } from '@/database/client/db';
import { MessageModel } from '@/database/models/message';
import { ThreadModel } from '@/database/models/thread';
import { BaseClientService } from '@/services/baseClientService';

import { IThreadService } from './type';

export class ClientService extends BaseClientService implements IThreadService {
  private get threadModel(): ThreadModel {
    return new ThreadModel(clientDB as any, this.userId);
  }

  private get messageModel(): MessageModel {
    return new MessageModel(clientDB as any, this.userId);
  }

  getThreads: IThreadService['getThreads'] = async (topicId) => {
    return this.threadModel.queryByTopicId(topicId);
  };

  createThreadWithMessage: IThreadService['createThreadWithMessage'] = async (input) => {
    let thread;
    try {
      thread = await this.threadModel.create({
        parentThreadId: input.parentThreadId,
        sourceMessageId: input.sourceMessageId,
        title: input.message.content.slice(0, 20),
        topicId: input.topicId,
        type: input.type,
      });
    } catch (error) {
      // Thread creation threw an error (non-conflict error)
      console.error('[createThreadWithMessage] Thread creation threw error:', error);
      throw error; // Re-throw to let caller handle it
    }

    // If thread creation failed (conflict - returns undefined), don't create the message
    if (!thread?.id) {
      console.error('[createThreadWithMessage] Thread creation failed (conflict), aborting message creation');
      return { messageId: undefined as any, threadId: undefined as any };
    }

    const message = await this.messageModel.create({
      ...input.message,
      sessionId: this.toDbSessionId(input.message.sessionId) as string,
      threadId: thread.id,
    });

    // If message creation failed, we still return the threadId but no messageId
    // The caller should handle this case
    if (!message?.id) {
      console.error('[createThreadWithMessage] Message creation failed after thread creation');
    }

    return { messageId: message?.id || (undefined as any), threadId: thread.id };
  };

  updateThread: IThreadService['updateThread'] = async (id, data) => {
    return this.threadModel.update(id, data);
  };

  removeThread: IThreadService['removeThread'] = async (id) => {
    return this.threadModel.delete(id);
  };

  private toDbSessionId = (sessionId: string | undefined) => {
    return sessionId === INBOX_SESSION_ID ? null : sessionId;
  };
}
