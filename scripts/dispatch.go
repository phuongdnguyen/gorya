//go:build ignore
// +build ignore

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/nduyphuong/gorya/internal/constants"
	queueOpts "github.com/nduyphuong/gorya/internal/queue/options"
	"github.com/nduyphuong/gorya/internal/worker"
)

func main() {
	dispatch := flag.Bool("dispatch", true, "run dispatcher")
	process := flag.Bool("process", true, "run processor")
	queueName := flag.String("queueName", "", "queue name")
	action := flag.Int("action", 0, "start/stop action")
	flag.Parse()
	if *queueName == "" {
		fmt.Println("nothing to do")
		return
	}
	ctx := context.TODO()
	w := worker.NewClient(worker.Options{
		QueueOpts: queueOpts.Options{
			Addr:        "localhost:6379",
			Name:        *queueName,
			PopInterval: 2 * time.Second,
		},
	})
	numWorkers := 10
	var input []worker.QueueElem
	awsAssumeRoleItem := worker.QueueElem{
		Provider:      constants.PROVIDER_AWS,
		Project:       "test-aws-account",
		TagKey:        "foo",
		TagValue:      "bar",
		Action:        *action,
		CredentialRef: "arn:aws:iam::043159268388:role/test",
	}
	awsDefaultCredentialItem := worker.QueueElem{
		Provider: constants.PROVIDER_AWS,
		Project:  "test-aws-account",
		TagKey:   "foo",
		TagValue: "bar",
		Action:   *action,
	}
	gcpImpersonateItem := worker.QueueElem{
		Provider:      constants.PROVIDER_GCP,
		Project:       "target-project-397310.iam.gserviceaccount.com",
		TagKey:        "foo",
		TagValue:      "bar",
		Action:        *action,
		CredentialRef: "priv-sa@target-project-397310.iam.gserviceaccount.com",
	}
	input = append(input, awsAssumeRoleItem)
	input = append(input, awsDefaultCredentialItem)
	input = append(input, gcpImpersonateItem)
	if *dispatch {
		for _, v := range input {
			if err := w.Dispatch(ctx, &v); err != nil {
				fmt.Printf("err: %v\n", err)
			}
		}
		return
	}
	if *process {
		errCh := make(chan error)
		for i := 0; i < numWorkers; i++ {
			go w.Process(ctx, ctx.Done(), errCh)
		}
		var resErr error
		for err := range errCh {
			resErr = errors.Join(resErr, err)
		}
		fmt.Printf("errs: %v", resErr)
		return
	}

}
