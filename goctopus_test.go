package goctopus

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestOrchestrate_SuccessAllTasks(test *testing.T) {
	defer goleak.VerifyNone(test)

	task1 := NewTask[bool](func(ctx context.Context) (bool, error) {
		t := time.NewTicker(10 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-t.C:
				return true, nil
			}
		}
	})
	task2 := NewTask[string](func(ctx context.Context) (string, error) {
		t := time.NewTicker(10 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-t.C:
				return "test", nil
			}
		}
	})
	task3 := NewTask[bool](func(ctx context.Context) (bool, error) {
		t := time.NewTicker(10 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-t.C:
				return true, nil
			}
		}
	})

	start := time.Now()
	err := Orchestrate(
		context.Background(),
		task1.Run(),
		task2.Run(),
		task3.Run(),
	)()

	diff := time.Now().Sub(start)
	if diff > 100*time.Millisecond {
		test.Errorf("Too late")
	}

	if err != nil || !(task1.result && task2.Result() == "test" && task3.Result()) {
		test.Errorf("One or more failed")
	}
}

func TestOrchestrate_TimeOut(test *testing.T) {
	defer goleak.VerifyNone(test)

	task1 := NewTask[bool](func(ctx context.Context) (bool, error) {
		t := time.NewTicker(10 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-t.C:
				return true, nil
			}
		}
	})
	task2 := NewTask[bool](func(ctx context.Context) (bool, error) {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-t.C:
				return true, nil
			}
		}
	})

	err := Orchestrate(
		context.Background(),
		task1.Run(),
		task2.Run(),
	)(TimeOut{
		Duration: 1 * time.Second,
	})
	if err == nil {
		test.Errorf("Failed to handle error")
	}

	if task2.Result() {
		test.Errorf("Failed to cancel task2")
	}
}

func TestOrchestrate_CancelTask(test *testing.T) {
	defer goleak.VerifyNone(test)

	task1 := NewTask[bool](func(ctx context.Context) (bool, error) {
		t := time.NewTicker(10 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-t.C:
				return true, nil
			}
		}
	})
	task2 := NewTask[string](func(ctx context.Context) (string, error) {
		t := time.NewTicker(20 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-t.C:
				return "", errors.New("task2 error occurred")
			}
		}
	})
	task3 := NewTask[bool](func(ctx context.Context) (bool, error) {
		t := time.NewTicker(50 * time.Millisecond)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return false, ctx.Err()
			case <-t.C:
				return true, nil
			}
		}
	})

	err := Orchestrate(
		context.Background(),
		task1.Run(),
		task2.Run(),
		task3.Run(),
	)()
	if err == nil {
		test.Errorf("Failed to handle error")
	}

	if task3.Result() {
		test.Errorf("Failed to cancel task3")
	}
}
