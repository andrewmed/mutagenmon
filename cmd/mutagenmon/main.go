package main

import (
	"go.andmed.org/mutagenmon"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var mm *mutagenmon.MutagenMon
	var err error
	for {
		mm, err = mutagenmon.New()
		if err == nil {
			break
		}
		time.Sleep(mutagenmon.IntervalSec * time.Second)
	}
	mm.Run()
}
