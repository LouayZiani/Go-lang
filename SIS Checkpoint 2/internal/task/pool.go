package task

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// defaults used by the demo. NewPool takes explicit values.
const (
	DefaultWorkers       = 4
	DefaultQueueCapacity = 16
	minWorkers           = 2
)

type Pool struct {
	svc      *Service
	queue    chan string
	workers  int
	wg       sync.WaitGroup
	running  atomic.Int32 // worker goroutines started and not yet exited
	mu       sync.RWMutex // guards `stopped` and, above all, "send on queue" vs "close(queue)"
	stopped  bool
	capacity int
}

func NewPool(svc *Service, workers, queueCapacity int) (*Pool, error) {
	if svc == nil {
		return nil, fmt.Errorf("%w: service is nil", ErrInvalidConfig)
	}
	if workers < minWorkers {
		return nil, fmt.Errorf("%w: need at least %d workers, got %d", ErrInvalidConfig, minWorkers, workers)
	}
	if queueCapacity < 1 {
		return nil, fmt.Errorf("%w: queue capacity must be at least 1, got %d", ErrInvalidConfig, queueCapacity)
	}

	p := &Pool{
		svc:      svc,
		queue:    make(chan string, queueCapacity),
		workers:  workers,
		capacity: queueCapacity,
	}
	p.wg.Add(workers)
	for i := 0; i < workers; i++ {
		p.running.Add(1)
		go p.worker()
	}
	return p, nil
}

func (p *Pool) worker() {
	defer p.wg.Done()
	defer p.running.Add(-1)
	for id := range p.queue {
		p.svc.runQueued(id)
	}
}

func (p *Pool) Submit(id string) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.stopped {
		return ErrPoolStopped
	}

	if _, err := p.svc.move(id, Pending, Queued); err != nil {
		return err
	}

	select {
	case p.queue <- id:
		return nil
	default:
		if _, err := p.svc.move(id, Queued, Pending); err != nil {
			return fmt.Errorf("rolling back after full queue: %w", err)
		}
		return ErrQueueFull
	}
}

func (p *Pool) Stop() {
	p.mu.Lock()
	if !p.stopped {
		p.stopped = true
		close(p.queue)
	}
	p.mu.Unlock()

	p.wg.Wait()
}

func (p *Pool) Stopped() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.stopped
}

func (p *Pool) Workers() int { return p.workers }

func (p *Pool) QueueCapacity() int { return p.capacity }

func (p *Pool) QueueLen() int { return len(p.queue) }

func (p *Pool) Running() int { return int(p.running.Load()) }
