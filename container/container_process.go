package container

import (
	"syscall"
	"os/exec"
	"os"
	"fmt"
)

func NewParentProcess(tty bool, command string) *exec.Cmd {
	fmt.Println("new parent")
	args := []string{"init",command}
	cmd := exec.Command("/proc/self/exe",args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
		syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET |syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: syscall.Getuid(),
				HostID:      syscall.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: syscall.Getgid(),
				HostID:      syscall.Getgid(),
				Size:        1,
			},
		},
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}