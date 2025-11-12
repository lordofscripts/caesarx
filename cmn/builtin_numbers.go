/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Both Western arabic and Eastern arabic numbers based on the
 * decimal system.
 *-----------------------------------------------------------------*/
package cmn

/* ----------------------------------------------------------------
 *							L o c a l s
 *-----------------------------------------------------------------*/

const (
	ALPHA_NAME_NUMBERS_ARABIC          = "numbers"
	ALPHA_NAME_NUMBERS_EASTERN         = "numbers_east"
	ALPHA_NAME_NUMBERS_ARABIC_EXTENDED = "numbers_ext"

	numbers_DISK         string = "0123456789" // Western Arabic Numerals 10 runes 10 bytes
	numbers_DISK_EASTERN string = "٠١٢٣٤٥٦٧٨٩" // Eastern Arabic Numerals 10 runes 20 bytes

	// The symbols after the numbers are in their ASCII code order!
	numbers_DISK_EXT string = "0123456789 #$%+-@" // Western Numbers & Basic Symbols 17 runes 17 bytes
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

var (
	NUMBERS_DISK         *Alphabet = &Alphabet{"Numbers (West)", numbers_DISK, false, false, true, nil, false, PSO_NUM_DEC}
	NUMBERS_EASTERN_DISK *Alphabet = &Alphabet{"Numbers (East)", numbers_DISK_EASTERN, false, false, true, nil, false, PSO_NUM_HIN}
	NUMBERS_DISK_EXT     *Alphabet = &Alphabet{"Numbers (West) Ext", numbers_DISK_EXT, false, false, true, nil, false, PSO_NUM_DEC_EXT}
)
