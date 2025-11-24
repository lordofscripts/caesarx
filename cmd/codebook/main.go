/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Codebook tool for CaesarX supported ciphers.
 *-----------------------------------------------------------------*/
package main

import (
	"flag"
	"fmt"
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmd"
	"lordofscripts/caesarx/internal/bip39"
	"os"
	"strings"
	"time"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	APP_VERSION string = "1.0"

	OUT_TEXT_PLAIN string = "text"
	OUT_TEXT_HTML  string = "html"
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/
func init() {
	caesarx.Copyright(caesarx.CO1, true)
	caesarx.BuyMeCoffee()
	fmt.Println("\t=========================================")
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// sample Command-line variations (a few)
func Usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s\n", os.Args[0])
	fmt.Println("\tcodebook [OPTIONS] -date 2025 -full")
	fmt.Println("\tcodebook [OPTIONS] -date 2025-12")
	fmt.Println("\tcodebook [OPTIONS] -date today -bip39")
}

// Help about using this
func Help() {
	flag.Usage()
	fmt.Println("Options:")
	flag.PrintDefaults()
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/
func main() {
	//var err error
	defer mlog.CloseLogFiles()

	// -------	CLI FLAGS ------
	var flgHelp, flgFullBook, flgBip39 bool
	var flgDate *cmd.DateFlag = cmd.NewDateVar("2006-01", "2006", "2006-Jan")
	var flgTitle, flgVariant, flgAlphabet, flgRecipient, flgOutFormat string

	flgOutFormat = OUT_TEXT_PLAIN
	flag.Usage = Usage
	flag.BoolVar(&flgHelp, "help", false, "This help")
	flag.BoolVar(&flgFullBook, "full", false, "Produce entire book when date is a year alone")
	flag.BoolVar(&flgBip39, "bip39", false, "Generate BIP39 mnemonic for a recoverable Caesarium")
	flag.StringVar(&flgTitle, "title", "Caesarium", "Codebook Title")
	flag.StringVar(&flgAlphabet, "alpha", "english", "Primary alphabet name")
	flag.StringVar(&flgVariant, "variant", "caesar", "Cipher variant")
	flag.StringVar(&flgRecipient, "for", "you@bitbucket.com", "The recipient of messages from this codebook")
	flag.Var(flgDate, "date", "now|today|ahora|hoy| date format such as 2006-01-02 or 2025-12")
	flag.Parse()

	// -------	CLI VALIDATION ------
	// .1 date (year or month & year)
	var bookDate time.Time
	if flgDate.IsSet {
		bookDate = flgDate.Value
	} else {
		bookDate = cmd.ThisMonth()
	}
	// .2 cipher variant
	selectedCipher, err := caesarx.NoCipher.Parse(flgVariant)
	if err != nil {
		app.DieWithError(err, caesarx.ERR_BAD_CIPHER)
	}
	// .3 Mnemonic recovery
	if flgBip39 {
		println("Generating a Recoverable Caesarium codebook...")
	}

	// .4 Recovery & Mnemonic
	var mnemonics string = ""
	if flgBip39 { // use BIP39 as recovery phrase instead of user-provided
		bip := bip39.NewBip39(bip39.Bip39Words12, ' ')
		if words, err := bip.GenerateMnemonic(); err != nil {
			app.DieWithError(err, caesarx.ERR_INTERNAL)
		} else {
			mnemonics = strings.Join(words, " ")
		}
	}

	// .3 retrieve the built-in alphabet requested by the user
	alphabet, _ := cmd.SelectAlphabet(flgAlphabet)
	// .4 books to be generated
	dtWithYear, dtWithMonth, _ := flgDate.Has()
	needYearBook := dtWithYear && !dtWithMonth
	needMonthBook := dtWithYear && dtWithMonth

	// -------	EXECUTION ------
	// .1 terminal options
	if flgHelp {
		Help()
		os.Exit(0)
	}

	// https://patorjk.com/software/taag/#p=display&f=Future&t=Caesarium&x=cppComment&v=4&h=4&w=80&we=false
	fmt.Println("    ┏━╸┏━┓┏━╸┏━┓┏━┓┏━┓╻╻ ╻┏┳┓")
	fmt.Println("    ┃  ┣━┫┣╸ ┗━┓┣━┫┣┳┛┃┃ ┃┃┃┃")
	fmt.Println("    ┗━╸╹ ╹┗━╸┗━┛╹ ╹╹┗╸╹┗━┛╹ ╹")
	fmt.Println("    By Lord-of-Scripts™")

	// .2 create a Caesarium engine with the selected renderer
	var renderer ICodebookRenderer = nil
	switch strings.ToLower(flgOutFormat) {
	case OUT_TEXT_PLAIN:
		renderer = NewConsoleCodebookRenderer(uint(bookDate.Year()), alphabet, flgTitle, mnemonics)

	case OUT_TEXT_HTML:

	}

	// .4 prepare a full year's Caesarium. An entire 13-page codebook, or just the month's cipher name
	if needYearBook {
		if flgFullBook {
			renderer.RenderYearBook(flgRecipient)
		} else {
			renderer.RenderYearPage(nil, bookDate.Year())
		}
	}

	// .5 prepare a one-month Caesarium
	if needMonthBook {
		renderer.RenderBookHead(flgRecipient)
		renderer.RenderMonthPage(selectedCipher, bookDate.Month(), bookDate.Year())
	}

	// .6 Render it with the selected format
	fmt.Println(renderer.GetDocument())

	caesarx.BuyMeCoffee()
}
