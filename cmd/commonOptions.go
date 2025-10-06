/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package cmd

import (
	"flag"
	"fmt"
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/cmn"
	"slices"
	"strings"
	"unicode"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	defaultLanguage    string = "english"
	supportedAlphabets string = "english|latin|spanish|german|greek|cyrillic|custom"
	supportedNumbers   string = "(N)one (A)rabic (E)xtended (H)indi"
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type IAppOptions interface {
	ShowUsage(string)
	Validate() (int, error)
	IsReady() bool
}

var _ IAppOptions = (*CommonOptions)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type CommonOptions struct {
	DefaultPhrase string
	demo          bool
	help          bool
	list          bool
	version       bool
	encodeSpace   bool
	alpha         string
	numeric       RuneFlag
	isReady       bool
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * Register options common to all applications/tools and provide
 * methods to access them.
 * NOTE: Requires flag.Parse()
 */
func NewCommonOptions() *CommonOptions {
	copts := &CommonOptions{}
	copts.DefaultPhrase = "Let's encrypt!"
	copts.initialize()
	return copts
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (c *CommonOptions) NeedsHelp() bool {
	return c.help
}

func (c *CommonOptions) NeedsDemo() bool {
	return c.demo
}

func (c *CommonOptions) NeedsList() bool {
	return c.list
}

func (c *CommonOptions) NeedsVersion() bool {
	return c.version
}

func (c *CommonOptions) EncodeSpaces() bool { // @audit deprecate
	return c.encodeSpace
}

func (c *CommonOptions) Alphabet() *cmn.Alphabet {
	var alphabet *cmn.Alphabet

	switch strings.ToLower(c.alpha) {
	case cmn.ALPHA_NAME_ENGLISH:
		alphabet = cmn.ALPHA_DISK
		c.DefaultPhrase = "I love cryptography"

	case cmn.ALPHA_NAME_SPANISH:
		fallthrough
	case cmn.ALPHA_NAME_LATIN:
		alphabet = cmn.ALPHA_DISK_LATIN
		c.DefaultPhrase = "Amo la criptografía"

	case cmn.ALPHA_NAME_GREEK:
		alphabet = cmn.ALPHA_DISK_GREEK
		c.DefaultPhrase = "Λατρεύω την κρυπτογραφία"

	case cmn.ALPHA_NAME_GERMAN:
		alphabet = cmn.ALPHA_DISK_GERMAN
		c.DefaultPhrase = "Daß liebe hübschen Mädchen"

	case cmn.ALPHA_NAME_UKRANIAN:
		fallthrough
	case cmn.ALPHA_NAME_RUSSIAN:
		fallthrough
	case cmn.ALPHA_NAME_CYRILLIC:
		alphabet = cmn.ALPHA_DISK_CYRILLIC
		c.DefaultPhrase = "Я люблю криптографию"

	case cmn.ALPHA_NAME_BINARY:
		alphabet = cmn.BINARY_DISK

	case "custom":
		if flag.NArg() == 1 {
			alphabet = cmn.NewAlphabet("Custom", flag.Arg(0), false, false)
		} else {
			app.Die("With '-alpha custom' specify strings of characters as custom alphabet", caesarx.ERR_NO_ALPHABET)
		}

	default:
		msg := fmt.Sprintf("Valid alphabets are: %s", supportedAlphabets)
		app.Die(msg, caesarx.ERR_PARAMETER)
	}

	return alphabet
}

func (c *CommonOptions) WantsSlave() (string, bool) {
	var slaveName string = "(None)"
	var wants bool = true

	if c.numeric.IsSet {
		switch c.numeric.Value { // converted to uppercase in Validate()
		case 'A': // Arabic Numbers only
			slaveName = cmn.NUMBERS_DISK.Name

		case 'H': // Hindi Numbers only
			slaveName = cmn.NUMBERS_EASTERN_DISK.Name

		case 'E': // Arabic numbers, space and number-related chars
			slaveName = cmn.NUMBERS_DISK_EXT.Name

		case 'N':
			slaveName = "(None)"
			fallthrough

		default:
			wants = false
		}
	}

	return slaveName, wants
}

func (c *CommonOptions) Numbers() *cmn.Alphabet {
	var numerics *cmn.Alphabet = nil
	if c.numeric.IsSet {
		switch c.numeric.Value { // converted to uppercase in Validate()
		case 'A': // Arabic Numbers only
			numerics = cmn.NUMBERS_DISK.Clone()

		case 'H': // Hindi Numbers only
			numerics = cmn.NUMBERS_EASTERN_DISK.Clone()

		case 'E': // Arabic numbers, space and number-related chars
			numerics = cmn.NUMBERS_DISK_EXT.Clone()

		case 'N':

		default:
			msg := fmt.Sprintf("Valid Number tables are: %s", supportedNumbers)
			app.Die(msg, caesarx.ERR_PARAMETER)
		}
	}

	return numerics
}

func (c *CommonOptions) ShowUsage(name string) {
	fmt.Printf("\t%s [-help|-demo|-list|-version]\n", name)
	fmt.Printf("\t%s -alpha {%s}\n", name, supportedAlphabets)
	fmt.Printf("\t%s -num {%s}\n", name, supportedNumbers)
}

func (c *CommonOptions) Validate() (int, error) {
	// for these options no further validation is needed
	if c.NeedsDemo() || c.NeedsHelp() || c.NeedsList() || c.NeedsVersion() {
		c.isReady = true
		return caesarx.EXIT_CODE_SUCCESS, nil
	}

	validNumberIDs := []rune{'N', 'A', 'H', 'E'}
	c.numeric.Value = unicode.ToUpper(c.numeric.Value)

	if c.numeric.IsSet && !slices.Contains(validNumberIDs, c.numeric.Value) {
		return caesarx.ERR_CLI_OPTIONS, fmt.Errorf("-num requires any of A|H|E|N")
	}

	if !c.NeedsDemo() && !c.NeedsHelp() && !c.NeedsList() && flag.NArg() != 1 {
		app.Die("for encode/decode the free argument must be the text.", caesarx.ERR_PARAMETER)
	}

	return caesarx.EXIT_CODE_SUCCESS, nil
}

func (c *CommonOptions) IsReady() bool {
	return c.isReady
}

func (c *CommonOptions) initialize() {
	flag.BoolVar(&c.help, "help", false, "Show help")
	flag.BoolVar(&c.demo, "demo", false, "Demonstration mode")
	flag.BoolVar(&c.list, "list", false, "List all cipher variants")
	flag.BoolVar(&c.version, "version", false, "Show version number")
	flag.StringVar(&c.alpha, "alpha", defaultLanguage, "Choose alphabet")
	flag.Var(&c.numeric, "num", "Include Numbers disk: (N)one, (A)rabic, (H)indi (E)xtended")
	c.isReady = false
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/
