package goctopus

import (
	"time"
)

type Options struct {
	TimeOut time.Duration
}

type Option interface {
	apply(*Options)
}

type TimeOut struct {
	Duration time.Duration
}

func (t TimeOut) apply(o *Options) {
	o.TimeOut = t.Duration
}
