package main

import (
	"log"
)

func main() {
	src := "/home/netrunner/College"
	archive := "/home/netrunner/College.zip"

	err := ZipArchive(src, "")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("done archiving files")

	err = Unzip(archive)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("done extracting zip file")
}
