package model

import "time"

type Report struct {
	TotalTime     time.Duration
	TotalRequests int
	HTTPCodes     map[int]int
}
