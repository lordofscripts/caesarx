package ciphers

import (
	"fmt"
	"lordofscripts/caesarx/cmn"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	GROUP_CHAR rune = '·'
)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type grouping struct {
	groupSize int
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// NewGroupingCommand(4) = Quartets()
// NewGroupingCommand(5) = Quintets()
func NewGroupingCommand(qty int) *grouping {
	return &grouping{qty}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (g *grouping) Execute(input string) (string, error) {
	switch g.groupSize {
	case 3:
		return cmn.Trigram(input, GROUP_CHAR), nil
	case 4:
		return cmn.Quartets(input, GROUP_CHAR), nil
	case 5:
		return cmn.Quintets(input, GROUP_CHAR), nil
	}

	return "", fmt.Errorf("invalid grouping size=%d", g.groupSize)
}
