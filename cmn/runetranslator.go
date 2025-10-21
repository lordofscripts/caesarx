/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * RuneTranslator transliterates an alphabet into another alphabet.
 * It does forward and reverse lookups and works fine with Unicode.
 * The index lookup in string objects was too unreliable, even with
 * utf8.RunesInString().
 *  If the rune is not found in the map, the transliterator returns
 * the same rune that was given as input (no translation).
 *-----------------------------------------------------------------*/
package cmn

import (
	"fmt"
	"unicode/utf8"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type RuneTranslator struct {
	Title       string
	table       map[rune]rune
	source      string
	target      string
	specialCase *SpecialCaseHandler
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// A case-insensitive Rune Translator that takes into consideration
// any special casing rules of the alphabet. Transliteration does
// preserve case though...
func NewRuneTranslator(name string, alphabet, altAlphabet *Alphabet) *RuneTranslator {
	if alphabet.Size() != altAlphabet.Size() {
		panic("Length of alphabets is not the same!")
	}

	if alphabet.specialCase != altAlphabet.specialCase {
		panic("Alphabets must have the same SpecialCase handlers") //@audit store both handlers (next version)
	}

	runesA := []rune(alphabet.ToUpperString(alphabet.Chars))
	runesB := []rune(altAlphabet.ToUpperString(altAlphabet.Chars))
	table := make(map[rune]rune, len(runesA))

	for index, key := range runesA {
		table[key] = runesB[index]
	}

	return &RuneTranslator{name, table, alphabet.Chars, altAlphabet.Chars, alphabet.specialCase}
}

// A case-insensitive Rune Translator that takes into consideration
// any special casing rules of the alphabet. Transliteration does
// preserve case though...
func NewSimpleRuneTranslator(name string, alphabet, altAlphabet string, caser *SpecialCaseHandler) *RuneTranslator {
	if utf8.RuneCountInString(alphabet) != utf8.RuneCountInString(altAlphabet) {
		panic("Length of alphabets is not the same!")
	}

	runesA := []rune(caser.ToUpperString(alphabet))
	runesB := []rune(caser.ToUpperString(altAlphabet))
	table := make(map[rune]rune, len(runesA))

	for index, key := range runesA {
		table[key] = runesB[index]
	}

	return &RuneTranslator{name, table, alphabet, altAlphabet, caser}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (r *RuneTranslator) String() string {
	return fmt.Sprintf("%s Qty:%d %s", r.Title, len(r.table), r.target)
}

func (r *RuneTranslator) GetSource() string {
	return r.source
}

func (r *RuneTranslator) GetTarget() string {
	return r.target
}

// Check whether this rune exists in the catalog. Case-insensitive.
func (r *RuneTranslator) Exists(char rune) bool {
	var charCF rune = char
	if r.specialCase != nil && r.specialCase.ToUpperRune != nil {
		charCF = r.specialCase.ToUpperRune(char)
	}

	_, ok := r.table[charCF]
	return ok
}

func (r *RuneTranslator) IsValid() bool {
	size := len(r.source)
	return size > 0 && size == len(r.target) && len(r.table) == size
}

// Given a rune, it returns the transliterated rune. If it is
// not found, it returns the input rune untraslated with error.
// Case-insensitive lookup.
func (r *RuneTranslator) Lookup(char rune) (rune, error) {
	var mkey rune = char
	isLower := r.specialCase.IsLowerRune(char)
	if isLower {
		mkey = r.specialCase.ToUpperRune(char)
	}

	if val, ok := r.table[mkey]; ok { // case-insensitive lookup
		if isLower { // the map contains only Uppercases
			return r.specialCase.ToLowerRune(val), nil
		}
		return val, nil
	}

	return char, fmt.Errorf("couldn't find %s:%c during forward lookup", r.Title, char)
}

// Given a transliterated rune, it returns the original rune. If it is
// not found, it returns the input rune untraslated with error.
// Case-insensitive lookup (preserves case)
func (r *RuneTranslator) ReverseLookup(char rune) (rune, error) {
	var mChar rune = char
	isLower := r.specialCase.IsLowerRune(char)
	if isLower {
		mChar = r.specialCase.ToUpperRune(char)
	}

	for key, value := range r.table {
		if value == mChar {
			if isLower {
				return r.specialCase.ToLowerRune(key), nil
			}
			return key, nil
		}
	}

	return char, fmt.Errorf("couldn't find %s:%c during reverse lookup", r.Title, char)
}

// The original (not transliterated) string
func (r *RuneTranslator) SourceString() string { //@audit deprecate
	return r.source
}

// The transliterated string
func (r *RuneTranslator) TransliteratedString() string { // @audit deprecate
	return r.target
}

// Source & Transliteratd strings separated by newline
func (r *RuneTranslator) TapeString() string {
	return fmt.Sprintf("%s\n%s\n", r.source, r.target)
}
