package main

import (
	"fmt"
	"runtime"
	"sync"

	"example.com/golang_test/web_crawler/crawler"
)

// 실제 작업을 처리하는 worker 함수
func worker(done <-chan struct{}, urls chan string, c chan<- crawler.Result) {
	for url := range urls { // urls 채널에서 url 을 가져옴
		select {
		case <-done: // 채널이 닫히면 worker 함수를 빠져나옴
			return
		default:
			crawler.Crawler(url, urls, c) // url 처리
		}
	}
}

func main() {
	numCPUS := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUS)

	urls := make(chan string)      // 작업을 요청할 채널
	done := make(chan struct{})    // 작업 고루틴에 정지 신호를 보낼 채널
	c := make(chan crawler.Result) // 결괏값을 저장할 채널

	var wg sync.WaitGroup
	const numWorkers = 10
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ { // 작업을 처리할 고루틴을 10개 생성
		go func() {
			worker(done, urls, c)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait() // 고루틴이 끝날떄까지 대기
		close(c)  // 대기가 끝나면 결괏값 채널을 닫음
	}()

	urls <- "https://github.com/mg5566?tab=following" // 최초 URL 을 보냄

	count := 0
	for r := range c { // 결과 채널에 값이 들어올 때까지 대기한 뒤 값을 가져옴
		fmt.Println("name:", r.Name)
		fmt.Println("url:", r.Url)

		count++
		if count > 100 { // 100명만 출력한 뒤 done 을 닫아서 worker 고루틴을 종료
			close(done)
			break
		}
	}
}
