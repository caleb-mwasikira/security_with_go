package main

import (
	"log"
)

func main() {
	src := "/home/netrunner/College"

	err := ZipArchive(src, "")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("done archiving files")
}
