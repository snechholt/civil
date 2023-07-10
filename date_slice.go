package civil

type DateSlice []Date

func (s DateSlice) Contains(d Date) bool {
	for _, item := range s {
		if d.Equal(item) {
			return true
		}
	}
	return false
}

func (s DateSlice) Filter(fn func(d Date) bool) DateSlice {
	var dates DateSlice
	for _, d := range s {
		if fn(d) {
			dates = append(dates, d)
		}
	}
	return dates
}

func (s DateSlice) Find(fn func(d Date) bool) Date {
	for _, d := range s {
		if fn(d) {
			return d
		}
	}
	return Date{}
}

func (s DateSlice) Some(fn func(d Date) bool) bool {
	for _, d := range s {
		if fn(d) {
			return true
		}
	}
	return false
}

func (s DateSlice) NonZero() DateSlice {
	return s.Filter(func(date Date) bool { return !date.IsZero() })
}

func (s DateSlice) Min() Date {
	return MinDate(s...)
}

func (s DateSlice) Max() Date {
	return MaxDate(s...)
}

func (s DateSlice) Len() int {
	return len(s)
}

func (s DateSlice) Less(i, j int) bool {
	return s[i].Before(s[j])
}

func (s DateSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
