package scheduler

import (
	"time"

	"go.uber.org/atomic"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
)

type ScheduleJob interface {
	Run(time.Time) error
	Name() string
}

type schedule struct {
	job      ScheduleJob
	interval time.Duration
	last     time.Time

	running *atomic.Bool
	counter *atomic.Int64
}

type Scheduler struct {
	jobs    []*schedule
	stop    chan chan struct{}
	log     *zap.Logger
	counter *atomic.Int64
}

func New(logger *zap.Logger) *Scheduler {
	return &Scheduler{
		stop:    make(chan chan struct{}, 1),
		log:     logger,
		counter: atomic.NewInt64(0),
	}
}

func (sched *Scheduler) Add(interval time.Duration, job ScheduleJob) {
	(*sched).jobs = append(sched.jobs, &schedule{
		job:      job,
		interval: interval,
		last:     time.Now(),
		running:  atomic.NewBool(false),
		counter:  atomic.NewInt64(0),
	})
}

func (sched *Scheduler) Run() {
	sched.log.Info("scheduler.run")
	for _, job := range sched.jobs {
		job.last = time.Now()
	}

	rl := ratelimit.New(1)
	for {
		now := rl.Take()

		select {
		case stopMessage := <-sched.stop:
			stopMessage <- struct{}{}

		default:
			for _, job := range sched.jobs {
				if now.Before(job.last.Add(job.interval)) {
					continue
				}
				go sched.runJob(job, now)
			}
		}
	}
}

func (sched *Scheduler) runJob(job *schedule, now time.Time) {
	mainCounter := sched.counter.Inc()
	jobCounter := job.counter.Inc()
	name := job.job.Name()

	sched.log.Debug("job.run",
		zap.String("name", name),
		zap.Int64("counter.main", mainCounter),
		zap.Int64("counter.job", jobCounter),
	)

	job.last = now
	if job.running.Load() {
		sched.log.Warn("job.isRunning",
			zap.String("name", name),
			zap.Int64("counter.main", mainCounter),
			zap.Int64("counter.job", jobCounter),
		)
		return
	}
	job.running.Store(true)
	err := job.job.Run(now)
	job.running.Store(false)
	if err != nil {
		sched.log.Error("job.error",
			zap.String("name", name),
			zap.Int64("counter.main", mainCounter),
			zap.Int64("counter.job", jobCounter),
			zap.Error(err))
		return
	}
	sched.log.Info("job.ok",
		zap.String("name", name),
		zap.Int64("counter.main", mainCounter),
		zap.Int64("counter.job", jobCounter),
	)
}

func (sched *Scheduler) Stop() {
	sched.log.Info("scheduler.stop.request")
	stopMessage := make(chan struct{}, 1)
	sched.stop <- stopMessage
	<-stopMessage
	sched.log.Info("scheduler.stop.ok")
}
