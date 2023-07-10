package civil

import (
	"fmt"
)

func (d Date) MarshalJSON() ([]byte, error) {
	s := `"` + d.Encode() + `"`
	return []byte(s), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
	n := len(b)
	if len(b) == 0 || b[0] != '"' || b[n-1] != '"' {
		return fmt.Errorf("invalid date string: %s", string(b))
	}
	s := string(b[1 : n-1])
	return d.Decode(s)
}
