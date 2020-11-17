package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

/*
RunContainerInitProcess function:
mounts the essential environment,
read the command stored by the NewParentProcess,
run the commands
*/
func RunContainerInitProcess(command string, args []string) error {
	fmt.Printf("the command is %s\n", command)
	fmt.Println("mount start")
	err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	command, err = exec.LookPath(command)
	if err != nil {
		fmt.Printf("error in finding %s \n", command)
		fmt.Println(err)
		return err
	}
	argv := append([]string{command}, args...)
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		fmt.Printf("exec error is %v\n", err)
		return err
	}
	return nil
}
