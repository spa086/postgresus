import { StorageType } from './StorageType';

export const getStorageNameFromType = (type: StorageType) => {
  switch (type) {
    case StorageType.LOCAL:
      return 'local storage';
    case StorageType.S3:
      return 'S3';
    case StorageType.GOOGLE_DRIVE:
      return 'Google Drive';
    default:
      return '';
  }
};
