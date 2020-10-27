package container

import (
	"fmt"
	"syscall"
	"os"
)


func RunContainerInitProcess(command string, args []string) error{
	fmt.Printf("command %s\n",command)
	fmt.Println("mount start")
	syscall.Mount("","/","",syscall.MS_PRIVATE|syscall.MS_REC,"")
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc","proc", uintptr(defaultMountFlags),"")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()) ; err!=nil {
		fmt.Println(err)
	}
	return nil
}