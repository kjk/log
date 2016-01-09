package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// File describes a log file
type File struct {
	sync.Mutex
	dir     string
	suffix  string
	currDay int
	path    string
	file    *os.File
}

func (l *File) close() {
	if l.file != nil {
		l.file.Close()
		l.file = nil
	}
}

func (l *File) open() (err error) {
	flag := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	l.file, err = os.OpenFile(l.path, flag, 0644)
	if err != nil {
		return err
	}
	_, err = l.file.Stat()
	if err != nil {
		l.close()
	}
	return err
}

// NewFile opens a new log file (creates if doesn't exist, will append if exists)
func NewFile(dir, suffix string) (*File, error) {
	res := &File{
		dir:    dir,
		suffix: suffix,
		path:   filepath.Join(dir, suffix), // TODO: add time
	}
	if err := res.open(); err != nil {
		return nil, err
	}
	return res, nil
}

// Close closes a log file
func (l *File) Close() {
	l.Lock()
	defer l.Unlock()
	if l == nil {
		return
	}
	l.close()
}

// Print writes to log file
func (l *File) Print(s string) {
	if l == nil {
		return
	}
	l.Lock()
	defer l.Unlock()
	l.file.Write([]byte(s))
}

// Printf formats and writes to log file
func (l *File) Printf(format string, arg ...interface{}) {
	l.Print(fmt.Sprintf(format, arg...))
}
