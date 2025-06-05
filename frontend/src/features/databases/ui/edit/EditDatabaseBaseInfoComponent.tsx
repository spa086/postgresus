import { InfoCircleOutlined } from '@ant-design/icons';
import { Button, Input, InputNumber, Select, TimePicker, Tooltip } from 'antd';
import dayjs, { Dayjs } from 'dayjs';
import { useEffect, useMemo, useState } from 'react';

import { type Database, databaseApi } from '../../../../entity/databases';
import { Period } from '../../../../entity/databases/model/Period';
import { type Interval, IntervalType } from '../../../../entity/intervals';

interface Props {
  database: Database;

  isShowName?: boolean;

  isShowCancelButton?: boolean;
  onCancel: () => void;

  saveButtonText?: string;
  isSaveToApi: boolean;
  onSaved: (database: Database) => void;
}

const weekdayOptions = [
  { value: 1, label: 'Mon' },
  { value: 2, label: 'Tue' },
  { value: 3, label: 'Wed' },
  { value: 4, label: 'Thu' },
  { value: 5, label: 'Fri' },
  { value: 6, label: 'Sat' },
  { value: 7, label: 'Sun' },
];

// Function to detect if user prefers 12-hour format based on their locale
const getUserTimeFormat = () => {
  const locale = navigator.language || 'en-US';
  const testDate = new Date(2023, 0, 1, 13, 0, 0); // 1 PM
  const timeString = testDate.toLocaleTimeString(locale, { hour: 'numeric' });
  return timeString.includes('PM') || timeString.includes('AM');
};

export const EditDatabaseBaseInfoComponent = ({
  database,

  isShowName,

  isShowCancelButton,
  onCancel,

  saveButtonText,
  isSaveToApi,
  onSaved,
}: Props) => {
  const [editingDatabase, setEditingDatabase] = useState<Database>();
  const [isUnsaved, setIsUnsaved] = useState(false);
  const [isSaving, setIsSaving] = useState(false);

  // Detect user's preferred time format (12-hour vs 24-hour)
  const timeFormat = useMemo(() => {
    const is12Hour = getUserTimeFormat();
    return {
      use12Hours: is12Hour,
      format: is12Hour ? 'h:mm A' : 'HH:mm',
    };
  }, []);

  const updateDatabase = (patch: Partial<Database>) => {
    if (!editingDatabase) return;
    setEditingDatabase({ ...editingDatabase, ...patch });
    setIsUnsaved(true);
  };

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

  const saveInterval = (patch: Partial<Interval>) => {
    if (!editingDatabase) return;
    const current = editingDatabase.backupInterval ?? ({} as Interval);
    updateDatabase({ backupInterval: { ...current, ...patch } });
  };

  useEffect(() => {
    setIsSaving(false);
    setIsUnsaved(false);

    setEditingDatabase({ ...database });
  }, [database]);

  if (!editingDatabase) return null;

  const { backupInterval } = editingDatabase;

  const localTime: Dayjs | undefined = backupInterval?.timeOfDay
    ? dayjs.utc(backupInterval.timeOfDay, 'HH:mm').local() /* cast to user tz */
    : undefined;

  let isAllFieldsFilled = true;

  if (!editingDatabase.name) isAllFieldsFilled = false;
  if (!editingDatabase.storePeriod) isAllFieldsFilled = false;

  if (!editingDatabase.backupInterval?.interval) isAllFieldsFilled = false;
  if (editingDatabase.backupInterval?.interval === IntervalType.WEEKLY) {
    if (!editingDatabase.backupInterval?.weekday) isAllFieldsFilled = false;
  }
  if (editingDatabase.backupInterval?.interval === IntervalType.MONTHLY) {
    if (!editingDatabase.backupInterval.dayOfMonth) isAllFieldsFilled = false;
  }

  return (
    <div>
      {isShowName && (
        <div className="mb-1 flex w-full items-center">
          <div className="min-w-[150px]">Name</div>
          <Input
            value={editingDatabase.name || ''}
            onChange={(e) => {
              updateDatabase({ name: e.target.value });
            }}
            size="small"
            placeholder="My favourite DB"
            className="max-w-[200px] grow"
          />
        </div>
      )}

      <div className="mt-4 mb-1 flex w-full items-center">
        <div className="min-w-[150px]">Backup interval</div>
        <Select
          value={backupInterval?.interval}
          onChange={(v) => {
            saveInterval({ interval: v });
          }}
          size="small"
          className="max-w-[200px] grow"
          options={[
            { label: 'Hourly', value: IntervalType.HOURLY },
            { label: 'Daily', value: IntervalType.DAILY },
            { label: 'Weekly', value: IntervalType.WEEKLY },
            { label: 'Monthly', value: IntervalType.MONTHLY },
          ]}
          placeholder="Select backup interval"
        />
      </div>

      {backupInterval?.interval === IntervalType.WEEKLY && (
        <div className="mb-1 flex w-full items-center">
          <div className="min-w-[150px]">Backup weekday</div>
          <Select
            value={backupInterval.weekday}
            onChange={(v) => {
              saveInterval({ weekday: v });
            }}
            size="small"
            className="max-w-[200px] grow"
            options={weekdayOptions}
            placeholder="Select backup weekday"
          />
        </div>
      )}

      {backupInterval?.interval === IntervalType.MONTHLY && (
        <div className="mb-1 flex w-full items-center">
          <div className="min-w-[150px]">Backup day of month</div>
          <InputNumber
            min={1}
            max={31}
            value={backupInterval.dayOfMonth}
            onChange={(v) => {
              saveInterval({ dayOfMonth: v ?? 1 });
            }}
            size="small"
            className="max-w-[200px] grow"
            placeholder="Select backup day of month"
          />
        </div>
      )}

      {backupInterval?.interval !== IntervalType.HOURLY && (
        <div className="mb-1 flex w-full items-center">
          <div className="min-w-[150px]">Backup time of day</div>
          <TimePicker
            value={localTime}
            format={timeFormat.format}
            use12Hours={timeFormat.use12Hours}
            onChange={(t) => {
              if (!t) return;
              // convert local picker value → UTC "HH:mm"
              const utcString = t.utc().format('HH:mm');
              saveInterval({ timeOfDay: utcString });
            }}
            allowClear={false}
            size="small"
            className="max-w-[200px] grow"
          />
        </div>
      )}

      <div className="mt-4 mb-1 flex w-full items-center">
        <div className="min-w-[150px]">Store period</div>
        <Select
          value={editingDatabase.storePeriod}
          onChange={(v) => {
            updateDatabase({ storePeriod: v });
          }}
          size="small"
          className="max-w-[200px] grow"
          options={[
            { label: '1 day', value: Period.DAY },
            { label: '1 week', value: Period.WEEK },
            { label: '1 month', value: Period.MONTH },
            { label: '3 months', value: Period.THREE_MONTH },
            { label: '6 months', value: Period.SIX_MONTH },
            { label: '1 year', value: Period.YEAR },
            { label: '2 years', value: Period.TWO_YEARS },
            { label: '3 years', value: Period.THREE_YEARS },
            { label: '4 years', value: Period.FOUR_YEARS },
            { label: '5 years', value: Period.FIVE_YEARS },
            { label: 'Forever', value: Period.FOREVER },
          ]}
        />
        <Tooltip
          className="cursor-pointer"
          title="How long to keep the backups? Make sure that you have enough space on the storage you are using (local, S3, Goole Drive, etc.)."
        >
          <InfoCircleOutlined className="ml-2" style={{ color: 'gray' }} />
        </Tooltip>
      </div>

      <div className="mt-5 flex">
        {isShowCancelButton && (
          <Button className="mr-1" danger ghost onClick={() => onCancel()}>
            Cancel
          </Button>
        )}

        <Button
          className={`${isShowCancelButton ? 'ml-1' : 'ml-auto'} mr-5`}
          type="primary"
          onClick={() => saveDatabase()}
          loading={isSaving}
          disabled={!isUnsaved || !isAllFieldsFilled}
        >
          {saveButtonText || 'Save'}
        </Button>
      </div>
    </div>
  );
};
