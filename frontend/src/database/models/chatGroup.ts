import { nanoid } from 'nanoid';

import {
  DB,
  toNullString,
  toNullJSON,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
  boolToInt,
  intToBool,
} from '@/types/database';

export class ChatGroupModel {
  private userId: string;

  constructor(_db: any, userId: string) {
    this.userId = userId;
  }

  // ******* Query Methods ******* //

  async findById(id: string): Promise<any | undefined> {
    return await DB.GetChatGroup({
      id,
      userId: this.userId,
    });
  }

  async query(): Promise<any[]> {
    return await DB.ListChatGroups({
      userId: this.userId,
    });
  }

  /**
   * OPTIMIZED: Uses single JOIN query (3A) to fetch groups with agents
   * Much faster than N+1 approach
   */
  async queryWithMemberDetails(): Promise<any[]> {
    // Single query with JOINs
    const results = await DB.ListChatGroupsWithAgents({
      userId: this.userId,
    });

    // Group by group_id
    const groupsMap = new Map<string, any>();

    for (const row of results) {
      const groupId = row.groupId;

      if (!groupsMap.has(groupId)) {
        groupsMap.set(groupId, {
          id: groupId,
          title: getNullableString(row.groupTitle as any),
          description: getNullableString(row.groupDescription as any),
          config: parseNullableJSON(row.groupConfig as any),
          pinned: intToBool(row.groupPinned || 0),
          createdAt: row.groupCreatedAt,
          updatedAt: row.groupUpdatedAt,
          userId: this.userId,
          members: [],
        });
      }

      // Add agent if exists
      if (row.agentId) {
        const group = groupsMap.get(groupId);
        group.members.push({
          id: row.agentId,
          title: getNullableString(row.agentTitle as any),
          description: getNullableString(row.agentDescription as any),
          avatar: getNullableString(row.agentAvatar as any),
          backgroundColor: getNullableString(row.agentBgColor as any),
          chatConfig: parseNullableJSON(row.agentChatConfig as any),
          params: parseNullableJSON(row.agentParams as any),
          systemRole: getNullableString(row.agentSystemRole as any),
          tts: parseNullableJSON(row.agentTts as any),
          model: getNullableString(row.agentModel as any),
          provider: getNullableString(row.agentProvider as any),
          createdAt: row.agentCreatedAt,
          updatedAt: row.agentUpdatedAt,
        });
      }
    }

    return Array.from(groupsMap.values());
  }

  /**
   * OPTIMIZED: Uses JOIN query to fetch group with agents
   */
  async findGroupWithAgents(groupId: string): Promise<{
    agents: any[];
    group: any;
  } | null> {
    const results = await DB.GetChatGroupWithAgents({
      id: toNullString(groupId),
      userId: this.userId,
    });

    if (results.length === 0) return null;

    // First row contains group data
    const firstRow = results[0];
    const group = {
      id: firstRow.groupId,
      title: getNullableString(firstRow.groupTitle as any),
      description: getNullableString(firstRow.groupDescription as any),
      config: parseNullableJSON(firstRow.groupConfig as any),
      pinned: intToBool(firstRow.groupPinned || 0),
      createdAt: firstRow.groupCreatedAt,
      updatedAt: firstRow.groupUpdatedAt,
      userId: this.userId,
    };

    // Extract agents
    const agents = results
      .filter((row) => row.agentId)
      .map((row) => ({
        agentId: row.agentId,
        chatGroupId: groupId,
        order: row.agentSortOrder || 0,
        role: getNullableString(row.agentRole as any),
        enabled: intToBool(row.agentEnabled || 1),
        userId: this.userId,
      }));

    return { agents, group };
  }

  // ******* Create Methods ******* //

  async create(params: any): Promise<any> {
    const id = params.id || nanoid();
    const now = currentTimestampMs();

    const result = await DB.CreateChatGroup({
      id,
      title: toNullString(params.title),
      description: toNullString(params.description),
      config: toNullJSON(params.config),
      clientId: toNullString(params.clientId),
      userId: this.userId,
      pinned: boolToInt(params.pinned || false),
      createdAt: now,
      updatedAt: now,
    });

    return result;
  }

  async createWithAgents(
    groupParams: any,
    agentIds: string[],
  ): Promise<{ agents: any[]; group: any }> {
    const group = await this.create(groupParams);

    if (agentIds.length === 0) {
      return { agents: [], group };
    }

    const agents = [];
    const now = currentTimestampMs();

    for (let i = 0; i < agentIds.length; i++) {
      await DB.LinkChatGroupToAgent({
        chatGroupId: toNullString(group.id),
        agentId: toNullString(agentIds[i]),
        userId: this.userId,
        enabled: boolToInt(true),
        sortOrder: i,
        role: toNullString('assistant'),
        createdAt: now,
        updatedAt: now,
      });

      agents.push({
        agentId: agentIds[i],
        chatGroupId: group.id,
        order: i,
        role: 'assistant',
        userId: this.userId,
      });
    }

    return { agents, group };
  }

  // ******* Update Methods ******* //

  async update(id: string, value: any): Promise<any> {
    const result = await DB.UpdateChatGroup({
      id,
      userId: this.userId,
      title: toNullString(value.title),
      description: toNullString(value.description),
      config: toNullJSON(value.config),
      pinned: boolToInt(value.pinned ?? false),
      updatedAt: currentTimestampMs(),
    });

    if (!result) {
      throw new Error('Chat group not found or access denied');
    }

    return result;
  }

  async addAgentToGroup(
    groupId: string,
    agentId: string,
    options?: { order?: number; role?: string },
  ): Promise<any> {
    const now = currentTimestampMs();

    await DB.LinkChatGroupToAgent({
      chatGroupId: toNullString(groupId),
      agentId: toNullString(agentId),
      userId: this.userId,
      enabled: boolToInt(true),
      sortOrder: options?.order || 0,
      role: toNullString(options?.role || 'assistant'),
      createdAt: now,
      updatedAt: now,
    });

    return {
      agentId,
      chatGroupId: groupId,
      order: options?.order || 0,
      role: options?.role || 'assistant',
      userId: this.userId,
    };
  }

  async addAgentsToGroup(groupId: string, agentIds: string[]): Promise<any[]> {
    const group = await this.findById(groupId);
    if (!group) throw new Error('Group not found');

    const existingAgents = await this.getGroupAgents(groupId);
    const existingAgentIds = new Set(existingAgents.map((a: any) => a.agentId));

    const newAgentIds = agentIds.filter((id) => !existingAgentIds.has(id));

    if (newAgentIds.length === 0) {
      return [];
    }

    const newAgents = [];
    const now = currentTimestampMs();

    for (const agentId of newAgentIds) {
      await DB.LinkChatGroupToAgent({
        chatGroupId: toNullString(groupId),
        agentId: toNullString(agentId),
        userId: this.userId,
        enabled: boolToInt(true),
        sortOrder: 0,
        role: toNullString('assistant'),
        createdAt: now,
        updatedAt: now,
      });

      newAgents.push({
        agentId,
        chatGroupId: groupId,
        enabled: true,
        userId: this.userId,
      });
    }

    return newAgents;
  }

  async removeAgentFromGroup(groupId: string, agentId: string): Promise<void> {
    await DB.UnlinkChatGroupFromAgent({
      chatGroupId: toNullString(groupId),
      agentId: toNullString(agentId),
      userId: this.userId,
    });
  }

  async updateAgentInGroup(
    groupId: string,
    agentId: string,
    updates: Partial<any>,
  ): Promise<any> {
    const result = await DB.UpdateChatGroupAgentLink({
      chatGroupId: toNullString(groupId),
      agentId: toNullString(agentId),
      userId: this.userId,
      sortOrder: updates.order ?? 0,
      role: toNullString(updates.role || 'assistant'),
      enabled: boolToInt(updates.enabled ?? true),
      updatedAt: currentTimestampMs(),
    });

    return result;
  }

  // ******* Delete Methods ******* //

  async delete(id: string): Promise<any> {
    // Get group first to return it
    const group = await this.findById(id);
    if (!group) {
      throw new Error('Chat group not found or access denied');
    }

    // Delete (agents are automatically deleted due to CASCADE)
    await DB.DeleteChatGroup({
      id,
      userId: this.userId,
    });

    return group;
  }

  async deleteAll(): Promise<void> {
    await DB.DeleteAllChatGroups({
      userId: this.userId,
    });
  }

  // ******* Agent Query Methods ******* //

  async getGroupAgents(groupId: string): Promise<any[]> {
    return await DB.GetChatGroupAgentLinks({
      chatGroupId: toNullString(groupId),
      userId: this.userId,
    });
  }

  async getEnabledGroupAgents(groupId: string): Promise<any[]> {
    return await DB.GetEnabledChatGroupAgentLinks({
      chatGroupId: toNullString(groupId),
      userId: this.userId,
    });
  }

  async getGroupsWithAgents(agentIds?: string[]): Promise<any[]> {
    if (!agentIds || agentIds.length === 0) {
      return this.query();
    }

    // Get all groups, then filter by agents
    // Note: This is not fully optimized - would need specific query for this
    const allGroups = await this.queryWithMemberDetails();
    
    return allGroups.filter((group) =>
      group.members.some((member: any) => agentIds.includes(member.id)),
    );
  }
}

