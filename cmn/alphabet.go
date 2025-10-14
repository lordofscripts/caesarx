/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Alphabet is the basis of all plain text.
 *-----------------------------------------------------------------*/
package cmn

import (
	"fmt"
	"lordofscripts/caesarx/app/mlog"
	"strings"
	"unicode"
	"unicode/utf8"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	CaseInsensitive bool = true
)

// For all alphabets where the standard library works. German
// however needs GermanSpecialCase as pre-configured in the
// built-in German alphabet of this module.
var NoSpecialCase SpecialCaseHandler = SpecialCaseHandler{
	ToUpperString: strings.ToUpper,
	ToLowerString: strings.ToLower,
	ToUpperRune:   unicode.ToUpper,
	ToLowerRune:   unicode.ToLower,
	IsUpperRune:   unicode.IsUpper,
	IsLowerRune:   unicode.IsLower,
}

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type IAlphabet interface {
	WithLangCode(string) *Alphabet
	WithSpecialCase(handlers *SpecialCaseHandler) *Alphabet
	BorrowSpecialCase() *SpecialCaseHandler
	Size() uint
	ToUpper() *Alphabet
	ToLower() *Alphabet
	PositionOf(rune) int
	GetRuneAt(int) rune
	Rename(string) *Alphabet
	Renumber(uint) *Alphabet
	Contains(r rune, ignoreCase bool) bool
	From(string) *Alphabet
	Rotate(rotateQty int) *Alphabet
	Check() bool
	IsBinary() bool
	Clone() *Alphabet

	LangCodeISO() string
}

var _ IAlphabet = (*Alphabet)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// Any alphabet. Those not based on the standard ASCII English
// alphabet have the 'Foreign' property set.
type Alphabet struct {
	Name        string // usually the alphabet's language (or Symbol class)
	Chars       string // the actual unicode character set
	Foreign     bool   // whether it represents a non-English character set (accents, multi-byte)
	Unicode     bool   // whether it contains Unicode characters (not just ASCII)
	OnlySymbols bool   // whether it is purely symbols, punctuation, digits, etc.
	specialCase *SpecialCaseHandler
	isBinary    bool
	langCode    string
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// Alphabet constructor
func NewAlphabet(name, letters string, isEnglish, isPureSymbols bool) *Alphabet {
	return &Alphabet{name, letters, !isEnglish, !IsExtendedASCII(letters), false, nil, false, ""}
}

// Alphabet constructor. Use exclusively for SYMBOLS (punctuation, digits, only)
func NewSymbolAlphabet(name, symbols string) *Alphabet {
	return &Alphabet{name, symbols, false, !IsExtendedASCII(symbols), true, nil, false, ""}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Alphabet with its name
func (a *Alphabet) String() string {
	return fmt.Sprintf("%10s %s", a.Name, a.Chars)
}

func (a *Alphabet) IsBinary() bool {
	return a.isBinary
}

/**
 * If nil it removes the special case handlers. Else it first invokes
 * its Normalize() method to ensure non-provided handlers are mapped
 * to library functions.
 */
func (a *Alphabet) WithSpecialCase(handlers *SpecialCaseHandler) *Alphabet {
	if handlers != nil {
		handlers.Normalize()
	}
	a.specialCase = handlers
	return a
}

func (a *Alphabet) WithLangCode(iso string) *Alphabet {
	a.langCode = iso
	return a
}

func (a *Alphabet) LangCodeISO() string {
	return a.langCode
}

/**
 * Get a copy of the Special Case Handlers (if any). Changing them does NOT
 * alter the alphabet's handlers. If the return value is not nil, it
 * guarantees that all handlers are set to a custom or standard library
 * equivalent.
 */
func (a *Alphabet) BorrowSpecialCase() *SpecialCaseHandler {
	if a.specialCase == nil {
		return nil
	}

	var handlers SpecialCaseHandler = SpecialCaseHandler{
		ToUpperString: a.specialCase.ToUpperString,
		ToLowerString: a.specialCase.ToLowerString,
		ToUpperRune:   a.specialCase.ToUpperRune,
		ToLowerRune:   a.specialCase.ToLowerRune,
	}

	handlersPtr := &handlers
	handlersPtr.Normalize()
	return handlersPtr
}

/**
 * The number of runes (characters/symbols) in the alphabet.
 * NOTE: It counts runes not bytes.
 */
func (a *Alphabet) Size() uint {
	return uint(utf8.RuneCountInString(a.Chars))
}

func (a *Alphabet) ToUpper() *Alphabet {
	/*
		// Another Hard-to-find Quirk BUG! strings.ToUpper('ß')
		// Found in DuckDuck.com Search Assist:
		// "The strings.ToUpper() function in Go does not convert the German
		// character "ß" to uppercase because it traditionally does not have
		// an uppercase equivalent. However, since 2017, the uppercase "ẞ" is
		// accepted in some contexts, but it is not automatically produced by
		// the ToUpper() function.
		a.Chars = strings.ToUpper(a.Chars)

		const LowerCaseEsTzet string = "ß"
		const UpperCaseEsTzet string = "ẞ"
		if strings.Contains(a.Chars, LowerCaseEsTzet) {
			// Now handle the exception
			// In GO strings.ToUpper(LowerCaseEsTzet) == LowerCaseEsTzet !!!
			a.Chars = strings.ReplaceAll(a.Chars, LowerCaseEsTzet, UpperCaseEsTzet)
		}
	*/

	a.Chars = a.ToUpperString(a.Chars)

	return a
}

func (a *Alphabet) ToLower() *Alphabet {
	a.Chars = a.ToLowerString(a.Chars)

	return a
}

func (a *Alphabet) ToUpperString(text string) string {
	var result string
	if a.specialCase == nil {
		result = strings.ToUpper(text)
	} else {
		result = a.specialCase.ToUpperString(text)
	}

	return result
}

func (a *Alphabet) ToLowerString(text string) string {
	var result string
	if a.specialCase == nil {
		result = strings.ToLower(text)
	} else {
		result = a.specialCase.ToLowerString(text)
	}

	return result
}

/**
 * Get the (0-based) position of the letter within the alphabet.
 * It works with single and multi-byte rune alphabets.
 */
func (a *Alphabet) PositionOf(r rune) int {
	if !strings.Contains(a.Chars, string(r)) {
		return -1
	}

	// @note strings.Index() returns the BYTE position, that works ONLY
	// with single-byte rune alphabets (English). But here we deal with Foreign
	// alphabets, many of which have multi-byte runes. We need the
	// position of the letter in the alphabet, not its byte position.
	bytePos := strings.Index(a.Chars, string(r))
	charPos := utf8.RuneCountInString(a.Chars[:bytePos])

	return charPos
}

/**
 * Get the alphabet rune at the given position.
 * Panics if out of range.
 * Special case Pos=-1 which returns the last rune.
 */
func (a *Alphabet) GetRuneAt(pos int) rune { //@audit use cmn.RuneAt
	if pos < -1 || pos >= utf8.RuneCountInString(a.Chars) {
		mlog.ErrorT("Index out of Range", mlog.String("At", "GetRuneAt"), mlog.String("Alpha", a.Name))
		panic("GetRuneAt index out of range")
	}

	if pos == -1 {
		pos = int(a.Size()) - 1
	}
	return []rune(a.Chars)[pos]
}

// Rename the alphabet (convenience method)
func (a *Alphabet) Rename(name string) *Alphabet {
	a.Name = name
	return a
}

// Rename the alphabet by appending a number to the current name (convenience method)
func (a *Alphabet) Renumber(sequence uint) *Alphabet {
	a.Name = fmt.Sprintf("%s#%d", a.Name, sequence)
	return a
}

/**
 * Check if the rune r exists in the current alphabet. The
 * comparison can be case-insensitive if ignoreCase is set.
 */
func (a *Alphabet) Contains(r rune, ignoreCase bool) bool {
	for _, ch := range a.Chars {
		if ignoreCase { // Then handle in UPPERCASE
			//ch = unicode.ToUpper(ch)
			//r = unicode.ToUpper(r)
			if a.specialCase == nil {
				ch = unicode.ToUpper(ch)
				r = unicode.ToUpper(r)
			} else {
				ch = a.specialCase.ToUpperRune(ch)
				r = a.specialCase.ToUpperRune(r)
			}
		}

		if ch == r {
			return true
		}
	}

	return false
}

// Returns a NEW alphabet object based on the new
func (a *Alphabet) From(otherAlphabetChars string) *Alphabet {
	if a.OnlySymbols {
		return NewSymbolAlphabet(a.Name, otherAlphabetChars)
	} else {
		return NewAlphabet(a.Name, otherAlphabetChars, !a.Foreign, a.OnlySymbols)
	}
}

// The current alphabet but rotated left (negative) or right (positive) positions.
// Returns a new instance.
func (a *Alphabet) Rotate(rotateQty int) *Alphabet {
	var result *Alphabet = nil
	if rotateQty == 0 {
		result = a
	} else if rotateQty < 0 {
		result = a.From(RotateStringLeft(a.Chars, rotateQty*-1))
	} else if rotateQty > 0 {
		result = a.From(RotateStringRight(a.Chars, rotateQty))
	}
	return result
}

/**
 * The name must not be empty and the alphabet is not empty and
 * must has only unique runes (no duplicates).
 */
func (a *Alphabet) Check() bool {
	if a.specialCase != nil {
		a.specialCase.Normalize()
	}
	return IsNotBlank(a.Name) && IsNotBlank(a.Chars) && HasUniqueRunes(a.Chars)
}

/**
 * Clone an alphabet. Please use this to clone the built-in alphabets
 */
func (a *Alphabet) Clone() *Alphabet {
	return NewAlphabet(a.Name, a.Chars, !a.Foreign, a.OnlySymbols).WithSpecialCase(a.specialCase)
}
