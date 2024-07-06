package design

import (
	"just4play/design/option"
	"just4play/design/singleton"
	"testing"
)

func TestSingleton(t *testing.T) {
	//单例设计模式
	s := singleton.New()
	s["this"] = "that"
	s2 := singleton.New()
	t.Log("This is", s2["this"])
}

func TestOption(t *testing.T) {
	//选项设计模式 grpc_retry 就是通过这个机制实现的，可以实现自动重试功能。
	option.InitOptions(option.WithStringOption1("Kirin"), option.WithStringOption2("Ayuan"), option.WithStringOption3(5))
}
