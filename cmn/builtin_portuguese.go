/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (UTF8) The official Portuguese alphabet. 38 runes 50 bytes.
 *-----------------------------------------------------------------*/
package cmn

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

const (
	alpha_DISK_PORTUGUESE string = "ABCÇDEFGHIJKLMNOPQRSTUVWXYZÁÉÍÓÚÀÂÊÔÃÕ" // Portuguese 38 runes 50 bytes
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ALPHA_NAME_PORTUGUESE = "portuguese"
)

var (
	ALPHA_DISK_PORTUGUESE *Alphabet = &Alphabet{"Portuguese", alpha_DISK_PORTUGUESE, true, false, false, nil, false, "PT"}
)
