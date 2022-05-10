package goctopus

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestOrchestrate_SuccessAllTasks(test *testing.T) {
	start := time.Now()
	outputs, err := Orchestrate(
		context.Background(),
		Task(func() (bool, error) {
			time.Sleep(10 * time.Millisecond)
			return true, nil
		}),
		Task(func() (string, error) {
			time.Sleep(10 * time.Millisecond)
			return "test", nil
		}),
		Task(func() (bool, error) {
			time.Sleep(10 * time.Millisecond)
			return true, nil
		}),
	)()

	diff := time.Now().Sub(start)
	if diff > 100*time.Millisecond {
		test.Errorf("Too late")
	}
	res1, _ := outputs.GetResult(0)
	res2, _ := outputs.GetResult(1)
	res3, _ := outputs.GetResult(2)

	if err != nil || !(res1.(bool) && res2.(string) == "test" && res3.(bool)) {
		test.Errorf("One or more failed")
	}
}

func TestOrchestrate_TimeOut(test *testing.T) {
	outputs, err := Orchestrate(
		context.Background(),
		Task(func() (bool, error) {
			time.Sleep(10 * time.Millisecond)
			return true, nil
		}),
		Task(func() (bool, error) {
			time.Sleep(2 * time.Second)
			return true, nil
		}),
	)(TimeOut{
		Duration: 1 * time.Second,
	})
	if err == nil {
		test.Errorf("Failed to handle error")
	}

	_, err = outputs.GetResult(1)
	if err == nil {
		test.Errorf("Failed to cancel task2")
	}
}

func TestOrchestrate_CancelTask(test *testing.T) {
	outputs, err := Orchestrate(
		context.Background(),
		Task(func() (bool, error) {
			time.Sleep(10 * time.Millisecond)
			return true, nil
		}),
		Task(func() (string, error) {
			time.Sleep(20 * time.Millisecond)
			return "", errors.New("task2 error occurred")
		}),
		Task(func() (bool, error) {
			time.Sleep(120 * time.Millisecond)
			return true, nil
		}),
	)()
	if err == nil {
		test.Errorf("Failed to handle error")
	}

	_, err = outputs.GetResult(2)
	if err == nil {
		test.Errorf("Failed to cancel task3")
	}
}
