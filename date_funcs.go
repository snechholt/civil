package civil

func MinDate(dates ...Date) Date {
	if len(dates) == 0 {
		return Date{}
	}
	min := dates[0]
	for _, d := range dates[1:] {
		if d.Before(min) {
			min = d
		}
	}
	return min
}

func MaxDate(dates ...Date) Date {
	if len(dates) == 0 {
		return Date{}
	}
	min := dates[0]
	for _, d := range dates[1:] {
		if d.After(min) {
			min = d
		}
	}
	return min
}
