/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * A custom Byte Flag for the GO flag package. We can now use byte
 * valued flags in the command line. For that we implement the flag.Value
 * interface.
 * This implementation works with both single and multi-byte runes.
 *-----------------------------------------------------------------*/
package cmd

import (
	"flag"
	"fmt"
	"strconv"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/
var _ flag.Value = (*ByteFlag)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type ByteFlag struct {
	Value byte
	IsSet bool
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (r *ByteFlag) String() string {
	if r.IsSet {
		return string(r.Value)
	}
	return ""
}

func (r *ByteFlag) Set(value string) error {
	val, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("not an integer value: %s", value)
	}
	if val < 0 || val >= 256 {
		return fmt.Errorf("not a byte value: %d", val)
	}

	r.Value = byte(val)
	r.IsSet = true
	return nil
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func RegisterByteVarSet(r *ByteFlag, name string, value byte, usage string) {
	r.Value = value
	r.IsSet = false
	flag.Var(r, name, usage)
}

func RegisterByteVar(r *ByteFlag, name string, usage string) {
	r.Value = 0
	r.IsSet = false
	flag.Var(r, name, usage)
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

/*
func DemoByteFlag() {
	var myByte ByteFlag
	flag.Var(&myByte, "byte", "custom Byte value")
	flag.Parse()

	fmt.Printf("Byte value: %c\n", myByte.Value)
}
*/
