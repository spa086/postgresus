import { Button, Modal, Select, Spin } from 'antd';
import { useEffect, useState } from 'react';

import { type Database, databaseApi } from '../../../../entity/databases';
import { type Storage, storageApi } from '../../../../entity/storages';
import { ConfirmationComponent } from '../../../../shared/ui';
import { EditStorageComponent } from '../../../storages/ui/edit/EditStorageComponent';

interface Props {
  database: Database;

  isShowCancelButton?: boolean;
  onCancel: () => void;

  isShowBackButton: boolean;
  onBack: () => void;

  isShowSaveOnlyForUnsaved: boolean;
  saveButtonText?: string;
  isSaveToApi: boolean;
  onSaved: (database: Database) => void;
}

export const EditDatabaseStorageComponent = ({
  database,

  isShowCancelButton,
  onCancel,

  isShowBackButton,
  onBack,

  isShowSaveOnlyForUnsaved,
  saveButtonText,
  isSaveToApi,
  onSaved,
}: Props) => {
  const [editingDatabase, setEditingDatabase] = useState<Database>();
  const [isUnsaved, setIsUnsaved] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  const [storages, setStorages] = useState<Storage[]>([]);
  const [isStoragesLoading, setIsStoragesLoading] = useState(false);
  const [isShowCreateStorage, setShowCreateStorage] = useState(false);

  const [isShowWarn, setIsShowWarn] = useState(false);

  const saveDatabase = async () => {
    if (!editingDatabase) return;

    if (isSaveToApi) {
      setIsSaving(true);

      try {
        await databaseApi.updateDatabase(editingDatabase);
        setIsUnsaved(false);
      } catch (e) {
        alert((e as Error).message);
      }

      setIsSaving(false);
    }

    onSaved(editingDatabase);
  };

  const loadStorages = async () => {
    setIsStoragesLoading(true);

    try {
      const storages = await storageApi.getStorages();
      setStorages(storages);
    } catch (e) {
      alert((e as Error).message);
    }

    setIsStoragesLoading(false);
  };

  useEffect(() => {
    setIsSaving(false);
    setEditingDatabase({ ...database });
    loadStorages();

    if (database.storage.id) {
      setIsShowWarn(true);
    }
  }, [database]);

  if (!editingDatabase) return null;

  if (isStoragesLoading)
    return (
      <div className="mb-5 flex items-center">
        <Spin />
      </div>
    );

  return (
    <div>
      <div className="mb-5 max-w-[275px] text-gray-500">
        Storage - is a place where backups will be stored (local disk, S3, Google Drive, etc.)
      </div>

      <div className="mb-5 flex w-full items-center">
        <div className="min-w-[150px]">Storages</div>

        <Select
          value={editingDatabase.storage.id}
          onChange={(storageId) => {
            if (storageId.includes('create-new-storage')) {
              setShowCreateStorage(true);
              return;
            }

            setEditingDatabase({
              ...editingDatabase,
              storage: storages.find((s) => s.id === storageId),
            } as unknown as Database);

            setIsUnsaved(true);
          }}
          size="small"
          className="max-w-[200px] grow"
          options={[
            ...storages.map((s) => ({ label: s.name, value: s.id })),
            { label: 'Create new storage', value: 'create-new-storage' },
          ]}
          placeholder="Select storages"
        />
      </div>

      <div className="mt-5 flex">
        {isShowCancelButton && (
          <Button className="mr-1" danger ghost onClick={() => onCancel()}>
            Cancel
          </Button>
        )}

        {isShowBackButton && (
          <Button className="mr-auto" type="primary" ghost onClick={() => onBack()}>
            Back
          </Button>
        )}

        {(!isShowSaveOnlyForUnsaved || isUnsaved) && (
          <Button
            type="primary"
            onClick={() => saveDatabase()}
            loading={isSaving}
            disabled={isSaving}
            className="mr-5"
          >
            {saveButtonText || 'Save'}
          </Button>
        )}
      </div>

      {isShowCreateStorage && (
        <Modal
          title="Add storage"
          footer={<div />}
          open={isShowCreateStorage}
          onCancel={() => setShowCreateStorage(false)}
        >
          <div className="my-3 max-w-[275px] text-gray-500">
            Storage - is a place where backups will be stored (local disk, S3, Google Drive, etc.)
          </div>

          <EditStorageComponent
            isShowName
            isShowClose={false}
            onClose={() => setShowCreateStorage(false)}
            onChanged={() => {
              loadStorages();
              setShowCreateStorage(false);
            }}
          />
        </Modal>
      )}

      {isShowWarn && (
        <ConfirmationComponent
          onConfirm={() => {
            setIsShowWarn(false);
          }}
          onDecline={() => {
            setIsShowWarn(false);
          }}
          description="If you change the storage, all backups in this storage will be deleted."
          actionButtonColor="red"
          actionText="I understand"
          cancelText="Cancel"
          hideCancelButton
        />
      )}
    </div>
  );
};
