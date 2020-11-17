package linuxnamespace

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

/*
Mount namespace
Use mount -t proc proc /proc to check the works
*/
func Mount() {
	cmd := exec.Command("sh")
	fmt.Println("start mountnamespace")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("end mountnamespace")
}
