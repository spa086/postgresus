import { useState } from 'react';

import { backupsApi } from '../../../entity/backups';
import {
  type Database,
  DatabaseType,
  Period,
  type PostgresqlDatabase,
  databaseApi,
} from '../../../entity/databases';
import { EditDatabaseBaseInfoComponent } from './edit/EditDatabaseBaseInfoComponent';
import { EditDatabaseNotifiersComponent } from './edit/EditDatabaseNotifiersComponent';
import { EditDatabaseSpecificDataComponent } from './edit/EditDatabaseSpecificDataComponent';
import { EditDatabaseStorageComponent } from './edit/EditDatabaseStorageComponent';

interface Props {
  onCreated: () => void;

  onClose: () => void;
}

export const CreateDatabaseComponent = ({ onCreated, onClose }: Props) => {
  const [isCreating, setIsCreating] = useState(false);
  const [database, setDatabase] = useState<Database>({
    id: undefined as unknown as string,
    name: '',
    storePeriod: Period.WEEK,

    postgresql: {
      cpuCount: 1,
    } as unknown as PostgresqlDatabase,

    type: DatabaseType.POSTGRES,

    storage: {} as unknown as Storage,

    notifiers: [],
    sendNotificationsOn: [],
  } as Database);

  const [step, setStep] = useState<'base-info' | 'db-settings' | 'storages' | 'notifiers'>(
    'base-info',
  );

  const createDatabase = async (database: Database) => {
    setIsCreating(true);

    try {
      const createdDatabase = await databaseApi.createDatabase(database);
      setDatabase({ ...createdDatabase });

      await backupsApi.makeBackup(createdDatabase.id);
      onCreated();
      onClose();
    } catch (error) {
      alert(error);
    }

    setIsCreating(false);
  };

  if (step === 'base-info') {
    return (
      <div>
        <EditDatabaseBaseInfoComponent
          database={database}
          isShowName
          isSaveToApi={false}
          saveButtonText="Continue"
          onCancel={() => onClose()}
          onSaved={(database) => {
            setDatabase({ ...database });
            setStep('db-settings');
          }}
        />
      </div>
    );
  }

  if (step === 'db-settings') {
    return (
      <EditDatabaseSpecificDataComponent
        database={database}
        isShowCancelButton={false}
        onCancel={() => onClose()}
        isShowBackButton
        onBack={() => setStep('base-info')}
        saveButtonText="Continue"
        isSaveToApi={false}
        onSaved={(database) => {
          setDatabase({ ...database });
          setStep('storages');
        }}
      />
    );
  }

  if (step === 'storages') {
    return (
      <EditDatabaseStorageComponent
        database={database}
        isShowCancelButton={false}
        onCancel={() => onClose()}
        isShowBackButton
        onBack={() => setStep('db-settings')}
        isShowSaveOnlyForUnsaved={false}
        saveButtonText="Continue"
        isSaveToApi={false}
        onSaved={(database) => {
          setDatabase({ ...database });
          setStep('notifiers');
        }}
      />
    );
  }

  if (step === 'notifiers') {
    return (
      <EditDatabaseNotifiersComponent
        database={database}
        isShowCancelButton={false}
        onCancel={() => onClose()}
        isShowBackButton
        onBack={() => setStep('storages')}
        isShowSaveOnlyForUnsaved={false}
        saveButtonText="Complete"
        isSaveToApi={false}
        onSaved={(database) => {
          if (isCreating) return;

          setDatabase({ ...database });
          createDatabase(database);
        }}
      />
    );
  }
};
