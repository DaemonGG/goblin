package cron

import (
	"log"
	"strings"
	"time"
)

type SpecScheduler struct {
	// The interval time between every scheduled time.
	interval time.Duration

	// Whether this scheduler has been used or not.
	started bool

	// The last scheduled time.
	prev time.Time

	// The previous time get from time.Now().
	prev_now time.Time
}

func (spec *SpecScheduler) Next() time.Time {
	if interval == 0 {
		log.Fatalln("Interval of SpecScheduler cannot be 0.")
	}
	now := time.Now()
	if !spec.started {
		// If it's the first time to schedule a function, call it immediately.
		spec.started = true
		spec.prev = now
	} else if spec.prev_now.After(now) {
		// If the system time drifted back too much. The functions in this cron
		// might not be scheduled for a long while. In this case, schedule these
		// functions after the next duration.
		spec.prev = now.Add(spec.interval)
	} else if spec.prev.Add(spec.interval).Before(now) {
		spec.prev = now
	} else {
		spec.prev = spec.prev.Add(spec.interval)
	}
	spec.prev_now = now
	return spec.prev
}

func parse(spec string) (Schedule, error) {

}
