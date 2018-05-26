package foo

import "time"

func SleepForMinute() {
	// This is obviously wrong code, linter should throw a warning.
	time.Sleep(60)
}
