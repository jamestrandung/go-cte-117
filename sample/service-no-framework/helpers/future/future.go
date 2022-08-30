package future

import (
	"context"
	"errors"
	"sync"
)

type LoaderFunc func()

func Load(loaders ...LoaderFunc) {
	wg := sync.WaitGroup{}

	for _, loader := range loaders {
		wg.Add(1)
		go func(l LoaderFunc) {
			defer wg.Done()

			l()
		}(loader)
	}

	wg.Wait()
}

// ErrNoResult is returned when there was no result
var ErrNoResult = errors.New("no result in future")

// Future - a shortcut for waiting on the result channel for a result or error.
type Future struct {
	parent   context.Context
	resultCh <-chan *Result

	done   sync.WaitGroup
	mtx    sync.RWMutex
	result *Result
}

// Result - the result of an async call - either value or error.
type Result struct {
	Value interface{}
	Err   error
}

// NewFuture - instantiate a new future. The result provider must close resultCh
// after providing a single result.
func NewFuture(parent context.Context, resultCh <-chan *Result) *Future {
	f := &Future{
		parent:   parent,
		resultCh: resultCh,
	}

	f.done.Add(1)
	go f.receive()

	return f
}

// NewFutureError instantiate a new future but have it immediately return as an error
func NewFutureError(parent context.Context, err error) *Future {
	f := &Future{
		parent: parent,
		result: &Result{Err: err},
	}

	return f
}

// Get - get the result or error. Blocks until the future is done or the
// context is cancelled.
func (f *Future) Get() (interface{}, error) {
	f.done.Wait()

	f.mtx.RLock()
	defer f.mtx.RUnlock()

	return f.result.Value, f.result.Err
}

// IsDone - is the future done.
func (f *Future) IsDone() bool {
	f.mtx.RLock()
	defer f.mtx.RUnlock()

	return f.result != nil
}

func (f *Future) receive() {
	// When receive() completes, the future will be marked as done.
	// Note that this will happen after the result is set and the
	// lock released.
	defer f.done.Done()

	var result *Result
	select {
	case <-f.parent.Done():
		result = &Result{Err: f.parent.Err()}

	case result = <-f.resultCh:
		if result == nil {
			result = &Result{Err: ErrNoResult}
		}
	}

	f.mtx.Lock()
	defer f.mtx.Unlock()

	f.result = result
}
