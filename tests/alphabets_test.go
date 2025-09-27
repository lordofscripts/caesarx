/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Subject: caesarx.Alphabets{}
 * API: GetRune, PositionOf, Size, Contains, ToLower, ToUpper, Rotate
 * Languages: English, Spanish, German, Greek & Cyrillic.
 * Updated: 31 Aug 2025
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Alphabets{} is an integral object in this cipher module because
 * it lets us handle the built-in and user-defined alphabets.
 *-----------------------------------------------------------------*/
package tests

import (
	"fmt"
	"lordofscripts/caesarx/cmn"
	"testing"
	"unicode/utf8"
)

func Test_GetRune(t *testing.T) {
	type Vector struct {
		alpha    *cmn.Alphabet
		first    rune
		middle   rune
		last     rune
		middleAt int
	}

	allCases := []Vector{ // UPPERCASES!!!
		{cmn.ALPHA_DISK, 'A', 'M', 'Z', 12},
		{cmn.ALPHA_DISK_LATIN, 'A', 'Ñ', 'Ü', 14},
		{cmn.ALPHA_DISK_GREEK, 'Α', 'Λ', 'Ω', 10},
		{cmn.ALPHA_DISK_GERMAN, 'A', 'O', 'ẞ', 14},
		{cmn.ALPHA_DISK_CYRILLIC, 'А', 'П', 'Я', 16},
	}

	for i, tc := range allCases {
		fmt.Print("\t· ", tc.alpha.Name, " ")
		allGood := true

		got := tc.alpha.GetRuneAt(0)
		if got != tc.first {
			t.Errorf("#%d %s 1st rune expected: %c got: %c", i+1, tc.alpha.Name, tc.first, got)
			allGood = false
		}

		got = tc.alpha.GetRuneAt(tc.middleAt)
		if got != tc.middle {
			t.Errorf("#%d %s middle rune expected: %c got: %c", i+1, tc.alpha.Name, tc.middle, got)
			allGood = false
		}

		lastRuneNr := tc.alpha.Size() - 1
		got = tc.alpha.GetRuneAt(int(lastRuneNr))
		if got != tc.last {
			t.Errorf("#%d %s last rune expected: %c got: %c", i+1, tc.alpha.Name, tc.last, got)
			allGood = false
		}

		if allGood {
			fmt.Println("OK")
		}
	}
}

func Test_PositionOf(t *testing.T) {
	type Vector struct {
		alpha    *cmn.Alphabet
		middle   rune
		middleAt int
	}

	allCases := []Vector{
		{cmn.ALPHA_DISK, 'M', 12},
		{cmn.ALPHA_DISK_LATIN, 'Ñ', 14},
		{cmn.ALPHA_DISK_GREEK, 'Λ', 10},
		{cmn.ALPHA_DISK_GERMAN, 'O', 14},
		{cmn.ALPHA_DISK_CYRILLIC, 'П', 16},
	}

	for i, tc := range allCases {
		fmt.Print("\t· ", tc.alpha.Name, " ")

		got := tc.alpha.PositionOf(tc.middle)
		if got != tc.middleAt {
			t.Errorf("#%d %s mid rune %c expected at: %d got: %d", i+1, tc.alpha.Name, tc.middle, tc.middleAt, got)
		}
	}
}

func Test_Size(t *testing.T) {
	const MIXED = "ABCDÁÉÍÓЖЗИЙКЛΔΕΖΗΘÄÖÜẞ"
	alpha := cmn.NewAlphabet("Mixed", MIXED, false, false)

	if alpha.Size() != uint(utf8.RuneCountInString(MIXED)) {
		t.Errorf("wrong size of mixed alphabet")
	}
}

func Test_Contains(t *testing.T) {
	const MIXED1 = "abcdÁÉÍÓЖЗИЙКЛΔΕΖΗΘÄÖÜẞ"
	const MIXED2 = "ABCDÁÉÍÓЖЗИЙКЛΔΕΖΗΘÄÖÜẞ"

	alpha := cmn.NewAlphabet("Mixed LC", MIXED1, false, false)
	r := 'A'
	if alpha.Contains(r, false) != false {
		t.Errorf("(case sensitive) uppercase '%c' isn't in %s", r, MIXED1)
	}

	alpha = cmn.NewAlphabet("Mixed UC", MIXED2, false, false)
	r = 'a'
	if alpha.Contains(r, true) != true {
		t.Errorf("(case insensitive) rune '%c' isn't in %s", r, MIXED2)
	}
}

func Test_AlphabetToLower(t *testing.T) {
	const ALPHA_U = "ABCDEFGHIJKБΘΛΜẞÜ"
	const ALPHA_L = "abcdefghijkбθλμßü"

	// NOTE: The ẞ to ß ToLower works with standard library
	alpha := cmn.NewAlphabet("Upper", ALPHA_U, false, false)
	got := alpha.ToLower().ToLower().Chars
	if got != ALPHA_L {
		t.Errorf("improper lowercase %s != %s", got, ALPHA_L)
	}
}

func Test_AlphabetToUpper(t *testing.T) {
	const ALPHA_U = "ABCDEFGHIJKБΘΛΜẞÜ"
	const ALPHA_L = "abcdefghijkбθλμßü"

	// NOTE: There is a BUG in GO v1.25 where neither strings.ToUpper()
	// nor unicode.ToUpper() work reliably if ß (lowercase) is present!
	// Therefore, I use GermanSpecialCase
	alpha := cmn.NewAlphabet("Lower Has German", ALPHA_L, false, false)
	alpha.WithSpecialCase(cmn.GermanSpecialCase)
	// @note A special exception is made for GERMAN because strings.ToUpper()
	// does NOT convert 'ß' (lower) to 'ẞ' (upper). However, this quirky
	// exception is handled in alphabet.ToUpper(). KEEP in mind that in other
	// programming languages 'SS' becomes the uppercase!
	got := alpha.ToLower().ToUpper().Chars
	if got != ALPHA_U {
		t.Errorf("improper lowercase\n\tGot: %s\n\tExp: %s", got, ALPHA_U)
	}
}

func Test_Rotate(t *testing.T) {
	// Normal        : АБВГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ
	// Rotate Right 3: ЭЮЯАБВГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬ
	const RIGHT_3 = +3
	alphaRight3 := cmn.ALPHA_DISK_CYRILLIC.Rotate(RIGHT_3)
	if alphaRight3.GetRuneAt(0) != 'Э' {
		fmt.Println("Rotated Right: ", alphaRight3.Chars)
		t.Errorf("unexpected 1st rune after right rotation")
	}
	if alphaRight3.Chars != "ЭЮЯАБВГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬ" {
		fmt.Println("Rotated Right: ", alphaRight3.Chars)
		t.Errorf("unexpected alphabet after right rotation")
	}

	// Rotate Left 8 : ЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯАБВГДЕËЖ
	const LEFT_8 = -8
	alphaLeft8 := cmn.ALPHA_DISK_CYRILLIC.Rotate(LEFT_8)
	if alphaLeft8.GetRuneAt(8) != 'П' {
		fmt.Println("Rotated Left: ", alphaLeft8.Chars)
		t.Errorf("unexpected 8th rune after left rotation")
	}
	if alphaLeft8.Chars != "ЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯАБВГДЕËЖ" {
		fmt.Println("Rotated Left: ", alphaLeft8.Chars)
		t.Errorf("unexpected alphabet after left rotation")
	}
}

/**
 * This should have passed but exposed a BUG in GO v1.25
 * It was submitted to the Go Developers.
 */
// /*
func Test_GermanCharacter(t *testing.T) {
	const LOWER_RUNE rune = 'ß' // 223
	const UPPER_RUNE rune = 'ẞ' // 7338
	const LOWER_STRING string = "daß"
	const UPPER_STRING string = "DAẞ"

	// just to make sure, these two PASS
	if LOWER_RUNE == UPPER_RUNE { // passes
		t.Errorf("This isn't supposed to happen in GO v1.25")
	}
	if LOWER_STRING == UPPER_STRING { // passes
		t.Errorf("This isn't supposed to happen in GO v1.25")
	}

	fmt.Printf("Lowercase es-tzet: %c (%U)\n", LOWER_RUNE, LOWER_RUNE)
	fmt.Printf("Uppercase es-tzet: %c (%U)\n", UPPER_RUNE, UPPER_RUNE)

	// Now the STD Lib inconsistencies
	const TRY_SPECIAL = true // @audit wait if GO bug becomes fixed in GO v1.26
	var gotS, expectS string
	var special cmn.SpecialCaseHandler
	if TRY_SPECIAL {
		special = *cmn.GermanSpecialCase // Should work
	} else {
		special = cmn.NoSpecialCase // Would fail due to bug in GO
	}

	// (a) strings.ToUpper("daß") should be "DAẞ"
	expectS = UPPER_STRING
	gotS = special.ToUpperString(LOWER_STRING)
	if gotS != expectS { // FAILS with Std. Lib. v1.25
		fmt.Println("strings.ToUpper() FAIL")
		t.Errorf("strings.ToUpper failure Got:'%s' Expect:'%s'", gotS, expectS)
	} else {
		fmt.Println("strings.ToUpper() OK")
	}

	// (b) strings.ToLower("DAẞ") should be "daß"
	expectS = LOWER_STRING
	gotS = special.ToLowerString(UPPER_STRING)
	if gotS != expectS {
		fmt.Println("strings.ToLower() FAIL")
		t.Errorf("strings.ToLower failure Got:'%s' Expect:'%s'", gotS, expectS)
	} else {
		fmt.Println("strings.ToLower() OK")
	}

	var gotR, expectR rune

	// (c) unicode.ToUpper('ß') should be 'ẞ'
	// Here it works fine, but in a small program it fails to convert.
	expectR = UPPER_RUNE
	gotR = special.ToUpperRune(LOWER_RUNE)
	if gotS != expectS {
		fmt.Println("unicode.ToUpper() FAIL")
		t.Errorf("unicode.ToUpper failure Got:'%c' Expect:'%c'", gotR, expectR)
	} else {
		fmt.Println("unicode.ToUpper() OK")
	}

	// (d) unicode.ToLower('ẞ') should be 'ß'
	expectR = LOWER_RUNE
	gotR = special.ToLowerRune(UPPER_RUNE)
	if gotS != expectS {
		fmt.Println("unicode.ToLower() FAIL")
		t.Errorf("unicode.ToLower failure Got:'%c' Expect:'%c'", gotR, expectR)
	} else {
		fmt.Println("unicode.ToLower() OK")
	}
}

// */
