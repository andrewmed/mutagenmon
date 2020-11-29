package main

import (
	"go.andmed.org/mutagenmon"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	out, err := ioutil.TempFile("/tmp", "mutagenmon-log-*")
	if err == nil {
		log.SetOutput(out)
	}

	var mm *mutagenmon.MutagenMon
	for {
		mm, err = mutagenmon.New()
		if err == nil {
			break
		}
		time.Sleep(mutagenmon.IntervalSec * time.Second)
	}
	mm.Run()
}
