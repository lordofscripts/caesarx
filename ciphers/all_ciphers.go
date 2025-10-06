/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package ciphers

import (
	"fmt"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var (
	availableCiphers map[string]*CipherInfo = make(map[string]*CipherInfo)
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type CipherInfo struct {
	code        string
	name        string
	description string
	author      string
	revision    string
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewCipherInfo(code, revision, author, name, description string) *CipherInfo {
	return &CipherInfo{
		code:        strings.ToUpper(code),
		name:        name,
		description: description,
		author:      author,
		revision:    revision,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (a *CipherInfo) String() string {
	return fmt.Sprintf("[%5s] v%s (%s) '%s' - %s", a.code, a.revision, a.author, a.name, a.description)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/**
 * Registers a cipher
 */
//func RegisterCipher(name, code, descr string) {
func RegisterCipher(descriptor *CipherInfo) {
	if _, exists := availableCiphers[descriptor.code]; !exists {
		availableCiphers[descriptor.code] = descriptor
	}
}

/**
 * Produce a formatted list of each registered cipher and
 * their summary.
 */
func PrintAvailableCiphers() string {
	var sb strings.Builder

	for _, v := range availableCiphers {
		sb.WriteString(v.String() + "\n")
	}

	return sb.String()
}
