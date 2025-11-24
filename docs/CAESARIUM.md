# The Caesarium

Caesar cipher and derivatives taken to the next level, not only by the
enhanced algorithms of *CaesarX* but also through the addition of the
*Caesarium*, a codebook to follow the Enigma cipher tradition.

## ✨ Features of the Caesarium

- **Built-in Alphabets**: Can build a codebook for any built-in language.
- **Renderers**: Can support multiple output renderers such as console, HTML, etc.
- **Recoverable**: When necessary, it can generate recoverable codebooks
- **BIP39 support**: You can use 12-word mnemonic phrases following BIP39 specifications.
- **Fully random**: By default generates fully random one-use codebooks
- **User Profiles**: Supports user-profiles (recipients)
- **Monthy schedule**: Generate a random codebook for any month
- **Yearly schedule**: Generate a random codebook for an entire year

## Security

Let me stress out that none of the Caesar variant are very secure for
modern day applications. But as stated in the documentation, there are
some modern-day cases where it is *secure enough for your purpose*, be
it games, friends or family.

Having said that, if you decide to go the CaesarX way, you can have
codebooks to share with your communication partner. Security is as good
as the weakest link, so ensure your codebook remains secret and share
it in a secure manner as well.

## Recoverable or not?

The best are the one-time codebooks. These are *non-recoverable* and are
generated using a *true cryptographically random number generator* to
derive ciphers and parameters `crypto/rand`. Generate it, deliver to your party and
secure it at both ends.

In some cases, such as game scenarios or say you prepared a **Last Will
and Testament** with some Caesar-encrypted data and want to *sort of*
ensure the codebook can be generated at any time by the recipients 
(your survivors), then it is reasonable to use *recoverable codebooks*.

Recoverable codebooks can be re-generated through the application or
website at any time provided that the user specifies the same *recovery
phrase* used during the genesis (the `-bip39` CLI flag). So, 
provided the application or website exists, your survivor can reveal
the secret you left for them 15 years after your departure (as an
example). However, to ensure codebooks can be recovered, the internal
codebook generation algorithm use a random number generator that is NOT
cryptographically secure: `math/rand`. This alternative uses the same
BIP39 specification to generate recovery mnemonic sentences used to
create cryptocurrency wallets, except here it is for the codebook.

---

## 📋 Using the Caesarium

Now you can exchange messages with friends and family or game peers using a codebook. For example, for any given built-in language or
supported cipher (Caesar, Didimus, Fibonacci, Bellaso, Vigenère, Affine), you can generate codebooks. The idea is that you and your
communication partner agree on a codebook so that instead of agreeing
on a cipher or key every time, you can automatically choose the
encryption settings using the codebook. The same system used with the
Enigma machine.

### Cipher schedule for a whole year

If you want to generate a random list of ciphers for use throughout the
year, the following command would generate a list of one cipher per
month for 2026:

```
lordofscrips@bitbucket:$ codebook -date 2026
    ┏━╸┏━┓┏━╸┏━┓┏━┓┏━┓╻╻ ╻┏┳┓
    ┃  ┣━┫┣╸ ┗━┓┣━┫┣┳┛┃┃ ┃┃┃┃
    ┗━╸╹ ╹┗━╸┗━┛╹ ╹╹┗╸╹┗━┛╹ ╹
    By Lord-of-Scripts™
                            Cipher Schedule for 2026                            
┌──────────────────────────────────────────────────────────────────┐
│January    February   March      April      May        June       │
│Caesar     Affine     Fibonacci  Affine     Vigenere   Affine     │
├──────────────────────────────────────────────────────────────────┤
│July       August     September  October    November   December   │
│Affine     Affine     Didimus    Didimus    Affine     Caesar     │
└──────────────────────────────────────────────────────────────────┘
```

### Schedule for a month

You may generate a codebook for any given month of the year. This codebook would specify, for each day of the month, you agree
on a given cipher and the codebook gives you that cipher's random 
parameter for each day of that month. For example, this will generate a
codebook for *Didimus* cipher for February 2025:

```
lordofscrips@bitbucket:$ codebook -variant didimus -date 2025-02
                                   Caesarium                                    
                               you@bitbucket.com                                
                                                                2025-November-11
                      Didimus Daily Settings for Jun-2025                       
                                 English (N=26)                                 
┌─────────────────────────────────────────────────────────────────────────┐
│                                  June                                   │
├─────────────────────────────────────────────────────────────────────────┤
│ Day    │  Key Shift  Offset                    Notes                    │
├────────┼───────────────────┼────────────────────────────────────────────┤
│   1 Sun│    F    5    +25  │                                            │
│   2 Mon│    P   15    +17  │                                            │
│   3 Tue│    Z   25     +8  │                                            │
│   4 Wed│    Q   16    +16  │                                            │
│   5 Thu│    D    3     +7  │                                            │
│   6 Fri│    M   12     +9  │                                            │
│   7 Sat│    N   13    +14  │                                            │
│   8 Sun│    B    1    +18  │                                            │
│   9 Mon│    J    9    +15  │                                            │
│  10 Tue│    X   23    +16  │                                            │
│  11 Wed│    L   11    +13  │                                            │
│  12 Thu│    P   15    +13  │                                            │
│  13 Fri│    D    3     +1  │                                            │
│  14 Sat│    J    9    +22  │                                            │
│  15 Sun│    L   11     +9  │                                            │
│  16 Mon│    Q   16     +1  │                                            │
│  17 Tue│    F    5    +24  │                                            │
│  18 Wed│    E    4    +20  │                                            │
│  19 Thu│    W   22     +2  │                                            │
│  20 Fri│    U   20    +22  │                                            │
│  21 Sat│    Q   16    +17  │                                            │
│  22 Sun│    V   21    +22  │                                            │
│  23 Mon│    Z   25     +6  │                                            │
│  24 Tue│    L   11     +5  │                                            │
│  25 Wed│    Q   16    +23  │                                            │
│  26 Thu│    U   20     +3  │                                            │
│  27 Fri│    O   14    +10  │                                            │
│  28 Sat│    C    2     +4  │                                            │
│  29 Sun│    S   18    +21  │                                            │
│  30 Mon│    R   17    +13  │                                            │
└────────┴───────────────────┴────────────────────────────────────────────┘
  · Alphabet (EN) has 26 runes and 26 bytes
  · Alphabet Runes: ABCDEFGHIJKLMNOPQRSTUVWXYZ
  · The Shift column is the Caesar shift for the given Key
  · The Offset applies to the secondary key relative to the main Key
  · The Offset is required for Didimus, optional for Fibonacci
``` 

Each monthly schedule table is formatted according to the parameters or
settings needed for the cipher chosen for that month.

### Full Year Codebook

If you want to generate an entire year's worth of codebooks that will
include:

- The cover page
- The monthly schedule of ciphers
- The 12-month daily schedules with settings for that month's cipher

```
lordofscrips@bitbucket:$ codebook -date 2026 -full
```

## 📌 Other CLI options for the Caesarium 


- `-bip39` generates and use a self-generated random mnemonic recovery phrase 
   according to BIP39
- `-title STRING` The title (defaults to "Caesarium")
- `-variant STRING` Select the cipher variant (defaults to "caesar").
- `-for STRING` the recipient (default to "you@bitbucket.com")
- `-date DATE` the date in ISO format, it is either a 4-digit year `2025`, or 
  the year followed by the month number `2025-03`.
- `-alpha STRING` Primary built-in alphabet name (defaults to "english"). 
  See [Languages](./LANGUAGES.md).

---
Updated 20 November 2025.