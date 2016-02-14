package server

import (
	"testing"
	"time"
)

func TestToDate(t *testing.T) {
	a := toDate(time.Unix(0, 0))
	e := "19700101"
	if a != e {
		t.Errorf("got %v\nwant %v", a, e)
	}
}
