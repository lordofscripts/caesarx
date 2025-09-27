/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A custom Rune Flag for the GO flag package. We can now use rune
 * flags in the command line. For that we implement the flag.Value
 * interface.
 * This implementation works with both single and multi-byte runes.
 *-----------------------------------------------------------------*/
package cmd

import (
	"flag"
	"fmt"
	"unicode/utf8"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/
var _ flag.Value = (*RuneFlag)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type RuneFlag struct {
	Value rune
	IsSet bool
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (r *RuneFlag) String() string {
	if r.IsSet {
		return string(r.Value)
	}
	return ""
}

func (r *RuneFlag) Set(value string) error {
	if utf8.RuneCountInString(value) != 1 {
		return fmt.Errorf("invalid rune: %s", value)
	}

	r.Value = []rune(value)[0]
	r.IsSet = true
	return nil
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func RegisterRuneVar(r *RuneFlag, name string, value rune, usage string) {
	r.Value = value
	r.IsSet = false
	flag.Var(r, name, usage)
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

/*
func DemoRuneFlag() {
	var myRune RuneFlag
	flag.Var(&myRune, "rune", "custom Rune value")
	flag.Parse()

	fmt.Printf("Rune value: %c\n", myRune.Value)
}
*/
