package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nduyphuong/gorya/internal/logging"
	"github.com/nduyphuong/gorya/internal/os"
	"github.com/nduyphuong/gorya/internal/queue"
	queueOptions "github.com/nduyphuong/gorya/internal/queue/options"
	"github.com/nduyphuong/gorya/internal/types"
	"github.com/nduyphuong/gorya/pkg/api/service/v1alpha1"
)

//go:generate mockery --name Interface
type Interface interface {
	// Process periodically check if there is any item in the queue that needs to be processed
	Process(ctx context.Context, stop <-chan struct{}, errChan chan<- error)
	Dispatch(ctx context.Context, e *QueueElem) error
}

type client struct {
	queue queue.Interface
}

type Options struct {
	QueueOpts queueOptions.Options
}

type QueueElem struct {
	Project       string `json:"project"`
	CredentialRef string `json:"credentialref"`
	TagKey        string `json:"tagkey"`
	TagValue      string `json:"tagvalue"`
	Action        int    `json:"action"`
	Provider      string `json:"provider"`
}

func NewClient(opts Options) Interface {
	c := &client{
		queue: queue.NewQueue(
			queueOptions.WithFetchInterval(opts.QueueOpts.PopInterval),
			queueOptions.WithQueueName(opts.QueueOpts.Name),
			queueOptions.WithQueueAddr(opts.QueueOpts.Addr),
		),
	}
	return c
}

func (c *client) Dispatch(ctx context.Context, e *QueueElem) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	if err := c.queue.Enqueue(ctx, b); err != nil {
		return err
	}
	return nil
}

// Process periodically check if there is any item in the queue that needs to be processed
func (c *client) Process(ctx context.Context, stop <-chan struct{}, errChan chan<- error) {
	log := logging.LoggerFromContext(ctx)
	resultChan := make(chan string)
	go c.queue.Dequeue(ctx, resultChan, errChan)

	for {
		select {
		case <-stop:
			fmt.Println("stop background process")
			return
		case task := <-resultChan:
			fmt.Printf("popped item %v \n", task)
			var elem QueueElem
			err := json.Unmarshal([]byte(task), &elem)
			if err != nil {
				log.Errorf("unmarshal elem from queue %v", err)
				return
			}
			changeStateRequest := v1alpha1.ChangeStateRequest{
				Action:        elem.Action,
				Project:       elem.Project,
				TagKey:        elem.TagKey,
				TagValue:      elem.TagValue,
				CredentialRef: elem.CredentialRef,
				Provider:      elem.Provider,
			}
			requestURL := fmt.Sprintf("http://localhost:%d%s", types.MustParseInt(os.GetEnv("PORT",
				"8080")), v1alpha1.GoryaTaskChangeStageProcedure)
			b, err := json.Marshal(changeStateRequest)
			if err != nil {
				log.Errorf("unmarshal changeStateRequest %v", err)
				return
			}
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(b))
			if err != nil {
				log.Errorf("creating request %v", err)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Errorf("making request %v", err)
			}
			if resp != nil {
				if resp.StatusCode != http.StatusOK {
					body, _ := io.ReadAll(resp.Body)
					log.Errorf("change state request failed with status code: %v and resp %v",
						resp.StatusCode, string(body))
				}
			}
		}
	}
}
