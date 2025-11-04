import {
  DB,
  toNullString,
  currentTimestampMs,
  boolToInt,
  intToBool,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';

export class AgentModel {
  private userId: string;
  private logger = createModelLogger('Agent', 'AgentModel', 'database/models/agent');

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  getAgentConfigById = async (id: string) => {
    await this.logger.methodEntry('getAgentConfigById', { userId: this.userId, id });
    const agent = await DB.GetAgent({
      id,
      userId: this.userId,
    });

    const knowledge = await this.getAgentAssignedKnowledge(id);

    const result = { ...agent, ...knowledge };
    await this.logger.methodExit('getAgentConfigById', result);
    return result;
  };

  getAgentAssignedKnowledge = async (id: string) => {
    const knowledgeBaseResult = await DB.GetAgentKnowledgeBases({
      agentId: toNullString(id) as any,
      userId: this.userId,
    });

    const fileResult = await DB.GetAgentFilesWithEnabled({
      agentId: toNullString(id) as any,
      userId: this.userId,
    });

    return {
      files: fileResult.map((item) => ({
        ...item,
        enabled: intToBool(item.enabled),
      })),
      knowledgeBases: knowledgeBaseResult.map((item) => ({
        ...item,
        enabled: intToBool(item.enabled),
      })),
    };
  };

  /**
   * Find agent by session id
   */
  findBySessionId = async (sessionId: string) => {
    try {
      const agent = await DB.GetAgentBySessionId({
        sessionId: toNullString(sessionId) as any,
        userId: this.userId,
      });

      return await this.getAgentConfigById(agent.id);
    } catch {
      return undefined;
    }
  };

  createAgentKnowledgeBase = async (
    agentId: string,
    knowledgeBaseId: string,
    enabled: boolean = true,
  ) => {
    const now = currentTimestampMs();

    await DB.LinkAgentToKnowledgeBase({
      agentId: toNullString(agentId) as any,
      knowledgeBaseId: toNullString(knowledgeBaseId) as any,
      enabled: boolToInt(enabled),
      userId: this.userId,
      createdAt: now,
      updatedAt: now,
    });
  };

  deleteAgentKnowledgeBase = async (agentId: string, knowledgeBaseId: string) => {
    await DB.UnlinkAgentFromKnowledgeBase({
      agentId: toNullString(agentId) as any,
      knowledgeBaseId: toNullString(knowledgeBaseId) as any,
      userId: this.userId,
    });
  };

  toggleKnowledgeBase = async (agentId: string, knowledgeBaseId: string, enabled?: boolean) => {
    await DB.ToggleAgentKnowledgeBase({
      agentId: toNullString(agentId) as any,
      knowledgeBaseId: toNullString(knowledgeBaseId) as any,
      enabled: boolToInt(enabled || false),
      userId: this.userId,
    });
  };

  createAgentFiles = async (agentId: string, fileIds: string[], enabled: boolean = true) => {
    // Exclude the fileIds that already exist in agentsFiles, and then insert them
    const existingFiles = await DB.GetAgentFileIds({
      agentId: toNullString(agentId) as any,
      userId: this.userId,
      fileIds,
    });

    const existingFilesIds = new Set(existingFiles.map((item) => (item as any).fileId));

    const needToInsertFileIds = fileIds.filter((fileId) => !existingFilesIds.has(fileId));

    if (needToInsertFileIds.length === 0) return;

    const now = currentTimestampMs();

    // Note: No batch insert support - insert one by one
    await Promise.all(
      needToInsertFileIds.map((fileId) =>
        DB.BatchLinkAgentToFiles({
          agentId: toNullString(agentId) as any,
          fileId: toNullString(fileId) as any,
          enabled: boolToInt(enabled),
          userId: this.userId,
          createdAt: now,
          updatedAt: now,
        }),
      ),
    );
  };

  deleteAgentFile = async (agentId: string, fileId: string) => {
    await DB.UnlinkAgentFromFile({
      agentId: toNullString(agentId) as any,
      fileId: toNullString(fileId) as any,
      userId: this.userId,
    });
  };

  toggleFile = async (agentId: string, fileId: string, enabled?: boolean) => {
    await DB.ToggleAgentFile({
      agentId: toNullString(agentId) as any,
      fileId: toNullString(fileId) as any,
      enabled: boolToInt(enabled || false),
      userId: this.userId,
    });
  };
}

