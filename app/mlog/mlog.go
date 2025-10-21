/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * MLog is an enhanced version of the log/slog package with logging
 * levels, extra log tags such as String, Bool, YesNo, Int, Byte,
 * Rune and At. But MLog is not for structured logging (JSON, etc.)
 * because I like traditional logging messages without clutter.
 *   Log level can be set in an environment variable and an optional
 * log filename too. It also has a supplementary log file called
 * "catheter" which is not formatted.
 *   If a log file is used, the main log is appended whereas the
 * catheter is truncated.
 *-----------------------------------------------------------------*/
package mlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	defaultPrefix string   = ""
	defaultLevel  LogLevel = LevelError

	tagCATHE string = "[CAT] "
	tagTRACE string = "[TRC] "
	tagDEBUG string = "[DBG] "
	tagINFO  string = "[INF] "
	tagWARN  string = "[WRN] "
	tagERROR string = "[ERR] "
	tagFATAL string = "[DIE] "

	// environment variable that overrides the default (Error) Log Level for CaesarX
	LOG_LEVEL_ENV string = "LOG_LEVEL_CX"
	// environment variable that indicates the log output filename for CaesarX (default stderr)
	LOG_FILE_ENV string = "LOG_FILE_CX"

	// Logging level enumeration
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

var (
	logMutex    sync.Mutex
	minLogLevel LogLevel    = LevelDebug
	ilogger     *log.Logger = nil
	logFile     *os.File    = nil
	catFile     *os.File    = nil
	// UTF8 BOM (Byte Order Mark)
	UTF8_BOM []byte = []byte{0xEF, 0xBB, 0xBF}
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

func init() {
	levelString := os.Getenv(LOG_LEVEL_ENV)
	if levelString != "" {
		minLogLevel = parseLevel(levelString)
	} else {
		minLogLevel = defaultLevel
	}

	const CUSTOM_TIME_FORMAT = "2006-01-02 15:04:05"
	cw := newCustomLogWriter(os.Stderr, CUSTOM_TIME_FORMAT)

	//ilogger = log.New(os.Stderr, defaultPrefix, log.Ldate|log.Ltime|log.Lshortfile)
	ilogger = log.New(os.Stderr, defaultPrefix, log.Ldate|log.Ltime|log.Lmsgprefix)
	ilogger.SetFlags(log.Lmsgprefix)
	outputLogFilename := os.Getenv(LOG_FILE_ENV)
	if len(outputLogFilename) != 0 {
		if fd, err := openLogFile(outputLogFilename, true, false); err != nil {
			ilogger.SetOutput(cw) // fallback to stderr
		} else {
			ilogger.SetOutput(fd)
		}
	} else {
		ilogger.SetOutput(cw)
	}
}

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type LogLevel int

type LevelLogger struct {
	*log.Logger
	MinLevel LogLevel
}

type customLogWriter struct {
	writer io.Writer
	format string
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func newCustomLogWriter(w io.Writer, timeStampFormat string) *customLogWriter {
	return &customLogWriter{writer: w, format: timeStampFormat}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (clw *customLogWriter) Write(p []byte) (n int, err error) {
	timestamp := time.Now().Format(clw.format)
	formattedMessage := fmt.Sprintf("%s %s", timestamp, p)
	return clw.writer.Write([]byte(formattedMessage))
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// opens the log file and outputs the first message to delimit
// multiple application runs.
func openLogFile(filePath string, isMainLog, truncate bool) (*os.File, error) {
	fileFlags := os.O_CREATE | os.O_WRONLY
	if truncate {
		fileFlags |= os.O_TRUNC
	} else {
		fileFlags |= os.O_APPEND
	}

	logFileX, err := os.OpenFile(filePath, fileFlags, 0666)
	if err != nil {
		return nil, err
	}

	const LEADER string = "[BEG]\t> > > >   T h e   B e g i n n i n g   < < < <\n"
	if isMainLog {
		ilogger.Print(string(UTF8_BOM), LEADER)
	} else {
		logFileX.WriteString(string(UTF8_BOM) + LEADER)
	}

	return logFileX, nil
}

// Create a "catheter" log file. It is a supplementary lifeline for
// exceptional logging and contains no format. Use PrintCathether()
// for writing output.
func SetCatheterFile(filename string) bool {
	var err error = nil
	if catFile != nil {
		return false
	}
	catFile, err = openLogFile(filename, false, true)

	return err == nil
}

// CloseLogFiles to close the log file. Call this in a defer statement in your
// main() IF you specified a log filename in the LOG_FILENAME environment var.
// It does nothing if you used SetOutput() with your own file writer.
func CloseLogFiles() {
	const TRAILER string = "[END]\t> > > >   T h e   E n d   < < < <\n"
	if logFile != nil {
		ilogger.Print(TRAILER)
		err := logFile.Close()
		if err != nil {
			ilogger.Printf("Error closing log file: %v", err)
		}
	}

	if catFile != nil {
		catFile.WriteString(TRAILER)
		err := catFile.Close()
		if err != nil {
			ilogger.Printf("Error closing catheter file: %v", err)
		}
	}
}

// parse a string to convert it to a LogLevel value
func parseLevel(s string) LogLevel {
	var lvl LogLevel
	s = strings.Trim(s, " \t")

	switch {
	case strings.EqualFold(s, "trace"):
		lvl = LevelTrace

	case strings.EqualFold(s, "debug"):
		lvl = LevelDebug

	case strings.EqualFold(s, "info"):
		lvl = LevelInfo

	case strings.EqualFold(s, "warning"):
		fallthrough
	case strings.EqualFold(s, "warn"):
		lvl = LevelWarning

	case strings.EqualFold(s, "error"):
		lvl = LevelError

	case strings.EqualFold(s, "fatal"):
		lvl = LevelFatal

	default:
		lvl = LevelFatal
	}

	return lvl
}

// SetLevel sets the current logging level. Unlike log and slog
// the mlog package supports logging levels.
func SetLevel(newLevel LogLevel) LogLevel {
	logMutex.Lock()
	defer logMutex.Unlock()

	oldLevel := minLogLevel
	minLogLevel = newLevel
	return oldLevel
}

// SetPrefix sets the prefix to appear on all log entries
func SetPrefix(prefix string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	ilogger.SetPrefix(prefix)
}

// SetOutput sets the logging output writer instance. By
// default mlog uses stderr.
func SetOutput(w io.Writer) {
	logMutex.Lock()
	defer logMutex.Unlock()

	ilogger.SetOutput(w)
}

// Print to the catheter file.
func PrintCatheter(message string, v ...ILogKeyValuePair) {
	if catFile != nil {
		var sb strings.Builder
		sb.WriteString(tagCATHE)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" " + t.String())
		}

		catFile.WriteString(sb.String() + "\n")
	}
}

// Trace level with variadic parameters
func Trace(v ...any) {
	if minLogLevel <= LevelTrace {
		v1 := append([]any{tagTRACE}, v...)
		ilogger.Print(v1...)
	}
}

// Trace level with format string
func Tracef(format string, v ...any) {
	if minLogLevel <= LevelTrace {
		ilogger.Printf(tagTRACE+format, v...)
	}
}

// Trace level with message and variadic MLog tags.
func TraceT(message string, v ...ILogKeyValuePair) {
	if minLogLevel <= LevelTrace {
		var sb strings.Builder
		sb.WriteString(tagTRACE)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" " + t.String())
		}
		ilogger.Print(sb.String())
	}
}

// Debug level with variadic parameters
func Debug(v ...any) {
	if minLogLevel <= LevelDebug {
		v1 := append([]any{tagDEBUG}, v...)
		ilogger.Print(v1...)
	}
}

// Debug level with format string
func Debugf(format string, v ...any) {
	if minLogLevel <= LevelDebug {
		ilogger.Printf(tagDEBUG+format, v...)
	}
}

// Debug level with message and variadic MLog tags.
func DebugT(message string, v ...ILogKeyValuePair) {
	if minLogLevel <= LevelDebug {
		var sb strings.Builder
		sb.WriteString(tagDEBUG)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" " + t.String())
		}
		ilogger.Print(sb.String())
	}
}

// Information level with variadic parameters
func Info(v ...any) {
	if minLogLevel <= LevelInfo {
		v1 := append([]any{tagINFO}, v...)
		ilogger.Print(v1...)
	}
}

// Information level with format string
func Infof(format string, v ...any) {
	if minLogLevel <= LevelInfo {
		ilogger.Printf(tagINFO+format, v...)
	}
}

// Information level with message and variadic MLog tags.
func InfoT(message string, v ...ILogKeyValuePair) {
	if minLogLevel <= LevelInfo {
		var sb strings.Builder
		sb.WriteString(tagINFO)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" " + t.String())
		}
		ilogger.Print(sb.String())
	}
}

// Warning level with variadic parameters
func Warn(v ...any) {
	if minLogLevel <= LevelWarning {
		v1 := append([]any{tagWARN}, v...)
		ilogger.Print(v1...)
	}
}

// Warning level with format string
func Warnf(format string, v ...any) {
	if minLogLevel <= LevelWarning {
		ilogger.Printf(tagWARN+format, v...)
	}
}

// Warning level with message and variadic MLog tags.
func WarnT(message string, v ...ILogKeyValuePair) {
	if minLogLevel <= LevelWarning {
		var sb strings.Builder
		sb.WriteString(tagWARN)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" " + t.String())
		}
		ilogger.Print(sb.String())
	}
}

// Error level with variadic parameters
func Error(v ...any) {
	if minLogLevel <= LevelError {
		v1 := append([]any{tagERROR}, v...)
		ilogger.Print(v1...)
	}
}

// Error level with format string
func Errorf(format string, v ...any) {
	if minLogLevel <= LevelError {
		ilogger.Printf(tagERROR+format, v...)
	}
}

// Error level with message and variadic MLog tags.
func ErrorT(message string, v ...ILogKeyValuePair) {
	if minLogLevel <= LevelError {
		var sb strings.Builder
		sb.WriteString(tagERROR)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" " + t.String())
		}
		ilogger.Print(sb.String())
	}
}

// Error level limited to the error itself
func ErrorE(err error) {
	if minLogLevel <= LevelError {
		ilogger.Println(tagERROR, err.Error())
	}
}

// Fatal level with variadic parameters and exitCode
// for terminating the application.
func Fatal(exitCode int, v ...any) {
	if minLogLevel <= LevelFatal {
		v1 := append([]any{tagFATAL}, v...)
		ilogger.Print(v1...)
	}

	os.Exit(exitCode)
}

// Trace level with format string and exitCode for terminating
// the application.
func Fatalf(exitCode int, format string, v ...any) {
	if minLogLevel <= LevelFatal {
		ilogger.Printf(tagFATAL+format, v...)
	}

	os.Exit(exitCode)
}

// Fatal level with message and variadic MLog tags.
// it terminates execution with exitCode.
func FatalT(exitCode int, message string, v ...ILogKeyValuePair) {
	if minLogLevel <= LevelFatal {
		var sb strings.Builder
		sb.WriteString(tagFATAL)
		sb.WriteString(message)
		for _, t := range v {
			sb.WriteString(" " + t.String())
		}
		ilogger.Print(sb.String())
	}

	os.Exit(exitCode)
}
