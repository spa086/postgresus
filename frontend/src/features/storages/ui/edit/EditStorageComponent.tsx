import { Button, Input, Select } from 'antd';
import { useEffect, useState } from 'react';

import { type Storage, StorageType, storageApi } from '../../../../entity/storages';
import { ToastHelper } from '../../../../shared/toast';
import { EditS3StorageComponent } from './storages/EditS3StorageComponent';

interface Props {
  isShowClose: boolean;
  onClose: () => void;

  isShowName: boolean;

  editingStorage?: Storage;
  onChanged: (storage: Storage) => void;
}

export function EditStorageComponent({
  isShowClose,
  onClose,
  isShowName,
  editingStorage,
  onChanged,
}: Props) {
  const [storage, setStorage] = useState<Storage | undefined>();
  const [isUnsaved, setIsUnsaved] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  const [isTestingConnection, setIsTestingConnection] = useState(false);
  const [isTestConnectionSuccess, setIsTestConnectionSuccess] = useState(false);

  const save = async () => {
    if (!storage) return;

    setIsSaving(true);

    try {
      await storageApi.saveStorage(storage);
      onChanged(storage);
      setIsUnsaved(false);
    } catch (e) {
      alert((e as Error).message);
    }

    setIsSaving(false);
  };

  const testConnection = async () => {
    if (!storage) return;

    setIsTestingConnection(true);

    try {
      await storageApi.testStorageConnectionDirect(storage);
      setIsTestConnectionSuccess(true);
      ToastHelper.showToast({
        title: 'Connection test successful!',
        description: 'Storage connection tested successfully',
      });
    } catch (e) {
      alert((e as Error).message);
    }

    setIsTestingConnection(false);
  };

  const setStorageType = (type: StorageType) => {
    if (!storage) return;

    storage.localStorage = undefined;
    storage.s3Storage = undefined;

    if (type === StorageType.LOCAL) {
      storage.localStorage = {};
    }

    if (type === StorageType.S3) {
      storage.s3Storage = {
        s3Bucket: '',
        s3Region: '',
        s3AccessKey: '',
        s3SecretKey: '',
        s3Endpoint: '',
      };
    }

    setStorage(
      JSON.parse(
        JSON.stringify({
          ...storage,
          type: type,
        }),
      ),
    );
  };

  useEffect(() => {
    setIsUnsaved(false);
    setStorage(
      editingStorage
        ? JSON.parse(JSON.stringify(editingStorage))
        : {
            id: undefined as unknown as string,
            name: '',
            type: StorageType.LOCAL,
            localStorage: {},
          },
    );
  }, [editingStorage]);

  const isAllDataFilled = () => {
    if (!storage) return false;

    if (!storage.name) return false;

    if (storage.type === StorageType.LOCAL) {
      return true; // No additional settings required for local storage
    }

    if (storage.type === StorageType.S3) {
      return (
        storage.s3Storage?.s3Bucket &&
        storage.s3Storage?.s3Region &&
        storage.s3Storage?.s3AccessKey &&
        storage.s3Storage?.s3SecretKey
      );
    }

    return false;
  };

  if (!storage) return <div />;

  return (
    <div>
      {isShowName && (
        <div className="mb-1 flex items-center">
          <div className="min-w-[110px]">Name</div>

          <Input
            value={storage?.name || ''}
            onChange={(e) => {
              setStorage({ ...storage, name: e.target.value });
              setIsUnsaved(true);
            }}
            size="small"
            className="w-full max-w-[250px]"
            placeholder="My Storage"
          />
        </div>
      )}

      <div className="mb-1 flex items-center">
        <div className="min-w-[110px]">Type</div>

        <Select
          value={storage?.type}
          options={[
            { label: 'Local storage', value: StorageType.LOCAL },
            { label: 'S3', value: StorageType.S3 },
          ]}
          onChange={(value) => {
            setStorageType(value);
            setIsUnsaved(true);
          }}
          size="small"
          className="w-full max-w-[250px]"
        />
      </div>

      <div className="mt-5" />

      <div>
        {storage?.type === StorageType.S3 && (
          <EditS3StorageComponent
            storage={storage}
            setStorage={setStorage}
            setIsUnsaved={setIsUnsaved}
          />
        )}
      </div>

      <div className="mt-3 flex">
        {isUnsaved && !isTestConnectionSuccess ? (
          <Button
            className="mr-1"
            disabled={isTestingConnection || !isAllDataFilled()}
            loading={isTestingConnection}
            type="primary"
            onClick={testConnection}
          >
            Test connection
          </Button>
        ) : (
          <div />
        )}

        {isUnsaved && isTestConnectionSuccess ? (
          <Button
            className="mr-1"
            disabled={isSaving || !isAllDataFilled()}
            loading={isSaving}
            type="primary"
            onClick={save}
          >
            Save
          </Button>
        ) : (
          <div />
        )}

        {isShowClose ? (
          <Button
            className="mr-1"
            disabled={isSaving}
            type="primary"
            danger
            ghost
            onClick={onClose}
          >
            Cancel
          </Button>
        ) : (
          <div />
        )}
      </div>
    </div>
  );
}
