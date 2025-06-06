import { getApplicationServer } from '../../../constants';
import RequestOptions from '../../../shared/api/RequestOptions';
import { apiHelper } from '../../../shared/api/apiHelper';
import type { Backup } from '../model/Backup';

export const backupsApi = {
  async getBackups(databaseId: string) {
    return apiHelper.fetchGetJson<Backup[]>(
      `${getApplicationServer()}/api/v1/backups?database_id=${databaseId}`,
      undefined,
      true,
    );
  },

  async makeBackup(databaseId: string) {
    const requestOptions: RequestOptions = new RequestOptions();
    requestOptions.setBody(JSON.stringify({ database_id: databaseId }));
    return apiHelper.fetchPostJson<{ message: string }>(
      `${getApplicationServer()}/api/v1/backups`,
      requestOptions,
    );
  },

  async deleteBackup(id: string) {
    return apiHelper.fetchDeleteRaw(`${getApplicationServer()}/api/v1/backups/${id}`);
  },
};
