/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (UTF8) The official Spanish alphabet. 33 runes 40 bytes.
 * It is common to many Latin alphabets.
 *-----------------------------------------------------------------*/
package cmn

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

const (
	alpha_DISK_LATIN string = "ABCDEFGHIJKLMNÑOPQRSTUVWXYZÁÉÍÓÚÜ" // Latin 33 runes 40 bytes
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ALPHA_NAME_LATIN   = "latin"
	ALPHA_NAME_SPANISH = "spanish"
)

var (
	ALPHA_DISK_LATIN *Alphabet = &Alphabet{"Spanish", alpha_DISK_LATIN, true, false, false, nil, false, "ES"}
)
