/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package tests

import (
	"fmt"
	"lordofscripts/caesarx/ciphers/affine"
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/cmn"
	"slices"
	"testing"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

var (
	//  English N=26 [1 3 5 7 9 11 15 17 19 21 23 25]
	//  Spanish N=33 [1 2 4 5 7 8 10 13 14 16 17 19 20 23 25 26 28 29 31 32]
	//   German N=30 [1 7 11 13 17 19 23 29]
	//    Greek N=24 [1 5 7 11 13 17 19 23]
	// Cyrillic N=33 [1 2 4 5 7 8 10 13 14 16 17 19 20 23 25 26 28 29 31 32]
	validAffineParamsEN = &affine.AffineParams{A: 5, B: 3, Ap: 21, N: 26}
	validAffineParamsES = &affine.AffineParams{A: 7, B: 3, Ap: 19, N: 33}
	validAffineParamsDE = &affine.AffineParams{A: 7, B: 3, Ap: 13, N: 30}
	validAffineParamsGR = &affine.AffineParams{A: 7, B: 3, Ap: 7, N: 24}
	validAffineParamsRU = &affine.AffineParams{A: 7, B: 3, Ap: 19, N: 33}
)

/* ----------------------------------------------------------------
 *					T e s t s :: AffineHelper
 *-----------------------------------------------------------------*/

// Affine Parameter A should be Coprime of alphabet length N
func Test_AreCoprime(t *testing.T) {
	allCases := []struct {
		A          int
		N          int
		areCoprime bool
	}{
		{5, 26, true},
		{1, 26, true},
		{11, 26, true},
		{25, 26, true},
		{2, 26, false},
		{10, 26, false},
	}

	h := affine.NewAffineHelper()
	for vnum, tc := range allCases {
		got := h.AreCoprime(tc.A, tc.N)
		if got != tc.areCoprime {
			t.Errorf("#%d A:%d and %d should be coprimes: %t but got %t", vnum+1, tc.A, tc.N, tc.areCoprime, got)
		}
	}
}

func Test_ValidCoprimesUpTo(t *testing.T) {
	coprimes26 := []int{1, 3, 5, 7, 9, 11, 15, 17, 19, 21, 23, 25}

	h := affine.NewAffineHelper()
	if slices.Compare(h.ValidCoprimesUpTo(26), coprimes26) != 0 {
		t.Error("coprimes for N=26 failed")
	}
}

func Test_ModularInverse(t *testing.T) {
	const N int = 26
	allCases := []struct {
		A  int
		Ap int
	}{
		{1, 1},
		{3, 9},
		{5, 21},
		{7, 15},
		{9, 3},
		{11, 19},
	}

	h := affine.NewAffineHelper()
	for vnum, tc := range allCases {
		inverse, err := h.ModularInverse(tc.A, N)
		if err != nil {
			t.Error(err)
		}
		if tc.Ap != inverse {
			t.Errorf("#%d A=%d A'=%d should be modular inverses, but got %d", vnum+1, tc.A, tc.Ap, inverse)
		}
	}
}

func Test_SetParams(t *testing.T) {
	const N1 int = 26
	allCases := []struct {
		A, B, N int
		valid   bool
	}{
		{3, 5, N1, true},
		{19, 18, N1, true},
		{2, 5, N1, false},   // 2 is not a coprime of 26
		{-1, 5, N1, false},  // A cannot be zero
		{11, -1, N1, false}, // B should not be negative
		{19, 8, 0, false},   // N should be positive
	}

	h := affine.NewAffineHelper()
	for vnum, tc := range allCases {
		if err := h.SetParameters(tc.A, tc.B, tc.N); (err == nil) != tc.valid {
			t.Errorf("#%d SetParameters should be valid: %t but got opposite", vnum+1, tc.valid)
		}
	}
}

func Test_GetParams(t *testing.T) {
	h := affine.NewAffineHelper()
	err := h.SetParameters(validAffineParamsEN.A, validAffineParamsEN.B, validAffineParamsEN.N)
	if err != nil {
		t.Errorf("SetParameters should have been error-free:\n\tgot: %v", err)
	}

	pars := h.GetParams()
	if pars.A != validAffineParamsEN.A {
		t.Errorf("retrived param A not the same exp:%d got:%d", validAffineParamsEN.A, pars.A)
	}
	if pars.B != validAffineParamsEN.B {
		t.Errorf("retrived param B not the same exp:%d got:%d", validAffineParamsEN.B, pars.B)
	}
	if pars.Ap != validAffineParamsEN.Ap {
		t.Errorf("retrived param A' not the same exp:%d got:%d", validAffineParamsEN.Ap, pars.Ap)
	}
	if pars.N != validAffineParamsEN.N {
		t.Errorf("retrived param N not the same exp:%d got:%d", validAffineParamsEN.N, pars.N)
	}
}

func Test_VerifyParams(t *testing.T) {
	t.Skip()
	h := affine.NewAffineHelper()
	p1 := &affine.AffineParams{ // let's have Ap be fixed/recalculated
		A:  5,
		Ap: -1,
		B:  3,
		N:  26,
	}

	fmt.Printf("Before call, address of params: %p\n", p1)
	err := h.VerifyParams(p1)
	fmt.Printf("After call, address of params: %p\n", p1)
	if err != nil {
		t.Error("should not have returned error")
	}
	if _, ap, _, _ := h.GetParameters(); ap != validAffineParamsEN.Ap {
		t.Errorf("calculated parameter A' did not get fixed. exp:%d got:%d", validAffineParamsEN.Ap, ap)
	}

	err = h.VerifyParams(validAffineParamsEN)
	if err != nil {
		t.Error("unexpected error for perfect parameters")
	}
}

func Test_Encode_Helper(t *testing.T) {
	allCases := []struct {
		In, Out int
	}{
		{3, 18},  // D,S
		{14, 21}, // O,V
	}

	h := affine.NewAffineHelper()
	if err := h.SetParams(validAffineParamsEN); err != nil {
		t.Error(err)
	}

	for vnum, tc := range allCases {
		if out, err := h.Encode(tc.In); err != nil {
			t.Errorf("#%d encode error: %v", vnum+1, err)
		} else if out != tc.Out {
			t.Errorf("#%d encode failed, exp:%d got:%d", vnum+1, tc.Out, out)
		}
	}
}

func Test_Decode_Helper(t *testing.T) {
	allCases := []struct {
		In, Out int
	}{
		{18, 3},  // S,D
		{21, 14}, // V,O
	}

	h := affine.NewAffineHelper()
	if err := h.SetParams(validAffineParamsEN); err != nil {
		t.Error(err)
	}

	for vnum, tc := range allCases {
		if out, err := h.Decode(tc.In); err != nil {
			t.Errorf("#%d decode error: %v", vnum+1, err)
		} else if out != tc.Out {
			t.Errorf("#%d decode failed, exp:%d got:%d", vnum+1, tc.Out, out)
		}
	}
}

func Test_EncodeRuneFrom(t *testing.T) {
	allCases := []struct {
		In, Out rune
		Alpha   *cmn.Alphabet
	}{ // SELENIA -> PXGXRD for A=5, B=3, N=26
		{'S', 'P', cmn.ALPHA_DISK},
		{'L', 'G', cmn.ALPHA_DISK},
		{'N', 'Q', cmn.ALPHA_DISK},
		{'D', 'S', cmn.ALPHA_DISK}, // D,S
		{'O', 'V', cmn.ALPHA_DISK}, // O,V
	}

	h := affine.NewAffineHelper()
	if err := h.SetParams(validAffineParamsEN); err != nil {
		t.Error(err)
	}

	for vnum, tc := range allCases {
		if out, err := h.EncodeRuneFrom(tc.In, tc.Alpha.Chars); err != nil {
			t.Errorf("#%d encode error: %v", vnum+1, err)
		} else if out != tc.Out {
			t.Errorf("#%d encode failed, exp:'%c' got:'%c'", vnum+1, tc.Out, out)
		}
	}
}

func Test_DecodeRuneFrom(t *testing.T) {
	allCases := []struct {
		In, Out rune
		Alpha   *cmn.Alphabet
	}{
		{'P', 'S', cmn.ALPHA_DISK},
		{'G', 'L', cmn.ALPHA_DISK},
		{'Q', 'N', cmn.ALPHA_DISK},
		{'S', 'D', cmn.ALPHA_DISK},
		{'O', 'X', cmn.ALPHA_DISK},
	}

	h := affine.NewAffineHelper()
	if err := h.SetParams(validAffineParamsEN); err != nil {
		t.Error(err)
	}

	for vnum, tc := range allCases {
		if out, err := h.DecodeRuneFrom(tc.In, tc.Alpha.Chars); err != nil {
			t.Errorf("#%d decode error: %v", vnum+1, err)
		} else if out != tc.Out {
			t.Errorf("#%d decode failed, exp:'%c' got:'%c'", vnum+1, tc.Out, out)
		}
	}
}

// Not a test, simply a helper to list the Coprime list
func Test_Coprimes(t *testing.T) {
	t.Helper()
	alphas := []*cmn.Alphabet{
		cmn.ALPHA_DISK,          //  English N=26 [1 3 5 7 9 11 15 17 19 21 23 25]
		cmn.ALPHA_DISK_LATIN,    //  Spanish N=33 [1 2 4 5 7 8 10 13 14 16 17 19 20 23 25 26 28 29 31 32]
		cmn.ALPHA_DISK_GERMAN,   //   German N=30 [1 7 11 13 17 19 23 29]
		cmn.ALPHA_DISK_GREEK,    //    Greek N=24 [1 5 7 11 13 17 19 23]
		cmn.ALPHA_DISK_CYRILLIC, // Cyrillic N=33 [1 2 4 5 7 8 10 13 14 16 17 19 20 23 25 26 28 29 31 32]
	}

	h := affine.NewAffineHelper()
	for _, alpha := range alphas {
		n := alpha.Size()
		coprimes := h.ValidCoprimesUpTo(n)
		fmt.Printf("%20s N=%2d %v\n", alpha.Name, n, coprimes)
	}
}

/* ----------------------------------------------------------------
 *					T e s t s :: AffineParams
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *					T e s t s :: AffineEncoder
 *-----------------------------------------------------------------*/

// Subject: affine.AffineEncoder.Encode()
func Test_Encoder(t *testing.T) {
	allCases := []struct {
		Alpha  *cmn.Alphabet
		Params *affine.AffineParams
		In     string
		Out    string
	}{
		// ALP: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		// CNV: "DINSXCHMRWBGLQVAFKPUZEJOTY"
		{cmn.ALPHA_DISK, validAffineParamsEN, "I love cryptography", "R gvex nktauvhkdamt"},
		// ALP: "ABCDEFGHIJKLMNÑOPQRSTUVWXYZÁÉÍÓÚÜ"
		// CNV: "DKQXÚFMSZAHÑUÉCJPWÓELRYÜGNTÁBIOVÍ"
		{cmn.ALPHA_DISK_LATIN, validAffineParamsES, "Amo la criptografía", "Duj ñd qózpljmódfid"},
		// ALP: "ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜẞ"
		// CNV: "DKRYBIPWẞGNUÖELSZCJQXAHOVÜFMTÄ"
		{cmn.ALPHA_DISK_GERMAN, validAffineParamsDE, "Daß liebe hübschen Mädschen", "Ydä ußbkb wtkjrwbe Öfyjrwbe"},
		// ALP: "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ"
		// CNV: "ΔΛΣΑΘΟΧΕΜΤΒΙΠΨΖΝΥΓΚΡΩΗΞΦ"
		{cmn.ALPHA_DISK_GREEK, validAffineParamsGR, "Λατρεύω την κρυπτογραφία", "Βδκυθύφ κχπ τυρνκζσυδωίδ"},
		// ALP: "АБВГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
		// CNV: "ГЙРЧЮЕЛТЩАЖНФЫВИПЦЭДКСШЯËМУЪБЗОХЬ"
		{cmn.ALPHA_DISK_CYRILLIC, validAffineParamsRU, "Мы любим криптографию", "Ыб фхйаы нцапдичцгсах"},
	}

	for vnum, tc := range allCases {
		enc := affine.NewAffineEncoder(tc.Alpha, tc.Params)
		if out, err := enc.Encode(tc.In); err != nil {
			t.Errorf("#%d %s encoder error: %v", vnum+1, tc.Alpha.LangCodeISO(), err)
		} else if out != tc.Out {
			t.Errorf("#%d %s encoder failed\n\tin :'%s'\n\texp:'%s'\n\tgot:'%s'", vnum+1, tc.Alpha.LangCodeISO(), tc.In, tc.Out, out)
		}
	}
}

/* ----------------------------------------------------------------
 *					T e s t s :: AffineDecoder
 *-----------------------------------------------------------------*/

// Subject: affine.AffineDecoder.Decode()
func Test_Decoder(t *testing.T) {
	allCases := []struct {
		Alpha  *cmn.Alphabet
		Params *affine.AffineParams
		Out    string
		In     string
	}{
		// ALP: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		// CNV: "DINSXCHMRWBGLQVAFKPUZEJOTY"
		{cmn.ALPHA_DISK, validAffineParamsEN, "I love cryptography", "R gvex nktauvhkdamt"},
		// ALP: "ABCDEFGHIJKLMNÑOPQRSTUVWXYZÁÉÍÓÚÜ"
		// CNV: "DKQXÚFMSZAHÑUÉCJPWÓELRYÜGNTÁBIOVÍ"
		{cmn.ALPHA_DISK_LATIN, validAffineParamsES, "Amo la criptografía", "Duj ñd qózpljmódfid"},
		// ALP: "ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜẞ"
		// CNV: "DKRYBIPWẞGNUÖELSZCJQXAHOVÜFMTÄ"
		{cmn.ALPHA_DISK_GERMAN, validAffineParamsDE, "Daß liebe hübschen Mädschen", "Ydä ußbkb wtkjrwbe Öfyjrwbe"},
		// ALP: "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ"
		// CNV: "ΔΛΣΑΘΟΧΕΜΤΒΙΠΨΖΝΥΓΚΡΩΗΞΦ"
		{cmn.ALPHA_DISK_GREEK, validAffineParamsGR, "Λατρεύω την κρυπτογραφία", "Βδκυθύφ κχπ τυρνκζσυδωίδ"},
		// ALP: "АБВГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
		// CNV: "ГЙРЧЮЕЛТЩАЖНФЫВИПЦЭДКСШЯËМУЪБЗОХЬ"
		{cmn.ALPHA_DISK_CYRILLIC, validAffineParamsRU, "Мы любим криптографию", "Ыб фхйаы нцапдичцгсах"},
	}

	for vnum, tc := range allCases {
		enc := affine.NewAffineDecoder(tc.Alpha, tc.Params)
		if out, err := enc.Decode(tc.In); err != nil {
			t.Errorf("#%d %s decoder error: %v", vnum+1, tc.Alpha.LangCodeISO(), err)
		} else if out != tc.Out {
			t.Errorf("#%d %s decoder failed\n\tin :'%s'\n\texp:'%s'\n\tgot:'%s'", vnum+1, tc.Alpha.LangCodeISO(), tc.In, tc.Out, out)
		}
	}
}

/* ----------------------------------------------------------------
 *					T e s t s :: AffineCommand
 *-----------------------------------------------------------------*/

// Tests AffineCrypto Encode & Decode with only the master (letters)
// alphabet.
func Test_AffineCrypto(t *testing.T) {
	allCases := []struct {
		Alpha  *cmn.Alphabet
		Params *affine.AffineParams
		In     string
		Out    string
	}{
		// ALP: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		// CNV: "DINSXCHMRWBGLQVAFKPUZEJOTY"
		{cmn.ALPHA_DISK, validAffineParamsEN, "I love cryptography 2025", "R gvex nktauvhkdamt 2025"},
		// ALP: "ABCDEFGHIJKLMNÑOPQRSTUVWXYZÁÉÍÓÚÜ"
		// CNV: "DKQXÚFMSZAHÑUÉCJPWÓELRYÜGNTÁBIOVÍ"
		{cmn.ALPHA_DISK_LATIN, validAffineParamsES, "Amo la criptografía 2025", "Duj ñd qózpljmódfid 2025"},
		// ALP: "ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜẞ"
		// CNV: "DKRYBIPWẞGNUÖELSZCJQXAHOVÜFMTÄ"
		{cmn.ALPHA_DISK_GERMAN, validAffineParamsDE, "Daß liebe hübschen Mädschen 2025", "Ydä ußbkb wtkjrwbe Öfyjrwbe 2025"},
		// ALP: "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ"
		// CNV: "ΔΛΣΑΘΟΧΕΜΤΒΙΠΨΖΝΥΓΚΡΩΗΞΦ"
		{cmn.ALPHA_DISK_GREEK, validAffineParamsGR, "Λατρεύω την κρυπτογραφία 2025", "Βδκυθύφ κχπ τυρνκζσυδωίδ 2025"},
		// ALP: "АБВГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
		// CNV: "ГЙРЧЮЕЛТЩАЖНФЫВИПЦЭДКСШЯËМУЪБЗОХЬ"
		{cmn.ALPHA_DISK_CYRILLIC, validAffineParamsRU, "Мы любим криптографию 2025", "Ыб фхйаы нцапдичцгсах 2025"},
	}

	for vnum, tc := range allCases {
		acrypto := affine.NewAffineCrypto(tc.Alpha, tc.Params)

		var cipherStr string
		var err error
		if cipherStr, err = acrypto.Encode(tc.In); err != nil {
			t.Errorf("#%d %s encoder error: %v", vnum+1, tc.Alpha.LangCodeISO(), err)
		} else if cipherStr != tc.Out {
			t.Errorf("#%d %s encoder failed\n\tin :'%s'\n\texp:'%s'\n\tgot:'%s'", vnum+1, tc.Alpha.LangCodeISO(), tc.In, tc.Out, cipherStr)
		}

		var plain string
		if plain, err = acrypto.Decode(cipherStr); err != nil {
			t.Errorf("#%d %s decoder error: %v", vnum+1, tc.Alpha.LangCodeISO(), err)
		} else if plain != tc.In {
			t.Errorf("#%d %s decoder failed\n\tin :'%s'\n\texp:'%s'\n\tgot:'%s'", vnum+1, tc.Alpha.LangCodeISO(), cipherStr, tc.In, plain)
		}
	}
}

// Tests the AffineCommand Encode & Decode using the Master alphabet (letters)
// and a chained Slave alphabet (digits, space, some symbols)
func Test_AffineCommand_Chained(t *testing.T) {
	allCases := []struct {
		Alpha  *cmn.Alphabet
		Params *affine.AffineParams
		In     string
		Out    string
	}{
		// ALP: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		// CNV: "DINSXCHMRWBGLQVAFKPUZEJOTY"
		{cmn.ALPHA_DISK, validAffineParamsEN, "I love cryptography 2025", "R2gvex2nktauvhkdamt2%3%#"},
		// ALP: "ABCDEFGHIJKLMNÑOPQRSTUVWXYZÁÉÍÓÚÜ"
		// CNV: "DKQXÚFMSZAHÑUÉCJPWÓELRYÜGNTÁBIOVÍ"
		{cmn.ALPHA_DISK_LATIN, validAffineParamsES, "Amo la criptografía 2025", "Duj5ñd5qózpljmódfid50304"},
		// ALP: "ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÜẞ"
		// CNV: "DKRYBIPWẞGNUÖELSZCJQXAHOVÜFMTÄ"
		{cmn.ALPHA_DISK_GERMAN, validAffineParamsDE, "Daß liebe hübschen Mädschen 2025", "Ydä5ußbkb5wtkjrwbe5Öfyjrwbe50304"},
		// ALP: "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ"
		// CNV: "ΔΛΣΑΘΟΧΕΜΤΒΙΠΨΖΝΥΓΚΡΩΗΞΦ"
		{cmn.ALPHA_DISK_GREEK, validAffineParamsGR, "Λατρεύω την κρυπτογραφία 2025", "Βδκυθύφ5κχπ5τυρνκζσυδωίδ50304"},
		// ALP: "АБВГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ"
		// CNV: "ГЙРЧЮЕЛТЩАЖНФЫВИПЦЭДКСШЯËМУЪБЗОХЬ"
		{cmn.ALPHA_DISK_CYRILLIC, validAffineParamsRU, "Мы любим криптографию 2025", "Ыб5фхйаы5нцапдичцгсах50304"},
	}

	for vnum, tc := range allCases {
		acrypto := commands.NewAffineCommandExt(tc.Alpha, tc.Params)
		// add a slave alphabet
		acrypto.WithChain(cmn.NUMBERS_DISK_EXT)

		var cipherStr string
		var err error
		if cipherStr, err = acrypto.Encode(tc.In); err != nil {
			t.Errorf("#%d %s encoder error: %v", vnum+1, tc.Alpha.LangCodeISO(), err)
		} else if cipherStr != tc.Out {
			t.Errorf("#%d %s encoder failed\n\tin :'%s'\n\texp:'%s'\n\tgot:'%s'", vnum+1, tc.Alpha.LangCodeISO(), tc.In, tc.Out, cipherStr)
		}

		var plain string
		if plain, err = acrypto.Decode(cipherStr); err != nil {
			t.Errorf("#%d %s decoder error: %v", vnum+1, tc.Alpha.LangCodeISO(), err)
		} else if plain != tc.In {
			t.Errorf("#%d %s decoder failed\n\tin :'%s'\n\texp:'%s'\n\tgot:'%s'", vnum+1, tc.Alpha.LangCodeISO(), cipherStr, tc.In, plain)
		}
	}
}
