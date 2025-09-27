/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package caesar

import (
	"fmt"
	"lordofscripts/caesarx/cmn"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type CaesarOptions struct {
	Variant      CaesarCipherMode
	Initial      rune
	Supplemental uint
	CleanAccents bool
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// CAESAR
// Standard ancient Caesar cipher. ALPHA_DISK only.
func NewCaesarOpts(key rune) *CaesarOptions {
	return &CaesarOptions{CAESAR, key, 0, false}
}

// CAESAR_EXTENDED
// Extended Caesar cipher with support for decimal digits and the most-used
// symbols & punctuations. Uses the SYMBOL_DISK and a single ALPHA_DISK.
func NewCaesarExtendedOpts(key rune) *CaesarOptions {
	return &CaesarOptions{CAESAR_EXTENDED, key, 0, false}
}

// CAESAR_AUGUSTUS
// Extended Caesar cipher with same character set as CAESAR_EXTENDED but the
// Inner disk oscillates back and forth between the Initial key and the
// Initial key plus an offset depending on the odd/even input character position.
// Uses the SYMBOL_DISK and a single oscillating ALPHA_DISK.
func NewCaesarAugustusOpts(key rune, alternateOffset uint) *CaesarOptions {
	return &CaesarOptions{CAESAR_AUGUSTUS, key, alternateOffset, false}
}

// CAESAR_AUGUSTUS
// Extended Caesar cipher with same character set as CAESAR_EXTENDED but the
// Inner disk oscillates back and forth between the Initial (outer) key and the
// Initial key plus an offset (innerKey) depending on the odd/even input character position.
// Uses the SYMBOL_DISK and a single oscillating ALPHA_DISK.
func NewCaesarTiberiusOpts(outerKey rune, innerKey uint) *CaesarOptions { // @audit does not work with Foreign alphabets as-is
	return &CaesarOptions{CAESAR_TIBERIUS, outerKey, innerKey % cmn.ALPHA_DISK.Size(), false}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (c *CaesarOptions) String() string {
	var readable = ""

	switch c.Variant {
	case CAESAR, CAESAR_EXTENDED:
		readable = fmt.Sprintf("%s (Key = %d)", c.Variant, c.Initial)

	case CAESAR_AUGUSTUS, CAESAR_TIBERIUS:
		readable = fmt.Sprintf("%s (Key = %d, Suplement = %d)", c.Variant, c.Initial, c.Supplemental)

	}

	return readable
}

func (c *CaesarOptions) LeaderString() string {
	var prefix []string = []string{"XXX", "CAES", "CAEX", "CAEO", "CAET"}
	var sidx string = ""
	if c.Variant == CAESAR_AUGUSTUS || c.Variant == CAESAR_TIBERIUS {
		sidx = fmt.Sprintf("%02d", c.Supplemental)
	}

	return fmt.Sprintf("%s%02d%s", prefix[int(c.Variant)], c.Initial, sidx)
}
