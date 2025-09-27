/*
 *  Module name: CaesarX
 *	Author     : Dídimo Emilio Grimaldo Tuñón
 *	Version	   : v1.0
 *	Created    : 16 Jan 2025
 *	Copyright (C)2025 Didimo Grimaldo
 * ------------------------------------------------------------------------
 *   This is all about having fun with ancient ciphers for those who love
 * (or want to get started with) cryptography.
 *   This module implements the Caesar monoalphabetic cipher as well as
 * the Affine cipher. Additionally, it also introduces a couple of new
 * polialphabetic ciphers like Didimus & Fibonacci, and XVIII century
 * ciphers like Bellaso & Vigenère which are also polyalphabetic.
 *   Due to their nature NONE of these ciphers are strong by modern
 * standards. However, they may be useful in certain circles or applications.
 * But above all this is nothing more than an educational experiment.
 *   I felt there was no purpose in re-doing Caesar and variants if they
 * were going to be (severely) limited like most implementations out there
 * that only work with a limited ASCII set (A..Z). Therefore, I went my
 * own way into making sure the time spent on it would serve a purpose.
 *   In addition to the English ASCII alphabet, I added support for other
 * built-in alphabets like Spanish/Latin, German, Greek and Cyrillic
 * alphabets. Most implementations in the wild simply break when fed
 * any of these multi-byte character sets.
 *   The Caesar alphabet was also limited to letters, an the use of spaces
 * and other modern-day tell-tale signs like @ $ % etc. would make it
 * very easy to attack the ciphers. Now, as I said, this is not secure,
 * but this module allows the user to use Slave alphabets to complement
 * the original letter alphabet, so you can also encode numbers, spaces
 * and some symbols.
 * CAESAR CIPHER BASICS
 *   The Caesar cipher consists of two discs or strings. Both have the same
 * order of letters/digits/symbols. To encode, the sender and recipient agree
 * on a positive number that ranges from 0 to the length of the Caesar
 * alphabet (A..Z), thus 26. This "key" or "shift quantity" indicates how
 * many characters to the right the disc is displaced. Zero (0) means no
 * transliteration. But a key of 5 means the "A" becomes "G", and so on:
 *                                   1            2
 *	   Key Number      : 0123 4567 8901 2345 6789 0123 45
 *	   Original (Outer): ABCD EFGH IJKL MNOP QRST UVWX YZ
 *	   Keyed (Inner)   : UVWX YZAB CDEF GHIJ KLMN OPQR ST
 *
 * If Using a device with Fixed outer ring (FOR) and Movable Inner Ring (MIR)
 * then the Caesar cipher key indicates how many steps to the left the INNER
 * disk is moved. A left rotation of 3 is the same as a right rotation of 23.
 */
package caesarx
