package salt

import (
	"strings"
	"time"
)

type saltTime struct {
	time.Time
}

func (t *saltTime) UnmarshalJSON(input []byte) error {
	s := string(input)
	s = strings.Trim(s, "\"")
	v, err := time.Parse("2006, Jan 02 15:04:05.000000", s)
	if err != nil {
		return err
	}
	t.Time = v
	return nil
}
