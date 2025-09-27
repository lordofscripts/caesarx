package main

import (
	"flag"
	"fmt"
	. "lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/ciphers/caesar"
	"lordofscripts/caesarx/cmd"
	"lordofscripts/caesarx/cmn"
	"os"
	"strings"
	"unicode"
)

const (
	APP_NAME = "tabularecta"
)

func Help() {
	fmt.Println("Usage:")
	fmt.Printf("\t%s -demo\n", APP_NAME)
	fmt.Printf("\t%s -alpha latin [-foldcase]\n", APP_NAME)
	fmt.Printf("\t%s [-foldcase] -alpha custom ABCDÉFÏJKLMNÖPQRST\n", APP_NAME)
	fmt.Printf("\t%s -alpha NAME [-foldcase] -key CHAR -e CHAR\n", APP_NAME)
	fmt.Printf("\t%s -alpha NAME [-foldcase] -key CHAR -d CHAR\n", APP_NAME)
	flag.PrintDefaults()
}

func main() {
	var actHelp, actDemo, optCaseFold bool
	var optAlpha string
	var optKey, actEncode, actDecode, optNumbers cmd.RuneFlag
	flag.BoolVar(&actHelp, "help", false, "Show help")
	flag.BoolVar(&actDemo, "demo", false, "Demonstration")
	flag.BoolVar(&optCaseFold, "foldcase", true, "Use case folding (preserves case)")
	flag.StringVar(&optAlpha, "alpha", "english", "Alphabet: english|latin|german|greek|cyrillic|custom")
	flag.Var(&optKey, "key", "Caesar Key")
	flag.Var(&actEncode, "e", "Rune to encode")
	flag.Var(&actDecode, "d", "Rune to decode")
	flag.Var(&optNumbers, "num", "Include Numbers disk: (N)one, (A)rabic, (H)indi (E)xtended")
	flag.Parse()

	Copyright(CO1, true)

	if actHelp {
		Help()
		os.Exit(0)
	}

	if actEncode.IsSet && actDecode.IsSet {
		app.Die("-e and -d are mutually exclusive", ERR_PARAMETER)
	}

	if (actEncode.IsSet || actDecode.IsSet) && !optKey.IsSet {
		app.Die("-e and -d flags REQUIRE -key", ERR_PARAMETER)
	}

	if !actEncode.IsSet && !actDecode.IsSet && optKey.IsSet {
		app.Die("-key REQUIRE either -e or-d flags", ERR_PARAMETER)
	}

	var alphabet *cmn.Alphabet
	var demoMsg1, demoMsg2 string
	optAlpha = strings.ToLower(optAlpha)
	switch optAlpha {
	case cmn.ALPHA_NAME_ENGLISH:
		alphabet = cmn.ALPHA_DISK.Clone()
		demoMsg1 = "ABCD0123wxyz"
		demoMsg2 = "I love cryptography"

	case cmn.ALPHA_NAME_SPANISH:
		fallthrough
	case cmn.ALPHA_NAME_LATIN:
		alphabet = cmn.ALPHA_DISK_LATIN.Clone()
		demoMsg1 = "ABÍÑ0123wxyz"
		demoMsg2 = "Amo la criptografía"

	case cmn.ALPHA_NAME_GREEK:
		alphabet = cmn.ALPHA_DISK_GREEK.Clone()
		demoMsg1 = "ΑΒΓΔ0123φχψω"
		demoMsg2 = "Λατρεύω την κρυπτογραφία"

	case cmn.ALPHA_NAME_GERMAN:
		alphabet = cmn.ALPHA_DISK_GERMAN.Clone()
		demoMsg1 = "ABCD0123äöüß"
		demoMsg2 = "Daß liebe hübschen Mädchen"

	case cmn.ALPHA_NAME_UKRANIAN:
		fallthrough
	case cmn.ALPHA_NAME_RUSSIAN:
		fallthrough
	case cmn.ALPHA_NAME_CYRILLIC:
		alphabet = cmn.ALPHA_DISK_CYRILLIC.Clone()
		demoMsg1 = "АБВГ0123юяьъ"
		demoMsg2 = "Я люблю криптографию"

	case "custom":
		if flag.NArg() == 1 {
			alphabet = cmn.NewAlphabet("Custom", flag.Arg(0), false, false)
		} else {
			fmt.Println("When using '-alpha custom' you must specify a non-spaced strings of characters")
			os.Exit(2)
		}

	default:
		fmt.Println("Valid alphabets are: english|latin|german|greek|cyrillic|custom")
		os.Exit(1)
	}

	if actDemo {
		//Demo(alphabet)
		caesar.DemoCaesarPlain(alphabet, demoMsg1)
		caesar.DemoCaesarPlain(alphabet, demoMsg2)
		os.Exit(0)
	} else if actEncode.IsSet {
		trAlpha := ciphers.NewTabulaRecta(alphabet, optCaseFold)
		result := trAlpha.EncodeRune(actEncode.Value, optKey.Value)
		fmt.Printf("(Plain Caesar) ƒ𝓍Enc(char:%c, key:%c) 🡪 %c\n", actEncode.Value, optKey.Value, result)
	} else if actDecode.IsSet {
		trAlpha := ciphers.NewTabulaRecta(alphabet, optCaseFold)
		result := trAlpha.DecodeRune(actDecode.Value, optKey.Value)
		fmt.Printf("(Plain Caesar) ƒ𝓍Dec(char:%c, key:%c) 🡪 %c\n", actDecode.Value, optKey.Value, result)
	} else {
		// Prepare Numbers disk if solicited
		var numerics *cmn.Alphabet = nil
		if optNumbers.IsSet {
			switch unicode.ToUpper(optNumbers.Value) {
			case 'A': // Arabic Numbers only
				numerics = cmn.NUMBERS_DISK.Clone()

			case 'H': // Hindi Numbers only
				numerics = cmn.NUMBERS_EASTERN_DISK.Clone()

			case 'E': // Arabic numbers, space and number-related chars
				numerics = cmn.NUMBERS_DISK_EXT.Clone()

			default:
			}
		}

		ciphers.DemoTabulaRecta(alphabet, optCaseFold, numerics)
		fmt.Println(("\tSee -help for more options!"))
	}

	BuyMeCoffee()
}
