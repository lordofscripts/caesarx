/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (Command Line application)
 * A supplementary CLI application for outputing reference Tabula Rectas.
 *-----------------------------------------------------------------*/
package main

import (
	"flag"
	"fmt"
	z "lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/ciphers/caesar"
	"lordofscripts/caesarx/cmd"
	"lordofscripts/caesarx/cmn"
	"os"
	"strings"
	"unicode"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	APP_NAME = "tabularecta"
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Help about using the tabularecta CLI application
func Help() {
	fmt.Println("Usage:")
	fmt.Printf("\t%s -demo\n", APP_NAME)
	fmt.Printf("\t%s -alpha latin [-foldcase]\n", APP_NAME)
	fmt.Printf("\t%s [-foldcase] -alpha custom 'ABCDÉFÏJKLMNÖPQRST'\n", APP_NAME)
	fmt.Printf("\t%s -alpha NAME [-foldcase] -key CHAR -e CHAR\n", APP_NAME)
	fmt.Printf("\t%s -alpha NAME [-foldcase] -key CHAR -d CHAR\n", APP_NAME)
	flag.PrintDefaults()
}

// Prints the Binary Tabula Recta as an ASCII (0..255) lookup table.
func PrintBinaryTabulaRecta(key byte, asHex bool) {

	// utilitary function to print a row of 10 integers
	printRow := func(s []byte, asHex, lastRow bool) {
		for pos, v := range s {
			if lastRow && pos >= 6 { // 256..259
				break
			}

			if asHex {
				fmt.Printf("%3x ", v)
			} else {
				fmt.Printf("%3d ", v)
			}
		}
		fmt.Println()
	}

	trB := ciphers.NewBinaryTabulaRecta()

	const LEADER string = "        "
	lidIndex := 10                     // open set (value not included) highest value
	values := trB.GetTabulaForKey(key) // 0..255
	keyChar := ""
	if unicode.IsPrint(rune(key)) {
		keyChar = "· '" + string(key) + "'"
	}
	fmt.Printf("%s   Binary Tabula for Key = %d · %02Xh %s\n\n", LEADER, key, key, keyChar)

	fmt.Print(LEADER + "      ")
	printRow([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, asHex, false)
	start := 0
	for i := range 26 {
		decade := i * 10
		if asHex {
			fmt.Printf("%s%03X : ", LEADER, decade)
		} else {
			fmt.Printf("%s%03d : ", LEADER, decade)
		}

		/*
			for j := range 10 {
				if decade+j < 256 {
					if isDecode {
						values[j], err = helper.Decode(values[j])
					} else {
						values[j], err = helper.Encode(values[j])
					}

					if err != nil {
						fmt.Println("Encode error ", err)
					}
				} else {
					values[j] = -1
				}
			}
		*/

		printRow(values[start:lidIndex], asHex, i == 25)
		/*
			for j := 0; j < 10; j++ {
				decade = (i + 1) * 10
				values[j] = decade + j
			}
		*/
		start = lidIndex
		lidIndex += 10
		if lidIndex > 256 {
			lidIndex = 256
		}
	}
	fmt.Println()
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

// This application does NOT contain any Affine cipher functionality.
// It serves only the (plain) Caesar ciphers (Caesar, Didimus, Fibonacci).
// For Bellaso & Vigenère it can be used for the individual runes
// that compose the secret/password.
//
// Print an English alphabet tabula recta for ALL keys:
//
//	tabularecta -demo -alpha english
//
// Print the Binary alphabet tabula recta:
//
//	tabularecta -demo -alpha binary
//
// Encode a rune
//
//	tabularecta -alpha spanish -key M -e Ñ
//
// Decode a rune
//
//	tabularecta -alpha german -key Ü -d ẞ
//
// Print the Binary Tabula Recta for shift key
//
//	tabularecta -alpha binary -key M
//	tabularecta -alpha binary -shift 129
func main() {
	var actHelp, actDemo, optCaseFold bool
	var optAlpha string
	var optKey, actEncode, actDecode, optNumbers cmd.RuneFlag
	var optShift cmd.ByteFlag
	defer mlog.CloseLogFiles()

	// -- Command-line flags definition & parsing
	flag.BoolVar(&actHelp, "help", false, "Show help")
	flag.BoolVar(&actDemo, "demo", false, "Demonstration")
	flag.BoolVar(&optCaseFold, "foldcase", true, "Use case folding (preserves case)")
	flag.StringVar(&optAlpha, "alpha", cmn.ALPHA_NAME_ENGLISH, "Alphabet: english|latin|german|greek|cyrillic|custom")
	flag.Var(&optKey, "key", "Caesar (Main) Key")
	flag.Var(&optShift, "shift", "Caesar (Main) key as shift value 0..255")
	flag.Var(&actEncode, "e", "Rune to encode")
	flag.Var(&actDecode, "d", "Rune to decode")
	flag.Var(&optNumbers, "num", "Include Numbers disk: (N)one, (A)rabic, (H)indi (E)xtended")
	flag.Parse()

	// -- Because the time we spend is the time others are lazy
	z.Copyright(z.CO1, true)

	// -- Determine which alphabet will be used
	var alphabet *cmn.Alphabet
	var demoMsg1, demoMsg2 string
	var isBinary bool = false
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

	case cmn.ALPHA_NAME_UKRAINIAN:
		fallthrough
	case cmn.ALPHA_NAME_RUSSIAN:
		fallthrough
	case cmn.ALPHA_NAME_CYRILLIC:
		alphabet = cmn.ALPHA_DISK_CYRILLIC.Clone()
		demoMsg1 = "АБВГ0123юяьъ"
		demoMsg2 = "Я люблю криптографию"

	case cmn.ALPHA_NAME_BINARY:
		alphabet = cmn.BINARY_DISK.Clone()
		isBinary = true

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

	// -- Validate command-line flags
	if actHelp {
		Help()
		os.Exit(0)
	}

	if actEncode.IsSet && actDecode.IsSet {
		app.Die("-e and -d are mutually exclusive", z.ERR_PARAMETER)
	}

	if (actEncode.IsSet || actDecode.IsSet) && !optKey.IsSet {
		app.Die("-e and -d flags REQUIRE -key", z.ERR_PARAMETER)
	}

	if !isBinary && !actEncode.IsSet && !actDecode.IsSet && optKey.IsSet {
		app.Die("-key REQUIRE either -e or-d flags", z.ERR_PARAMETER)
	}

	if optShift.IsSet && optKey.IsSet {
		app.Die("-key and -shift are mutually exclusive", z.ERR_PARAMETER)
	}

	// further validations now that we know it is binary
	if isBinary {
		if !optShift.IsSet && !optKey.IsSet {
			app.Die("-alpha binary REQUIRES either -key or -shift flags", z.ERR_PARAMETER)
		}
		if optKey.IsSet && int(optKey.Value) >= 256 {
			app.Die("-alpha binary with -key REQUIRES that the rune be within (extended)ASCII range", z.ERR_PARAMETER)
		}
	}

	// -- Execution
	switch {
	case actDemo:
		//Demo(alphabet)
		if !isBinary {
			caesar.DemoCaesarPlain(alphabet, demoMsg1)
			caesar.DemoCaesarPlain(alphabet, demoMsg2)
		} else {
			ciphers.DemoBinaryTabulaRecta()
		}
		os.Exit(0)

	case isBinary:
		var binaryShift byte
		if optShift.IsSet {
			binaryShift = optShift.Value
		} else if optKey.IsSet {
			binaryShift = byte(optKey.Value) // we already validated that it was < 256
		}
		PrintBinaryTabulaRecta(binaryShift, true)

	case actEncode.IsSet:
		trAlpha := ciphers.NewTabulaRecta(alphabet, optCaseFold)
		result := trAlpha.EncodeRune(actEncode.Value, optKey.Value)
		fmt.Printf("(Plain Caesar) ƒ𝓍Enc(char:%c, key:%c) 🡪 %c\n", actEncode.Value, optKey.Value, result)

	case actDecode.IsSet:
		trAlpha := ciphers.NewTabulaRecta(alphabet, optCaseFold)
		result := trAlpha.DecodeRune(actDecode.Value, optKey.Value)
		fmt.Printf("(Plain Caesar) ƒ𝓍Dec(char:%c, key:%c) 🡪 %c\n", actDecode.Value, optKey.Value, result)

	case !isBinary:
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

	// -- because our time is money, you didn't have to do it, right?
	z.BuyMeCoffee()
}
