/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							 go-caesarx
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Numeric error (exit) codes
 *-----------------------------------------------------------------*/
package caesarx

import "fmt"

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// v1.20 due to the nature of the math.Rand generator
	GO_MIN_REQUIRED = "1.20"

	EXIT_CODE_SUCCESS = 0
	ERR_PARAMETER     = 1
	ERR_NO_ALPHABET   = 2
	ERR_CLI_OPTIONS   = 3
	ERR_GO_VERSION    = 4
	ERR_DEMO_ERROR    = 9
	ERR_BAD_ALPHABET  = 10
	ERR_SEQUENCER     = 50
	ERR_CIPHER        = 51
	ERR_INTERNAL      = 126
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

func init() {
	current, ok := GoVersionMin(GO_MIN_REQUIRED)
	if !ok {
		fmt.Printf("This application requires GO v%s+ we have v%s\n", GO_MIN_REQUIRED, current)
	}
}
