package cron

import (
	"log"
	"math"
	"time"

	"goblin/util/prio"
)

const (
	Forever time.Duration = math.MaxInt64 * time.Hour
)

type Function interface {
	// Define the method of this job.
	Run()
}

// A wrapper that turns a func() into a cron.Job
type FuncJob func()

func (f FuncJob) Run() { f() }

type Schedule interface {
	// Given the current time, return the next run time.
	Next(time.Time) time.Time

	// Provide the first execution time.
	// After() time.Duration
}

type Job struct {
	// The function to execute periodically.
	function Function

	// Represents how to schedule this function.
	scheduler Schedule

	// The next time this function will be executed.
	next time.Time

	// The most recent time this function was executed.
	prev time.Time

	// Index of this Job object in priority queue. Index -1 means this object is
	// not in priority queue. Default value is -1.
	index int
}

func (job_1 *Job) Less(job_2 *Job) bool {
	return job_1.next.Before(job_2.next)
}

func (job *Job) Index(i int) bool {
	job.index = i
}

func (job *Job) Run() {
	job.function.Run()
}

func (job *Job) NextRunTime() time.Time {
	return job.next
}

func (job *Job) setNextRunTime() {
	job.next = job.scheduler.Next()
	job.prev = job.next
}

type Cron struct {
	jobs      Queue
	add_job   chan *Job
	stop      chan bool
	destroyed bool
}

func NewCron() *Cron {
	return &Cron{
		jobs:    prio.New([]Job{}),
		add_job: make(chan *Job),
		stop:    make(chan bool),
	}
}

func (cron *Cron) run() {
	for {
		now := time.Now()
		var timer *time.Timer
		if cron.jobs.Empty() {
			timer = time.NewTimer(Forever)
		} else {
			first_job := cron.jobs.Peek()
			timer = time.NewTimer(first_job.next.Sub(now))
		}

		for {
			select {
			case <-timer.C:
				for !cron.jobs.Peek().NextRunTime().After(now) {
					one_job := cron.jobs.Pop()
					if one_job.next.IsZero() {
						log.Fatalln("Job scheduled without next time set.")
					}
					go one_job.Run()
					one_job.setNextRunTime()
					jobs.Push(one_job)
				}

			case new_job := <-cron.add_job:
				new_job.setNextRunTime(now)
				cron.jobs.Push(new_job)

			case <-cron.stop:
				timer.Stop()
				return
			}
			break
		}
	}
}

func (cron *Cron) Add(spec string, f func()) {
	if cron.destroyed {
		log.Fatalln("Cron has been destroyed")
	}
	if parsed_scheduler, err := parse(spec); err {
		log.Fatalf("Error when parsing spec[%s]:%s\n", spec, err)
	}
	new_job := &Job{
		function:  FuncJob(f),
		scheduler: parsed_scheduler,
		index:     -1,
	}
	cron.add_job <- new_job
}

func (cron *Cron) DestroySelf() {
	c.stop <- true
	// TODO: Need to release memory in Queue
}
