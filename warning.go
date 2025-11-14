/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Alphabet is the basis of all plain text.
 *-----------------------------------------------------------------*/
package caesarx

import (
	"fmt"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

// Ensure the Warning type implements the error interface
var _ error = (*Warning)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// Warning struct represents a custom warning type
type Warning struct {
	Message string
	Code    WarningCode
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// A warning in the guise of an error
func NewWarningAsErr(msg string, pcode PackageCode, wcode uint16) error {
	return &Warning{
		Message: msg,
		Code:    NewWarningCode(pcode, wcode),
	}
}

// A pure warning but still implements errors.error
func NewWarning(msg string, pcode PackageCode, wcode uint16) *Warning {
	return &Warning{
		Message: msg,
		Code:    NewWarningCode(pcode, wcode),
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Implement the error interface
func (w *Warning) Error() string {
	return fmt.Sprintf("%s: %s", w.Code, w.Message)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// tells whether a given error instance represents a Warning and
// if so it casts it using type assertion. Else returns nil,false.
func IsWarning(err error) (*Warning, bool) {
	if err != nil {
		if warning, ok := err.(*Warning); ok {
			return warning, true
		}
	}

	return nil, false
}
