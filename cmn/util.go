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

import (
	"bytes"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Normalize a string by sanitizing the vocals removing the accents.
func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}

/**
 * Check that a string is composed of unique runes (characters)
 */
func HasUniqueRunes(s string) bool {
	charMap := make(map[rune]bool)

	for _, char := range s {
		if charMap[char] {
			return false
		}
		charMap[char] = true
	}
	charMap = nil // tip GC

	return true
}

/**
 * Check that the string is not made of only white-space.
 */
func IsNotBlank(s string) bool {
	return strings.TrimSpace(s) != ""
}

/**
 * For fun, maps every letter in A-Z to a Runic alphabet.
 */
func RuneString(latin string) string {
	const (
		LOOKUP_STD string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		RUNES      string = "\u16ab\u16d2\u16b3\u16de\u16d6\u16a0\u16b7\u16bb\u16c1\u16c3\u16f1\u16da\u16d7\u16be\u16a9\u16c8\u16e9\u16b1\u16cb\u16cf\u16a2\u16a1\u16b9\u16ea\u16e6\u16ce"
	)

	chars := []rune(strings.ToUpper(latin))
	runesLookup := []rune(RUNES)
	result := make([]rune, len(chars))

	for index, char := range chars {
		if strings.ContainsRune(LOOKUP_STD, char) {
			at := strings.IndexRune(LOOKUP_STD, char)
			result[index] = runesLookup[at]
		} else {
			result[index] = char
		}
	}

	return string(result)
}

// Checks whether the string 's' only contains ASCII characters (0..127).
func IsASCII(s string) bool {
	for _, char := range s { // we can range over the string because we don't care about the index
		if char > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// Like IsASCII except it also returns true for character codes
// between 0xC0 (À) and 0xFF (ÿ). And for GO we make the exception
// of the German Es-zet character lowercase ß 223 and uppercase ẞ 7838
func IsExtendedASCII(s string) bool {
	//for _, char := range []rune(s) { // here the index is a rune counter
	for _, char := range s { // here the index is the 1st byte of the (multi-byte) rune
		if char == CHAR_UPPER_ESZET {
			continue
		}
		if char > unicode.MaxASCII {
			if char < '\u00C0' || char > unicode.MaxLatin1 {
				return false
			}
		}
	}
	return true
}

/**
 * Find out whether the string 's' has at least one multi-byte
 * character (GO rune). In that case functions like strings.Index()
 * won't work!
 */
func IsMultiByteString(s string) bool {
	return len(s) != utf8.RuneCountInString(s)
}

/**
 * Returns the CHARACTER index of the given rune within the STRING.
 * This works for multi-byte characters as well. Using strings.Index()
 * on multi-byte strings does NOT return the index of the character.
 * Therefore, when suspecting a multi-byte string use this method.
 * @returns -1 if not found, else the zero-based character position.
 */
func RuneIndex(s string, r rune) int {
	var where int = -1
	if byteIndex := strings.Index(s, string(r)); byteIndex > -1 {
		// so we convert it to rune index here
		where = utf8.RuneCountInString(s[:byteIndex])
	}

	return where
}

func RuneIndexFold(s string, r rune, handlers *SpecialCaseHandler) int { // @note needs test case
	var where int = -1
	var sx string
	var R rune

	if handlers != nil {
		sx = handlers.ToUpperString(s)
		R = handlers.ToUpperRune(r)
	} else {
		sx = strings.ToUpper(s)
		R = unicode.ToUpper(r)
	}

	if byteIndex := strings.Index(sx, string(R)); byteIndex > -1 {
		// so we convert it to rune index here
		where = utf8.RuneCountInString(s[:byteIndex])
	}

	return where
}

/**
 * Get the alphabet rune at the given position.
 * Panics if out of range.
 * Special case Pos=-1 which returns the last rune.
 */
func RuneAt(s string, pos int) rune {
	length := utf8.RuneCountInString(s)

	if pos < 0 && (pos*-1 > length) {
		panic("RuneAt index out of range")
	}

	var result rune
	ruined := []rune(s)
	if pos >= 0 {
		result = ruined[pos]
	} else {
		result = ruined[length+pos]
	}

	return result
}

// Rotate a string to the left N characters, wrapping at the end.
// It is corrected to work with Unicode strings. The length of
// the string is counted in UTF runes, not in bytes/simple characters.
// Additionally, rather than duplicating the code of RotateStringRight()
// it takes a shortcut by using the complementary number. That is
// IF the string has M runes in length, rotating N characters to the
// RIGHT is equivalent to rotating M - N with LEFT rotation.
func RotateStringLeft(s string, shift int) string {
	complementShift := utf8.RuneCountInString(s) - shift
	return RotateStringRight(s, complementShift)
}

// Rotate a string to the right N characters, wrapping at the end.
// It is corrected to work with unicode strings.
func RotateStringRight(s string, shift int) string {
	//alphaSize := len(s)
	alphaSize := utf8.RuneCountInString(s)
	// only rotate right
	if shift < 0 {
		shift = shift * -1
	}
	// bigger than alphabet? then wrap it
	if shift > alphaSize {
		shift = shift % alphaSize
	}
	// no rotation?
	if shift == alphaSize || shift == 0 {
		return s
	}
	// rotate
	//return s[alphaSize-shift:] + s[0:alphaSize-shift]
	runic := []rune(s)
	result := string(runic[alphaSize-shift:])
	result += string(runic[0 : alphaSize-shift])
	return result
}

// Utility function for the Vigenère auto-key functionality.
// This only works with keys and messages composed entirely of
// single-byte runes, i.e. English; therefore, it is not suitable
// for foreign strings with accented characters or extended Unicode
// like Greek, German, Cyrillic, etc. Else use AutoKeyUTF8()
func AutoKeyASCII(plain, key string) string {
	keylen := len(key)
	textlen := len(plain)
	if keylen >= len(plain) {
		return key[0:textlen]
	} else {
		return key + plain[0:textlen-keylen]
	}
}

/**
 * Utility function used for Vigenère auto-key cipher. The key is prepended
 * to the message string and the result is truncated to the length of the
 * original message. Therefore, if the key is the same length as the message
 * the resulting autokey is the key itself.
 * NOTE: This function also works well with keys and messages that contain
 * multi-byte runes. Therefore it works fine with English, Spanish, German,
 * Greek, Cyrillic, etc.
 */
func AutoKeyUTF8(plain, key string) string {
	keylen := utf8.RuneCountInString(key)
	textlen := utf8.RuneCountInString(plain)
	if keylen >= len(plain) {
		//return key[0:textlen]
		return string([]rune(key)[0:textlen])
	} else {
		//return key + plain[0:textlen-keylen]
		return key + string([]rune(plain)[0:textlen-keylen])
	}
}

/**
 * Insert a separator character every N (N>0) characters in string s.
 * It also works with strings that contain multi-byte runes.
 */
func InsertNth(s string, n uint, char rune) string {
	if n == 0 {
		return s
	}

	var buffer bytes.Buffer
	var n1 = int(n - 1)
	var l1 = utf8.RuneCountInString(s) - 1
	letters := []rune(s)
	for i, rune := range letters {
		buffer.WriteRune(rune)
		if i%int(n) == n1 && i != l1 {
			buffer.WriteRune(char)
		}
	}

	return buffer.String()
}

// Locate() gets the index of 'char' within string 'in' but it does
// so looking for multi-byte runes, not plain ASCII. The standard
// Index() functions in strings and unicode return the index of the
// first byte of the (multi-byte) rune. That value is unsuitable for
// looking up characters.
func Locate(char rune, in string) int {
	var at int = -1
	char = unicode.ToUpper(char) // @audit GO v1.25 ToUpper fails on some runes like ß
	for index, arune := range []rune(in) {
		if unicode.ToUpper(arune) == char {
			at = index
			break
		}
	}
	return at
}

// from golang.org/x/exp/slices.
// Compact replaces consecutive runs of equal elements with a single copy.
// It is like the uniq command found on Unix. Compact modifies the contents
// of the slice s; it does not create a new slice.
func Compact[S ~[]E, E comparable](s S) S {
	if len(s) == 0 {
		return s
	}

	// remove adjacent duplicates
	i := 1
	last := s[0]
	for _, v := range s[1:] {
		if v != last {
			s[i] = v
			i++
			last = v
		}
	}
	return s[:i]
}

func IntersectInt(set1, set2 []int) []int {
	sort.Ints(set1)
	sort.Ints(set2)

	common := make([]int, 0)

	for i, j := 0, 0; i < len(set1) && j < len(set2); {
		if set1[i] == set2[j] {
			common = append(common, set2[j])
			i++
			j++
		} else if set1[i] < set2[j] {
			i++
		} else {
			j++
		}
	}

	return common
}

// Removes all spaces and tabs from S, then after every Nth character
// it inserts the CHAR rune.
// Example: ToMessageTape("N buqx sjf dnom ytne kjmw 7<7@", 5, '·')
// Result : Nbuqx·sjfdn·omytn·ekjmw·7<7@
func ToMessageTape(s string, n uint, char rune) string {
	cleaned := strings.ReplaceAll(s, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "\t", "")
	return InsertNth(cleaned, n, char)
}

func Bigram(s string, char rune) string {
	return InsertNth(s, 2, char)
}

func Trigram(s string, char rune) string {
	return InsertNth(s, 3, char)
}

func Quartets(s string, char rune) string {
	return InsertNth(s, 4, char)
}

func Quintets(s string, char rune) string {
	return InsertNth(s, 5, char)
}
