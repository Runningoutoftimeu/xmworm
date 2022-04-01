package shell

// import (
// 	"fmt"
// 	"strings"
// )


type Shell interface {
	init()
	execute(args ...string) (stdOut string, stdErr string, err error)
}


func Run(cmd string, shellType string) (error){
	var shell Shell
	if shellType == "powershell" {
		shell = &PowerShell{}
	}else if shellType == "cmd" {
		// do cmd stuff
	}
	shell.init()
	//stdOut, stdErr, err := shell.execute(cmd)
	_, _, err := shell.execute(cmd)

	//fmt.Printf("\nRan::LnkFile.Cmd():\nStdOut : '%s'\nStdErr: '%s'\nErr: %s", strings.TrimSpace(stdOut), stdErr, err)
	return err
}