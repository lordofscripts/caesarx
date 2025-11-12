/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * NOTE: What I label as a Punctuation alphabet contains a mix of
 * what the Unicode library considers Punctuation & Symbols. For
 * example $+<=> are symbols whereas ¡!"#%&'()*,-./:;¿?@ are
 * punctuation.
 *-----------------------------------------------------------------*/
package cmn

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

const (
	ALPHA_NAME_SYMBOLS string = "symbols"

	// Punctuation: ¡!"#%&'()*,-./:;¿?@
	// Symbol     : $+<=>
	symbol_DISK string = "¡!\"#$%&'()*+,-./0123456789:;<=>¿?@[]" // Symbols 36 runes 38 bytes
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var (
	// The Symbols disk contains Decimal numerals, some punctuation and some symbols.
	// It is a more complete version of the Western Numerals Extended disk.
	SYMBOL_DISK *Alphabet = &Alphabet{"Symbols", symbol_DISK, false, false, true, nil, false, PSO_PUNCT_DEC}
)
