package container

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

/*
NewParentProcess function:
creates a child precess with namespaces,
run itself in the process,
create a pipe and file for the Run fucntion to write
*/
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	fmt.Println("New parent")
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		fmt.Printf("pipe error %v\n", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET, /*| syscall.CLONE_NEWUSER,
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
		},*/
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	mntURL := "/home/wqy/mnt/"
	rootURL := "/home/wqy/"
	if err := NewWorkSpace(rootURL, mntURL); err != nil {
		fmt.Println(err)
		return nil, nil
	}
	cmd.Dir = mntURL
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}
