# Scheduler

[![GoDoc](https://godoc.org/github.com/anieri/scheduler?status.png)](https://godoc.org/github.com/anieri/scheduler)

Really small library to help you run tasks at predefined rates.

## Getting started

Implement your tasks as the following interface:

```go
type ScheduleJob interface {
    Run(time.Time) error
    Name() string
}
```

Then configure your tasks:

```go
logger, _ := zap.NewProduction()
jobScheduler := scheduler.New(logger.Named("scheduler"))
// Run this task every minute
jobScheduler.Add(time.Minute, yourTask)
// Run this one every ten milliseconds
jobScheduler.Add(time.Millisecond*10, fastTask)

// Configure the scheduler to check for tasks 10 times a second.
go jobScheduler.Run(10)

// Stop will block until all remaining tasks are completed
defer jobScheduler.Stop()
```

### Also check out

- [rakanalh/scheduler](https://github.com/rakanalh/scheduler).
