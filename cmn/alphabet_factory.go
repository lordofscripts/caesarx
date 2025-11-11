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

const (
	// CLI alphabet composer concatenation operator
	ALPHA_COMPOSER_SEP string = "+"
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

// When given a string containing an alphabet, it looks
// up all built-in alphabets and returns the correct
// alphabet instance. The comparison of alphabet characters
// is case-insensitive but the order of the characters
// must be the same.
func IdentifyAlphabet(alphaStr string) *Alphabet {
	var all = []*Alphabet{
		ALPHA_DISK,
		ALPHA_DISK_LATIN,
		ALPHA_DISK_GERMAN,
		ALPHA_DISK_GREEK,
		ALPHA_DISK_CYRILLIC,
		NUMBERS_DISK,
		NUMBERS_DISK_EXT,
		NUMBERS_EASTERN_DISK,
		SYMBOL_DISK,
		PUNCTUATION_DISK,
	}

	for _, candidate := range all {
		if strings.EqualFold(candidate.Chars, alphaStr) {
			return candidate
		}
	}

	return nil
}

// the spec contains a list of built-in alphabet names separated by "+"
// which are then used to compose a single alphabet. If there is just one
// then that is used. It verifies the composition has no duplicates.
func AlphabetComposer(spec string) IAlphabet {
	var alpha IAlphabet = nil

	if strings.Contains(spec, ALPHA_COMPOSER_SEP) {
		// get all valid alphabet names in user input
		alphabet_names := strings.Split(spec, ALPHA_COMPOSER_SEP)
		var allRunes, name, iso string = "", "", ""
		for _, name := range alphabet_names {
			candidate := AlphabetFactory(name)
			if candidate != nil {
				// collect individual runes in a single alphabet
				allRunes = allRunes + candidate.Clone().Chars
				isoCode := candidate.LangCodeISO()
				if isoCode != "" {
					name = candidate.Clone().Name + " (Composed)"
					iso = isoCode
				}
			}
		}

		// check there are no duplicate runes
		if !HasUniqueRunes(allRunes) {
			mlog.Error("There are duplicate runes in the concatenation of ", mlog.String("Alphas", spec))
		} else {
			alpha = NewAlphabet("Custom", allRunes, false, false)
			alpha.WithLangCode(iso)
			alpha.Rename(name)
		}
	} else {
		alpha = AlphabetFactory(spec)
	}

	return alpha
}
