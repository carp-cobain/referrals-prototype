package model

import "time"

// Time is used to store timestamps as INT in SQLite
type Time int64

func (t *Time) FromUnix() time.Time {
	return time.Unix(int64(*t), 0)
}
