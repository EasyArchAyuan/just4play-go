package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", indexHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	number := rand.Intn(2)
	if number == 0 {
		time.Sleep(time.Second * 10)
		fmt.Fprintf(writer, "slow response")
		return
	}
	fmt.Fprint(writer, "quick response")

}
