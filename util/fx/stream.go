package fx

import (
	"just4play/util/collection"
	"just4play/util/lang"
	"just4play/util/thread"
	"sort"
	"sync"
)

const (
	defaultWorkers = 16
	minWorkers     = 1
)

type (
	rxOptions struct {
		unlimitedWorkers bool //是否使用无限数量的工作者
		workers          int  //指定的工作者数量
	}

	// FilterFunc 用于过滤流中的元素。它接收一个元素并返回一个布尔值，用于指示是否保留该元素
	FilterFunc func(item any) bool
	// ForAllFunc 用于处理流中的所有元素。它接收一个通道，可以从该通道中接收元素进行处理
	ForAllFunc func(pipe <-chan any)
	// ForEachFunc 用于处理流中的每个元素。它接收一个元素作为参数
	ForEachFunc func(item any)
	// GenerateFunc 用于向流中发送元素。它接收一个通道，可以向该通道发送元素
	GenerateFunc func(source chan<- any)
	// KeyFunc 用于为流中的元素生成键。它接收一个元素作为参数，并返回一个键
	KeyFunc func(item any) any
	// LessFunc 用于比较流中元素的大小。它接收两个元素作为参数，并返回一个布尔值，用于指示第一个元素是否小于第二个元素
	LessFunc func(a, b any) bool
	// MapFunc 用于将流中的每个元素映射到另一个对象。它接收一个元素作为参数，并返回一个映射后的元素
	MapFunc func(item any) any
	// Option 用于自定义流的选项。它接收一个 rxOptions 指针作为参数
	Option func(opts *rxOptions)
	// ParallelFunc 用于并行处理流中的元素。它接收一个元素作为参数
	ParallelFunc func(item any)
	// ReduceFunc 用于将流中的所有元素进行归约。它接收一个通道作为参数，并返回一个归约后的结果和可能的错误
	ReduceFunc func(pipe <-chan any) (any, error)
	// WalkFunc 用于遍历流中的所有元素。它接收一个元素和一个通道作为参数，可以将处理后的元素发送到通道中
	WalkFunc func(item any, pipe chan<- any)

	// Stream 定义了一个流，它包含一个源通道，用于从源通道中接收元素进行处理
	Stream struct {
		source <-chan any
	}
)

// Concat 合并Stream
func Concat(s Stream, others ...Stream) Stream {
	return s.Concat(others...)
}

// From 通过From函数构建流并返回Stream，流数据通过channel进行存储
func From(generate GenerateFunc) Stream {
	source := make(chan any)

	thread.SafeGoroutine(func() {
		defer close(source)
		// 构造流数据写入channel
		generate(source)
	})

	return Range(source)
}

// Just 批量传入元素构建stream流
func Just(items ...any) Stream {
	source := make(chan any, len(items))
	for _, item := range items {
		source <- item
	}
	close(source)

	return Range(source)
}

// Range converts the given channel to a Stream.
func Range(source <-chan any) Stream {
	return Stream{
		source: source,
	}
}

// AllMach 一个stream中的所有元素是否都满足给定的条件
func (s Stream) AllMach(predicate func(item any) bool) bool {
	for item := range s.source {
		if !predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return false
		}
	}

	return true
}

// AnyMach 一个Stream中是否存在任何满足给定条件的元素
func (s Stream) AnyMach(predicate func(item any) bool) bool {
	for item := range s.source {
		if predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return true
		}
	}

	return false
}

// Buffer Stream中的元素缓存到一个大小为n的队列
// 平衡生产者和消费者之间的处理吞吐量不匹配问题
func (s Stream) Buffer(n int) Stream {
	if n < 0 {
		n = 0
	}

	source := make(chan any, n)
	go func() {
		for item := range s.source {
			source <- item
		}
		close(source)
	}()

	return Range(source)
}

// Concat 合并流
func (s Stream) Concat(others ...Stream) Stream {
	source := make(chan any)

	go func() {
		group := thread.NewRoutineGroup()
		group.Run(func() {
			for item := range s.source {
				source <- item
			}
		})

		for _, each := range others {
			each := each
			group.Run(func() {
				for item := range each.source {
					source <- item
				}
			})
		}

		group.Wait()
		close(source)
	}()

	return Range(source)
}

// Count 统计流中元素个数
func (s Stream) Count() (count int) {
	for range s.source {
		count++
	}
	return
}

// Distinct distinct对流中元素进行去重，去重在业务开发中比较常用，经常需要对用户id等做去重操作
func (s Stream) Distinct(fn KeyFunc) Stream {
	source := make(chan any)

	thread.SafeGoroutine(func() {
		defer close(source)
		// 通过key进行去重，相同key只保留一个
		keys := make(map[any]lang.PlaceholderType)
		for item := range s.source {
			key := fn(item)
			// key存在则不保留
			if _, ok := keys[key]; !ok {
				source <- item
				keys[key] = lang.Placeholder
			}
		}
	})

	return Range(source)
}

// Done 是等待所有上游操作完成
func (s Stream) Done() {
	drain(s.source)
}

// Filter 过滤不满足条件的item
func (s Stream) Filter(fn FilterFunc, opts ...Option) Stream {
	return s.Walk(func(item any, pipe chan<- any) {
		if fn(item) {
			pipe <- item
		}
	}, opts...)
}

// First 返回第一个元素
func (s Stream) First() any {
	for item := range s.source {
		// make sure the former goroutine not block, and current func returns fast.
		go drain(s.source)
		return item
	}

	return nil
}

// ForAll handles the streaming elements from the source and no later streams.
func (s Stream) ForAll(fn ForAllFunc) {
	fn(s.source)
	// avoid goroutine leak on fn not consuming all items.
	go drain(s.source)
}

// ForEach 遍历流中所有元素
func (s Stream) ForEach(fn ForEachFunc) {
	for item := range s.source {
		fn(item)
	}
}

// Group Group对流数据进行分组，需定义分组的key，数据分组后以slice存入channel:
func (s Stream) Group(fn KeyFunc) Stream {
	// 定义分组存储map
	groups := make(map[any][]any)
	for item := range s.source {
		// 用户自定义分组key
		key := fn(item)
		// key相同分到一组
		groups[key] = append(groups[key], item)
	}

	source := make(chan any)
	go func() {
		for _, group := range groups {
			// 相同key的一组数据写入到channel
			source <- group
		}
		close(source)
	}()

	return Range(source)
}

// Head 取出前n个item，返回新stream
func (s Stream) Head(n int64) Stream {
	if n < 1 {
		panic("n must be greater than 0")
	}

	source := make(chan any)

	go func() {
		for item := range s.source {
			n--
			if n >= 0 {
				source <- item
			}
			if n == 0 {
				// let successive method go ASAP even we have more items to skip
				close(source)
				// why we don't just break the loop, and drain to consume all items.
				// because if breaks, this former goroutine will block forever,
				// which will cause goroutine leak.
				drain(s.source)
			}
		}
		// not enough items in s.source, but we need to let successive method to go ASAP.
		if n > 0 {
			close(source)
		}
	}()

	return Range(source)
}

// Last 返回最后一个元素
func (s Stream) Last() (item any) {
	for item = range s.source {
	}
	return
}

// Map 对象转换
func (s Stream) Map(fn MapFunc, opts ...Option) Stream {
	return s.Walk(func(item any, pipe chan<- any) {
		pipe <- fn(item)
	}, opts...)
}

// Max 返回Stream中item的最大值
func (s Stream) Max(less LessFunc) any {
	var max any
	for item := range s.source {
		if max == nil || less(max, item) {
			max = item
		}
	}

	return max
}

// Merge 合并item到slice并生成新stream
func (s Stream) Merge() Stream {
	var items []any
	for item := range s.source {
		items = append(items, item)
	}

	source := make(chan any, 1)
	source <- items
	close(source)

	return Range(source)
}

// Min 返回Stream中item的最小值
func (s Stream) Min(less LessFunc) any {
	var min any
	for item := range s.source {
		if min == nil || less(item, min) {
			min = item
		}
	}

	return min
}

// NoneMatch Stream中的所有元素是否都不满足给定的条件
func (s Stream) NoneMatch(predicate func(item any) bool) bool {
	for item := range s.source {
		if predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return false
		}
	}

	return true
}

// Parallel applies the given ParallelFunc to each item concurrently with given number of workers.
func (s Stream) Parallel(fn ParallelFunc, opts ...Option) {
	s.Walk(func(item any, pipe chan<- any) {
		fn(item)
	}, opts...).Done()
}

// Reduce 汇总
func (s Stream) Reduce(fn ReduceFunc) (any, error) {
	return fn(s.source)
}

// Reverse reverse可以对流中元素进行反转处理
func (s Stream) Reverse() Stream {
	var items []any
	// 获取流中数据
	for item := range s.source {
		items = append(items, item)
	}
	// 反转算法
	for i := len(items)/2 - 1; i >= 0; i-- {
		opp := len(items) - 1 - i
		items[i], items[opp] = items[opp], items[i]
	}
	// 写入流
	return Just(items...)
}

// Skip 跳过前n个item，返回新stream
func (s Stream) Skip(n int64) Stream {
	if n < 0 {
		panic("n must not be negative")
	}
	if n == 0 {
		return s
	}

	source := make(chan any)

	go func() {
		for item := range s.source {
			n--
			if n >= 0 {
				continue
			} else {
				source <- item
			}
		}
		close(source)
	}()

	return Range(source)
}

// Sort  对item进行排序
func (s Stream) Sort(less LessFunc) Stream {
	var items []any
	for item := range s.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})

	return Just(items...)
}

// Split 分割对流数据进行分割
func (s Stream) Split(n int) Stream {
	if n < 1 {
		panic("n should be greater than 0")
	}

	source := make(chan any)
	go func() {
		var chunk []any
		for item := range s.source {
			chunk = append(chunk, item)
			if len(chunk) == n {
				source <- chunk
				chunk = nil
			}
		}
		if chunk != nil {
			source <- chunk
		}
		close(source)
	}()

	return Range(source)
}

// Tail 与Head功能类似，取出后n个item组成新stream
func (s Stream) Tail(n int64) Stream {
	if n < 1 {
		panic("n should be greater than 0")
	}

	source := make(chan any)

	go func() {
		ring := collection.NewRing(int(n))
		for item := range s.source {
			ring.Add(item)
		}
		for _, item := range ring.Take() {
			source <- item
		}
		close(source)
	}()

	return Range(source)
}

// Walk Walk函数并发的作用在流中每一个item上，可以通过WithWorkers设置并发数，默认并发数为16，
// 最小并发数为1，如设置unlimitedWorkers为true则并发数无限制，但并发写入流中的数据由defaultWorkers限制，
// WalkFunc中用户可以自定义后续写入流中的元素，可以不写入也可以写入多个元素
func (s Stream) Walk(fn WalkFunc, opts ...Option) Stream {
	option := buildOptions(opts...)
	if option.unlimitedWorkers {
		return s.walkUnlimited(fn, option)
	}

	return s.walkLimited(fn, option)
}

func (s Stream) walkLimited(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan any, option.workers)

	go func() {
		var wg sync.WaitGroup
		pool := make(chan lang.PlaceholderType, option.workers)

		for item := range s.source {
			// important, used in another goroutine
			val := item
			pool <- lang.Placeholder
			wg.Add(1)

			// better to safely run caller defined method
			thread.SafeGoroutine(func() {
				defer func() {
					wg.Done()
					<-pool
				}()

				fn(val, pipe)
			})
		}

		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}

func (s Stream) walkUnlimited(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan any, option.workers)

	go func() {
		var wg sync.WaitGroup

		for item := range s.source {
			// important, used in another goroutine
			val := item
			wg.Add(1)
			// better to safely run caller defined method
			thread.SafeGoroutine(func() {
				defer wg.Done()
				fn(val, pipe)
			})
		}

		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}

// UnlimitedWorkers lets the caller use as many workers as the tasks.
func UnlimitedWorkers() Option {
	return func(opts *rxOptions) {
		opts.unlimitedWorkers = true
	}
}

// WithWorkers lets the caller customize the concurrent workers.
func WithWorkers(workers int) Option {
	return func(opts *rxOptions) {
		if workers < minWorkers {
			opts.workers = minWorkers
		} else {
			opts.workers = workers
		}
	}
}

// buildOptions returns a rxOptions with given customizations.
func buildOptions(opts ...Option) *rxOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

// drain 清空给定的channel
func drain(channel <-chan any) {
	for range channel {
	}
}

// newOptions 返回一个默认的rxOptions指针
func newOptions() *rxOptions {
	return &rxOptions{
		workers: defaultWorkers,
	}
}
