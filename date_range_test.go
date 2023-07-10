package civil

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewDateRangeDuration(t *testing.T) {
	tests := []struct {
		date Date
		dur  int
		unit Unit
		want DateRange
	}{
		{date: NewDate(2000, 1, 1), dur: 1, unit: Day, want: NewDateRange(NewDate(2000, 1, 1), NewDate(2000, 1, 2))},
		{date: NewDate(2000, 1, 1), dur: 366, unit: Day, want: NewDateRange(NewDate(2000, 1, 1), NewDate(2001, 1, 1))}, // leap year
		{date: NewDate(2001, 1, 1), dur: 365, unit: Day, want: NewDateRange(NewDate(2001, 1, 1), NewDate(2002, 1, 1))}, // non-leap year
		{date: NewDate(2000, 1, 1), dur: 1, unit: Month, want: NewDateRange(NewDate(2000, 1, 1), NewDate(2000, 2, 1))},
		{date: NewDate(2000, 1, 1), dur: 12, unit: Month, want: NewDateRange(NewDate(2000, 1, 1), NewDate(2001, 1, 1))},
		{date: NewDate(2000, 1, 1), dur: 1, unit: Year, want: NewDateRange(NewDate(2000, 1, 1), NewDate(2001, 1, 1))},
		{date: NewDate(2000, 1, 1), dur: 2, unit: Year, want: NewDateRange(NewDate(2000, 1, 1), NewDate(2002, 1, 1))},
	}
	for _, test := range tests {
		got := NewDateRangeDuration(test.date, test.dur, test.unit)
		if !got.Equal(test.want) {
			t.Errorf("NewDateRangeDuration(%v, %v, %v) = %v, want %v",
				test.date, test.dur, test.unit, got, test.want)
		}
	}
}

func TestDateRangeContains(t *testing.T) {
	dateRange := DateRange{start: NewDate(2000, 1, 5), end: NewDate(2000, 1, 7)}
	tests := map[Date]bool{
		NewDate(2000, 1, 1):  false,
		NewDate(2000, 1, 4):  false,
		NewDate(2000, 1, 5):  true,
		NewDate(2000, 1, 6):  true,
		NewDate(2000, 1, 7):  false,
		NewDate(2000, 1, 10): false,
	}
	for date, want := range tests {
		if got := dateRange.Contains(date); got != want {
			t.Errorf("Contains(%v) = %v, want %v", date, got, want)
		}
	}
}

func TestDateRangeIntersects(t *testing.T) {
	dateRange := NewDateRange(NewDate(2000, 1, 5), NewDate(2000, 1, 7))
	tests := map[DateRange]bool{
		NewDateRange(NewDate(2000, 1, 3), NewDate(2000, 1, 4)): false, // end before
		NewDateRange(NewDate(2000, 1, 4), NewDate(2000, 1, 5)): false, // end before
		NewDateRange(NewDate(2000, 1, 4), NewDate(2000, 1, 6)): true,  // start before, end inside
		NewDateRange(NewDate(2000, 1, 5), NewDate(2000, 1, 6)): true,  // inside
		NewDateRange(NewDate(2000, 1, 6), NewDate(2000, 1, 7)): true,  // inside
		NewDateRange(NewDate(2000, 1, 6), NewDate(2000, 1, 8)): true,  // start inside, end after
		NewDateRange(NewDate(2000, 1, 7), NewDate(2000, 1, 8)): false, // start after
		NewDateRange(NewDate(2000, 1, 8), NewDate(2000, 1, 9)): false, // start after
		NewDateRange(NewDate(2000, 1, 1), NewDate(2000, 1, 9)): true,  // start before, end after
	}
	for otherRange, want := range tests {
		if got := dateRange.Intersects(otherRange); got != want {
			t.Errorf("%v.Insersect(%v) = %v, want %v", dateRange, otherRange, got, want)
		}
		if got := otherRange.Intersects(dateRange); got != want {
			t.Errorf("%v.Insersect(%v) = %v, want %v", otherRange, dateRange, got, want)
		}
	}
}

func TestDateRangeIsZero(t *testing.T) {
	tests := map[DateRange]bool{
		DateRange{}: true,
		NewDateRange(Date{}, NewDate(2000, 1, 1)):              false,
		NewDateRange(NewDate(2000, 1, 1), Date{}):              false,
		NewDateRange(NewDate(2000, 1, 1), NewDate(2000, 1, 1)): false,
		NewDateRange(NewDate(2000, 1, 1), NewDate(2000, 1, 2)): false,
	}
	for dateRange, want := range tests {
		if got := dateRange.IsZero(); got != want {
			t.Errorf("%s.IsZero() = %v, want %v", dateRange, got, want)
		}
	}
}

func TestDateRangeEqual(t *testing.T) {
	dateRange := NewDateRange(NewDate(2015, 1, 5), NewDate(2015, 1, 7))
	tests := map[DateRange]bool{
		DateRange{}: false,
		NewDateRange(Date{}, NewDate(2000, 1, 1)):              false,
		NewDateRange(NewDate(2000, 1, 1), Date{}):              false,
		NewDateRange(NewDate(2000, 1, 1), NewDate(2000, 1, 1)): false,
		NewDateRange(NewDate(2000, 1, 1), NewDate(2000, 1, 2)): false,
		dateRange: true,
		NewDateRange(dateRange.start, dateRange.end): true,
	}
	for otherDateRange, want := range tests {
		if got := dateRange.Equal(otherDateRange); got != want {
			t.Errorf("%s.Equal(%s) = %v, want %v", dateRange, otherDateRange, got, want)
		}
	}
}

func TestDateRangeForEach(t *testing.T) {
	date, dr := NewDate, NewDateRange
	tests := []struct {
		dateRange DateRange
		unit      Unit
		want      []Date
	}{
		{
			dateRange: dr(date(2010, 1, 1), date(2010, 1, 4)),
			unit:      Day,
			want:      []Date{date(2010, 1, 1), date(2010, 1, 2), date(2010, 1, 3)},
		},
	}
	for _, test := range tests {
		var got []Date
		fn := func(d Date) { got = append(got, d) }
		test.dateRange.ForEach(test.unit, fn)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s.ForEach(%s) iterated over wrong dates\nWant %v\nGot  %v", test.dateRange, test.unit, test.want, got)
		}
	}
}

func TestDateRangeCount(t *testing.T) {
	date, dr := NewDate, NewDateRange
	tests := []struct {
		dateRange DateRange
		unit      Unit
		want      int
	}{
		// Count days
		{dateRange: dr(date(2010, 1, 1), date(2010, 1, 4)), unit: Day, want: 3},
		{dateRange: dr(date(2010, 1, 4), date(2010, 1, 1)), unit: Day, want: -3},
		{dateRange: dr(date(2010, 1, 1), date(2010, 1, 1)), unit: Day, want: 0},
		{dateRange: dr(date(2015, 1, 1), date(2016, 1, 1)), unit: Day, want: 365}, // non-leap year
		{dateRange: dr(date(2016, 1, 1), date(2017, 1, 1)), unit: Day, want: 366}, // leap year
	}
	for _, test := range tests {
		got := test.dateRange.Count(test.unit)
		if got != test.want {
			t.Errorf("%s.Count(%s) = %v, want %v", test.dateRange, test.unit, got, test.want)
		}
	}
}

func TestDateRangeSplit(t *testing.T) {
	dateRange := NewDateRange(NewDate(2000, 1, 1), NewDate(2000, 1, 4))
	tests := map[Date][2]DateRange{
		NewDate(2000, 1, 1): [2]DateRange{
			DateRange{},
			NewDateRange(NewDate(2000, 1, 1), NewDate(2000, 1, 4)),
		},
		NewDate(2000, 1, 2): [2]DateRange{
			NewDateRange(NewDate(2000, 1, 1), NewDate(2000, 1, 2)),
			NewDateRange(NewDate(2000, 1, 2), NewDate(2000, 1, 4)),
		},
		NewDate(2000, 1, 3): [2]DateRange{
			NewDateRange(NewDate(2000, 1, 1), NewDate(2000, 1, 3)),
			NewDateRange(NewDate(2000, 1, 3), NewDate(2000, 1, 4)),
		},
		NewDate(2000, 1, 4): [2]DateRange{
			NewDateRange(NewDate(2000, 1, 1), NewDate(2000, 1, 4)),
			DateRange{},
		},
	}
	for time, want := range tests {
		got := dateRange.Split(time)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s.Split(): wrong return value\nWant %v\nGot  %v", time, want, got)
		}
	}
}

func TestDateRangeEncoding(t *testing.T) {
	dateRange := NewDateRange(NewDate(2010, 1, 1), NewDate(2010, 1, 4))
	enc := dateRange.Encode()
	var got DateRange
	if err := got.Decode(enc); err != nil {
		t.Fatal(err)
	}
	if !got.Equal(dateRange) {
		t.Errorf("Wrong date range after encode/decode\nWant %v\nGot  %v", dateRange, got)
	}
}

func TestDateRangeMarshallJSON(t *testing.T) {
	tests := []DateRange{
		NewDateRange(NewDate(2000, 1, 1), NewDate(2010, 2, 1)),
		DateRange{},
	}
	for _, dateRange := range tests {
		type Obj struct {
			DateRange DateRange
		}
		obj := Obj{DateRange: dateRange}
		b, err := json.Marshal(obj)
		if err != nil {
			t.Errorf("Error marshalling %v: %v", dateRange, err)
		}
		var got Obj
		if err := json.Unmarshal(b, &got); err != nil {
			t.Fatalf("Error unmarshalling %v: %v", dateRange, err)
		}
		if !got.DateRange.Equal(dateRange) {
			t.Errorf("Wrong value marshalled/unmarshalled for %v. Got %v", dateRange, got.DateRange)
		}
	}
}

