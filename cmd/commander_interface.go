/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Yet another Command design pattern
 *-----------------------------------------------------------------*/
package cmd

import (
	"fmt"
	"strings"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type ICommander interface {
	// execute the command with the input provided in the constructor.
	// returns nil on success.
	Execute() error

	// get command output and return it. If print is set it gets printed
	// on the console
	GetOutput(print bool) string

	// implements fmt.Stringer giving general info about the command
	String() string
}

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type CommanderBase struct {
	out strings.Builder
}

// appends the formatted string to the output buffer that will
// be fetched by calling GetOutput()
func (cb *CommanderBase) Output(format string, params ...any) {
	cb.out.WriteString(fmt.Sprintf(format, params...))
}

// Retrieves all the command execution output that was created by
// successive calls to Output()
func (cb *CommanderBase) GetOutput(print bool) string {
	if print {
		fmt.Println(cb.out.String())
	}

	return cb.out.String()
}
