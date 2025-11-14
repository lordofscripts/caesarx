/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Command-Line option processor for flags that are common to all
 * the cmd/* applications. However, it allows the possibility to
 * exclude the flags on an individual basis by passing the black list
 * in the constructor.
 *-----------------------------------------------------------------*/
package cmd

import (
	"flag"
	"fmt"
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmn"
	"slices"
	"strings"
	"unicode"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	OPT_SLAVE_NONE     rune = 'N'
	OPT_SLAVE_ARABIC   rune = 'A' // id:PSO_NUM_DEC name:ALPHA_NAME_NUMBERS_ARABIC
	OPT_SLAVE_HINDI    rune = 'H' // id:PSO_NUM_HIN name:ALPHA_NAME_NUMBERS_EASTERN
	OPT_SLAVE_EXTENDED rune = 'E' // id:PSO_NUM_DEC_EXT name:ALPHA_NAME_NUMBERS_ARABIC_EXTENDED
	OPT_SLAVE_PUNCT    rune = 'P' // id:PSO_PUNCT name:ALPHA_NAME_PUNCTUATION
	OPT_SLAVE_SYMBL    rune = 'S' // id:PSO_PUNCT_DEC name:ALPHA_NAME_SYMBOLS

	defaultLanguage    string = "english"
	supportedAlphabets string = "english|latin|spanish|german|greek|cyrillic|italian|portuguese|czech|custom|binary"
	supportedNumbers   string = "(N)one (A)rabic (E)xtended (H)indi"
)

const (
	// Common CLI flags. Each may be excluded on an app basis
	FLAG_HELP    string = "help"
	FLAG_DEMO    string = "demo"
	FLAG_LIST    string = "list"
	FLAG_VERSION string = "version"
	FLAG_ALPHA   string = "alpha"
	FLAG_NUM     string = "num"
	FLAG_PROFILE string = "profile" // (optional) Select profile
)

var AppConfig = NewConfiguration()

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
	FileExt() string
}

var _ IAppOptions = (*CommonOptions)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type CommonOptions struct {
	DefaultPhrase string
	optProfile    string
	demo          bool
	help          bool
	list          bool
	version       bool
	encodeSpace   bool
	alpha         string
	numeric       RuneFlag
	isReady       bool
}

type FileOptions struct {
	Input  string
	Output string
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * Register options common to all applications/tools and provide
 * methods to access them.
 * NOTE: Requires flag.Parse()
 */
func NewCommonOptions(skipFlags ...string) *CommonOptions {
	copts := &CommonOptions{}
	copts.optProfile = ""
	copts.DefaultPhrase = "Let's encrypt!"
	copts.initialize(skipFlags...)
	AppConfig.InitConfiguration()
	return copts
}

func NewFileOptions(inp, out string) *FileOptions {
	return &FileOptions{inp, out}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// The -help CLI flag is given
func (c *CommonOptions) NeedsHelp() bool {
	return c.help
}

// The -demo CLI flag is given
func (c *CommonOptions) NeedsDemo() bool {
	return c.demo
}

// The -list CLI flag is given
func (c *CommonOptions) NeedsList() bool {
	return c.list
}

// the -version CLI flag is given
func (c *CommonOptions) NeedsVersion() bool {
	return c.version
}

// The -profile CLI flag is given.
// The user requests encryption presets from a user profile.
// See also GetProfileID()
func (c *CommonOptions) RequestsProfile() bool {
	return len(c.optProfile) != 0
}

// The requested profile ID if -profile was given
func (c *CommonOptions) GetRequestedProfile() string {
	return c.optProfile
}

func (c *CommonOptions) EncodeSpaces() bool { // @audit deprecate
	return c.encodeSpace
}

// get the built-in alphabet represented by the CLI option -alpha
func (c *CommonOptions) Alphabet() *cmn.Alphabet {
	alpha, phrase := SelectAlphabet(c.alpha)
	c.DefaultPhrase = phrase
	return alpha
}

// the selected alphabet is Binary
func (c *CommonOptions) IsBinary() bool {
	return strings.ToLower(c.alpha) == cmn.ALPHA_NAME_BINARY
}

func (c *CommonOptions) WantsSlave() (string, bool) {
	var slaveName string = "(None)"
	var wants bool = true

	if c.numeric.IsSet {
		switch c.numeric.Value { // converted to uppercase in Validate()
		case OPT_SLAVE_ARABIC: // Arabic Numbers only
			slaveName = cmn.NUMBERS_DISK.Name

		case OPT_SLAVE_HINDI: // Hindi Numbers only
			slaveName = cmn.NUMBERS_EASTERN_DISK.Name

		case OPT_SLAVE_EXTENDED: // Arabic numbers, space and number-related chars
			slaveName = cmn.NUMBERS_DISK_EXT.Name

		case OPT_SLAVE_PUNCT: // Punctuation according to UTF8
			slaveName = cmn.PUNCTUATION_DISK.Name

		case OPT_SLAVE_SYMBL: // Symbols according to UTF8
			slaveName = cmn.SYMBOL_DISK.Name

		case OPT_SLAVE_NONE:
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
		case OPT_SLAVE_ARABIC: // Arabic Numbers only
			numerics = cmn.NUMBERS_DISK.Clone()

		case OPT_SLAVE_HINDI: // Hindi Numbers only
			numerics = cmn.NUMBERS_EASTERN_DISK.Clone()

		case OPT_SLAVE_EXTENDED: // Arabic numbers, space and number-related chars
			numerics = cmn.NUMBERS_DISK_EXT.Clone()

		case OPT_SLAVE_PUNCT:
			numerics = cmn.PUNCTUATION_DISK.Clone()

		case OPT_SLAVE_SYMBL:
			numerics = cmn.SYMBOL_DISK.Clone()

		case OPT_SLAVE_NONE:

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

// FileExt of CommonOptions returns an empty string
func (c *CommonOptions) FileExt() string {
	return ""
}

// Validates -num option to any of N|A|E|H|P|S
func (c *CommonOptions) Validate() (int, error) {
	// for these options no further validation is needed
	if c.NeedsDemo() || c.NeedsHelp() || c.NeedsList() || c.NeedsVersion() {
		c.isReady = true
		return caesarx.EXIT_CODE_SUCCESS, nil
	}

	validNumberIDs := []rune{
		OPT_SLAVE_NONE,
		OPT_SLAVE_ARABIC,
		OPT_SLAVE_HINDI,
		OPT_SLAVE_EXTENDED,
		OPT_SLAVE_PUNCT,
		OPT_SLAVE_SYMBL,
	}
	c.numeric.Value = unicode.ToUpper(c.numeric.Value)

	if c.numeric.IsSet && !slices.Contains(validNumberIDs, c.numeric.Value) {
		return caesarx.ERR_CLI_OPTIONS, fmt.Errorf("-num requires any of A|H|E|N")
	}

	return caesarx.EXIT_CODE_SUCCESS, nil
}

func (c *CommonOptions) IsReady() bool {
	return c.isReady
}

// Used for presetting the primary alphabet (from Profile config).
// The name must be any of the ALPHA_NAME_*
func (c *CommonOptions) PresetPrimaryAlphabet(name string) {
	c.alpha = name
}

// Used for presetting the secondary alphabet (from Profile config).
// The name must be any of the
func (c *CommonOptions) PresetSecondaryAlphabet(name string) {
	c.numeric.IsSet = true
	switch name { // converted to uppercase in Validate()
	case cmn.ALPHA_NAME_NUMBERS_ARABIC: // Arabic Numbers only
		c.numeric.Value = OPT_SLAVE_ARABIC

	case cmn.ALPHA_NAME_NUMBERS_EASTERN: // Hindi Numbers only
		c.numeric.Value = OPT_SLAVE_HINDI

	case cmn.ALPHA_NAME_NUMBERS_ARABIC_EXTENDED: // Arabic numbers, space and number-related chars
		c.numeric.Value = OPT_SLAVE_EXTENDED

	case cmn.ALPHA_NAME_PUNCTUATION:
		c.numeric.Value = OPT_SLAVE_PUNCT

	case cmn.ALPHA_NAME_SYMBOLS:
		c.numeric.Value = OPT_SLAVE_SYMBL

	case "":
		c.numeric.Value = OPT_SLAVE_NONE

	default:
		c.numeric.IsSet = false
		mlog.Error("unable to preset slave alphabet", mlog.String("Name", name), mlog.At())
	}
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

// initializes common options by registering the CLI flags that are
// not present in the skip list
func (c *CommonOptions) initialize(skipFlags ...string) {
	// user-configuration overrides: Primary alphabet
	defLang := defaultLanguage
	if AppConfig.IsGood() {
		defLang = AppConfig.Configuration.Defaults.AlphaName
	}

	// register Common flags that are NOT in the skip list
	if skipFlags == nil {
		skipFlags = make([]string, 0)
	}

	if !slices.Contains(skipFlags, FLAG_HELP) {
		flag.BoolVar(&c.help, FLAG_HELP, false, "Show help")
	}
	if !slices.Contains(skipFlags, FLAG_DEMO) {
		flag.BoolVar(&c.demo, FLAG_DEMO, false, "Demonstration mode")
	}
	if !slices.Contains(skipFlags, FLAG_LIST) {
		flag.BoolVar(&c.list, FLAG_LIST, false, "List all cipher variants")
	}
	if !slices.Contains(skipFlags, FLAG_VERSION) {
		flag.BoolVar(&c.version, FLAG_VERSION, false, "Show version number")
	}
	if !slices.Contains(skipFlags, FLAG_ALPHA) {
		flag.StringVar(&c.alpha, FLAG_ALPHA, defLang, "Choose alphabet")
	}
	if !slices.Contains(skipFlags, FLAG_NUM) {
		flag.Var(&c.numeric, FLAG_NUM, "Include Numbers disk: (N)one, (A)rabic, (H)indi (E)xtended")
	}
	if !slices.Contains(skipFlags, FLAG_PROFILE) {
		flag.StringVar(&c.optProfile, FLAG_PROFILE, "", "Profile selector for cipher presets")
	}

	c.isReady = false
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func SelectAlphabet(name string) (*cmn.Alphabet, string) {
	var alphabet *cmn.Alphabet
	var phrase string

	switch strings.ToLower(name) {
	case cmn.ALPHA_NAME_ENGLISH:
		alphabet = cmn.ALPHA_DISK
		phrase = "I love cryptography"

	case cmn.ALPHA_NAME_SPANISH:
		fallthrough
	case cmn.ALPHA_NAME_LATIN:
		alphabet = cmn.ALPHA_DISK_LATIN
		phrase = "Amo la criptografía"

	case cmn.ALPHA_NAME_GREEK:
		alphabet = cmn.ALPHA_DISK_GREEK
		phrase = "Λατρεύω την κρυπτογραφία"

	case cmn.ALPHA_NAME_GERMAN:
		alphabet = cmn.ALPHA_DISK_GERMAN
		phrase = "Daß liebe hübschen Mädchen"

	case cmn.ALPHA_NAME_UKRAINIAN:
		fallthrough
	case cmn.ALPHA_NAME_RUSSIAN:
		fallthrough
	case cmn.ALPHA_NAME_CYRILLIC:
		alphabet = cmn.ALPHA_DISK_CYRILLIC
		phrase = "Я люблю криптографию"

	case cmn.ALPHA_NAME_ITALIAN:
		alphabet = cmn.ALPHA_DISK_ITALIAN
		phrase = "Amo la crittografia"

	case cmn.ALPHA_NAME_PORTUGUESE:
		alphabet = cmn.ALPHA_DISK_PORTUGUESE
		phrase = "Eu amo criptografia"

	case cmn.ALPHA_NAME_CZECH:
		alphabet = cmn.ALPHA_DISK_CZECH
		phrase = "Miluji kryptografii"

	case cmn.ALPHA_NAME_BINARY:
		alphabet = cmn.BINARY_DISK
		phrase = "love ántaño Λατρ Daß люблю"

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

	return alphabet, phrase
}
