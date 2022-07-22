package number

import "time"

func DurationMs(duration time.Duration) int64 {
	return duration.Nanoseconds() / 1e6
}
