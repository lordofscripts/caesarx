/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (ASCII) The official English alphabet. 26 runes 26 bytes.
 * Just plain letters without accent. The simple life of ASCII.
 *-----------------------------------------------------------------*/
package cmn

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

const (
	alpha_DISK string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" // English  26 runes 26 bytes
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ALPHA_NAME_ENGLISH = "english"
)

var (
	ALPHA_DISK *Alphabet = &Alphabet{"English", alpha_DISK, false, false, false, nil, false, ISO_EN}
)
