package main

import (
	"context"
	"errors"
	"fmt"
	asyncjob "github.com/coderconquerer/social-todo/common/async_job"
	"log"
)

func main() {
	job1 := asyncjob.NewJob(func(ctx context.Context) error {
		fmt.Println("I am job 1")
		return errors.New("something went wrong at job1")
	}, asyncjob.WithName("he he he"))

	if err := job1.Execute(context.Background()); err != nil {
		log.Println(err)
	}

	for {
		err := job1.Retry(context.Background())

		if err == nil || job1.State() == asyncjob.StateRetryFailed {
			break
		}
	}
}
