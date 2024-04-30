package main

import (
	"fmt"
	"time"
)

func work(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Print("done")
	// 发送一个值来通知我们已经完工啦。
	done <- true
}

func main() {

	done := make(chan bool, 1)
	go work(done)
	<-done

}
