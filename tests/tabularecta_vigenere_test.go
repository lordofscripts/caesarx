package tests

import (
	"encoding/binary"
	"fmt"
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/ciphers/vigenere"
	"lordofscripts/caesarx/cmn"
	"os"
	"slices"
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
func Test_VigenereCmd_RoundTrip(t *testing.T) {
	var allCases = []struct {
		Alpha     *cmn.Alphabet
		Secret    string
		Input     string
		ExpCipher string
	}{
		// @audit Decode failing when decoding a SPACE, Encoded positions with SPACE are gotten wrong!
		// 1st variant has offset within domain
		// 2nd variant has offset that wraps over and gets promoted to 1st+1
		{cmn.ALPHA_DISK, "KEY", "MESSAGE", "WIQEEYW"},
		{cmn.ALPHA_DISK, "Key", "I love cryptography", "S+jwvp7xvyrkmvkovyy"},
		{cmn.ALPHA_DISK, "Keys", "I love cryptography", "S+jgde4qmcpvfegtdnp"},
		{cmn.ALPHA_DISK_LATIN, "Llave", "Años amé la criptografía", "Lyoi+azk$la5órsptqxzpylg"},
		{cmn.ALPHA_DISK_LATIN, "LlavesMaestras", "Años amé la", "Lyoi+sxé+ót"},
		{cmn.ALPHA_DISK_GERMAN, "Schlüßel", "Daß liebe hübschen Mädchen", "Vcg4jhimh gümägiin0Köveoiä"},
		{cmn.ALPHA_DISK_GERMAN, "Ein", "Daß liebe hübschen Mädchen", "Him%lhemm+icbzaiwp0Qjdodhp"},
		{cmn.ALPHA_DISK_GREEK, "τηνκρυπ", "Λατρεύω την κρυπτογραφία", "Εηηβφύτ8εηη9ξπυκαγγβρπίτ"},
		{cmn.ALPHA_DISK_CYRILLIC, "кр", "Я люблю криптографию", "Й кюмйя5ируаыюхягеит"},
		{cmn.ALPHA_DISK_CYRILLIC, "кфюЖ", "Я люблю криптографию", "Й+йеалй7льжпэялатгло"},
	}

	defer func() {
		if r := recover(); r == nil {
		} else {
			if msg, ok := r.(string); !ok || msg != "message message" {
				t.Errorf("unexpected panic value %v", r)
			}
		}
	}()

	for vnum, v := range allCases {
		alg := commands.NewVigenereCommand(v.Alpha, v.Secret)
		var cipherStr, decodedStr string
		var err error
		var outcome string = "Fail"
		// Test Encode
		fmt.Printf("Encode 𝑽𝑒('%s','%s') ", v.Secret, v.Input)
		cipherStr, err = alg.Encode(v.Input)
		if cipherStr == v.ExpCipher {
			outcome = "OK"
		}
		fmt.Println(outcome)

		if err != nil {
			t.Errorf("#%d unexpected encode error: %v", vnum+1, err)
		} else if cipherStr != v.ExpCipher {
			t.Errorf("#%d Encode fail\n\tIn : %s\n\texp: %s\n\tgot: %s", vnum+1, v.Input, v.ExpCipher, cipherStr)
		}

		// Test Decode
		fmt.Printf("Decode 𝑽𝑑('%s','%s') ", v.Secret, cipherStr)
		decodedStr, err = alg.Decode(cipherStr)
		if decodedStr == v.Input {
			outcome = "OK"
		} else {
			outcome = "Fail"
		}
		fmt.Println(outcome)
		if err != nil {
			t.Errorf("#%d unexpected decode error: %v", vnum+1, err)
		} else if decodedStr != v.Input {
			t.Errorf("#%d Decode %s fail\n\tIn : %s\n\texp: %s\n\tgot: %s", vnum+1, cipherStr, v.Alpha.Name, v.Input, decodedStr)
		}
	}
}

// Tests text file Vigenère encryption with round-trip
// EncryptTextFile followed by DecryptTextFile
func Test_VigenereCmd_EncryptTextFile(t *testing.T) {
	// Make test file
	var fdIn *os.File
	var err error
	FILE_IN := "/tmp/test_vigenere.txt"
	FILE_OUT := cmn.NewNameExtOnly(FILE_IN, commands.FILE_EXT_VIGENERE, true)
	FILE_RET := "/tmp/test_vigenere_rt.txt"
	if fdIn, err = os.Create(FILE_IN); err != nil {
		t.Error(err)
	} else {
		fdIn.WriteString("I love cryptography" + "\n")
	}

	const SECRET string = "ORALE"
	ctr := commands.NewVigenereCommand(cmn.ALPHA_DISK, SECRET)
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

// Plain Caesar encryption and decryption of a BINARY buffer.
// It tests the underlying BinaryTabula EncodeBytes() & DecodeBytes()
// which are the low-level functions for binary FILE encryption.
func Test_EncodeDecodeBytes_Vigenere(t *testing.T) {
	allCases := []struct {
		Secret string
		Plain  uint32
		Cipher uint32
	}{
		{"Amor", 0x44332211, 0x85807163}, // BigEndian
	}

	var start time.Time
	var elapsed time.Duration

	for _, tc := range allCases {
		ctr := vigenere.NewVigenereTabulaRecta(cmn.BINARY_DISK, tc.Secret)

		fmt.Printf("DataIn  0x%08x\n", tc.Plain)
		dataIn := make([]byte, 4)
		binary.BigEndian.PutUint32(dataIn, tc.Plain)

		start = time.Now()
		dataOut := ctr.EncodeBytes(dataIn)
		elapsed = time.Since(start)
		value := binary.BigEndian.Uint32(dataOut)
		fmt.Printf("DataOut 0x%08x Took: %s\n", value, elapsed)
		if value != tc.Cipher {
			t.Errorf("encodeBytes failed. Exp: %08xh Got: %08xh", tc.Cipher, value)
		}

		start = time.Now()
		dataRet := ctr.DecodeBytes(dataOut)
		elapsed = time.Since(start)
		value = binary.BigEndian.Uint32(dataRet)
		fmt.Printf("DataRet 0x%08x Took: %s\n", value, elapsed)
		if value != tc.Plain {
			t.Errorf("decodeBytes failed. Exp: %08xh Got: %08xh", tc.Plain, value)
		}

		if slices.Compare(dataIn, dataRet) != 0 {
			t.Errorf("Decrypted data does not match plain binary input")
		}
	}
}

// Tests Vigenere round-trip encryption of a BINARY FILE.
func Test_VigenereCmd_EncryptBinFile(t *testing.T) {
	// this depends on the encryption algorithm
	const ENC_FILE_EXT string = commands.FILE_EXT_VIGENERE

	allCases := []struct {
		Secret        string
		InputFilename string // plain binary file to be encrypted
		TwinFilename  string // plain binary file after round-trip encrypt-decrypt
	}{
		{"Amor", "input.bin", "output_V.bin"}, // 0x11223344556677880a => 0x526f82966688aacc5f
		{"Detox", "caesar-silver-coin.png", "caesar-silver-coin-V-ret.png"},
	}

	for i, tc := range allCases {
		var err error
		var start time.Time
		var elapsed time.Duration

		assetIn := getAssetFilename(t, TEST_ASSETS, tc.InputFilename)
		assetOut := cmn.NewNameExtOnly(assetIn, ENC_FILE_EXT, true)
		assetRet := getAssetFilename(t, TEST_ASSETS, tc.TwinFilename)
		fmt.Printf("Binary file #%d - %s\n", i+1, assetIn)

		ctr := commands.NewVigenereCommand(cmn.BINARY_DISK, tc.Secret)
		// generate encrypted binary named assetOut
		// assetIn -> assetOut
		start = time.Now()
		if err = ctr.EncryptBinFile(assetIn); err != nil {
			t.Errorf("#%d failed EncryptBinFile: %v", i+1, err)
		}
		elapsed = time.Since(start)
		fmt.Printf("· EncryptBinFile #%d took: %s\n", i+1, elapsed)

		// assetOut -> assetRet where to succedd assetRet == assetIn
		start = time.Now()
		if err = ctr.DecryptBinFile(assetOut, assetRet); err != nil {
			t.Errorf("#%d failed DecryptBinFile: %v", i+1, err)
		}
		elapsed = time.Since(start)
		fmt.Printf("· DecryptBinFile #%d took: %s\n", i+1, elapsed)

		md5In, _ := cmn.CalculateFileMD5(assetIn)
		md5Out, _ := cmn.CalculateFileMD5(assetRet)
		if md5In != md5Out {
			t.Errorf("round-trip decrypted file not the same as input. %s vs %s", md5In, md5Out)
		}

		os.Remove(assetOut)
		os.Remove(assetRet)
	}
}
