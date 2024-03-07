package internal

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/luis-olivetti/go-stresstest/internal/model"
)

var (
	index      int
	indexMutex sync.Mutex
)

func TestStress(t *testing.T) {
	// Arrange
	urlStr := "http://localhost:8080"
	concurrency := 2
	requests := 4

	index = 0
	server := newServerHTTP([]int{200, 400, 500, 200})
	go server.Run(":8080")

	stress := NewStress(urlStr, concurrency, requests)
	// Act
	report := stress.Run()

	// Assert
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

func TestPrintReport(t *testing.T) {
	// Arrange
	urlStr := "http://localhost:8080"
	concurrency := 2
	requests := 100

	stress := NewStress(urlStr, concurrency, requests)

	report := &model.Report{
		TotalTime:     10,
		TotalRequests: 100,
		HTTPCodes: map[int]int{
			200: 50,
			404: 20,
			500: 30,
		},
	}

	file, err := os.CreateTemp("", "output.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	old := os.Stdout
	os.Stdout = file
	defer func() {
		os.Stdout = old
	}()

	// Act
	stress.PrintReport(report)

	// Assert
	expectedOutput := fmt.Sprintf("\nReport:\n\n> Total time: %s\n\n> Total requests: %d\n\n> HTTP Code %d: %d req(s)\n> HTTP Code %d: %d req(s)\n> HTTP Code %d: %d req(s)\n",
		report.TotalTime, report.TotalRequests,
		200, report.HTTPCodes[200],
		404, report.HTTPCodes[404],
		500, report.HTTPCodes[500])

	fileContent, err := os.ReadFile(file.Name())
	if err != nil {
		t.Errorf("Failed to read file: %s", err)
		return
	}

	if string(fileContent) != expectedOutput {
		t.Errorf("PrintReport output is incorrect, got: %s, want: %s", string(fileContent), expectedOutput)
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
