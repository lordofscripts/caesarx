/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Codebook console renderer and interface that opens the way to add
 * future renderers (HTML, PDF, ColorConsole).
 *-----------------------------------------------------------------*/
package main

import (
	"fmt"
	"io"
	"lordofscripts/caesarx"
	"lordofscripts/caesarx/cmn"
	"lordofscripts/caesarx/internal/sched"
	"slices"
	"strings"
	"time"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	maxWIDTH        = 80
	newLINE         = "\n"
	LEADER_LEN      = 5
	PAGE_BREAK rune = '\f'
	// https://www.unicode.org/charts/nameslist/n_2500.html
	UPLEFT      rune = rune(0x250c)
	UPRIGHT     rune = rune(0x2510)
	HORIZ       rune = rune(0x2500)
	HORIZ_HEAVY rune = rune(0x2501)
	MIDHORIZ    rune = '─'          // 0x2528 ¿?
	HORIZ_UP    rune = rune(0x2534) // ┴
	HORIZ_DOWN  rune = rune(0x252c) // ┬
	HORIZ_CROSS rune = rune(0x253c) // ┼
	VERT        rune = rune(0x2502)
	VERT_HEAVY  rune = rune(0x2503)
	DOWNLEFT    rune = rune(0x2514)
	DOWNRIGHT   rune = rune(0x2518)
	MIDLEFT     rune = '├'
	MIDRIGHT    rune = '┤'
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type ICodebookRenderer interface {
	// Renders the entire codebook for a year book. Equivalent to
	// Book head, Year page and 12 Month pages.
	RenderYearBook(recipient string)
	// Renders the Caesarium's book heading
	RenderBookHead(recipient string)
	// Renders the year's cipher schedule, specifying which cipher
	// should be used for any given month.
	RenderYearPage(ciphers []caesarx.CipherVariant, year int)
	// Renders the month's daily schedule of cipher settings
	RenderMonthPage(cipher caesarx.CipherVariant, month time.Month, year int)
	// Get the document text that has been generated
	GetDocument() string
}

var _ ICodebookRenderer = (*ConsoleCodebookRenderer)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type ConsoleCodebookRenderer struct {
	// ** User-Provided values
	date  time.Time
	alpha *cmn.Alphabet
	title string

	// ** Internal members
	sb  *strings.Builder
	hlp *sched.Caesarium
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// Create a new instance of a Codebook renderer for a Console (tty).
// Very suited for command-line applications.
// The date is adjusted so that it corresponds to January 1st of that
// year and no time component (midnight) Local time.
func NewConsoleCodebookRenderer(year uint, alpha *cmn.Alphabet, title, recovery string) *ConsoleCodebookRenderer {
	var builder strings.Builder

	dateY := time.Date(int(year), time.January, 1, 0, 0, 0, 0, time.Local)
	genZ := sched.NewCaesarium(title, alpha, dateY, 0)
	// is recoverability requested? @todo Use Bip39
	if len(recovery) != 0 {
		genZ.MakeRecoverable(recovery)
	}

	return &ConsoleCodebookRenderer{
		// ** User-Provided values
		date:  dateY,
		alpha: alpha,
		title: title,
		// ** Internal members
		sb:  &builder,
		hlp: genZ,
	}
}

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// Renders the entire codebook for a year book. Equivalent to
// Book head, Year page and 12 Month pages.
func (r *ConsoleCodebookRenderer) RenderYearBook(recipient string) {
	ciphers := r.hlp.CompileYearBook()

	r.RenderBookHead(recipient)
	r.RenderYearPage(ciphers, r.date.Year())

	r.sb.WriteRune(PAGE_BREAK)
	for m := range time.December {
		r.RenderMonthPage(ciphers[m], r.date.Month(), r.date.Year())

		// advance the month in the date object
		r.date = r.date.AddDate(0, 1, 0)
	}

	// restore out date to JANUARY as set in the ctor.
	r.date = r.date.AddDate(0, -11, 0)
}

// Renders the Caesarium's book heading
func (r *ConsoleCodebookRenderer) RenderBookHead(recipient string) {
	fmt.Fprintln(r.sb, centerString(r.title, maxWIDTH))
	fmt.Fprintln(r.sb, centerString(recipient, maxWIDTH))
	fmt.Fprintf(r.sb, "%+*s\n", maxWIDTH, time.Now().Format("2006-January-02"))
}

// Renders the year's cipher schedule, specifying which cipher
// should be used for any given month. If the ciphers parameter is
// empty or nil, we internally compile the list.
func (r *ConsoleCodebookRenderer) RenderYearPage(ciphers []caesarx.CipherVariant, year int) {
	const MONTH_WIDTH = 11
	var HEMI1 = []string{"January", "February", "March", "April", "May", "June"}
	var HEMI2 = []string{"July", "August", "September", "October", "November", "December"}
	var LINE_WIDTH = 11*6 + 2
	var Leader string = strings.Repeat(" ", LEADER_LEN)

	if len(ciphers) == 0 {
		ciphers = r.hlp.CompileYearBook()
	}

	//                    Cipher Schedule for %Year%
	// ┌──────────────────────────────────────────────────────────────────┐
	text := fmt.Sprintf("Cipher Schedule for %d", year) // @todo Localize
	fmt.Fprintln(r.sb, centerString(text, maxWIDTH))
	fmt.Fprint(r.sb, boxLine(UPLEFT, UPRIGHT, HORIZ, LEADER_LEN, LINE_WIDTH))

	// ┌──────────────────────────────────────────────────────────────────┐
	// │January    February   March      April      May        June       │
	for i, str := range HEMI1 {
		if i == 0 {
			fmt.Fprintf(r.sb, "%s%c", Leader, VERT)
		}
		fmt.Fprintf(r.sb, "%-*s", MONTH_WIDTH, str)
	}
	fmt.Fprintf(r.sb, "%c\n", VERT)

	// │Fibonacci  Caesar     None       Fibonacci  Didimus    None       │
	// ├──────────────────────────────────────────────────────────────────┤
	for i, str := range ciphers[0:6] {
		if i == 0 {
			fmt.Fprintf(r.sb, "%s%c", Leader, VERT)
		}
		fmt.Fprintf(r.sb, "%-*s", MONTH_WIDTH, str)
	}
	r.sb.WriteRune(VERT)
	r.sb.WriteRune('\n')
	r.sb.WriteString(boxLine(MIDLEFT, MIDRIGHT, MIDHORIZ, LEADER_LEN, LINE_WIDTH))

	// │July       August     September  October    November   December   │
	for i, str := range HEMI2 {
		if i == 0 {
			fmt.Fprintf(r.sb, "%s%c", Leader, VERT)
		}
		fmt.Fprintf(r.sb, "%-*s", MONTH_WIDTH, str)
	}
	fmt.Fprintf(r.sb, "%c\n", VERT)

	// │Fibonacci  Bellaso    Caesar     Bellaso    Fibonacci  Vigenere   │
	// └──────────────────────────────────────────────────────────────────┘
	for i, str := range ciphers[6:] {
		if i == 0 {
			fmt.Fprintf(r.sb, "%s%c", Leader, VERT)
		}
		fmt.Fprintf(r.sb, "%-*s", MONTH_WIDTH, str)
	}
	fmt.Fprintf(r.sb, "%c\n", VERT)

	fmt.Fprint(r.sb, boxLine(DOWNLEFT, DOWNRIGHT, HORIZ, LEADER_LEN, LINE_WIDTH))
}

// Renders the month's daily schedule of cipher settings
func (r *ConsoleCodebookRenderer) RenderMonthPage(cipher caesarx.CipherVariant, month time.Month, year int) {
	const LINE_WIDTH = 75 // from left bar to right bar excluding leading space (leader)
	// .1 Table Title
	//                       %Cipher% Daily Settings for %Month%-%Year%
	date := time.Date(r.date.Year(), month, 1, 0, 0, 0, 0, time.Local)
	tableTitle := fmt.Sprintf("%s Daily Settings for %s", cipher, date.Format("Jan-2006")) //@note Localize

	// .2 Table Header row
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │                               %Month%                                   │
	// ├─────────────────────────────────────────────────────────────────────────┤
	fmt.Fprintln(r.sb, centerString(tableTitle, maxWIDTH))
	fmt.Fprintln(r.sb, centerString(fmt.Sprintf("%s (N=%d)", r.alpha.Name, r.alpha.Size()), maxWIDTH))
	fmt.Fprint(r.sb, boxLine(UPLEFT, UPRIGHT, HORIZ, LEADER_LEN, LINE_WIDTH))
	fmt.Fprintf(r.sb, "%*s%c%s%c\n", LEADER_LEN, "", VERT, centerString(date.Month().String(), LINE_WIDTH-2), VERT)
	fmt.Fprintf(r.sb, "%s", boxLine(MIDLEFT, MIDRIGHT, HORIZ, LEADER_LEN, LINE_WIDTH))

	// .2 Generate rows
	var remnant int = LINE_WIDTH
	// .3.1 Day# DayName
	// │ 31 Mon │ %Cipher_Parameters%...
	const EMPTY_NOTE string = ""
	const DAY_PART_FORMAT_H = "%*s%c%4s%4s%c " // LEADER_LEN + 11
	const DAY_PART_FORMAT_D = "%*s%c%4d%4s%c "
	var cellDividers []int = make([]int, 0)
	var footnotes []string = make([]string, 0)

	fmt.Fprintf(r.sb, DAY_PART_FORMAT_H, LEADER_LEN, "", VERT, "Day", " ", VERT)
	remnant -= (LEADER_LEN + 11)

	lastDay := sched.LastDay(date).Day()
	switch cipher {
	// (Cipher) Affine:	[Day Part] [A] [B]	[A'] [Notes]
	case caesarx.AffineCipher:
		// │ %Day% %Weekday% │ %CoefA% %CoefB% %CoefAP% %Notes%                      │
		// ├─────────────────┼─────────────────────────┼─────────────────────────────┤
		fmt.Fprintf(r.sb, "%5s%5s%6s%s%c\n", "A", "B", "A'", centerString("Notes", remnant-12), VERT)
		cellDividers = []int{9, 30}
		fmt.Fprintf(r.sb, "%s", boxLineCell(MIDLEFT, MIDRIGHT, HORIZ, HORIZ_CROSS, cellDividers, LEADER_LEN, LINE_WIDTH))
		// │ 31 Mon │  A    B    A'   │   Notes                                      │
		affine := r.hlp.CompileAffineBook()
		for i := range lastDay {
			// %Day% %Weekday% │ ...
			weekday := date.Weekday().String()[0:3] // 1st three letters of the Weekday
			fmt.Fprintf(r.sb, DAY_PART_FORMAT_D, LEADER_LEN, "", VERT, date.Day(), weekday, VERT)
			// ... %A%  %B%  %AP% │     Notes                                        │
			fmt.Fprintf(r.sb, "%5d%5d%5d    %c%s%c\n", affine[i].A, affine[i].B, affine[i].C, VERT, centerString(EMPTY_NOTE, remnant-16), VERT)
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
		fmt.Fprintf(r.sb, "%-4s%6s%s%c\n", "Key", " Shift", centerString("Notes", remnant-6), VERT)
		fmt.Fprintf(r.sb, "%s", boxLineCell(MIDLEFT, MIDRIGHT, HORIZ, HORIZ_CROSS, cellDividers, LEADER_LEN, LINE_WIDTH))
		// │ 31 Mon │  X   23   Notes                                                │
		keysC := r.hlp.CompileCaesarBook()
		for i := range lastDay {
			// %Day% %Weekday% │
			weekday := date.Weekday().String()[0:3] // 1st three letters of the Weekday
			keyShift := keysC[i]
			fmt.Fprintf(r.sb, DAY_PART_FORMAT_D, LEADER_LEN, "", VERT, date.Day(), weekday, VERT)
			fmt.Fprintf(r.sb, "%-4c%5d  %c%s%c\n", cmn.RuneAt(r.alpha.Chars, keyShift), keyShift, VERT, centerString(EMPTY_NOTE, remnant-8), VERT)
			date = date.AddDate(0, 0, 1)
		}
		footnotes = append(footnotes, "The Shift column is the Caesar shift for the given Key")

		// (Cipher) Didimus/Fibonacci: [Day Part] [Key] [Shift] [Offset] [Notes]
	case caesarx.DidimusCipher, caesarx.FibonacciCipher:
		// │ %Day% %Weekday% │ %Key% %Shift% %Offset%   %Notes%                      │
		// ├─────────────────┼────────────────────────┼──────────────────────────────┤
		cellDividers = []int{9, 29}
		fmt.Fprintf(r.sb, "%4s%6s%8s%s%c\n", "Key", " Shift", "Offset", centerString("Notes", remnant-14), VERT)
		fmt.Fprintf(r.sb, "%s", boxLineCell(MIDLEFT, MIDRIGHT, HORIZ, HORIZ_CROSS, cellDividers, LEADER_LEN, LINE_WIDTH))
		keysDF := r.hlp.CompileBiAlphabeticBook()
		for i := range lastDay {
			// %Day% %Weekday% │
			weekday := date.Weekday().String()[0:3] // 1st three letters of the Weekday
			fmt.Fprintf(r.sb, DAY_PART_FORMAT_D, LEADER_LEN, "", VERT, date.Day(), weekday, VERT)
			fmt.Fprintf(r.sb, "%4c%5d%+7d  %c%s%c\n", cmn.RuneAt(r.alpha.Chars, keysDF[i].A), keysDF[i].A, keysDF[i].B, VERT, centerString(EMPTY_NOTE, remnant-15), VERT)
			date = date.AddDate(0, 0, 1)
		}
		footnotes = append(footnotes, "The Shift column is the Caesar shift for the given Key")
		footnotes = append(footnotes, "The Offset applies to the secondary key relative to the main Key")
		footnotes = append(footnotes, "The Offset is required for Didimus, optional for Fibonacci")

		// (Cipher) Bellaso/Vigenère: [Day Part] [Secret] [Notes]
	case caesarx.BellasoCipher, caesarx.VigenereCipher:
		// │ %Day% %Weekday% │ %Secret%                 %Notes%                      │
		// ├─────────────────┼───────────────────────┼───────────────────────────────┤
		cellDividers = []int{9, 38}
		fmt.Fprintf(r.sb, "%-16s%s%c\n", "Secret", centerString("Notes", remnant-12), VERT)
		fmt.Fprintf(r.sb, "%s", boxLineCell(MIDLEFT, MIDRIGHT, HORIZ, HORIZ_CROSS, cellDividers, LEADER_LEN, LINE_WIDTH))
		keysBV := r.hlp.CompileWordBook(26)
		for i := range lastDay {
			// %Day% %Weekday% │
			weekday := date.Weekday().String()[0:3] // 1st three letters of the Weekday
			fmt.Fprintf(r.sb, DAY_PART_FORMAT_D, LEADER_LEN, "", VERT, date.Day(), weekday, VERT)
			fmt.Fprintf(r.sb, "%-27s%c%s%c\n", keysBV[i], VERT, centerString(EMPTY_NOTE, remnant-24), VERT)
			date = date.AddDate(0, 0, 1)
		}
		footnotes = append(footnotes, "For Bellaso the Secret is repeated over the input")
		footnotes = append(footnotes, "For Vigenere Auto-key the Secret is only used once")
	}

	// .4 Box footer
	addFootnote := func(fd io.Writer, format string, args ...any) {
		fmt.Fprintf(fd, "%*s· ", LEADER_LEN, "")
		fmt.Fprintf(fd, format, args...)
		fmt.Fprint(fd, "\n")
	}

	// draw the bottom line of the table with the cell dividers every cipher table needs
	lenRunes, lenBytes := r.alpha.SizeExt()
	fmt.Fprint(r.sb, boxLineCell(DOWNLEFT, DOWNRIGHT, HORIZ, HORIZ_UP, cellDividers, LEADER_LEN, LINE_WIDTH))

	// add general and cipher-specific footnotes
	addFootnote(r.sb, "Alphabet (%2s) has %d runes and %d bytes", r.alpha.LangCodeISO(), lenRunes, lenBytes)
	addFootnote(r.sb, "Alphabet Runes: %s", r.alpha.Chars)
	if len(footnotes) > 0 {
		for _, footnote := range footnotes {
			addFootnote(r.sb, footnote)
		}
	}

	r.sb.WriteRune(PAGE_BREAK)
}

// Get the document text that has been generated
func (r *ConsoleCodebookRenderer) GetDocument() string {
	return r.sb.String()
}

/* ----------------------------------------------------------------
 *				P r i v a t e		M e t h o d s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// centers a string in the width
func centerString(s string, width int) string {
	if len(s) >= width {
		return s // Return original string if it's longer than the width
	}

	totalSpaces := width - len(s)
	leftSpaces := totalSpaces / 2
	rightSpaces := totalSpaces - leftSpaces

	return strings.Repeat(" ", leftSpaces) + s + strings.Repeat(" ", rightSpaces)
}

// Uses Unicode Box Drawing characters to draw a multi-functional line using
// the left and right corners and a middle character. This line can be used
// at the top, middle or bottom depending on the chosen left, right & middle
// runes. The total line length is given and how many leader spaces to use
// (or zero if none). The line is terminated with a newline.
func boxLine(left, right rune, middle rune, leader, total int) string {
	return fmt.Sprintf("%*s%c%s%c\n", leader, "", left, strings.Repeat(string(middle), total-2), right)
}

// Same as boxLine() with the addition of a character/rune to use mid-line
// (such as horizontal up/down/both) and their 0-based positions.
func boxLineCell(left, right rune, middle rune, atChar rune, atPos []int, leader, total int) string {
	midline := make([]rune, total)
	for pos := range total {
		switch pos {
		case 0:
			midline[0] = left
		case total - 1:
			midline[pos] = right
		default:
			if slices.Contains(atPos, pos) {
				midline[pos] = atChar
			} else {
				midline[pos] = middle
			}
		}
	}

	return strings.Repeat(" ", leader) + string(midline) + newLINE
}
