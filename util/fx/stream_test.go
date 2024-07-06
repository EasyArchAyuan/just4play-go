package fx

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"
)

func inputStream(ch chan int) {
	count := 0
	for {
		ch <- count
		time.Sleep(time.Millisecond * 500)
		count++
	}
}

func outputStream(ch chan int) {
	From(func(source chan<- interface{}) {
		for c := range ch {
			source <- c
		}
	}).Walk(func(item interface{}, pipe chan<- interface{}) {
		count := item.(int)
		pipe <- count
	}).Filter(func(item interface{}) bool {
		itemInt := item.(int)
		if itemInt%2 == 0 {
			return true
		}
		return false
	}).ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func TestFx(t *testing.T) {
	ch := make(chan int)

	go inputStream(ch)
	go outputStream(ch)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c
}

func TestFrom(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	From(func(source chan<- any) {
		for _, v := range s {
			source <- v
		}
	})
	t.Log(s)
}

func TestFilter(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	From(func(source chan<- interface{}) {
		for _, v := range s {
			source <- v
		}
	}).Filter(func(item interface{}) bool {
		//保留偶数
		if item.(int)%2 == 0 {
			return true
		}
		return false
	}).ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func TestGroup(t *testing.T) {
	ss := []string{"golang", "google", "php", "python", "java", "c++"}
	From(func(source chan<- interface{}) {
		for _, s := range ss {
			source <- s
		}
	}).Group(func(item interface{}) interface{} {
		// 按照首字符"g"或者"p"分组，没有则分到另一组
		if strings.HasPrefix(item.(string), "g") {
			return "g"
		} else if strings.HasPrefix(item.(string), "p") {
			return "p"
		}
		return ""
	}).ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func TestReverse(t *testing.T) {
	Just(1, 2, 3, 4, 5).Reverse().ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func TestDistinct(t *testing.T) {
	Just(1, 2, 2, 2, 3, 3, 4, 5, 6).Distinct(func(item interface{}) interface{} {
		return item
	}).ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func TestWalk(t *testing.T) {
	Just("aaa", "bbb", "ccc").Walk(func(item interface{}, pipe chan<- interface{}) {
		newItem := strings.ToUpper(item.(string))
		pipe <- newItem
	}).ForEach(func(item interface{}) {
		fmt.Println(item)
	})
}

func TestParallel(t *testing.T) {
	Parallel(func() {
		fmt.Println("aaa")
	}, func() {
		fmt.Println("bbb")
	}, func() {
		fmt.Println("ccc")
	})
}
