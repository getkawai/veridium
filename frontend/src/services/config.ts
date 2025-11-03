import { BRANDING_NAME } from '@/const';
import { downloadFile, exportJSONFile } from '@/utils/client';
import dayjs from 'dayjs';

import { ImportPgDataStructure } from '@/types/export';

import { exportService } from './export';

class ConfigService {
  exportAll = async () => {
    const { data, url } = await exportService.exportData();
    const filename = `${dayjs().format('YYYY-MM-DD-hh-mm')}_${BRANDING_NAME}-data.json`;

    // if url exists, means export data from server and upload the data to S3
    // just need to download the file
    if (url) {
      await downloadFile(url, filename);
      return;
    }

    // or export to file with the data
    const result = await this.createDataStructure(data, 'pglite');

    exportJSONFile(result, filename);
  };

  exportAgents = async () => {
  };

  exportSingleAgent = async (agentId: string) => {
  };

  exportSessions = async () => {
  };

  exportSettings = async () => {
  };

  exportSingleSession = async (sessionId: string) => {
  };

  private createDataStructure = async (
    data: any,
    mode: 'pglite' | 'postgres',
  ): Promise<ImportPgDataStructure> => {
    return { data, mode, schemaHash: 'empty' };
  };
}

export const configService = new ConfigService();
