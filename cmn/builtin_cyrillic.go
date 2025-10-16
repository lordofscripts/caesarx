/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (UTF8) The official Cyrillic alphabet. 33 runes 66 bytes.
 *-----------------------------------------------------------------*/
package cmn

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

const (
	alpha_DISK_CYRILLIC string = "АБВГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ" // Cyrillic 33 runes 66 bytes
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ALPHA_NAME_CYRILLIC  = "cyrillic"
	ALPHA_NAME_UKRAINIAN = "ukranian"
	ALPHA_NAME_RUSSIAN   = "russian"
)

var (
	ALPHA_DISK_CYRILLIC *Alphabet = &Alphabet{"Cyrillic", alpha_DISK_CYRILLIC, true, true, false, nil, false, "RU"}
)
