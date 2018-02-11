// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Level represents the log level.
type Level int8

const (
	// DBG is the Level for debug messages.
	DBG Level = iota
	// INF is the Level for informational messages.
	INF
	// WRN is the Level for warning messages.
	WRN
	// ERR is the Level for error messages.
	ERR
	// NUL is the Level corresponding to no log output.
	NUL
)

type Flag int8

const (
	FlagSourceInfo = 1 << iota
	FlagFunctionInfo
)

const FunctionWidth int = 32

// String returns a string representation of the log level for use in traces.
func (l Level) String() string {
	switch l {
	case DBG:
		return "[D]"
	case INF:
		return "[I]"
	case WRN:
		return "[W]"
	case ERR:
		return "[E]"
	}
	return ""
}

// logln is the prototype of log functions writing a line to a stream.
type logln func(writer io.Writer, args ...interface{}) (int, error)

// logf is the prototype of log functions writing a formatted output to a stream.
type logf func(writer io.Writer, format string, args ...interface{}) (int, error)

var (
	logLevel          Level
	logLevelLock      sync.RWMutex
	logStream         io.Writer
	logStreamLock     sync.RWMutex
	logTimeFormat     string
	logTimeFormatLock sync.RWMutex
	logColorise       bool
	logColoriseLock   sync.RWMutex
	logFlags          int8
	logFlagsLock      sync.RWMutex
	logDebugf         logf
	logInfof          logf
	logWarnf          logf
	logErrorf         logf
	logDebugln        logln
	logInfoln         logln
	logWarnln         logln
	logErrorln        logln
)

func init() {
	SetLevel(DBG)
	SetStream(os.Stderr)
	SetTimeFormat("2006-01-02@15:04:05.000")
	if runtime.GOOS == "windows" {
		SetColorise(false)
	} else {
		SetColorise(true)
	}
	SetFlags(FlagSourceInfo | FlagFunctionInfo)
}

// SetLevel sets the log level for the application.
func SetLevel(level Level) {
	logLevelLock.Lock()
	defer logLevelLock.Unlock()
	logLevel = level
}

// GetLevel retur s the current log level.
func GetLevel() Level {
	logLevelLock.RLock()
	defer logLevelLock.RUnlock()
	return logLevel
}

// SetStream sets the stream to write messages to.
func SetStream(stream io.Writer) {
	logStreamLock.Lock()
	defer logStreamLock.Unlock()
	logStream = stream
}

// GetStream returns the current log stream.
func GetStream() io.Writer {
	logStreamLock.RLock()
	defer logStreamLock.RUnlock()
	return logStream
}

// SetTimeFormat sets the format for log messages time.
func SetTimeFormat(format string) {
	logTimeFormatLock.Lock()
	defer logTimeFormatLock.Unlock()
	logTimeFormat = format
}

// GetTimeFormat returns the current format of log messages time.
func GetTimeFormat() string {
	logTimeFormatLock.RLock()
	defer logTimeFormatLock.RUnlock()
	return logTimeFormat
}

// SetColorise enables or disables the colouring of the log messages
// according to their severity. By default this is disabled on
// Windows and enabled on *NIX systems; this function is the way
// to toggle it.
func SetColorise(enabled bool) {
	logColoriseLock.Lock()
	defer logColoriseLock.Unlock()
	if enabled {
		logDebugf = color.New(color.FgWhite).Fprintf
		logInfof = color.New(color.FgGreen).Fprintf
		logWarnf = color.New(color.FgYellow).Fprintf
		logErrorf = color.New(color.FgRed).Fprintf
		logDebugln = color.New(color.FgWhite).Fprintln
		logInfoln = color.New(color.FgGreen).Fprintln
		logWarnln = color.New(color.FgYellow).Fprintln
		logErrorln = color.New(color.FgRed).Fprintln
	} else if !enabled {
		logDebugf = fmt.Fprintf
		logInfof = fmt.Fprintf
		logWarnf = fmt.Fprintf
		logErrorf = fmt.Fprintf
		logDebugln = fmt.Fprintln
		logInfoln = fmt.Fprintln
		logWarnln = fmt.Fprintln
		logErrorln = fmt.Fprintln
	}
	logColorise = enabled
}

// SetFlags sets the flags governing how much information is gathered at runtime.
func SetFlags(flags int8) {
	logFlagsLock.Lock()
	defer logFlagsLock.Unlock()
	logFlags = flags
}

// GetFlags returns the flags governing how much information is gathered at runtime.
func GetFlags() int8 {
	logFlagsLock.RLock()
	defer logFlagsLock.RUnlock()
	return logFlags
}

// IsDebug returns whether the debug (DBG) log elevel is enabled.
func IsDebug() bool {
	return GetLevel() <= DBG
}

// IsInfo returns whether the informational (INF) log elevel is enabled.
func IsInfo() bool {
	return GetLevel() <= INF
}

// IsWarning returns whether the warning (WRN) log elevel is enabled.
func IsWarning() bool {
	return GetLevel() <= WRN
}

// IsError returns whether the error (ERR) log elevel is enabled.
func IsError() bool {
	return GetLevel() <= ERR
}

// IsDisabled returns whether the log is disabled.
func IsDisabled() bool {
	return GetLevel() <= NUL
}

// Debugln writes a debug message to the current output stream,
// appending a new line.
func Debugln(args ...interface{}) (int, error) {
	if IsDebug() {
		_, args = prepareFormatAndArgs(DBG, FunctionWidth, "", args)
		return logDebugln(GetStream(), args...)
	}
	return 0, nil
}

// Infoln writes an informational message to the current output stream,
// appending a new line.
func Infoln(args ...interface{}) (int, error) {
	if IsInfo() {
		_, args = prepareFormatAndArgs(INF, FunctionWidth, "", args)
		return logInfoln(GetStream(), args...)
	}
	return 0, nil
}

// Warnln writes a warning message to the current output stream,
// appending a new line.
func Warnln(args ...interface{}) (int, error) {
	if IsWarning() {
		_, args = prepareFormatAndArgs(WRN, FunctionWidth, "", args)
		return logWarnln(GetStream(), args...)
	}
	return 0, nil
}

// Errorln writes an error message to the current output stream,
// appending a new line.
func Errorln(args ...interface{}) (int, error) {
	if IsError() {
		_, args = prepareFormatAndArgs(ERR, FunctionWidth, "", args)
		return logErrorln(GetStream(), args...)
	}
	return 0, nil
}

// Debugf writes a debug message to the current output stream,
// appending a new line.
func Debugf(format string, args ...interface{}) (int, error) {
	if IsDebug() {
		format, args = prepareFormatAndArgs(DBG, FunctionWidth, format, args...)
		if !strings.HasSuffix(format, "\n") && !strings.HasSuffix(format, "\r") {
			format = format + "\n"
		}
		return logDebugf(GetStream(), format, args...)
	}
	return 0, nil
}

// Infof writes an informational message to the current output stream,
// appending a new line.
func Infof(format string, args ...interface{}) (int, error) {
	if IsInfo() {
		format, args = prepareFormatAndArgs(INF, FunctionWidth, format, args...)
		if !strings.HasSuffix(format, "\n") && !strings.HasSuffix(format, "\r") {
			format = format + "\n"
		}
		return logInfof(GetStream(), format, args...)
	}
	return 0, nil
}

// Warnf writes a warning message to the current output stream,
// appending a new line.
func Warnf(format string, args ...interface{}) (int, error) {
	if IsWarning() {
		format, args = prepareFormatAndArgs(WRN, FunctionWidth, format, args...)
		if !strings.HasSuffix(format, "\n") && !strings.HasSuffix(format, "\r") {
			format = format + "\n"
		}
		return logWarnf(GetStream(), format, args...)
	}
	return 0, nil
}

// Errorf writes an error message to the current output stream,
// appending a new line.
func Errorf(format string, args ...interface{}) (int, error) {
	if IsError() {
		format, args = prepareFormatAndArgs(ERR, FunctionWidth, format, args...)
		if !strings.HasSuffix(format, "\n") && !strings.HasSuffix(format, "\r") {
			format = format + "\n"
		}
		return logErrorf(GetStream(), format, args...)
	}
	return 0, nil
}

// Println is a raw version of the debug functions; it tries to interpret
// the message by checking if it starts with anthing like "[D]" or "[W]";
// if so, it delegates to the corresponding logging function, otherwise it
// just prints to the log stream as is, with no additional formatting.
func Println(args ...interface{}) (int, error) {
	if len(args) > 0 {
		if value, ok := args[0].(string); ok {
			switch {
			case strings.HasPrefix(value, "[D]"):
				return Debugln(args[1:]...)
			case strings.HasPrefix(value, "[I]"):
				return Infoln(args[1:]...)
			case strings.HasPrefix(value, "[W]"):
				return Warnln(args[1:]...)
			case strings.HasPrefix(value, "[E]"):
				return Errorln(args[1:]...)
			}
		}
	}
	return fmt.Fprintln(GetStream(), args...)
}

// Printf is a raw version of the debug functions; it tries to interpret
// the message by checking if it starts with anything like "[D]" or "[W]";
// if so, it delegates to the corresponding logging function, otherwise it
// just prints to the log stream as is, with no additional formatting.
func Printf(format string, args ...interface{}) (int, error) {
	re := regexp.MustCompile(`^\[(D|I|W|E)\]\s`)
	switch {
	case strings.HasPrefix(format, "[D]"):
		return Debugf(re.ReplaceAllString(format, ""), args...)
	case strings.HasPrefix(format, "[I]"):
		return Infof(re.ReplaceAllString(format, ""), args...)
	case strings.HasPrefix(format, "[W]"):
		return Warnf(re.ReplaceAllString(format, ""), args...)
	case strings.HasPrefix(format, "[E]"):
		return Errorf(re.ReplaceAllString(format, ""), args...)
	}
	return fmt.Fprintf(GetStream(), format, args...)
}

func prepareFormatAndArgs(level Level, length int, format string, args ...interface{}) (string, []interface{}) {

	leadFormat := "%s %s - "
	tailFormat := ""
	leadArgs := []interface{}{level.String(), time.Now().Format(GetTimeFormat())}
	tailArgs := []interface{}{}

	flags := GetFlags()
	if flags != 0 {
		var fun, file string
		var line int
		pc, file, line, ok := runtime.Caller(2)
		if !ok {
			fun = "<unknown function>"
			file = "<unknown source>"
			line = 0
		} else {
			if flags&FlagFunctionInfo != 0 {
				f := runtime.FuncForPC(pc)
				if f == nil {
					fun = "<unknown function>"
				} else {
					fun = f.Name()
				}
				fun = fun[strings.LastIndex(fun, "/")+1:]
				if len(fun) >= 3 && len(fun) > length {
					fun = "..." + fun[len(fun)-length+3:]
				}
				leadFormat = fmt.Sprintf("%s%%-%ds: ", leadFormat, length)
				leadArgs = append(leadArgs, fun)
			}
			if flags&FlagSourceInfo != 0 {
				tailFormat = " (%s:%d)"
				tailArgs = append(tailArgs, []interface{}{file, line}...)
			}
		}
	}
	format = fmt.Sprintf("%s%s%s", leadFormat, format, tailFormat)
	args = append(leadArgs, append(args, tailArgs...)...)
	return format, args
}

// // methodName returns the name of the calling method, assumed to be two stack
// // frames above; the name can be fully qualified, including the whole import
// // path (e.g. "github.com/myrepo/myproject/mypackage.MyMethod").
// func methodName() string {
// 	pc, _, _, _ := runtime.Caller(2)
// 	f := runtime.FuncForPC(pc)
// 	if f == nil {
// 		return "<?>"
// 	}
// 	return f.Name()
// }

// // unqulified ensures that the input string (typically representing a method name)
// // is not qualified with the full import path.
// func unqualified(name string) string {
// 	// index := strings.LastIndex(name, "/")
// 	// if index != -1 {
// 	// 	return name[index:]
// 	// }
// 	// return name
// 	return name[strings.LastIndex(name, "/")+1:]
// }

// // truncateBefore ensures that the input string is no longer than the given
// // length; if it is, it is truncated and prepended with three dots (...).
// func truncateBefore(s string, length int) string {
// 	if len(s) >= 3 && len(s) > length {
// 		s = "..." + s[len(s)-length+3:]
// 	}
// 	return s
// }

// // truncateBefore ensures that the input string is no longer than the given
// // length; if it is, it is truncated and prepended with three dots (...).
// func truncateAfter(s string, length int) string {
// 	if len(s) >= 3 && len(s) > length {
// 		s = s[0:len(s)-length-3] + "..."
// 	}
// 	return s
// }

// ToJSON converts an object into pretty-printed JSON format.
func ToJSON(object interface{}) string {
	if bytes, err := json.MarshalIndent(object, "", "  "); err == nil {
		return string(bytes)
	}
	return ""
}
