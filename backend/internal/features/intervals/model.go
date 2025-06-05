package intervals

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Interval struct {
	ID       uuid.UUID    `json:"id"       gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Interval IntervalType `json:"interval" gorm:"type:text;not null"`

	TimeOfDay *string `json:"timeOfDay"            gorm:"type:text;"`
	// only for WEEKLY
	Weekday *int `json:"weekday,omitempty"    gorm:"type:int"`
	// only for MONTHLY
	DayOfMonth *int `json:"dayOfMonth,omitempty" gorm:"type:int"`
}

func (i *Interval) BeforeSave(tx *gorm.DB) error {
	return i.Validate()
}

func (i *Interval) Validate() error {
	// for daily, weekly and monthly intervals time of day is required
	if (i.Interval == IntervalDaily || i.Interval == IntervalWeekly || i.Interval == IntervalMonthly) &&
		i.TimeOfDay == nil {
		return errors.New("time of day is required for daily, weekly and monthly intervals")
	}

	// for weekly interval weekday is required
	if i.Interval == IntervalWeekly && i.Weekday == nil {
		return errors.New("weekday is required for weekly intervals")
	}

	// for monthly interval day of month is required
	if i.Interval == IntervalMonthly && i.DayOfMonth == nil {
		return errors.New("day of month is required for monthly intervals")
	}

	return nil
}

// ShouldTriggerBackup checks if a backup should be triggered based on the interval and last backup time
func (i *Interval) ShouldTriggerBackup(now time.Time, lastBackupTime *time.Time) bool {
	// If no backup has been made yet, trigger immediately
	if lastBackupTime == nil {
		return true
	}

	switch i.Interval {
	case IntervalHourly:
		return now.Sub(*lastBackupTime) >= time.Hour
	case IntervalDaily:
		return i.shouldTriggerDaily(now, *lastBackupTime)
	case IntervalWeekly:
		return i.shouldTriggerWeekly(now, *lastBackupTime)
	case IntervalMonthly:
		return i.shouldTriggerMonthly(now, *lastBackupTime)
	default:
		return false
	}
}

// daily trigger: calendar-based if TimeOfDay set, otherwise next calendar day
func (i *Interval) shouldTriggerDaily(now, lastBackup time.Time) bool {
	if i.TimeOfDay != nil {
		target, err := time.Parse("15:04", *i.TimeOfDay)
		if err == nil {
			todayTarget := time.Date(
				now.Year(),
				now.Month(),
				now.Day(),
				target.Hour(),
				target.Minute(),
				0,
				0,
				now.Location(),
			)

			// if it's past today's target time and we haven't backed up today
			if now.After(todayTarget) && !isSameDay(lastBackup, now) {
				return true
			}

			// if it's exactly the target time and we haven't backed up today
			if now.Equal(todayTarget) && !isSameDay(lastBackup, now) {
				return true
			}

			// if it's before today's target time, don't trigger yet
			if now.Before(todayTarget) {
				return false
			}
		}
	}
	// no TimeOfDay: if it's a new calendar day
	return !isSameDay(lastBackup, now)
}

// weekly trigger: on specified weekday/calendar week, otherwise ≥7 days
func (i *Interval) shouldTriggerWeekly(now, lastBackup time.Time) bool {
	if i.Weekday != nil {
		targetWd := time.Weekday(*i.Weekday)
		startOfWeek := getStartOfWeek(now)

		// today is target weekday and no backup this week
		if now.Weekday() == targetWd && lastBackup.Before(startOfWeek) {
			if i.TimeOfDay != nil {
				t, err := time.Parse("15:04", *i.TimeOfDay)
				if err == nil {
					todayT := time.Date(
						now.Year(),
						now.Month(),
						now.Day(),
						t.Hour(),
						t.Minute(),
						0,
						0,
						now.Location(),
					)
					return now.After(todayT) || now.Equal(todayT)
				}
			}
			return true
		}
		// passed this week's slot and missed entirely
		targetThisWeek := startOfWeek.AddDate(0, 0, int(targetWd))
		if now.After(targetThisWeek) && lastBackup.Before(startOfWeek) {
			return true
		}
		return false
	}
	// no Weekday: generic 7-day interval
	return now.Sub(lastBackup) >= 7*24*time.Hour
}

// monthly trigger: on specified day/calendar month, otherwise next calendar month
func (i *Interval) shouldTriggerMonthly(now, lastBackup time.Time) bool {
	if i.DayOfMonth != nil {
		day := *i.DayOfMonth
		startOfMonth := getStartOfMonth(now)

		// today is target day and no backup this month
		if now.Day() == day && lastBackup.Before(startOfMonth) {
			if i.TimeOfDay != nil {
				t, err := time.Parse("15:04", *i.TimeOfDay)
				if err == nil {
					todayT := time.Date(
						now.Year(),
						now.Month(),
						now.Day(),
						t.Hour(),
						t.Minute(),
						0,
						0,
						now.Location(),
					)
					return now.After(todayT) || now.Equal(todayT)
				}
			}
			return true
		}
		// passed this month's slot and missed entirely
		if now.Day() > day && lastBackup.Before(startOfMonth) {
			return true
		}
		return false
	}
	// no DayOfMonth: if we're in a new calendar month
	return lastBackup.Before(getStartOfMonth(now))
}

func isSameDay(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func getStartOfWeek(t time.Time) time.Time {
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	return time.Date(t.Year(), t.Month(), t.Day()-wd+1, 0, 0, 0, 0, t.Location())
}

func getStartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}
