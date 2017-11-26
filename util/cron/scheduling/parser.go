package scheduling

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Schedule interface {
	// Given the current time, return the next run time.
	Next() time.Time

	// Provide the first execution time.
	// After() time.Duration
}

const (
	SpecSchedulerStrPrefix = "@every"
)

type specScheduler struct {
	// The interval time between every scheduled time.
	interval time.Duration

	// Whether this scheduler has been used or not.
	started bool

	// The last scheduled time.
	prev time.Time

	// The previous time get from time.Now().
	prev_now time.Time
}

func (spec *specScheduler) Next() time.Time {
	if spec.interval == 0 {
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

func parseSpecScheduler(duration_str string) (*specScheduler, error) {
	if len(duration_str) == 0 {
		return nil, errors.New(`No duration info is provided. Unable to create a
                           scheduler`)
	}
	// ParseDuration parses a duration string.
	// A duration string is a possibly signed sequence of
	// decimal numbers, each with optional fraction and a unit suffix,
	// such as "300ms", "-1.5h" or "2h45m".
	// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
	if d, err := time.ParseDuration(duration_str); err != nil {
		return nil, fmt.Errorf("Un-supported duration: %s", duration_str)
	} else {
		return &specScheduler{
			interval: d,
			started:  false,
		}, nil
	}
}

func Parse(spec string) (Schedule, error) {
	if len(spec) == 0 {
		return nil, errors.New(`No schedule string provided. Unable to create a
                           scheduler`)
	}
	str_splits := strings.Split(spec, " ")
	descriptor := str_splits[0]

	switch descriptor {
	case SpecSchedulerStrPrefix:
		return parseSpecScheduler(str_splits[1])
	default:
		return nil, errors.New(`No valid descriptor provided. Unable to create a
                           scheduler`)
	}
}
