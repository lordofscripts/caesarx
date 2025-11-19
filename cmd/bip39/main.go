/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Command-line utility to generate BIP39 recovery phrases.
 *-----------------------------------------------------------------*/
package main

import (
	"flag"
	"fmt"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/internal/bip39"
	"os"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// The number of words in the BIP39 mnemonic sentence
	Bip39Words12 Bip39Length = iota
	Bip39Words15
	Bip39Words18
	Bip39Words21
	Bip39Words24
)

const (
	WORD_SEP rune = ' '
)

var (
	RenderMap = map[Bip39Length]Renderer{ // E-BitSize E-ByteSize
		Bip39Words12: {Table{4, 3}, Table{4, 4}}, // 128	16
		Bip39Words15: {Table{5, 3}, Table{5, 4}}, // 160 20
		Bip39Words18: {Table{6, 3}, Table{6, 4}}, // 192 24
		Bip39Words21: {Table{7, 3}, Table{7, 4}}, // 224 28
		Bip39Words24: {Table{6, 4}, Table{8, 4}}, // 256 32
	}
)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type Bip39Length int

type Table struct {
	Rows int
	Cols int
}
type Renderer struct {
	Sentence Table
	Entropy  Table
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// get the BIP sentence length
func getBipLength(size int) Bip39Length {
	var result Bip39Length
	switch size {
	case 12:
		result = Bip39Words12

	case 15:
		result = Bip39Words15

	case 18:
		result = Bip39Words18

	case 21:
		result = Bip39Words21

	case 24:
		result = Bip39Words24

	default:
		msg := fmt.Sprintf("not a valid BIP sentence size: %d", size)
		app.Die(msg, 10)
	}

	return result
}

// renders the mnemonic sentence as a table
func renderMnemonic(mnemonic []string) {
	bipLen := getBipLength(len(mnemonic))
	table := RenderMap[bipLen].Sentence

	fmt.Println("M n e m o n i c s:")
	for row := range table.Rows {
		fmt.Print("\t")
		for col := range table.Cols {
			// In BIP39 English the maximum word length is 8
			offset := table.Cols * row
			fmt.Printf("%-10s", mnemonic[offset+col])
		}
		fmt.Println()
	}
}

// render the entropy as a table
func renderEntropy(modeBIP Bip39Length, entropy []byte) {
	table := RenderMap[modeBIP].Sentence

	fmt.Println("E n t r o p y:")
	for row := range table.Rows {
		fmt.Print("\t")
		for col := range table.Cols {
			// In BIP39 English the maximum word length is 8
			offset := table.Cols * row
			fmt.Printf("%5d", entropy[offset+col])
		}
		fmt.Println()
	}
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
		flag.Usage()
		os.Exit(0)
	}

	var modeBIP Bip39Length
	if flgSize != 0 && flgVerify {
		app.Die("generate and verify are mutually exclusive", 1)
	}
	if flgSize == 0 && !flgVerify {
		app.Die("generate OR verify must be given", 2)
	}
	if flgVerify {
		if flag.NArg() != 1 {
			app.Die("mnemonic sentence must be given as argument", 3)
		}
		if len(flgPassphrase) > 0 {
			println("ignoring -passphrase")
		}
		if flgSeed {
			println("ignoring -seed")
		}

		words := strings.Split(flag.Arg(0), string(WORD_SEP))
		bipSize := len(words)
		modeBIP = getBipLength(bipSize)

		bip := bip39.NewBip39(bipSize, WORD_SEP)
		if entropy, err := bip.EntropyFromMnemonic(words); err != nil {
			app.DieWithError(err, 22)
		} else {
			renderMnemonic(words)
			renderEntropy(modeBIP, entropy)
		}
	} else if flgSize > 0 {
		modeBIP = getBipLength(flgSize)
		bip := bip39.NewBip39(flgSize, WORD_SEP)
		if mnemonics, err := bip.GenerateMnemonic(); err != nil {
			app.DieWithError(err, 20)
		} else {
			mnemonicStr := bip.String()
			fmt.Println("BIP-39 Mnemonic:")
			if flgPlain {
				fmt.Println("\t", bip.String())
			} else {
				renderMnemonic(mnemonics)
			}
			renderEntropy(modeBIP, bip.GetEntropy())

			if flgSeed {
				fmt.Println("BIP-39 Seed Passphrase:")
				fmt.Println("\tPassphrase:", flgPassphrase)
				fmt.Println("BIP-39 Hex Seed:")
				fmt.Println("\t", bip.ToSeedHex(mnemonicStr, flgPassphrase))
			}
		}
	}
}
