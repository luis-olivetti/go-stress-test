package internal

import (
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	index      int
	indexMutex sync.Mutex
)

func TestStress(t *testing.T) {
	urlStr := "http://localhost:8080"
	concurrency := 2
	requests := 4

	index = 0
	server := newServerHTTP([]int{200, 400, 500, 200})
	go server.Run(":8080")

	stress := NewStress(urlStr, concurrency, requests)
	report := stress.Run()

	if report.TotalRequests != requests {
		t.Errorf("Expected %d, got %d", requests, report.TotalRequests)
	}

	if report.TotalTime.Seconds() == 0 {
		t.Errorf("Expected > 0, got %f", report.TotalTime.Seconds())
	}

	if report.HTTPCodes[200] != 2 {
		t.Errorf("Expected 2, got %d", report.HTTPCodes[200])
	}

	if report.HTTPCodes[400] != 1 {
		t.Errorf("Expected 1, got %d", report.HTTPCodes[400])
	}

	if report.HTTPCodes[500] != 1 {
		t.Errorf("Expected 1, got %d", report.HTTPCodes[500])
	}
}

func newServerHTTP(httpCodes []int) *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		indexMutex.Lock()
		defer indexMutex.Unlock()

		if index >= len(httpCodes) {
			index = 0
		}

		c.String(httpCodes[index], "Hello, World!")
		index++
	})

	return router
}
