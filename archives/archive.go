package archives

import (
	"archive/zip"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ZipArchive(source string, dest string) (string, error) {
	// check if source is a valid directory
	stat, err := os.Stat(source)
	if err != nil {
		return "", err
	}

	if !stat.IsDir() {
		return "", fmt.Errorf("failed to archive path %s as it is not a directory", source)
	}

	// check if dest is a valid directory
	_, err = os.Stat(dest)
	if err != nil {
		// if dest directory does not exist create zip file
		// at the parent of source dir
		if os.IsNotExist(err) {
			parent := filepath.Dir(source)
			dest = filepath.Join(parent, fmt.Sprintf("%v.zip", stat.Name()))
		} else {
			return "", err
		}
	}

	// create zip file to write to
	zip_file, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer zip_file.Close()

	// create a zip writer on top of the zip file
	zip_writer := zip.NewWriter(zip_file)
	defer zip_writer.Close()

	// loop though all files in source dir adding them to the archive
	err = filepath.WalkDir(
		source,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// skip adding to archive if we are still on the source dir
			if path == source {
				return nil
			}

			zip_path, found := strings.CutPrefix(path, source)
			if !found {
				return fmt.Errorf("directory %s has no prefix %s", path, source)
			}

			// if dirEntry is a directory, create new directory
			// inside the archive
			if d.IsDir() {
				zip_path = fmt.Sprintf("%s/", zip_path)
				log.Printf("creating directory %s in archive ...", zip_path)

				_, err := zip_writer.Create(zip_path)
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
			defer in_file.Close()

			_, err = io.Copy(out_file, in_file)
			if err != nil {
				return err
			}

			return nil
		})

	if err != nil {
		return "", err
	}

	return dest, nil
}

func genNewFilename() (string, error) {
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		return "", fmt.Errorf("failed to generate new filename")
	}

	new_filename := string(md5.New().Sum(data))
	return new_filename, nil
}

func deleteFileIfExists(filepath string) bool {
	stat, _ := os.Stat(filepath)
	if stat != nil {
		err := os.Remove(filepath)
		return err == nil
	}

	return true
}

func Unzip(source string) error {
	// open the archive with with a zip reader
	zip_reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer zip_reader.Close()

	// define the dest root directory
	parent := filepath.Dir(source)
	dirname := filepath.Base(source)

	if filepath.Ext(dirname) != ".zip" {
		return fmt.Errorf("expected archive in ZIP format")
	}

	// strip .zip suffix from dirname to get new dirname
	new_dirname, found := strings.CutSuffix(dirname, ".zip")
	if !found {
		return fmt.Errorf("extension .zip not found in archive filename")
	}

	if len(new_dirname) == 0 {
		new_dirname, err = genNewFilename()
		if err != nil {
			return fmt.Errorf("failed to generate new dirname for archive")
		}
	}
	dest := filepath.Join(parent, new_dirname)

	// delete destination dir if already exists
	success := deleteFileIfExists(dest)
	if !success {
		return fmt.Errorf("cannot extract archive on existing directory %v", dest)
	}

	// iterate through each of the files found in the archive
	for _, file := range zip_reader.Reader.File {
		// open the file like a normal file
		zipped_file, err := file.Open()
		if err != nil {
			return err
		}
		defer zipped_file.Close()

		// specify what the extracted filename should be
		extracted_filepath := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			log.Printf("creating directory %v ...\n", extracted_filepath)
			err := os.MkdirAll(extracted_filepath, 0700)
			if err != nil {
				return fmt.Errorf("failed to create directory %v: %v", extracted_filepath, err)
			}
		} else {
			log.Printf("extracting file %v ...\n", file.Name)

			out_file, err := os.OpenFile(
				extracted_filepath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				0700,
			)
			if err != nil {
				return err
			}
			defer out_file.Close()

			// copy contents of archive file into out_file
			_, err = io.Copy(out_file, zipped_file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
