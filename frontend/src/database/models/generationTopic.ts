import { GenerationAsset, ImageGenerationTopic } from '@/types';
import { nanoid } from 'nanoid';
import {
  DB,
  toNullString,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
  GenerationTopic,
} from '@/types/database';
import { createModelLogger } from '@/utils/logger';
import { NotificationService, NotificationOptions } from '@@/github.com/wailsapp/wails/v3/pkg/services/notifications';
import { GetFullFileUrl } from '@@/github.com/kawai-network/veridium/internal/services/fileservice';

export class GenerationTopicModel {
  private logger = createModelLogger('GenerationTopic', 'GenerationTopicModel', 'database/models/generationTopic');

  /**
   * Show error notification to user
   */
  private async showErrorNotification(title: string, message: string) {
    try {
      await NotificationService.SendNotification(
        new NotificationOptions({
          id: `generation-topic-error-${Date.now()}`,
          title: `Generation Topic Error: ${title}`,
          body: message,
        })
      );
    } catch (notifError) {
      // Silently fail if notification fails - don't want notification errors to break the app
      console.error('Failed to show notification:', notifError);
    }
  }

  queryAll = async (): Promise<ImageGenerationTopic[]> => {
    try {
      const topics: GenerationTopic[] = await DB.ListGenerationTopics();

      return Promise.all(
        topics.map(async (topic): Promise<ImageGenerationTopic> => {
          const coverUrl = getNullableString(topic.coverUrl as any);
          const fullCoverUrl = coverUrl ? await GetFullFileUrl(coverUrl) : null;

          return {
            id: topic.id,
            title: getNullableString(topic.title as any),
            coverUrl: fullCoverUrl,
            createdAt: new Date(topic.createdAt),
            updatedAt: new Date(topic.updatedAt),
          };
        }),
      );
    } catch (error) {
      await this.logger.error('Failed to query generation topics', { error });
      throw error;
    }
  };

  create = async (title: string): Promise<GenerationTopic> => {
    await this.logger.methodEntry('create', { title });

    try {
      const id = nanoid();
      const now = currentTimestampMs();

      const newGenerationTopic = await DB.CreateGenerationTopic({
        id,
        title: toNullString(title),
        coverUrl: toNullString(null),
        createdAt: now,
        updatedAt: now,
      });

      await this.logger.methodExit('create', { id: newGenerationTopic.id });
      return newGenerationTopic;
    } catch (error) {
      await this.logger.error('Failed to create generation topic', { error, title });
      await this.showErrorNotification(
        'Create Failed',
        `Failed to create generation topic "${title}". Please try again.`
      );
      throw error;
    }
  };

  update = async (
    id: string,
    data: Partial<ImageGenerationTopic>,
  ): Promise<GenerationTopic | undefined> => {
    await this.logger.methodEntry('update', { id, data });

    try {
      // 1. Fetch existing topic to preserve fields
      const existing = await DB.GetGenerationTopic(id);
      if (!existing) {
        throw new Error(`Topic with id ${id} not found`);
      }

      // 2. Prepare values (use existing if not provided in update)
      // Note: check for undefined strictly to allow setting null/empty if intended
      const titleToSave = data.title !== undefined ? data.title : getNullableString(existing.title as any);
      const coverUrlToSave = data.coverUrl !== undefined ? data.coverUrl : getNullableString(existing.coverUrl as any);

      const updatedTopic: GenerationTopic = await DB.UpdateGenerationTopic({
        id,
        title: toNullString(titleToSave),
        coverUrl: toNullString(coverUrlToSave),
        updatedAt: currentTimestampMs(),
      });

      await this.logger.methodExit('update', { id });
      return updatedTopic;
    } catch (error) {
      await this.logger.error('Failed to update generation topic', { error, id, data });
      await this.showErrorNotification(
        'Update Failed',
        `Failed to update generation topic "${data.title || ''}". Please try again.`
      );
      throw error;
    }
  };

  /**
   * OPTIMIZED: Uses query to fetch all assets before deletion
   * Returns file URLs for cleanup
   *
   * This method follows the "database first, files second" deletion principle:
   * 1. First queries the topic with all its batches and generations to collect file URLs
   * 2. Then deletes the database record (cascade delete handles related batches and generations)
   * 3. Returns the deleted topic data and file URLs for cleanup
   */
  delete = async (
    id: string,
  ): Promise<{ deletedTopic: GenerationTopic; filesToDelete: string[] } | undefined> => {
    await this.logger.methodEntry('delete', { id });

    try {
      // 1. Get topic to verify ownership
      const topic = await DB.GetGenerationTopic(id);

      if (!topic) {
        await this.logger.warn('Generation topic not found', { id });
        return undefined;
      }

      // 2. Get all assets from generations under this topic
      const assets = await DB.GetGenerationTopicAssets(id);

      // 3. Collect all file URLs
      const filesToDelete: string[] = [];

      // Add cover image URL if exists
      const coverUrl = getNullableString(topic.coverUrl as any);
      if (coverUrl) {
        filesToDelete.push(coverUrl);
      }

      // Add thumbnail URLs from all generations
      for (const row of assets) {
        const asset = parseNullableJSON(row.asset as any) as GenerationAsset;
        if (asset?.thumbnailUrl) {
          filesToDelete.push(asset.thumbnailUrl);
        }
      }

      // 4. Delete the topic record (cascade will handle batches and generations)
      await DB.DeleteGenerationTopic(id);

      await this.logger.methodExit('delete', { id, filesCount: filesToDelete.length });
      return {
        deletedTopic: topic,
        filesToDelete,
      };
    } catch (error) {
      await this.logger.error('Failed to delete generation topic', { error, id });
      await this.showErrorNotification(
        'Delete Failed',
        `Failed to delete generation topic. Please try again.`
      );
      throw error;
    }
  };
}

