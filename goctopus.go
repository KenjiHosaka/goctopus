package goctopus

import (
	"context"
)

type OrchestrationFunc func(...Option) error

func OrchestrateRuns(ctx context.Context, tasks ...AsyncTask) OrchestrationFunc {
	return func(opts ...Option) error {
		options := Options{}
		for _, opt := range opts {
			opt.apply(&options)
		}

		var c context.Context
		var cancel context.CancelFunc

		if options.TimeOut > 0 {
			c, cancel = context.WithTimeout(ctx, options.TimeOut)
		} else {
			c, cancel = context.WithCancel(ctx)
		}
		defer cancel()

		recoverCh := make(chan interface{}, len(tasks))
		errCh := make(chan error, len(tasks))
		doneCh := make(chan struct{}, len(tasks))

		for _, task := range tasks {
			go func(t AsyncTask) {
				defer func() {
					r := recover()
					if r != nil {
						recoverCh <- r
					}
				}()

				err := t(c)

				if err != nil {
					errCh <- err
					return
				}

				doneCh <- struct{}{}
			}(task)
		}

		for i := 0; i < len(tasks); i++ {
			select {
			case <-c.Done():
				return c.Err()
			case <-doneCh:
			case err := <-errCh:
				return err
			case r := <-recoverCh:
				panic(r)
			}
		}

		return nil
	}
}

func OrchestrateTasks[T any](ctx context.Context, tasks []*Task[T]) OrchestrationFunc {
	return func(opts ...Option) error {
		options := Options{}
		for _, opt := range opts {
			opt.apply(&options)
		}

		var c context.Context
		var cancel context.CancelFunc

		if options.TimeOut > 0 {
			c, cancel = context.WithTimeout(ctx, options.TimeOut)
		} else {
			c, cancel = context.WithCancel(ctx)
		}
		defer cancel()

		recoverCh := make(chan interface{}, len(tasks))
		errCh := make(chan error, len(tasks))
		doneCh := make(chan struct{}, len(tasks))

		for _, task := range tasks {
			go func(t AsyncTask) {
				defer func() {
					r := recover()
					if r != nil {
						recoverCh <- r
					}
				}()

				err := t(c)

				if err != nil {
					errCh <- err
					return
				}

				doneCh <- struct{}{}
			}(task.Run())
		}

		for i := 0; i < len(tasks); i++ {
			select {
			case <-c.Done():
				return c.Err()
			case <-doneCh:
			case err := <-errCh:
				return err
			case r := <-recoverCh:
				panic(r)
			}
		}

		return nil
	}
}
