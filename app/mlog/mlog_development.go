//go:build mlog
// +build mlog

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * MLog Development Build. A full-featured version of MLog.
 *-----------------------------------------------------------------*/
package mlog

import "strings"

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

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

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *				P r i v i l e g e d   L e v e l s
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

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
