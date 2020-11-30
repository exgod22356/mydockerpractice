package main

import (
	"fmt"
	"io/ioutil"
	cgroupmanager "mydocker/cgroupmanager"
	container "mydocker/container"
	subsystem "mydocker/subsystem"
	"os"
	"os/exec"
	"strings"

	_ "github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `mydocker is a simple container.`

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
		commitCommand,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

var runCommand = cli.Command{
	Name:  "run",
	Usage: "Create a container",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "limit the memory",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
	},

	Action: func(context *cli.Context) error {
		fmt.Println("start runCommand")
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("ti")
		resConf := &subsystem.ResourceConfig{
			MemoryLimit: context.String("m"),
		}
		volume := context.String("v")
		Run(tty, cmdArray, resConf, volume)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init the process",

	Action: func(context *cli.Context) error {
		fmt.Println("start initCommand")
		cmd := readUserCommand()
		err := container.RunContainerInitProcess(cmd[0], cmd[1:])
		return err
	},
}

var commitCommand = cli.Command{
	Name:  "commit",
	Usage: "commit a container image",
	Action: func(context *cli.Context) error {
		fmt.Println("start commit")
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container name")

		}
		imageName := context.Args().Get(0)
		if err := commitContainer(imageName); err != nil {
			return fmt.Errorf("error in commitContainer: %v", err)
		}
		return nil
	},
}

/*
Run command
start a parentprocess with namespace
start a CgroupManager
set the config
store the command
*/
func Run(tty bool, commandArray []string, resConf *subsystem.ResourceConfig, volume string) {
	parent, writePipe := container.NewParentProcess(tty, volume)
	if err := parent.Start(); err != nil {
		fmt.Println(err)
		return
	}
	cgroupManager := cgroupmanager.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(resConf)
	cgroupManager.Apply(parent.Process.Pid)
	sendInitCommand(commandArray, writePipe)
	parent.Wait()
	mntURL := "/home/wqy/mnt/"
	rootURL := "/home/wqy/"
	fmt.Println("start cleaning")
	if err := container.DeleteWorkSpace(rootURL, mntURL, volume); err != nil {
		fmt.Println(err)
	}
	os.Exit(0)
}

//Write the commands
func sendInitCommand(commandArray []string, writePipe *os.File) {
	command := strings.Join(commandArray, " ")
	fmt.Printf("your command is %v\n", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

//Hard code fd, to read the command
func readUserCommand() []string {
	fmt.Println("reading user command")
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		fmt.Println(err)
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

func commitContainer(imageName string) error {
	mntURL := "/home/wqy/mnt/"
	rootURL := "/home/wqy/"
	imageTar := rootURL + imageName + ".tar"
	fmt.Println(imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		return fmt.Errorf("tar error: %v ", err)
	}
	return nil
}
