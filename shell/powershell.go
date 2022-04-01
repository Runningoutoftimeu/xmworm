package shell

import (
	"bytes"
	"os/exec"
)

// PowerShell struct
type PowerShell struct {
	pathname string
}

func (sh *PowerShell) init() {
	ps, _ := exec.LookPath("powershell.exe")
	sh.pathname = ps
}

func (p *PowerShell) execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.pathname, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}