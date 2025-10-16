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
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	APP_NAME    = "affine"
	APP_VERSION = "1.1"
)

var (
	nameMasterAlphabet   string
	nameSlaveAlphabet    string
	namePrimaryOperation string
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

// Help shows help about using the Affine CLI application
func Help(co *cmd.CommonOptions) int {
	fmt.Println("Usage:")
	co.ShowUsage(APP_NAME)

	flag.Usage()
	return z.EXIT_CODE_SUCCESS
}

// PrintCoprimes will list all coprimes. By default it uses N from the alphabet
// size, i.e. the chosen alphabet. But if the "n" parameter is greater
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
	aparams, err := affine.NewAffineParams(a, b, int(alpha.Size()))
	if err != nil {
		app.DieWithError(err, z.ERR_INTERNAL)
	}

	// print Affine cipher (given & calculated) parameters
	fmt.Printf("\t%s\n", aparams)

	helper := affine.NewAffineHelper()
	err = helper.SetParams(aparams)
	if err != nil {
		app.DieWithError(err, z.ERR_INTERNAL)
	}

	var tabula string
	if tabula, err = helper.GetTabulaString(alpha.Chars, "\t"); err != nil {
		app.DieWithError(err, z.ERR_INTERNAL)
	} else {
		fmt.Println(tabula)
	}

	return z.EXIT_CODE_SUCCESS
}

// ExecuteMessage performs encryption or decryption of a (single) text string
// specified as a free CLI argument. The user has named the alpha Master alphabet,
// and an optional numbers Slave alphabet.
func ExecuteMessage(alpha, numbers *cmn.Alphabet, opts *AffineCliOptions, input string) (int, error) {
	var err error
	var exitCode int = z.EXIT_CODE_SUCCESS
	var output string

	cmdCipher := setupAffineCrypto(alpha, numbers, opts)

	if opts.ActIsDecode {
		output, err = cmdCipher.Decode(input)
	} else {
		output, err = cmdCipher.Encode(input)
	}

	if err != nil {
		exitCode = z.ERR_CIPHER
	} else {
		//var params *affine.AffineParams
		//params, err = affine.NewAffineParams(a, b, int(alpha.Size()))
		paramsM, paramsS := cmdCipher.GetParams()

		fmt.Println("Operation: ", namePrimaryOperation)
		fmt.Printf("Alphabet : %s (Master/Primary)\n", nameMasterAlphabet)
		fmt.Printf("Alphabet : %s (Slave/Secondary)\n", nameSlaveAlphabet)
		fmt.Println("Params  M: ", paramsM)
		fmt.Println("Params  S: ", paramsS)
		fmt.Println("Algorithm: ", cmdCipher.String())
		if opts.ActIsDecode {
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

// ExecuteFile performs encryption or decryption of a file. In this case for encryption
// the user specifies the input filename and the output filename is derived by the
// application (see documentation). For decryption the user must specify two free
// arguments, the input and output filenames respectively.
func ExecuteFile(alpha, numbers *cmn.Alphabet, opts *AffineCliOptions) (int, error) {
	var err error
	var exitCode int = z.EXIT_CODE_SUCCESS

	cmdCipher := setupAffineCrypto(alpha, numbers, opts)

	if opts.ActIsDecode {
		if !opts.Common.IsBinary() {
			err = cmdCipher.DecryptTextFile(opts.Files.Input, opts.Files.Output)
		} else {
			err = cmdCipher.DecryptBinFile(opts.Files.Input, opts.Files.Output)
		}
	} else {
		if !opts.Common.IsBinary() {
			err = cmdCipher.EncryptTextFile(opts.Files.Input)
		} else {
			err = cmdCipher.EncryptBinFile(opts.Files.Input)
		}
	}

	if err != nil {
		exitCode = z.ERR_CIPHER
	} else {
		//var params *affine.AffineParams
		//params, err = affine.NewAffineParams(a, b, int(alpha.Size()))
		paramsM, paramsS := cmdCipher.GetParams()

		fmt.Println("Operation: ", namePrimaryOperation)
		fmt.Printf("Alphabet : %s (Master/Primary)\n", nameMasterAlphabet)
		fmt.Printf("Alphabet : %s (Slave/Secondary)\n", nameSlaveAlphabet)
		fmt.Println("Params  M: ", paramsM)
		fmt.Println("Params  S: ", paramsS)
		fmt.Println("Algorithm: ", cmdCipher.String())
		if opts.ActIsDecode {
			fmt.Println("Encoded  : ", opts.Files.Input)
			fmt.Println("Decoded  : ", opts.Files.Output)
		} else {
			fmt.Println("Plain    : ", opts.Files.Input)
			fmt.Println("Encoded  : ", opts.Files.Output)
		}
		fmt.Println()
	}

	return exitCode, err
}

// setupAffineCrypto does the preliminary setup for the cryptographic operation.
// and returns an Affine cryptographic command object capable of performing the
// actual encryption/decryption.
func setupAffineCrypto(alpha, numbers *cmn.Alphabet, opts *AffineCliOptions) *commands.AffineCommand {
	// the main Affine cipher engine only has a language/letters alphabet and no slave/chain
	// cmdCipher implements ciphers.ICipherCommand
	cmdCipher := commands.NewAffineCommand(alpha, opts.CoefficientA, opts.CoefficientB)
	// attach any optional alphabet if any
	nameSlaveAlphabet = "(None)"
	if numbers != nil {
		cmdCipher.WithChain(numbers)
		nameSlaveAlphabet = numbers.Name
	}

	if opts.OptNgramSize > 0 && !opts.ActIsDecode {
		cmdNGram := cmn.NewNgramFormatter(uint8(opts.OptNgramSize), '·') // is cmn.ICommand
		cmdCipher.WithPipe(cmdNGram)
	}

	if opts.ActIsDecode {
		namePrimaryOperation = "Decrypt"
	} else {
		namePrimaryOperation = "Encrypt"
	}

	return cmdCipher
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/
func main() {
	// -------	CLI FLAGS ------
	copts := cmd.NewCommonOptions() // -help|-demo|-alpha ALPHA|-num N
	aopts := NewAffineOptions(copts)

	// -------	CLI VALIDATION ------
	if exitCode, err := copts.Validate(); err != nil {
		app.DieWithError(err, exitCode)
	}

	if !copts.IsReady() { // only if the common option is NOT terminal
		if exitCode, err := aopts.Validate(); err != nil {
			app.DieWithError(err, exitCode)
		}
	}

	// -------	EXECUTION ------
	var exitCode int = z.EXIT_CODE_SUCCESS
	var err error = nil

	switch {
	/*
	 * Common terminal arguments
	 */
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

	/*
	 * Affine-specific terminal arguments
	 */
	case aopts.ActListCoprimes:
		exitCode = PrintCoprimes(copts.Alphabet(), aopts.OptModulo)

	case aopts.ActPrintTabula:
		exitCode = PrintAffineTabula(copts.Alphabet(), aopts.CoefficientA, aopts.CoefficientB)

	/*
	 * Cryptographic operations
	 */
	case aopts.UseFiles():
		// the input is a filename specified in the CLI arguments
		exitCode, err = ExecuteFile(copts.Alphabet(), copts.Numbers(), aopts)

	default:
		// the input is a short message given as a CLI argument
		exitCode, err = ExecuteMessage(copts.Alphabet(), copts.Numbers(), aopts, flag.Arg(0))
	}

	if err != nil {
		app.AnnounceError(err, exitCode)
		os.Exit(exitCode)
	}

	z.BuyMeCoffee()
}
