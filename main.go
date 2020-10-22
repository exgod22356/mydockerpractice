package main

import (
	"os/exec"
	"fmt"
	"os"
	"time"
)



func main(){
	if os.Args[0] == "/proc/self/exe" {
		fmt.Println("restart process")
		fmt.Println("sleep")
		time.Sleep(time.Second*2)
	}
	fmt.Println("start process")
	cmd := exec.Command("/proc/self/exe")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err:=cmd.Run();err!=nil{
		fmt.Println(err)
		os.Exit(-1)
	}else {
		fmt.Println("??? this should not be printed")
	}
	cmd.Process.Wait()
}