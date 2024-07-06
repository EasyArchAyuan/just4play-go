package mp

import (
	"fmt"
	"testing"
	"time"
)

func TrackTime() func() {
	pre := time.Now()

	return func() {
		elapsed := time.Since(pre)
		fmt.Println("耗时:", elapsed)
	}
}

func TestMp(t *testing.T) {
	// 待处理数据
	uid := []int{1, 2, 3, 4, 5, 6}
	// 传递数据
	a := func(source chan<- interface{}) {
		for _, v := range uid {
			source <- v
		}
	}
	// 数据处理
	b := func(item interface{}, writer Writer, cancel func(err error)) {
		tmp := item.(int) + 1
		writer.Writer(tmp)
	}
	// 数据合并
	c := func(pipe <-chan interface{}, writer Writer, cancel func(err error)) {
		var uid []int
		for v := range pipe {
			uid = append(uid, v.(int))
		}
		fmt.Println(uid)
		writer.Writer(uid)
	}
	// 并发调用
	res, err := MapReduce(a, b, c)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
