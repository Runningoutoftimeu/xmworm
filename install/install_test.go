package install

import (
	"testing"
)

// installing means that the worm
// copies itself from the infected drive
// into the startup folder of the victim pc.

type MockUserNameGetter struct {
	username string
}

func (user MockUserNameGetter) Username() (string){
	return user.username
}

func TestGetInstallDir(t* testing.T) {
	// it should test wether Test Function
	// if the install function can dynamically generated
	// the correct startup folder for the target computer
	// on which its run

	// testing the copying function was the responsiblity of the
	// files package.

	userObj := MockUserNameGetter{ username: "Rasta"}

	want := `C:\Users\Rasta\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\`
	got := getInstallDir(userObj)
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}