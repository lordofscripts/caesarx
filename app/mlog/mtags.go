/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * MLog variadic tags: String, Rune, Int, Bool, YesNo, Byte & At.
 * Tags like the log/slog package but enhanced.
 *-----------------------------------------------------------------*/
package mlog

import (
	"fmt"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

// ILogKeyValuePair defines the interface for mlog tags that can
// appear as variadic parameters to the mlog logging functions.
type ILogKeyValuePair interface {
	fmt.Stringer
}

var _ ILogKeyValuePair = (*kvString)(nil)
var _ ILogKeyValuePair = (*kvRune)(nil)
var _ ILogKeyValuePair = (*kvInt)(nil)
var _ ILogKeyValuePair = (*kvBool)(nil)
var _ ILogKeyValuePair = (*kvYesNo)(nil)
var _ ILogKeyValuePair = (*kvByte)(nil)
var _ ILogKeyValuePair = (*kvAt)(nil)
var _ ILogKeyValuePair = (*kvError)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type kvString struct {
	k string
	v string
}

type kvRune struct {
	k string
	v rune
}

type kvInt struct {
	k string
	v int
}

type kvBool struct {
	k string
	v bool
}

type kvYesNo struct {
	k string
	v bool
}

type kvByte struct {
	k string
	v byte
}

type kvAt struct {
	v *CallerInfo
}

type kvError struct {
	v error
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer for mlog.String()
func (k *kvString) String() string {
	return fmt.Sprintf("%s='%s'", k.k, k.v)
}

// implements fmt.Stringer for mlog.Rune()
func (k *kvRune) String() string {
	return fmt.Sprintf("%s='%c'", k.k, k.v)
}

// implements fmt.Stringer for mlog.Int()
func (k *kvInt) String() string {
	return fmt.Sprintf("%s=%d", k.k, k.v)
}

// implements fmt.Stringer for mlog.Bool()
func (k *kvBool) String() string {
	return fmt.Sprintf("%s=%t", k.k, k.v)
}

// implements fmt.Stringer for mlog.YesNo()
func (k *kvYesNo) String() string {
	s := "No"
	if k.v {
		s = "Yes"
	}
	return fmt.Sprintf("%s=%s", k.k, s)
}

// implements fmt.Stringer for mlog.Byte()
func (k *kvByte) String() string {
	return fmt.Sprintf("%s=0x%02X", k.k, k.v)
}

// implements fmt.Stringer for mlog.At()
func (k *kvAt) String() string {
	return fmt.Sprintf("At=%s", k.v)
}

// implements fmt.Stringer for mlog.At()
func (k *kvError) String() string {
	return fmt.Sprintf("Error=%T=>%s", k.v, k.v)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// log the string key=value pair
func String(key, value string) ILogKeyValuePair {
	return &kvString{key, value}
}

// log the Rune=value
func Rune(key string, value rune) ILogKeyValuePair {
	return &kvRune{key, value}
}

// log the integer key=value pair
func Int(key string, value int) ILogKeyValuePair {
	return &kvInt{key, value}
}

// log the boolean key=value pair
func Bool(key string, value bool) ILogKeyValuePair {
	return &kvBool{key, value}
}

// log the boolean key=value pair as a Yes/No value
func YesNo(key string, value bool) ILogKeyValuePair {
	return &kvYesNo{key, value}
}

// log the byte as a key=value pair
func Byte(key string, value byte) ILogKeyValuePair {
	return &kvByte{key, value}
}

// log the current package/method/function/line location
func At() ILogKeyValuePair {
	return &kvAt{RetrieveCallerInfo(FRAMENR_THIS + 1)}
}

// log an error by its type and message
func Err(err error) ILogKeyValuePair {
	return &kvError{err}
}

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/
/*
func DemoMLog() {
	mlog.SetLevel(mlog.LevelDebug)
	mlog.Info("Useful information")
	mlog.Infof("For your info %c", '⥖')
	mlog.Error("Error happened")
	mlog.Errorf("%s with %d", message, value)
	err := fmt.Errorf("random error")
	mlog.ErrorE(err)
	mlog.Fatal(-5, "Terrible thing happened")
	mlog.DebugT("lazy programmer", mlog.String("Key","value"),
					mlog.Int("Key", 5),
					mlog.Rune("Rune", 'x'),
					mlog.Bool("Key", true),
					mlog.YesNo("Key", false),
					mlog.Err(err),
					mlog.At())

mlog.SetLevel(mlog.LevelDebug)
	mlog.Info("DidimusCommand R/T")
	mlog.Infof("This is %c", '⥖')
	mlog.InfoT("Tagged error", mlog.Rune("Rune", 'E'), mlog.YesNo("Bad", true), mlog.Int("Value", 5), mlog.String("String", "text here"))
}
*/
