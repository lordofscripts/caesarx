/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Logging to Console with ANSI Color codes. But it always logs to the
 * console regardless of logging level, colorful simplicity.
 *-----------------------------------------------------------------*/
package mlog

import (
	"fmt"
	"os"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var Console *ColorConsole = &ColorConsole{}

const ( // some from https://azrael.digipen.edu/~mmead/www/mg/ansicolors/index.html
	ColorBlack       Color = "\u001b[30m"
	ColorRed         Color = "\u001b[31m"
	ColorLightRed    Color = "\u001b[91m"
	ColorGreen       Color = "\u001b[32m"
	ColorBrown       Color = "\u001b[33m"
	ColorYellow      Color = "\u001b[93m"
	ColorPurple      Color = "\u001b[35m"
	ColorLightPurple Color = "\u001b[95m"
	ColorBlue        Color = "\u001b[34m"
	ColorMagenta     Color = "\u001b[35m"
	ColorCyan        Color = "\u001b[36m"
	SlowBlink        Color = "\u001b[5m"
	BlinkOff         Color = "\u001b[25m"
	ColorReset       Color = "\u001b[0m"
	BoldOn           Color = "\u001b[1m"
	BoldOff          Color = "\u001b[0m" //"\u001b[21m"
	UnderlineOn      Color = "\u001b[4m"
	UnderlineOff     Color = "\u001b[24m"

	face1 string = "⚆ _ ⚆"
	face2 string = "¯\\_(ツ)_/¯"
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type Color string

type ColorConsole struct {
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Trace level with format string
func (c *ColorConsole) Trace(format string, args ...any) {
	fmt.Print(ColorLightPurple, tagTRACE)
	fmt.Printf(format, args...)
	fmt.Print(ColorReset)
}

// Debug level with format string
func (c *ColorConsole) Debug(format string, args ...any) {
	fmt.Print(ColorBrown, tagDEBUG)
	fmt.Printf(format, args...)
	fmt.Print(ColorReset)
}

// Information level with format string
func (c *ColorConsole) Info(format string, args ...any) {
	fmt.Print(ColorGreen, tagINFO)
	fmt.Printf(format, args...)
	fmt.Print(ColorReset)
}

// Warning level with format string
func (c *ColorConsole) Warn(format string, args ...any) {
	fmt.Print(ColorYellow, tagWARN)
	fmt.Printf(format, args...)
	fmt.Print(ColorReset)
}

// Error level with format string
func (c *ColorConsole) Error(format string, args ...any) {
	fmt.Print(ColorPurple, tagERROR)
	fmt.Printf(format, args...)
	fmt.Print(ColorReset)
}

// Fatal level with format string
func (c *ColorConsole) Fatal(exitCode int, format string, args ...any) {
	fmt.Print(ColorRed, tagFATAL)
	fmt.Printf(format, args...)
	fmt.Println("\t\t", face2, ColorReset)
	os.Exit(exitCode)
}

/*
func demo() {
	mlog.Console.Trace("Trace %d\n", 1)
	mlog.Console.Debug("Debug %d\n", 2)
	mlog.Console.Info("Info %d\n", 3)
	mlog.Console.Warn("Warn %d\n", 4)
	mlog.Console.Error("Error %d\n", 5)
	mlog.Console.Fatal(110, "Fatal %d %s\n", 6, mlog.At())
}
*/
