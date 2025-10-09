/* -----------------------------------------------------------------
 *					Copyright LordOfScripts
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *						U n i t   T e s t
 * https://cryptii.com/pipes/caesar-cipher
 *-----------------------------------------------------------------*/
package tests

import (
	"lordofscripts/caesarx/cmn"
	"slices"
	"testing"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/*
	  Here are a few Caesar Cipher conversions, including borderline cases.
	 Using the original standard alphabet of 26 latin characters.
		 K  0         1         2
		 E  012345678901234567890123456
		 Y  ---------------------------
		 ~  ABCDEFGHIJKLMNOPQRSTUVWXYZ
		 0  ABCDEFGHIJKLMNOPQRSTUVWXYZ
		 5  FGHIJKLMNOPQRSTUVWXYZABCDE
		 23 XYZABCDEFGHIJKLMNOPQRSTUVW
		 26 ABCDEFGHIJKLMNOPQRSTUVWXYZ
		 27 BCDEFGHIJKLMNOPQRSTUVWXYZA
*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *				U n i t  T e s t   F u n c t i o n s
 *-----------------------------------------------------------------*/

// Tests RotateStringRight() which rotates a string to the right
// and works with strings that have UTF8 unicode characters.
func Test_RotateStringRight(t *testing.T) {
	vectors := []struct {
		Key    int
		Input  string
		Expect string
	}{
		// @note If the built-in alphabets are modified, this would NEED revision!
		{3, cmn.ALPHA_DISK.Chars, "XYZABCDEFGHIJKLMNOPQRSTUVW"},
		{3, cmn.ALPHA_DISK_LATIN.Chars, "ÓÚÜABCDEFGHIJKLMNÑOPQRSTUVWXYZÁÉÍ"},
		{3, cmn.ALPHA_DISK_GREEK.Chars, "ΧΨΩΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦ"},
		{3, cmn.ALPHA_DISK_CYRILLIC.Chars, "ЭЮЯАБВГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬ"},
		{3, cmn.ALPHA_DISK_GERMAN.Chars, "ÖÜẞABCDEFGHIJKLMNOPQRSTUVWXYZÄ"},
	}

	for nr, vector := range vectors {
		if rot := cmn.RotateStringRight(vector.Input, vector.Key); rot != vector.Expect {
			t.Errorf("#%d\n\texpected %q\n\tgot      %q\n", nr+1, vector.Expect, rot)
		}
	}
}

// Tests RotateStringLeft() which rotates a string to the left
// and works with strings that have UTF8 unicode characters.
func Test_RotateStringLeft(t *testing.T) {
	vectors := []struct {
		Key    int
		Input  string
		Expect string
	}{
		{3, cmn.ALPHA_DISK.Chars, "DEFGHIJKLMNOPQRSTUVWXYZABC"},
		{3, cmn.ALPHA_DISK_LATIN.Chars, "DEFGHIJKLMNÑOPQRSTUVWXYZÁÉÍÓÚÜABC"},
		{3, cmn.ALPHA_DISK_GREEK.Chars, "ΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩΑΒΓ"},
		{3, cmn.ALPHA_DISK_CYRILLIC.Chars, "ГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯАБВ"},
		{3, cmn.ALPHA_DISK_GERMAN.Chars, "DEFGHIJKLMNOPQRSTUVWXYZÄÖÜẞABC"},
	}

	for nr, vector := range vectors {
		if rot := cmn.RotateStringLeft(vector.Input, vector.Key); rot != vector.Expect {
			t.Errorf("#%d\n\texpected %q\n\tgot      %q\n", nr+1, vector.Expect, rot)
		}
	}
}

func Test_RotateSliceRight(t *testing.T) {
	Y0 := []byte{0, 1, 2, 3}
	Y1 := []byte{3, 0, 1, 2}
	Y3 := []byte{1, 2, 3, 0}

	y0 := cmn.RotateSliceRight(Y0, 0) // no rotation
	if slices.Compare(y0, Y0) != 0 {
		t.Errorf("zero rotation should return same slice")
	}

	y1 := cmn.RotateSliceRight(Y0, 1) // no rotation
	if slices.Compare(y1, Y1) != 0 {
		t.Errorf("rotate=1 exp:%v got:%v", Y1, y1)
	}

	y3 := cmn.RotateSliceRight(Y0, 3) // no rotation
	if slices.Compare(y3, Y3) != 0 {
		t.Errorf("rotate=3 exp:%v got:%v", Y3, y3)
	}
}

func Test_InsertNth(t *testing.T) {
	vectors := []struct {
		Input  string
		Pos    uint
		Char   rune
		Expect string
	}{
		{"ABCDEFGHIJKL", 4, '-', "ABCD-EFGH-IJKL"},
		{"ГДЕËЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯАБВ", 5, '-', "ГДЕËЖ-ЗИЙКЛ-МНОПР-СТУФХ-ЦЧШЩЪ-ЫЬЭЮЯ-АБВ"},
	}

	for vnum, v := range vectors {
		if newS := cmn.InsertNth(v.Input, v.Pos, v.Char); newS != v.Expect {
			t.Errorf("case #%d\n\texp %q\n\tgot %q", vnum+1, v.Expect, newS)
		}
	}

}

func Test_NgramFormatter(t *testing.T) {
	vectors := []struct {
		Input     string
		Len       uint8
		Separator rune
		Expect    string
	}{
		{"ABCDEFGH", 0, '·', "ABCDEFGH"},
		{"ABCDEFGH", 2, '·', "AB·CD·EF·GH"},
		{"ABCDEFGH", 3, '·', "ABC·DEF·GH"},
		{"AB CD EF GH", 3, '·', "ABC·DEF·GH"},
		{"ABCDEFGH", 4, '·', "ABCD·EFGH"},
		{"ABCDEFGH", 5, '·', "ABCDE·FGH"},
	}

	for vnum, v := range vectors {
		ngram := cmn.NewNgramFormatter(v.Len, v.Separator)
		got, err := ngram.Execute(v.Input)
		if err != nil {
			t.Errorf("vector #%d %v", vnum+1, err)
		}
		if got != v.Expect {
			t.Errorf("vector #%d NGram(%d,'%c') = '%s' got '%s'", vnum+1, v.Len, v.Separator, v.Expect, got)
		}
	}
}

func Test_RemoveAccents(t *testing.T) {
	const (
		INPUT    = "AÄÁBCÇDEËÉIÏÍOÖÓUÜÚ"
		MODIFIED = "AAABCCDEEEIIIOOOUUU"
	)
	if newS := cmn.RemoveAccents(INPUT); newS != MODIFIED {
		t.Errorf("expected %q got %q", MODIFIED, newS)
	}
}

func Test_Locate(t *testing.T) { // @note needs revision IF built-in alpha modified
	const (
		C1 rune = 'Й' // 10 in Cyrillic
		C2 rune = 'Я' // 32 in Cyrillic
		C3 rune = 'Λ' // 10 in Greek
		C4 rune = 'Ω' // 23 in Greek
	)

	if cmn.Locate(C1, cmn.ALPHA_DISK_CYRILLIC.Chars) != 10 {
		t.Errorf("cyrillic rune %c", C1)
	}

	if cmn.Locate(C2, cmn.ALPHA_DISK_CYRILLIC.Chars) != 32 {
		t.Errorf("cyrillic rune %c", C2)
	}

	if cmn.Locate(C3, cmn.ALPHA_DISK_GREEK.Chars) != 10 {
		t.Errorf("greek rune %c", C3)
	}

	if cmn.Locate(C4, cmn.ALPHA_DISK_GREEK.Chars) != 23 {
		t.Errorf("greek rune %c", C4)
	}

	if cmn.Locate(C1, cmn.ALPHA_DISK_GREEK.Chars) != -1 {
		t.Errorf("cyrillic rune %c shouldn't be found", C1)
	}

	if cmn.Locate(C3, cmn.ALPHA_DISK_CYRILLIC.Chars) != -1 {
		t.Errorf("greek rune %c shouldn't be found", C3)
	}
}

func Test_Compact(t *testing.T) {
	data := []string{"abc", "abc", "abc", "def", "def", "ghi"}
	if len(cmn.Compact(data)) != 3 {
		t.Error("compact failed")
	}
}

func Test_IsASCII(t *testing.T) {
	vectors := []struct {
		Input  *cmn.Alphabet
		Expect bool
	}{
		// plain English alphabet from the birth of computing
		{cmn.ALPHA_DISK, true},
		// these alphabets contain accented characters
		{cmn.ALPHA_DISK_LATIN, false},
		{cmn.ALPHA_DISK_GERMAN, false},
		// these alphabets contain Unicode characters, some multi-byte
		{cmn.ALPHA_DISK_GREEK, false},
		{cmn.ALPHA_DISK_CYRILLIC, false},
	}

	for nr, v := range vectors {
		if cmn.IsASCII(v.Input.Chars) != v.Expect {
			t.Errorf("#%d for %q expected IsASCII %t", nr+1, v.Input.Name, v.Expect)
		}
	}
}

func Test_IsExtendedASCII(t *testing.T) {
	vectors := []struct {
		Input  *cmn.Alphabet
		Expect bool
	}{
		{cmn.ALPHA_DISK, true},
		{cmn.ALPHA_DISK_LATIN, true},
		{cmn.ALPHA_DISK_GERMAN, true},
		{cmn.ALPHA_DISK_GREEK, false},
		{cmn.ALPHA_DISK_CYRILLIC, false},
	}

	for nr, v := range vectors {
		if cmn.IsExtendedASCII(v.Input.Chars) != v.Expect {
			t.Errorf("#%d for %q expected IsExtendedASCII %t", nr+1, v.Input.Name, v.Expect)
		}
	}
}

func Test_RuneIndex(t *testing.T) {
	const S string = "ABCΓΔΕΖËЖЗИЙXYZ"
	type Vector struct {
		R      rune
		Expect int
	}

	allCases := []Vector{
		{'A', 0},
		{'Γ', 3},
		{'Ж', 8},
		{'X', 12},
	}

	var got int
	for vnum, v := range allCases {
		got = cmn.RuneIndex(S, v.R)
		if got != v.Expect {
			t.Errorf("in #%d:%s find('%c') exp:%d got:%d", vnum+1, S, v.R, v.Expect, got)
		}
	}
}

func Test_IsMultiByteString(t *testing.T) {
	const MULTI1 string = "ABCΓΔΕΖËЖЗИЙXYZ"
	const MULTI2 string = "ABCDEFGHIJKÄÑZ"
	const SINGLE string = "ABCDEFGHIJKLMNZ"

	if !cmn.IsMultiByteString(MULTI1) {
		t.Errorf("not a multi-byte rune string? %s", MULTI1)
	}

	if !cmn.IsMultiByteString(MULTI2) {
		t.Errorf("not a multi-byte rune string? %s", MULTI2)
	}

	if cmn.IsMultiByteString(SINGLE) {
		t.Errorf("not a single-byte/rune string? %s", SINGLE)
	}
}

func Test_RuneAt(t *testing.T) {
	vectors := []struct {
		Input  string
		Expect rune
		Seek   int
	}{
		{"ABCD0123XYZ", '1', 5},
		{"ΓΔΕΖ0123ËЖЗИЙ", '2', 6},
		{"ΓΔΕΖËЖЗИЙ", 'Δ', 1},
		{"ΓΔΕΖËЖЗИЙ", 'Ж', 5},
		{"ΓΔΕΖËЖЗИЙ", 'Й', -1},
		{"ΓΔΕΖËЖЗИЙ", 'Ë', -5},
		{"ΓΔΕΖËЖЗИЙ", 'Γ', -9},
	}

	for vnum, v := range vectors {
		got := cmn.RuneAt(v.Input, v.Seek)
		if got != v.Expect {
			t.Errorf("#%d Seek('%s',%d) Exp:%c Got:%c", vnum+1, v.Input, v.Seek, v.Expect, got)
		}
	}
}

func Test_AutoKey(t *testing.T) {
	vectors := []struct {
		Input  string
		Secret string
		Expect string
	}{
		// secret smaller than plain text: autokey is secret prepended to initial part of plain
		{"abcdefghi", "xyz", "xyzabcdef"},
		{"abcdefghi", "xyz123kl", "xyz123kla"},
		// secret same length as plain : autokey is same as secret
		{"abcdefghi", "xyz123klm", "xyz123klm"},
		// secret longer than plain text : truncate to length of plain text
		{"abcdefghi", "xyz123klmn", "xyz123klm"},
		// Now the same but with Multi-byte runes
		{"abcdefghi", "xäz", "xäzabcdef"},
		{"abcüefghi", "xyz", "xyzabcüef"}, // #6
		{"abcüefghi", "xßz", "xßzabcüef"}, // #7
		// secret same length as plain : autokey is same as secret
		{"abcdefghi", "xyß123klm", "xyß123klm"}, // #8
		// secret longer than plain text : truncate to length of plain text
		{"abcdëfghi", "xyß123klmn", "xyß123klm"}, // #9
	}

	for nr, v := range vectors {
		got := cmn.AutoKeyUTF8(v.Input, v.Secret)
		if got != v.Expect {
			t.Errorf("#%d autokey(%q,%q)\n\tExp: %q\n\tGot: %q", nr+1, v.Input, v.Secret, v.Expect, got)
		}
	}
}

func Test_IntersectInt(t *testing.T) {
	set1 := []int{1, 3, 5, 7, 9}
	set2 := []int{0, 2, 4, 6, 8}
	set3 := []int{2, 3, 7}

	if len(cmn.IntersectInt(set1, set2)) != 0 {
		t.Errorf("set 1 & 2 should have nothing in common")
	}
	if len(cmn.IntersectInt(set1, set3)) != 2 {
		t.Errorf("set 1 & 3 should only have 2 in common")
	}
	if len(cmn.IntersectInt(set2, set3)) != 1 {
		t.Errorf("set 2 & 3 should only have 1 in common")
	}
}

func Test_NewName(t *testing.T) {
	allCases := []struct {
		In     string
		For    string
		Expect string
	}{
		{"/Home/Smarty/shrimp.txt", "lobster.doc", "/Home/Smarty/lobster.doc"},
		{"/Home/Smarty/shrimp", "lobster.txt", "/Home/Smarty/lobster.txt"},
	}

	for i, tc := range allCases {
		got := cmn.NewName(tc.In, tc.For)
		if got != tc.Expect {
			t.Errorf("#%d Exp: '%s' Got: '%s'", i+1, tc.Expect, got)
		}
	}
}

func Test_NewNameExtOnly(t *testing.T) {
	allCases := []struct {
		In       string
		NewExt   string
		Preserve bool
		Expect   string
	}{
		{"/Home/Smarty/shrimp.txt", ".doc", true, "/Home/Smarty/shrimp_txt.doc"},
		{"/Home/Smarty/shrimp.txt", ".doc", false, "/Home/Smarty/shrimp.doc"},
		{"/Home/Smarty/shrimp.txt", "doc", true, "/Home/Smarty/shrimp_txt.doc"},
		{"/Home/Smarty/shrimp.txt", "doc", false, "/Home/Smarty/shrimp.doc"},
		{"/Home/Smarty/shrimp.txt", "", false, "/Home/Smarty/shrimp"},
		{"/Home/Smarty/shrimp", ".doc", true, "/Home/Smarty/shrimp.doc"},
		{"/Home/Smarty/shrimp", ".doc", false, "/Home/Smarty/shrimp.doc"},
		{"/Home/Smarty/shrimp", "", false, "/Home/Smarty/shrimp.new"},
	}

	for i, tc := range allCases {
		got := cmn.NewNameExtOnly(tc.In, tc.NewExt, tc.Preserve)
		if got != tc.Expect {
			t.Errorf("#%d Exp: '%s' Got: '%s'", i+1, tc.Expect, got)
		}
	}
}

/* ----------------------------------------------------------------
 *					H e l p e r   F u n c t i o n s
 *-----------------------------------------------------------------*/
