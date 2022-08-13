// Package log implements a level logger using standard Go log library.
package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Standard flags for no verbose logging.
const stdFlags = log.LstdFlags | log.Lmicroseconds

// Available logging levels.
const (
	LevelInfo    Severity = iota // Lower
	LevelWarning                 // Medium
	LevelError                   // High
)

// Severity represents logging level.
type Severity int

// Logger is the logger structure.
type Logger struct {
	out       *log.Logger
	level     Severity
	calldepth int
}

// New instantiates a new Logger.
// level is the minimum logging level message to be printed.
// By default all logs are printed on standard output.
func New(level Severity) *Logger {
	return &Logger{
		out:       log.New(os.Stdout, "", stdFlags),
		level:     level,
		calldepth: 2,
	}
}

var std = newStd()

// Info logs an Info level message on the standard output.
// Arguments are handled in the manner of fmt.Print.
// Log message is emitted only if the current logging level is equal or less than LevelInfo.
func Info(v ...interface{}) {
	std.Info(v...)
}

// Infof logs an Info level message on the standard output.
// Arguments are handled in the manner of fmt.Printf.
// Log message is emitted only if the current logging level is equal or less than LevelInfo.
func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

// Infoln logs an Info level message on the standard output.
// Arguments are handled in the manner of fmt.Println.
// Log message is emitted only if the current logging level is equal or less than LevelInfo.
func Infoln(v ...interface{}) {
	std.Infoln(v...)
}

// Warning logs a Warning level message on the standard output.
// Arguments are handled in the manner of fmt.Print.
// Log message is emitted only if the current logging level is equal or less than LevelWarning.
func Warning(v ...interface{}) {
	std.Warning(v...)
}

// Warningf logs a Warning level message on the standard output.
// Arguments are handled in the manner of fmt.Printf.
// Log message is emitted only if the current logging level is equal or less than LevelWarning.
func Warningf(format string, v ...interface{}) {
	std.Warningf(format, v...)
}

// Warningln logs a Warning level message on the standard output.
// Arguments are handled in the manner of fmt.Println.
// Log message is emitted only if the current logging level is equal or less than LevelWarning.
func Warningln(v ...interface{}) {
	std.Warningln(v...)
}

// Error logs an Error level message on the standard error.
// Arguments are handled in the manner of fmt.Print.
func Error(v ...interface{}) {
	std.Error(v...)
}

// Errorf logs an Error level message on the standard error.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

// Errorln logs an Error level message on the standard error.
// Arguments are handled in the manner of fmt.Println.
func Errorln(v ...interface{}) {
	std.Errorln(v...)
}

// Fatal logs an Error level message on the standard error and calls os.Exit(1).
// Arguments are handled in the manner of fmt.Print.
func Fatal(v ...interface{}) {
	std.Fatal(v...)
}

// Fatalf logs an Error level message on the standard error and calls os.Exit(1).
// Arguments are handled in the manner of fmt.Printf.
func Fatalf(format string, v ...interface{}) {
	std.Fatalf(format, v...)
}

// Fatalln logs an Error level message on the standard error and calls os.Exit(1).
// Arguments are handled in the manner of fmt.Println.
func Fatalln(v ...interface{}) {
	std.Fatalln(v...)
}

// Verbose selects between short or verbose prefix (currently adds file and line number).
func Verbose(v bool) {
	std.Verbose(v)
}

// SetLevel selects the minimum logging level to print.
func SetLevel(level Severity) {
	std.SetLevel(level)
}

// Level returns the log level currently set.
func Level() Severity {
	return std.Level()
}

// SetWriter sets the logger's output stream for messages.
func SetWriter(w io.Writer) {
	std.SetWriter(w)
}

// Writer returns the output stream for the logger.
func Writer() io.Writer {
	return std.Writer()
}

var prefix = [...]string{LevelInfo: "INFO> ", LevelWarning: "WARN> ", LevelError: "ERROR> "}

// Info logs an Info level message on the standard output.
// Arguments are handled in the manner of fmt.Print.
// Log message is emitted only if the current logging level is equal or less than LevelInfo.
func (l *Logger) Info(v ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	l.out.Output(l.calldepth, prefix[LevelInfo]+fmt.Sprint(v...)) // #nosec
}

// Infof logs an Info level message on the standard output.
// Arguments are handled in the manner of fmt.Printf.
// Log message is emitted only if the current logging level is equal or less than LevelInfo.
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	l.out.Output(l.calldepth, fmt.Sprintf(prefix[LevelInfo]+format, v...)) // #nosec
}

// Infoln logs an Info level message on the standard output.
// Arguments are handled in the manner of fmt.Println.
// Log message is emitted only if the current logging level is equal or less than LevelInfo.
func (l *Logger) Infoln(v ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	l.out.Output(l.calldepth, prefix[LevelInfo]+fmt.Sprintln(v...)) // #nosec
}

// Warning logs a Warning level message on the standard output.
// Arguments are handled in the manner of fmt.Print.
// Log message is emitted only if the current logging level is equal or less than LevelWarning.
func (l *Logger) Warning(v ...interface{}) {
	if l.level > LevelWarning {
		return
	}
	l.out.Output(l.calldepth, prefix[LevelWarning]+fmt.Sprint(v...)) // #nosec
}

// Warningf logs a Warning level message on the standard output.
// Arguments are handled in the manner of fmt.Printf.
// Log message is emitted only if the current logging level is equal or less than LevelWarning.
func (l *Logger) Warningf(format string, v ...interface{}) {
	if l.level > LevelWarning {
		return
	}
	l.out.Output(l.calldepth, fmt.Sprintf(prefix[LevelWarning]+format, v...)) // #nosec
}

// Warningln logs a Warning level message on the standard output.
// Arguments are handled in the manner of fmt.Println.
// Log message is emitted only if the current logging level is equal or less than LevelWarning.
func (l *Logger) Warningln(v ...interface{}) {
	if l.level > LevelWarning {
		return
	}
	l.out.Output(l.calldepth, prefix[LevelWarning]+fmt.Sprintln(v...)) // #nosec
}

// Error logs an Error level message on the standard error.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Error(v ...interface{}) {
	l.out.Output(l.calldepth, prefix[LevelError]+fmt.Sprint(v...)) // #nosec
}

// Errorf logs an Error level message on the standard error.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.out.Output(l.calldepth, fmt.Sprintf(prefix[LevelError]+format, v...)) // #nosec
}

// Errorln logs an Error level message on the standard error.
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Errorln(v ...interface{}) {
	l.out.Output(l.calldepth, prefix[LevelError]+fmt.Sprintln(v...)) // #nosec
}

// Fatal logs an Error level message on the standard error and calls os.Exit(1).
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Fatal(v ...interface{}) {
	l.out.Output(l.calldepth, prefix[LevelError]+fmt.Sprint(v...)) // #nosec
	os.Exit(1)
}

// Fatalf logs an Error level message on the standard error and calls os.Exit(1).
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.out.Output(l.calldepth, fmt.Sprintf(prefix[LevelError]+format, v...)) // #nosec
	os.Exit(1)
}

// Fatalln logs an Error level message on the standard error and calls os.Exit(1).
// Arguments are handled in the manner of fmt.Println.
func (l *Logger) Fatalln(v ...interface{}) {
	l.out.Output(l.calldepth, prefix[LevelError]+fmt.Sprintln(v...)) // #nosec
	os.Exit(1)
}

// Verbose selects between short or verbose prefix (currently adds file and line number).
func (l *Logger) Verbose(v bool) {
	flags := stdFlags
	if v {
		flags |= log.Lshortfile
	}

	l.out.SetFlags(flags)
}

// SetLevel selects the minimum logging level to print.
func (l *Logger) SetLevel(level Severity) {
	l.level = level
}

// Level returns the log level currently set.
func (l *Logger) Level() Severity {
	return l.level
}

// SetWriter sets the logger's output stream for messages.
func (l *Logger) SetWriter(w io.Writer) {
	l.out.SetOutput(w)
}

// Writer returns the output stream for the logger.
func (l *Logger) Writer() io.Writer {
	return l.out.Writer()
}

// newStd is used to initializes the default logger.
func newStd() *Logger {
	v := New(LevelInfo)
	v.calldepth = 3
	return v
}
