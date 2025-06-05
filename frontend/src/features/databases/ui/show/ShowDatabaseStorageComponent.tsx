import { type Database } from '../../../../entity/databases';
import { getStorageLogoFromType } from '../../../../entity/storages/models/getStorageLogoFromType';

interface Props {
  database: Database;
}

export const ShowDatabaseStorageComponent = ({ database }: Props) => {
  return (
    <div>
      <div className="mb-5 flex w-full items-center">
        <div className="min-w-[150px]">Storage</div>
        <div>{database.storage?.name || ''}</div>{' '}
        <img
          src={getStorageLogoFromType(database.storage?.type)}
          alt="storageIcon"
          className="ml-1 h-4 w-4"
        />
      </div>
    </div>
  );
};
