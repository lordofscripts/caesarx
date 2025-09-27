/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							 go-caesarx
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Numeric error (exit) codes
 *-----------------------------------------------------------------*/
package caesarx

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	EXIT_CODE_SUCCESS = 0
	ERR_PARAMETER     = 1
	ERR_NO_ALPHABET   = 2
	ERR_CLI_OPTIONS   = 3
	ERR_DEMO_ERROR    = 9
	ERR_BAD_ALPHABET  = 10
	ERR_SEQUENCER     = 50
	ERR_CIPHER        = 51
	ERR_INTERNAL      = 126
)
