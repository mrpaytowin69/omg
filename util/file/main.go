package file

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// IsProtected returns true if the file is too critical to alter or remove
func IsProtected(fpath string) bool {
	if fpath == "" {
		return true
	}
	absPath, err := filepath.Abs(fpath)
	if err != nil {
		return false
	}
	switch absPath {
	case "/", "/bin", "/boot", "/dev", "/dev/pts", "/dev/shm", "/home", "/opt", "/proc", "/sys", "/tmp", "/usr", "/var":
		return true
	default:
		return false
	}
}

// Exists returns true if the file path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// ExistsNotDir returns true if the file path exists and is not a directory.
func ExistsNotDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ExistsAndDir returns true if the file path exists and is a directory.
func ExistsAndDir(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// ExistsAndRegular returns true if the file path exists and is a regular file.
func ExistsAndRegular(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.Mode().IsRegular()
}

//
// Copy copies the file content from src file path to dst file path.
// If dst does not exist, it is created.
//
func Copy(src string, dst string) (err error) {
	var (
		r *os.File
		w *os.File
	)
	if r, err = os.Open(src); err != nil {
		return err
	}
	defer r.Close()
	if w, err = os.Create(dst); err != nil {
		return err
	}
	defer w.Close()
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	r.Close()
	w.Close()
	return nil
}

//
// ReadAll reads and return all content of the file at src
//
func ReadAll(src string) ([]byte, error) {
	var (
		r   *os.File
		err error
	)
	if r, err = os.Open(src); err != nil {
		return []byte{}, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

//
// ModTime returns the file modification time or a zero time.
//
func ModTime(p string) (mtime time.Time) {
	fi, err := os.Stat(p)
	if err != nil {
		return
	}
	mtime = fi.ModTime()
	return
}
