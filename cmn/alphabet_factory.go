/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package cmn

import (
	"lordofscripts/caesarx/app/mlog"
	"strings"
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// An alphabet factory that given an alphabet name *_NAME_* it returns
// the appropriate alphabet if there is a match. It supports all
// built-in alphabets: English, Latin (Spanish), German, Greek,
// Cyrillic (Ukrainian, Russian), Arabic numbers, Hindi numbers,
// Extended numbers, Symbols and Punctuation.
func AlphabetFactory(language string) IAlphabet {
	var na IAlphabet = nil

	switch strings.ToLower(language) {
	case ALPHA_NAME_ENGLISH: // English  26 runes 26 bytes
		na = ALPHA_DISK.Clone().WithSpecialCase(ALPHA_DISK.specialCase)

	case ALPHA_NAME_LATIN: // Latin 33 runes 40 bytes
		fallthrough
	case ALPHA_NAME_SPANISH:
		na = ALPHA_DISK_LATIN.Clone().WithSpecialCase(ALPHA_DISK_LATIN.specialCase)

	case ALPHA_NAME_GERMAN: // German 30 runes 34 bytes
		na = ALPHA_DISK_GERMAN.Clone().WithSpecialCase(ALPHA_DISK_GERMAN.specialCase)

	case ALPHA_NAME_GREEK: // Greek 24 runes 48 bytes
		na = ALPHA_DISK_GREEK.Clone().WithSpecialCase(ALPHA_DISK_GREEK.specialCase)

	case ALPHA_NAME_CYRILLIC: // Cyrillic 33 runes 66 bytes
		fallthrough
	case ALPHA_NAME_UKRAINIAN:
		fallthrough
	case ALPHA_NAME_RUSSIAN:
		na = ALPHA_DISK_CYRILLIC.Clone().WithSpecialCase(ALPHA_DISK_CYRILLIC.specialCase)

	case ALPHA_NAME_NUMBERS_ARABIC: // Numbers 10 runes 10 bytes
		na = NUMBERS_DISK

	case ALPHA_NAME_NUMBERS_EASTERN: // Eastern Numbers 10 runes 20 bytes
		na = NUMBERS_EASTERN_DISK

	case ALPHA_NAME_NUMBERS_ARABIC_EXTENDED: // Numbers Extended 17 runes 17 bytes
		na = NUMBERS_DISK_EXT

	case ALPHA_NAME_SYMBOLS: // Symbols 36 runes 38 bytes
		na = SYMBOL_DISK

	case ALPHA_NAME_PUNCTUATION: // Punctuation 26 runes 28 bytes
		na = PUNCTUATION_DISK

	default:
		mlog.ErrorT("unknown alphabet name", mlog.String("Language", language))
		panic("Bad language in factory")
	}

	return na
}
