/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (UTF8) The official Czech alphabet. 41 runes 56 bytes.
 *-----------------------------------------------------------------*/
package cmn

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

const (
	alpha_DISK_CZECH string = "ABCČDĎEFGHIJKLMNŇOPQRŘSŠTŤUVWXYÝZŽÁÉÍÓÚĚŮ" // Czech 41 runes 56 bytes
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ALPHA_NAME_CZECH = "czech"
)

var (
	ALPHA_DISK_CZECH *Alphabet = &Alphabet{"Czech", alpha_DISK_CZECH, true, false, false, nil, false, "CZ"}
)
