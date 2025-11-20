/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Command-line utility to generate BIP39 recovery phrases. It uses
 * the internal BIP39 String Renderer. To enhance with PDF or HTML
 * output, simply implement your own bip39.IBip39Renderer.
 *-----------------------------------------------------------------*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/cmn"
	"lordofscripts/caesarx/internal/bip39"
	"os"
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

// Given a list of mnemonics it validates it and shows its entropy.
func Verify(mnemonics []string) error {
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
			fmt.Print(renderBIP.FormatMnemonic(mnemonics))
			fmt.Println("E n t r o p y:")
			fmt.Print(renderBIP.FormatEntropy(entropy))
			//renderMnemonic(mnemonics, false)
			//renderEntropy(modeBIP, entropy)
		}
	}

	return err
}

func Generate(length int, showSeed bool, withPassphrase string, showPlainList bool) error {
	var err error = nil

	modeBIP := getBipLength(length)
	bip := bip39.NewBip39(modeBIP, WORD_SEP)

	var mnemonics []string
	if mnemonics, err = bip.GenerateMnemonic(); err == nil {
		renderBIP := bip39.NewBip39StringRenderer(modeBIP)

		mnemonicStr := bip.String()
		fmt.Println("BIP-39 M n e m o n i c:")
		if showPlainList {
			fmt.Println("\t", bip.String())
		} else {
			fmt.Print(renderBIP.FormatMnemonic(mnemonics))
			//renderMnemonic(mnemonics, true)
		}
		//renderEntropy(modeBIP, bip.GetEntropy())
		fmt.Println("E n t r o p y:")
		fmt.Print(renderBIP.FormatEntropy(bip.GetEntropy()))

		if showSeed {
			fmt.Println("BIP-39 S e e d:")
			bip39Seed := bip.ToSeed(mnemonicStr, withPassphrase)
			reducedSeed := cmn.CalculateCRC64(bip39Seed)
			fmt.Printf("\tReduced Seed: %16x (%d)\n", reducedSeed, reducedSeed)
			fmt.Printf("\tSeed Passphrase: '%s'\n", withPassphrase)
			if showPlainList {
				fmt.Println("\t", bip.ToSeedHex(mnemonicStr, withPassphrase))
			} else {
				fmt.Println("BIP-39 Hex Seed:")
				fmt.Print(renderBIP.FormatSeed(bip39Seed))
			}
		}
	}

	return err
}

func Help() {
	flag.Usage()
	fmt.Println("Examples:")
	fmt.Println("\tbip39 -generate {12|15|18|21|24} [-seed [-passphrase 'TEXT']] [-plain]")
	fmt.Println("\tbip39 -verify 'MNEMONIC LIST'")
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

func main() {
	var flgHelp, flgVerify, flgSeed, flgPlain bool
	var flgSize int
	var flgPassphrase string
	// Global
	flag.BoolVar(&flgHelp, "help", false, "This help")
	flag.BoolVar(&flgPlain, "plain", false, "Show mnemonics as one string")
	// Generate
	flag.IntVar(&flgSize, "generate", 0, "Mnemonic sentence length (12/15/18/21/24)")
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

		words := strings.Split(flag.Arg(0), string(WORD_SEP))
		if err := Verify(words); err != nil {
			app.DieWithError(err, EXITCODE_VERIFY)
		}
	} else if flgSize > 0 {
		if err := Generate(flgSize, flgSeed, flgPassphrase, flgPlain); err != nil {
			app.DieWithError(err, EXITCODE_GENERATE)
		}
	}
}
