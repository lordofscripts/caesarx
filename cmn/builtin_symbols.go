//go:build exclude

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
	// Punctuation: ¡!"#%&'()*,-./:;¿?@
	// Symbol     : $+<=>
	symbol_DISK string = "¡!\"#$%&'()*+,-./0123456789:;<=>¿?@[]"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var (
	SYMBOL_DISK *Alphabet = &Alphabet{"Symbols", symbol_DISK, false, false, true}
)
