package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DailyRotateFile describes a file that gets rotated daily
type DailyRotateFile struct {
	sync.Mutex
	pathFormat string

	// info about currently opened file
	day  int
	path string
	file *os.File
}

func (f *DailyRotateFile) close() error {
	var err error
	if f.file != nil {
		err = f.file.Close()
		f.file = nil
	}
	return err
}

func (f *DailyRotateFile) open() error {
	t := time.Now()
	f.path = t.Format(f.pathFormat)
	f.day = t.YearDay()

	// we can't assume that the dir for the file already exists
	dir := filepath.Dir(f.path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	flag := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	f.file, err = os.OpenFile(f.path, flag, 0644)
	return err
}

// rotate on new day
func (f *DailyRotateFile) reopenIfNeeded() error {
	t := time.Now()
	if t.YearDay() == f.day {
		return nil
	}
	err := f.close()
	if err != nil {
		return err
	}
	return f.open()
}

// NewDailyRotateFile opens a new log file (creates if doesn't exist, will append if exists)
func NewDailyRotateFile(pathFormat string) (*DailyRotateFile, error) {
	res := &DailyRotateFile{
		pathFormat: pathFormat,
	}
	if err := res.open(); err != nil {
		return nil, err
	}
	return res, nil
}

// Close closes the file
func (f *DailyRotateFile) Close() error {
	var err error
	if f != nil {
		f.Lock()
		err = f.close()
		f.Unlock()
	}
	return err
}

// Write writes data to a file
func (f *DailyRotateFile) Write(d []byte) (int, error) {
	if f == nil {
		return 0, errors.New("File not opened")
	}
	f.Lock()
	f.Unlock()
	err := f.reopenIfNeeded()
	if err != nil {
		return 0, err
	}
	return f.file.Write(d)
}

// WriteString writes a string to a file
func (f *DailyRotateFile) WriteString(s string) (int, error) {
	return f.Write([]byte(s))
}

// Printf formats and writes to the file
func (f *DailyRotateFile) Printf(format string, arg ...interface{}) {
	f.WriteString(fmt.Sprintf(format, arg...))
}
