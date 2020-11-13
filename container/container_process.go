package container

import (
	"syscall"
	"os/exec"
	"os"
	"fmt"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File){
	fmt.Println("new parent")
	readPipe, writePipe, err := NewPipe()
	if err!=nil {
		fmt.Println("pipe error %v",err)
		return nil,nil
	}
	//args := []string{"init",}
	cmd := exec.Command("/proc/self/exe","init")
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
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}

func NewPipe()(*os.File, *os.File, error){
	read, write, err := os.Pipe()
	if err!=nil {
		return nil, nil, err
	}
	return read, write, nil
}