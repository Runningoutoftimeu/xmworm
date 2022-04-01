package shell

import (
	"testing"
	"reflect"
	"strings"
)

func TestInit(t *testing.T) {
	want := &PowerShell{`C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe`}
	got := &PowerShell{}
	got.init()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}

}

func TestExecute(t *testing.T) {
	want := "64-bit"

	psh := &PowerShell{}
	psh.init()
	stdOut, _, err := psh.execute(`(Get-WmiObject win32_operatingsystem).osarchitecture`) 
	if err != nil {
		panic(err)
	}

	got := strings.TrimSpace(stdOut)
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}