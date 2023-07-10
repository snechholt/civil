package civil

import (
	"testing"
	"time"
)

func TestNewDate(t *testing.T) {
	type yearMonthDay struct {
		year  int
		month time.Month
		day   int
	}
	ymd := func(year int, month time.Month, day int) yearMonthDay {
		return yearMonthDay{year: year, month: month, day: day}
	}
	tests := []struct {
		in   yearMonthDay
		want yearMonthDay
	}{
		{in: ymd(2000, 1, 1), want: ymd(2000, 1, 1)},
		{in: ymd(2000, 1, 32), want: ymd(2000, 2, 1)}, // Handles overflowing day
		{in: ymd(2000, 13, 1), want: ymd(2001, 1, 1)}, // Handles overflowing month
	}
	for _, test := range tests {
		date := NewDate(test.in.year, test.in.month, test.in.day)
		year, month, day := date.Date()
		if year != test.want.year || month != test.want.month || day != test.want.day {
			t.Errorf("NewDate(%v, %v, %v).Date() = %v, %v, %v, want %v, %v, %v",
				test.in.year, test.in.month, test.in.day, year, month, day, test.want.year, test.want.month, test.want.day)
		}
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		date Date
		want bool
	}{
		{date: NewDate(2000, 1, 1), want: false},
		{date: NewDate(0001, 1, 1), want: true},
		{date: DateFromTime(time.Time{}), want: true},
		{date: Date{}, want: true},
	}
	for _, test := range tests {
		if got := test.date.IsZero(); got != test.want {
			t.Errorf("Date(%v).IsZero() = %v, want %v", test.date, got, test.want)
		}
	}
}

func TestDateSub(t *testing.T) {
	tests := []struct {
		date1 Date
		date2 Date
		want  int
	}{
		{date1: NewDate(2000, 1, 1), date2: NewDate(2000, 1, 1), want: 0},
		{date1: NewDate(2000, 1, 2), date2: NewDate(2000, 1, 1), want: 1},
		{date1: NewDate(2000, 2, 1), date2: NewDate(2000, 1, 1), want: 31},
		{date1: NewDate(2001, 1, 1), date2: NewDate(2000, 1, 1), want: 366}, // 366 days due to leap year
	}
	for _, test := range tests {
		if want, got := test.want, test.date1.Sub(test.date2); got != want {
			t.Errorf("(%s).Sub(%s) = %d, want %d", test.date1, test.date2, got, want)
		}
		if want, got := -test.want, test.date2.Sub(test.date1); got != want {
			t.Errorf("(%s).Sub(%s) = %d, want %d", test.date2, test.date1, got, want)
		}
	}
}
