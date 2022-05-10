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
outputs, err := goctopus.Orchestrate(
	context.Background(), 
	goctopus.Task(func() (bool, error) {
		return true, nil
	}), 
	goctopus.Task(func() (string, error) {
		return "result", nil
	}), 
	goctopus.Task(func() (int, error) {
		return 0, nil
	}), 
)()

// get result
task1Res, err := outputs.GetResult(0)
task1Res.(bool)


outputs, err := goctopus.Orchestrate(
	context.Background(), 
	goctopus.Task(func() (bool, error) {
		// ...
		return true, nil
	}), 
	goctopus.Task(func() (bool, error) {
		// ...
		return true, nil
	}), 
)(goctopus.TimeOut{
	Duration: 1 * time.Second,
})
```
