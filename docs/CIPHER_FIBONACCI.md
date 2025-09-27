# Caesar-Fibonacci Cipher

[![Go Reference](https://pkg.go.dev/badge/github.com/lordofscripts/caesarx.svg)](https://pkg.go.dev/github.com/lordofscripts/caesarx)
[![GitHub release (with filter)](https://img.shields.io/github/v/release/lordofscripts/caesarx)](https://github.com/lordofscripts/caesarx/releases/latest)
[![GitHub License](https://img.shields.io/github/license/lordofscripts/caesarx)](https://github.com/lordofscripts/caesarx/blob/master/LICENSE)
[![Go Report](https://goreportcard.com/badge/github.com/lordofscripts/caesarx)](https://goreportcard.com/report/github.com/lordofscripts/caesarx)

![](./assets/caesarx_header.jpg)


## History

Once again no history here! This was another toy cipher I created in an attempt to make
encryption for the masses (family, friends, games) and yet make it more difficult for
the average Mary or the casual thief/overlooker.

I figured, if *Caesar* is fun but too basic (1 key), and *Didimus* is more fun (2 keys), 
why not do some more? I came up with this *Fibonacci Caesar* variant based on a fixed
10-term Fibonacci series.

## Strengths & Weaknesses

Strengths:
* Polyalphabetic substitution cipher based on a 10-term Fibonacci series

Weaknesses:
* Similar weaknesses to Caesar but to a minor extent

## Encryption & Decryption

First of all, what is a Fibonacci series? It is a series of numbers that start with 0
and all subsequent numbers are the sum of the previous two numbers.

The 10-term Fibonacci series is composed of the numbers: `0, 1, 1, 2, 3, 5, 8, 13, 21, and 34`.
We also have the standard *Primary Caesar Key* to which we create a series of keys by
adding the corresponding Fibonacci number to the *shift* of the primary key. Therefore,
if the Primary Caesar Key is "G" (shift=6) the series of keys to use would be 
"G" (6+0), "H" (6+1), "H", "I" (6+2), "J" (6+3), "L" (6+5), "O" (6+8), "T" (0+13), "B" (6+21),
and "O" (6+34 mod 26).

The keen reader certainly notices that for the English alphabet (N=26) the last term of
the Fibonacci series is 34 (bigger than the number of characters in the alphabet) and thus
the result has a Modulo operation with the *length of the primary alphabet*. 

Since there are several built-in alphabets with multiple lengths, that last term for
any given Primary Key would be different —due to the modulo operation— depending on the
primary alphabet.

## Using it with GoCaesarX

To use the Fibonacci Caesar variant, you use the same parameters as in the
plain Caesar invocation, except:

* You give `-variant fibonacci` as parameter this time
* You still need the `key <LETTER>` as in plain Caesar, it is your Primary Key.

As an example: `caesarx -variant fibonacci -alpha spanish -key G "Éste año fué malo"` results in:

```
	🔱 Go CaesarX v0.1.0-Alpha-0 (C)2025 Didimo Grimaldo 🔱
				 ⚞◕͜ ◕⚟
	☕ Buy me a Coffee? https://www.buymeacoffee/lostinwriting
	=========================================
Alphabet :  Spanish
Key      : G (shift=6)
Algorithm:  𝑭ƒ𝓍 (10)
Plain    :  Éste año 2025 fué malo
Encoded  :  Bzám2léb3969$1ñüj$ghqv
NGram-5  :  Bzám2·léb39·69$1ñ·üj$gh·qv

	☕ Buy me a Coffee? https://www.buymeacoffee/lostinwriting
```

***
Copyright &copy;2025 Lord of Scripts
