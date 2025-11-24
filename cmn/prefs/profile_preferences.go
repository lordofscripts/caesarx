/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package prefs

import (
	"encoding/hex"
	"fmt"
	"lordofscripts/caesarx"

	"gopkg.in/yaml.v3"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	// @note changing these require updating user YAML profiles!
	itemTypeKey          string = "withKey"
	itemTypeSecret       string = "withSecret"
	itemTypeCaesarium    string = "withCodebook"
	itemTypeCoefficients string = "withCoefficients"
)

/* ----------------------------------------------------------------
 *				M o d u l e   I n i t i a l i z a t i o n
 *-----------------------------------------------------------------*/

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

// ICipherItem is an interface for our items
type ICipherItem interface {
	ItemType() string
}

var _ ICipherItem = (*CaesarModel)(nil)
var _ ICipherItem = (*AffineModel)(nil)
var _ ICipherItem = (*SecretsModel)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// CipherItemContainer holds the ICipherItem for marshaling
type CipherItemContainer struct {
	Item ICipherItem
}

type Recipient struct {
	// The recipient's email address (profile identifier)
	Email string `yaml:"email"`
	// (optional) The recipient's (full) name
	Name string `yaml:"name,omitempty"`
	// The cipher variant (fixed) or NoCipher to rely on a Caesarium
	Variant caesarx.CipherVariant `yaml:"variant"`
	// The 2-letter ISO language code to select primary alphabet
	LangCode string `yaml:"lang_iso"`
	// The name of the chained/slave alphabet, usually symbol,numeric,punctuation types
	Chained string `yaml:"chained,omitempty"`
	// The cipher-specific encryption parameters
	Params CipherItemContainer `yaml:"params"` // using the polymorphic wrapper
}

// Caesar/Didimus/Fibonacci cipher parameter(s) model
type CaesarModel struct {
	// the main encryption key letter (determines shift)
	Key Rune `yaml:"key"`
	// (optional) offset used to derive secondary key in Didimus & Fibonacci
	Offset uint `yaml:"offset,omitempty"`
}

// Affine parameter model
type AffineModel struct {
	A  uint `yaml:"a"`
	B  uint `yaml:"b"`
	Ap uint `yaml:"ap"`
}

// The encryption secret's model for Bellaso & Vigenère ciphers
type SecretsModel struct {
	Secret string `yaml:"secret"`
}

type CaesariumModel struct {
	Mnemonics string `yaml:"mnemonics,omitempty"`
	Entropy   string `yaml:"entropy,omitempty"`
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// a recipient profile with a specific cipher and its parameters
func NewProfileWithCipher(id, name string, cipher caesarx.CipherVariant, langIso, chained string, params ICipherItem) *Recipient {
	wrapped := CipherItemContainer{Item: params}
	return &Recipient{
		Email:    id,
		Name:     name,
		Variant:  cipher,
		LangCode: langIso,
		Chained:  chained,
		Params:   wrapped,
	}
}

// a recipient profile that makes use of a Caesarium (policipher)
func NewProfileWithCaesarium(id, name string, langIso, chained string, params ICipherItem) *Recipient {
	wrapped := CipherItemContainer{Item: params}
	return &Recipient{
		Email:    id,
		Name:     name,
		Variant:  caesarx.NoCipher,
		LangCode: langIso,
		Chained:  chained,
		Params:   wrapped,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (cm *CaesarModel) ItemType() string {
	return itemTypeKey
}

func (cm *CaesarModel) String() string {
	return fmt.Sprintf("CaesarModel Key:%c (optional)Offset:%d", cm.Key, cm.Offset)
}

func (am *AffineModel) ItemType() string {
	return itemTypeCoefficients
}

func (am *AffineModel) String() string {
	return fmt.Sprintf("AffineModel A:%d B:%d Ap:%d", am.A, am.B, am.Ap)
}

func (sm *SecretsModel) ItemType() string {
	return itemTypeSecret
}

func (sm *SecretsModel) String() string {
	return fmt.Sprintf("SecretsModel Secret:%s", sm.Secret)
}

func (csm *CaesariumModel) ItemType() string {
	return itemTypeCaesarium
}

func (csm *CaesariumModel) String() string {
	if len(csm.Entropy) != 0 {
		return fmt.Sprintf("CodebookModel/E:%s", csm.Entropy)
	} else {
		return fmt.Sprintf("CodebookModel/M:%s", csm.Mnemonics)
	}
}

// either it has mnemonics or entropy (hex string). It just checks
// the length as the values are validated when instantiated.
func (csm *CaesariumModel) HasMnemonics() bool {
	return len(csm.Mnemonics) != 0
}

func (csm *CaesariumModel) GetEntropy() []byte {
	if val, err := hex.DecodeString(csm.Entropy); err == nil {
		return val
	} else {
		return []byte{}
	}
}

// implements yaml.Marshaler interface
func (c CipherItemContainer) MarshalYAML() (any, error) {
	return map[string]any{
		"type": c.Item.ItemType(),
		"data": c.Item,
	}, nil
}

// implements yaml.Unmarshaler interface
func (c *CipherItemContainer) UnmarshalYAML(unmarshal func(any) error) error {
	var item struct {
		Type string    `yaml:"type"`
		Data yaml.Node `yaml:"data"`
	}

	if err := unmarshal(&item); err != nil {
		return err
	}

	switch item.Type {
	case itemTypeKey:
		var params CaesarModel
		if err := item.Data.Decode(&params); err != nil {
			return err
		}
		c.Item = &params

	case itemTypeCoefficients:
		var params AffineModel
		if err := item.Data.Decode(&params); err != nil {
			return err
		}
		c.Item = &params

	case itemTypeSecret:
		var params SecretsModel
		if err := item.Data.Decode(&params); err != nil {
			return err
		}
		c.Item = &params

	case itemTypeCaesarium:
		var params CaesariumModel
		if err := item.Data.Decode(&params); err != nil {
			return err
		}
		c.Item = &params

	default:
		return fmt.Errorf("unknown cipher type: %s", item.Type)
	}
	return nil
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/
