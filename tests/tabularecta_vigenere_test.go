package tests

import (
	"fmt"
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
