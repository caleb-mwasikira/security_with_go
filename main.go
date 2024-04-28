package main

import (
	"log"

	arch "example.com/backup/archives"
)

func main() {
	source := "/home/netrunner/College"
	dest := "/home/netrunner/College.zip"

	dest, err := arch.ZipArchive(source, "")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("done archiving source %v into destination %v", source, dest)

	archive := "~/Videos.zip"
	err = arch.Unzip(archive)
	if err != nil {
		log.Fatal(err)
	}
}
