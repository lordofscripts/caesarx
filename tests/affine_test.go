/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package tests

import (
	"fmt"
	z "lordofscripts/caesarx"
	"lordofscripts/caesarx/ciphers/affine"
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/cmn"
	"os"
	"os/exec"
	"path"
	"runtime"
	"slices"
	"testing"
	"time"
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
	validAffineParamsEN  = &affine.AffineParams{A: 5, B: 3, Ap: 21, N: 26}
	validAffineParamsES  = &affine.AffineParams{A: 7, B: 3, Ap: 19, N: 33}
	validAffineParamsDE  = &affine.AffineParams{A: 7, B: 3, Ap: 13, N: 30}
	validAffineParamsGR  = &affine.AffineParams{A: 7, B: 3, Ap: 7, N: 24}
	validAffineParamsRU  = &affine.AffineParams{A: 7, B: 3, Ap: 19, N: 33}
	validAffineParamsBIN = &affine.AffineParams{A: 7, B: 3, Ap: 19, N: 256}
)

/* ----------------------------------------------------------------
 *					T e s t s :: AffineHelper
 *-----------------------------------------------------------------*/

// Subject: affine.AffineHelper.AreCoprime()
// Coefficient A should be Coprime of alphabet length N
func Test_Helper_AreCoprime(t *testing.T) {
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

// Subject: affine.AffineHelper.ValidCoprimesUpTo
// Compare against a known list of valid coprimes for N=26
func Test_Helper_ValidCoprimesUpTo(t *testing.T) {
	coprimes26 := []int{1, 3, 5, 7, 9, 11, 15, 17, 19, 21, 23, 25}

	h := affine.NewAffineHelper()
	if slices.Compare(h.ValidCoprimesUpTo(26), coprimes26) != 0 {
		t.Error("coprimes for N=26 failed")
	}
}

// Subject: affine.AffineHelper.ModularInverse()
// Try several A coefficients with their known A' for N=26
func Test_Helper_ModularInverse(t *testing.T) {
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

// Subject: affine.AffineHelper.VerifyParams() && affine.AffineHelper.GetParams()
func Test_Helper_VerifyParams(t *testing.T) {
	h := affine.NewAffineHelper()
	p1 := &affine.AffineParams{ // let's have Ap be fixed/recalculated
		A:  5,
		Ap: -1,
		B:  3,
		N:  26,
	}

	// we set them as well so that GetParams does not fail
	// after this p1.Ap contains a corrected value.
	err := h.VerifyParams(p1, true)
	if err != nil {
		t.Error("should not have returned error")
	}
	if p2 := h.GetParams(); p2.Ap != validAffineParamsEN.Ap {
		t.Errorf("calculated parameter A' did not get fixed. exp:%d got:%d", validAffineParamsEN.Ap, p2.Ap)
	}

	err = h.VerifyParams(validAffineParamsEN, false)
	if err != nil {
		t.Error("unexpected error for perfect parameters")
	}
}

// Subject: affine.AffineHelper.Encode() && affine.AffileHelper.SetParams()
func Test_Helper_Encode(t *testing.T) {
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

// Subject: affine.AffineHelper.Decode()
func Test_Helper_Decode(t *testing.T) {
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

// Subject: affine.AffineHelper.EncodeRuneFrom()
func Test_Helper_EncodeRuneFrom(t *testing.T) {
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

// Subject: affine.AffineHelper.DecodeRuneFrom()
func Test_Helper_DecodeRuneFrom(t *testing.T) {
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
func Test_Helper_Coprimes(t *testing.T) {
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

// Subject: affine.AffineParams{} ctor.
// Try several good & bad combinations of Affine coefficients and
// check whether the constructor returns an error or not.
func Test_AffineParams_Ctor(t *testing.T) {
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

	var err error
	var pars *affine.AffineParams
	for vnum, tc := range allCases {
		pars, err = affine.NewAffineParams(tc.A, tc.B, tc.N)
		if err != nil && tc.valid {
			t.Errorf("#%d the parameter combination should be valid. %v\n\t%s", vnum+1, err, pars)
		} else if err == nil && !tc.valid {
			t.Errorf("#%d the parameter combination should be invalid.\n\t%s", vnum+1, pars)
		}
	}
}

func Test_AffineParams_Get(t *testing.T) {
	params, err := affine.NewAffineParams(validAffineParamsEN.A,
		validAffineParamsEN.B,
		validAffineParamsEN.N)
	if err != nil {
		t.Errorf("These Affine coefficients should have been error-free:\n\tgot: %v", err)
	}

	if params.A != validAffineParamsEN.A {
		t.Errorf("retrived param A not the same exp:%d got:%d", validAffineParamsEN.A, params.A)
	}
	if params.B != validAffineParamsEN.B {
		t.Errorf("retrived param B not the same exp:%d got:%d", validAffineParamsEN.B, params.B)
	}
	if params.Ap != validAffineParamsEN.Ap {
		t.Errorf("retrived param A' not the same exp:%d got:%d", validAffineParamsEN.Ap, params.Ap)
	}
	if params.N != validAffineParamsEN.N {
		t.Errorf("retrived param N not the same exp:%d got:%d", validAffineParamsEN.N, params.N)
	}
}

/* ----------------------------------------------------------------
 *					T e s t s :: AffineCrypto
 *-----------------------------------------------------------------*/

// Subject: affine.AffineCrypto.Encode()
func Test_AffineCrypto_Encoder(t *testing.T) {
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
		enc := affine.NewAffineCrypto(tc.Alpha, tc.Params)
		if out, err := enc.Encode(tc.In); err != nil {
			t.Errorf("#%d %s encoder error: %v", vnum+1, tc.Alpha.LangCodeISO(), err)
		} else if out != tc.Out {
			t.Errorf("#%d %s encoder failed\n\tin :'%s'\n\texp:'%s'\n\tgot:'%s'", vnum+1, tc.Alpha.LangCodeISO(), tc.In, tc.Out, out)
		}
	}
}

// Subject: affine.AffineCrypto.Decode()
func Test_AffineCrypto_Decoder(t *testing.T) {
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
		enc := affine.NewAffineCrypto(tc.Alpha, tc.Params)
		if out, err := enc.Decode(tc.In); err != nil {
			t.Errorf("#%d %s decoder error: %v", vnum+1, tc.Alpha.LangCodeISO(), err)
		} else if out != tc.Out {
			t.Errorf("#%d %s decoder failed\n\tin :'%s'\n\texp:'%s'\n\tgot:'%s'", vnum+1, tc.Alpha.LangCodeISO(), tc.In, tc.Out, out)
		}
	}
}

// Tests AffineCrypto Encode & Decode with only the master (letters)
// alphabet.
func Test_AffineCrypto_Master(t *testing.T) {
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

/* ----------------------------------------------------------------
 *					T e s t s :: AffineCommand
 *-----------------------------------------------------------------*/

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

// Tests text file Affine encryption with round-trip
// EncryptTextFile followed by DecryptTextFile
func Test_AffineCommand_EncryptTextFile(t *testing.T) {
	// Make test file
	var fdIn *os.File
	var err error
	FILE_IN := "/tmp/test_affine.txt"
	FILE_OUT := cmn.NewNameExtOnly(FILE_IN, commands.FILE_EXT_AFFINE, true)
	FILE_RET := "/tmp/test_affine_rt.txt"
	if fdIn, err = os.Create(FILE_IN); err != nil {
		t.Error(err)
	} else {
		fdIn.WriteString("I love cryptography" + "\n")
	}

	ctr := commands.NewAffineCommand(cmn.ALPHA_DISK, 7, 23)
	err = ctr.EncryptTextFile(FILE_IN)
	if err != nil {
		t.Errorf("failed EncryptTextFile: %v", err)
	}

	err = ctr.DecryptTextFile(FILE_OUT, FILE_RET)
	if err != nil {
		t.Errorf("failed DecryptTextFile: %v", err)
	}

	md5In, _ := cmn.CalculateFileMD5(FILE_IN)
	md5Out, _ := cmn.CalculateFileMD5(FILE_RET)
	if md5In != md5Out {
		t.Errorf("rount-trip decrypted file not the same as input. %s vs %s", md5In, md5Out)
	}

	os.Remove(FILE_IN)
	os.Remove(FILE_OUT)
	os.Remove(FILE_RET)
}

// Tests plain Caesar round-trip encryption of a BINARY FILE. If the
// underlying EncodeBytes/DecodeBytes test do not work, then this won't either.
func Test_AffineCommand_EncryptBinFile(t *testing.T) {
	// this depends on the encryption algorithm
	const ENC_FILE_EXT string = commands.FILE_EXT_AFFINE

	allCases := []struct {
		Key           rune
		InputFilename string // plain binary file to be encrypted
		TwinFilename  string // plain binary file after round-trip encrypt-decrypt
	}{
		{'M', "input.bin", "output.bin"},
		{'Z', "caesar-silver-coin.png", "caesar-silver-coin-ret.png"},
	}

	for i, tc := range allCases {
		var err error
		var start time.Time
		var elapsed time.Duration

		assetIn := getAssetFilename(t, TEST_ASSETS, tc.InputFilename)
		assetOut := cmn.NewNameExtOnly(assetIn, ENC_FILE_EXT, true)
		assetRet := getAssetFilename(t, TEST_ASSETS, tc.TwinFilename)

		ctr := commands.NewAffineCommand(cmn.BINARY_DISK, validAffineParamsBIN.A, validAffineParamsBIN.B)
		fmt.Println("Binary File #", i+1, assetIn)
		// generate encrypted binary named assetOut
		// assetIn -> assetOut
		start = time.Now()
		if err = ctr.EncryptBinFile(assetIn); err != nil {
			t.Errorf("#%d failed EncryptBinFile: %v", i+1, err)
		}
		elapsed = time.Since(start)
		fmt.Printf("· EncryptBinFile took: %s\n", elapsed)

		// assetOut -> assetRet where to succedd assetRet == assetIn
		start = time.Now()
		if err = ctr.DecryptBinFile(assetOut, assetRet); err != nil {
			t.Errorf("#%d failed DecryptBinFile: %v", i+1, err)
		}
		elapsed = time.Since(start)
		fmt.Printf("· DecryptBinFile took: %s\n", elapsed)

		md5In, _ := cmn.CalculateFileMD5(assetIn)
		md5Out, _ := cmn.CalculateFileMD5(assetRet)
		passed := md5In == md5Out
		if passed {
			fmt.Println("· Round-trip Binary OK")
		} else {
			t.Errorf("#%d round-trip decrypted file not the same as input. %s vs %s", i+1, md5In, md5Out)
		}

		os.Remove(assetOut)
		os.Remove(assetRet)
	}
}

/* ----------------------------------------------------------------
 *				T e s t s :: Affine CLI Application
 *-----------------------------------------------------------------*/

// Test_Affine_Exit exercises the Affine executable with various CLI
// parameter/argument combinations for both valid and invalid invocations
// to check the return value. It helps ensuring the application complies
// with the documentation.
func Test_Affine_Exit(t *testing.T) {
	const OUT_PLAIN_FILE = "testdata/text_EN.txt"          // part of the repository!
	const OUT_CIPHER_FILE = "testdata/text_EN_txt.afi"     // generated
	const OUT_DECODED_FILE = "testdata/text_EN_afi_rt.txt" // generated
	// test cases for CLI execution
	allCases := []struct {
		Title    string
		ExitCode int
		Args     []string
	}{
		// common: terminal cases
		{"Help", z.EXIT_CODE_SUCCESS, []string{"-help"}},
		{"Demo", z.EXIT_CODE_SUCCESS, []string{"-demo"}},
		{"Version", z.EXIT_CODE_SUCCESS, []string{"-version"}},
		// application terminal cases
		{"Print Tabula", z.EXIT_CODE_SUCCESS, []string{"-tabula", "-A", "7", "-B", "23"}},
		{"List Coprimes", z.EXIT_CODE_SUCCESS, []string{"-coprime"}},
		{"List Coprimes For N", z.EXIT_CODE_SUCCESS, []string{"-coprime", "-N", "20"}},
		// common: chained alphabets
		{"Chained None", z.EXIT_CODE_SUCCESS, []string{"-num", "N", "-A", "7", "-B", "20", "'plain text'"}},
		{"Chained Arabic", z.EXIT_CODE_SUCCESS, []string{"-num", "A", "-A", "7", "-B", "20", "'plain text'"}},
		{"Chained Hindi", z.EXIT_CODE_SUCCESS, []string{"-num", "H", "-A", "7", "-B", "20", "'plain text'"}},
		{"Chained Extended", z.EXIT_CODE_SUCCESS, []string{"-num", "E", "-A", "7", "-B", "20", "'plain text'"}},
		{"Chained invalid", z.ERR_CLI_OPTIONS, []string{"-num", "X", "-A", "7", "-B", "20", "'plain text'"}},
		// common: NGrams
		{"NGram 2", z.EXIT_CODE_SUCCESS, []string{"-ngram", "2", "-A", "7", "-B", "20", "'plain text'"}},
		{"NGram 5", z.EXIT_CODE_SUCCESS, []string{"-ngram", "5", "-A", "7", "-B", "20", "'plain text'"}},
		{"NGram invalid", z.ERR_PARAMETER, []string{"-ngram", "6", "-A", "7", "-B", "20", "'plain text'"}},
		// application: encode cases
		{"Encode message", z.EXIT_CODE_SUCCESS, []string{"-A", "7", "-B", "20", "'plain text'"}},
		{"Encode message missing -B", z.ERR_PARAMETER, []string{"-A", "7", "'plain text'"}},
		{"Encode message missing -A", z.ERR_PARAMETER, []string{"-B", "20", "'plain text'"}},
		{"Encode message missing -A -B", z.ERR_PARAMETER, []string{"'plain text'"}},
		{"Encode file", z.EXIT_CODE_SUCCESS, []string{"-A", "7", "-B", "20", "-F", OUT_PLAIN_FILE}},
		// application: decode cases
		{"Decode message", z.EXIT_CODE_SUCCESS, []string{"-A", "7", "-B", "20", "-d", "'cipher text'"}},
		{"Decode message missing -B", z.ERR_PARAMETER, []string{"-A", "7", "-d", "'plain text'"}},
		{"Decode message missing -A", z.ERR_PARAMETER, []string{"-B", "20", "-d", "'plain text'"}},
		{"Decode message missing -A -B", z.ERR_PARAMETER, []string{"-d", "'plain text'"}},
		{"Decode file missing output", z.ERR_PARAMETER, []string{"-A", "7", "-B", "20", "-d", "-F", OUT_CIPHER_FILE}},
		{"Decode file", z.EXIT_CODE_SUCCESS, []string{"-A", "7", "-B", "20", "-d", "-F", OUT_CIPHER_FILE, OUT_DECODED_FILE}},
	}

	// @note We set this on go.yml so that this test is SKIPPED on GitHub servers
	if os.Getenv("GITHUBLOS") != "" {
		t.Skip("Skipping working test due to missing executable")
	}

	application := getAffineExecutable(t)
	for i, tc := range allCases {
		cmd := exec.Command(application, tc.Args...) // Adjust path if needed
		err := cmd.Run()

		// Check exit code
		if err == nil && tc.ExitCode == z.EXIT_CODE_SUCCESS {
			fmt.Printf("Affine %s OK\n", tc.Title)
		} else {
			if e, ok := err.(*exec.ExitError); !ok {
				t.Errorf("general failure, that is not a CLI exit error")
			} else {
				if e.ExitCode() != tc.ExitCode {
					t.Errorf("#%d [%15s] exp:%d got:%d\nError: %v", i+1, tc.Title, tc.ExitCode, e.ExitCode(), err)
				}
			}
		}
	}

	os.Remove(OUT_CIPHER_FILE)
	os.Remove(OUT_DECODED_FILE)
}

/* - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						H e l p e r s
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - */

// Get the fully-qualified path to the Affine CLI executable.
// It adjusts the name by adding ".exe" if we are on (God forbid!) Windows.
func getAffineExecutable(t *testing.T) string {
	t.Helper()

	appName := "affine" // Linux & MacOS
	if runtime.GOOS == "windows" {
		appName = "affine.exe"
	}

	// get path of current test file
	_, filename, _, _ := runtime.Caller(0)
	testDir := path.Dir(filename)

	// my projects have their own staging BIN directory
	appExecutable := path.Join(cmn.Conjoin(testDir, "../bin"), appName)

	return appExecutable
}
