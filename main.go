package main

import (
	"github.com/urfave/cli"
	"fmt"
	"os"
	_ "github.com/Sirupsen/logrus"
	container "mydocker/container"
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
	},

	Action: func(context *cli.Context) error {
		fmt.Println("start runCommand")
		if len(context.Args()) < 1 {
			return fmt.Errorf("Missing container command")
		}
		cmd := context.Args().Get(0)
		tty := context.Bool("ti")
		Run(tty, cmd)
		return nil
	},
}

var initCommand = cli.Command {
	Name : "init",
	Usage : "Init the process",

	Action: func(context *cli.Context) error {
		fmt.Println("start initCommand")
		fmt.Println("init come on")
		cmd := context.Args().Get(0)
		fmt.Printf("command %s\n",cmd)
		err := container.RunContainerInitProcess(cmd,nil)
		return err
	},
}


func Run(tty bool, command string){
	parent := container.NewParentProcess(tty, command)
	if err:= parent.Start(); err!=nil {
		fmt.Println("az")
		fmt.Println(err)
	}
	parent.Wait()
	os.Exit(-1)
}

