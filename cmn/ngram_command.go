/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package cmn

import (
	"fmt"
	"strings"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/
var _ ICommand = (*NgramCmd)(nil)

type NgramCmd struct {
	length uint8
	sep    rune
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

func NewNgramFormatter(length uint8, separator rune) *NgramCmd {
	if length == 1 || length > 5 {
		panic("nGramFormatter only supports 0,2,3,4 & 5")
	}

	return &NgramCmd{length, separator}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (f *NgramCmd) Execute(s string) (string, error) {
	var err error = nil
	var result string

	s = strings.Replace(s, " ", "", -1)
	switch f.length {
	case 0:
		result = s

	case 2:
		result = Bigram(s, f.sep)

	case 3:
		result = Trigram(s, f.sep)

	case 4:
		result = Quartets(s, f.sep)

	case 5:
		result = Quintets(s, f.sep)

	default:
		result = s
		err = fmt.Errorf("unsupported grouping")
	}

	return result, err
}

/**
 * A friendly representation of the command in the pipe
 */
func (f *NgramCmd) String() string {
	return fmt.Sprintf("NGram(%d%c)", f.length, f.sep)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						M A I N | E X A M P L E
 *-----------------------------------------------------------------*/

/*
func DemoNgramFormatter() {
	bigram,_ := NewFormatter(2, '·').Execute("ABCDEFGH") // AB·CD·EF·GH
	trigram,_ := NewFormatter(3, '·').Execute("ABCDEFGH") // ABC·DEF·GH
	quartet,_ := NewFormatter(4, '·').Execute("ABCDEFGH") // ABCD·EFGH
	quintet,_ := NewFormatter(5, '·').Execute("ABCDEFGH") // ABCDE·FGH
}
*/
