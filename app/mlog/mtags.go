/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package mlog

import (
	"fmt"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type ILogKeyValuePair interface {
	fmt.Stringer
}

var _ ILogKeyValuePair = (*kvString)(nil)
var _ ILogKeyValuePair = (*kvRune)(nil)
var _ ILogKeyValuePair = (*kvInt)(nil)
var _ ILogKeyValuePair = (*kvBool)(nil)
var _ ILogKeyValuePair = (*kvYesNo)(nil)
var _ ILogKeyValuePair = (*kvByte)(nil)

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

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (k *kvString) String() string {
	return fmt.Sprintf("%s='%s'", k.k, k.v)
}

func (k *kvRune) String() string {
	return fmt.Sprintf("%s='%c'", k.k, k.v)
}

func (k *kvInt) String() string {
	return fmt.Sprintf("%s=%d", k.k, k.v)
}

func (k *kvBool) String() string {
	return fmt.Sprintf("%s=%t", k.k, k.v)
}

func (k *kvYesNo) String() string {
	s := "No"
	if k.v {
		s = "Yes"
	}
	return fmt.Sprintf("%s=%s", k.k, s)
}

func (k *kvByte) String() string {
	return fmt.Sprintf("%s=0x%02X", k.k, k.v)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func String(key, value string) ILogKeyValuePair {
	return &kvString{key, value}
}

func Rune(key string, value rune) ILogKeyValuePair {
	return &kvRune{key, value}
}

func Int(key string, value int) ILogKeyValuePair {
	return &kvInt{key, value}
}

func Bool(key string, value bool) ILogKeyValuePair {
	return &kvBool{key, value}
}

func YesNo(key string, value bool) ILogKeyValuePair {
	return &kvYesNo{key, value}
}

func Byte(key string, value byte) ILogKeyValuePair {
	return &kvByte{key, value}
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
					mlog.YesNo("Key", false))

mlog.SetLevel(mlog.LevelDebug)
	mlog.Info("DidimusCommand R/T")
	mlog.Infof("This is %c", '⥖')
	mlog.InfoT("Tagged error", mlog.Rune("Rune", 'E'), mlog.YesNo("Bad", true), mlog.Int("Value", 5), mlog.String("String", "text here"))
}
*/
