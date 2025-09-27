/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Affine cipher application. This is a simple Caesar-like cipher
 * based on a linear formula y = (Ax + B) % N. Despite the appearance
 * of the formula, it remains a monoalphabetic cipher; therefore, it
 * inherits the weaknesses of those ciphers. But it is a fun,
 * educational experiment.
 *-----------------------------------------------------------------*/
package main

import (
	"flag"
	"fmt"
	z "lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/ciphers/affine"
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/cmd"
	"lordofscripts/caesarx/cmn"
	"os"
	"strconv"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	APP_NAME    = "affine"
	APP_VERSION = "1.0"
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/
func init() {
	//fmt.Println("GoCaesarPlus v1.0 (C)2025 Didimo Grimaldo \u2720 " + caesarx.RuneString("LordOfScripts"))
	z.Copyright(z.CO1, true)
	z.BuyMeCoffee()
	fmt.Println("\t=========================================")
}

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// List all coprimes. By default it uses N from the alphabet size,
// i.e. the chosen alphabet. But if the "n" parameter is greater
// than zero, it is assumed that the caller wants to experiment with
// an arbitrary length other than the chosen alphabet.
func PrintCoprimes(alpha *cmn.Alphabet, n int) int {
	if n <= 0 {
		n = int(alpha.Size()) // default to alphabet length
	}

	fmt.Println("COPRIMES")
	fmt.Printf("\tAlphabet: %s\n\tCharacters (N): %d\n", alpha.Name, n)
	fmt.Println("\tValid Coprimes, i.e. gcd(A,N) = 1")

	helper := affine.NewAffineHelper()
	coprimes := helper.ValidCoprimesUpTo(uint(n))
	for _, a := range coprimes {
		fmt.Printf("\tA = %d\n", a)
	}
	fmt.Println("\tCoefficient 'B' can be any positive integer.")

	return z.EXIT_CODE_SUCCESS
}

func PrintAffineTabula(alpha *cmn.Alphabet, a, b int) int {
	helper := affine.NewAffineHelper()
	if err := helper.SetParameters(a, b, int(alpha.Size())); err != nil {
		app.DieWithError(err, z.ERR_INTERNAL)
	}

	// print Affine cipher (given & calculated) parameters
	aparams := helper.GetParams()
	fmt.Printf("\t%s\n", aparams)

	if tabula, err := helper.GetTabulaString(alpha.Chars, "\t"); err != nil {
		app.DieWithError(err, z.ERR_INTERNAL)
	} else {
		fmt.Println(tabula)
	}

	return z.EXIT_CODE_SUCCESS
}

// Execute Affine encryption OR decryption.
// @param alpha (*cmnt.Alphabet) Primary/Master Reference Alphabet (letters)
// @param numbers (*cmnt.Alphabet) optional Secondary/Slave Reference Alphabet (numbers and/or symbols)
// @param decode (bool) true for decryption, false for encryption
// @param a (int) Affine A coefficient
// @param b (int) Affine B coefficient
// @param ngram (int) ignored if 0, else group CIPHERED output in 2/3/4/5 letters
// @param input (string) Text to encrypt or decrypt
// @returns (int) application exit code, 0 for success
// @returns (error) nil on success, else error
func Execute(alpha, numbers *cmn.Alphabet, decode bool, a, b, ngram int, input string) (int, error) {
	// the main Affine cipher engine only has a language/letters alphabet and no slave/chain
	cmdCipher := commands.NewAffineCommand(alpha, a, b) // is ciphers.ICipherCommand
	// attach any optional alphabet if any
	slaveName := "(None)"
	if numbers != nil {
		cmdCipher.WithChain(numbers)
		slaveName = numbers.Name
	}

	if ngram > 0 && !decode {
		cmdNGram := cmn.NewNgramFormatter(uint8(ngram), '·') // is cmn.ICommand
		cmdCipher.WithPipe(cmdNGram)
	}

	var err error
	var exitCode int = z.EXIT_CODE_SUCCESS
	var output string
	var operation string

	if decode {
		operation = "Decrypt"
		output, err = cmdCipher.Decode(input)
	} else {
		operation = "Encrypt"
		output, err = cmdCipher.Encode(input)
	}

	if err != nil {
		exitCode = z.ERR_CIPHER
	} else {
		//var params *affine.AffineParams
		//params, err = affine.NewAffineParams(a, b, int(alpha.Size()))
		paramsM, paramsS := cmdCipher.GetParams()

		fmt.Println("Operation: ", operation)
		fmt.Printf("Alphabet : %s (Master/Primary)\n", alpha.Name)
		fmt.Printf("Alphabet : %s (Slave/Secondary)\n", slaveName)
		fmt.Println("Params  M: ", paramsM)
		fmt.Println("Params  S: ", paramsS)
		fmt.Println("Algorithm: ", cmdCipher.String())
		if decode {
			fmt.Println("Encoded  : ", input)
			fmt.Println("Decoded  : ", output)
		} else {
			fmt.Println("Plain    : ", input)
			fmt.Println("Encoded  : ", output)
		}
		fmt.Println()
	}

	return exitCode, err
}

func Help(co *cmd.CommonOptions) int {
	fmt.Println("Usage:")
	co.ShowUsage(APP_NAME)

	flag.Usage()
	return z.EXIT_CODE_SUCCESS
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/
func main() {
	const (
		FLAG_COEFF_A  = "A"
		FLAG_NGRAM    = "ngram" // (optional) only if encrypting
		FLAG_COEFF_B  = "B"
		FLAG_DECODE   = "d"
		FLAG_COPRIMES = "coprime"
		FLAG_MODULO   = "N" // (optional) only if -coprime is given
		FLAG_TABULA   = "tabula"
	)
	// -------	CLI FLAGS ------
	copts := cmd.NewCommonOptions() // -help|-demo|-alpha ALPHA|-num N
	var optNgram, optCoeffA, optCoeffB, optModulo int
	var decode, listCoprimes, printTabula bool

	flag.IntVar(&optCoeffA, FLAG_COEFF_A, 1, "Affine coefficient A")
	flag.IntVar(&optCoeffB, FLAG_COEFF_B, 0, "Affine coefficient B")
	flag.IntVar(&optModulo, FLAG_MODULO, 0, "Affine module N (only if -coprime is used), else derived from alpha")
	flag.IntVar(&optNgram, FLAG_NGRAM, 0, "Format encoded output as NGram")
	flag.BoolVar(&decode, FLAG_DECODE, false, "Decode text")
	flag.BoolVar(&listCoprimes, FLAG_COPRIMES, false, "List coprimes for 'A' for the chosen alphabet")
	flag.BoolVar(&printTabula, FLAG_TABULA, false, "Print Tabula for chosen parameters")
	flag.Parse()

	// -------	CLI VALIDATION ------
	if !(optNgram == 0 || (optNgram >= 2 && optNgram <= 5)) {
		app.Die("Ngram size must be 2,3,4 or 5 not "+strconv.Itoa(optNgram), z.ERR_PARAMETER)
	}

	if optModulo != 0 && !listCoprimes {
		app.Die("optional -N can only be used with -coprime but it defaults to alphabet length", z.ERR_PARAMETER)
	}

	if optModulo < 0 {
		app.Die("-N should be positive integer", z.ERR_PARAMETER)
	}

	if !copts.NeedsDemo() && !copts.NeedsHelp() && !listCoprimes && !copts.NeedsVersion() && !printTabula && flag.NArg() != 1 {
		app.Die("for encode/decode the free argument must be the text.", z.ERR_PARAMETER)
	}

	// -------	EXECUTION ------
	var exitCode int = z.EXIT_CODE_SUCCESS
	var err error = nil

	switch {
	case copts.NeedsHelp():
		exitCode = Help(copts)

	case copts.NeedsDemo():
		passed := affine.DemoAffine()
		if !passed {
			exitCode = z.ERR_DEMO_ERROR
		}

	case copts.NeedsVersion():
		fmt.Printf("\tAffine cipher app. v%s\n", APP_VERSION)
		fmt.Println("\t", affine.Info)

	case listCoprimes:
		exitCode = PrintCoprimes(copts.Alphabet(), optModulo)

	case printTabula:
		exitCode = PrintAffineTabula(copts.Alphabet(), optCoeffA, optCoeffB)

	default:
		exitCode, err = Execute(copts.Alphabet(), copts.Numbers(), decode, optCoeffA, optCoeffB, optNgram, flag.Arg(0))
	}

	if err != nil {
		app.AnnounceError(err, exitCode)
		os.Exit(exitCode)
	}

	z.BuyMeCoffee()
}
