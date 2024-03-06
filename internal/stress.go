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

	var wg sync.WaitGroup
	for i := 0; i < s.concurrency; i++ {
		wg.Add(1)
		go func(c int) {
			defer wg.Done()

			httpClient := &http.Client{}

			for j := 0; j < s.requests; j++ {
				//fmt.Printf("Concurrency %d -> Fazendo solicitação %d de %d\n", c, j+1, s.requests)

				resp, err := makeRequest(httpClient, &wg, s.url.String())

				if err != nil {
					//fmt.Printf("Erro ao fazer solicitação: %s\n", err)
					continue
				}

				s.mu.Lock()
				s.responseData[resp.StatusCode]++
				s.mu.Unlock()
			}
		}(i)
	}
	wg.Wait()

	finalTime := time.Now()

	spin.Stop()

	fmt.Println("Relatório de solicitações:")
	fmt.Printf("Tempo total decorrido: %s\n", finalTime.Sub(initTime))

	for statusCode, count := range s.responseData {
		fmt.Printf("Código de status %d: %d solicitações\n", statusCode, count)
	}
}

func makeRequest(client *http.Client, wg *sync.WaitGroup, url string) (*http.Response, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}
