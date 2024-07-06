package mp

import (
	"errors"
	"fmt"
	"log"
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
	t.Logf("res:{%v}; err:{%v}\n", res, err)
}

func TestFinish(t *testing.T) {

	a := func() error {
		log.Println("aaaa")
		return nil
	}

	b := func() error {
		log.Println("bbbb")
		return errors.New("err about bbbb")
		//return nil
	}

	c := func() error {
		log.Println("cccc")
		return nil
	}

	err := Finish(a, b, c)

	t.Log("finish err:", err)
}
