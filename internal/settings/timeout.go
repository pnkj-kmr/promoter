package settings

import "time"

const (
	DefaultSleep = 30 * time.Second
	CheckTimeout = 5 * time.Second
	QueryTimeout = 30 * time.Second
	LongTimeout  = 90 * time.Second
)
