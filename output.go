package goctopus

import "errors"

type Output struct {
	result any
}

type Outputs map[int]Output

func (os Outputs) GetResult(taskIndex int) (any, error) {
	o, exist := os[taskIndex]
	if !exist {
		return nil, errors.New("not exist index")
	}

	return o.result, nil
}
