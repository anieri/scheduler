package scheduler

import "time"

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
	jobs []schedule
}

func New() *Scheduler {
	return &Scheduler{}
}

func (sched *Scheduler) Add(interval time.Duration, job ScheduleJob) {
	(*sched).jobs = append(sched.jobs, schedule{job, interval, time.Now()})
}

func (sched *Scheduler) Run() {
}

func (sched *Scheduler) Stop() {
}
