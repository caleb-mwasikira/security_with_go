package main

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ZipArchive(source string, dest string) error {
	// check if source is a valid directory
	stat, err := os.Stat(source)
	if err != nil {
		return err
	}

	// check if dest is a valid directory
	_, err = os.Stat(dest)
	if err != nil {
		// if dest directory does not exist create zip file
		// at the parent of source dir
		if os.IsNotExist(err) {
			parent := filepath.Dir(source)
			dest = filepath.Join(parent, fmt.Sprintf("%s.zip", stat.Name()))
		} else {
			return err
		}
	}

	// create zip file to write to
	log.Printf("creating zip file %s ...\n", dest)
	zip_file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zip_file.Close()

	// create a zip writer on top of the zip file
	zip_writer := zip.NewWriter(zip_file)

	// loop though all files in source dir adding them to the archive
	err = filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// skip adding to archive if we are still on the source dir
		if path == source {
			return nil
		}

		// get sub-directories after source dir
		// for example: we are archiving the dir /home/user/foo
		// that has the structure /home/user/foo/bar/baz.
		// we remove the prefix path (/home/user/foo) so that we can
		// remain with /bar/baz, that we can place in our archive as
		// foo.zip/bar/baz

		zip_path, found := strings.CutPrefix(path, source)
		if !found {
			return fmt.Errorf("directory %s has no prefix %s", path, source)
		}

		// if dirEntry is a directory, create new directory
		// inside the archive
		if d.IsDir() {
			zip_path = fmt.Sprintf("%s/", zip_path)
			log.Printf("creating directory %s in archive ...", zip_path)

			_, err := zip_writer.Create(fmt.Sprintf("%s/", zip_path))
			if err != nil {
				return err
			}

			return nil
		}

		// copy file into the zip archive
		log.Printf("adding file %s into archive ...", zip_path)

		out_file, err := zip_writer.Create(zip_path)
		if err != nil {
			return err
		}

		in_file, err := os.Open(path)
		if err != nil {
			return err
		}

		// close in_file
		defer func() {
			if err := in_file.Close(); err != nil {
				log.Panic(err)
			}
		}()

		err = CopyFile(in_file, out_file)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// close our zip writer
	err = zip_writer.Close()
	if err != nil {
		return err
	}

	return nil
}
