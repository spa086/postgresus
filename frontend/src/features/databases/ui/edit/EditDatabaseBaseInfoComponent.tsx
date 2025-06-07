import { InfoCircleOutlined } from '@ant-design/icons';
import { Button, Input, InputNumber, Select, TimePicker, Tooltip } from 'antd';
import dayjs, { Dayjs } from 'dayjs';
import { useEffect, useMemo, useState } from 'react';

import { type Database, databaseApi } from '../../../../entity/databases';
import { Period } from '../../../../entity/databases/model/Period';
import { type Interval, IntervalType } from '../../../../entity/intervals';
import {
  getLocalDayOfMonth,
  getLocalWeekday,
  getUserTimeFormat,
  getUtcDayOfMonth,
  getUtcWeekday,
} from '../../../../shared/time/utils';

const weekdayOptions = [
  { value: 1, label: 'Mon' },
  { value: 2, label: 'Tue' },
  { value: 3, label: 'Wed' },
  { value: 4, label: 'Thu' },
  { value: 5, label: 'Fri' },
  { value: 6, label: 'Sat' },
  { value: 7, label: 'Sun' },
];

interface Props {
  database: Database;

  isShowName?: boolean;
  isShowCancelButton?: boolean;
  onCancel: () => void;

  saveButtonText?: string;
  isSaveToApi: boolean;
  onSaved: (db: Database) => void;
}

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

  const timeFormat = useMemo(() => {
    const is12 = getUserTimeFormat();
    return { use12Hours: is12, format: is12 ? 'h:mm A' : 'HH:mm' };
  }, []);

  const updateDatabase = (patch: Partial<Database>) => {
    setEditingDatabase((prev) => (prev ? { ...prev, ...patch } : prev));
    setIsUnsaved(true);
  };

  const saveInterval = (patch: Partial<Interval>) => {
    setEditingDatabase((prev) => {
      if (!prev) return prev;

      const updatedBackupInterval = { ...(prev.backupInterval ?? {}), ...patch };

      if (!updatedBackupInterval.id && prev.backupInterval?.id) {
        updatedBackupInterval.id = prev.backupInterval.id;
      }

      return { ...prev, backupInterval: updatedBackupInterval as Interval };
    });

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

  useEffect(() => {
    setIsSaving(false);
    setIsUnsaved(false);
    setEditingDatabase({ ...database });
  }, [database]);

  if (!editingDatabase) return null;
  const { backupInterval } = editingDatabase;

  // UTC → local conversions for display
  const localTime: Dayjs | undefined = backupInterval?.timeOfDay
    ? dayjs.utc(backupInterval.timeOfDay, 'HH:mm').local()
    : undefined;

  const displayedWeekday: number | undefined =
    backupInterval?.interval === IntervalType.WEEKLY &&
    backupInterval.weekday &&
    backupInterval.timeOfDay
      ? getLocalWeekday(backupInterval.weekday, backupInterval.timeOfDay)
      : backupInterval?.weekday;

  const displayedDayOfMonth: number | undefined =
    backupInterval?.interval === IntervalType.MONTHLY &&
    backupInterval.dayOfMonth &&
    backupInterval.timeOfDay
      ? getLocalDayOfMonth(backupInterval.dayOfMonth, backupInterval.timeOfDay)
      : backupInterval?.dayOfMonth;

  // mandatory-field check
  const isAllFieldsFilled =
    Boolean(editingDatabase.name) &&
    Boolean(editingDatabase.storePeriod) &&
    Boolean(backupInterval?.interval) &&
    (!backupInterval ||
      ((backupInterval.interval !== IntervalType.WEEKLY || displayedWeekday) &&
        (backupInterval.interval !== IntervalType.MONTHLY || displayedDayOfMonth)));

  return (
    <div>
      {isShowName && (
        <div className="mb-1 flex w-full items-center">
          <div className="min-w-[150px]">Name</div>
          <Input
            value={editingDatabase.name || ''}
            onChange={(e) => updateDatabase({ name: e.target.value })}
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
          onChange={(v) => saveInterval({ interval: v })}
          size="small"
          className="max-w-[200px] grow"
          options={[
            { label: 'Hourly', value: IntervalType.HOURLY },
            { label: 'Daily', value: IntervalType.DAILY },
            { label: 'Weekly', value: IntervalType.WEEKLY },
            { label: 'Monthly', value: IntervalType.MONTHLY },
          ]}
        />
      </div>

      {backupInterval?.interval === IntervalType.WEEKLY && (
        <div className="mb-1 flex w-full items-center">
          <div className="min-w-[150px]">Backup weekday</div>
          <Select
            value={displayedWeekday}
            onChange={(localWeekday) => {
              if (!localWeekday) return;
              const ref = localTime ?? dayjs();
              saveInterval({ weekday: getUtcWeekday(localWeekday, ref) });
            }}
            size="small"
            className="max-w-[200px] grow"
            options={weekdayOptions}
          />
        </div>
      )}

      {backupInterval?.interval === IntervalType.MONTHLY && (
        <div className="mb-1 flex w-full items-center">
          <div className="min-w-[150px]">Backup day of month</div>
          <InputNumber
            min={1}
            max={31}
            value={displayedDayOfMonth}
            onChange={(localDom) => {
              if (!localDom) return;
              const ref = localTime ?? dayjs();
              saveInterval({ dayOfMonth: getUtcDayOfMonth(localDom, ref) });
            }}
            size="small"
            className="max-w-[200px] grow"
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
            allowClear={false}
            size="small"
            className="max-w-[200px] grow"
            onChange={(t) => {
              if (!t) return;
              const patch: Partial<Interval> = { timeOfDay: t.utc().format('HH:mm') };

              if (backupInterval?.interval === IntervalType.WEEKLY && displayedWeekday) {
                patch.weekday = getUtcWeekday(displayedWeekday, t);
              }
              if (backupInterval?.interval === IntervalType.MONTHLY && displayedDayOfMonth) {
                patch.dayOfMonth = getUtcDayOfMonth(displayedDayOfMonth, t);
              }

              saveInterval(patch);
            }}
          />
        </div>
      )}

      <div className="mt-4 mb-1 flex w-full items-center">
        <div className="min-w-[150px]">Store period</div>
        <Select
          value={editingDatabase.storePeriod}
          onChange={(v) => updateDatabase({ storePeriod: v })}
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
          title="How long to keep the backups? Make sure you have enough storage space."
        >
          <InfoCircleOutlined className="ml-2" style={{ color: 'gray' }} />
        </Tooltip>
      </div>

      <div className="mt-5 flex">
        {isShowCancelButton && (
          <Button danger ghost className="mr-1" onClick={onCancel}>
            Cancel
          </Button>
        )}
        <Button
          type="primary"
          className={`${isShowCancelButton ? 'ml-1' : 'ml-auto'} mr-5`}
          onClick={saveDatabase}
          loading={isSaving}
          disabled={!isUnsaved || !isAllFieldsFilled}
        >
          {saveButtonText || 'Save'}
        </Button>
      </div>
    </div>
  );
};
