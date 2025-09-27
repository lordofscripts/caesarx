/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   go-caesar
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Application-related functions.
 *-----------------------------------------------------------------*/
package app

import (
	"fmt"
	"lordofscripts/caesarx/app/mlog"
)

/* ----------------------------------------------------------------
 *						G l o b a l s
 *-----------------------------------------------------------------*/

const (
	UC_RED_EXCLAMATION = rune(0x2757) // Dingbats
)

/* ----------------------------------------------------------------
 *					F u n c t i o n s
 *-----------------------------------------------------------------*/

/**
 * Death of an application by outputting a good-bye and setting
 * the OS exit code.
 */
func Die(message string, exitCode int) {
	fmt.Println("\n", "\t💀 x 💀 x 💀\n\t", message, "\n\tExit code: ", exitCode)
	mlog.FatalT(exitCode, message, mlog.YesNo("Died", true), mlog.Int("Code", exitCode))
}

func DieWithError(err error, exitCode int) {
	fmt.Println("\n", "\t💀 x 💀 x 💀\n\t", err.Error(), "\n\tExit code: ", exitCode)
	mlog.FatalT(exitCode, err.Error(), mlog.YesNo("Died", true), mlog.Int("Code", exitCode))
}

func Assert(condition bool, warnMessage string) {
	if condition {
		fmt.Printf("\n\t%c Assertion Failed:\n\t%s\n", UC_RED_EXCLAMATION, warnMessage)
	}
}

func AssertOrDie(condition bool, deathMessage string, exitCode int) {
	if condition {
		fmt.Printf("\n\t%c Assertion Failed:", UC_RED_EXCLAMATION)
		Die(deathMessage, exitCode)
	}
}

func AnnounceErrorMessage(message string, exitCode int) {
	fmt.Println("\n", "\t💀 x 💀 x 💀\n\t", message, "\n\tExit code: ", exitCode)
}

func AnnounceError(err error, exitCode int) {
	fmt.Println("\n", "\t💀 x 💀 x 💀\n\t", err.Error(), "\n\tExit code: ", exitCode)
}
