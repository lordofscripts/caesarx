package tests

import (
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/cmn"
	"testing"
)

/**
 * Cipher: Didimus (polialphabetic Caesar variant).
 * Languages: English (ASCII), Spanish, German, Greek, Cyrillic.
 * Type : Round-trip (Encode-Decode)
 * Variants: The 1st has an offset key within range, the 2nd has an
 * 			 offset key wrapped to the 2nd letter in the alphabet.
 */
func Test_FibonacciCmd_RoundTrip(t *testing.T) {
	var allCases = []struct {
		Alpha     *cmn.Alphabet
		PrimeKey  rune
		Input     string
		ExpCipher string
	}{
		// @audit Decode failing when decoding a SPACE, Encoded positions with SPACE are gotten wrong!
		// 1st variant has offset within domain
		// 2nd variant has offset that wraps over and gets promoted to 1st+1
		{cmn.ALPHA_DISK, 'M', "I love cryptography", "U6yckv%bysbgbugrjgf"},
		{cmn.ALPHA_DISK, 'Z', "I love cryptography", "H#mpxi0olfouphtewts"},
		{cmn.ALPHA_DISK_LATIN, 'M', "Años amé la criptografía", "Máéa8qüt#xm6oüwahhhúmrjñ"},
		{cmn.ALPHA_DISK_LATIN, 'Z', "Años amé la criptografía", "Zijñ5únb7fz3íneñuutmzüwé"},
		{cmn.ALPHA_DISK_GERMAN, 'M', "Daß liebe hübschen Mädchen", "Pnm7äzyäh9tlocryyi%Üiqpvta"},
		{cmn.ALPHA_DISK_GERMAN, 'ẞ', "Daß liebe hübschen Mädchen", "Cba#nmlny%gßctellz%Pzedigr"},
		{cmn.ALPHA_DISK_GREEK, 'Λ', "Λατρεύω την κρυπτογραφία", "Φμζεσύσ@βγψ4φειηνξκνλθίν"},
		{cmn.ALPHA_DISK_GREEK, 'Ω', "Λατρεύω την κρυπτογραφία", "Κβυσηύη5οπμ#λσχυβγψβωχίβ"},
		{cmn.ALPHA_DISK_CYRILLIC, 'Л', "Я люблю криптографию", "К6шлпьс1лэфьяьсбумйк"},
		{cmn.ALPHA_DISK_CYRILLIC, 'Я', "Я люблю криптографию", "Ю#мягпе5юсзрупефжаья"},
	}

	const DUMMY = 'x'
	for vnum, v := range allCases {
		alg := commands.NewFibonacciCommand(v.Alpha, v.PrimeKey)
		var cipherStr, decodedStr string
		var err error
		// Test Encode
		cipherStr, err = alg.Encode(v.Input)
		if err != nil {
			t.Errorf("#%d unexpected encode error: %v", vnum+1, err)
		} else if cipherStr != v.ExpCipher {
			t.Errorf("#%d Encode fail\n\texp: '%s'\n\tgot: '%s'", vnum+1, v.ExpCipher, cipherStr)
		}

		// Test Decode
		decodedStr, err = alg.Decode(cipherStr)
		if err != nil {
			t.Errorf("#%d unexpected decode error: %v", vnum+1, err)
		} else if decodedStr != v.Input {
			t.Errorf("#%d Decode %s fail\n\texp: '%s'\n\tgot: '%s'", vnum+1, v.Alpha.Name, v.Input, decodedStr)
		}
	}
}
