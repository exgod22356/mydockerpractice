package main

import (
	"github.com/urfave/cli"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	_ "github.com/Sirupsen/logrus"
	container "mydocker/container"
	subsystem "mydocker/subsystem"
	cgroupmanager "mydocker/cgroupmanager"
)

const usage = `mydocker is a simple container.`

func main(){
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage
	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}
	
	if err:= app.Run(os.Args); err!=nil {
		fmt.Println(err)
	}	
}

var runCommand = cli.Command{
	Name : "run",
	Usage : "Create a container",
	Flags : []cli.Flag{
		cli.BoolFlag{
			Name : "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name : "m",
			Usage : "limit the memory",
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
		Run(tty, cmdArray, resConf)
		return nil
	},
}

var initCommand = cli.Command {
	Name : "init",
	Usage : "Init the process",

	Action: func(context *cli.Context) error {
		fmt.Println("start initCommand")
		cmd:=readUserCommand()
		err := container.RunContainerInitProcess(cmd[0],cmd[1:])
		return err
	},
}

/*
Run command
start a CgroupManager
set the config
store the command
*/
func Run(tty bool, commandArray []string, resConf *subsystem.ResourceConfig){
	parent, writePipe := container.NewParentProcess(tty)
	if err:= parent.Start(); err!=nil {
		fmt.Println(err)
		return
	}
	cgroupManager := cgroupmanager.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(resConf)
	cgroupManager.Apply(parent.Process.Pid)
	sendInitCommand(commandArray, writePipe)
	parent.Wait()
}

func sendInitCommand(commandArray []string, writePipe *os.File) {
	command := strings.Join(commandArray, " ")
	fmt.Printf("your command is %v\n",command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3),"pipe")
	msg,err := ioutil.ReadAll(pipe)
	if err!=nil {
		fmt.Println(err)
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
