package main

import (
	"fmt"
	"sort"
)

// Orderable 定义一个约束，允许比较整数、浮点数和字符串
type Orderable interface {
	int | int8 | int16 | int32 | int64 |
		float32 | float64 |
		string
}

// SortSliceStable 泛型函数，用于稳定排序切片
func SortSliceStable[T any, V Orderable](slice []T, getV func(T) V, asc bool) {
	less := func(i, j int) bool {
		v1 := getV(slice[i])
		v2 := getV(slice[j])
		if asc {
			return v1 < v2
		}
		return v1 > v2
	}
	sort.SliceStable(slice, less)
}

type StudentX struct {
	Id    int     `sort:"id,asc"`    // 升序
	Name  string  `sort:"name,desc"` // 降序
	Score float64 `sort:"score,asc"`
}

// 实现getV函数，用于获取Student结构体的特定字段
func getId(s StudentX) int {
	return s.Id
}

func getName(s StudentX) string {
	return s.Name
}

func getScore(s StudentX) float64 {
	return s.Score
}

func main() {
	students := []StudentX{
		{1, "zhangsan", 90.0},
		{2, "lisi", 80.0},
		{3, "wangwu", 70.0},
	}

	// 按照Name字段降序排序
	SortSliceStable(students, getName, false)
	fmt.Println("Sorted by Name (desc):", students)

	// 按照Id字段升序排序
	SortSliceStable(students, getId, true)
	fmt.Println("Sorted by Id (asc):", students)

	// 按照Score字段升序排序
	SortSliceStable(students, getScore, true)
	fmt.Println("Sorted by Score (asc):", students)
}
