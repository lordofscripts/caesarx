/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2023 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Caller's package information useful for logging
 *-----------------------------------------------------------------*/
package mlog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// Package information descriptor.
type PackageInfo struct {
	Package  string
	Filename string
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Fullname returns the full_package_name|full_filename.go
func (p *PackageInfo) Fullname() string {
	return fmt.Sprintf("%s|%s", p.Package, p.Filename)
}

// Base returns the base package name, i.e. "logy" instead of the
// full package name ("lordofscripts/logy")
func (p *PackageInfo) Base() string {
	return filepath.Base(p.Package)
}

// The fmt.Stringer returns the package name
func (p *PackageInfo) String() string {
	return p.Package
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// RetrievePackageInfo gets the package info (full package name & filename).
// By default (no parameters) it returns the caller's package. If a single
// integer number is given, it refers to how many stack frames to retrieve
// to get to the caller. The default (no parameter) is ONE frame back (caller).
// @return (*PackageInfo) nil on error.
func RetrievePackageInfo(frameNr ...FrameNr) *PackageInfo {
	frames := FRAMENR_CALLER
	if len(frameNr) != 0 {
		frames = frameNr[0]
	}

	pc, fileName, _, ok := runtime.Caller(frames)
	if ok {
		pif := &PackageInfo{"", fileName}

		funcName := runtime.FuncForPC(pc).Name()
		lastSlash := strings.LastIndexByte(funcName, '/')
		if lastSlash < 0 {
			lastSlash = 0
		}
		lastDot := strings.LastIndexByte(funcName[lastSlash:], '.') + lastSlash
		if idx := strings.IndexByte(funcName[:lastDot], '.'); idx > -1 {
			pif.Package = funcName[:idx]
		} else {
			pif.Package = funcName[:lastDot]
		}

		return pif
	}

	return nil
}
