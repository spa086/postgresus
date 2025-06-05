import { getApplicationServer } from '../../../constants';
import RequestOptions from '../../../shared/api/RequestOptions';
import { apiHelper } from '../../../shared/api/apiHelper';
import type { Database } from '../model/Database';

export const databaseApi = {
  async createDatabase(database: Database) {
    const requestOptions: RequestOptions = new RequestOptions();
    requestOptions.setBody(JSON.stringify(database));
    return apiHelper.fetchPostJson<Database>(
      `${getApplicationServer()}/api/v1/databases/create`,
      requestOptions,
    );
  },

  async updateDatabase(database: Database) {
    const requestOptions: RequestOptions = new RequestOptions();
    requestOptions.setBody(JSON.stringify(database));
    return apiHelper.fetchPostJson<Database>(
      `${getApplicationServer()}/api/v1/databases/update`,
      requestOptions,
    );
  },

  async getDatabase(id: string) {
    const requestOptions: RequestOptions = new RequestOptions();
    return apiHelper.fetchGetJson<Database>(
      `${getApplicationServer()}/api/v1/databases/${id}`,
      requestOptions,
    );
  },

  async getDatabases() {
    const requestOptions: RequestOptions = new RequestOptions();
    return apiHelper.fetchGetJson<Database[]>(
      `${getApplicationServer()}/api/v1/databases`,
      requestOptions,
    );
  },

  async deleteDatabase(id: string) {
    const requestOptions: RequestOptions = new RequestOptions();
    return apiHelper.fetchDeleteRaw(
      `${getApplicationServer()}/api/v1/databases/${id}`,
      requestOptions,
    );
  },

  async testDatabaseConnection(id: string) {
    const requestOptions: RequestOptions = new RequestOptions();
    return apiHelper.fetchPostJson(
      `${getApplicationServer()}/api/v1/databases/${id}/test-connection`,
      requestOptions,
    );
  },

  async testDatabaseConnectionDirect(database: Database) {
    const requestOptions: RequestOptions = new RequestOptions();
    requestOptions.setBody(JSON.stringify(database));
    return apiHelper.fetchPostJson(
      `${getApplicationServer()}/api/v1/databases/test-connection-direct`,
      requestOptions,
    );
  },

  async isNotifierUsing(notifierId: string) {
    const requestOptions: RequestOptions = new RequestOptions();
    return apiHelper
      .fetchGetJson<{
        isUsing: boolean;
      }>(
        `${getApplicationServer()}/api/v1/databases/notifier/${notifierId}/is-using`,
        requestOptions,
      )
      .then((res) => res.isUsing);
  },

  async isStorageUsing(storageId: string): Promise<boolean> {
    const requestOptions: RequestOptions = new RequestOptions();
    return apiHelper
      .fetchGetJson<{
        isUsing: boolean;
      }>(`${getApplicationServer()}/api/v1/databases/storage/${storageId}/is-using`, requestOptions)
      .then((res) => res.isUsing);
  },
};
