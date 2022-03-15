package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type counter struct {
	mu    sync.Mutex
	count int
}

func main() {
	var (
		workers = make(chan struct{}, 5)
		count = counter{}
	)

	for num := range processUrls(readUrls(), workers, httpClient()) {
		count.mu.Lock()
		count.count += num
		count.mu.Unlock()
	}

	fmt.Printf("Total: %d\n", count.count)

}

func readUrls() <-chan string {
	output := make(chan string)

	go func() {
		defer close(output)
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			output <- scanner.Text()
		}

	}()

	return output
}

func processUrls(input <-chan string, limiter chan struct{}, client *http.Client) <-chan int {
	output := make(chan int)
	wg := new(sync.WaitGroup)

	go func() {
		defer close(output)
		defer wg.Wait()

		for line := range input {
			limiter <- struct{}{}
			wg.Add(1)
			go func(line string) {
				defer wg.Done()
				defer func() {
					<- limiter
				}()

				if _, err := url.ParseRequestURI(line); err != nil {
					logError(line, err)
					return
				}

				req, err := http.NewRequest("GET", line, nil); if err != nil {
					logError(line, err)
					return
				}

				resp, err := client.Do(req); if err != nil {
					logError(line, err)
					return
				}

				defer resp.Body.Close()

				bodyBytes, err := io.ReadAll(resp.Body); if err != nil {
					logError(line, err)
					return
				}

				count := strings.Count(strings.ToLower(string(bodyBytes)), "go")

				fmt.Printf("Count for %s: %d\n", line, count)

				output <- count
			}(line)
		}
	}()

	return output
}

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func logError(line string, err error) {
	log.Printf("skipped: %s, error: %s\n", line, err.Error())
}
