package main

import (
	"fmt"
	str "strings"
)

// 我们给 `fmt.Println` 一个短名字的别名，我们随后将会经常用到。
var p = fmt.Println

func main() {
	// 这是一些 `strings` 中的函数例子。注意他们都是包中的
	// 函数，不是字符串对象自身的方法，这意味着我们需要考
	// 虑在调用时将字符串作为第一个参数进行传递。
	p("Count:     ", str.Count("test", "t"))
	p("Contains:  ", str.Contains("test", "es"))
	p("HasPrefix: ", str.HasPrefix("test", "te"))
	p("HasSuffix: ", str.HasSuffix("test", "st"))
	p("Index:     ", str.Index("test", "e"))
	p("Join:      ", str.Join([]string{"a", "b"}, "-"))
	p("Repeat:    ", str.Repeat("a", 5))
	p("Replace:   ", str.Replace("foo", "o", "0", -1))
	p("Replace:   ", str.Replace("foo", "o", "0", 1))
	p("Split:     ", str.Split("a-b-c-d-e", "-"))
	p("ToLower:   ", str.ToLower("TEST"))
	p("ToUpper:   ", str.ToUpper("test"))
	p()
	// 虽然不是 `strings` 的一部分，但是仍然值得一提的是获
	// 取字符串长度和通过索引获取一个字符的机制。
	p("Len: ", len("hello"))
	p("Char:", "hello"[1])
}
