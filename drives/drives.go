package drives

import (
	"golang.org/x/sys/windows"
	"io/ioutil"
)

type Drive struct {
	name string
	removable bool
	files []string
}

func (d Drive) Name() (string) {
	return d.name
}

type WindowsDriveScanner interface {
	GetLogicalDrives()(drivesBitMask uint32, err error)
}

type DefaultWindowsDriveScanner struct {}

func (sc DefaultWindowsDriveScanner) GetLogicalDrives() (drivesBitMask uint32, err error) {
	drivesBitMask, err = windows.GetLogicalDrives()
	return
}

type DriveTypeDetector interface {
	GetDriveType(drive *Drive)
}

type DefaultWindowsDriveTypeDetector struct {}

func (d DefaultWindowsDriveTypeDetector) GetDriveType(drive *Drive) {
	driveType := windows.GetDriveType(windows.StringToUTF16Ptr(drive.name))
	if driveType == windows.DRIVE_REMOVABLE {
		drive.removable = true
	}else {
		drive.removable = false
	}
}



func NewDrive(driveName string) Drive {
	return Drive{name: driveName, removable: false}
}

func bitsToDrives(bitMap uint32) (drives []string) {
    availableDrives := []string{`A:\`, `B:\`, `C:\`, `D:\`, `E:\`, `F:\`, `G:\`, `H:\`, `I:\`, `J:\`, `K:\`, `L:\`, `M:\`, `N:\`, `O:\`, `P:\`, `Q:\`, `R:\`, `S:\`, `T:\`, `U:\`, `V:\`, `W:\`, `X:\`, `Y:\`, `Z:\`}

    for i := range availableDrives {
        if bitMap & 1 == 1 {
            drives = append(drives, availableDrives[i])
        }
        bitMap >>= 1
    }
    return drives
}

func ScanDrives(sc WindowsDriveScanner) (drives []Drive) {

	drivesBitMask, err := sc.GetLogicalDrives()
	if err != nil {
		panic(err)
	}
	drivesNames := bitsToDrives(drivesBitMask)

	for _, driveName := range drivesNames {
		drives = append(drives, NewDrive(driveName))
	}

	return
}
func ScanDrivesEx() (drives []Drive) {
	driveScanner := &DefaultWindowsDriveScanner{}
	drives = ScanDrives(driveScanner)
	return
}

func ScanRemovableDrives(drives []Drive, driveTypeDetector DriveTypeDetector) (removableDrives []Drive) {
	for _, drive := range drives {
		driveTypeDetector.GetDriveType(&drive)
		if drive.removable == true {
			removableDrives = append(removableDrives, drive)
		}
	}
	return
}
func ScanRemovableDrivesEx(drives []Drive) (removableDrives []Drive) {
	driverTypeDetector := DefaultWindowsDriveTypeDetector{}
	removableDrives = ScanRemovableDrives(drives, driverTypeDetector)
	return removableDrives
}


type DriveReader interface {
	ReadDir(drive Drive) ([]string, error)
}

type DefaultDriveReader struct {}

func (drv DefaultDriveReader) ReadDir(drive Drive) (files []string, err error) {
	fs, err := ioutil.ReadDir(drive.name)
	if err != nil {
		panic(err)
	}
	for _, file := range fs {
		files = append(files, file.Name())
	}
	return
}

// list file names found in the drive object
func (d Drive) ListFiles(driveReader DriveReader) (files []string){
	files, err := driveReader.ReadDir(d)
	if err != nil {
		panic(err)
	}
	return
}

// Exported version
func (drive Drive) ListFilesEx() (files[] string) {
	driveReader := DefaultDriveReader{}
	files = drive.ListFiles(driveReader)
	return
}