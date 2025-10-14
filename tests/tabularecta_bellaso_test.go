package tests

import (
	"fmt"
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/cmn"
	"os"
	"testing"
	"time"
)

/**
 * Cipher: Didimus (polialphabetic Caesar variant).
 * Languages: English (ASCII), Spanish, German, Greek, Cyrillic.
 * Type : Round-trip (Encode-Decode)
 * Variants: The 1st has an offset key within range, the 2nd has an
 * 			 offset key wrapped to the 2nd letter in the alphabet.
 */
func Test_BellasoCmd_RoundTrip(t *testing.T) {
	var allCases = []struct {
		Alpha     *cmn.Alphabet
		Secret    string
		Input     string
		ExpCipher string
	}{
		// @audit Decode failing when decoding a SPACE, Encoded positions with SPACE are gotten wrong!
		// 1st variant has offset within domain
		// 2nd variant has offset that wraps over and gets promoted to 1st+1
		{cmn.ALPHA_DISK, "Key", "I love cryptography", "S+jyzc3gpitrykpktfi"}, // S+jyzc3gpitrykpktfi
		{cmn.ALPHA_DISK, "Keys", "I love cryptography", "S+jgfi0ubcnlykpszlw"},
		{cmn.ALPHA_DISK_LATIN, "Key", "Años amé la criptografía", "Krhí+yvü1ue1mvazxhpvyoay"},
		{cmn.ALPHA_DISK_LATIN, "Keys", "Años amé la criptografía", "Krhf3eeñ3oy$mvacósúekjus"},
		{cmn.ALPHA_DISK_GERMAN, "Schlüßel", "Daß liebe hübschen Mädchen", "Vcg4jhimw$ojßrgswp0Xycgswp"},
		{cmn.ALPHA_DISK_GERMAN, "Eins", "Daß liebe hübschen Mädchen", "Him#pqrti1uqfäpziv6Aalpziv"},
		{cmn.ALPHA_DISK_GREEK, "τηνκρυπ", "Λατρεύω την κρυπτογραφία", "Εηηβφύτ8ννα2βμλκαγμιυμίτ"},
		{cmn.ALPHA_DISK_CYRILLIC, "кр", "Я люблю криптографию", "Й цольи хбуаэянбкеуо"},
		{cmn.ALPHA_DISK_CYRILLIC, "к", "Я люблю криптографию", "Й4цилци4хыуъэщныкяуи"},
	}

	for vnum, v := range allCases {
		alg := commands.NewBellasoCommand(v.Alpha, v.Secret)
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

// Tests text file Bellaso encryption with round-trip
// EncryptTextFile followed by DecryptTextFile
func Test_BellasoCmd_EncryptTextFile(t *testing.T) {
	// Make test file
	var fdIn *os.File
	var err error
	FILE_IN := "/tmp/test_bellaso.txt"
	FILE_OUT := cmn.NewNameExtOnly(FILE_IN, commands.FILE_EXT_BELLASO, true)
	FILE_RET := "/tmp/test_bellaso_rt.txt"
	if fdIn, err = os.Create(FILE_IN); err != nil {
		t.Error(err)
	} else {
		fdIn.WriteString("I love cryptography" + "\n")
	}

	const SECRET string = "ORALE"
	ctr := commands.NewBellasoCommand(cmn.ALPHA_DISK, SECRET)
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

// Tests Bellaso round-trip encryption of a BINARY FILE.
func Test_BellasoCmd_EncryptBinFile(t *testing.T) {
	// this depends on the encryption algorithm
	const ENC_FILE_EXT string = commands.FILE_EXT_BELLASO

	allCases := []struct {
		Secret        string
		InputFilename string // plain binary file to be encrypted
		TwinFilename  string // plain binary file after round-trip encrypt-decrypt
	}{
		{"Amor", "input.bin", "output_B.bin"},
		{"Detox", "caesar-silver-coin.png", "caesar-silver-coin-B-ret.png"},
	}

	for i, tc := range allCases {
		var err error
		var start time.Time
		var elapsed time.Duration

		assetIn := getAssetFilename(t, TEST_ASSETS, tc.InputFilename)
		assetOut := cmn.NewNameExtOnly(assetIn, ENC_FILE_EXT, true)
		assetRet := getAssetFilename(t, TEST_ASSETS, tc.TwinFilename)

		ctr := commands.NewBellasoCommand(cmn.BINARY_DISK, tc.Secret)
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
		if md5In != md5Out {
			t.Errorf("round-trip decrypted file not the same as input. %s vs %s", md5In, md5Out)
		}

		os.Remove(assetOut)
		os.Remove(assetRet)
	}
}
