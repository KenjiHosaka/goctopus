package goctopus

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

type MutexBool struct {
	mu    sync.RWMutex
	value bool
}

type MutexTime struct {
	mu    sync.RWMutex
	value time.Time
}

func TestOrchestrate_SuccessAllTasks(test *testing.T) {
	start := time.Now()
	var res1, res2, res3 MutexBool

	err := Orchestrate(
		context.Background(),
		Task(func() error {
			res1.mu.Lock()
			defer res1.mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			res1.value = true
			return nil
		}),
		Task(func() error {
			res2.mu.Lock()
			defer res2.mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			res2.value = true
			return nil
		}),
		Task(func() error {
			res3.mu.Lock()
			defer res3.mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			res3.value = true
			return nil
		}),
	)()

	diff := time.Now().Sub(start)
	if diff > 100*time.Millisecond {
		test.Errorf("Too late")
	}
	if err != nil || !(res1.value && res2.value && res3.value) {
		test.Errorf("One or more failed")
	}
}

func TestOrchestrate_TimeOut(test *testing.T) {
	var res1, res2 MutexBool

	err := Orchestrate(
		context.Background(),
		Task(func() error {
			res1.mu.Lock()
			defer res1.mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			res1.value = true
			return nil
		}),
		Task(func() error {
			res2.mu.Lock()
			defer res2.mu.Unlock()
			time.Sleep(5 * time.Second)
			res2.value = true
			return nil
		}),
	)(TimeOut{
		Duration: 1 * time.Second,
	})

	if err == nil {
		test.Errorf("Failed to handle error")
	}

	if !(res1.value && !res2.value) {
		test.Errorf("Failed to cancel task2")
	}
}

func TestOrchestrate_CancelTask(test *testing.T) {
	var res1, res2, res3 MutexBool

	err := Orchestrate(
		context.Background(),
		Task(func() error {
			res1.mu.Lock()
			defer res1.mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			res1.value = true
			return nil
		}),
		Task(func() error {
			res2.mu.Lock()
			defer res2.mu.Unlock()
			time.Sleep(20 * time.Millisecond)
			return errors.New("task2 error occurred")
		}),
		Task(func() error {
			res3.mu.Lock()
			defer res3.mu.Unlock()
			time.Sleep(120 * time.Millisecond)
			res3.value = true
			return nil
		}),
	)()

	if err == nil {
		test.Errorf("Failed to handle error")
	}

	if !(res1.value && !res2.value && !res3.value) {
		test.Errorf("Failed to cancel task3")
	}
}

func TestTasks(test *testing.T) {
	var task2Start, task3Start MutexTime
	_ = Orchestrate(
		context.Background(),
		Task(func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}),
		Tasks(
			Task(func() error {
				task2Start.mu.Lock()
				defer task2Start.mu.Unlock()
				task2Start.value = time.Now()
				time.Sleep(10 * time.Millisecond)
				return nil
			}),
			Task(func() error {
				task3Start.mu.Lock()
				defer task3Start.mu.Unlock()
				task3Start.value = time.Now()
				time.Sleep(10 * time.Millisecond)
				return nil
			}),
		),
	)()

	diff := task2Start.value.Sub(task3Start.value)
	if diff >= 10*time.Millisecond {
		test.Errorf("Not processed in order")
	}
}

func TestTasks_CancelTask(test *testing.T) {
	var res1, res2, res3, res4 MutexBool
	err := Orchestrate(
		context.Background(),
		Task(func() error {
			res1.mu.Lock()
			defer res1.mu.Unlock()
			time.Sleep(10 * time.Millisecond)
			res1.value = true
			return nil
		}),
		Tasks(
			Task(func() error {
				time.Sleep(20 * time.Millisecond)
				return errors.New("task2 error occurred")
			}),
			Task(func() error {
				res3.mu.Lock()
				defer res3.mu.Unlock()
				time.Sleep(10 * time.Millisecond)
				res3.value = true
				return nil
			}),
		),
		Task(func() error {
			res4.mu.Lock()
			defer res4.mu.Unlock()
			time.Sleep(120 * time.Millisecond)
			res4.value = true
			return nil
		}),
	)()

	if err == nil {
		test.Errorf("Failed to handle error")
	}

	if !(res1.value && !res2.value && !res3.value && !res4.value) {
		test.Errorf("Failed to cancel task4")
	}
}
