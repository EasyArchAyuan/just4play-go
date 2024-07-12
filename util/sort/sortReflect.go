package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type Student struct {
	Id    int     `sort:"id,asc"`    // 升序
	Name  string  `sort:"name,desc"` // 降序
	Score float64 `sort:"score,asc"`
}

func SortByField(slice any, sortField string) error {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return fmt.Errorf("SortSliceStable called with non-slice type")
	}
	if rv.Len() == 0 {
		return nil
	}
	// 获取结构体的元素类型
	elemType := rv.Type().Elem()
	// 获取排序字段
	field, ok := elemType.FieldByName(sortField)
	if !ok {
		return fmt.Errorf("SortSliceStable called with unknown or unexported struct field name: " + sortField)
	}
	// 获取 less 函数
	less := func(i, j int) bool {
		v1 := rv.Index(i).FieldByName(field.Name)
		v2 := rv.Index(j).FieldByName(field.Name)
		if !v1.IsValid() || !v2.IsValid() {
			return false
		}
		switch v1.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return v1.Int() < v2.Int()
		case reflect.Float32, reflect.Float64:
			return v1.Float() < v2.Float()
		case reflect.String:
			return v1.String() < v2.String()
		default:
			return false
		}
	}
	// 使用 sort.SliceStable 进行排序
	sort.SliceStable(slice, less)
	return nil
}

func SortByTag(slice any, sortTag string) error {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return fmt.Errorf("SortSliceStable called with non-slice type")
	}
	if rv.Len() == 0 {
		return nil
	}

	elemType := rv.Type().Elem()
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("sort")
		if strings.HasPrefix(tag, sortTag) {
			// 解析排序方向
			parts := strings.Split(tag, ",")
			sortDirection := "asc" // 默认为升序
			if len(parts) > 1 {
				sortDirection = parts[1]
			}

			// 获取 less 函数
			less := func(i, j int) bool {
				v1 := rv.Index(i).Field(i)
				v2 := rv.Index(j).Field(i)
				if !v1.IsValid() || !v2.IsValid() {
					return false
				}
				switch v1.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if sortDirection == "asc" {
						return v1.Int() < v2.Int()
					} else {
						return v1.Int() > v2.Int()
					}
				case reflect.Float32, reflect.Float64:
					if sortDirection == "asc" {
						return v1.Float() < v2.Float()
					} else {
						return v1.Float() > v2.Float()
					}
				case reflect.String:
					if sortDirection == "asc" {
						return v1.String() < v2.String()
					} else {
						return v1.String() > v2.String()
					}
				default:
					return false
				}
			}
			// 使用 sort.SliceStable 进行排序
			sort.SliceStable(slice, less)
			return nil
		}
	}
	return fmt.Errorf("SortSliceStable called with unknown or unexported sort tag: " + sortTag)
}
func main() {
	students := []Student{
		{1, "zhangsan", 90.0},
		{2, "lisi", 80.0},
		{3, "wangwu", 70.0},
	}

	// 错误处理的示例
	err := SortByTag(students, "name")
	if err != nil {
		fmt.Println("排序错误:", err)
		return
	}
	fmt.Println(students)

	//// 其他排序逻辑保持不变
	//err = SortByField(students, "id")
	//if err != nil {
	//	fmt.Println("排序错误:", err)
	//	return
	//}
	//fmt.Println(students)
	//
	//err = SortByField(students, "score")
	//if err != nil {
	//	fmt.Println("排序错误:", err)
	//	return
	//}
	//fmt.Println(students)
}
