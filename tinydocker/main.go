package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {

	switch os.Args[1] {
	case "run":
		// proc/self/exe 是一个软链接，程序内部读取到的链接是自身可执行文件的路径
		initCmd, err := os.Readlink("/proc/self/exe")
		if err != nil {
			fmt.Println("get init process error ", err)
			return
		}

		os.Args[1] = "init"
		cmd := exec.Command(initCmd, os.Args[1:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWUTS,
		}
		cmd.Env = os.Environ()
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("init proc end", initCmd)
	case "init":
		cmd := os.Args[2]
		err := syscall.Exec(cmd, os.Args[2:], os.Environ())
		if err != nil {
			fmt.Println("exec proc fail ", err)
			return
		}
		fmt.Println("forever exec it ")
	default:
		fmt.Println("not valid cmd")
	}

}
