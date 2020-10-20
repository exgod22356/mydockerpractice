package linuxnamespace

import (
	"os/exec"
	"syscall"
	"os"
	"log"
	"fmt")

/*
Uts namespace
Use hosename -b xxx to change the hostname
*/
func Uts(){
	fmt.Println("start uts namespace")
	cmd:=exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	
	if err:=cmd.Run();err!=nil {
		log.Fatal(err)
	}
	fmt.Println("end utsnamespace")
}
