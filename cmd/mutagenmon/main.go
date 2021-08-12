package main

import (
	"go.andmed.org/mutagenmon"
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

func main() {
	var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	out, err := ioutil.TempFile("/tmp", "mutagenmon-log-*")
	if err == nil {
		log.SetOutput(out)
	}

	cmd := exec.Command("/usr/local/bin/mutagen", "list")


	var mm *mutagenmon.MutagenMon
	for {
		err := cmd.Run()
		if err != nil {
			log.Printf("[ERROR] %s\n", err)
		}

		mm, err = mutagenmon.New()
		if err == nil {
			break
		}
		log.Printf("[Info] waiting for initialization\n")
		time.Sleep(mutagenmon.IntervalSec * time.Second)
	}
	mm.Run()
}
