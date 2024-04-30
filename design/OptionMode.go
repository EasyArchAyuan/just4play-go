package main

import "fmt"

// 结构体
type Options struct {
	str1 string
	str2 string
	int1 int
	int2 int
}

// 传参用
type Option func(*Options)

func InitOptions(opts ...Option) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}
	fmt.Printf("options:%#v\n", options)
}

func WithStringOption1(str string) Option {
	return func(opts *Options) {
		opts.str1 = str
	}
}

func WithStringOption2(str string) Option {
	return func(opts *Options) {
		opts.str2 = str
	}
}
func WithStringOption3(int1 int) Option {
	return func(opts *Options) {
		opts.int1 = int1
	}
}
func WithStringOption4(int1 int) Option {
	return func(opts *Options) {
		opts.int2 = int1
	}
}

// 选项设计模式
func main() {
	InitOptions(WithStringOption1("Kirin"), WithStringOption2("Ayuan"), WithStringOption3(5), WithStringOption4(6))
	//grpc_retry 就是通过这个机制实现的，可以实现自动重试功能。
}
