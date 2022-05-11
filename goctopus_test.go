package goctopus

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestOrchestrate_SuccessAllTasks(test *testing.T) {
	task1 := NewTask[bool](func() (bool, error) {
		time.Sleep(10 * time.Millisecond)
		return true, nil
	})
	task2 := NewTask[string](func() (string, error) {
		time.Sleep(10 * time.Millisecond)
		return "test", nil
	})
	task3 := NewTask[bool](func() (bool, error) {
		time.Sleep(10 * time.Millisecond)
		return true, nil
	})

	start := time.Now()
	outputs, err := Orchestrate(
		context.Background(),
		task1.Run(),
		task2.Run(),
		task3.Run(),
	)()

	diff := time.Now().Sub(start)
	if diff > 100*time.Millisecond {
		test.Errorf("Too late")
	}
	res1, _ := FindResult(outputs, task1)
	res2, _ := FindResult(outputs, task2)
	res3, _ := FindResult(outputs, task3)

	if err != nil || !(res1 && res2 == "test" && res3) {
		test.Errorf("One or more failed")
	}
}

func TestOrchestrate_TimeOut(test *testing.T) {
	task1 := NewTask[bool](func() (bool, error) {
		time.Sleep(10 * time.Millisecond)
		return true, nil
	})
	task2 := NewTask[bool](func() (bool, error) {
		time.Sleep(2 * time.Second)
		return true, nil
	})

	outputs, err := Orchestrate(
		context.Background(),
		task1.Run(),
		task2.Run(),
	)(TimeOut{
		Duration: 1 * time.Second,
	})
	if err == nil {
		test.Errorf("Failed to handle error")
	}

	_, exist := FindResult(outputs, task2)
	if exist {
		test.Errorf("Failed to cancel task2")
	}
}

func TestOrchestrate_CancelTask(test *testing.T) {
	task1 := NewTask[bool](func() (bool, error) {
		time.Sleep(10 * time.Millisecond)
		return true, nil
	})
	task2 := NewTask[string](func() (string, error) {
		time.Sleep(20 * time.Millisecond)
		return "", errors.New("task2 error occurred")
	})
	task3 := NewTask[bool](func() (bool, error) {
		time.Sleep(10 * time.Millisecond)
		return true, nil
	})

	outputs, err := Orchestrate(
		context.Background(),
		task1.Run(),
		task2.Run(),
		task3.Run(),
	)()
	if err == nil {
		test.Errorf("Failed to handle error")
	}

	_, exist := FindResult(outputs, task3)
	if exist {
		test.Errorf("Failed to cancel task3")
	}
}
