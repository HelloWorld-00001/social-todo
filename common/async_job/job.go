package async_job

import (
	"context"
	"log"
	"time"
)

// Job requirement
// 1. Job can do sth (Handler)
// 2. Job can retry
// 2.1 Job can config retry times & duration
// 3. Job should be stateful: todo-explain
// 4. Need a job manager to handle & coordinate jobs
type Job interface {
	Execute(ctx context.Context) error
	Retry(ctx context.Context) error
	State() JobState
	SetRetryDuration(times []time.Duration)
}

const (
	DefaultMaxTimeout = time.Second * 10
)

var (
	defaultRetryTime = []time.Duration{time.Second, time.Second * 2, time.Second * 4}
)

type JobState int

type JobHandler func(ctx context.Context) error

const (
	StateInit JobState = iota
	StateRunning
	StateCompleted
	StateFailed
	StateRetryFailed
	StateTimeout
)

func (js JobState) String() string {
	return [6]string{"Init", "Running", "Completed", "Failed", "RetryFailed", "Timeout"}[js]
}

type JobConfig struct {
	Name       string
	MaxTimeout time.Duration
	Retries    []time.Duration
}

type job struct {
	jobConfig  JobConfig
	state      JobState
	handler    JobHandler
	retryIndex int
	stopChan   chan bool
}

func NewJob(handler JobHandler, opts ...OptionHdl) *job {
	j := job{
		jobConfig: JobConfig{
			MaxTimeout: DefaultMaxTimeout,
			Retries:    defaultRetryTime,
		},
		handler:    handler,
		retryIndex: -1,
		stopChan:   make(chan bool),
	}

	for i := range opts {
		opts[i](&j.jobConfig)
	}

	return &j
}

func (j *job) StateHasChange(state JobState) {
	j.state = state
}

type OptionHdl func(cfg *JobConfig)

func (j *job) Execute(ctx context.Context) error {
	log.Printf("start job %s", j.jobConfig.Name)
	j.StateHasChange(StateInit)

	var err error
	err = j.handler(ctx)
	if err != nil {
		j.StateHasChange(StateFailed)
		return err
	}

	j.StateHasChange(StateCompleted)
	return nil
}

func (j *job) Retry(ctx context.Context) error {
	if j.retryIndex >= len(j.jobConfig.Retries)-1 {
		j.StateHasChange(StateFailed)
		return nil // todo: find a way to return last failed err
	}
	j.retryIndex++
	time.Sleep(j.jobConfig.Retries[j.retryIndex])

	err := j.Execute(ctx)

	if err == nil {
		j.StateHasChange(StateCompleted)
		return nil
	}

	if j.retryIndex == len(j.jobConfig.Retries)-1 {
		j.StateHasChange(StateRetryFailed)
		return err
	}

	j.StateHasChange(StateFailed)
	return err
}

func (j *job) State() JobState {
	return j.state
}

func (j *job) RetryIndex() int {
	return j.retryIndex
}
func (j *job) SetRetryDuration(times []time.Duration) {
	if len(times) == 0 {
		return
	}
	j.jobConfig.Retries = times
}

func WithName(name string) OptionHdl {
	return func(cfg *JobConfig) {
		cfg.Name = name
	}
}

func WithMaxTimeout(maxTimeout time.Duration) OptionHdl {
	return func(cfg *JobConfig) {
		cfg.MaxTimeout = maxTimeout
	}
}

func WithRetries(retries ...time.Duration) OptionHdl {
	if len(retries) == 0 {
		retries = defaultRetryTime
	}
	return func(cfg *JobConfig) {
		cfg.Retries = append(cfg.Retries, retries...)
	}
}
