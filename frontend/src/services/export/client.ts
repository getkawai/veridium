import { DataExporterRepos } from '@/database/repositories/dataExporter';
import { BaseClientService } from '@/services/baseClientService';
import { DB } from '@/types/database';

export class ClientService extends BaseClientService {
  private get dataExporterRepos(): DataExporterRepos {
    return new DataExporterRepos(DB as any, this.userId);
  }

  exportData = async () => {
    const data = await this.dataExporterRepos.export();

    return { data, url: undefined };
  };
}
