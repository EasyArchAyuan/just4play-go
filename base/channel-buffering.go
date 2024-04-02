package main

import "fmt"

func main() {

	// 这里我们创建了一个字符串通道，最多允许缓存 2 个值。
	messages := make(chan string, 2)

	messages <- "buffered"
	messages <- "channel"

	fmt.Println(<-messages)
	go func() { messages <- "ping" }()

	fmt.Println(<-messages)
}
