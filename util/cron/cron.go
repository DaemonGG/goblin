package cron

import (
	"fmt"
	"log"
	"math"
	"time"

	"goblin/util/cron/scheduling"
	"goblin/util/prio"
)

const (
	Forever time.Duration = math.MaxInt64
)

type Function interface {
	// Define the method of this job.
	Run()
}

// A wrapper that turns a func() into a cron.Job
type FuncJob func()

func (f FuncJob) Run() { f() }

type Job struct {
	// The function to execute periodically.
	function Function

	// Represents how to schedule this function.
	scheduler scheduling.Schedule

	// The next time this function will be executed.
	next time.Time

	// The most recent time this function was executed.
	actual_prev time.Time

	// Index of this Job object in priority queue. Index -1 means this object is
	// not in priority queue. Default value is -1.
	index int
}

func (job_1 *Job) Less(job_2 prio.Interface) bool {
	job_2_ := job_2.(*Job)
	return job_1.next.Before(job_2_.next)
}

func (job *Job) Index(i int) {
	job.index = i
}

func (job *Job) Run() {
	job.actual_prev = time.Now()
	job.function.Run()
}

func (job *Job) NextRunTime() time.Time {
	return job.next
}

func (job *Job) setNextRunTime() {
	job.next = job.scheduler.Next()
}

type Cron struct {
	jobs      prio.Queue
	add_job   chan *Job
	stop      chan bool
	destroyed bool
}

func New() *Cron {
	c := &Cron{
		jobs:    prio.New(),
		add_job: make(chan *Job),
		stop:    make(chan bool),
	}
	go c.run()
	return c
}

func (cron *Cron) run() {
	for {
		now := time.Now()
		var timer *time.Timer
		if cron.jobs.Empty() {
			timer = time.NewTimer(Forever)
		} else {
			first_job := cron.jobs.Peek().(*Job)
			fmt.Println("Next time: ", first_job.NextRunTime())
			timer = time.NewTimer(first_job.NextRunTime().Sub(now))
		}

		for {
			select {
			case <-timer.C:
				for !cron.jobs.Peek().(*Job).NextRunTime().After(time.Now()) {
					one_job := cron.jobs.Pop().(*Job)
					if one_job.NextRunTime().IsZero() {
						log.Fatalln("Job scheduled without next time set.")
					}
					fmt.Println("Scheduled time: ", time.Now())
					go one_job.Run()
					one_job.setNextRunTime()
					cron.jobs.Push(one_job)
				}

			case new_job := <-cron.add_job:
				fmt.Println("Get new job, ", time.Now())
				new_job.setNextRunTime()
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
	if parsed_scheduler, err := scheduling.Parse(spec); err != nil {
		log.Fatalf("Error when parsing spec[%s]:%s\n", spec, err)
	} else {
		new_job := &Job{
			function:  FuncJob(f),
			scheduler: parsed_scheduler,
			index:     -1,
		}
		cron.add_job <- new_job
	}
}

func (cron *Cron) DestroySelf() {
	if !cron.destroyed {
		cron.stop <- true
		// TODO: Need to release memory in Queue
	}
}
