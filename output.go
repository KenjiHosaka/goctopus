package goctopus

type Output struct {
	result any
}

type Outputs map[uint32]Output

func FindResult[T any](outputs Outputs, task Task[T]) (T, bool) {
	val, exist := outputs[task.id].result.(T)
	return val, exist
}
