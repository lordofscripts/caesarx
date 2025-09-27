/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Subject: ciphers.TabulaRecta{}
 * API: EncodeRune, DecodeRune, HasRune
 * Languages: English, Spanish, German, Greek & Cyrillic.
 * Updated: 31 Aug 2025
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * TabulaRecta is the basis for multiple ciphers. In particular the
 * Caesar cipher.
 * We ensure this building block works. We especially take care that
 * it works well using multi-byte UTF8 alphabets. Particularly, we
 * take into account the strange exception of the German 'ß' letter.
 *	In GO it is the same in upper/lowercase but utf8.IsLower()
 * reports it as lower and there is no utf8.IsUpper(). However, in at
 * least many other programming languages 'ß' is lowercase and 'SS'
 * its uppercase counterpart.
 *-----------------------------------------------------------------*/
package tests

import (
	"fmt"
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/cmn"
	"testing"
)

// **** Test TabulaRecta in plain Caesar mode  ****
// A simple monosylabic translation where there is
// only ONE key (a rune) used for all text.

/**
 * Test TabulaRecta.EncodeRune() with all built-in alphabets.
 */
func Test_TabulaRecta_EncodeRune(t *testing.T) {
	type Vector struct {
		alpha  *cmn.Alphabet
		key    rune
		input  rune
		expect rune
	}

	allCases := []Vector{
		// English - Plain ASCII
		{cmn.ALPHA_DISK, 'A', 'P', 'P'},
		{cmn.ALPHA_DISK, 'M', 'P', 'B'},
		{cmn.ALPHA_DISK, 'Z', 'P', 'O'},
		// Spanish - Extended ASCII
		{cmn.ALPHA_DISK_LATIN, 'A', 'P', 'P'},
		{cmn.ALPHA_DISK_LATIN, 'M', 'P', 'É'},
		{cmn.ALPHA_DISK_LATIN, 'Z', 'P', 'J'},
		{cmn.ALPHA_DISK_LATIN, 'Ü', 'P', 'O'},
		// German - UTF8
		{cmn.ALPHA_DISK_GERMAN, 'A', 'P', 'P'},
		{cmn.ALPHA_DISK_GERMAN, 'M', 'P', 'Ö'},
		{cmn.ALPHA_DISK_GERMAN, 'Z', 'P', 'K'},
		{cmn.ALPHA_DISK_GERMAN, 'Ü', 'P', 'N'},
		{cmn.ALPHA_DISK_GERMAN, cmn.CHAR_UPPER_ESZET, 'P', 'O'}, // Beware: uppercase ẞ !!!
		// Greek - UTF8
		{cmn.ALPHA_DISK_GREEK, 'Α', 'Ξ', 'Ξ'},
		{cmn.ALPHA_DISK_GREEK, 'Λ', 'Ξ', 'Ω'},
		{cmn.ALPHA_DISK_GREEK, 'Ω', 'Ξ', 'Ν'},
		// Cyrillic - UTF8
		{cmn.ALPHA_DISK_CYRILLIC, 'А', 'Ж', 'Ж'},
		{cmn.ALPHA_DISK_CYRILLIC, 'Л', 'Ж', 'Т'},
		{cmn.ALPHA_DISK_CYRILLIC, 'Я', 'Ж', 'Ë'},
	}

	for i, v := range allCases {
		tr := ciphers.NewTabulaRecta(v.alpha, cmn.CaseInsensitive)
		got := tr.EncodeRune(v.input, v.key)
		if got != v.expect {
			t.Errorf("#%d %s ENC failed %c != %c", i+1, v.alpha.Name, got, v.expect)
		}
	}
}

func Test_TabulaRecta_DecodeRune(t *testing.T) {
	type Vector struct {
		alpha  *cmn.Alphabet
		key    rune
		input  rune
		expect rune
	}

	allCases := []Vector{
		// English - Plain ASCII
		{cmn.ALPHA_DISK, 'A', 'P', 'P'},
		{cmn.ALPHA_DISK, 'M', 'B', 'P'},
		{cmn.ALPHA_DISK, 'Z', 'O', 'P'},
		// Spanish - Extended ASCII
		{cmn.ALPHA_DISK_LATIN, 'A', 'P', 'P'},
		{cmn.ALPHA_DISK_LATIN, 'M', 'É', 'P'},
		{cmn.ALPHA_DISK_LATIN, 'Z', 'J', 'P'},
		{cmn.ALPHA_DISK_LATIN, 'Ü', 'O', 'P'},
		// German - UTF8
		{cmn.ALPHA_DISK_GERMAN, 'A', 'P', 'P'},
		{cmn.ALPHA_DISK_GERMAN, 'M', 'V', 'J'},
		{cmn.ALPHA_DISK_GERMAN, 'Z', 'K', 'P'},
		{cmn.ALPHA_DISK_GERMAN, 'Ü', 'M', 'O'},
		{cmn.ALPHA_DISK_GERMAN, cmn.CHAR_UPPER_ESZET, 'O', 'P'}, // 'ẞ'
		// Greek - UTF8
		{cmn.ALPHA_DISK_GREEK, 'Α', 'Ξ', 'Ξ'},
		{cmn.ALPHA_DISK_GREEK, 'Λ', 'Ω', 'Ξ'},
		{cmn.ALPHA_DISK_GREEK, 'Ω', 'Ν', 'Ξ'},
		// Cyrillic - UTF8
		{cmn.ALPHA_DISK_CYRILLIC, 'А', 'Ж', 'Ж'},
		{cmn.ALPHA_DISK_CYRILLIC, 'П', 'Я', 'П'},
		{cmn.ALPHA_DISK_CYRILLIC, 'Я', 'Ц', 'Ч'},
	}

	for i, v := range allCases {
		tr := ciphers.NewTabulaRecta(v.alpha, cmn.CaseInsensitive)
		got := tr.DecodeRune(v.input, v.key)
		if got != v.expect {
			t.Errorf("#%d %s DEC failed %c != %c for key %c", i+1, v.alpha.Name, got, v.expect, v.key)
		}
	}
}

/**
 * Does a round-trip EncodeRune+DecodeRune on the entire table
 * to detect possible anomalies.
 */
func Test_RoundTripAll(t *testing.T) {
	alphabets := []*cmn.Alphabet{
		cmn.ALPHA_DISK,
		cmn.NUMBERS_DISK,
		cmn.ALPHA_DISK_LATIN,
		cmn.ALPHA_DISK_GREEK,
		cmn.ALPHA_DISK_CYRILLIC,
		cmn.ALPHA_DISK_GERMAN,
	}

	for _, alpha := range alphabets {
		fmt.Print("\t·", alpha.Name, " ")
		tr := ciphers.NewTabulaRecta(alpha, cmn.CaseInsensitive)
		allGood := true
		for _, key := range alpha.Chars {
			for _, char := range alpha.Chars {
				encR := tr.EncodeRune(char, key)
				decR := tr.DecodeRune(encR, key)
				if decR != char {
					allGood = false
					t.Errorf("%s R/T failed for Key: %c Char: %c Got: %c != %c", alpha.Name, key, char, decR, char)
				}
			}
		}
		if allGood {
			fmt.Println("OK")
		}
	}
}

// Verify the existance and location (shift offset) of a rune
// in the TabulaRecta alphabet. We try both existing and non-existing.
// We also ensure to test those beyond the position of a multi-byte
// rune, thus to avoid the UTF8 pitfalls.
func Test_HasRune(t *testing.T) {
	type Vector struct {
		alpha  *cmn.Alphabet
		input  rune
		exists bool
		where  int
	}

	// Ensure we test runes that are AFTER multi-byte runes
	allCases := []Vector{
		{cmn.ALPHA_DISK, 'Y', true, 24},
		{cmn.ALPHA_DISK, 'Ñ', false, -1},
		{cmn.ALPHA_DISK_LATIN, 'Ú', true, 31},
		{cmn.ALPHA_DISK_LATIN, 'ß', false, -1},
		{cmn.ALPHA_DISK_GREEK, 'Π', true, 15},
		{cmn.ALPHA_DISK_GREEK, 'Ñ', false, -1},
		{cmn.ALPHA_DISK_GERMAN, 'Ö', true, 27},
		{cmn.ALPHA_DISK_GERMAN, 'Ñ', false, -1},
		{cmn.ALPHA_DISK_CYRILLIC, 'Ш', true, 25},
		{cmn.ALPHA_DISK_CYRILLIC, 'ß', false, -1},
	}

	for i, tc := range allCases {
		tr := ciphers.NewTabulaRecta(tc.alpha, cmn.CaseInsensitive)
		exists, where := tr.HasRune(tc.input)
		if exists != tc.exists {
			t.Errorf("#%d %s Char: %c Exists: %t Got: %t", i+1, tc.alpha.Name, tc.input, tc.exists, exists)
		}
		if where != tc.where {
			t.Errorf("#%d %s Char: %c Where: %d Got: %d", i+1, tc.alpha.Name, tc.input, tc.where, where)
		}
	}
}
