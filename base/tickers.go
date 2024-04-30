package main

import (
	"fmt"
	"time"
)

func main() {

	ticker := time.NewTicker(time.Millisecond * 500)

	go func() {
		for ticker := range ticker.C {
			fmt.Println("tick at", ticker)
		}
	}()

}
