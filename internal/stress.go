package internal

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/luis-olivetti/go-stresstest/internal/model"
)

type Stress struct {
	url          *url.URL
	concurrency  int
	requests     int
	responseData map[int]int
}

func NewStress(urlStr string, concurrency, requests int) *Stress {
	url, _ := url.Parse(urlStr)

	return &Stress{
		url:          url,
		concurrency:  concurrency,
		requests:     requests,
		responseData: make(map[int]int),
	}
}

func (s *Stress) Run() *model.Report {
	spin := spinner.New(spinner.CharSets[43], 100*time.Millisecond)
	spin.Color("red")
	spin.Start()

	initTime := time.Now()

	chanReq := make(chan string, s.requests)
	chanResp := make(chan *http.Response, s.requests)

	var wg sync.WaitGroup
	for i := 0; i < s.concurrency; i++ {
		wg.Add(1)
		go execute(chanReq, chanResp, &wg)
	}

	for i := 0; i < s.requests; i++ {
		chanReq <- s.url.String()
	}
	// Sinaliza as go routines que não terá mais requisições
	close(chanReq)

	wg.Wait()
	close(chanResp)

	var respCount int
	for resp := range chanResp {
		if resp != nil {
			s.responseData[resp.StatusCode]++
		}
		respCount++
	}

	finalTime := time.Now()

	spin.Stop()

	return &model.Report{
		TotalTime:     finalTime.Sub(initTime),
		TotalRequests: respCount,
		HTTPCodes:     s.responseData,
	}
}

func (s *Stress) PrintReport(report *model.Report) {
	fmt.Println("\nReport:")
	fmt.Printf("\n> Total time: %s\n", report.TotalTime)
	fmt.Printf("\n> Total requests: %d\n\n", report.TotalRequests)

	for statusCode, count := range report.HTTPCodes {
		fmt.Printf("> HTTP Code %d: %d req(s)\n", statusCode, count)
	}
}

func makeRequest(client *http.Client, url string) (*http.Response, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func execute(requests chan string, responses chan *http.Response, wg *sync.WaitGroup) {
	defer wg.Done()

	httpClient := &http.Client{}

	// Quando vai ler o channel requests, ele fica bloqueado até que tenha algo para ler
	for url := range requests {
		resp, err := makeRequest(httpClient, url)
		if err != nil {
			resp = &http.Response{
				StatusCode: 500,
			}
		}
		responses <- resp
	}
}
