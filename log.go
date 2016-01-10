package log

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"runtime"
	"sync/atomic"
)

var (
	logInfo  *File
	logError *File

	dot       = []byte(".")
	centerDot = []byte("·")

	// LogToStdout tells to log to stdout if true
	LogToStdout    bool
	verbosityLevel int32
)

// IncVerbosity increases verbosity level.
// the idea of verbose logging is to provide a way to turn detailed logging
// on a per-request basis. This is an approximate solution: since there is
// no per-gorutine context, we use a shared variable that is increased at request
// beginning and decreased at end. We might get additional logging from other
// gorutines. It's much simpler than an alternative, like passing a logger
// to every function that needs to log
func IncVerbosity() {
	atomic.AddInt32(&verbosityLevel, 1)
}

// DecVerbosity decreases verbosity level
func DecVerbosity() {
	atomic.AddInt32(&verbosityLevel, -1)
}

// IsVerbose returns true if we're doing verbose logging
func IsVerbose() bool {
	return atomic.LoadInt32(&verbosityLevel) > 0
}

/*
StartVerboseForURL will start verbose logging if the url has vl= arg in it.
The intended usage is:

if StartVerboseForURL(r.URL) {
  defer StopVerboseForURL()
}
*/
func StartVerboseForURL(u *url.URL) bool {
	// "vl" stands for "verbose logging" and any value other than empty string
	// truns it on
	if u.Query().Get("vl") != "" {
		IncVerbosity()
		return true
	}
	return false
}

// StopVerboseForURL is for name parity with StartVerboseForURL()
func StopVerboseForURL() {
	DecVerbosity()
}

func open(dir, suffix string, fileOut **File) error {
	lf, err := NewFile(dir, suffix)
	if err != nil {
		return err
	}
	*fileOut = lf
	return nil
}

// Open opens a standard log file
func Open(dir, suffix string) error {
	return open(dir, suffix, &logInfo)
}

// OpenError opens a log file for
func OpenError(dir, suffix string) error {
	return open(dir, suffix, &logError)
}

// Close closes all log files
func Close() {
	logInfo.Close()
	logInfo = nil
	logError.Close()
	logError = nil
}

func functionFromPc(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//      runtime/debug.*T·ptrmethod
	// and want
	//      *T.ptrmethod
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return string(name)
}

func p(info *File, err *File, s string) {
	if LogToStdout {
		fmt.Print(s)
	}
	if err != nil {
		err.Print(s)
		return
	}
	info.Print(s)
}

// Fatalf is like log.Fatalf() but also pre-pends name of the caller,
// so that we don't have to do that manually in every log statement
func Fatalf(format string, arg ...interface{}) {
	s := fmt.Sprintf(format, arg...)
	if pc, _, _, ok := runtime.Caller(1); ok {
		s = functionFromPc(pc) + ": " + s
	}
	p(logInfo, logError, s)
	fmt.Print(s)
	log.Fatal(s)
}

// Errorf logs an error to error log (if not available, to info log)
// Prepends name of the function that called it
func Errorf(format string, arg ...interface{}) {
	s := fmt.Sprintf(format, arg...)
	if pc, _, _, ok := runtime.Caller(1); ok {
		s = functionFromPc(pc) + ": " + s
	}
	p(logInfo, logError, s)
}

// Error logs error to error log (if not available, to info log)
func Error(err error) {
	s := err.Error() + "\n"
	if pc, _, _, ok := runtime.Caller(1); ok {
		s = functionFromPc(pc) + ": " + s
	}
	p(logInfo, logError, s)
}

// Infof logs non-error things
func Infof(format string, arg ...interface{}) {
	s := fmt.Sprintf(format, arg...)
	if pc, _, _, ok := runtime.Caller(1); ok {
		s = functionFromPc(pc) + ": " + s
	}
	p(logInfo, nil, s)
}

// Verbosef logs more detailed information if verbose logging
// is turned on
func Verbosef(format string, arg ...interface{}) {
	if !IsVerbose() {
		return
	}
	s := fmt.Sprintf(format, arg...)
	if pc, _, _, ok := runtime.Caller(1); ok {
		s = functionFromPc(pc) + ": " + s
	}
	p(logInfo, nil, s)
}
