import {
  CheckCircleOutlined,
  CloudUploadOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
  InfoCircleOutlined,
  SyncOutlined,
} from '@ant-design/icons';
import { Button, Modal, Table, Tooltip } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import dayjs from 'dayjs';
import { useEffect, useRef, useState } from 'react';

import { type Backup, BackupStatus, backupsApi } from '../../../entity/backups';
import type { Database } from '../../../entity/databases';
import { getUserTimeFormat } from '../../../shared/time';
import { ConfirmationComponent } from '../../../shared/ui';
import { RestoresComponent } from '../../restores';

interface Props {
  database: Database;
}

export const BackupsComponent = ({ database }: Props) => {
  const [isLoading, setIsLoading] = useState(false);
  const [backups, setBackups] = useState<Backup[]>([]);

  const [isMakeBackupRequestLoading, setIsMakeBackupRequestLoading] = useState(false);

  const [showingBackupError, setShowingBackupError] = useState<Backup | undefined>();

  const [deleteConfimationId, setDeleteConfimationId] = useState<string | undefined>();
  const [deletingBackupId, setDeletingBackupId] = useState<string | undefined>();

  const [showingRestoresBackupId, setShowingRestoresBackupId] = useState<string | undefined>();

  const isReloadInProgress = useRef(false);

  const loadBackups = async () => {
    if (isReloadInProgress.current) {
      return;
    }

    isReloadInProgress.current = true;

    try {
      const backups = await backupsApi.getBackups(database.id);
      setBackups(backups);
    } catch (e) {
      alert((e as Error).message);
    }

    isReloadInProgress.current = false;
  };

  const makeBackup = async () => {
    setIsMakeBackupRequestLoading(true);

    try {
      await backupsApi.makeBackup(database.id);
      await loadBackups();
    } catch (e) {
      alert((e as Error).message);
    }

    setIsMakeBackupRequestLoading(false);
  };

  const deleteBackup = async () => {
    if (!deleteConfimationId) {
      return;
    }

    setDeleteConfimationId(undefined);
    setDeletingBackupId(deleteConfimationId);

    try {
      await backupsApi.deleteBackup(deleteConfimationId);
      await loadBackups();
    } catch (e) {
      alert((e as Error).message);
    }

    setDeletingBackupId(undefined);
    setDeleteConfimationId(undefined);
  };

  useEffect(() => {
    setIsLoading(true);
    loadBackups().then(() => setIsLoading(false));

    const interval = setInterval(() => {
      loadBackups();
    }, 1_000);

    return () => clearInterval(interval);
  }, [database]);

  const columns: ColumnsType<Backup> = [
    {
      title: 'Created at',
      dataIndex: 'createdAt',
      key: 'createdAt',
      render: (createdAt: string) => (
        <div>
          {dayjs.utc(createdAt).local().format(getUserTimeFormat().format)} <br />
          <span className="text-gray-500">({dayjs.utc(createdAt).local().fromNow()})</span>
        </div>
      ),
      sorter: (a, b) => dayjs(a.createdAt).unix() - dayjs(b.createdAt).unix(),
      defaultSortOrder: 'descend',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status: BackupStatus, record: Backup) => {
        if (status === BackupStatus.FAILED) {
          return (
            <Tooltip title="Click to see error details">
              <div
                className="flex cursor-pointer items-center text-red-600 underline"
                onClick={() => setShowingBackupError(record)}
              >
                <ExclamationCircleOutlined className="mr-2" style={{ fontSize: 16 }} />

                <div>Failed</div>
              </div>
            </Tooltip>
          );
        }

        if (status === BackupStatus.COMPLETED) {
          return (
            <div className="flex items-center text-green-600">
              <CheckCircleOutlined className="mr-2" style={{ fontSize: 16 }} />
              <div>Successful</div>
            </div>
          );
        }

        if (status === BackupStatus.DELETED) {
          return (
            <div className="flex items-center text-gray-600">
              <DeleteOutlined className="mr-2" style={{ fontSize: 16 }} />
              <div>Deleted</div>
            </div>
          );
        }

        if (status === BackupStatus.IN_PROGRESS) {
          return (
            <div className="flex items-center font-bold text-blue-600">
              <SyncOutlined spin />
              <span className="ml-2">In progress</span>
            </div>
          );
        }

        return <span className="font-bold">{status}</span>;
      },
      filters: [
        {
          value: BackupStatus.IN_PROGRESS,
          text: 'In progress',
        },
        {
          value: BackupStatus.FAILED,
          text: 'Failed',
        },
        {
          value: BackupStatus.COMPLETED,
          text: 'Successful',
        },
        {
          value: BackupStatus.DELETED,
          text: 'Deleted',
        },
      ],
      onFilter: (value, record) => record.status === value,
    },
    {
      title: (
        <div className="flex items-center">
          Size
          <Tooltip
            className="ml-1"
            title="The file size we actually store in the storage (local, S3, Google Drive, etc.), usually compressed in ~5x times"
          >
            <InfoCircleOutlined />
          </Tooltip>
        </div>
      ),
      dataIndex: 'backupSizeMb',
      key: 'backupSizeMb',
      width: 150,
      render: (sizeMb: number) => {
        if (sizeMb >= 1024) {
          const sizeGb = sizeMb / 1024;
          return `${Number(sizeGb.toFixed(2)).toLocaleString()} GB`;
        }
        return `${Number(sizeMb?.toFixed(2)).toLocaleString()} MB`;
      },
    },
    {
      title: 'Duration',
      dataIndex: 'backupDurationMs',
      key: 'backupDurationMs',
      width: 150,
      render: (durationMs: number) => {
        const hours = Math.floor(durationMs / 3600000);
        const minutes = Math.floor((durationMs % 3600000) / 60000);
        const seconds = Math.floor((durationMs % 60000) / 1000);

        if (hours > 0) {
          return `${hours}h ${minutes}m ${seconds}s`;
        }

        return `${minutes}m ${seconds}s`;
      },
    },
    {
      title: 'Actions',
      dataIndex: '',
      key: '',
      render: (_, record: Backup) => {
        return (
          <div className="flex gap-2 text-lg">
            {record.status === BackupStatus.COMPLETED && (
              <div>
                {deletingBackupId === record.id ? (
                  <SyncOutlined spin />
                ) : (
                  <>
                    <Tooltip title="Delete backup">
                      <DeleteOutlined
                        className="cursor-pointer"
                        onClick={() => {
                          if (deletingBackupId) return;
                          setDeleteConfimationId(record.id);
                        }}
                        style={{ color: '#ff0000', opacity: deletingBackupId ? 0.2 : 1 }}
                      />
                    </Tooltip>

                    <Tooltip className="ml-3" title="Restore from backup">
                      <CloudUploadOutlined
                        className="cursor-pointer"
                        onClick={() => {
                          setShowingRestoresBackupId(record.id);
                        }}
                        style={{
                          color: '#0d6efd',
                        }}
                      />
                    </Tooltip>
                  </>
                )}
              </div>
            )}
          </div>
        );
      },
    },
  ];

  return (
    <div>
      <h2 className="text-xl font-bold">Backups</h2>

      <div className="mt-5" />

      <div className="flex">
        <Button
          onClick={makeBackup}
          className="mr-1"
          type="primary"
          disabled={isMakeBackupRequestLoading}
          loading={isMakeBackupRequestLoading}
        >
          Make backup right now
        </Button>
      </div>

      <div className="mt-5 max-w-[850px]">
        <Table
          bordered
          columns={columns}
          dataSource={backups}
          rowKey="id"
          loading={isLoading}
          size="small"
          pagination={false}
        />
      </div>

      {deleteConfimationId && (
        <ConfirmationComponent
          onConfirm={deleteBackup}
          onDecline={() => setDeleteConfimationId(undefined)}
          description="Are you sure you want to delete this backup?"
          actionButtonColor="red"
          actionText="Delete"
        />
      )}

      {showingRestoresBackupId && (
        <Modal
          width={400}
          open={!!showingRestoresBackupId}
          onCancel={() => setShowingRestoresBackupId(undefined)}
          title="Restore from backup"
          footer={null}
        >
          <RestoresComponent
            database={database}
            backup={backups.find((b) => b.id === showingRestoresBackupId) as Backup}
          />
        </Modal>
      )}

      {showingBackupError && (
        <Modal
          title="Backup error details"
          open={!!showingBackupError}
          onCancel={() => setShowingBackupError(undefined)}
          footer={null}
        >
          <div className="text-sm">{showingBackupError.failMessage}</div>
        </Modal>
      )}
    </div>
  );
};
