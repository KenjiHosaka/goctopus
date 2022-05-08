package goctopus

import (
	"context"
)

type AsyncTask func(ctx context.Context) error

func Orchestrate(ctx context.Context, tasks ...AsyncTask) error {
	c, cancel := context.WithCancel(ctx)
	defer cancel()

	recoverCh := make(chan interface{}, len(tasks))
	errCh := make(chan error, len(tasks))
	doneCh := make(chan bool, len(tasks))

	for _, task := range tasks {
		go func(t AsyncTask) {
			defer func() {
				r := recover()
				if r != nil {
					recoverCh <- r
				}
			}()

			if err := t(c); err != nil {
				errCh <- err
				return
			}

			doneCh <- true
		}(task)
	}

	for i := 0; i < len(tasks); i++ {
		select {
		case <-doneCh:
		case err := <-errCh:
			return err
		case r := <-recoverCh:
			panic(r)
		}
	}

	return nil
}

func Task(f func() error) AsyncTask {
	return func(ctx context.Context) error {
		ch := make(chan error, 1)
		go func() {
			ch <- f()
		}()

		select {
		case <-ctx.Done():
			<-ch
			return ctx.Err()
		case err := <-ch:
			return err
		}
	}
}

func Tasks(tasks ...AsyncTask) AsyncTask {
	return func(ctx context.Context) error {
		for _, task := range tasks {
			if err := task(ctx); err != nil {
				return err
			}
		}

		return nil
	}
}