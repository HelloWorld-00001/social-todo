package subscribers

import (
	"context"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/coderconquerer/social-todo/common"
	asyncjob "github.com/coderconquerer/social-todo/common/async_job"
	"github.com/coderconquerer/social-todo/pubsub"
	"log"
)

// subJob represents a subscription job with a title and a handler function.
type subJob struct {
	Title string // Name or description of the job
	// SHandler is a handler function that processes incoming pubsub messages.
	SHandler func(ctx context.Context, message *pubsub.Message) error
}

// pbEngine is the main engine struct that holds the service context.
type pbEngine struct {
	serviceCtx goservice.ServiceContext // Shared service context for configuration, logging, etc.
}

// NewEngine creates and returns a new pbEngine instance using the provided service context.
func NewEngine(serviceCtx goservice.ServiceContext) *pbEngine {
	return &pbEngine{serviceCtx: serviceCtx}
}

// Start starts the pbEngine's operations.
// Currently, it does nothing and just returns nil, but this is likely a placeholder
// for initializing subscriptions, starting background tasks, or setting up listeners.
func (engine *pbEngine) Start() error {
	_ = engine.startSubTopic(common.TopicIncreaseTotalReact, true,
		IncreaseTotalReactionCount(engine.serviceCtx),
		NotifyUserReactTodoItem(engine.serviceCtx))
	_ = engine.startSubTopic(common.TopicDecreaseTotalReact, false,
		DecreaseTotalReactionCount(engine.serviceCtx))
	// todo: handle error with side jobs
	return nil
}

// GroupJob is an interface representing a group of jobs that can be run with a context
type GroupJob interface {
	Run(ctx context.Context) error
}

// startSubTopic sets up a subscriber for a given topic and starts processing messages
func (engine *pbEngine) startSubTopic(topic pubsub.Topic, isConcurrent bool, jobs ...subJob) error {
	// Subscribe to the topic using the pubsub system
	ps := engine.serviceCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)
	c, _ := ps.Subscribe(context.Background(), topic)

	// Log the setup of each job subscriber
	for _, item := range jobs {
		log.Printf("Setup subscriber for: %s", item.Title)
	}

	// Helper function to wrap job execution logic into an async job handler
	getJobHandler := func(job *subJob, message *pubsub.Message) asyncjob.JobHandler {
		return func(ctx context.Context) error {
			log.Printf("Running job for: %s, Value: %s", job.Title, message.Data())
			return job.SHandler(ctx, message) // Execute the job's handler
		}
	}

	// Goroutine to continuously listen for new messages and trigger jobs
	go func() {
		for {
			msg := <-c // Receive a new message from the subscription channel

			// Prepare an array of async jobs
			jobHdlArr := make([]asyncjob.Job, len(jobs))

			// Wrap each job with its message context
			for i := range jobs {
				jobHdl := getJobHandler(&jobs[i], msg)
				jobHdlArr[i] = asyncjob.NewJob(jobHdl, asyncjob.WithName(jobs[i].Title))
			}

			// Create a job group (can be concurrent or sequential)
			group := asyncjob.NewJobGroup(isConcurrent, jobHdlArr...)

			// Execute the group of jobs
			if err := group.Run(context.Background()); err != nil {
				log.Println(err)
			}
		}
	}()

	return nil
}
