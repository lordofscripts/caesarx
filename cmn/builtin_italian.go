/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (UTF8) The official Italian alphabet. 28 runes 35 bytes.
 *-----------------------------------------------------------------*/
package cmn

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

const (
	alpha_DISK_ITALIAN string = "ABCDEFGHILMNOPQRSTUVZÉÓÀÈÌÒÙ" // Italian 28 runes 35 bytes
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ALPHA_NAME_ITALIAN = "italian"
)

var (
	ALPHA_DISK_ITALIAN *Alphabet = &Alphabet{"Italian", alpha_DISK_ITALIAN, true, false, false, nil, false, "IT"}
)
