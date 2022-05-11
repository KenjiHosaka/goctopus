package goctopus

import (
	"context"
	"math/rand"
	"time"
)

type taskResult struct {
	taskID uint32
	value  any
}

type AsyncTask func(ctx context.Context) (taskResult, error)

type OrchestrationFunc func(...Option) (Outputs, error)

func Orchestrate(ctx context.Context, tasks ...AsyncTask) OrchestrationFunc {
	return func(opts ...Option) (Outputs, error) {
		options := Options{}
		for _, opt := range opts {
			opt.apply(&options)
		}

		var c context.Context
		var cancel context.CancelFunc

		if options.TimeOut > 100*time.Millisecond {
			c, cancel = context.WithTimeout(ctx, options.TimeOut)
		} else {
			c, cancel = context.WithCancel(ctx)
		}
		defer cancel()

		recoverCh := make(chan interface{}, len(tasks))
		errCh := make(chan error, len(tasks))
		doneCh := make(chan taskResult, len(tasks))

		for _, task := range tasks {
			go func(t AsyncTask) {
				defer func() {
					r := recover()
					if r != nil {
						recoverCh <- r
					}
				}()

				res, err := t(c)

				if err != nil {
					errCh <- err
					return
				}

				doneCh <- res
			}(task)
		}

		results := make(map[uint32]Output, len(tasks))
		for i := 0; i < len(tasks); i++ {
			select {
			case <-c.Done():
				return nil, c.Err()
			case res := <-doneCh:
				results[res.taskID] = Output{
					result: res.value,
				}
			case err := <-errCh:
				return nil, err
			case r := <-recoverCh:
				panic(r)
			}
		}

		return results, nil
	}
}

func NewTask[T any](f func() (T, error)) Task[T] {
	return Task[T]{
		id: rand.Uint32(),
		f:  f,
	}
}

type Task[T any] struct {
	id uint32
	f  func() (T, error)
}

func (t Task[T]) Run() AsyncTask {
	return func(ctx context.Context) (taskResult, error) {
		errCh := make(chan error, 1)
		resCh := make(chan T, 1)
		go func() {
			res, e := t.f()
			if e != nil {
				errCh <- e
			}

			resCh <- res
		}()

		select {
		case <-ctx.Done():
			res := <-resCh
			return taskResult{taskID: t.id, value: res}, ctx.Err()
		case err := <-errCh:
			return taskResult{}, err
		case res := <-resCh:
			return taskResult{taskID: t.id, value: res}, nil
		}
	}
}
