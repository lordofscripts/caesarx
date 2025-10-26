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
	"os"
	"path/filepath"
	"runtime"
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

// Death of an application by outputting a good-bye and setting
// the OS exit code. It is logged as fatal.
func Die(message string, exitCode int) {
	fmt.Println("\n", "\t💀 x 💀 x 💀\n\t", message, "\n\tExit code: ", exitCode)
	mlog.FatalT(exitCode, message, mlog.YesNo("Died", true), mlog.Int("Code", exitCode))
}

// display the error and die with an exit code, logging it as Fatal.
func DieWithError(err error, exitCode int) {
	fmt.Println("\n", "\t💀 x 💀 x 💀\n\t", err.Error(), "\n\tExit code: ", exitCode)
	mlog.FatalT(exitCode, err.Error(), mlog.YesNo("Died", true), mlog.Int("Code", exitCode))
}

// When the condition is met the warning message is printed
func Assert(condition bool, warnMessage string) {
	if condition {
		fmt.Printf("\n\t%c Assertion Failed:\n\t%s\n", UC_RED_EXCLAMATION, warnMessage)
	}
}

// If the condition is met, the death message is printed and the
// application terminates with the exit code.
func AssertOrDie(condition bool, deathMessage string, exitCode int) {
	if condition {
		fmt.Printf("\n\t%c Assertion Failed:", UC_RED_EXCLAMATION)
		Die(deathMessage, exitCode)
	}
}

// prints the error emssage with the exit code but does NOT exit.
func AnnounceErrorMessage(message string, exitCode int) {
	fmt.Println("\n", "\t💀 x 💀 x 💀\n\t", message, "\n\tExit code: ", exitCode)
}

// prints the error and exit code but does NOT exit the application.
func AnnounceError(err error, exitCode int) {
	fmt.Println("\n", "\t💀 x 💀 x 💀\n\t", err.Error(), "\n\tExit code: ", exitCode)
}

// Returns true if the application input is not from a character device (tty)
// but instead from a piped input like "cat textfile.txt | yourapp -encrypt".
// When true you can use a bufio.Scanner to read text lines one by one and
// process them accordingly.
func IsPipedInput() bool {
	fi, _ := os.Stdin.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}

// platform-agnostic function to obtain the user's configuration directory.
// In Linux "~/.config/appName", Windows "APPDATA/appName" and
// MacOS "~/Library/Application Support/appName"
func GetConfigDir(orgName, appName string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), orgName, appName)
	case "darwin": // macOS
		return filepath.Join(homeDir, "Library", "Application Support", orgName, appName)
	default: // Other platforms (Linux, etc.)
		return filepath.Join(homeDir, ".config", orgName, appName)
	}
}

// Ensures a directory and all its parents exist and create them if necessary.
// Default permissions is 0755.
func EnsureConfigDir(path string) error {
	// Create the config directory if it doesn't exist
	err := os.MkdirAll(path, 0755) // 0755 permissions: rwxr-xr-x
	if err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	return nil
}

// Checks whether the file exists and is readable.
func CheckFileExistsAndReadable(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}
		return fmt.Errorf("error checking file: %w", err)
	}
	return nil
}
