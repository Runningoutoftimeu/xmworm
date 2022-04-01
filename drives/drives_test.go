package drives

import (
	"testing"
	"reflect"
)


type MockWindowsDriveScanner struct {
	drivesBitMask uint32
}
func (sc MockWindowsDriveScanner) GetLogicalDrives() (drivesBitMask uint32, err error) {
	return sc.drivesBitMask, nil
}


func TestScanDrives(t *testing.T) {

	t.Run("it should return a list of Drives on the system", func(t *testing.T) {
		driveScanner := &DefaultWindowsDriveScanner{}
		want := reflect.TypeOf([]Drive{})
		got := reflect.TypeOf(ScanDrives(driveScanner))

		if want != got {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	t.Run("it should return only drives E:\\, F:\\ and Z:\\", func (t *testing.T) {
		
		driveScanner := &MockWindowsDriveScanner{
			drivesBitMask: 0b10000000000000000000110000,
		}

		want := []Drive {
			Drive{name: `E:\`, removable: false},
			Drive{name: `F:\`, removable: false},
			Drive{name: `Z:\`, removable: false},
		}

		got := ScanDrives(driveScanner)

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}

type MockDriverTypeDetector struct {
	config map[string]bool
}

func (d MockDriverTypeDetector) GetDriveType(drive *Drive) {
	if d.config[drive.name] == true {
		drive.removable = true
	}else {
		drive.removable = false
	}
}

func TestScanRemovableDrives(t *testing.T) {
	driverTypeDetector := MockDriverTypeDetector{
		config: map[string]bool{
			`A:\`: false, `B:\`: true, `C:\`: false,
			`D:\`: false, `E:\`: true, `F:\`: false,
		    `G:\`: true, `H:\`: false, `I:\`: false,
			`J:\`: true, `K:\`: false,
		},
	}
	input := []Drive {
		Drive{name: `A:\`, removable: false},
		Drive{name: `B:\`, removable: false},
		Drive{name: `C:\`, removable: false},
		Drive{name: `D:\`, removable: false},
		Drive{name: `E:\`, removable: false},
		Drive{name: `F:\`, removable: false},
		Drive{name: `G:\`, removable: false},
		Drive{name: `H:\`, removable: false},
		Drive{name: `I:\`, removable: false},
		Drive{name: `J:\`, removable: false},
		Drive{name: `K:\`, removable: false},
	}

	want := []Drive {
		Drive{name: `B:\`, removable: true},
		Drive{name: `E:\`, removable: true},
		Drive{name: `G:\`, removable: true},
		Drive{name: `J:\`, removable: true},
	}

	got := ScanRemovableDrives(input, driverTypeDetector)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}

}


/* Works, but requires you that usb be  plugged and labelled by system
as G:\

func TestDefaultWindowsDriveTypeDetector(t *testing.T) {

	driverTypeDetector := DefaultWindowsDriveTypeDetector{}
	want := &Drive{name: `G:\`, removable: true}
	got := &Drive{name: `G:\`}
	driverTypeDetector.GetDriveType(got)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}

}

*/

/*

TestScanWritableRemovableDrives
*/


type MockDriveReader struct {}

func (drv MockDriveReader) ReadDir(drive Drive) (files []string, err error){
	err = nil
	files = drive.files
	return 
}

func TestListFiles(t *testing.T) {
	t.Run("it should return a list of filenames inside the drive", func(t *testing.T) {
		want := []string{"1.exe", "2.exe", "foobar.bak", "3.txt"}
		drive := Drive{name: "G:/", files: want}
		driveReader := MockDriveReader{}
		got := drive.ListFiles(driveReader)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}