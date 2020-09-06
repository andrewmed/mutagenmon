package main

import (
	"go.andmed.org/mutagenmon"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mm, err := mutagenmon.New()
	if err != nil {
		log.Fatal(err)
	}

	mm.Run()
}
