//go:build !mlog
// +build !mlog

/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * MLog Release Build. A stripped-off version of MLog which only
 * logs errors, fatals, and warnings with everything else disabled.
 *-----------------------------------------------------------------*/
package mlog

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// The "catheter" feature is not enabled.
func SetCatheterFile(filename string) bool {
	return false
}

// The "catheter" feature is not enabled.
func PrintCatheter(message string, v ...ILogKeyValuePair) {}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *				P r i v i l e g e d   L e v e l s
 *- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -*/

// Trace level with variadic parameters
func Trace(v ...any) {}

// Trace level with format string
func Tracef(format string, v ...any) {}

// Trace level with message and variadic MLog tags.
func TraceT(message string, v ...ILogKeyValuePair) {}

// Debug level with variadic parameters
func Debug(v ...any) {}

// Debug level with format string
func Debugf(format string, v ...any) {}

// Debug level with message and variadic MLog tags.
func DebugT(message string, v ...ILogKeyValuePair) {}

// Information level with variadic parameters
func Info(v ...any) {}

// Information level with format string
func Infof(format string, v ...any) {}

// Information level with message and variadic MLog tags.
func InfoT(message string, v ...ILogKeyValuePair) {}
