package test

import "time"

const (
	pollingInterval     = 5 * time.Second
	pollingTimeout      = 10 * time.Minute
	resourceWaitTimeout = 1 * time.Minute
)

var (
	retries = int(resourceWaitTimeout.Seconds() / pollingInterval.Seconds())
)
