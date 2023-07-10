package civil

import (
	"encoding/json"
	"fmt"
)

type Unit int

const (
	Day Unit = iota + 1
	Month
	Year
)

func (u Unit) String() string {
	switch u {
	case Day:
		return "Day"
	case Month:
		return "Month"
	case Year:
		return "Year"
	default:
		return "<invalid unit>"
	}
}

type DateRange struct {
	start Date
	end   Date
}

func NewDateRange(start, end Date) DateRange {
	return DateRange{start: start, end: end}
}

func NewDateRangeDuration(start Date, duration int, unit Unit) DateRange {
	dr := DateRange{start: start}
	switch unit {
	case Day:
		dr.end = start.AddDate(0, 0, duration)
	case Month:
		dr.end = start.AddDate(0, duration, 0)
	case Year:
		dr.end = start.AddDate(duration, 0, 0)
	default:
		panic(fmt.Sprintf("NewDateRangeDuration: Unknown duration unit: '%v'", unit))
	}
	return dr
}

func (dr DateRange) Start() Date {
	return dr.start
}

func (dr DateRange) End() Date {
	return dr.end
}

func (dr DateRange) IsZero() bool {
	return dr == DateRange{}
}

func (dr DateRange) Equal(dr2 DateRange) bool {
	return dr.start.Equal(dr2.start) && dr.end.Equal(dr2.end)
}

func (dr DateRange) Contains(d Date) bool {
	return !dr.start.After(d) && d.Before(dr.end)
}

func (dr DateRange) Intersects(other DateRange) bool {
	return dr.Contains(other.start) || other.Contains(dr.start)
}

func (dr DateRange) Count(unit Unit) int {
	switch unit {
	case Day:
		return int(dr.end.UTC().Sub(dr.start.UTC()).Minutes() / 24 / 60)
	default:
		panic(fmt.Sprintf("DateRange.Count(%v) not implemented", unit))
	}
}

func (dr DateRange) ForEachIndexErr(unit Unit, fn func(date Date, index int) error) error {
	var years, months, days int
	switch unit {
	case Day:
		days = 1
	case Month:
		months = 1
	case Year:
		years = 1
	}
	var index int
	for date := dr.start; date.Before(dr.end); date = date.AddDate(years, months, days) {
		if err := fn(date, index); err != nil {
			return err
		}
		index++
	}
	return nil
}

func (dr DateRange) ForEachIndex(unit Unit, fn func(date Date, index int)) {
	_ = dr.ForEachIndexErr(unit, func(date Date, index int) error {
		fn(date, index)
		return nil
	})
}

func (dr DateRange) ForEachErr(unit Unit, fn func(date Date) error) error {
	return dr.ForEachIndexErr(unit, func(date Date, index int) error {
		return fn(date)
	})
}

func (dr DateRange) ForEach(unit Unit, fn func(date Date)) {
	_ = dr.ForEachIndexErr(unit, func(date Date, index int) error {
		fn(date)
		return nil
	})
}

func (dr DateRange) Split(d Date) [2]DateRange {
	split := [2]DateRange{}
	switch {
	case !d.After(dr.start):
		split[1] = dr
	case !d.Before(dr.end):
		split[0] = dr
	default:
		split[0] = NewDateRange(dr.start, d)
		split[1] = NewDateRange(d, dr.end)
	}
	return split
}

func (dr DateRange) Split2(d Date) (DateRange, DateRange) {
	switch {
	case !d.After(dr.start):
		return DateRange{}, dr
	case !d.Before(dr.end):
		return dr, DateRange{}
	default:
		return NewDateRange(dr.start, d), NewDateRange(d, dr.end)
	}
}

func (dr DateRange) ToSlice(unit Unit) DateSlice {
	var dates []Date
	dr.ForEach(unit, func(d Date) {
		dates = append(dates, d)
	})
	return dates
}

func (dr DateRange) Encode() string {
	return fmt.Sprintf("[%s,%s)", dr.start, dr.end)
}

func (dr *DateRange) Decode(src string) error {
	errInvalid := fmt.Errorf("DateRange.Decode(): invalid src: '%s'", src)
	if len(src) != 23 || src[0] != '[' || src[11] != ',' || src[22] != ')' {
		return errInvalid
	}
	var start, end Date
	if err := start.Decode(src[1:11]); err != nil {
		return errInvalid
	}
	if err := end.Decode(src[12:22]); err != nil {
		return errInvalid
	}
	dr.start = start
	dr.end = end
	return nil
}

func (dr DateRange) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`{"start":"%s","end":"%s"}`, dr.start.Encode(), dr.end.Encode())
	return []byte(s), nil
}

func (dr *DateRange) UnmarshalJSON(b []byte) error {
	var dst struct {
		Start Date `json:"start"`
		End   Date `json:"end"`
	}
	if err := json.Unmarshal(b, &dst); err != nil {
		return err
	}
	dr.start = dst.Start
	dr.end = dst.End
	return nil
}

func (dr DateRange) String() string {
	return fmt.Sprintf("[%s, %s>", dr.start, dr.end)
}
