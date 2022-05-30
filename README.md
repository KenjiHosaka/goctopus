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
task1 := goctopus.NewTask[bool](func(ctx context.Context) (bool, error) {
	return true, nil
})

task2 := goctopus.NewTask[string](func(ctx context.Context) (string, error) {
	return "result", nil
})

err := goctopus.OrchestrateRuns(
	context.Background(), 
	task1.Run(), 
	task2.Run(), 
)()

// same type results
tasks := []*goctopus.Task{&task1}
err := goctopus.OrchestrateTasks(
	context.Background(), 
	tasks,
)()

// Get Result
res := task2.Result()
// res: result

```

### Set Timeout
```golang
err := goctopus.OrchestrateRuns(
	// ...
)(goctopus.TimeOut{
	Duration: 1 * time.Second,
})
```


