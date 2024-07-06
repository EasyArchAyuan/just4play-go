package mp

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

//实际业务场景中多个依赖如果有一个出错我们期望能立即返回而不是等所有依赖都执行完再返回结果

var (
	cancelWithNil = errors.New("reduce end with nil error")
)

/**
* 1.source channel 无缓冲chan，generate 生成数据写入source channel，mapper读取source channel数据
* 2.collector channel 有缓冲chan,mapper处理完写入collector channel,reducer读取collector channel数据
* 3.output channel 无缓冲chan,reducer处理完写入output channel
* 4.done channel 无缓冲chan,执行异样时候通知其他goroutine退出
**/
type (
	GenerateFunc func(source chan<- interface{})

	MapFunc    func(item interface{}, writer Writer)
	MapperFunc func(item interface{}, writer Writer, cancel func(err error))

	ReducerFunc     func(pipe <-chan interface{}, writer Writer, cancel func(err error))
	VoidReducerFunc func(pipe <-chan interface{}, cancel func(error))

	Writer interface {
		Writer(val interface{})
	}
)

type writeChan struct {
	write chan interface{}
	done  chan struct{}
}

func newWriteChan(write chan interface{}, done chan struct{}) writeChan {
	return writeChan{
		write: write,
		done:  done,
	}
}

func (w writeChan) Writer(val interface{}) {
	select {
	case <-w.done:
		return
	default:
		w.write <- val
	}
}

func (w writeChan) Load() interface{} {
	return <-w.write
}

// 通过传入的generate方法产生数据，并返回source提供给mapper读取
func buildSource(generate GenerateFunc) chan interface{} {
	source := make(chan interface{})
	go func() {
		defer func() {
			close(source)
		}()
		generate(source)
	}()
	return source
}

// 消费generate产生的数据，并写入collector，mapper默认最大并发数为16
func executeMappers(mapper MapFunc, collector chan interface{}, done chan struct{}, source <-chan interface{}) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
		//只关collector就行，done可能已经关闭了
		close(collector)
	}()

	writer := newWriteChan(collector, done)
	pool := make(chan struct{}, 16)
	for {
		select {
		case <-done:
			return
		case pool <- struct{}{}:
			item, ok := <-source
			if !ok {
				//说明管道已关闭
				<-pool
				return
			}

			wg.Add(1)
			go func() {
				defer func() {
					wg.Done()
					// 在这里关闭, 以保证最多有 16 个在进行
					<-pool
				}()
			}()
			//运行自定义处理函数
			mapper(item, writer)
		}
	}
}

// MapReduce 并发执行任务
func MapReduce(generate GenerateFunc, mapper MapperFunc, reducer ReducerFunc) (interface{}, error) {
	source := buildSource(generate)
	var errVal atomic.Value
	done := make(chan struct{})
	reduceChan := make(chan interface{})
	var (
		cancelOnce  sync.Once
		reduceClose sync.Once
	)
	finish := func() {
		cancelOnce.Do(func() {
			close(done)
		})
		reduceClose.Do(func() {
			close(reduceChan)
		})
	}

	cancel := func(err error) {

		if err != nil {
			errVal.Store(err)
		} else {
			errVal.Store(cancelWithNil)
		}
		defer func() {
			// 把资源管道里的数据清空
			drain(source)
		}()

		finish()
	}
	write := newWriteChan(reduceChan, done)
	// 存放处理结果, mapper 往这个里面写入， reduce 从这个里面读出
	resChan := make(chan interface{})
	go func() {
		defer func() {
			finish()
			drain(resChan)
		}()
		// 在这里可能遇到错误就结束运行了, reschan 可能还有数据, 所以要在 defer 中把数据都给读取完
		reducer(resChan, write, cancel)
	}()

	// 现在开始从执行管道里读取数据处理
	go executeMappers(func(item interface{}, writer Writer) {
		mapper(item, writer, cancel)
	}, resChan, done, source)

	// 此时我们应该取出错误 和 结果
	res, ok := <-reduceChan
	errIntF := errVal.Load()
	var err error
	if errIntF != nil {
		err = errIntF.(error)
		return nil, err
	}
	if !ok {
		return nil, err
	} else {
		return res, nil
	}

}

func drain(channel <-chan interface{}) {
	for range channel {

	}
}

func Finish(fns ...func() error) error {
	_, err := MapReduce(func(source chan<- interface{}) {
		for _, fn := range fns {
			source <- fn
		}
	}, func(item interface{}, writer Writer, cancel func(err error)) {
		fmt.Println(item)
		f := item.(func() error)
		if err := f(); err != nil {
			cancel(err)
		}
	}, func(pipe <-chan interface{}, writer Writer, cancel func(err error)) {
		drain(pipe)
	})
	return err
}
