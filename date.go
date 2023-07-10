// Package civil provides functionality for working with civil dates.
//
// A civil date is the date part of a timestamp without any timezone
// information attached. Types from this package aim to ease working
// with dates where timezone concerns are irrelevant and the standard
// library time package is too complex.

package civil

import (
	"time"
)

type Date struct {
	t time.Time
}

func NewDate(year int, month time.Month, day int) Date {
	t := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return Date{t: t}
}

func DateFromTime(t time.Time) Date {
	return NewDate(t.Date())
}

func TodayIn(loc *time.Location) Date {
	return DateFromTime(time.Now().In(loc))
}

func Today() Date {
	return TodayIn(time.Local)
}

func ParseDate(layout, value string) (Date, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return Date{}, err
	}
	return DateFromTime(t), nil
}

func (d Date) Year() int {
	return d.t.Year()
}

func (d Date) Month() time.Month {
	return d.t.Month()
}

func (d Date) Day() int {
	return d.t.Day()
}

func (d Date) Date() (year int, month time.Month, day int) {
	return d.t.Date()
}

func (d Date) Weekday() time.Weekday {
	return d.t.Weekday()
}

func (d Date) ISOWeek() (year, week int) {
	return d.t.ISOWeek()
}

func (d Date) AddDate(years, months, days int) Date {
	return DateFromTime(d.UTC().AddDate(years, months, days))
}

func (d Date) Sub(d2 Date) int {
	return NewDateRange(d2, d).Count(Day)
}

func (d Date) IsZero() bool {
	return d == Date{}
}

func (d Date) Before(d2 Date) bool {
	return d.UTC().Before(d2.UTC())
}

func (d Date) After(d2 Date) bool {
	return d.UTC().After(d2.UTC())
}

func (d Date) Equal(d2 Date) bool {
	return d.t.Equal(d2.t)
}

func (d Date) In(loc *time.Location) time.Time {
	year, month, day := d.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, loc)
}

func (d Date) Local() time.Time {
	return d.In(time.Local)
}

func (d Date) UTC() time.Time {
	return d.In(time.UTC)
}

func (d Date) Format(layout string) string {
	return d.Local().Format(layout)
}

func (d Date) String() string {
	return d.Format("2006-01-02")
}

func (d Date) Encode() string {
	return d.Format("2006-01-02")
}

func (d *Date) Decode(value string) error {
	parsed, err := ParseDate("2006-01-02", value)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}
