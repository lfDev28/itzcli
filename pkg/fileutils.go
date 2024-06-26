package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	logger "github.com/sirupsen/logrus"
)

var re = regexp.MustCompile(`^file:\/{0,2}`)

// GetITZHomeDir returns the home directory or the ITZ command
func GetITZHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", os.ErrNotExist
	}
	return filepath.Join(home, ".itz"), nil
}

func MustITZHomeDir() string {
	home, err := GetITZHomeDir()
	if err != nil {
		logger.Fatal(err)
	}
	return home
}


// Use the templates folder for local development
func GetLocalDir() string {
	// Retunr the pwd
	dir, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}
	return dir
}

// AppendToFilename appends the suffix to the name of the file. The file is expected
// to be a URL
func AppendToFilename(fn string, suffix string) (string, error) {
	u, err := url.Parse(fn)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(u.Path)
	base := filepath.Base(u.Path)
	baseParts := strings.Split(base, ".")
	ext := filepath.Ext(u.Path)
	u.Path = fmt.Sprintf("%s/%s%s%s", dir, strings.Join(baseParts[:len(baseParts)-1], "."), suffix, ext)
	return u.String(), nil
}

// ReadFile reads the given file into the byte array
func ReadFile(path string) ([]byte, error) {
	// Yes, for now this is just a wrapper around os.ReadFile
	p := re.ReplaceAllString(path, ``)
	return os.ReadFile(p)
}

func WriteFile(path string, data []byte) error {
	// Open a new file for writing only
	file, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0600,
	)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func AppendToFile(source string, to string) error {
	// Open a new file for writing only
	file, err := os.OpenFile(
		source,
		os.O_RDONLY,
		0600,
	)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	// open the other file...
	file2, err := os.OpenFile(
		to,
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0600,
	)

	if err != nil {
		return err
	}
	defer func() {
		if err := file2.Close(); err != nil {
			panic(err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			return err
		}
		_, err = file2.Write(buf[:n])
		if err != nil {
			break
		}
	}

	return err
}

func StringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Unzip extracts the zip archive to the specified directory
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		logger.Errorf("could not create directory %s: %v", dest, err)
		return err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				logger.Errorf("could not create directory %s: %v", path, err)
				return err
			}
		} else {
			err = os.MkdirAll(filepath.Dir(path), f.Mode())
			if err != nil {
				logger.Errorf("could not create directory %s: %v", path, err)
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
