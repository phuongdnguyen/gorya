package worker

import (
	"context"
	"encoding/json"

	"github.com/nduyphuong/gorya/internal/queue"
	queueOptions "github.com/nduyphuong/gorya/internal/queue/options"
)

//go:generate mockery --name Interface --output fake --outpkg fake
type Interface interface {
	Process(ctx context.Context, resultChan chan<- string, errChan chan<- error)
	Dispatch(ctx context.Context, e *QueueElem) error
}

type client struct {
	queue queue.Interface
}

type Options struct {
	QueueOpts queueOptions.Options
}

type QueueElem struct {
	RequestURI string `json:"relative_uri"`
	Project    string `json:"project"`
	TagKey     string `json:"tagkey"`
	TagValue   string `json:"tagvalue"`
	Action     int    `json:"action"`
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

func (c *client) Process(ctx context.Context, resultChan chan<- string, errChan chan<- error) {
	c.queue.Dequeue(ctx, resultChan, errChan)
}
