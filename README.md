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
### Simple
```golang
outputs, err := goctopus.Orchestrate(
	context.Background(), 
	goctopus.NewTask[bool](func() (bool, error) {
		return true, nil
	}).Run(), 
	goctopus.NewTask[string](func() (string, error) {
		return "result", nil
	}).Run(), 
	goctopus.NewTask[int](func() (int, error) {
		return 0, nil
	}).Run(), 
)()
```

### Get Result
```golang
task1 := goctopus.NewTask[bool](func() (bool, error) {
	time.Sleep(10 * time.Millisecond)
	return true, nil
})

outputs, err := goctopus.Orchestrate(
	context.Background(), 
	task1.Run(),
)()
res, exist := goctopus.FindResult(outputs, task1)
// res: true
```

### Set Timeout
```golang
outputs, err := goctopus.Orchestrate(
	context.Background(), 
	goctopus.NewTask[bool](func() (bool, error) {
		return true, nil
	}).Run(), 
)(goctopus.TimeOut{
	Duration: 1 * time.Second,
})
```


