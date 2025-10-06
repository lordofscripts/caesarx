/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * (ASCII) The official Binary alphabet. 256 runes 384 bytes.
 * The ASCII table as a binary alphabet where extended ASCII codes
 * (128..255) are 2-byte runes.
 *-----------------------------------------------------------------*/
package cmn

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

var (
	binary_DISK string = string(makeRuneASCII()) // Binary  256 runes 384 bytes
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	ALPHA_NAME_BINARY = "binary"
)

var (
	BINARY_DISK *Alphabet = &Alphabet{
		Name:        ALPHA_NAME_BINARY,
		Chars:       binary_DISK,
		Foreign:     false,
		Unicode:     true,
		OnlySymbols: true,
		specialCase: nil,
		langCode:    "BX", // does not conflict with any ISO-639 language code
	}
)

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

func makeByteASCII() []byte {
	ASCII := make([]byte, 256)
	for i := 0; i < len(ASCII); i++ {
		ASCII[i] = byte(i) // 128..255 are 2-byte runes
	}
	return ASCII
}

func makeRuneASCII() []rune {
	ASCII := make([]rune, 256)
	for i := 0; i < len(ASCII); i++ {
		ASCII[i] = rune(i) // 128..255 are 2-byte runes
	}
	return ASCII
}
