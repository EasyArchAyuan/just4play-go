package main

import "fmt"

func vals() (int, int) {
	return 3, 7
}

func main() {
	i, i2 := vals()
	fmt.Println(i, i2)

	_, c := vals()
	fmt.Println(c)
}
