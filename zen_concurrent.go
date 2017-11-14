package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func worker(in <-chan bool, out chan<- string, putback chan<- bool) {
	for job := range in {
		client := &http.Client{}
		res, err := client.Get("https://api.github.com/zen")
		if err != nil {
			putback <- job
			break
		}
		if res.StatusCode == http.StatusOK {
			body, _ := ioutil.ReadAll(res.Body)
			out <- string(body)
		} else {
			putback <- job
		}
		res.Body.Close()
	}
}

func main() {
	count := 3

	wg := sync.WaitGroup{}
	jobs := make(chan bool, count)
	resp := make(chan string)
	messages := make(map[string]string)

	for i := 0; i < count; i++ {
		jobs <- true
		wg.Add(1)
	}

	for i := 0; i < count; i++ {
		go worker(jobs, resp, jobs)
	}

	go func() {
		for res := range resp {
			if _, ok := messages[res]; ok {
				jobs <- true
			} else {
				fmt.Println(res)
				messages[res] = res
				wg.Done()
			}
		}
	}()

	wg.Wait()
}
