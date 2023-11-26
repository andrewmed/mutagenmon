package main

import (
	"go.andmed.org/mutagenmon"
	"log"
	"time"
)

func main() {
	var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var mm *mutagenmon.MutagenMon
	for {
		mm, err = mutagenmon.New()
		if err == nil {
			break
		}
		log.Printf("[Info] waiting for initialization\n")
		time.Sleep(3 * time.Second)
	}
	mm.Run()
}
