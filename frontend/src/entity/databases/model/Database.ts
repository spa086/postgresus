import type { Interval } from '../../intervals';
import type { Notifier } from '../../notifiers';
import type { BackupNotificationType } from './BackupNotificationType';
import type { DatabaseType } from './DatabaseType';
import type { HealthStatus } from './HealthStatus';
import type { Period } from './Period';
import type { PostgresqlDatabase } from './postgresql/PostgresqlDatabase';

export interface Database {
  id: string;
  name: string;

  type: DatabaseType;

  backupInterval?: Interval;
  storePeriod: Period;

  postgresql?: PostgresqlDatabase;

  storage: Storage;

  notifiers: Notifier[];
  sendNotificationsOn: BackupNotificationType[];

  lastBackupTime?: Date;
  lastBackupErrorMessage?: string;

  healthStatus?: HealthStatus;
}
