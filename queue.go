package robot

import (
	"errors"
	"time"

	"golang.org/x/time/rate"
	"k8s.io/client-go/util/workqueue"
)

type queue interface {
	// Push push an object in queue
	push(QueueObject)

	// Pop return an object in the queue, if queue is closed
	// it will return an error
	Pop() (QueueObject, error)

	// ReQueue indicates if an object fails to execute, he should be rejoined in the queue.
	// But it can`t be ReQueue for more than 3 times.
	ReQueue(QueueObject) error

	// Finish indicates that an object has been successfully processed.
	Finish(QueueObject)

	// Close will cause queue to ignore all new items added to it. As soon as the
	// worker goroutines have drained the existing items in the queue, they will be
	// instructed to exit.
	close()
}

type wq struct {
	workqueue.RateLimitingInterface
}

var _ queue = &wq{}

func newWorkQueue() *wq {
	return &wq{
		workqueue.NewRateLimitingQueue(workqueue.NewMaxOfRateLimiter(
			workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
			// 100 qps, 100 bucket size.  This is only for retry speed and its only the overall factor (not per item)
			&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(100), 100)},
		)),
	}
}

func (c *wq) push(obj QueueObject) {
	c.AddRateLimited(obj)
}

func (c *wq) Pop() (QueueObject, error) {
	obj, quit := c.Get()
	if quit {
		return QueueObject{}, errors.New("Controller has been stoped. ")
	}

	return obj.(QueueObject), nil
}

func (c *wq) Finish(obj QueueObject) {
	c.Forget(obj)
	c.Done(obj)
}

func (c *wq) ReQueue(obj QueueObject) error {
	if c.NumRequeues(obj) < 3 {
		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.AddRateLimited(obj)
		return nil
	}

	c.Forget(obj)

	c.Done(obj)

	return errors.New("This object has been requeued for many times, but still fails. ")
}

func (c *wq) close() {
	c.ShutDown()
}
