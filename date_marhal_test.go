package civil

import (
	"encoding/json"
	"testing"
)

func TestDateMarshallJSON(t *testing.T) {
	tests := []Date{
		NewDate(2000, 1, 1),
		{},
	}
	for _, date := range tests {
		type Obj struct {
			Date Date
		}
		obj := Obj{Date: date}
		b, err := json.Marshal(obj)
		if err != nil {
			t.Errorf("Error marshalling %v: %v", date, err)
		}
		var got Obj
		if err := json.Unmarshal(b, &got); err != nil {
			t.Fatalf("Error unmarshalling %v: %v", date, err)
		}
		if !got.Date.Equal(date) {
			t.Errorf("Wrong value marshalled/unmarshalled for %v. Got %v", date, got.Date)
		}
	}
}
