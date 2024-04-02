package main

import (
	"fmt"
	"os"
)

func main() {

	f := createFile("E:/temp/defer.txt")
	//类似finally
	defer closeFile(f)
	writeFile(f)
}

func closeFile(f *os.File) {
	fmt.Println("closing")
	err := f.Close()
	if err != nil {
		return
	}
}

func writeFile(f *os.File) {
	fmt.Println("writing")
	_, err := fmt.Fprintln(f, "data")
	if err != nil {
		return
	}
}

func createFile(p string) *os.File {
	fmt.Println("creating file :", p)
	f, err := os.Create(p)
	if err != nil {
		panic(err)
	}
	return f
}
