/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * All useful punctuation characters.
 *-----------------------------------------------------------------*/
package cmn

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

const (
	ALPHA_NAME_PUNCTUATION string = "punctuation"

	punctuation_DISK string = "¡!\"#$%&'()*+,-./:;<=>¿?@[]" // Punctuation 26 runes 28 bytes
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var (
	// The Punctuation disk only contains the most common punctuation characters.
	// It does not contain any numerals.
	PUNCTUATION_DISK *Alphabet = &Alphabet{"Punctuation", punctuation_DISK, false, false, true, nil, false, PSO_PUNCT}
)
