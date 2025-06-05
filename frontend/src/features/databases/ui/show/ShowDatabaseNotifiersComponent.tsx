import { type Database } from '../../../../entity/databases';
import { BackupNotificationType } from '../../../../entity/databases/model/BackupNotificationType';
import { getNotifierLogoFromType } from '../../../../entity/notifiers/models/getNotifierLogoFromType';

interface Props {
  database: Database;
}

const notificationTypeLabels = {
  [BackupNotificationType.BACKUP_FAILED]: 'backup failed',
  [BackupNotificationType.BACKUP_SUCCESS]: 'backup success',
};

export const ShowDatabaseNotifiersComponent = ({ database }: Props) => {
  const notificationLabels =
    database.sendNotificationsOn?.map((type) => notificationTypeLabels[type]).join(', ') || '';

  return (
    <div>
      <div className="mb-2 flex w-full">
        <div className="min-w-[150px]">Send notification when</div>
        <div>
          {notificationLabels.split(', ').map((label) => (
            <div key={label}>- {label}</div>
          ))}
        </div>
      </div>

      <div className="flex w-full">
        <div className="min-w-[150px]">Notify to</div>
        <div>
          {database.notifiers?.map((notifier) => (
            <div className="flex items-center" key={notifier.id}>
              <div>- {notifier.name}</div>
              <img src={getNotifierLogoFromType(notifier?.notifierType)} className="ml-1 h-4 w-4" />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
