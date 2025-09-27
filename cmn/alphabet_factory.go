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

func AlphabetFactory(language string) IAlphabet {
	var na IAlphabet = nil

	switch strings.ToLower(language) {
	case ALPHA_NAME_ENGLISH:
		na = ALPHA_DISK.Clone().WithSpecialCase(ALPHA_DISK.specialCase)

	case ALPHA_NAME_LATIN:
		fallthrough
	case ALPHA_NAME_SPANISH:
		na = ALPHA_DISK_LATIN.Clone().WithSpecialCase(ALPHA_DISK_LATIN.specialCase)

	case ALPHA_NAME_GERMAN:
		na = ALPHA_DISK_GERMAN.Clone().WithSpecialCase(ALPHA_DISK_GERMAN.specialCase)

	case ALPHA_NAME_GREEK:
		na = ALPHA_DISK_GREEK.Clone().WithSpecialCase(ALPHA_DISK_GREEK.specialCase)

	case ALPHA_NAME_CYRILLIC:
		fallthrough
	case ALPHA_NAME_UKRANIAN:
		fallthrough
	case ALPHA_NAME_RUSSIAN:
		na = ALPHA_DISK_CYRILLIC.Clone().WithSpecialCase(ALPHA_DISK_CYRILLIC.specialCase)

	default:
		mlog.ErrorT("unknown alphabet name", mlog.String("Language", language))
		panic("Bad language in factory")
	}

	return na
}
