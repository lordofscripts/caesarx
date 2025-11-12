/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (UTF8) The official German alphabet. 30 runes 34 bytes.
 * NOTE: Contrary to other programming languages, in GO the 'ß'
 * (es-zet) character represents both lowercase and uppercase!
 * In the other languages 'SS' is the uppercase but that would break
 * logic because it has 2 runes instead of 1. Therefore in GO
 * unicode.IsLower('ß') returns true.
 *-----------------------------------------------------------------*/
package cmn

import (
	"strings"
	"unicode"
)

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

// **** @note IMPORTANT DEVELOPER NOTES FOR ALL END-USERS ****
// This has an uppercase ẞ which looks very similar to the
// lowercase ß. The GO strings.ToUpper() does NOT make that
// conversion for THAT character! After spending a LOT of time
// debugging, I found it out! I had to take special measures in
// Alphabet.ToUpper()
const (
	alpha_DISK_GERMAN string = "ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜẞ" // German 30 runes 34 bytes UPPERCASES!
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ALPHA_NAME_GERMAN  = "german"
	CHAR_UPPER_ESZET   = 'ẞ'
	CHAR_LOWER_ESZET   = 'ß'
	STRING_UPPER_ESZET = "ẞ"
	STRING_LOWER_ESZET = "ß"
)

var (
	// @note when some are not given Normalize() must be called.
	GermanSpecialCase *SpecialCaseHandler = &SpecialCaseHandler{
		ToUpperString: germanToUpperString,
		ToLowerString: germanToLowerString,
		ToUpperRune:   germanToUpperRune,
		ToLowerRune:   germanToLowerRune,
		IsUpperRune:   unicode.IsUpper,
		IsLowerRune:   unicode.IsLower,
	}
	ALPHA_DISK_GERMAN *Alphabet = &Alphabet{"German", alpha_DISK_GERMAN, true, false, false, GermanSpecialCase, false, ISO_DE}
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/**
 * Workaround for GO strings package bug (submitted 2 Sep 2025).
 */
func germanToUpperString(s string) string {
	// Another Hard-to-find Quirk BUG! strings.ToUpper('ß')
	// Found in DuckDuck.com Search Assist:
	// "The strings.ToUpper() function in Go does not convert the German
	// character "ß" to uppercase because it traditionally does not have
	// an uppercase equivalent. However, since 2017, the uppercase "ẞ" is
	// accepted in some contexts, but it is not automatically produced by
	// the ToUpper() function.
	// @note GO v1.25 strings.ToUpper() does NOT convert ß (lower) into ẞ (upper)
	s = strings.ToUpper(s)

	if strings.Contains(s, STRING_LOWER_ESZET) {
		// Now handle the exception
		// In GO strings.ToUpper(LowerCaseEsTzet) == LowerCaseEsTzet !!!
		s = strings.ReplaceAll(s, STRING_LOWER_ESZET, STRING_UPPER_ESZET)
	}

	return s
}

/**
 * Workaround for GO unicode package bug (submitted 2 Sep 2025).
 */
func germanToUpperRune(r rune) rune {
	// in GO v1.25 unicode.ToUpper(ß) does not yield ẞ (upper!), not always!
	r = unicode.ToUpper(r)
	if r == CHAR_LOWER_ESZET { // @audit GO Bug Workaround
		r = CHAR_UPPER_ESZET
	}

	return r
}

// Just in case although it is not broken
func germanToLowerString(s string) string {
	s = strings.ToLower(s)

	if strings.Contains(s, STRING_UPPER_ESZET) {
		s = strings.ReplaceAll(s, STRING_UPPER_ESZET, STRING_LOWER_ESZET)
	}

	return s
}

// Just in case although it is not broken
func germanToLowerRune(r rune) rune {
	r = unicode.ToLower(r)
	if r == CHAR_UPPER_ESZET {
		r = CHAR_LOWER_ESZET
	}

	return r
}
