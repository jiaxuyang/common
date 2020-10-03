package parallel

import (
	"container/list"
	"context"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

// Pool
type Pool struct {
	options  *Options
	queue    *list.List
	weighted *semaphore.Weighted
	wg       sync.WaitGroup
	lock     sync.Mutex
	err      error
	errOnce  sync.Once
	// TODO: implement cancel based on context, refer to errgroup
}

// NewPool create new pool
func NewPool(options *Options) *Pool {
	queue := list.New()
	weighted := semaphore.NewWeighted(int64(options.MaxPoolSize))
	pool := &Pool{
		options:  options,
		queue:    queue,
		weighted: weighted,
	}
	go pool.schedule()
	return pool
}

func (p *Pool) schedule() {
	for {
		for p.queue.Len() <= 0 {
			time.Sleep(time.Millisecond)
		}
		_ = p.weighted.Acquire(context.Background(), 1)
		f := p.dequeue()
		go func() {
			defer func() {
				p.weighted.Release(1)
				p.wg.Done()
			}()
			if err := f(); err != nil {
				p.errOnce.Do(func() {
					p.err = err
				})
			}
		}()
	}
}

func (p *Pool) dequeue() func() error {
	p.lock.Lock()
	defer p.lock.Unlock()
	item := p.queue.Front()
	p.queue.Remove(item)
	f := item.Value.(func() error)
	return f
}

// Go feed a closure without error return
func (p *Pool) Go(f func()) {
	p.Feed(func() error {
		f()
		return nil
	})
}

// Feed a closure, non-block.
func (p *Pool) Feed(f func() error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.queue.PushBack(f)
	p.wg.Add(1)
}

// Wait wait for all closures to finish.
func (p *Pool) Wait() {
	p.wg.Wait()
}

// Err wait and get first error
func (p *Pool) Err() error {
	p.Wait()
	return p.err
}
