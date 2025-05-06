package main

// fmt 패키지를 가져와 표준 출력 기능을 사용
import "fmt"

func main() {
	// 문자열을 전송할 수 있는 채널 생성
	messages := make(chan string)
	// 불리언 값을 전송할 수 있는 채널 생성
	signals := make(chan bool)

	// 고루틴 생성 - 별도의 스레드에서 실행되는 함수
	go func() {
		// messages 채널에서 메시지를 수신할 때까지 대기
		// 이 고루틴은 메인 함수와 병렬로 실행됨
		msg := <-messages
		fmt.Println("received goroutine message", msg)
	}()

	// select 문을 사용하여 채널 작업 처리
	// 논블로킹 방식으로 채널을 확인 (default 케이스가 있으므로)
	select {
	case msg := <-messages:
		// messages 채널에서 메시지를 수신할 수 있으면 실행
		fmt.Println("received message", msg)
	default:
		// 즉시 수신할 메시지가 없으면 실행 (논블로킹)
		fmt.Println("no message received")
	}

	// "hi" 메시지 정의
	msg := "hi"
	// select 문을 사용하여 논블로킹 방식으로 메시지 전송 시도
	select {
	case messages <- msg:
		// 메시지 전송이 성공하면 실행
		fmt.Println("sent message", msg)
	default:
		// 즉시 메시지를 전송할 수 없으면 실행 (채널이 꽉 찼거나 수신자가 없는 경우)
		fmt.Println("no message sent")
	}

	// 다시 select 문을 사용하여 여러 채널에서 데이터를 기다림
	select {
	case msg := <-messages:
		// messages 채널에서 데이터를 수신하면 실행
		fmt.Println("received message", msg)
	case sig := <-signals:
		// signals 채널에서 데이터를 수신하면 실행
		fmt.Println("received signal", sig)
	default:
		// 즉시 수신할 데이터가 없으면 실행 (논블로킹)
		fmt.Println("no activity")
	}
}