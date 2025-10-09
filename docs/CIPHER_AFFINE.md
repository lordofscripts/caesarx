# Affine Cipher

[![Go Reference](https://pkg.go.dev/badge/github.com/lordofscripts/caesarx.svg)](https://pkg.go.dev/github.com/lordofscripts/caesarx)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/lordofscripts/caesarx)](https://github.com/lordofscripts/caesarx/releases/latest)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)
[![Go Report](https://goreportcard.com/badge/github.com/lordofscripts/caesarx)](https://goreportcard.com/report/github.com/lordofscripts/caesarx)

![](./assets/caesarx_header.jpg)


## History

It is unknown when this cipher was invented, but it most certainly was way after Caesar's.
Why I think that? because it is based on a mathematical formula and has certain restrictions
that required certain mathematical knowledge to be in the public domain at the time.

In fact, despite being a former version, The Caesar cipher is a simplified version of
the Affine cipher. More on that in the Encryption/Decryption section.

## Strengths & Weaknesses

Strengths:
* Many combinations of parameters that together with the alphabet length define the corresponding Tabula Recta.

Weaknesses:
* A monoalphabetic substitution cipher.
* Succeptibility to frequency analysis
* If two plaintext-ciphertext character pairs are knnown, it is relatively easy to determine the keys.
* Incorrect parameters cause round-trip (Encode followed by Decode) to fail

Therefore, this one too is not cryptographically secure despite it being defined by a
mathematical function.


## Encryption & Decryption

This cipher is a bit more difficult to do on the field without computing equipment.
Additionally, it requires certain mathematical knowledge (or a tool in lieu of that)
to choose the correct parameters to feed the cipher. Basically, the output is based
on a **linear function**, such a thing is called an affine and therefore the name
of the cipher.

Like plain Caesar, the basic concept is that each letter in the alphabet gets
assigned a shift number starting with 0 for the first letter in the alphabet.

In the Encrypt and Decrypt subsections, we will introduce certain (mathematical)
concepts. The chosen A (or A') parameters are not arbitrary. The consequence of
using incorrect parameters that do not fulfill the conditions, is that you may
not be able to recover the original message during decryption!

### Encrypting with Affine

Given a Primary Reference Alphabet with `N` characters (i.e. 26 in English, 33 in Cyrillic),
the character substitution during the *encryption* process is dictated by the 
following formula (where "%" is the *modulo* operation):

```
		𝓎 = ( 𝓐 × 𝓍 + 𝓑) % 𝓝
```

A and B can be any value but we shall always choose positive integers. But unlike 
Caesar, things get complicated here.

>
> 🌻 “Two integers are said to be coprimes if their Greatest Common Divisor is 1❞ 
>

While B can be any number (no restrictions), the value of A is **not arbitrary**.
For the cipher to work, the value of `A` must be chosen so that `A` and the alphabet
length `N` are *coprimes* (relatively prime or mutually prime). That means that the
only positive number that divides **both** of them is the number 1. Thus, we have
that:

```
		gcd(𝓐,𝓝) = 1
```

That means that for any given built-in alphabet, the allowed set of values for the
"A" equation parameter is finite and limited by the length of the alphabet.

> 💡 Speaking of coprimes, did you know that in mechanical systems, for the teeth of
> two gears ⚙ to mesh together —and avoid grinding each other,— the number of
> teeth of the gears must be coprimes? Just in case you want to design your own
> Enigma-like machine for the Affine cipher.


| LangCode    | Alphabet    |  N    | Valid "A" coefficient values    |
| --- | --- | --- | --- |
|  EN   | English    |  26   |  1, 3, 5, 7, 9, 11, 15, 17, 19, 21, 23, 25   |
|  ES   | Spanish    |  33   |  1, 2, 4, 5, 7, 8, 10, 13, 14, 16, 17, 19, 20, 23, 25, 26, 28, 29, 31, 32   |
|  DE   | German     |  30   |  1, 7, 11, 13, 17, 19, 23, 29   |
|  GR   | Greek      |  24   |  1, 5, 7, 11, 13, 17, 19, 23   |
|  RU   | Cyrillic   |  33   |  1, 2, 4, 5, 7, 8, 10, 13, 14, 16, 17, 19, 20, 23, 25, 26, 28, 29, 31, 32   |


>
> 💡 Did you notice that these built-in alphabets have the following valid A coefficients
> in common: 1, 7, 17, 19 and 23.
>

Since Affine is monoalphabetic, once we have chosen the alphabet (and thus N), we can pick
any "A" coefficient value from the corresponding list. Finally we choose any arbitrary
"B" coefficient.

Once that is done, our encryption function or affine is fully defined. We then take
each of the characters in the chosen alphabet (first letter is shift=0) and feed that
*shift* value as the "x" in the affine. The result "y" is a shift number that we use
to look up the corresponding letter in that same alphabet. This way, for each of the
letters we map it to its sibling. We end up with a monoalphabetic Tabula Recta like 
the one I showed in [Caesar](./CIPHER_CAESAR.md).

After that, the encryption process is easy, just like plain Caesar. First look up
the character from the plain-text message in the top row of the Tabula, and once
found, use the character right below (the result of "y" in the equation).

Did you notice that the Caesar cipher is a special case of the Affine cipher where
"A=1"?

### Decrypting with Affine

For the decryption process we have the Cipher message. If we are the legitimate
recipient, we know the source alphabet and thus "N". We should also know "A"
and "B" because those were the *Affine parameters*. The function or affine
that will help us build the Tabula is as follows:

```
		𝓍 = 𝓐' × ( 𝓎 - 𝓑 ) % 𝓝
	where:
		gcd(𝓐,𝓝) = 1
		𝓐 × 𝓐' % 𝓝 = 1
```

Again, some mathematics needed because we have the "B" and "N" values and we know
"A" as well (which is not in this equation), but we do not have "A'". As it turns 
out, life isn't any easier because A' (A prime) must be the *multiplicative inverse* 
of the value `A modulo N`. This inverse only exists if "A" and "N" are coprimes.

>
> 🌻 “The multiplicative inverse of two number A and A' is such that the result
>	  of their multiplication modulo N is equal to the number 1. For that number
>	  to exist, A and N must be coprimes.❞ 
>

We already know that `gcd(A,N) = 1` because it was a requirement for encryption.
Therefore, given A and N we must calculate A' so that:

```
	( 𝓐 × 𝓐' ) % 𝓝 = 1
```

Sounds complicated, huh? Mathematicians and Lawyers have that in common, and
I by the grace of God am neither of them. This concept is easily explained
and (graphically) depicted with clear examples 
[here](https://www.andreaminini.net/math/modular-inverse-of-a-number).
So don't let that scare you away from playing with the *Affine cipher*.

If we take the built-in English alphabet (N=26), from the previous section
we have the list of valid A coefficients. I will give you the (A,A') pairs
of the first A coefficients in that list: (1,1), (3,9), (5,21), (7,15) and
leave you to determine the rest.

Ready? For doing Affine decryption it is best if you take this brief moment
of joy and accomplishment to use that solved decryption formula to generate
the *transliterated shift value* —and therefore the alphabet character— 
for each of the characters in the Reference Alphabet. This way you have
beforehand, a Tabula for quick conversion of an encoded character to its
decoded version.

## Using it with GoCaesarX

Due to a slightly different implementation, the `affine` command is not
a command alias for `caesarx` but its own standalone application. Here
are the various options.

For encoding/decoding short texts you can specify the message in the
CLI as a free argument:

>
> `affine [options] {parameters} "user text to be processed"
>

However, if you have long data to be processed, it is better to encode/decode
**files** instead by using the `-F` option. The number of free arguments
depends on the operation.

For encrypting the *text* `secret_file.txt` with Affine:

>
> `affine [options] {parameters} -F secret_file.txt
>

The output filename is automatically generated, and in this example it
would be `secret_file_txt.afi`.

For decrypting the *text* `secret_file.txt` with Affine you must also
specify the output (plain text) file as the 2nd free argument:

>
> `affine [options] {parameters} -F secret_file_txt.afi plain_file.txt
>

In my implementation, by default the **Affine cipher** defaults to the English
built-in alphabet *without* slave/secondary/supplementary alphabet. That means
only letters will be processed, the rest (numbers, symbols, spaces, etc.) is
passed through as-is.

Optionally, you can improve the ciphered output by adding a Slave/Secondary
alphabet that contains for example numbers, basic symbols and the space
character. There are several [numeric alphabets](./LANGUAGES.md) to choose
from, but I suggest the *Extended* version.

### Supplementary options

* `-help` produces a short guide for using the application, including all the options.
* `-version` shows the CaesarX global version plus the Affine application version.
* `-demo` a demonstration of a round-trip encryption/decryption
* `-ngram M` (**encryption only**) remove spaces from encrypted output and make groups
   of "M" characters separated by "·". "M" can be an integer between 2 and 5.
* `-F` to indicate the free argument(s) is/are filenames (encrypt/decrypt FILES). Without
	this option, the free argument is a string to be encrypted/decrypted.

### Operational parameters

As you should have learned "A" and "B" are encryption coeficcients. Of these "A" is
constricted to a finite set of values for any given "N", where "N" is the length of
the alphabet that will be used during the cipher. "B" is arbitrary but in effect
it is a modulo of "N" to a great extent.

* `-alpha ID` select the primary alphabet of letters: english (default), spanish/latin,
   german, greek and cyrillic.
* `-A value` set the "A" coefficient. It is NOT an arbitrary value but a coprime of "N".
* `-B value` set the "B" coefficient. It is an arbitrary value.
* `-coprime` gives you a list of valid "A" values for the selected alphabet (N).
* `-N value` set the alphabet length "N" for listing coprimes (only with `-coprime`)
* `-d` indicates a Decryption operation (text follows), else it encrypts.
* `-num value` optionally add a supplementary [alphabet](./LANGUAGES.md) containing numbers,
   space, punctuation, etc. Value can be "N" (none), "A" (Arabic decimals), "H" (Hindi numbers)
   or "E" (extended). Extended contains 0..9, space and 6 basic tell-tale common punctuation.

When Decrypting (`-d`) and encrypting you must always provide both `-A` and `-B` as well
as the text to be decrypted/encrypted.

As you learned earlier, for decrypting "A'" is used instead of "A" but you should always
give "A" because "A'" is its *modular multiplicative inverse* that is derived from
both "A" and "N".

### Sample outputs

The `-tabula` option is good if you want to keep a paper strip of the Rune Transliteration
used for both encryption and decryption with the **selected** parameters.

```
lordofscripts@lisbon:$ affine -tabula -A 7 -B 23
	🔱 Go CaesarX v0.5.0-Beta-0 (C)2025 Didimo Grimaldo 🔱
				 ⚞◕͜ ◕⚟
	☕ Buy me a Coffee? https://www.buymeacoffee/lostinwriting
	=========================================
	Affine ::= A=7 B=23 N=26 A'=15
	0.........1.........2.....
	01234567890123456789012345
	--------------------------
	ABCDEFGHIJKLMNOPQRSTUVWXYZ
	XELSZGNUBIPWDKRYFMTAHOVCJQ
```

Now let's say we want to encrypt using Spanish as alphabet, but without 
encoding numbers, symbols or spaces:

```
lordofscripts@lisbon:$ affine -alpha spanish -A 7 -B 23 "2025 fué un mal año"
	🔱 Go CaesarX v0.5.0-Beta-0 (C)2025 Didimo Grimaldo 🔱
				 ⚞◕͜ ◕⚟
	☕ Buy me a Coffee? https://www.buymeacoffee/lostinwriting
	=========================================
Operation:  Encrypt
Alphabet :  Spanish
Params  M:  Affine ::= A=7 B=23 N=33 A'=19
Params  S:  <nil>
Algorithm:  AFIN ES
Plain    :  2025 fué un mal año
Encoded  :  2025 yfu fo iwb wví
NGram-5  :  2025y·fufoi·wbwví
```

As you see, the numbers and spaces passed through. Depending on your
application, it might leak information, and since the spaces are known,
it may help revealing the keys by guessing the short words within spaces.
Therefore, you *might* want to post-process the encryption with the
`-ngram` option (example in the output).

Now let us enhance the encryption by removing as much tell-tale signs
to make the cipher stronger against attacks:

```
didi@lisbon:$ affine -alpha spanish -A 7 -B 23 -num E "2025 fué un mal año"
	🔱 Go CaesarX v0.5.0-Beta-0 (C)2025 Didimo Grimaldo 🔱
				 ⚞◕͜ ◕⚟
	☕ Buy me a Coffee? https://www.buymeacoffee/lostinwriting
	=========================================
Operation:  Encrypt
Alphabet :  Spanish
Params  M:  Affine ::= A=7 B=23 N=33 A'=19
Params  S:  Affine ::= A=7 B=23 N=17 A'=19
Algorithm:  AFIN ES
Plain    :  2025 fué un mal año
Encoded  :  36378yfu8fo8iwb8wví
NGram-5  :  36378·yfu8f·o8iwb·8wví
```

***
Copyright &copy;2025 Lord of Scripts
