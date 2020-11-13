package container

import (
	"fmt"
	"syscall"
	"os"
	"os/exec"
)


func RunContainerInitProcess(command string, args []string) error{
	fmt.Printf("the command is %s\n",command)
	fmt.Println("mount start")
	err :=syscall.Mount("","/","",syscall.MS_PRIVATE|syscall.MS_REC,"")
	if err!=nil {
		fmt.Println(err)
		return err
	}
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc","proc", uintptr(defaultMountFlags),"")
	if err!=nil {
		fmt.Println(err)
	}
	command,err = exec.LookPath(command)
	if(err!=nil){
		fmt.Printf("error in finding %s \n",command)
		fmt.Println(err)
		return nil;
	}
	argv := append([]string{command},args...)
	fmt.Println(argv)
	if err := syscall.Exec(command, argv, os.Environ()) ; err!=nil {
		fmt.Printf("mount error is %v\n",err)
	}
	fmt.Println("mount over")
	return nil
}