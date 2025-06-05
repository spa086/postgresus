import { StorageType } from './StorageType';

export const getStorageLogoFromType = (type: StorageType) => {
  switch (type) {
    case StorageType.LOCAL:
      return '/icons/storages/local.svg';
    case StorageType.S3:
      return '/icons/storages/s3.svg';
    default:
      return '';
  }
};
