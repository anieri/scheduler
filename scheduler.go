package scheduler

import (
	"time"

	"go.uber.org/ratelimit"
)

type ScheduleJob interface {
	Run(time.Time) error
	Name() string
}

type schedule struct {
	job      ScheduleJob
	interval time.Duration
	last     time.Time
}

type Scheduler struct {
	jobs []*schedule
}

func New() *Scheduler {
	return &Scheduler{}
}

func (sched *Scheduler) Add(interval time.Duration, job ScheduleJob) {
	(*sched).jobs = append(sched.jobs, &schedule{job, interval, time.Now()})
}

func (sched *Scheduler) Run() {
	for _, job := range sched.jobs {
		(*job).last = time.Now()
	}
	rl := ratelimit.New(1)
	for {
		now := rl.Take()

		for _, job := range sched.jobs {
			if now.Before(job.last.Add(job.interval)) {
				continue
			}
			go job.job.Run(now)
		}
	}
}

func (sched *Scheduler) Stop() {
}
