package linuxnamespace

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

/*
Ipc namespace
use ipcmk -Q and ipcs -q to check
*/
func Ipc() {
	cmd := exec.Command("sh")
	fmt.Println("start ipcnamespace")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("end ipcnamespace")
}
