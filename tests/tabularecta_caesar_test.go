package tests

import (
	"encoding/binary"
	"fmt"
	z "lordofscripts/caesarx"
	"lordofscripts/caesarx/ciphers/caesar"
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/cmn"
	"os"
	"os/exec"
	"path"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"
)

const (
	DOC_ASSETS          string = "../docs/assets"
	TEST_ASSETS         string = "test_data"
	ERR_FORMAT_VECTORED string = "#%02d (%s) Expected %q Got %q"
)

/*
  TODO: Caesar WithAlphabet(), Encode(), Decode()
*/

/*
CORE VERSIONS BASED DIRECTLY ON TabulaRecta
*/
func Test_Caesar_Encode_Decode(t *testing.T) {
	type Vector struct {
		Input  string
		Expect string
	}

	for _, alpha := range AllAlphabets {
		lastKey := alpha.GetRuneAt(-1)
		ctr := caesar.NewCaesarTabulaRecta(alpha, lastKey)

		// @note ALL of these will be done with the last letter in the alphabet!
		// It ensures that we pass through all multi-byte runes.
		// @note you can confirm the test vectors at https://cryptii.com/
		var V1, V2, V3 Vector
		if IsEnglish(alpha) {
			V1 = Vector{"CMZ", "BLY"}
			V2 = Vector{"C M Z", "B L Y"}
			V3 = Vector{"I love cryptography", "H knud bqxosnfqzogx"}
		} else if IsSpanish(alpha) {
			V1 = Vector{"CÑÜ", "BNÚ"}
			V2 = Vector{"C Ñ Ü", "B N Ú"}
			V3 = Vector{"Amo la criptografía", "Ülñ kü bqhosñfqüeéü"}
		} else if IsGerman(alpha) {
			V1 = Vector{"CMẞ", "BLÜ"}
			V2 = Vector{"C M ẞ", "B L Ü"}
			V3 = Vector{"Daß liebe hübschen Mädschen", "Cßü khdad göarbgdm Lzcrbgdm"}
		} else if IsGreek(alpha) {
			V1 = Vector{"ΓΞΩ", "ΒΝΨ"}
			V2 = Vector{"Γ Ξ Ω", "Β Ν Ψ"}
			V3 = Vector{"Λατρεύω την κρυπτογραφία", "Κωσπδύψ σζμ ιπτοσξβπωυίω"}
		} else if IsCyrillic(alpha) {
			V1 = Vector{"ЖПЯ", "ËОЮ"}
			V2 = Vector{"Ж П Я", "Ë О Ю"}
			V3 = Vector{"Мы любим криптографию", "Лъ кэазл йпзоснвпяузэ"}
		}

		for vnum, tv := range []Vector{V1, V2, V3} {
			got := ctr.Encode(tv.Input)
			if got != tv.Expect {
				t.Errorf("«%s» vector #%d ENC\n\tKey:%c\n\tExp:'%s'\n\tGot:'%s'", alpha.Name, vnum+1, lastKey, tv.Expect, got)
			}

			plain := ctr.Decode(got)
			if plain != tv.Input {
				t.Errorf("«%s» vector #%d DEC\n\tKey:%c\n\tExp:'%s'\n\tGot:'%s'", alpha.Name, vnum+1, lastKey, tv.Input, plain)
			}
		}
	}
}

/*
  WRAPPER VERSIONS BASED ON COMMAND PATTERN
*/

func Test_Caesar_Encode_Decode_CommandPattern(t *testing.T) {
	type Vector struct {
		Input  string
		Expect string
	}

	for _, alpha := range AllAlphabets {
		lastKey := alpha.GetRuneAt(-1)
		ctr := commands.NewCaesarCommand(alpha, lastKey)

		// @note ALL of these will be done with the last letter in the alphabet!
		// It ensures that we pass through all multi-byte runes.
		// @note you can confirm the test vectors at https://cryptii.com/
		var V1, V2, V3 Vector
		if IsEnglish(alpha) {
			V1 = Vector{"CMZ", "BLY"}
			V2 = Vector{"C M Z", "B L Y"}
			V3 = Vector{"I love cryptography", "H knud bqxosnfqzogx"}
		} else if IsSpanish(alpha) {
			V1 = Vector{"CÑÜ", "BNÚ"}
			V2 = Vector{"C Ñ Ü", "B N Ú"}
			V3 = Vector{"Amo la criptografía", "Ülñ kü bqhosñfqüeéü"}
		} else if IsGerman(alpha) {
			V1 = Vector{"CMẞ", "BLÜ"}
			V2 = Vector{"C M ẞ", "B L Ü"}
			V3 = Vector{"Daß liebe hübschen Mädschen", "Cßü khdad göarbgdm Lzcrbgdm"}
		} else if IsGreek(alpha) {
			V1 = Vector{"ΓΞΩ", "ΒΝΨ"}
			V2 = Vector{"Γ Ξ Ω", "Β Ν Ψ"}
			V3 = Vector{"Λατρεύω την κρυπτογραφία", "Κωσπδύψ σζμ ιπτοσξβπωυίω"}
		} else if IsCyrillic(alpha) {
			V1 = Vector{"ЖПЯ", "ËОЮ"}
			V2 = Vector{"Ж П Я", "Ë О Ю"}
			V3 = Vector{"Мы любим криптографию", "Лъ кэазл йпзоснвпяузэ"}
		}

		for vnum, tv := range []Vector{V1, V2, V3} {
			got, err := ctr.Encode(tv.Input)
			if err != nil {
				t.Errorf("«%s» vector #%d ENC Key:%c\n\tError: %v", alpha.Name, vnum+1, lastKey, err)
			}
			if got != tv.Expect {
				t.Errorf("«%s» vector #%d ENC Key:%c\n\tExp:'%s'\n\tGot:'%s'", alpha.Name, vnum+1, lastKey, tv.Expect, got)
			}

			var plain string
			plain, err = ctr.Decode(got)
			if err != nil {
				t.Errorf("«%s» vector #%d DEC Key:%c\n\tError: %v", alpha.Name, vnum+1, lastKey, err)
			}
			if plain != tv.Input {
				t.Errorf("«%s» vector #%d DEC Key:%c\n\tExp:'%s'\n\tGot:'%s'", alpha.Name, vnum+1, lastKey, tv.Input, plain)
			}
		}
	}
}

// Tests text file Caesar encryption with round-trip
// EncryptTextFile followed by DecryptTextFile
func Test_CaesarCommand_EncryptTextFile(t *testing.T) {
	// Make test file
	var fdIn *os.File
	var err error
	FILE_IN := "/tmp/test_caesar.txt"
	FILE_OUT := cmn.NewNameExtOnly(FILE_IN, commands.FILE_EXT_CAESAR, true)
	FILE_RET := "/tmp/test_caesar_rt.txt"
	if fdIn, err = os.Create(FILE_IN); err != nil {
		t.Error(err)
	} else {
		fdIn.WriteString("I love cryptography" + "\n")
	}

	const KEY rune = 'Z'
	ctr := commands.NewCaesarCommand(cmn.ALPHA_DISK, KEY)
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
func Test_EncodeDecodeBytes_Caesar(t *testing.T) {
	allCases := []struct {
		Key    rune
		Plain  uint32
		Cipher uint32
	}{
		{'M', 0x44332211, 0x91806f5e}, // BigEndian
	}

	var start time.Time
	var elapsed time.Duration

	for _, tc := range allCases {
		ctr := caesar.NewCaesarTabulaRecta(cmn.BINARY_DISK, tc.Key)

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

// Tests plain Caesar round-trip encryption of a BINARY FILE. If the
// underlying EncodeBytes/DecodeBytes test do not work, then this won't either.
func Test_CaesarCommand_EncryptBinFile(t *testing.T) {
	// this depends on the encryption algorithm
	const ENC_FILE_EXT string = commands.FILE_EXT_CAESAR

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

		ctr := commands.NewCaesarCommand(cmn.BINARY_DISK, tc.Key)
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

/*
func Test_WithAlphabet(t *testing.T) {
	const KEY = 3
	const GREEK_RESULT = "ΔΕΖΗΘ"
	const LATIN_RESULT = "DóEFGHúüabc"
	const CYRILLIC_RESULT = "ГДЕЖЗ"

	// test vectors with various foreign alphabets based on Key=3
	vectors := [3]struct {
		Alpha  *cmn.Alphabet
		Input  string
		Expect string
	}{
		{cmn.ALPHA_DISK_LATIN, "AáBCDEéíóúü", LATIN_RESULT},
		{cmn.ALPHA_DISK_GREEK, "ΑΒΓΔΕ", GREEK_RESULT},
		{cmn.ALPHA_DISK_CYRILLIC, "АБВГД", CYRILLIC_RESULT},
	}

	// prepare
	options := caesar.CaesarOptions{Variant: caesar.CAESAR, Initial: KEY, Supplemental: 0, CleanAccents: false}
	cmd := caesar.NewCaesarCommand(options)

	// try subcases
	for _, vector := range vectors {
		fmt.Println(cmd.TabulaRecta())

		cmd.WithAlphabet(vector.Alpha)
		if cmd.Alphabet() != vector.Alpha.Name {
			t.Errorf("· Disc name mismatch %s != %s\n", cmd.Alphabet(), vector.Alpha.Name)
		}

		if encoded, err := cmd.Encode(vector.Input); err == nil {
			if encoded != vector.Expect {
				t.Errorf("%s Alphabet() Exp. %s Got %s\n", vector.Alpha.Name, vector.Expect, encoded)
			} else {
				fmt.Printf("· %s alphabet: OK", vector.Alpha.Name)
			}
		} else {
			t.Error(err)
		}
	}
}

// / Test CaesarCipherCommand.Encode() & Decode()
func Test_Caesar_RoundTrip(t *testing.T) {
	var vectors []Vector = []Vector{
		{"This is a test", "This is a test", caesar.CaesarOptions{Variant: caesar.CAESAR, Initial: 0, Supplemental: 0, CleanAccents: false}},
		{"This is a test", "Wklv lv d whvw", caesar.CaesarOptions{Variant: caesar.CAESAR, Initial: 3, Supplemental: 0, CleanAccents: false}},
		{"This is a test", "Aopz pz h alza", caesar.CaesarOptions{Variant: caesar.CAESAR, Initial: 7, Supplemental: 0, CleanAccents: false}}, // borderline
		{"This is a test", "Sghr hr z sdrs", caesar.CaesarOptions{Variant: caesar.CAESAR, Initial: 25, Supplemental: 0, CleanAccents: false}},
		{"This is a test", "This is a test", caesar.CaesarOptions{Variant: caesar.CAESAR, Initial: 26, Supplemental: 0, CleanAccents: false}}, // borderline
		{"This is a test", "Aopz pz h alza", caesar.CaesarOptions{Variant: caesar.CAESAR, Initial: 7, Supplemental: 0, CleanAccents: false}},
	}

	for idx, vector := range vectors {
		cmd := caesar.NewCaesarCommand(vector.Params)

		fmt.Printf("Vector #%02d %s\n", idx+1, vector.Params.String())
		fmt.Println("\tIn : ", vector.Input)
		if encoded, err := cmd.Encode(vector.Input); err == nil {
			fmt.Println("\tExp: ", vector.Expect)
			fmt.Println("\tOut: ", encoded)
			if encoded != vector.Expect {
				t.Errorf(ERR_FORMAT_VECTORED, idx+1, vector.Params.LeaderString(), vector.Expect, encoded)
			}

			decoded, _ := cmd.Decode(encoded)
			if decoded != vector.Input {
				t.Errorf(ERR_FORMAT_VECTORED, idx+1, vector.Params.LeaderString(), vector.Input, decoded)
			}
		} else {
			t.Error(err)
		}
	}
}
*/

// Test_Affine_Exit exercises the Affine executable with various CLI
// parameter/argument combinations for both valid and invalid invocations
// to check the return value. It helps ensuring the application complies
// with the documentation.
// @note something odd happening with this, at times it reports the wrong
// exitCode even though the constant is correct!
func Test_Caesar_Exit(t *testing.T) {
	const OUT_PLAIN_FILE = "test_data/text_EN.txt"              // part of the repository!
	const OUT_CIPHER_FILE_CAE = "test_data/text_EN_txt.cae"     // generated
	const OUT_DECODED_FILE_CAE = "test_data/text_EN_cae_rt.txt" // generated
	const OUT_CIPHER_FILE_VIG = "test_data/text_EN_txt.vig"     // generated
	const OUT_CIPHER_FILE_BEL = "test_data/text_EN_txt.bel"     // generated

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
		{"List ciphers", z.EXIT_CODE_SUCCESS, []string{"-list"}},
		// common: chained alphabets
		{"Chained None", z.EXIT_CODE_SUCCESS, []string{"-num", "N", "-key", "L", "'plain text'"}},
		{"Chained Arabic", z.EXIT_CODE_SUCCESS, []string{"-num", "A", "-key", "L", "'plain text'"}},
		{"Chained Hindi", z.EXIT_CODE_SUCCESS, []string{"-num", "H", "-key", "L", "'plain text'"}},
		{"Chained Extended", z.EXIT_CODE_SUCCESS, []string{"-num", "E", "-key", "L", "'plain text'"}},
		{"Chained invalid", z.ERR_CLI_OPTIONS, []string{"-num", "V", "-key", "L", "'plain text'"}},
		// common: NGrams
		{"NGram 2", z.EXIT_CODE_SUCCESS, []string{"-ngram", "2", "-key", "L", "'plain text'"}},
		{"NGram 5", z.EXIT_CODE_SUCCESS, []string{"-ngram", "5", "-key", "L", "'plain text'"}},
		{"NGram invalid", z.ERR_PARAMETER, []string{"-ngram", "6", "-key", "L", "'plain text'"}}, // @audit 1 got 2
		// application: encode cases
		{"Encode Caesar message", z.EXIT_CODE_SUCCESS, []string{"-key", "L", "'plain text'"}},
		{"Encode Caesar message missing -key", z.ERR_CLI_OPTIONS, []string{"'plain text'"}},
		{"Encode Caesar file", z.EXIT_CODE_SUCCESS, []string{"-key", "L", "-F", OUT_PLAIN_FILE}},
		// application: decode cases
		{"Decode Caesar message", z.EXIT_CODE_SUCCESS, []string{"-key", "L", "-d", "'cipher text'"}},
		{"Decode Caesar message missing -key", z.ERR_CLI_OPTIONS, []string{"-d", "'plain text'"}},
		{"Decode Caesar file missing output", z.ERR_PARAMETER, []string{"-key", "L", "-d", "-F", OUT_CIPHER_FILE_CAE}},
		{"Decode Caesar file", z.EXIT_CODE_SUCCESS, []string{"-key", "L", "-d", "-F", OUT_CIPHER_FILE_CAE, OUT_DECODED_FILE_CAE}},
		// application: other abnormal cases
		{"Redirect -variant affine", z.ERR_CLI_OPTIONS, []string{"-variant", "affine", "'plain text'"}},
		// other...
		{"Didimus", z.EXIT_CODE_SUCCESS, []string{"-num", "E", "-variant", "didimus", "-key", "L", "-offset", "3", "'plain text'"}},
		{"Didimus missing -offset", z.ERR_CLI_OPTIONS, []string{"-num", "E", "-variant", "didimus", "-key", "L", "'plain text'"}},
		{"Fibonacci", z.EXIT_CODE_SUCCESS, []string{"-num", "E", "-variant", "fibonacci", "-key", "L", "'plain text'"}},
		{"Bellaso Message", z.EXIT_CODE_SUCCESS, []string{"-num", "E", "-variant", "bellaso", "-secret", "PASSWD", "'plain text'"}},
		{"Bellaso File", z.EXIT_CODE_SUCCESS, []string{"-num", "E", "-variant", "bellaso", "-secret", "PASSWD", "-F", OUT_PLAIN_FILE}},
		{"Bellaso missing -secret", z.ERR_CLI_OPTIONS, []string{"-num", "E", "-variant", "bellaso", "-key", "P", "'plain text'"}},
		{"Vigenere Message", z.EXIT_CODE_SUCCESS, []string{"-num", "E", "-variant", "vigenere", "-secret", "PASSWD", "'plain text'"}},
		{"Vigenere missing -secret", z.ERR_CLI_OPTIONS, []string{"-num", "E", "-variant", "vigenere", "-key", "P", "'plain text'"}},
		{"Vigenere File", z.EXIT_CODE_SUCCESS, []string{"-num", "E", "-variant", "vigenere", "-secret", "PASSWD", "-F", OUT_PLAIN_FILE}},
	}

	application := getCaesarExecutable(t, "caesarx")
	for i, tc := range allCases {
		cmd := exec.Command(application, tc.Args...) // Adjust path if needed
		err := cmd.Run()

		// Check exit code
		if err == nil && tc.ExitCode == z.EXIT_CODE_SUCCESS {
			fmt.Printf("Caesar %02d %12s OK\n", i+1, tc.Title)
		} else {
			if e, ok := err.(*exec.ExitError); !ok {
				t.Errorf("general failure, that is not a CLI exit error. %v", err)
			} else {
				if e.ExitCode() != tc.ExitCode {
					t.Errorf("#%d [%15s] exp:%d got:%d\nArgs: %v\nError: %v", i+1, tc.Title, tc.ExitCode, e.ExitCode(), tc.Args, err)
				}
			}
		}
	}

	os.Remove(OUT_CIPHER_FILE_CAE)
	os.Remove(OUT_DECODED_FILE_CAE)
	os.Remove(OUT_CIPHER_FILE_VIG)
	os.Remove(OUT_CIPHER_FILE_BEL)
}

func getAssetFilename(t *testing.T, where, asset string) string {
	t.Helper()

	_, filename, _, _ := runtime.Caller(0)
	testDir := path.Dir(filename)

	// my projects have their own staging BIN directory
	filename = path.Join(cmn.Conjoin(testDir, where), asset)
	return filename
}

// Get the fully-qualified path to the Caesar CLI executable.
// It adjusts the name by adding ".exe" if we are on (God forbid!) Windows.
// Keeping in mind that vigenere, bellaso, didimus and fibonacci are symlinks
// to the base caesarx executable
func getCaesarExecutable(t *testing.T, appName string) string {
	t.Helper()

	appName = strings.ToLower(appName)
	switch appName {
	case "caesarx":
	case "didimus":
	case "fibonacci":
	case "bellaso":
	case "vigenere":

	default:
		panic(appName + " is not a valid executable name")
	}

	if runtime.GOOS == "windows" {
		appName = appName + ".exe"
	}

	// get path of current test file
	_, filename, _, _ := runtime.Caller(0)
	testDir := path.Dir(filename)

	// my projects have their own staging BIN directory
	appExecutable := path.Join(cmn.Conjoin(testDir, "../bin"), appName)

	return appExecutable
}
