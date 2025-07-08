import { type Database } from '../../../../entity/databases';

interface Props {
  database: Database;
  isShowName?: boolean;
}

export const ShowDatabaseBaseInfoComponent = ({ database, isShowName }: Props) => {
  return (
    <div>
      {isShowName && (
        <div className="mb-1 flex w-full items-center">
          <div className="min-w-[150px]">Name</div>
          <div>{database.name || ''}</div>
        </div>
      )}
    </div>
  );
};
