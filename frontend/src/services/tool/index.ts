// Tool service client
// This provides client-side functionality for tool operations

import { getToolManifest } from '@/utils/toolManifest';

export interface ToolManifest {
  identifier: string;
  name: string;
  description?: string;
  version?: string;
  api?: any[];
  meta?: any;
  type?: string;
}

export interface OldPluginListParams {
  // Define parameters as needed
  [key: string]: any;
}

export interface OldPluginItem {
  identifier: string;
  name: string;
  description?: string;
  version?: string;
  [key: string]: any;
}

export class ToolServiceClient {
  async getToolManifest(manifest: any, useProxy?: boolean): Promise<ToolManifest> {
    console.warn('ToolServiceClient.getToolManifest not fully implemented', manifest);
    // Use the existing utility function
    const lobeManifest = await getToolManifest(manifest, useProxy);
    
    // Convert LobeChatPluginManifest to ToolManifest
    return {
      identifier: lobeManifest.identifier,
      name: lobeManifest.meta?.title || lobeManifest.identifier,
      description: lobeManifest.meta?.description,
      version: lobeManifest.version,
      api: lobeManifest.api,
      meta: lobeManifest.meta,
      type: lobeManifest.type,
    };
  }

  async getOldPluginList(params?: OldPluginListParams): Promise<OldPluginItem[]> {
    console.warn('ToolServiceClient.getOldPluginList not implemented', params);
    // Return empty array for build purposes
    return [];
  }
}

export const toolService = new ToolServiceClient();
