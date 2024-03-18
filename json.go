package main

// 下面我们将使用这两个结构体来演示自定义类型的编码和解码。

type Response1 struct {
	Page   int
	Fruits []string
}
type Response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

func main() {
	
}
