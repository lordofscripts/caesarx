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
	"lordofscripts/caesarx/ciphers"
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

/*
// Renders the entire year book comprising title, cipher schedule for every month and the
// corresponding month daily schedule for every designated cipher.
func RenderYearBook(sb *strings.Builder, date time.Time, title, recipient string, alpha *cmn.Alphabet, hlp *sched.Caesarium) {
	if date.Month() != time.January {
		date = date.AddDate(0, -1*int(date.Month()), 0) // make it January
	}

	RenderBookHead(sb, title, recipient, time.Now())
	ciphers := hlp.CompileYearBook()
	RenderYearBookPage(sb, ciphers, date.Year())
	sb.WriteRune(PAGE_BREAK)

	for m := range 12 {
		RenderMonthPage(sb, ciphers[m], date, alpha, hlp)
		sb.WriteRune(PAGE_BREAK)
		// advance the month in the date object
		date = date.AddDate(0, 1, 0)
	}
}

// Renders the Caesarium's book heading
func RenderBookHead(sb *strings.Builder, title, recipient string, genDate time.Time) {
	sb.WriteString(centerString(title, maxWIDTH) + newLINE)
	sb.WriteString(centerString(recipient, maxWIDTH) + newLINE)
	sb.WriteString(fmt.Sprintf("%*s", maxWIDTH-8, genDate.Format("2006-January-02")) + newLINE)
}

// Renders a Caesarium's year schedule of ciphers, one random cipher for every month
func RenderYearBookPage(sb *strings.Builder, ciphers []caesarx.CipherVariant, year int) {
	const MONTH_WIDTH = 11
	var HEMI1 = []string{"January", "February", "March", "April", "May", "June"}
	var HEMI2 = []string{"July", "August", "September", "October", "November", "December"}
	var LINE_WIDTH = 11*6 + 2

	//                    Cipher Schedule for %Year%
	// ┌──────────────────────────────────────────────────────────────────┐
	text := fmt.Sprintf("Cipher Schedule for %d", year)
	sb.WriteString(centerString(text, maxWIDTH) + newLINE)
	sb.WriteString(boxLine(UPLEFT, UPRIGHT, HORIZ, LEADER_LEN, LINE_WIDTH))

	// ┌──────────────────────────────────────────────────────────────────┐
	// │January    February   March      April      May        June       │
	for i, str := range HEMI1 {
		if i == 0 {
			sb.WriteString(strings.Repeat(" ", LEADER_LEN))
			sb.WriteRune(VERT)
		}
		fmt.Fprintf(sb, "%-*s", MONTH_WIDTH, str)
	}
	sb.WriteRune(VERT)
	sb.WriteRune('\n')

	// │Fibonacci  Caesar     None       Fibonacci  Didimus    None       │
	// ├──────────────────────────────────────────────────────────────────┤
	for i, str := range ciphers[0:6] {
		if i == 0 {
			sb.WriteString(strings.Repeat(" ", LEADER_LEN))
			sb.WriteRune(VERT)
		}
		fmt.Fprintf(sb, "%-*s", MONTH_WIDTH, str)
	}
	sb.WriteRune(VERT)
	sb.WriteRune('\n')
	sb.WriteString(boxLine(MIDLEFT, MIDRIGHT, MIDHORIZ, LEADER_LEN, LINE_WIDTH))

	// │July       August     September  October    November   December   │
	for i, str := range HEMI2 {
		if i == 0 {
			sb.WriteString(strings.Repeat(" ", LEADER_LEN))
			sb.WriteRune(VERT)
		}
		fmt.Fprintf(sb, "%-*s", MONTH_WIDTH, str)
	}
	sb.WriteRune(VERT)
	sb.WriteRune('\n')

	// │Fibonacci  Bellaso    Caesar     Bellaso    Fibonacci  Vigenere   │
	// └──────────────────────────────────────────────────────────────────┘
	for i, str := range ciphers[6:] {
		if i == 0 {
			sb.WriteString(strings.Repeat(" ", LEADER_LEN))
			sb.WriteRune(VERT)
		}
		fmt.Fprintf(sb, "%-*s", MONTH_WIDTH, str)
	}
	sb.WriteRune(VERT)
	sb.WriteRune('\n')

	sb.WriteString(boxLine(DOWNLEFT, DOWNRIGHT, HORIZ, LEADER_LEN, LINE_WIDTH))
}

// Renders a page with a daily schedule of parameters for the named month and
// selected cipher.
func RenderMonthPage(sb *strings.Builder, cipher caesarx.CipherVariant, date time.Time, alpha *cmn.Alphabet, hlp *sched.Caesarium) {
	const LINE_WIDTH = 75 // from left bar to right bar excluding leading space (leader)
	// .1 Table Title
	//                       %Cipher% Daily Settings for %Month%-%Year%
	tableTitle := fmt.Sprintf("%s Daily Settings for %s", cipher, date.Format("Jan-2006"))

	// .2 Table Header row
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │                               %Month%                                   │
	// ├─────────────────────────────────────────────────────────────────────────┤
	//var NOTES_LEN = LINE_WIDTH - 19
	sb.WriteString(centerString(tableTitle, maxWIDTH) + newLINE)
	fmt.Fprintln(sb, centerString(fmt.Sprintf("%s (N=%d)", alpha.Name, alpha.Size()), maxWIDTH))
	sb.WriteString(boxLine(UPLEFT, UPRIGHT, HORIZ, LEADER_LEN, LINE_WIDTH))
	fmt.Fprintf(sb, "%*s%c%s%c\n", LEADER_LEN, "", VERT, centerString(date.Month().String(), LINE_WIDTH-2), VERT)
	fmt.Fprintf(sb, "%s", boxLine(MIDLEFT, MIDRIGHT, HORIZ, LEADER_LEN, LINE_WIDTH))

	// .2 Generate rows
	var remnant int = LINE_WIDTH
	// .3.1 Day# DayName
	// │ 31 Mon │ %Cipher_Parameters%...
	const EMPTY_NOTE string = ""
	const DAY_PART_FORMAT_H = "%*s%c%4s%4s%c " // LEADER_LEN + 11
	const DAY_PART_FORMAT_D = "%*s%c%4d%4s%c "
	var cellDividers []int = make([]int, 0)
	var footnotes []string = make([]string, 0)

	fmt.Fprintf(sb, DAY_PART_FORMAT_H, LEADER_LEN, "", VERT, "Day", " ", VERT)
	remnant -= (LEADER_LEN + 11)

	lastDay := sched.LastDay(date).Day()
	switch cipher {
	// (Cipher) Affine:	[Day Part] [A] [B]	[A'] [Notes]
	case caesarx.AffineCipher:
		// │ %Day% %Weekday% │ %CoefA% %CoefB% %CoefAP% %Notes%                      │
		// ├─────────────────┼─────────────────────────┼─────────────────────────────┤
		fmt.Fprintf(sb, "%5s%5s%6s%s%c\n", "A", "B", "A'", centerString("Notes", remnant-12), VERT)
		cellDividers = []int{9, 30}
		fmt.Fprintf(sb, "%s", boxLineCell(MIDLEFT, MIDRIGHT, HORIZ, HORIZ_CROSS, cellDividers, LEADER_LEN, LINE_WIDTH))
		// │ 31 Mon │  A    B    A'   │   Notes                                      │
		affine := hlp.CompileAffineBook()
		for i := range lastDay {
			// %Day% %Weekday% │ ...
			weekday := date.Weekday().String()[0:3] // 1st three letters of the Weekday
			fmt.Fprintf(sb, DAY_PART_FORMAT_D, LEADER_LEN, "", VERT, date.Day(), weekday, VERT)
			// ... %A%  %B%  %AP% │     Notes                                        │
			fmt.Fprintf(sb, "%5d%5d%5d    %c%s%c\n", affine[i].A, affine[i].B, affine[i].C, VERT, centerString(EMPTY_NOTE, remnant-16), VERT)
			date = date.AddDate(0, 0, 1)
		}
		footnotes = append(footnotes, "Affine A coefficient used during encryption")
		footnotes = append(footnotes, "Affine A' coefficient used during decryption")
		footnotes = append(footnotes, "Affine B used for both encryption & decryption")

		// (Cipher) Caesar:	[Day Part] [Key] [Shift] [Notes]
	case caesarx.CaesarCipher:
		// │ %Day% %Weekday% │ %Key% %Shift%            %Notes%                      │
		// ├─────────────────┼─────────────────┼─────────────────────────────────────┤
		cellDividers = []int{9, 22}
		fmt.Fprintf(sb, "%-4s%6s%s%c\n", "Key", " Shift", centerString("Notes", remnant-6), VERT)
		fmt.Fprintf(sb, "%s", boxLineCell(MIDLEFT, MIDRIGHT, HORIZ, HORIZ_CROSS, cellDividers, LEADER_LEN, LINE_WIDTH))
		// │ 31 Mon │  X   23   Notes                                                │
		keysC := hlp.CompileCaesarBook()
		for i := range lastDay {
			// %Day% %Weekday% │
			weekday := date.Weekday().String()[0:3] // 1st three letters of the Weekday
			keyShift := keysC[i]
			fmt.Fprintf(sb, DAY_PART_FORMAT_D, LEADER_LEN, "", VERT, date.Day(), weekday, VERT)
			fmt.Fprintf(sb, "%-4c%5d  %c%s%c\n", cmn.RuneAt(alpha.Chars, keyShift), keyShift, VERT, centerString(EMPTY_NOTE, remnant-8), VERT)
			date = date.AddDate(0, 0, 1)
		}

		// (Cipher) Didimus/Fibonacci: [Day Part] [Key] [Shift] [Offset] [Notes]
	case caesarx.DidimusCipher, caesarx.FibonacciCipher:
		// │ %Day% %Weekday% │ %Key% %Shift% %Offset%   %Notes%                      │
		// ├─────────────────┼────────────────────────┼──────────────────────────────┤
		cellDividers = []int{9, 29}
		fmt.Fprintf(sb, "%4s%6s%8s%s%c\n", "Key", " Shift", "Offset", centerString("Notes", remnant-14), VERT)
		fmt.Fprintf(sb, "%s", boxLineCell(MIDLEFT, MIDRIGHT, HORIZ, HORIZ_CROSS, cellDividers, LEADER_LEN, LINE_WIDTH))
		keysDF := hlp.CompileBiAlphabeticBook()
		for i := range lastDay {
			// %Day% %Weekday% │
			weekday := date.Weekday().String()[0:3] // 1st three letters of the Weekday
			fmt.Fprintf(sb, DAY_PART_FORMAT_D, LEADER_LEN, "", VERT, date.Day(), weekday, VERT)
			fmt.Fprintf(sb, "%4c%5d%+7d  %c%s%c\n", cmn.RuneAt(alpha.Chars, keysDF[i].A), keysDF[i].A, keysDF[i].B, VERT, centerString(EMPTY_NOTE, remnant-15), VERT)
			date = date.AddDate(0, 0, 1)
		}

		// (Cipher) Bellaso/Vigenère: [Day Part] [Secret] [Notes]
	case caesarx.BellasoCipher, caesarx.VigenereCipher:
		// │ %Day% %Weekday% │ %Secret%                 %Notes%                      │
		// ├─────────────────┼───────────────────────┼───────────────────────────────┤
		cellDividers = []int{9, 38}
		fmt.Fprintf(sb, "%-16s%s%c\n", "Secret", centerString("Notes", remnant-12), VERT)
		fmt.Fprintf(sb, "%s", boxLineCell(MIDLEFT, MIDRIGHT, HORIZ, HORIZ_CROSS, cellDividers, LEADER_LEN, LINE_WIDTH))

	}

	// .4 Box footer
	addFootnote := func(fd io.Writer, format string, args ...any) {
		fmt.Fprintf(fd, "%*s· ", LEADER_LEN, "")
		fmt.Fprintf(fd, format, args...)
		fmt.Fprintln(fd)
	}

	// draw the bottom line of the table with the cell dividers every cipher table needs
	lenRunes, lenBytes := alpha.SizeExt()
	fmt.Fprint(sb, boxLineCell(DOWNLEFT, DOWNRIGHT, HORIZ, HORIZ_UP, cellDividers, LEADER_LEN, LINE_WIDTH))

	// add general and cipher-specific footnotes
	addFootnote(sb, "Alphabet (%2s) has %d runes and %d bytes", alpha.LangCodeISO(), lenRunes, lenBytes)
	addFootnote(sb, "Alphabet Runes: %s", alpha.Chars)
	if len(footnotes) > 0 {
		for _, footnote := range footnotes {
			addFootnote(sb, footnote)
		}
	}
}
*/

// Help about using this
func Help() {
	flag.Usage()
	flag.PrintDefaults()
	fmt.Println(ciphers.PrintAvailableCiphers()) // @audit the init() functions of ciphers are not getting called6se
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
	var flgTitle, flgVariant, flgAlphabet, flgRecovery, flgRecipient, flgOutFormat string

	flgOutFormat = OUT_TEXT_PLAIN
	flag.BoolVar(&flgHelp, "help", false, "This help")
	flag.BoolVar(&flgFullBook, "full", false, "Produce entire book when date is a year alone")
	flag.BoolVar(&flgBip39, "bip39", false, "Generate BIP39 mnemonic (use i.s.o. -recovery)")
	flag.StringVar(&flgTitle, "title", "Caesarium", "Codebook Title")
	flag.StringVar(&flgAlphabet, "alpha", "english", "Primary alphabet name")
	flag.StringVar(&flgVariant, "variant", "caesar", "Cipher variant")
	flag.StringVar(&flgRecovery, "recovery", "", "Recovery phrase to generate a recoverable Caesarium (else use -bip39)")
	flag.StringVar(&flgRecipient, "for", "you@bitbucket.com", "The recipient of messages from this codebook")
	flag.Var(flgDate, "date", "now|today|ahora|hoy| date format such as 2006-01-02")
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
	if flgBip39 && len(flgRecovery) != 0 {
		app.Die("options -bip39 and -recovery are mutually exclusive", caesarx.ERR_PARAMETER)
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
	} else {
		mnemonics = flgRecovery
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

	if needYearBook {
		if flgFullBook {
			renderer.RenderYearBook(flgRecipient)
		} else {
			renderer.RenderYearPage(nil, bookDate.Year())
		}
	}

	if needMonthBook {
		renderer.RenderBookHead(flgRecipient)
		renderer.RenderMonthPage(selectedCipher, bookDate.Month(), bookDate.Year())
	}

	fmt.Println(renderer.GetDocument())

	caesarx.BuyMeCoffee()
}
