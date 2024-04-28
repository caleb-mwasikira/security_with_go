package main

import (
	"log"

	"example.com/backup/codecs"
)

func main() {
	source := "/home/netrunner/College"
	dest := "/home/netrunner/College.zip"

	// dest, err := arch.ZipArchive(source, "")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("done archiving source %v into destination %v", source, dest)

	// archive := "~/Videos.zip"
	// err = arch.Unzip(archive)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := codecs.GzipCompress(dest, codecs.PRESERVE_DIR)
	if err != nil {
		log.Fatalf("GzipCompress failed to compress file %v - error: %v", source, err)
	}
}
