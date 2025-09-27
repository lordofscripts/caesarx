/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * CaesarX application CLI options
 *-----------------------------------------------------------------*/
package main

import (
	"flag"
	"fmt"
	z "lordofscripts/caesarx"
	"lordofscripts/caesarx/ciphers/bellaso"
	"lordofscripts/caesarx/ciphers/caesar"
	"lordofscripts/caesarx/ciphers/vigenere"
	"lordofscripts/caesarx/cmd"
	"lordofscripts/caesarx/cmn"
	iciphers "lordofscripts/caesarx/internal/ciphers"
	"os"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	FLAG_VARIANT = "variant"
	FLAG_NGRAM   = "ngram"
	FLAG_OFFSET  = "offset"
	FLAG_DECODE  = "d"
	FLAG_KEY     = "key"
	FLAG_SECRET  = "secret"
)

const (
	NeedNone         Needs = iota // -key
	NeedKey                       // -key
	NeedCompositeKey              // -key -offset
	NeedsSecret                   // -secret
	NeedOther
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ cmd.IAppOptions = (*CaesarxOptions)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type CaesarxOptions struct {
	VariantTag     string
	VariantVersion string
	MainKey        cmd.RuneFlag
	Secret         string
	NGramSize      int
	Offset         int
	IsDecode       bool
	// derived values
	ItNeeds   Needs
	VariantID CaesarVariant

	Common *cmd.CommonOptions
}

type Needs uint

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * Register options common to all applications/tools and provide
 * methods to access them.
 * NOTE: Requires flag.Parse()
 */
func NewCaesarxOptions(common *cmd.CommonOptions) *CaesarxOptions {
	opts := &CaesarxOptions{
		ItNeeds:        NeedNone,
		VariantID:      VariantCaesar,
		VariantVersion: "",
		Common:         common,
	}
	opts.initialize()
	return opts
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (c *CaesarxOptions) initialize() {
	flag.StringVar(&c.VariantTag, FLAG_VARIANT, iciphers.ALG_NAME_CAESAR, "Caesar variant (caesar|didimus)")
	flag.IntVar(&c.NGramSize, FLAG_NGRAM, 0, "Format encoded output as NGram")
	flag.IntVar(&c.Offset, FLAG_OFFSET, 0, "Alternate key offset (Didimus)")
	flag.BoolVar(&c.IsDecode, FLAG_DECODE, false, "Decode text")
	flag.Var(&c.MainKey, FLAG_KEY, "Prime key")
	flag.StringVar(&c.Secret, "secret", "", "Secret word/phrase used in Bellaso & Vigenere variants")
	flag.Parse()
}

// the executable MAY have various names that would pre-configure
// the application in some way, else the default caesarx program.
func (c *CaesarxOptions) checkAlterEgo() {
	switch os.Args[0] {
	case APP_NAME_ALT1:
		c.VariantID = VariantBellaso
		c.ItNeeds = NeedsSecret
		c.VariantVersion = bellaso.Info.String()

	case APP_NAME_ALT2:
		c.VariantID = VariantVigenere
		c.ItNeeds = NeedsSecret
		c.VariantVersion = vigenere.Info.String()

	case APP_NAME_ALT3:
		c.VariantID = VariantDidimus
		c.ItNeeds = NeedCompositeKey
		c.VariantVersion = caesar.InfoDidimus.String()

	case APP_NAME:
		c.VariantID = VariantCaesar
		c.ItNeeds = NeedKey
		c.VariantVersion = caesar.Info.String()
		fallthrough

	default:
		switch strings.ToLower(c.VariantTag) {
		case strings.ToLower(iciphers.ALG_NAME_CAESAR):
			c.VariantID = VariantCaesar
			c.ItNeeds = NeedKey

		case strings.ToLower(iciphers.ALG_NAME_DIDIMUS):
			c.VariantID = VariantDidimus
			c.ItNeeds = NeedCompositeKey

		case strings.ToLower(iciphers.ALG_NAME_FIBONACCI):
			c.VariantID = VariantFibonacci
			c.ItNeeds = NeedKey

		case strings.ToLower(iciphers.ALG_NAME_BELLASO):
			c.VariantID = VariantBellaso
			c.ItNeeds = NeedsSecret

		case strings.ToLower(cmn.RemoveAccents(iciphers.ALG_NAME_VIGENERE)):
			c.VariantID = VariantVigenere
			c.ItNeeds = NeedsSecret

		case strings.ToLower(iciphers.ALG_NAME_AFFINE):
			c.VariantID = VariantAffine
			c.ItNeeds = NeedNone
		}
	}
}

func (c *CaesarxOptions) ShowUsage(name string) {
	fmt.Println("Options for ALL variants:")
	fmt.Println("\t[-alpha ALPHABET] [-ngram SIZE] [-d]")
	fmt.Println("Caesar & Fibonacci variants")
	fmt.Printf("\t%s -variant NAME -key LETTER [other options] 'user text'", name)
	fmt.Println("Didimus variant")
	fmt.Printf("\t%s -variant didimus -key LETTER -offset NUMBER [other options] 'user text'", name)
	fmt.Println("Bellaso & Vigenère variants")
	fmt.Printf("\t%s -variant NAME -secret 'password' [other options] 'user text'", name)
}

func (c *CaesarxOptions) Validate() (int, error) {
	c.checkAlterEgo()

	var err error = nil
	var exitCode int = z.EXIT_CODE_SUCCESS

	// firewall
	if c.VariantID == VariantAffine {
		// Affine not supported by caesarx executable but by its own affine program
		return z.ERR_CLI_OPTIONS, fmt.Errorf("please use the 'affine' program instead")
	}

	// check basic needs. Except in DEMO mode
	if !c.Common.NeedsDemo() {
		switch c.ItNeeds {
		case NeedCompositeKey:
			if c.Offset <= 0 {
				err = fmt.Errorf("needs offset '%s INTEGER' for composite key", FLAG_OFFSET)
				exitCode = z.ERR_CLI_OPTIONS
				break
			}
			fallthrough

		case NeedKey:
			if !c.MainKey.IsSet {
				err = fmt.Errorf("needs main key '%s LETTER'", FLAG_KEY)
				exitCode = z.ERR_CLI_OPTIONS
			}

		case NeedsSecret:
			if len(c.Secret) == 0 {
				err = fmt.Errorf("needs a secret password or phrase '%s 'SECRET'", FLAG_SECRET)
				exitCode = z.ERR_CLI_OPTIONS
			}

		case NeedOther:

		case NeedNone:

		default:
			err = fmt.Errorf("unknown necessity")
			exitCode = z.ERR_INTERNAL
		}

		if err != nil {
			return exitCode, err
		}
	}

	return exitCode, nil
}
