# goctopus
[![codecov](https://codecov.io/gh/KenjiHosaka/goctopus/branch/main/graph/badge.svg?token=ET0SRXKUKZ)](https://codecov.io/gh/KenjiHosaka/goctopus)

Easy to run goroutine

## Features
- Execute functions at concurrent
- Cancel functions if an error occurs in a function

## Installation
```
go get github.com/KenjiHosaka/goctopus
```

## How to use
```golang
type MutexBool struct {
	mu    sync.RWMutex
	value bool
}
var res1, res2, res3 MutexBool

err := goctopus.Orchestrate(
	context.Background(), 
	goctopus.Task(func() error {
		res1.mu.Lock()
		defer res1.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
		res1.value = true
		return nil
	}), 
	goctopus.Task(func() error {
		res2.mu.Lock()
		defer res2.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
		res2.value = true
		return nil
	}), 
	goctopus.Task(func() error {
		res3.mu.Lock()
		defer res3.mu.Unlock()
		time.Sleep(10 * time.Millisecond)
		res3.value = true
		return nil
	}), 
)()

err := Orchestrate(
	context.Background(), 
	Task(func() error {
		// ...
		return nil
	}), 
	Task(func() error {
		// ...
		return nil
	}), 
)(TimeOut{
	Duration: 1 * time.Second,
})
```
