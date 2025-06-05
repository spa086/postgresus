import dayjs from 'dayjs';
import { useMemo } from 'react';

import { type Database } from '../../../../entity/databases';
import { Period } from '../../../../entity/databases/model/Period';
import { IntervalType } from '../../../../entity/intervals';

interface Props {
  database: Database;
  isShowName?: boolean;
}

const weekdayLabels = {
  1: 'Mon',
  2: 'Tue',
  3: 'Wed',
  4: 'Thu',
  5: 'Fri',
  6: 'Sat',
  7: 'Sun',
};

const intervalLabels = {
  [IntervalType.HOURLY]: 'Hourly',
  [IntervalType.DAILY]: 'Daily',
  [IntervalType.WEEKLY]: 'Weekly',
  [IntervalType.MONTHLY]: 'Monthly',
};

const periodLabels = {
  [Period.DAY]: '1 day',
  [Period.WEEK]: '1 week',
  [Period.MONTH]: '1 month',
  [Period.THREE_MONTH]: '3 months',
  [Period.SIX_MONTH]: '6 months',
  [Period.YEAR]: '1 year',
  [Period.TWO_YEARS]: '2 years',
  [Period.THREE_YEARS]: '3 years',
  [Period.FOUR_YEARS]: '4 years',
  [Period.FIVE_YEARS]: '5 years',
  [Period.FOREVER]: 'Forever',
};

// Function to detect if user prefers 12-hour format based on their locale
const getUserTimeFormat = () => {
  const locale = navigator.language || 'en-US';
  const testDate = new Date(2023, 0, 1, 13, 0, 0); // 1 PM
  const timeString = testDate.toLocaleTimeString(locale, { hour: 'numeric' });
  return timeString.includes('PM') || timeString.includes('AM');
};

export const ShowDatabaseBaseInfoComponent = ({ database, isShowName }: Props) => {
  // Detect user's preferred time format (12-hour vs 24-hour)
  const timeFormat = useMemo(() => {
    const is12Hour = getUserTimeFormat();
    return {
      use12Hours: is12Hour,
      format: is12Hour ? 'h:mm A' : 'HH:mm',
    };
  }, []);

  const { backupInterval } = database;

  const localTime = backupInterval?.timeOfDay
    ? dayjs.utc(backupInterval.timeOfDay, 'HH:mm').local()
    : undefined;

  const formattedTime = localTime ? localTime.format(timeFormat.format) : '';

  return (
    <div>
      {isShowName && (
        <div className="mb-1 flex w-full items-center">
          <div className="min-w-[150px]">Name</div>
          <div>{database.name || ''}</div>
        </div>
      )}

      <div className="mt-4 mb-1 flex w-full items-center">
        <div className="min-w-[150px]">Backup interval</div>
        <div>{backupInterval?.interval ? intervalLabels[backupInterval.interval] : ''}</div>
      </div>

      {backupInterval?.interval === IntervalType.WEEKLY && (
        <div className="mb-1 flex w-full items-center">
          <div className="min-w-[150px]">Backup weekday</div>
          <div>
            {backupInterval.weekday
              ? weekdayLabels[backupInterval.weekday as keyof typeof weekdayLabels]
              : ''}
          </div>
        </div>
      )}

      {backupInterval?.interval === IntervalType.MONTHLY && (
        <div className="mb-1 flex w-full items-center">
          <div className="min-w-[150px]">Backup day of month</div>
          <div>{backupInterval.dayOfMonth || ''}</div>
        </div>
      )}

      {backupInterval?.interval !== IntervalType.HOURLY && (
        <div className="mb-1 flex w-full items-center">
          <div className="min-w-[150px]">Backup time of day</div>
          <div>{formattedTime}</div>
        </div>
      )}

      <div className="mt-4 mb-1 flex w-full items-center">
        <div className="min-w-[150px]">Store period</div>
        <div>{database.storePeriod ? periodLabels[database.storePeriod] : ''}</div>
      </div>
    </div>
  );
};
