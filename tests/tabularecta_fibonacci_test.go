package tests

import (
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/cmn"
	"os"
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

// Tests text file Fibonacci encryption with round-trip
// EncryptTextFile followed by DecryptTextFile
func Test_FibonacciCommand_EncryptTextFile(t *testing.T) {
	// Make test file
	var fdIn *os.File
	var err error
	FILE_IN := "/tmp/test_fibonacci.txt"
	FILE_OUT := cmn.NewNameExtOnly(FILE_IN, commands.FILE_EXT_FIBONACCI, true)
	FILE_RET := "/tmp/test_fibonacci_rt.txt"
	if fdIn, err = os.Create(FILE_IN); err != nil {
		t.Error(err)
	} else {
		fdIn.WriteString("I love cryptography" + "\n")
	}

	const KEY rune = 'Z'
	ctr := commands.NewFibonacciCommand(cmn.ALPHA_DISK, KEY)
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
