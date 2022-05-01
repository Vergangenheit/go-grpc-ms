package files

import (
	"io"
	"path/filepath"
)

// Local is an implementation of the Storage interface which works with the
// local disk on the current machine
type Local struct {
	maxFileSize int // maximum number of bytes for files
	basePath    string
}

// NewLocal creates a new Local filesytem with the given base path
// basePath is the base directory to save files to
// maxSize is the max number of bytes that a file can be
func NewLocal(basePath string, maxSize int) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &Local{basePath: p}, nil
}

// Save the contents of the Writer to the given path
// path is a relative path, basePath will be appended
func (l *Local) Save(path string, contents io.Reader) error {
	// get the full path for the file
	fp := l.fullPath(path)

	// gwt the directory and make sure it exists
	d := filepath.Dir(fp)

}

// returns the absolute path
func (l *Local) fullPath(path string) string {
	// append the given path to the base path
	return filepath.Join(l.basePath, path)
}
