/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Basic Caesar Cipher command-line application with in-lined (appended)
 * alphabets rather than using the chain mechanism.
 *-----------------------------------------------------------------*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	z "lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/cmd"
	"lordofscripts/caesarx/cmn"
	"os"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	APP_VERSION     string = "1.0"
	NGRAM_SEPARATOR rune   = '·'
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/
func init() {
	if !app.IsPipedInput() {
		z.Copyright(z.CO1, true)
		z.BuyMeCoffee()
		fmt.Println("\t=========================================")
	}
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Encrypt using plain Caesar cipher with the given alphabet and key.
func DoEncrypt(alpha *cmn.Alphabet, key rune, nGramSize int) (int, error) {
	var exitCode int = z.EXIT_CODE_SUCCESS
	// These commands implement IPipe, ICipherCommand & ICommand
	cmdCipher := commands.NewCaesarCommand(alpha, key)

	var cipher string
	var err error
	if app.IsPipedInput() {
		reader := bufio.NewReader(os.Stdin)
		scanner := bufio.NewScanner(reader)
		var lineIn string
		for scanner.Scan() {
			lineIn = scanner.Text()
			cipher, err = cmdCipher.Encode(lineIn)
			if err != nil {
				exitCode = z.ERR_CIPHER
				mlog.ErrorE(err)
				break
			}
			fmt.Println(cipher)
		}

		if err = scanner.Err(); err != nil {
			exitCode = z.ERR_FILE_IO
			mlog.ErrorE(err)
		}
	} else {
		plain := flag.Arg(0)
		cipher, err = cmdCipher.Encode(plain)

		if err == nil {
			fmt.Println("Alphabet : ", alpha.Chars)
			fmt.Println("Key      : ", key)
			fmt.Println("Plain    : ", plain)
			fmt.Println("Encoded  : ", cipher)
		} else {
			exitCode = z.ERR_CIPHER
		}

		// If NGram formatting wanted, create it as Pipe command, only for Encoding
		var ngram string = ""
		if err == nil && nGramSize > 0 {
			ngramCmd := cmn.NewNgramFormatter(uint8(nGramSize), NGRAM_SEPARATOR)
			if ngram, err = ngramCmd.Execute(cipher); err == nil {
				fmt.Println("\t", ngramCmd)
				fmt.Println("NGram    : ", ngram)
			} else {
				exitCode = z.ERR_POST_CMD
			}
		}
	}

	return exitCode, err
}

// Decrypt a text using plain Caesar cipher with the given alphabet and key.
func DoDecrypt(alpha *cmn.Alphabet, key rune) (int, error) {
	var exitCode int = z.EXIT_CODE_SUCCESS
	// These commands implement IPipe, ICipherCommand & ICommand
	cmdCipher := commands.NewCaesarCommand(alpha, key)

	var plain string
	var err error
	if app.IsPipedInput() {
		reader := bufio.NewReader(os.Stdin)
		scanner := bufio.NewScanner(reader)
		var lineIn string
		for scanner.Scan() {
			lineIn = scanner.Text()
			plain, err = cmdCipher.Decode(lineIn)
			if err != nil {
				exitCode = z.ERR_CIPHER
				mlog.ErrorE(err)
				break
			}
			fmt.Println(plain)
		}

		if err = scanner.Err(); err != nil {
			exitCode = z.ERR_FILE_IO
			mlog.ErrorE(err)
		}
	} else { // short messages that can be given on the CLI
		cipher := flag.Arg(0)
		plain, err = cmdCipher.Decode(cipher)

		if err == nil {
			fmt.Println("Alphabet : ", alpha.Chars)
			fmt.Println("Key      : ", key)
			fmt.Println("Encoded  : ", cipher)
			fmt.Println("Plain    : ", plain)
		} else {
			exitCode = z.ERR_CIPHER
		}
	}

	return exitCode, err
}

// Prints the encoder/decoder tape (table) for the selected key and alphabet
func DoPrintTape(alpha *cmn.Alphabet, key rune, center, boxDrawing bool) {
	var sb strings.Builder
	// Prints a Row of Runes
	rowPrinterFunc := func(row []rune) {
		for _, char := range row {
			sb.WriteString(fmt.Sprintf("%c ", char))
		}

		sb.WriteRune('\n')
	}

	// convert key to a key shift value
	keyShift := alpha.PositionOf(key)

	// Generates a Space Leader to Center a string
	const MAX_WIDTH = 80
	centerLeaderFunc := func(length int) string {
		leaderLength := int((MAX_WIDTH - length) / 2)
		return strings.Repeat(" ", leaderLength)
	}

	var leader string = ""
	plain := []rune(alpha.Chars)
	if center {
		leader = centerLeaderFunc(len(plain)*2 + 2)
	}

	// Print Heading
	var bC, bH, bV rune
	if boxDrawing {
		bC = '\u250c' // ┌
		bH = '\u2500' // ─
		bV = '\u2502' // │
	} else {
		bC = '+'
		bH = '-'
		bV = '|'
	}

	title := fmt.Sprintf("%c %s %c\n", 0x00ab, alpha.Name, 0x00bb)
	sb.WriteString(centerLeaderFunc(len(title)))
	sb.WriteString(title)

	sb.WriteString(leader)
	sb.WriteString("  ")
	rowPrinterFunc(plain)
	sb.WriteString(fmt.Sprintf("%s %c%s\n", leader, bC, strings.Repeat(string(bH), 2*len(plain)-1)))

	encoded := []rune(cmn.RotateSliceRight(plain, keyShift))
	sb.WriteString(fmt.Sprintf("%s%c%c", leader, key, bV))
	rowPrinterFunc(encoded)

	fmt.Println(sb.String())
}

// Help about using this application
func Help() {
	fmt.Println("caesar -alpha {LANG} -key {LETTER} -tape")
	fmt.Println("caesar -alpha {LANG} -key {LETTER} [-ngram SIZE] 'plain text'")
	fmt.Println("caesar -alpha {LANG} -key {LETTER} -d 'cipher text'")
	flag.Usage()
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/
func main() {
	defer mlog.CloseLogFiles()

	// -------	CLI FLAGS ------
	var actHelp, actTape, actVersion, actDecrypt bool
	var optLang string
	var optKey cmd.RuneFlag
	var optNGram int
	flag.BoolVar(&actHelp, "help", false, "This help")
	flag.BoolVar(&actDecrypt, "d", false, "Decrypt text")
	flag.BoolVar(&actTape, "tape", false, "Print tabula for selected key")
	flag.BoolVar(&actVersion, "version", false, "Show application version")
	flag.StringVar(&optLang, "alpha", cmn.ALPHA_NAME_ENGLISH, "Encoding Alphabet")
	flag.IntVar(&optNGram, "ngram", -1, "Group encrypted output as NGram of 2|3|4|5")
	cmd.RegisterRuneVar(&optKey, "key", 0, "Encoding key")
	flag.Parse()

	// -------	CLI VALIDATION ------
	if !actHelp && !actVersion && !optKey.IsSet && (optKey.IsSet && optKey.Value == 0) {
		app.Die("A key must be provided", z.ERR_PARAMETER)
	}

	if optNGram != -1 && (optNGram < 2 || optNGram > 5) {
		app.Die("NGram size is 2|3|4|5", z.ERR_PARAMETER)
	}

	if optNGram != -1 && app.IsPipedInput() {
		app.Die("NGram only possible for short messages, not piped input.", z.ERR_CLI_OPTIONS)
	}

	var alphabet *cmn.Alphabet = nil
	if optLang != cmn.ALPHA_NAME_ENGLISH {
		// a composed alphabet
		alphabet = cmn.AlphabetComposer(optLang).(*cmn.Alphabet)
		if cmn.RuneIndex(alphabet.Chars, optKey.Value) == -1 {
			app.Die("That key is not part of the composed alphabet", z.ERR_PARAMETER)
		}
	} else {
		alphabet = cmn.AlphabetFactory(optLang).(*cmn.Alphabet)
	}

	// -------	EXECUTION ------
	var err error = nil
	var exitCode int = z.EXIT_CODE_SUCCESS
	switch {
	// -help
	case actHelp:
		Help()

	// -version
	case actVersion:
		fmt.Println("\tCaesar Lite version", APP_VERSION)

	case actTape:
		DoPrintTape(alphabet, optKey.Value, true, true)

	// -alpha {ALPHABET[+ALPHABET] -key {LETTER} -d
	case actDecrypt:
		exitCode, err = DoDecrypt(alphabet, optKey.Value)

	// -alpha {ALPHABET[+ALPHABET] -key {LETTER}
	case !actDecrypt:
		exitCode, err = DoEncrypt(alphabet, optKey.Value, optNGram)
	}

	// epilogue
	if exitCode != z.EXIT_CODE_SUCCESS {
		if err != nil {
			app.DieWithError(err, exitCode)
		} else {
			app.Die("an error occured", exitCode)
		}
	}

	if !app.IsPipedInput() {
		z.BuyMeCoffee()
	}
}
