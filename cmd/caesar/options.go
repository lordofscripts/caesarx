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
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers/bellaso"
	"lordofscripts/caesarx/ciphers/caesar"
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/ciphers/vigenere"
	"lordofscripts/caesarx/cmd"
	"lordofscripts/caesarx/cmn"
	"lordofscripts/caesarx/cmn/prefs"
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
	VariantID z.CipherVariant
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
		VariantID:      z.CaesarCipher,
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
	// user-defined configuration defaults
	var defaultNGram int = cmd.DEFAULT_UNSET_NGRAM
	if cmd.AppConfig.IsGood() {
		defaultNGram = cmd.AppConfig.Configuration.Defaults.NGramSize
	}

	flag.StringVar(&c.VariantTag, FLAG_VARIANT, crypto.ALG_NAME_CAESAR, "Algorithm (caesar|didimus|fibonacci|bellaso|vigenere)")
	flag.IntVar(&c.NGramSize, FLAG_NGRAM, defaultNGram, "Format encoded output as NGram")
	flag.IntVar(&c.Offset, FLAG_OFFSET, 0, "Alternate key offset (Didimus)")
	flag.BoolVar(&c.IsDecode, FLAG_DECODE, false, "Decode text")
	flag.BoolVar(&c.UseFiles, FLAG_FILE, false, "Free argument(s) are/is filename(s)")
	flag.BoolVar(&c.OptVerify, FLAG_VERIFY, false, "Verify operation (only if -F is used)")
	flag.Var(&c.MainKey, FLAG_KEY, "Main key")
	flag.StringVar(&c.Secret, FLAG_SECRET, "", "Secret word/phrase used in Bellaso & Vigenere variants")
	flag.Parse()

	// check that user is requesting presets from a profile and that the profile exists. @note perhaps move elsewhere
	if cmd.AppConfig.IsGood() && c.Common.RequestsProfile() {
		profileID := c.Common.GetRequestedProfile()
		if target := cmd.AppConfig.FindProfile(profileID); target != nil {
			mlog.InfoT("Presets from ", mlog.String("ProfileID", profileID))
			// Preset cipher variant
			c.VariantID = target.Variant
			// Preset primary alphabet
			if alpha, handle := cmn.AlphabetNameByPISO(target.LangCode); alpha != nil {
				c.Common.PresetPrimaryAlphabet(handle)
			}
			// Preset slave (optional) alphabet
			c.Common.PresetSecondaryAlphabet(target.Chained)
			println("user-profile", "slave", target.Chained)
			// Preset cipher-specific parameters
			switch v := target.Params.Item.(type) {
			case *prefs.CaesarModel:
				c.MainKey.Value = rune(v.Key)
				c.MainKey.IsSet = true
				if c.VariantID == z.DidimusCipher || c.VariantID == z.FibonacciCipher {
					c.Offset = int(v.Offset)
					c.ItNeeds = NeedCompositeKey
				} else {
					c.ItNeeds = NeedKey
				}

			case *prefs.SecretsModel:
				c.Secret = v.Secret
				c.ItNeeds = NeedsSecret

			case *prefs.AffineModel:
				mlog.Fatal(z.ERR_PROFILE_CONFIG, "cannot use caesarx app with Affine parameters. Use affine app instead.")

			default:
				msg := fmt.Sprintf("unknown polymorphic parameter type %T on profile %s", v, profileID)
				mlog.Fatal(z.ERR_PROFILE_CONFIG, msg)
			}
			c.setVersion(target.Variant)
		} else {
			msg := fmt.Sprintf("couldn't find requested profile '%s'", profileID)
			warn := z.NewWarning(msg, z.CommandPCode, 1)
			mlog.Warn(warn, mlog.At())
			fmt.Println(warn)
		}
	}
}

func (c *CaesarxOptions) setVersion(variant z.CipherVariant) {
	switch variant {
	case z.CaesarCipher:
		c.VariantID = z.CaesarCipher
		c.VariantTag = crypto.ALG_NAME_CAESAR
		c.fileExt = commands.FILE_EXT_CAESAR
		c.ItNeeds = NeedKey

	case z.DidimusCipher:
		c.VariantID = z.DidimusCipher
		c.VariantTag = crypto.ALG_NAME_DIDIMUS
		c.fileExt = commands.FILE_EXT_DIDIMUS
		c.ItNeeds = NeedCompositeKey

	case z.FibonacciCipher:
		c.VariantID = z.FibonacciCipher
		c.VariantTag = crypto.ALG_NAME_FIBONACCI
		c.fileExt = commands.FILE_EXT_FIBONACCI
		c.ItNeeds = NeedKey

	case z.BellasoCipher:
		c.VariantID = z.BellasoCipher
		c.VariantTag = crypto.ALG_NAME_BELLASO
		c.fileExt = commands.FILE_EXT_BELLASO
		c.ItNeeds = NeedsSecret

	case z.VigenereCipher:
		c.VariantID = z.VigenereCipher
		c.VariantTag = crypto.ALG_NAME_VIGENERE
		c.fileExt = commands.FILE_EXT_VIGENERE
		c.ItNeeds = NeedsSecret

	case z.AffineCipher:
		c.VariantID = z.AffineCipher
		c.VariantTag = crypto.ALG_NAME_AFFINE
		c.fileExt = commands.FILE_EXT_AFFINE
		c.ItNeeds = NeedNone
	}
}

// the executable MAY have various names that would pre-configure
// the application in some way, else the default caesarx program.
func (c *CaesarxOptions) checkAlterEgo() {
	switch os.Args[0] {
	case APP_NAME_ALT1:
		c.VariantID = z.BellasoCipher
		c.ItNeeds = NeedsSecret
		c.VariantVersion = bellaso.Info.String()

	case APP_NAME_ALT2:
		c.VariantID = z.VigenereCipher
		c.ItNeeds = NeedsSecret
		c.VariantVersion = vigenere.Info.String()

	case APP_NAME_ALT3:
		c.VariantID = z.DidimusCipher
		c.ItNeeds = NeedCompositeKey
		c.VariantVersion = caesar.InfoDidimus.String()

	case APP_NAME_ALT4:
		c.VariantID = z.FibonacciCipher
		c.ItNeeds = NeedKey
		c.VariantVersion = caesar.InfoFibonacci.String()

	case APP_NAME:
		c.VariantID = z.CaesarCipher
		c.ItNeeds = NeedKey
		c.VariantVersion = caesar.Info.String()
		fallthrough

	default:
		switch strings.ToLower(c.VariantTag) {
		case strings.ToLower(crypto.ALG_NAME_CAESAR):
			c.VariantID = z.CaesarCipher
			c.fileExt = commands.FILE_EXT_CAESAR
			c.ItNeeds = NeedKey

		case strings.ToLower(crypto.ALG_NAME_DIDIMUS):
			c.VariantID = z.DidimusCipher
			c.fileExt = commands.FILE_EXT_DIDIMUS
			c.ItNeeds = NeedCompositeKey

		case strings.ToLower(crypto.ALG_NAME_FIBONACCI):
			c.VariantID = z.FibonacciCipher
			c.fileExt = commands.FILE_EXT_FIBONACCI
			c.ItNeeds = NeedKey

		case strings.ToLower(crypto.ALG_NAME_BELLASO):
			c.VariantID = z.BellasoCipher
			c.fileExt = commands.FILE_EXT_BELLASO
			c.ItNeeds = NeedsSecret

		case strings.ToLower(cmn.RemoveAccents(crypto.ALG_NAME_VIGENERE)):
			c.VariantID = z.VigenereCipher
			c.fileExt = commands.FILE_EXT_VIGENERE
			c.ItNeeds = NeedsSecret

		case strings.ToLower(crypto.ALG_NAME_AFFINE):
			c.VariantID = z.AffineCipher
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
	if c.VariantID == z.AffineCipher {
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
	valid := true
	if !c.IsDecode && c.NGramSize != cmd.DEFAULT_UNSET_NGRAM {
		if c.NGramSize < 2 || c.NGramSize > 5 {
			valid = false
		}
	}

	return valid
}
