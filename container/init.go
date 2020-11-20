package container

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	err := setUpMount()
	if err != nil {
		fmt.Println("setUpMount error")
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

func setUpMount() error {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("get wd error: %v\n", err)
		return err
	}
	fmt.Println(pwd)
	if err = pivotRoot(pwd); err != nil {
		fmt.Printf("pivot root error: %v\n", err)
		return err
	}

	err = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		fmt.Printf("mount / error: %v\n", err)
		return err
	}
	pwd, err = os.Getwd()
	fmt.Printf("the current wd is %v\n", pwd)
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		fmt.Printf("mount proc error: %v\n", err)
		return err
	}

	err = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		fmt.Printf("mount / error: %v\n", err)
		return err
	}
	/*fmt.Println("-------------------")
	files, _ := ioutil.ReadDir("./tmpfs")
	for _, f := range files {
		fmt.Println(f.Name())
	}
	fmt.Println("-------------------")*/
	/*err = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=0755")
	if err != nil {
		fmt.Printf("mount tmpfs error: %v\n", err)
		return err
	}*/
	return nil
}

func pivotRoot(root string) error {
	fmt.Println("pivot rooting")
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount bind error: %v", err)
	}

	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("mkdir error: %v", err)
	}

	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		fmt.Printf("mount / error: %v\n", err)
		return err
	}
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivotroot error: %v", err)
	}

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir error: %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("umount error: %v", err)
	}
	return os.Remove(pivotDir)
}
