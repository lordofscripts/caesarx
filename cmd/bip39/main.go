/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Command-line utility to generate BIP39 recovery phrases. It uses
 * the internal BIP39 String Renderer. To enhance with PDF or HTML
 * output, simply implement your own bip39.IBip39Renderer.
 * It can also verify them using a Mnemonic sentence or a hex Entropy
 * string as source.
 *-----------------------------------------------------------------*/
package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	z "lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/internal/bip39"
	"os"
	"path"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	WORD_SEP rune = ' '
)

const (
	EXITCODE_EXCLUSIVE            int = 1
	EXITCODE_MISSING_ACTION       int = 2
	EXITCODE_NO_MNEMONIC_SENTENCE int = 3
	EXITCODE_BIP_LENGTH           int = 4
	EXITCODE_VERIFY               int = 20
	EXITCODE_GENERATE             int = 30
)

var (
	ErrHexDecode = errors.New("cannot decode HEX string")
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/
func init() {
	if !app.IsPipedInput() {
		z.Copyright(z.CO1, true)
		z.BuyMeCoffee()
		fmt.Println("\t=========================================")
	}
}

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// get the BIP sentence length
func getBipLength(size int) bip39.Bip39Length {
	var dummy bip39.Bip39Length
	result, err := dummy.Convert(size)
	if err != nil {
		msg := fmt.Sprintf("not a valid BIP sentence size: %d", size)
		app.Die(msg, EXITCODE_BIP_LENGTH)
	}

	return result
}

// (console output) show all known BIP39 information
func showBIP39(modeBIP bip39.Bip39Length, mnemonics []string, entropy []byte,
	renderBIP bip39.IBip39Renderer, optHex, optPlain bool) {
	fmt.Println("BIP-39 Mode:")
	fmt.Println("\t", modeBIP)
	fmt.Println("BIP-39 M n e m o n i c:")
	if optPlain {
		fmt.Println("\t", strings.Join(mnemonics, " "))
	} else {
		fmt.Print(renderBIP.FormatMnemonic(mnemonics))
	}
	fmt.Println("BIP-39 E n t r o p y:")
	fmt.Print(renderBIP.FormatEntropy(entropy, optHex))
}

// (console output) show the BIP39 seed generated from the original entropy
// and the derived reduced seed that could be used for Recoverable Codebooks.
func showBIPSeed(bip39Seed []byte, reducedSeed uint64, seedPassphrase string,
	renderBIP bip39.IBip39Renderer, optHex, optPlain bool) {
	fmt.Println("BIP-39 S e e d:")
	fmt.Printf("\tReduced Seed: %16x (%d)\n", reducedSeed, reducedSeed)
	fmt.Printf("\tSeed Passphrase: '%s'\n", seedPassphrase)
	if optPlain {
		fmt.Println("\t", hex.EncodeToString(bip39Seed))
	} else {
		fmt.Println("BIP-39 Hex Seed:")
		fmt.Print(renderBIP.FormatSeed(bip39Seed, optHex))
	}
}

// It validates the hexadecimal string and if it is valid hex, it converts it
// to an entropy. From there it generates the list of mnemonics associated with
// that entropy.
// Input:
// · Entropy as a hexadecimal string
// Output:
// · Error or nil on success
func VerifyEntropy(entropyStr string) error {
	if entropy, err := hex.DecodeString(entropyStr); err != nil {
		mlog.Err(err)
		return ErrHexDecode
	} else {
		// nr. of bytes in a BIP39 entropy
		modeBIP := bip39.BipWordCountFromEntropy(entropy)
		if !modeBIP.IsValid() {
			return fmt.Errorf("that is not a valid entropy length")
		}

		bip := bip39.NewBip39(modeBIP, WORD_SEP)
		if mnemonics, err := bip.GenerateMnemonicFromEntropy(entropy); err != nil {
			return err
		} else {
			renderBIP := bip39.NewBip39StringRenderer(modeBIP)
			showBIP39(modeBIP, mnemonics, entropy, renderBIP, false, false)
		}
	}

	return nil
}

// Validates the list of mnemonics to ensure the quantity is correct and that
// they all exist in the reference BIP39 (English) wordlist. From that list
// it derives the original entropy.
// Input:
// · list of mnemonics as a slice of strings
// Output:
// · Error or nil on success
func Verify(mnemonics []string, asHex bool) error {
	bipSize := len(mnemonics)
	modeBIP := getBipLength(bipSize)

	var err error = nil
	bip := bip39.NewBip39(modeBIP, WORD_SEP)
	if bip == nil {
		err = errors.New("invalid BIP size")
	} else if err = bip.ValidateMnemonics(mnemonics); err == nil {
		var entropy []byte
		if entropy, err = bip.EntropyFromMnemonic(mnemonics); err == nil {
			renderBIP := bip39.NewBip39StringRenderer(modeBIP)
			showBIP39(modeBIP, mnemonics, entropy, renderBIP, asHex, false)
		}
	}

	return err
}

func Generate(length int, showSeed bool, withPassphrase string, showPlainList, asHex bool) error {
	var err error = nil

	modeBIP := getBipLength(length)
	bip := bip39.NewBip39(modeBIP, WORD_SEP)

	var mnemonics []string
	if mnemonics, err = bip.GenerateMnemonic(); err == nil {
		renderBIP := bip39.NewBip39StringRenderer(modeBIP)

		mnemonicStr := bip.String()
		entropy := bip.GetEntropy()
		showBIP39(modeBIP, mnemonics, entropy, renderBIP, asHex, showPlainList)

		if showSeed {
			bip39Seed, reducedSeed := bip.ToSeed(mnemonicStr, withPassphrase)
			showBIPSeed(bip39Seed, reducedSeed, withPassphrase, renderBIP, asHex, showPlainList)
		}
	}

	return err
}

// Usage information in the form of command-line examples.
func Usage() {
	name := path.Base(os.Args[0])
	fmt.Printf("Usage of %s\n", name)
	fmt.Printf("\t%s [OPTIONS] -generate {12|15|18|21|24} [-seed [-passphrase 'TEXT']]\n", name)
	fmt.Printf("\t%s [OPTIONS] -verify 'MNEMONIC LIST'\n", name)
	fmt.Printf("\t%s [OPTIONS] -verify 'HEX_ENTROPY_STRING'\n", name)
}

// Shows usage information and information about every parameter.
func Help() {
	flag.Usage()
	flag.PrintDefaults()
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

func main() {
	var flgHelp, flgVerify, flgSeed, flgPlain, flgHex bool
	var flgSize int
	var flgPassphrase string
	flag.Usage = Usage
	// Global
	flag.BoolVar(&flgHelp, "help", false, "This help")
	flag.BoolVar(&flgPlain, "plain", false, "Show mnemonics as one string")
	flag.BoolVar(&flgHex, "hex", false, "Display entropy as Hex string")
	// Generate
	flag.IntVar(&flgSize, "generate", 0, "Mnemonic sentence length (12|15|18|21|24)")
	flag.BoolVar(&flgSeed, "seed", false, "Show seed (for -generate)")
	flag.StringVar(&flgPassphrase, "passphrase", "", "Passphrase to protect seed (with -generate and -seed)")
	// Verify
	flag.BoolVar(&flgVerify, "verify", false, "Verify a mnemonic sentence")
	flag.Parse()

	if flgHelp {
		Help()
		os.Exit(0)
	}

	if flgSize != 0 && flgVerify {
		app.Die("generate and verify are mutually exclusive", EXITCODE_EXCLUSIVE)
	}
	if flgSize == 0 && !flgVerify {
		app.Die("generate OR verify must be given", EXITCODE_MISSING_ACTION)
	}
	if flgVerify {
		if flag.NArg() != 1 {
			app.Die("mnemonic sentence must be given as argument", EXITCODE_NO_MNEMONIC_SENTENCE)
		}
		if len(flgPassphrase) > 0 {
			println("ignoring -passphrase")
		}
		if flgSeed {
			println("ignoring -seed")
		}

		if err := VerifyEntropy(flag.Arg(0)); errors.Is(err, ErrHexDecode) {
			// Argument 0 is not a hex string, thus not an Entropy string.
			// Proceed as if it is a Mnemonic sentence
			words := strings.Fields(flag.Arg(0))
			if err := Verify(words, flgHex); err != nil {
				app.DieWithError(err, EXITCODE_VERIFY)
			}
		} else if err != nil {
			app.DieWithError(err, EXITCODE_VERIFY)
		}
	} else if flgSize > 0 {
		if err := Generate(flgSize, flgSeed, flgPassphrase, flgPlain, flgHex); err != nil {
			app.DieWithError(err, EXITCODE_GENERATE)
		}
	}

	if !app.IsPipedInput() {
		z.BuyMeCoffee()
	}
}
