/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2023 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package mlog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	FRAMENR_THIS   FrameNr = 1
	FRAMENR_CALLER FrameNr = 2
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/
type CallerInfo struct {
	packageN  string
	structure string
	function  string
	filename  string
	lineno    int
}

type FrameNr = int

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Source file information.
// @returns (string) filename.go:line#
func (c *CallerInfo) SourceInfo() string {
	base := filepath.Base(c.filename)
	return fmt.Sprintf("%s:%d", base, c.lineno)
}

// Source object information.
// @returns (string) package.struct{}.method() OR package.function()
func (c *CallerInfo) ObjectInfo() string {
	const SEP string = "."
	var sb strings.Builder

	// clean package name
	sb.WriteString(c.packageN + SEP)

	// object if any (a method info)
	if len(c.structure) > 0 {
		sb.WriteString(c.structure + "{}.")
	}

	// function or method
	sb.WriteString(c.function + "()")
	return sb.String()
}

// implements the fmt.Stringer interface for *CallerInfo . For non-structs
// (functions) it returns "package.function()" and for structs the return
// value is in the format "package.object{}.method()".
func (c *CallerInfo) String() string {
	if len(c.structure) == 0 {
		//return fmt.Sprintf("%s.%s()", c.packageN, c.function)
		return c.StringF("%B")
	} else {
		return c.StringF("%c") // or %C
	}
}

// (CallerInfo) Formatted stringify takes any of the following format specifiers:
// %P package name
// %S structure name (or empty if none), i.e. Event{}
// %F or %M function/method name (aliased), i.e. Sum()
// %L line nr.
// %A short version (%P)
// %B median version (%P.%F)
// %C long version (%P.%S.%M)
func (c *CallerInfo) StringF(format string) string {
	const (
		IFUNC = "()"
		ISTRU = "{}"
	)
	//out := strings.ToUpper(format)
	out := format
	// macro replacements
	out = strings.Replace(out, "%A", "%P", 1)
	out = strings.Replace(out, "%B", "%P.%F#%L", 1)
	out = strings.Replace(out, "%C", "%P.%S.%M#%L", 1)
	out = strings.Replace(out, "%c", "%p.%S.%M#%L", 1)
	var struN string = ""
	if len(c.structure) > 0 {
		struN = c.structure + ISTRU
	}
	// atomic replacements
	out = strings.Replace(out, "%P", c.packageN, 1)
	out = strings.Replace(out, "%S", struN, 1)
	out = strings.Replace(out, "%F", c.function+IFUNC, 1)
	out = strings.Replace(out, "%M", c.function+IFUNC, 1)
	out = strings.Replace(out, "%L", strconv.Itoa(c.lineno), 1)
	if strings.Contains(out, "%p") && strings.LastIndexByte(c.packageN, '/') > 0 {
		pname := c.packageN[strings.LastIndexByte(c.packageN, '/')+1:]
		out = strings.Replace(out, "%p", pname, 1)
	}
	return out
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Retrieve caller info returnin these values:
// @return string pkg : package name
// @return string stru: structure name if caller is a method, else ""
// @return string fun : function or method name
func RetrieveCallerInfo(frame FrameNr) *CallerInfo {
	const FrameLevel = 1 // 1 for in-file demo, 2 from elsewhere
	if frame < FRAMENR_THIS {
		frame = FrameLevel
	}
	pc, fileName, lineNo, ok := runtime.Caller(frame)

	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		//fmt.Println(funcName)
		lastSlash := strings.LastIndexByte(funcName, '/')
		if lastSlash < 0 {
			lastSlash = 0
		}
		lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash

		//var title string
		var pkg, stru, fun string
		//fmt.Printf("RCI ***%s\n", funcName)
		if idx := strings.IndexByte(funcName[:lastDot], '.'); idx > -1 {
			// Value struct:  Event
			// Pointer to struct: (*Event)
			// Package Init(): init
			stru = funcName[idx+1 : lastDot]
			pkg = funcName[:idx]
			if stru == "init" {
				// A package.init() comes as stru:init func:0
				stru = ""
			}
			//title = "Method "
		} else {
			stru = ""
			pkg = funcName[:lastDot]
			//title = "Func   "
		}

		fun = strings.Trim(funcName[lastDot+1:], " ")
		if fun == "0" {
			fun = "init"
		}
		ci := &CallerInfo{packageN: pkg, structure: stru, function: fun, filename: fileName, lineno: lineNo}

		//fmt.Printf("  Package: %s\n", funcName[:firstDot])
		//fmt.Printf("  Package: %s\n", pkg)
		//fmt.Printf("  Object : %s\n", stru)
		//fmt.Printf("  %s: %s\n", title, fun)
		return ci
	}

	return nil
}
