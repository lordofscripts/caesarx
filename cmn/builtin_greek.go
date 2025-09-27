/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (UTF8) The official Greek alphabet. 24 runes 48 bytes.
 *-----------------------------------------------------------------*/
package cmn

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

const (
	alpha_DISK_GREEK string = "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ" // Greek 24 runes 48 bytes
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ALPHA_NAME_GREEK = "greek"
)

var (
	ALPHA_DISK_GREEK *Alphabet = &Alphabet{"Greek", alpha_DISK_GREEK, true, true, false, nil, "GR"}
)
