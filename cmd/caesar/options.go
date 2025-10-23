/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * CaesarX application CLI options
 *-----------------------------------------------------------------*/
package main

import (
	"errors"
	"flag"
	"fmt"
	z "lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/ciphers/bellaso"
	"lordofscripts/caesarx/ciphers/caesar"
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/ciphers/vigenere"
	"lordofscripts/caesarx/cmd"
	"lordofscripts/caesarx/cmn"
	"lordofscripts/caesarx/internal/crypto"
	"os"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	FLAG_VARIANT = "variant" // select encoding algorithm
	FLAG_NGRAM   = "ngram"   // (only for ENCODE) format output as NGram
	FLAG_OFFSET  = "offset"  // (only for Didimus) numeric offset to main key
	FLAG_DECODE  = "d"       // operation: DECODE, if not given operation is ENCODE
	FLAG_KEY     = "key"     // (only for Caesar, Didimus & Fibonacci) main encoding key
	FLAG_SECRET  = "secret"  // (only for Vigenère & Bellaso) secret password/phrase
	FLAG_FILE    = "F"       // ENCODE or DECODE files, free argument(s) are filenames
	FLAG_VERIFY  = "verify"  // (optional) ignored unless -F is used
)

const (
	// CLI application name and its alter-egos
	APP_NAME      = "caesarx"
	APP_NAME_ALT1 = "bellaso"
	APP_NAME_ALT2 = "vigenere"
	APP_NAME_ALT3 = "didimus"
	APP_NAME_ALT4 = "fibonacci"
)

const (
	NeedNone         Needs = iota // -key
	NeedKey                       // -key
	NeedCompositeKey              // -key -offset
	NeedsSecret                   // -secret
	NeedOther
)

var (
	ErrPipeTextOnly     = errors.New("for pipe input only text operations allowed")
	ErrPipeOutOnly      = errors.New("for pipe input only piped output allowed")
	ErrFreeTextRequired = errors.New("encode/decode the SINGLE free parameter must be a text string")
	ErrFilesRequired    = errors.New("encode/decode the free parameter(s) must be filename(s)")
	ErrNGramSize        = errors.New("size of NGram should be 2,3,4 or 5")
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
	UseFiles       bool
	OptVerify      bool // ignored unless -F is used
	// derived values
	ItNeeds   Needs
	VariantID CaesarVariant
	Files     *cmd.FileOptions
	fileExt   string
	isReady   bool

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
		UseFiles:       false,
		OptVerify:      false,
		Common:         common,
		isReady:        false,
		Files:          nil,
		fileExt:        "",
	}
	opts.initialize()
	return opts
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (c *CaesarxOptions) initialize() {
	flag.StringVar(&c.VariantTag, FLAG_VARIANT, crypto.ALG_NAME_CAESAR, "Algorithm (caesar|didimus|fibonacci|bellaso|vigenere)")
	flag.IntVar(&c.NGramSize, FLAG_NGRAM, -1, "Format encoded output as NGram")
	flag.IntVar(&c.Offset, FLAG_OFFSET, 0, "Alternate key offset (Didimus)")
	flag.BoolVar(&c.IsDecode, FLAG_DECODE, false, "Decode text")
	flag.BoolVar(&c.UseFiles, FLAG_FILE, false, "Free argument(s) are/is filename(s)")
	flag.BoolVar(&c.OptVerify, FLAG_VERIFY, false, "Verify operation (only if -F is used)")
	flag.Var(&c.MainKey, FLAG_KEY, "Main key")
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

	case APP_NAME_ALT4:
		c.VariantID = VariantFibonacci
		c.ItNeeds = NeedKey
		c.VariantVersion = caesar.InfoFibonacci.String()

	case APP_NAME:
		c.VariantID = VariantCaesar
		c.ItNeeds = NeedKey
		c.VariantVersion = caesar.Info.String()
		fallthrough

	default:
		switch strings.ToLower(c.VariantTag) {
		case strings.ToLower(crypto.ALG_NAME_CAESAR):
			c.VariantID = VariantCaesar
			c.fileExt = commands.FILE_EXT_CAESAR
			c.ItNeeds = NeedKey

		case strings.ToLower(crypto.ALG_NAME_DIDIMUS):
			c.VariantID = VariantDidimus
			c.fileExt = commands.FILE_EXT_DIDIMUS
			c.ItNeeds = NeedCompositeKey

		case strings.ToLower(crypto.ALG_NAME_FIBONACCI):
			c.VariantID = VariantFibonacci
			c.fileExt = commands.FILE_EXT_FIBONACCI
			c.ItNeeds = NeedKey

		case strings.ToLower(crypto.ALG_NAME_BELLASO):
			c.VariantID = VariantBellaso
			c.fileExt = commands.FILE_EXT_BELLASO
			c.ItNeeds = NeedsSecret

		case strings.ToLower(cmn.RemoveAccents(crypto.ALG_NAME_VIGENERE)):
			c.VariantID = VariantVigenere
			c.fileExt = commands.FILE_EXT_VIGENERE
			c.ItNeeds = NeedsSecret

		case strings.ToLower(crypto.ALG_NAME_AFFINE):
			c.VariantID = VariantAffine
			c.fileExt = commands.FILE_EXT_AFFINE
			c.ItNeeds = NeedNone
		}
	}
}

func (c *CaesarxOptions) ShowUsage(name string) {
	fmt.Println("Options for ALL variants:")
	fmt.Println("\t[-alpha ALPHABET] [-ngram SIZE] [-F [-verify]] [-d]")
	fmt.Println("Caesar & Fibonacci variants")
	fmt.Printf("\t%s -variant NAME -key LETTER [other options] 'user text'", name)
	fmt.Println("Didimus variant")
	fmt.Printf("\t%s -variant didimus -key LETTER -offset NUMBER [other options] 'user text'", name)
	fmt.Println("Bellaso & Vigenère variants")
	fmt.Printf("\t%s -variant NAME -secret 'password' [other options] 'user text'", name)
}

func (c *CaesarxOptions) IsReady() bool {
	return c.isReady
}

func (c *CaesarxOptions) FileExt() string {
	return c.fileExt
}

func (c *CaesarxOptions) Validate() (int, error) {
	c.checkAlterEgo()

	var err error = nil
	var exitCode int = z.EXIT_CODE_SUCCESS

	// firewall
	if c.VariantID == VariantAffine {
		// Affine not supported by caesarx executable but by its own affine program
		return z.ERR_CLI_OPTIONS, fmt.Errorf("please use the 'affine' program instead")
	} else {

	}

	// check basic needs. Except in DEMO mode
	if !c.Common.NeedsDemo() {
		// check nr. of free arguments
		if !c.Common.IsReady() { // the common options are NOT terminal
			if !app.IsPipedInput() {
				if !c.UseFiles {
					if flag.NArg() != 1 {
						return z.ERR_PARAMETER, ErrFreeTextRequired
					}
				} else { // free arguments are filenames
					numargs := 1
					if c.IsDecode {
						numargs = 2
					}

					if flag.NArg() != numargs {
						return z.ERR_PARAMETER, ErrFilesRequired
					} else {
						// now we know we have sufficient free args
						if c.IsDecode {
							c.Files = cmd.NewFileOptions(flag.Arg(0), flag.Arg(1))
						} else {
							// @note in Ring 1 the encrypted filename is auto-generated, we use the same spec here
							outputFilename := cmn.NewNameExtOnly(flag.Arg(0), c.fileExt, true)
							c.Files = cmd.NewFileOptions(flag.Arg(0), outputFilename)
						}
					}
				}
			} else { // Piped input validations
				if c.Common.IsBinary() {
					return z.ERR_PARAMETER, ErrPipeTextOnly
				}

				if c.UseFiles {
					return z.ERR_PARAMETER, ErrPipeOutOnly
				}
			}

		}

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

		// validate NGramSize
		if !c.isValidNGram() {
			err = ErrNGramSize
			exitCode = z.ERR_PARAMETER
		}

		if err != nil {
			return exitCode, err
		}
	}

	c.isReady = true
	return exitCode, nil
}

// isValidNGram verifies the NGram size validity. If it is
// out of context it is ignored (returns true). It is only
// checked for Encoding operations provided it has been set
// via the CLI
func (c *CaesarxOptions) isValidNGram() bool {
	const NOT_SET = -1
	valid := true
	if !c.IsDecode && c.NGramSize != NOT_SET {
		if c.NGramSize < 2 || c.NGramSize > 5 {
			valid = false
		}
	}

	return valid
}
