package goctopus

import (
	"context"
	"math/rand"
)

type AsyncTask func(ctx context.Context) error

type Task[T any] struct {
	id     uint32
	f      func() (T, error)
	result T
}

func NewTask[T any](f func() (T, error)) Task[T] {
	return Task[T]{
		id: rand.Uint32(),
		f:  f,
	}
}

func (t *Task[T]) Run() AsyncTask {
	return func(ctx context.Context) error {
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
			t.result = res
			return ctx.Err()
		case err := <-errCh:
			return err
		case res := <-resCh:
			t.result = res
			return nil
		}
	}
}

func (t Task[T]) Result() T {
	return t.result
}
