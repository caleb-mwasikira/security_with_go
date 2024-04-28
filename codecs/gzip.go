package codecs

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*
creates a file's path and its parent dirs if they dont exist
*/
func ensureFilePath(_path string) error {
	err := os.MkdirAll(filepath.Dir(_path), 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(_path)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func gzipCompressFile(source string, dest string) error {
	// open source file for reading
	infile, err := os.OpenFile(source, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer infile.Close()

	err = ensureFilePath(dest)
	if err != nil {
		return err
	}

	log.Printf("compressing file %v into destination %v ...", source, dest)

	// create/open the file to write our compressed contents to
	outfile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil && err != os.ErrNotExist {
		return err
	}
	defer outfile.Close()

	// create a gzip writer on top of the file writer
	gzip_writer := gzip.NewWriter(outfile)
	defer gzip_writer.Close()

	// compress file
	buffer := make([]byte, 1024)

	for {
		// read a chunk
		n, err := infile.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// now we write onto the gzip writer.
		// whatever data we write onto the gzip writer
		// will be compressed and written to the underlying
		// file as well
		_, err = gzip_writer.Write(buffer)
		if err != nil {
			return err
		}
	}

	return nil
}

func getDestinationFileName(source string, flag int) (string, error) {
	// verify is source is a valid path
	stat, err := os.Stat(source)
	if err != nil {
		return "", err
	}

	if stat.IsDir() {
		return "", fmt.Errorf("codecs.getDestinationFileName: expected source path to be a file but found a folder instead")
	}

	var (
		dest    string = BACKUP_DIR
		subdirs string
		found   bool = false
	)

	fname := filepath.Base(source)
	parent := filepath.Dir(source)

	if flag == PRESERVE_DIR {
		subdirs, found = strings.CutPrefix(parent, HOME_DIR)
		if !found {
			return "", fmt.Errorf("codecs.getDestinationFileName: failed to preserve destination sub directories. error: %v", err)
		}
	}

	dest = filepath.Join(dest, subdirs, fmt.Sprintf("%v.gz", fname))
	return dest, nil
}

func GzipCompress(source string, flag int) error {
	// check if source is a valid path
	stat, err := os.Stat(source)
	if err != nil {
		return err
	}

	if stat.IsDir() {
		files_to_compress := []string{}

		// recursively walk dir tree listing all files
		// due for compression
		err = filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// skip root dir
			if path == source {
				return nil
			}

			if !d.IsDir() {
				files_to_compress = append(files_to_compress, path)
			}

			return err
		})

		for _, fpath := range files_to_compress {
			dest, err := getDestinationFileName(fpath, flag)
			if err != nil {
				return err
			}

			err = gzipCompressFile(fpath, dest)
			if err != nil {
				return fmt.Errorf("codecs.GzipCompress: failed to compress file. error: %v", err)
			}
		}

		return err
	} else {
		dest, err := getDestinationFileName(source, flag)
		if err != nil {
			return fmt.Errorf("codec.GzipCompress: failed to generate destination directory.\n%v", err)
		}

		err = gzipCompressFile(source, dest)
		return err
	}
}
