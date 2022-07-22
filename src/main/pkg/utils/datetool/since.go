package datetool

import "time"

// Since returns ms from t to now
func Since(t time.Time) int64 {
	return time.Since(t).Nanoseconds() / 1e6
}
