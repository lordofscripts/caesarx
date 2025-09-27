/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * An Upper/Lowercase special handling interface. It allows built-in
 * languages to do special rune case conversions. At the time of this
 * development, neither English/Latin/Greek/Cyrillic need it; However,
 * the German languages NEEDS it because there are special GO rules
 * (and a GO bug!) for the ß conversion to ẞ.
 *-----------------------------------------------------------------*/
package cmn

import (
	"strings"
	"unicode"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var (
	// When there are no special case rules and we can use the
	// standard GO library functions.
	DefaultCaseHandler *SpecialCaseHandler = (&SpecialCaseHandler{}).Normalize()
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type ISpecialCaseHandlers interface {
	ToUpperString(string) string
	ToLowerString(string) string
	ToUpperRune(rune) rune
	ToLowerRune(rune) rune
	IsUpperRune(rune) bool
	IsLowerRune(rune) bool
}

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type SpecialCaseHandler struct { // @audit consider IsUpperRune? & IsLowerRune?
	ToUpperString func(string) string
	ToLowerString func(string) string
	ToUpperRune   func(rune) rune
	ToLowerRune   func(rune) rune
	IsUpperRune   func(rune) bool
	IsLowerRune   func(rune) bool
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

/**
 * Ensure than non-provided case handlers are mapped to the
 * standard GO library equivalents.
 */
func (s *SpecialCaseHandler) Normalize() *SpecialCaseHandler {
	if s.ToLowerRune == nil {
		s.ToLowerRune = unicode.ToLower
	}
	if s.ToUpperRune == nil {
		s.ToUpperRune = unicode.ToUpper
	}
	if s.ToLowerString == nil {
		s.ToLowerString = strings.ToLower
	}
	if s.ToUpperString == nil {
		s.ToUpperString = strings.ToUpper
	}
	if s.IsLowerRune == nil {
		s.IsLowerRune = unicode.IsLower
	}
	if s.IsUpperRune == nil {
		s.IsUpperRune = unicode.IsUpper
	}

	return s
}

/**
 * Ensure all handlers point to a custom function or to a standard library
 * function. If anything is missing Then die (panic)
 */
func (s *SpecialCaseHandler) Assert() {
	if s.ToLowerRune == nil {
		panic("ToLowerRune not set")
	}
	if s.ToUpperRune == nil {
		panic("ToUpperRune not set")
	}
	if s.ToLowerString == nil {
		panic("ToLowerString not set")
	}
	if s.ToUpperString == nil {
		panic("ToUpperString not set")
	}
	if s.IsLowerRune == nil {
		panic("IsLowerRune not set")
	}
	if s.IsUpperRune == nil {
		panic("IsUpperrRune not set")
	}
}

func (s *SpecialCaseHandler) EqualRuneFold(a, b rune) bool {
	aU := s.ToUpperRune(a)
	bU := s.ToUpperRune(b)

	return aU == bU
}
