/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Affine application CLI options
 *-----------------------------------------------------------------*/
package main

import (
	"errors"
	"flag"
	"fmt"
	z "lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/cmd"
	"lordofscripts/caesarx/cmn"
	"lordofscripts/caesarx/cmn/prefs"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	FLAG_COEFF_A  = "A"
	FLAG_NGRAM    = "ngram" // (optional) only if encrypting
	FLAG_COEFF_B  = "B"
	FLAG_DECODE   = "d"
	FLAG_COPRIMES = "coprime"
	FLAG_MODULO   = "N" // (optional) only if -coprime is given
	FLAG_TABULA   = "tabula"
	FLAG_FILE     = "F"      // (optional) free args are filenames and not strings
	FLAG_VERIFY   = "verify" // (optional) ignored unless -F is used
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ cmd.IAppOptions = (*AffineCliOptions)(nil)

var (
	ErrNeedAffineCoefficients = errors.New("need -A and -B Affine coefficients")
	ErrModuloNotNeeded        = errors.New("modulo -N defaults to alphabet length, but can only be used with -coprime")
	ErrInvalidModulo          = errors.New("modulo -N must be a positive integer")
	ErrNGramSize              = errors.New("size of NGram should be 2,3,4 or 5")
	ErrFreeTextRequired       = errors.New("for encode/decode text the SINGLE free parameter must be a string")
	ErrFilesRequired          = errors.New("for encode/decode a file the 2 free parameters must be input and output filenames")
	ErrPipeTextOnly           = errors.New("for pipe input only text operations allowed")
	ErrPipeOutOnly            = errors.New("for pipe input only piped output allowed")
)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type AffineCliOptions struct {
	CoefficientA    int
	CoefficientB    int
	OptModulo       int
	OptNgramSize    int
	OptUseFiles     bool
	OptVerify       bool // ignored unless -F is used
	ActListCoprimes bool
	ActPrintTabula  bool
	ActIsDecode     bool

	isReady bool
	Files   *cmd.FileOptions
	Common  *cmd.CommonOptions
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * Register options common to all applications/tools and provide
 * methods to access them.
 * NOTE: Requires flag.Parse()
 */
func NewAffineOptions(common *cmd.CommonOptions) *AffineCliOptions {
	opts := &AffineCliOptions{
		CoefficientA:    1,
		CoefficientB:    -1,
		OptModulo:       0,
		OptNgramSize:    0,
		OptUseFiles:     false,
		OptVerify:       false,
		ActListCoprimes: false,
		ActPrintTabula:  false,
		ActIsDecode:     false,
		isReady:         false,
		Files:           nil,
		Common:          common,
	}
	opts.initialize()
	return opts
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (c *AffineCliOptions) initialize() {
	// user-defined configuration defaults
	var defaultNGram int = -1
	if cmd.AppConfig.IsGood() {
		defaultNGram = cmd.AppConfig.Configuration.Defaults.NGramSize
	}

	flag.IntVar(&c.CoefficientA, FLAG_COEFF_A, -1, "Affine coefficient A")
	flag.IntVar(&c.CoefficientB, FLAG_COEFF_B, -1, "Affine coefficient B")
	flag.IntVar(&c.OptModulo, FLAG_MODULO, 0, "Affine module N (only if -coprime is used), else derived from alpha")
	flag.IntVar(&c.OptNgramSize, FLAG_NGRAM, defaultNGram, "Format encoded output as NGram")
	flag.BoolVar(&c.OptUseFiles, FLAG_FILE, false, "Free argument(s) are filenames")
	flag.BoolVar(&c.OptVerify, FLAG_VERIFY, false, "Verify operation (only if -F is used)")
	flag.BoolVar(&c.ActIsDecode, FLAG_DECODE, false, "Decode text")
	flag.BoolVar(&c.ActListCoprimes, FLAG_COPRIMES, false, "List coprimes for 'A' for the chosen alphabet")
	flag.BoolVar(&c.ActPrintTabula, FLAG_TABULA, false, "Print Tabula for chosen parameters")
	flag.Parse()

	// check that user is requesting presets from a profile and that the profile exists. @note perhaps move elsewhere
	if cmd.AppConfig.IsGood() && c.Common.RequestsProfile() {
		profileID := c.Common.GetRequestedProfile()
		if target := cmd.AppConfig.FindProfile(profileID); target != nil {
			fmt.Println("PRsetting from", profileID)
			mlog.InfoT("Presets from ", mlog.String("ProfileID", profileID))
			// Preset cipher variant should be fixed to Affine, else abort
			if target.Variant != z.AffineCipher {
				mlog.Fatalf(z.ERR_PROFILE_CONFIG, "preset specifies %s cipher and this app is Affine-specific", target.Variant)
			}

			// Preset primary alphabet
			if alpha, handle := cmn.AlphabetNameByPISO(target.LangCode); alpha != nil {
				c.Common.PresetPrimaryAlphabet(handle)
			}
			// Preset slave (optional) alphabet
			c.Common.PresetSecondaryAlphabet(target.Chained)
			// Preset cipher-specific parameters
			switch v := target.Params.Item.(type) {
			case *prefs.AffineModel:
				c.CoefficientA = int(v.A)
				c.CoefficientB = int(v.B)

			case *prefs.SecretsModel, *prefs.CaesarModel:
				mlog.Fatalf(z.ERR_PROFILE_CONFIG, "found non-Affine parameter model in configuration profile")

			default:
				msg := fmt.Sprintf("unknown polymorphic parameter type %T on profile %s", v, profileID)
				mlog.Fatal(z.ERR_PROFILE_CONFIG, msg)
			}
		} else {
			msg := fmt.Sprintf("couldn't find requested profile '%s'", profileID)
			warn := z.NewWarning(msg, z.CommandPCode, 1)
			mlog.Warn(warn, mlog.At())
			fmt.Println(warn)
		}
	}
}

func (c *AffineCliOptions) ShowUsage(name string) {
	flag.Usage()
}

// IsReady indicates whether a previous Validate invocation was successful.
func (c *AffineCliOptions) IsReady() bool {
	return c.isReady
}

func (c *AffineCliOptions) FileExt() string {
	return commands.FILE_EXT_AFFINE
}

// Validate validates application CLI parameters. If it is successful
// it returns EXIT_CODE_SUCCESS with nil error.
func (c *AffineCliOptions) Validate() (int, error) {

	var err error = nil
	var exitCode int = z.EXIT_CODE_SUCCESS

	// check for terminal options that don't require anything else
	if c.Common.NeedsDemo() || c.Common.NeedsHelp() ||
		c.Common.NeedsVersion() || c.ActListCoprimes {
		c.isReady = true
	} else { // non-terminal arguments
		if !(c.OptNgramSize == 0 || (c.OptNgramSize >= 2 && c.OptNgramSize <= 5)) {
			err = ErrNGramSize
		} else if c.OptModulo != 0 && !c.ActListCoprimes {
			err = ErrModuloNotNeeded
		} else if c.OptModulo < 0 {
			err = ErrInvalidModulo
		} else if c.CoefficientA == -1 || c.CoefficientB == -1 { // for -tabula and -d we require -A and -B
			err = ErrNeedAffineCoefficients
		} else if !c.ActPrintTabula {
			// encode OR decode (-d) operation requested
			if !app.IsPipedInput() {
				if c.OptUseFiles { // -F given
					// 2 free arguments are input & output filenames respectively
					if c.ActIsDecode { // -F -d ciphered_filename output_filename
						if flag.NArg() != 2 {
							err = ErrFilesRequired
						} else {
							c.Files = cmd.NewFileOptions(flag.Arg(0), flag.Arg(1))
						}
					} else { // -F plain_filename
						if flag.NArg() != 1 {
							err = ErrFilesRequired
						} else {
							// @note in Ring 1 the encrypted filename is auto-generated, we use the same spec here
							outputFilename := cmn.NewNameExtOnly(flag.Arg(0), commands.FILE_EXT_AFFINE, true)
							c.Files = cmd.NewFileOptions(flag.Arg(0), outputFilename)
						}
					}

				} else { // single free argument is plain OR cipher string
					if flag.NArg() != 1 {
						err = ErrFreeTextRequired
					}
				}
			} else {
				// Validations for exclusively Piped input
				if c.Common.IsBinary() {
					return z.ERR_PARAMETER, ErrPipeTextOnly
				}

				if c.OptUseFiles {
					return z.ERR_PARAMETER, ErrPipeOutOnly
				}
			}
		}
	}

	if err != nil {
		exitCode = z.ERR_PARAMETER
		c.isReady = false
	} else {
		c.isReady = true
	}

	return exitCode, err
}

// UseFiles indicates whether the encrypt/decrypt operation will work
// with input/output file instead of a (short) text string.
func (c *AffineCliOptions) UseFiles() bool {
	return c.Files != nil
}
