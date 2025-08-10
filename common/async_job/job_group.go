package async_job

import (
	"context"
	"github.com/coderconquerer/social-todo/common"
	"log"
	"sync"
)

type JobGroup struct {
	jobs          []Job
	isConcurrency bool
	waitGroup     *sync.WaitGroup
}

func NewJobGroup(isConcurrent bool, jobs ...Job) *JobGroup {
	return &JobGroup{
		jobs,
		isConcurrent,
		new(sync.WaitGroup),
	}
}

// Run executes all jobs in the JobGroup.
// If isConcurrency == true, jobs are run in parallel (concurrent execution).
// Otherwise, jobs are run sequentially (one after another).
func (jg *JobGroup) Run(ctx context.Context) error {
	// Add the total number of jobs to the WaitGroup counter.
	// This ensures we can wait until all jobs finish.
	jg.waitGroup.Add(len(jg.jobs))

	// Create an error channel to store the result of each job.
	// Buffered to len(jg.jobs) so sends won't block.
	errChan := make(chan error, len(jg.jobs))

	// Iterate over all jobs in the group
	for i := range jg.jobs {
		if jg.isConcurrency {
			// Run each job in its own goroutine for concurrent execution
			go func(asyncJob Job) {
				defer common.Recovery() // Catch and log any panic inside the goroutine

				// Run the job and send any error to errChan
				errChan <- jg.RunJob(ctx, asyncJob)

				// Signal that this job is done
				jg.waitGroup.Done()
			}(jg.jobs[i]) // Pass job as parameter to avoid loop variable capture
		} else {
			// Sequential execution: run the job immediately in the current goroutine
			toExecute := jg.jobs[i]

			// Execute the job
			err := jg.RunJob(ctx, toExecute)
			if err != nil {
				// Log error and stop execution — sequential mode exits early on failure
				log.Println(err)
				return err
			}

			// Send nil or error to the channel (in sequential case, this will always be nil here)
			errChan <- err

			// Mark this job as done
			jg.waitGroup.Done()
		}
	}

	// Wait for all jobs to finish before collecting results
	jg.waitGroup.Wait()

	// Check the collected errors
	var err error
	for i := 0; i < len(jg.jobs); i++ {
		if v := <-errChan; v != nil {
			// todo: tracks all the errs rather than return the first err
			return v
		}
	}

	// Return the last non-nil error encountered, or nil if all succeeded
	return err
}

// RunJob executes a single job with retry logic.
func (jg *JobGroup) RunJob(ctx context.Context, job Job) error {
	// Try executing the job
	if err := job.Execute(ctx); err != nil {
		// If the job failed, enter a retry loop
		for {
			log.Println(err)

			// Stop retrying if the job reached maximum retry attempts
			if job.State() == StateRetryFailed {
				return err
			}

			// Retry the job. If retry succeeds, exit with nil.
			if job.Retry(ctx) == nil {
				return nil
			}

			// If retry failed, loop again (and possibly fail later)
		}
	}

	// Job succeeded on the first try
	return nil
}
