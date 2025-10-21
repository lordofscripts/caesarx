/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Extended Caesar Cipher command-line application. It supports the
 * following ciphers: Caesar, Didimus, Fibonacci, Bellaso & Vigenère.
 *-----------------------------------------------------------------*/
package main

import (
	"flag"
	"fmt"
	z "lordofscripts/caesarx"
	"lordofscripts/caesarx/app"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/ciphers"
	"lordofscripts/caesarx/ciphers/affine"
	"lordofscripts/caesarx/ciphers/bellaso"
	"lordofscripts/caesarx/ciphers/caesar"
	"lordofscripts/caesarx/ciphers/commands"
	"lordofscripts/caesarx/ciphers/vigenere"
	"lordofscripts/caesarx/cmd"
	"lordofscripts/caesarx/cmn"
	"os"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	APP_VERSION = "1.1"

	VariantCaesar CaesarVariant = iota
	VariantDidimus
	VariantFibonacci
	VariantBellaso
	VariantVigenere
	VariantAffine
)

type CaesarVariant uint8

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/
func init() {
	z.Copyright(z.CO1, true)
	z.BuyMeCoffee()
	fmt.Println("\t=========================================")
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func Demo(copts *cmd.CommonOptions, aopts *CaesarxOptions) (int, error) {
	var passed bool // the demos return whether the Round-trip Encode/Decode was good
	switch aopts.VariantID {
	case VariantCaesar: // -variant caesar -alpha <ALPHABET_NAME> -key <LETTER>
		passed = commands.DemoCaesarCommand(copts.Alphabet(), copts.Numbers(), copts.DefaultPhrase)

	case VariantDidimus: // -variant didimus -alpha <ALPHABET_NAME> -key <LETTER> -offset <NUMBER>
		passed = commands.DemoDidimusCommand(copts.Alphabet(), copts.Numbers(), copts.DefaultPhrase)

	case VariantFibonacci: // -variant fibonacci -alpha <ALPHABET_NAME> -key <LETTER>
		passed = commands.DemoFibonacciCommand(copts.Alphabet(), copts.Numbers(), copts.DefaultPhrase)

	case VariantBellaso: // -variant bellaso -alpha <ALPHABET_NAME> -secret <SECRET_WORD>
		passed = commands.DemoBellasoCommand(copts.Alphabet(), copts.Numbers(), copts.DefaultPhrase)

	case VariantVigenere: // -variant vigenere -alpha <ALPHABET_NAME> -secret <SECRET_WORD>
		passed = commands.DemoVigenereCommand(copts.Alphabet(), copts.Numbers(), copts.DefaultPhrase)

	case VariantAffine:
		passed = affine.DemoAffine()
	}

	var err error = nil
	var exitCode int = z.EXIT_CODE_SUCCESS
	if !passed {
		err = fmt.Errorf("round-trip encryption/decryption FAILED")
		exitCode = z.ERR_DEMO_ERROR
	}

	return exitCode, err
}

func DoCrypto(co *cmd.CommonOptions, ao *CaesarxOptions) (int, error) {
	var tempOut string = "" // temporary filename IF used
	var postCmd cmd.ICommander = nil

	// These commands implement IPipe, ICipherCommand & ICommand
	var cmdCipher ciphers.ICipherCommand
	switch ao.VariantID {
	case VariantCaesar:
		// single key, space not included
		cmdCipher = commands.NewCaesarCommand(co.Alphabet(), ao.MainKey.Value)

	case VariantDidimus:
		// double alternating key, space, numbers and number-related symbols included
		cmdCipher = commands.NewDidimusCommand(co.Alphabet(), ao.MainKey.Value, uint8(ao.Offset))

	case VariantFibonacci:
		// key plus 10-term Fibonacci offsets, space, numbers and number-related symbols included
		cmdCipher = commands.NewFibonacciCommand(co.Alphabet(), ao.MainKey.Value)

	case VariantBellaso:
		// multi-letter secret, space, numbers and number-related symbols included
		cmdCipher = commands.NewBellasoCommand(co.Alphabet(), ao.Secret)

	case VariantVigenere:
		// multi-letter secret & autokey, space, numbers and number-related symbols included
		cmdCipher = commands.NewVigenereCommand(co.Alphabet(), ao.Secret)

	case VariantAffine:
		fmt.Println("Please use the affine (affine.exe) application.")
		fallthrough

	default:
		return z.ERR_PARAMETER, fmt.Errorf("unknown algorithm") // @audit other error code please
	}

	// Check if user wants to enhance cipher with Slave alphabet
	if _, wants := co.WantsSlave(); wants {
		cmdCipher.WithChain(co.Numbers())
	}

	// If NGram formatting wanted, create it as Pipe command, only for Encoding
	if ao.NGramSize > 0 && !ao.IsDecode {
		ngramCmd := cmn.NewNgramFormatter(uint8(ao.NGramSize), '·')
		cmdCipher.WithPipe(ngramCmd) // @audit add Tee Command to output regular and NGram
	}

	// Do the (de)cipher operation
	var plain, cipher, operation string
	var err error
	if ao.IsDecode {
		operation = "Decrypt"
		cipher = flag.Arg(0)
		if ao.UseFiles {
			if co.IsBinary() {
				err = cmdCipher.DecryptBinFile(ao.Files.Input, ao.Files.Output)
			} else {
				err = cmdCipher.DecryptTextFile(ao.Files.Input, ao.Files.Output)
			}

			// is file verification requested
			if err == nil && ao.OptVerify {
				postCmd = cmd.NewVerifyFileCommand(ao.Files.Output, cmd.HashCRC64)
			}
		} else { // short messages that can be given on the CLI
			plain, err = cmdCipher.Decode(cipher)
		}
	} else {
		operation = "Encrypt"
		plain = flag.Arg(0)
		if ao.UseFiles {
			if co.IsBinary() {
				err = cmdCipher.EncryptBinFile(ao.Files.Input)
			} else {
				err = cmdCipher.EncryptTextFile(ao.Files.Input)
			}

			// For round-trip verification if -verify is given
			if err == nil && ao.OptVerify {
				// temporary filename for round-trip
				tempOut = cmn.GenerateTemporaryFileName("tempfile-caesarx-*") //+ ao.FileExt()
				cipherFilename := cmdCipher.GetOutputFilename()

				if !ao.Common.IsBinary() {
					err = cmdCipher.DecryptTextFile(cipherFilename, tempOut)
				} else {
					err = cmdCipher.DecryptBinFile(cipherFilename, tempOut)
				}

				// issue Verify command ONLY if the temporary decrypted file exists
				if err == nil {
					postCmd = cmd.NewVerifyFilesCommand(ao.Files.Input, tempOut, cmd.HashMD5)
				} else {
					tempOut = ""
				}
			}
		} else {
			cipher, err = cmdCipher.Encode(plain)
		}
	}

	if err != nil {
		return z.ERR_INTERNAL, err
	} else {
		// common output
		numbersName, _ := co.WantsSlave()
		fmt.Println("Algorithm: ", cmdCipher.String())
		fmt.Println("Info     : ", ao.VariantVersion)
		fmt.Println("Operation: ", operation)
		fmt.Printf("Alphabet : %s (Master/Primary)\n", co.Alphabet().Name)
		fmt.Printf("Alphabet : %s (Slave/Secondary)\n", numbersName)
		// assorted parameters
		switch ao.ItNeeds {
		case NeedCompositeKey:
			fmt.Printf("Offset   : %d \n", ao.Offset)
			fallthrough // composite needs -key as well

		case NeedKey:
			fmt.Printf("Key      : %c (shift=%d)\n", ao.MainKey.Value, co.Alphabet().PositionOf(ao.MainKey.Value))

		case NeedsSecret:
			fmt.Printf("Secret   :  %s\n", ao.Secret)

		}
		// input/output relations
		if ao.IsDecode {
			if ao.UseFiles {
				cipher = ao.Files.Input
				plain = ao.Files.Output
			}
			fmt.Println("Encoded  : ", cipher)
			fmt.Println("Decoded  : ", plain)
		} else {
			if ao.UseFiles {
				plain = ao.Files.Input
				cipher = ao.Files.Output
			}
			fmt.Println("Plain    : ", plain)
			fmt.Println("Encoded  : ", cipher)
		}

		// any post-execution command?
		if postCmd != nil {
			// perform the file verification
			if err = postCmd.Execute(); err == nil {
				fmt.Println("\t", postCmd)
				postCmd.GetOutput(true)
			}

			// remove temporary file
			if tempOut != "" {
				os.Remove(tempOut)
			}
		}

		fmt.Println()
	}

	return z.EXIT_CODE_SUCCESS, nil
}

func Help(co *cmd.CommonOptions, ao *CaesarxOptions) {
	fmt.Println("Usage:")
	co.ShowUsage(APP_NAME)
	ao.ShowUsage(APP_NAME)

	flag.Usage()
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/
func main() {
	var exitCode int = -1
	var err error
	defer mlog.CloseLogFiles()

	// -------	CLI FLAGS ------
	copts := cmd.NewCommonOptions() // -help|-demo|-alpha ALPHA|-num N
	aopts := NewCaesarxOptions(copts)

	// -------	CLI VALIDATION ------
	exitCode, err = copts.Validate()
	if err != nil {
		app.DieWithError(err, exitCode)
	}

	if !copts.IsReady() {
		exitCode, err = aopts.Validate()
		if err != nil {
			app.DieWithError(err, exitCode)
		}
	}

	// -------	EXECUTION ------
	switch {
	// -demo
	case copts.NeedsDemo(): // @note pass Numeric disk to Demo functions
		exitCode, err = Demo(copts, aopts)

	// -list
	case copts.NeedsList():
		fmt.Println("List of available Caesar-cipher class variants:")
		fmt.Println(ciphers.PrintAvailableCiphers())
		exitCode = z.EXIT_CODE_SUCCESS

	// -help
	case copts.NeedsHelp():
		Help(copts, aopts)
		exitCode = z.EXIT_CODE_SUCCESS

	// -version
	case copts.NeedsVersion():
		fmt.Printf("\tCaesarX cipher app. v%s\n", APP_VERSION)
		fmt.Println("\t", caesar.Info)
		fmt.Println("\t", caesar.InfoDidimus)
		fmt.Println("\t", caesar.InfoFibonacci)
		fmt.Println("\t", bellaso.Info)
		fmt.Println("\t", vigenere.Info)
		exitCode = z.EXIT_CODE_SUCCESS

	// -d or encrypt
	default:
		exitCode, err = DoCrypto(copts, aopts)
	}

	// epilogue
	if exitCode != z.EXIT_CODE_SUCCESS {
		if err != nil {
			app.DieWithError(err, exitCode)
		} else {
			app.Die("an error ocurred", exitCode)
		}
	}

	z.BuyMeCoffee()
}
