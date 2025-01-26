package crawler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

func fetch(url string) (*html.Node, error) {
	res, err := http.Get(url) // URL 에서 HTML 데이터를 가져옴
	if err != nil {
		log.Printf("URL %s 가져오기 실패: %v", url, err)
		return nil, err
	}

	doc, err := html.Parse(res.Body) // res.Body 를 넣으면 파싱된 데이터가 반환됨
	if err != nil {
		log.Println("HTML 파싱 실패: ", err)
		return nil, err
	}

	return doc, nil
}

func parseFollowing(doc *html.Node) []string {
	var urls = make([]string, 0)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" { // img 태그
			for _, a := range n.Attr {
				if a.Key == "class" && a.Val == "avatar" { // class 가 avatar 인 경우
					for _, a := range n.Attr {
						if a.Key == "alt" {
							fmt.Println(a.Val) // 사용자 이름 출력
							user := a.Val[1:]
							urls = append(urls, "https://github.com/"+user+"?tab=following")
							break
						}
					}
				}
				if a.Key == "class" && a.Val == "avatar avatar-user" { // class 가 "avatar avatar-user" 인 경우
					for _, a := range n.Attr {
						if a.Key == "alt" {
							fmt.Println(a.Val)
							user := a.Val[1:]
							urls = append(urls, "https://github.com/"+user+"?tab=following")
							break
						}
					}
				}

				// if a.Key == "class" && a.Val == "gravatar" { // class 가 gravatar 인 경우
				// 	user := n.Parent.Attr[0].Val // 부모 요소의 첫번째 속성(href)

				// 	// 사용자 이름으로 팔로잉 URL 조합
				// 	urls = append(urls, "https://github.com"+user+"?tab=following")
				// 	break
				// }
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c) // 재귀호출로 자식과 형제를 모두 탐색
		}
	}
	f(doc)

	return urls
}

var fetched = struct {
	m map[string]error // 중복 검사를 위한 URL 과 에러 값 저장
	sync.Mutex
}{m: make(map[string]error)} // 변수를 선언하면서 이름이 없는 구조체를 정의, 초기값 생성하여 대입

func Crawler(url string) {
	fetched.Lock() // 맵은 뮤텍스로 보호
	if _, ok := fetched.m[url]; ok {
		fetched.Unlock()
		return
	}
	fetched.Unlock()

	doc, err := fetch(url) // URL 에서 파싱된 HTML 데이터를 가져옴

	fetched.Lock()
	fetched.m[url] = err // 가져온 URL 은 맵에 URL 과 에러 값 저장
	fetched.Unlock()

	urls := parseFollowing(doc) // 사용자 정보 출력, 팔로잉 URL 을 구함

	done := make(chan bool)
	for _, u := range urls {
		go func(url string) {
			// time.Sleep(1 * time.Second)
			Crawler(url)
			done <- true
		}(u)
	}

	for range urls {
		<-done
	}
}
