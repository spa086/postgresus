import type { LocalStorage } from './LocalStorage';
import type { S3Storage } from './S3Storage';
import type { StorageType } from './StorageType';

export interface Storage {
  id: string;
  type: StorageType;
  name: string;
  lastSaveError?: string;

  // specific storage types
  localStorage?: LocalStorage;
  s3Storage?: S3Storage;
}
