/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Sample BIP39 renderers. Renders Entroyy and Mnemonic as a string matrix.
 *-----------------------------------------------------------------*/
package bip39

import (
	"encoding/hex"
	"fmt"
	"strings"
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

type IBip39Renderer interface {
	// formats the mnemonic sentence as a matrix
	FormatMnemonic(mnemonic []string) string
	// formats the mnemonic entropy as a matrix
	FormatEntropy(entropy []byte, asHex bool) string
	// formats the seed as a hex string matrix
	FormatSeed(seed []byte, asHex bool) string
}

var _ IBip39Renderer = (*Bip39StringRenderer)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type Bip39StringRenderer struct {
	rowsMnem int
	colsMnem int
	rowsEnt  int
	colsEnt  int
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// An instance of a BIP39 mnemonic & entropy renderer using plain
// strings as a matrix of items.
func NewBip39StringRenderer(mode Bip39Length) IBip39Renderer {
	var rowsM, rowsE, colsE, colsM int = 0, 0, 4, 0
	//	MODE		  MNEMONIC		ENTROPY		  M-Bits 	E-Bits
	//	Bip39Words12: {Table{4, 3}, Table{4, 4}}  128		16
	//	Bip39Words15: {Table{5, 3}, Table{5, 4}}  160 		20
	//	Bip39Words18: {Table{6, 3}, Table{6, 4}}  192 		24
	//	Bip39Words21: {Table{7, 3}, Table{7, 4}}  224 		28
	//	Bip39Words24: {Table{6, 4}, Table{8, 4}}  256 		32
	switch mode {
	case Bip39Words12:
		rowsM = 4
		colsM = 3
		rowsE = 4

	case Bip39Words15:
		rowsM = 5
		colsM = 3
		rowsE = 5

	case Bip39Words18:
		rowsM = 6
		colsM = 3
		rowsE = 6

	case Bip39Words21:
		rowsM = 7
		colsM = 3
		rowsE = 7

	case Bip39Words24:
		rowsM = 6
		colsM = 4
		rowsE = 8
	}

	return &Bip39StringRenderer{
		rowsMnem: rowsM,
		colsMnem: colsM,
		rowsEnt:  rowsE,
		colsEnt:  colsE,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (r *Bip39StringRenderer) String() string {
	return "Bip39StringRenderer"
}

// renders the mnemonic sentence as a table
func (r *Bip39StringRenderer) FormatMnemonic(mnemonic []string) string {
	var sb strings.Builder
	const CR rune = '\n'

	if len(mnemonic) != r.rowsMnem*r.colsMnem {
		sb.WriteString(strings.Join(mnemonic, " "))
		sb.WriteRune(CR)
	} else {
		for row := range r.rowsMnem {
			sb.WriteRune('\t')
			for col := range r.colsMnem {
				// In BIP39 English the maximum word length is 8
				offset := r.colsMnem * row
				sb.WriteString(fmt.Sprintf("%-10s", mnemonic[offset+col]))
			}
			sb.WriteRune(CR)
		}
	}

	return sb.String()
}

// render the entropy as a table
func (r *Bip39StringRenderer) FormatEntropy(entropy []byte, asHex bool) string {
	var sb strings.Builder
	const CR rune = '\n'

	if len(entropy) != r.rowsEnt*r.colsEnt {
		sb.WriteString(fmt.Sprintln(string(entropy)))
		sb.WriteRune(CR)
	} else if asHex {
		sb.WriteString(hex.EncodeToString(entropy))
		sb.WriteRune(CR)
	} else {
		for row := range r.rowsEnt {
			sb.WriteRune('\t')
			for col := range r.colsEnt {
				// In BIP39 English the maximum word length is 8
				offset := r.colsEnt * row
				sb.WriteString(fmt.Sprintf("%5d", entropy[offset+col]))
			}
			sb.WriteRune(CR)
		}
	}

	return sb.String()
}

func (r *Bip39StringRenderer) FormatSeed(seed []byte, asHex bool) string {
	const CR rune = '\n'
	var sb strings.Builder
	var rows int = 4
	var cols int = 16

	if asHex {
		sb.WriteString(hex.EncodeToString(seed))
		sb.WriteRune(CR)
	} else {
		for row := range rows {
			sb.WriteRune('\t')
			for col := range cols {
				offset := cols * row
				sb.WriteString(fmt.Sprintf("%02x", seed[offset+col]))
			}
			sb.WriteRune(CR)
		}

	}

	return sb.String()
}
