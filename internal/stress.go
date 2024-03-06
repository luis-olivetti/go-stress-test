package internal

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/briandowns/spinner"
)

type Stress struct {
	url          *url.URL
	concurrency  int
	requests     int
	responseData map[int]int
	mu           sync.Mutex
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

func (s *Stress) Run() {
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
	close(chanReq)

	go func() {
		wg.Wait()
		close(chanResp)
	}()

	for resp := range chanResp {
		if resp != nil {
			s.mu.Lock()
			s.responseData[resp.StatusCode]++
			s.mu.Unlock()
		}
	}

	finalTime := time.Now()

	spin.Stop()

	fmt.Println("\nReport:")
	fmt.Printf("\nTotal time: %s\n", finalTime.Sub(initTime))

	for statusCode, count := range s.responseData {
		fmt.Printf("HTTP Code %d: %d req(s)\n", statusCode, count)
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
