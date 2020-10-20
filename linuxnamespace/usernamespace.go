package linuxnamespace

import (
	"os/exec"
	"syscall"
	"os"
	"log"
	"fmt")

/*
User namespace  
Use id to check the works
*/
func User(){
	cmd:=exec.Command("sh")
	fmt.Println("start usernamespace")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS|syscall.CLONE_NEWIPC|syscall.CLONE_NEWPID|syscall.CLONE_NEWNS|syscall.CLONE_NEWUSER,
	}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid:uint32(1),Gid:uint32(1)}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	
	if err:=cmd.Run();err!=nil {
		log.Fatal(err)
	}
	fmt.Println("end usernamespace")
}