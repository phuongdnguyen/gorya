package queue

import (
	"context"
	"time"

	libredis "github.com/go-redis/redis/v8"
	"github.com/nduyphuong/gorya/internal/queue/options"
)

// TODO: implement a more reliable queue

//go:generate mockery --name Interface
type Interface interface {
	Enqueue(ctx context.Context, v any) error
	Dequeue(ctx context.Context, out chan<- string, errChan chan<- error)
	IsEmpty(ctx context.Context) (bool, error)
}

type Client struct {
	Redis *libredis.Client
	opts  options.Options
}

func NewQueue(opts ...options.Option) Interface {
	return newQueue(opts...)
}

func newQueue(opts ...options.Option) *Client {
	c := &Client{}
	for _, o := range opts {
		o.Apply(&c.opts)
	}

	c.Redis = libredis.NewClient(&libredis.Options{
		Addr: c.opts.Addr,
	})
	return c
}

// Enqueue add element to the end of the queue
func (c *Client) Enqueue(ctx context.Context, v any) error {
	if err := c.Redis.LPush(ctx, c.opts.Name, v).Err(); err != nil {
		return err
	}
	return nil
}

// Dequeue pop head of the queue and return the popped value
func (c *Client) dequeue(ctx context.Context) (string, error) {
	resp, err := c.Redis.RPop(ctx, c.opts.Name).Result()
	if err != nil {
		if err != libredis.Nil {
			return "", err
		}
	}
	return resp, nil
}

func (c *Client) Dequeue(ctx context.Context, out chan<- string, errChan chan<- error) {
	ticker := time.NewTicker(c.opts.PopInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			resp, err := c.dequeue(ctx)
			if err != nil {
				errChan <- err
				return
			}
			if len(resp) != 0 {
				out <- resp
			}
		}
	}
}

func (c *Client) IsEmpty(ctx context.Context) (bool, error) {
	len, err := c.Redis.LLen(ctx, c.opts.Name).Result()
	if err != nil {
		return false, err
	}
	return len == 0, nil
}
