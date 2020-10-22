package cgroupgo

import (
	"os/exec"
	"path"
	"os"
	"fmt"
	"io/ioutil"
	"syscall"
	"strconv"
)

const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func cgroup(){
	if os.Args[0] == "/proc/self/exe" {
		//cgroup内部
		fmt.Printf("current pid is %d\n",syscall.Getpid())
		//执行stress程序
		cmd := exec.Command("sh","-c","stress --vm-bytes 200m --vm-keep  -m 1")
		cmd.SysProcAttr = &syscall.SysProcAttr{} 
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run();err!=nil {
			fmt.Println("error :")
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags : syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err:= cmd.Start(); err!=nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		//外部空间pid
		fmt.Printf("start docker %v\n",cmd.Process.Pid)
		//创建cgroup
		os.Mkdir(path.Join(cgroupMemoryHierarchyMount,"testmemorylimit"),0755)
		//挂载到cgroup中
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount,"testmemorylimit","tasks"), []byte(strconv.Itoa(cmd.Process.Pid)),0644)  
		//限制cgroup
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount,"testmemorylimit","memory.limit_in_bytes"),[]byte("100m"),0644)
		fmt.Println("time to start")
	}
	cmd.Process.Wait()
}