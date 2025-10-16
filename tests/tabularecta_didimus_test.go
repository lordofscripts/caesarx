package tests

import (
	"fmt"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/cmn"
	"lordofscripts/caesarx/internal/crypto"
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
func Test_DidimusCmd_RoundTrip(t *testing.T) {
	mlog.Info("info message", mlog.String("String", "texto"), mlog.YesNo("Valid", true), mlog.Int("Value", 23))
	var allCases = []struct {
		Alpha     *cmn.Alphabet
		PrimeKey  rune
		Offset    uint8
		Input     string
		ExpAltKey rune
		ExpCipher string
	}{
		// @audit Decode failing when decoding a SPACE, Encoded positions with SPACE are gotten wrong!
		// 1st variant has offset within domain
		// 2nd variant has offset that wraps over and gets promoted to 1st+1
		{cmn.ALPHA_DISK, 'M', 5, "I love cryptography", 'R', "U xfhv5tdpbkaxdrbyk"},
		{cmn.ALPHA_DISK, 'M', 15, "I love cryptography", 'B', "U#xphf5ddzbuahdbbik"},
		{cmn.ALPHA_DISK_LATIN, 'M', 5, "Años amé la criptografía", 'Q', "Múád5qxm5ém ñctaüürcmviq"},
		{cmn.ALPHA_DISK_LATIN, 'M', 22, "Años amé la criptografía", 'B', "Moát5bxí5mm#ñstqüprsmgib"},
		{cmn.ALPHA_DISK_GERMAN, 'M', 5, "Daß liebe hübschen Mädschen", 'R', "Prl xzqsq tpnfoyqa5ẞiuattvz"},
		{cmn.ALPHA_DISK_GERMAN, 'M', 19, "Daß liebe hübschen Mädschen", 'B', "Pbl#xjqcq#tßntoiqo5Nieadtfz"},
		{cmn.ALPHA_DISK_GREEK, 'Λ', 5, "Λατρεύω την κρυπτογραφία", 'Π', "Φπεθούο3κρδ3αγλβκασγπηίπ"},
		{cmn.ALPHA_DISK_GREEK, 'Λ', 15, "Λατρεύω την κρυπτογραφία", 'Β', "Φβεσούα3υρξ3λγφβυαδγβηίβ"},
		{cmn.ALPHA_DISK_CYRILLIC, 'Л', 5, "Я люблю криптографию", 'Р', "К чомьй цбфаюяоблефо"},
		{cmn.ALPHA_DISK_CYRILLIC, 'Л', 22, "Я люблю криптографию", 'Б', "К#чяммй#цсфрюпослхфя"},
	}

	const EVEN_POS = 0
	const ODD_POS = 1
	const DUMMY = 'x'
	for vnum, v := range allCases {
		var k rune
		seq := crypto.NewDidimusSequencer(v.PrimeKey, v.Offset, v.Alpha)
		// Test PrimeKey
		if k = seq.GetKey(EVEN_POS, DUMMY); k != v.PrimeKey {
			t.Errorf("#%d PrimeKey exp: %c got: %c", vnum+1, v.PrimeKey, k)
		}
		// Test AltKey
		if k = seq.GetKey(ODD_POS, DUMMY); k != v.ExpAltKey {
			t.Errorf("#%d AltKey exp: %c got: %c", vnum+1, v.ExpAltKey, k)
		}

		alg := commands.NewDidimusCommand(v.Alpha, v.PrimeKey, v.Offset)
		var cipherStr, decodedStr string
		var err error
		// Test Encode
		cipherStr, err = alg.Encode(v.Input)
		if err != nil {
			t.Errorf("#%d unexpected encode error: %v", vnum+1, err)
		} else if cipherStr != v.ExpCipher {
			t.Errorf("#%d Encode fail\n\tInp: %s\n\texp: %s\n\tgot: %s", vnum+1, v.Input, v.ExpCipher, cipherStr)
		}
		// Test Decode
		decodedStr, err = alg.Decode(cipherStr)
		if err != nil {
			t.Errorf("#%d unexpected decode error: %v", vnum+1, err)
		} else if decodedStr != v.Input {
			t.Errorf("#%d Decode fail\n\tInp: %s\n\texp: %s\n\tgot: %s ", vnum+1, cipherStr, v.Input, decodedStr)
		}
	}
}

// Tests text file Didimus encryption with round-trip
// EncryptTextFile followed by DecryptTextFile
func Test_DidimusCommand_EncryptTextFile(t *testing.T) {
	// Make test file
	var fdIn *os.File
	var err error
	FILE_IN := "/tmp/test_didimus.txt"
	FILE_OUT := cmn.NewNameExtOnly(FILE_IN, commands.FILE_EXT_DIDIMUS, true)
	FILE_RET := "/tmp/test_didimus_rt.txt"
	if fdIn, err = os.Create(FILE_IN); err != nil {
		t.Error(err)
	} else {
		fdIn.WriteString("I love cryptography" + "\n")
	}

	const KEY rune = 'Z'
	const OFFSET uint8 = 8
	ctr := commands.NewDidimusCommand(cmn.ALPHA_DISK, KEY, OFFSET)
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

// Tests Didimus round-trip encryption of a BINARY FILE.
func Test_DidimusCommand_EncryptBinFile(t *testing.T) {
	// this depends on the encryption algorithm
	const ENC_FILE_EXT string = commands.FILE_EXT_DIDIMUS

	allCases := []struct {
		Key           rune
		Offset        uint8
		InputFilename string // plain binary file to be encrypted
		TwinFilename  string // plain binary file after round-trip encrypt-decrypt
	}{
		{'M', 12, "input.bin", "output_D.bin"},
		{'Z', 20, "caesar-silver-coin.png", "caesar-silver-coin-D-ret.png"},
	}

	for i, tc := range allCases {
		var err error
		var start time.Time
		var elapsed time.Duration

		assetIn := getAssetFilename(t, TEST_ASSETS, tc.InputFilename)
		assetOut := cmn.NewNameExtOnly(assetIn, ENC_FILE_EXT, true)
		assetRet := getAssetFilename(t, TEST_ASSETS, tc.TwinFilename)

		ctr := commands.NewDidimusCommand(cmn.BINARY_DISK, tc.Key, tc.Offset)
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
