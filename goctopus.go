package goctopus

import (
	"context"
	"time"
)

type taskResult struct {
	value any
	index int
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

		for i, task := range tasks {
			go func(t AsyncTask, index int) {
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

				res.index = index
				doneCh <- res
			}(task, i)
		}

		results := make(map[int]Output, len(tasks))
		for i := 0; i < len(tasks); i++ {
			select {
			case <-c.Done():
				return nil, c.Err()
			case res := <-doneCh:
				results[res.index] = Output{
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

func Task[T any](f func() (T, error)) AsyncTask {
	return func(ctx context.Context) (taskResult, error) {
		errCh := make(chan error, 1)
		resCh := make(chan T, 1)
		go func() {
			res, e := f()
			if e != nil {
				errCh <- e
			}

			resCh <- res
		}()

		select {
		case <-ctx.Done():
			res := <-resCh
			return taskResult{value: res}, ctx.Err()
		case err := <-errCh:
			return taskResult{}, err
		case res := <-resCh:
			return taskResult{value: res}, nil
		}
	}
}
