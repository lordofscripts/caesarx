# Vigenère Cipher

[![Go Reference](https://pkg.go.dev/badge/github.com/lordofscripts/caesarx.svg)](https://pkg.go.dev/github.com/lordofscripts/caesarx)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/lordofscripts/caesarx)](https://github.com/lordofscripts/caesarx/releases/latest)
[![License: CC BY-NC-ND 4.0](https://img.shields.io/badge/License-CC_BY--NC--ND_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by-nc-nd/4.0/)
[![Go Report](https://goreportcard.com/badge/github.com/lordofscripts/caesarx)](https://goreportcard.com/report/github.com/lordofscripts/caesarx)

![](./assets/caesarx_header.jpg)


## History

The Vigenère *autokey* cipher was created by the 16th century French diplomat Blaise de Vigenère (1523-1596).
It is a further refinement of the [Bellaso](./CIPHER_BELLASO.md) cipher. So he didn't quite invent it,
he improved it. Sadly, historians made the grave mistake of attributing the Bellaso cipher to
Vigenère.

Just like Bellaso's the autokey was unbreakable for a very long time.


## Strengths & Weaknesses

Strengths:
* When used correctly, the Vigenère autokey can produce almost crack-proof ciphers if used correctly.
* Polyalphabetic with a secret that can be as long as the message

Weaknesses:
* It is still a substitution cipher.

## Encryption & Decryption

Encryption and decryption still follow the same workflow as the plain Caesar cipher, that
means you can always help yourself (if you don't have a program) with a Tabula Recta,
or a Caesar disk (as long as you rotate it for every key).

Just like the [Bellaso](./CIPHER_BELLASO.md) cipher, we need a **Secret**, not a (single-letter)
key but a secret. This secret is one or more words, a phrase and the longer and the more
different characters it has the better! In terms of Caesar jargon, each letter of this
*secret* becomes a *key*, that's why it is a polyalphabetic substitution cipher.

Unlike the Bellaso cipher, *the secret is NOT repeated over the input stream*. It is
called *autokey* because the message itself can become part of the secret.

As an example, lets use the [English](./data/english_tabula.txt) built-in alphabet and for 
the sake of simplicity, ignore the presence of the numeric/symbol Slave alphabet. Our 
input stream is composed entirely (for improved security, but not necessary for the 
application) of characters in the alphabet.

**Plain message:** "KISS AT DUSK"
**Secret:** "ADJX"  *ensure all its characters are in the alphabet*

Now, we will explain encryption and decryption separately because the workflow is
slightly different. For encryption we have complete knowledge of the Autokey.
However, for decryption we only know the Secret, that complicates the matter. 
We will omit the spaces for this example.

### Encrypting Vigenère Autokey

We have the input stream "KISSATDUSK" with secret "ADJX". The first thing to do is
compose the sequence of encryption keys this way:

1. Count the number of characters in the input stream
2. Start with laying down the **Secret** word/phrase "ADJX"
3. Next to the secret, lay down the **input message** up and until you have the same amount of characters as the original message.

```
    Plain:    "KISSATDUSK"
	Secret:   "ADJX"
    Autokey : "ADJXKISSAT"
```

Then we iterate through the plain message "KISSATDUSK" and for each character we lookup
the corresponding character (in that same position) in the generated *Autokey*. Did you
notice that we know the **entire** Autokey beforehand? So, that makes it straightforward
to use your favorite tool (Tabula Recta, Caesar disk, etc.) to look them up. We end with
this:

```
    Plain:    "KISSATDUSK"
	Secret:   "ADJX"
    Autokey : "ADJXKISSAT"
	Cipher  : "KLBPKBVMSD"
```

### Decrypting Vigenère Autokey

The workflow is similar to encryption but not the same. Why? because at this point
we only have the *ciphered message* and the *Secret*. 

```
    Cipher:   "KLBPKBVMSD"
	Secret:   "ADJX"
    Autokey : "ADJX-------"
```

Did you notice? we only have the initial part of the *required* autokey! the part
made of the *Secret* but, what about the rest?

First we decrypt as many characters as we have from the Secret. And each **decoded**
character we get, we *append it* to our incomplete autokey. Thus we progressively
build the Autokey as we decode the cipher. Here is the initial part, the known secret:

1. for "K" in the cipher we look up the "A" autokey in the Tabula to decode it and obtain "K". We append this "K" we just decoded to our autokey, which now becomes "ADJXK".
2. Then for the cipher's "L" we use the next autokey letter "D" to decode an "I", again we append it to the progressive autokey. The autokey is now "ADJXKI". By now our decoded message is "KI"
3. We do the same for the next two letters "BP" with known keys "JX" which decode as "SS" and our autokey becomes "ADJXKISS". Our decoded message so far is "KISS"

Did you notice that as we progress our incomplete Autokey looks more and more like the original 
Autokey we used during encryption?

That is an important part of the Vigenère decoding workflow, as we obtain a decoded character, we feed it back to our partial Autokey until we reach the end.

```
    Cipher:   "KLBPKBVMSD"
	Secret:   "ADJX"
    Autokey : "ADJXKISSAT-"
	Decoded : "KISSAT"
```

I leave it as an exercise to complete it using your favorite manual tools. This way you will 
become proficient at using Vigènere on the field without the need for a computer!

### For Best Results

These are suggestions to secure your message, but it is not compulsory.

* Remove from the plain message all characters that are not present in the Primary and Slave Reference alphabets
* To conceal the length of the original message, pad the plain text with random characters.
* Do not use words as the Secret Key, instead use a secret based on random characters
* The secret key should be as long as the message
* Don't use the same secret for other messages, always change the secret
* Protect the secret key with your life

## Using it with GoCaesarX

Use the same CLI options that you would with Bellaso except that for Vigenère
you would:

* Use `-variant vigenere` to select this cipher.
* Add the `-secret '<SECRET>'` option to specify a secret that is a word, or preferably a phrase.
	
Example:
	
```
	caesarx -variant vigenere -alpha english -secret "ADJX" "Kiss at Dusk"
```

Alternatively, you could use:

```
	vigenere -alpha english -secret "ADJX" "Kiss at Dusk"
```
	
The output would look like:
	
```
	🔱 Go CaesarX v0.1.0-Alpha-0 (C)2025 Didimo Grimaldo 🔱
				 ⚞◕͜ ◕⚟
	☕ Buy me a Coffee? https://www.buymeacoffee/lostinwriting
	=========================================
Alphabet :  English
Secret   :  ADJX
Algorithm:  Vigenère
Plain    :  Kiss at Dusk
Encoded  :  Klbp3il#Dulk

	☕ Buy me a Coffee? https://www.buymeacoffee/lostinwriting
```

***
Copyright &copy;2025 Lord of Scripts

