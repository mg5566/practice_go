package crawler

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/html"
)

func fetch(url string) (*html.Node, error) {
	// 0.5~1초 사이의 랜덤 딜레이 추가
	delay := time.Duration(500+rand.Intn(500)) * time.Millisecond
	time.Sleep(delay)

	// 요청 시작
	fmt.Println("요청 시작: ", url)
	res, err := http.Get(url)
	if err != nil {
		log.Printf("URL %s 가져오기 실패: %v", url, err)
		return nil, err
	}
	defer res.Body.Close() // response body를 반드시 닫아줌

	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Println("HTML 파싱 실패: ", err)
		return nil, err
	}

	// 요청 완료
	fmt.Println("요청 완료: ", url)
	return doc, nil
}

func parseFollowing(doc *html.Node, urls chan string) <-chan string {
	name := make(chan string)

	go func() { // 교착 상태가 되지 않도록 고루틴으로 실행
		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "img" { // img 태그
				for _, a := range n.Attr {
					if a.Key == "class" && a.Val == "avatar" { // class 가 avatar 인 경우
						for _, a := range n.Attr {
							if a.Key == "alt" {
								fmt.Println("attribute value: ", a.Val) // 사용자 이름 출력
								user := a.Val[1:]
								fmt.Println("user: ", user)
								// urls = append(urls, "https://github.com/"+user+"?tab=following")
								name <- user
								urls <- "https://github.com/" + user + "?tab=following"

								break
							}
						}
					}
					if a.Key == "class" && a.Val == "avatar avatar-user" { // class 가 "avatar avatar-user" 인 경우
						for _, a := range n.Attr {
							if a.Key == "alt" {
								fmt.Println("attribute value: ", a.Val) // 사용자 이름 출력
								user := a.Val[1:]
								fmt.Println("user: ", user)
								tempUrl := "https://github.com/" + user + "?tab=following"
								fmt.Println("tempUrl: ", tempUrl)
								name <- user
								urls <- tempUrl
								break
							}
						}
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c) // 재귀 호출로 자식과 형제를 모두 탐색
			}
		}
		f(doc)
	}()

	return name
}

var fetched = struct {
	m map[string]error // 중복 검사를 위한 URL 과 에러 값 저장
	sync.Mutex
}{m: make(map[string]error)} // 변수를 선언하면서 이름이 없는 구조체를 정의, 초기값 생성하여 대입

func Crawler(url string, urls chan string, c chan<- Result) {
	fetched.Lock()                   // 맵은 뮤텍스로 보호
	if _, ok := fetched.m[url]; ok { // URL 중복 처리 여부를 검사
		fetched.Unlock()
		return
	}
	fetched.Unlock()

	doc, err := fetch(url) // URL 에서 파싱된 HTML 데이터를 가져옴
	if err != nil {        // URL 을 가져오지 못했을 때
		go func(u string) { // 교착 상태가 되지 않도록 고루틴을 생성
			urls <- u // 채널에 URL 을 보냄
		}(url)
		return
	}

	fetched.Lock()
	fetched.m[url] = err // 가져온 URL 은 맵에 URL 과 에러 값 저장
	fetched.Unlock()

	name := <-parseFollowing(doc, urls) // 사용자 정보, 팔로잉 URL 을 구함
	c <- Result{
		Url:  url,
		Name: name,
	} // 가져온 URL 과 사용자 이름을 구조체 인스턴스로 생성하여 채널 c 에 보냄
}
