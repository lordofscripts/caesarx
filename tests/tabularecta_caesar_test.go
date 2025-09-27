package tests

import (
	"lordofscripts/caesarx/ciphers/caesar"
	"lordofscripts/caesarx/ciphers/commands"
	"testing"
)

const (
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
