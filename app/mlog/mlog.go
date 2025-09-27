/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
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

	tagTRACE string = "[TRC] "
	tagDEBUG string = "[DBG] "
	tagINFO  string = "[INF] "
	tagWARN  string = "[WRN] "
	tagERROR string = "[ERR] "
	tagFATAL string = "[DIE] "

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
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

func init() {
	levelString := os.Getenv("LOG_LEVEL")
	if levelString != "" {
		minLogLevel = parseLevel(levelString)
	} else {
		minLogLevel = defaultLevel
	}

	const CUSTOM_TIME_FORMAT = "2006-01-02 15:04:05"
	cw := newCustomLogWriter(os.Stderr, CUSTOM_TIME_FORMAT)

	//ilogger = log.New(os.Stderr, defaultPrefix, log.Ldate|log.Ltime|log.Lshortfile)
	ilogger = log.New(os.Stderr, defaultPrefix, log.Ldate|log.Ltime|log.Lmsgprefix)
	ilogger.SetOutput(cw)
	ilogger.SetFlags(log.Lmsgprefix)
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

func SetLevel(newLevel LogLevel) LogLevel {
	logMutex.Lock()
	defer logMutex.Unlock()

	oldLevel := minLogLevel
	minLogLevel = newLevel
	return oldLevel
}

func SetPrefix(prefix string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	ilogger.SetPrefix(prefix)
}

func SetOutput(w io.Writer) {
	logMutex.Lock()
	defer logMutex.Unlock()

	ilogger.SetOutput(w)
}

func Trace(v ...any) {
	if minLogLevel <= LevelTrace {
		v1 := append([]any{tagTRACE}, v...)
		ilogger.Print(v1...)
	}
}

func Tracef(format string, v ...any) {
	if minLogLevel <= LevelTrace {
		ilogger.Printf(tagTRACE+format, v...)
	}
}

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

func Debug(v ...any) {
	if minLogLevel <= LevelDebug {
		v1 := append([]any{tagDEBUG}, v...)
		ilogger.Print(v1...)
	}
}

func Debugf(format string, v ...any) {
	if minLogLevel <= LevelDebug {
		ilogger.Printf(tagDEBUG+format, v...)
	}
}

// @note actually, since all the tags are ILogKeyValuePairs and that
// interface includes the fmt.Stringer interface, we can also use the
// String/Int/Rune/Bool/YesNo tags as variadic parameters :)
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

func Info(v ...any) {
	if minLogLevel <= LevelInfo {
		v1 := append([]any{tagINFO}, v...)
		ilogger.Print(v1...)
	}
}

func Infof(format string, v ...any) {
	if minLogLevel <= LevelInfo {
		ilogger.Printf(tagINFO+format, v...)
	}
}

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

func Warn(v ...any) {
	if minLogLevel <= LevelWarning {
		v1 := append([]any{tagWARN}, v...)
		ilogger.Print(v1...)
	}
}

func Warnf(format string, v ...any) {
	if minLogLevel <= LevelWarning {
		ilogger.Printf(tagWARN+format, v...)
	}
}

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

func Error(v ...any) {
	if minLogLevel <= LevelError {
		v1 := append([]any{tagERROR}, v...)
		ilogger.Print(v1...)
	}
}

func Errorf(format string, v ...any) {
	if minLogLevel <= LevelError {
		ilogger.Printf(tagERROR+format, v...)
	}
}

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

func ErrorE(err error) {
	if minLogLevel <= LevelError {
		ilogger.Println(tagERROR, err.Error())
	}
}

func Fatal(exitCode int, v ...any) {
	if minLogLevel <= LevelFatal {
		v1 := append([]any{tagFATAL}, v...)
		ilogger.Print(v1...)
	}

	os.Exit(exitCode)
}

func Fatalf(exitCode int, format string, v ...any) {
	if minLogLevel <= LevelFatal {
		ilogger.Printf(tagFATAL+format, v...)
	}

	os.Exit(exitCode)
}

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
