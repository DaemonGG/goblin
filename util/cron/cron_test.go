package cron

import (
	"fmt"
	"testing"
	"time"
)

const (
	OneSecond  = time.Second + 2*time.Millisecond
	SmallWhile = 1000 * time.Nanosecond
)

func TestAddOneFunc(t *testing.T) {
	cron := New()
	defer cron.DestroySelf()
	var (
		num_call                          = 0
		expected_num_call                 = 0
		repeated_times_reach_test_success = 100
	)

	fmt.Println("Add one func at time: ", time.Now())
	cron.Add("@every 1s", func() { num_call += 1 })

	expected_num_call += 1
	<-time.After(SmallWhile)
	if num_call != expected_num_call {
		t.Fatalf("expected this function has been scheduled for one time. ("+
			"num_call=%d)", num_call)
	}
	for {
		expected_num_call += 1
		<-time.After(OneSecond)
		if num_call != expected_num_call {
			t.Fatalf("expected this function has been scheduled for %s times. "+
				"(num_call=%d)", expected_num_call, num_call)
		}
		if num_call == repeated_times_reach_test_success {
			fmt.Println("Success!")
			break
		}
	}
}
