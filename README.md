### log - small and opinionated logging library for Go

I've been using pretty much the same logging code in my projects so I extracted
it into a package.

This code is small and opinionated - it does logging the way I like to do
logging (opinionated) and has zero extra features (small).

How do I log in my projects:
* the context is logging for server software
* logs a are rotated daily
* log files are stored in their own directory
* log file names start with `YYYY-MM-DD` and you can append a suffix (e.g. `.txt`
  or something more unique if you have more than one thing creating log files in
  that directory)
* for easier debugging, each log line starts with the name of the function calling
  the log routine

How to use the logging library:
* `import github.com/kjk/log`
* call `log.Open(dir, suffix string)` to start logging to a file
* optionally call `log.OpenError(dir, suffix string)`. If not called, `log.Errorf()`
  will log to the same file as non-error logs. Don't forget to use a different
  suffix (e.g. `-err.txt`)
* set `log.LogToStdout` if you want logging also to be printed to stdout. Useful
  in dev mode, when working locally
* call `log.Infof(format string, arg ...interface{})` to log a formatted string.
  Unlike standard `log` package, we don't automatically add `\n` at the end,
  so it's your job to add it if desired
* call `log.Errorf()` to log an error formatted string. It will write to a
  separate error log file opened with `log.OpenError()`. This is useful if
  you want to look for errors (e.g. with `tail -f`) and not be drowned with
  potentially much noisier output from `log.Infof()`
* call `log.Error(err error)` is a convenience function, equivalent to `log.Errorf("%s\n", err)`
* call `log.Close()` when you're done (but it's fine if you don't; os will close
  the file anyway)
* call `log.Verbosef()` for that extra information that is only printed if
  verbosity level is > 0. Call `log.IncVerbosity()`, `log.DecVerbosity()` to
  manage verbosity level.

The idea behind verbosity logging is that when you investigate an issue, it's
useful to have more info, but the same info is noise in regular use.

Instead of using `log.Infof()`, which always log, you can add `log.Verbosef()`
which only print if verbosity is > 0.

In local development, you can increase verbosity for the whole duration of the
program (e.g. add `-verbose` flag to your program and call `log.IncVerbosity()`
at startup if flag is set).

In production it's useful to be able to turn verbose logging on and off. In web
servers you could turn verbosity on and off by adding simple http handlers e.g.
`/api/debug/incverbosity`, `/api/debug/decverbosity`.

In web servers it's often useful to turn verbosity for the duration of a single
http request. You can do it by adding this to your http handlers:
```go
if StartVerboseForURL(r.URL) {
  defer StopVerboseForURL()
}
```

This will turn verbose logging if url has `vl=1` argument in it. You can
easily implement a different scheme.
