package tests

import (
	"lordofscripts/caesarx/cmn"
	"strings"
)

var (
	AllAlphabets = []*cmn.Alphabet{ // all built-in alphabets in order of complexity
		cmn.ALPHA_DISK,
		cmn.ALPHA_DISK_LATIN,
		cmn.ALPHA_DISK_GERMAN,
		cmn.ALPHA_DISK_GREEK,
		cmn.ALPHA_DISK_CYRILLIC,
	}
)

func IsEnglish(a *cmn.Alphabet) bool {
	return strings.EqualFold(a.Chars, cmn.ALPHA_DISK.Chars)
}

func IsSpanish(a *cmn.Alphabet) bool {
	return strings.EqualFold(a.Chars, cmn.ALPHA_DISK_LATIN.Chars)
}

func IsGerman(a *cmn.Alphabet) bool {
	return strings.EqualFold(a.Chars, cmn.ALPHA_DISK_GERMAN.Chars)
}

func IsGreek(a *cmn.Alphabet) bool {
	return strings.EqualFold(a.Chars, cmn.ALPHA_DISK_GREEK.Chars)
}

func IsCyrillic(a *cmn.Alphabet) bool {
	return strings.EqualFold(a.Chars, cmn.ALPHA_DISK_CYRILLIC.Chars)
}
