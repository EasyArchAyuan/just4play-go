package main

import (
	"fmt"
	"sort"
)

// ByLength 为了在 Go 中使用自定义函数进行排序，我们需要一个对应的类型。
// 这里我们创建一个为内置 `[]string` 类型的别名的 `ByLength` 类型。
type ByLength []string

// 实现了 `sort.Interface` 的 `Len`，`Less` 和 `Swap` 方法

func (b ByLength) Len() int {
	return len(b)
}

func (b ByLength) Less(i, j int) bool {
	return len(b[i]) < len(b[j])
}

func (b ByLength) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func main() {
	fruits := []string{"peach", "banana", "kiwi"}
	sort.Sort(ByLength(fruits))
	fmt.Println(fruits)

}
