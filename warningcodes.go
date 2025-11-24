/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Application warning codes.
 *-----------------------------------------------------------------*/
package caesarx

import (
	"fmt"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ApplicationPCode PackageCode = iota
	CiphersPCode
	CommandPCode
	CommonPCode
	InternalPCode
	ConfigurationPCode
)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type PackageCode uint8

// Warning struct represents a custom warning type
type WarningCode struct {
	Package PackageCode
	Code    uint16
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// Compose a warning code with its package identifier
func NewWarningCode(pkg PackageCode, wcode uint16) WarningCode {
	return WarningCode{
		Package: pkg,
		Code:    wcode,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// Implement the Stringer interface
func (w WarningCode) String() string {
	return fmt.Sprintf("%s Warning W%03d-%d", toPackageName(w.Package), w.Package, w.Code)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// converts a package enum to its name
func toPackageName(pkg PackageCode) string { // @note must match iota definition
	pn := map[PackageCode]string{
		ApplicationPCode:   "Application",
		CiphersPCode:       "Ciphers",
		CommandPCode:       "Command",
		CommonPCode:        "Common",
		InternalPCode:      "Internal",
		ConfigurationPCode: "Configuration",
	}

	if name, ok := pn[pkg]; ok {
		return name
	}

	return ""
}

/*
func DemoWarning() {
	if err != nil {
		if warning, ok := err.(*Warning); ok {
			fmt.Println("Handled warning:", warning)
		} else {
			fmt.Println("Error:", err)
		}
	}
}
*/
