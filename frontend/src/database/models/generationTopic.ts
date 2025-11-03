import { GenerationAsset, ImageGenerationTopic } from  '@/types';
import { nanoid } from 'nanoid';

import { FileService } from '@/server/services/file';

import {
  DB,
  toNullString,
  parseNullableJSON,
  getNullableString,
  currentTimestampMs,
} from '@/types/database';

export class GenerationTopicModel {
  private userId: string;
  private fileService: FileService;

  constructor(_db: any, userId: string) {
    this.userId = userId;
    // Note: FileService still needs db for now, pass dummy
    this.fileService = new FileService(null as any, userId);
  }

  queryAll = async () => {
    const topics = await DB.ListGenerationTopics(this.userId);

    return Promise.all(
      topics.map(async (topic) => {
        const coverUrl = getNullableString(topic.coverUrl as any);
        if (coverUrl) {
          return {
            ...topic,
            coverUrl: await this.fileService.getFullFileUrl(coverUrl),
          };
        }
        return topic;
      }),
    );
  };

  create = async (title: string) => {
    const id = nanoid();
    const now = currentTimestampMs();

    const newGenerationTopic = await DB.CreateGenerationTopic({
      id,
      userId: this.userId,
      title: toNullString(title),
      coverUrl: toNullString(null),
      createdAt: now,
      updatedAt: now,
    });

    return newGenerationTopic;
  };

  update = async (
    id: string,
    data: Partial<ImageGenerationTopic>,
  ): Promise<any | undefined> => {
    const updatedTopic = await DB.UpdateGenerationTopic({
      id,
      userId: this.userId,
      title: toNullString(data.title),
      coverUrl: toNullString(data.coverUrl as any),
      updatedAt: currentTimestampMs(),
    });

    return updatedTopic;
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
  ): Promise<{ deletedTopic: any; filesToDelete: string[] } | undefined> => {
    // 1. Get topic to verify ownership
    const topic = await DB.GetGenerationTopic({
      id,
      userId: this.userId,
    });

    if (!topic) {
      return undefined;
    }

    // 2. Get all assets from generations under this topic
    const assets = await DB.GetGenerationTopicAssets({
      id: toNullString(id) as any,
      userId: this.userId,
    });

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
    await DB.DeleteGenerationTopic({
      id,
      userId: this.userId,
    });

    return {
      deletedTopic: topic,
      filesToDelete,
    };
  };
}

