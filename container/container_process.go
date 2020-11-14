package container

import (
	"syscall"
	"os/exec"
	"os"
	"fmt"
)

/*
NewParentProcess function:
creates a child precess with namespaces, 
run itself in the process, 
store the command, 
handle it to the container init process
*/
func NewParentProcess(tty bool) (*exec.Cmd, *os.File){
	fmt.Println("New parent")
	readPipe, writePipe, err := NewPipe()
	if err!=nil {
		fmt.Printf("pipe error %v\n",err)
		return nil,nil
	}
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