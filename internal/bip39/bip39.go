/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Basic BIP39 Mnemonic Sentence Generator (English).
 *-----------------------------------------------------------------*/
package bip39

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256" // for Seed
	"crypto/sha512" // for Seed
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"     // for Seed
	"golang.org/x/text/unicode/norm" // NFKD normalization for Seed
)

/**
Wordlists:
	· Spanish: https://github.com/bitcoin/bips/blob/master/bip-0039/spanish.txt
	· English: https://github.com/bitcoin/bips/blob/master/bip-0039/english.txt
I. Generate the Mnemonic using the BIP39_WORDS list of words with unique first 4-letter combinations
1.1 BIP39 uses either 12-word or 24-word sentence. Most systems use English mnemonics so it is
	advised to keep it in English for global compability.
			CS = ENT / 32 (bits)
			MS = (ENT + CS) / 11 (# of words in mnemonic)

			|  ENT  | CS | ENT+CS |  MS  |
			+-------+----+--------+------+
			|  128  |  4 |   132  |  12  |
			|  160  |  5 |   165  |  15  |
			|  192  |  6 |   198  |  18  |
			|  224  |  7 |   231  |  21  |
			|  256  |  8 |   264  |  24  |
1.2 Select the entropy size (in bits) for your algorithm, this determines the mnemonic
	sentence length (12-24 words).
1.3 Generate the entropy ENT of the selected length in bits
    entropy, _ := hex.DecodeString("066dca1a2bb7e8a1db2832148ce9933eea0f3ac9548d793112d9a95c9407efad")
1.4 Generate the SHA256 hash of the ENTropy and take only the first ENT/32 bits of that hash
	That gives you the ENT+CS bit size and the number of words MS of the mnemonic sentence.
	Those CS bits are APPENDED to the ENTropy bits.
1.5 Divide that ENT+CS bit stream into groups of 11 bits. Each group would result in a
	number between 0..2047. Each of these 11-bit numbers represent an INDEX into the wordlist.
1.6 For each generated index, pick the corresponding BIP39_WORDS word from the Wordlist.

II. Generate the seed
2.1 use PBKDF2 function with Mnemonic sentence in UTF8 NFKD as password, and a salt composed of
	 "mnemonic" + passphrase. If no passphrase is given by the user "" is used. The iteration
	 count is set to 2048 and HMAC-SHA512 is used as the pseudo-random function. The length of
	 the derived key is 512 bits (= 64 bytes).
*/
// https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki
// https://github.com/dongri/go-mnemonic
// https://github.com/vedhavyas/go-mnemonic
// https://github.com/gofika/bip39

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	PBKDF2_ROUNDS int = 2048

	BIP39_WORDS string = "abandon ability able about above absent absorb abstract absurd abuse access accident account accuse achieve acid acoustic acquire across act action actor actress actual adapt add addict address adjust admit adult advance advice aerobic affair afford afraid again age agent agree ahead aim air airport aisle alarm album alcohol alert alien all alley allow almost alone alpha already also alter always amateur amazing among amount amused analyst anchor ancient anger angle angry animal ankle announce annual another answer antenna antique anxiety any apart apology appear apple approve april arch arctic area arena argue arm armed armor army around arrange arrest arrive arrow art artefact artist artwork ask aspect assault asset assist assume asthma athlete atom attack attend attitude attract auction audit august aunt author auto autumn average avocado avoid awake aware away awesome awful awkward axis baby bachelor bacon badge bag balance balcony ball bamboo banana banner bar barely bargain barrel base basic basket battle beach bean beauty because become beef before begin behave behind believe below belt bench benefit best betray better between beyond bicycle bid bike bind biology bird birth bitter black blade blame blanket blast bleak bless blind blood blossom blouse blue blur blush board boat body boil bomb bone bonus book boost border boring borrow boss bottom bounce box boy bracket brain brand brass brave bread breeze brick bridge brief bright bring brisk broccoli broken bronze broom brother brown brush bubble buddy budget buffalo build bulb bulk bullet bundle bunker burden burger burst bus business busy butter buyer buzz cabbage cabin cable cactus cage cake call calm camera camp can canal cancel candy cannon canoe canvas canyon capable capital captain car carbon card cargo carpet carry cart case cash casino castle casual cat catalog catch category cattle caught cause caution cave ceiling celery cement census century cereal certain chair chalk champion change chaos chapter charge chase chat cheap check cheese chef cherry chest chicken chief child chimney choice choose chronic chuckle chunk churn cigar cinnamon circle citizen city civil claim clap clarify claw clay clean clerk clever click client cliff climb clinic clip clock clog close cloth cloud clown club clump cluster clutch coach coast coconut code coffee coil coin collect color column combine come comfort comic common company concert conduct confirm congress connect consider control convince cook cool copper copy coral core corn correct cost cotton couch country couple course cousin cover coyote crack cradle craft cram crane crash crater crawl crazy cream credit creek crew cricket crime crisp critic crop cross crouch crowd crucial cruel cruise crumble crunch crush cry crystal cube culture cup cupboard curious current curtain curve cushion custom cute cycle dad damage damp dance danger daring dash daughter dawn day deal debate debris decade december decide decline decorate decrease deer defense define defy degree delay deliver demand demise denial dentist deny depart depend deposit depth deputy derive describe desert design desk despair destroy detail detect develop device devote diagram dial diamond diary dice diesel diet differ digital dignity dilemma dinner dinosaur direct dirt disagree discover disease dish dismiss disorder display distance divert divide divorce dizzy doctor document dog doll dolphin domain donate donkey donor door dose double dove draft dragon drama drastic draw dream dress drift drill drink drip drive drop drum dry duck dumb dune during dust dutch duty dwarf dynamic eager eagle early earn earth easily east easy echo ecology economy edge edit educate effort egg eight either elbow elder electric elegant element elephant elevator elite else embark embody embrace emerge emotion employ empower empty enable enact end endless endorse enemy energy enforce engage engine enhance enjoy enlist enough enrich enroll ensure enter entire entry envelope episode equal equip era erase erode erosion error erupt escape essay essence estate eternal ethics evidence evil evoke evolve exact example excess exchange excite exclude excuse execute exercise exhaust exhibit exile exist exit exotic expand expect expire explain expose express extend extra eye eyebrow fabric face faculty fade faint faith fall false fame family famous fan fancy fantasy farm fashion fat fatal father fatigue fault favorite feature february federal fee feed feel female fence festival fetch fever few fiber fiction field figure file film filter final find fine finger finish fire firm first fiscal fish fit fitness fix flag flame flash flat flavor flee flight flip float flock floor flower fluid flush fly foam focus fog foil fold follow food foot force forest forget fork fortune forum forward fossil foster found fox fragile frame frequent fresh friend fringe frog front frost frown frozen fruit fuel fun funny furnace fury future gadget gain galaxy gallery game gap garage garbage garden garlic garment gas gasp gate gather gauge gaze general genius genre gentle genuine gesture ghost giant gift giggle ginger giraffe girl give glad glance glare glass glide glimpse globe gloom glory glove glow glue goat goddess gold good goose gorilla gospel gossip govern gown grab grace grain grant grape grass gravity great green grid grief grit grocery group grow grunt guard guess guide guilt guitar gun gym habit hair half hammer hamster hand happy harbor hard harsh harvest hat have hawk hazard head health heart heavy hedgehog height hello helmet help hen hero hidden high hill hint hip hire history hobby hockey hold hole holiday hollow home honey hood hope horn horror horse hospital host hotel hour hover hub huge human humble humor hundred hungry hunt hurdle hurry hurt husband hybrid ice icon idea identify idle ignore ill illegal illness image imitate immense immune impact impose improve impulse inch include income increase index indicate indoor industry infant inflict inform inhale inherit initial inject injury inmate inner innocent input inquiry insane insect inside inspire install intact interest into invest invite involve iron island isolate issue item ivory jacket jaguar jar jazz jealous jeans jelly jewel job join joke journey joy judge juice jump jungle junior junk just kangaroo keen keep ketchup key kick kid kidney kind kingdom kiss kit kitchen kite kitten kiwi knee knife knock know lab label labor ladder lady lake lamp language laptop large later latin laugh laundry lava law lawn lawsuit layer lazy leader leaf learn leave lecture left leg legal legend leisure lemon lend length lens leopard lesson letter level liar liberty library license life lift light like limb limit link lion liquid list little live lizard load loan lobster local lock logic lonely long loop lottery loud lounge love loyal lucky luggage lumber lunar lunch luxury lyrics machine mad magic magnet maid mail main major make mammal man manage mandate mango mansion manual maple marble march margin marine market marriage mask mass master match material math matrix matter maximum maze meadow mean measure meat mechanic medal media melody melt member memory mention menu mercy merge merit merry mesh message metal method middle midnight milk million mimic mind minimum minor minute miracle mirror misery miss mistake mix mixed mixture mobile model modify mom moment monitor monkey monster month moon moral more morning mosquito mother motion motor mountain mouse move movie much muffin mule multiply muscle museum mushroom music must mutual myself mystery myth naive name napkin narrow nasty nation nature near neck need negative neglect neither nephew nerve nest net network neutral never news next nice night noble noise nominee noodle normal north nose notable note nothing notice novel now nuclear number nurse nut oak obey object oblige obscure observe obtain obvious occur ocean october odor off offer office often oil okay old olive olympic omit once one onion online only open opera opinion oppose option orange orbit orchard order ordinary organ orient original orphan ostrich other outdoor outer output outside oval oven over own owner oxygen oyster ozone pact paddle page pair palace palm panda panel panic panther paper parade parent park parrot party pass patch path patient patrol pattern pause pave payment peace peanut pear peasant pelican pen penalty pencil people pepper perfect permit person pet phone photo phrase physical piano picnic picture piece pig pigeon pill pilot pink pioneer pipe pistol pitch pizza place planet plastic plate play please pledge pluck plug plunge poem poet point polar pole police pond pony pool popular portion position possible post potato pottery poverty powder power practice praise predict prefer prepare present pretty prevent price pride primary print priority prison private prize problem process produce profit program project promote proof property prosper protect proud provide public pudding pull pulp pulse pumpkin punch pupil puppy purchase purity purpose purse push put puzzle pyramid quality quantum quarter question quick quit quiz quote rabbit raccoon race rack radar radio rail rain raise rally ramp ranch random range rapid rare rate rather raven raw razor ready real reason rebel rebuild recall receive recipe record recycle reduce reflect reform refuse region regret regular reject relax release relief rely remain remember remind remove render renew rent reopen repair repeat replace report require rescue resemble resist resource response result retire retreat return reunion reveal review reward rhythm rib ribbon rice rich ride ridge rifle right rigid ring riot ripple risk ritual rival river road roast robot robust rocket romance roof rookie room rose rotate rough round route royal rubber rude rug rule run runway rural sad saddle sadness safe sail salad salmon salon salt salute same sample sand satisfy satoshi sauce sausage save say scale scan scare scatter scene scheme school science scissors scorpion scout scrap screen script scrub sea search season seat second secret section security seed seek segment select sell seminar senior sense sentence series service session settle setup seven shadow shaft shallow share shed shell sheriff shield shift shine ship shiver shock shoe shoot shop short shoulder shove shrimp shrug shuffle shy sibling sick side siege sight sign silent silk silly silver similar simple since sing siren sister situate six size skate sketch ski skill skin skirt skull slab slam sleep slender slice slide slight slim slogan slot slow slush small smart smile smoke smooth snack snake snap sniff snow soap soccer social sock soda soft solar soldier solid solution solve someone song soon sorry sort soul sound soup source south space spare spatial spawn speak special speed spell spend sphere spice spider spike spin spirit split spoil sponsor spoon sport spot spray spread spring spy square squeeze squirrel stable stadium staff stage stairs stamp stand start state stay steak steel stem step stereo stick still sting stock stomach stone stool story stove strategy street strike strong struggle student stuff stumble style subject submit subway success such sudden suffer sugar suggest suit summer sun sunny sunset super supply supreme sure surface surge surprise surround survey suspect sustain swallow swamp swap swarm swear sweet swift swim swing switch sword symbol symptom syrup system table tackle tag tail talent talk tank tape target task taste tattoo taxi teach team tell ten tenant tennis tent term test text thank that theme then theory there they thing this thought three thrive throw thumb thunder ticket tide tiger tilt timber time tiny tip tired tissue title toast tobacco today toddler toe together toilet token tomato tomorrow tone tongue tonight tool tooth top topic topple torch tornado tortoise toss total tourist toward tower town toy track trade traffic tragic train transfer trap trash travel tray treat tree trend trial tribe trick trigger trim trip trophy trouble truck true truly trumpet trust truth try tube tuition tumble tuna tunnel turkey turn turtle twelve twenty twice twin twist two type typical ugly umbrella unable unaware uncle uncover under undo unfair unfold unhappy uniform unique unit universe unknown unlock until unusual unveil update upgrade uphold upon upper upset urban urge usage use used useful useless usual utility vacant vacuum vague valid valley valve van vanish vapor various vast vault vehicle velvet vendor venture venue verb verify version very vessel veteran viable vibrant vicious victory video view village vintage violin virtual virus visa visit visual vital vivid vocal voice void volcano volume vote voyage wage wagon wait walk wall walnut want warfare warm warrior wash wasp waste water wave way wealth weapon wear weasel weather web wedding weekend weird welcome west wet whale what wheat wheel when where whip whisper wide width wife wild will win window wine wing wink winner winter wire wisdom wise wish witness wolf woman wonder wood wool word work world worry worth wrap wreck wrestle wrist write wrong yard year yellow you young youth zebra zero zone zoo"
)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

// BIP39 Mnemonic Sentence Generator (12/15/18/21/24 words).
type Bip39 struct {
	SentenceLength int // size of Mnemonic sentence length (12|24)
	Entropy        uint16
	ChecksumBits   uint8
	separator      rune
	mnemonic       []string
	entropy        []byte
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// (Constructor) instantiates a BIP39 Mnemonic Sentence Generator
func NewBip39(length int, separator rune) *Bip39 {
	var ent uint16
	var cs uint8
	switch length {
	case 12:
		ent = 128
		cs = 4
	case 15:
		ent = 160
		cs = 5
	case 18:
		ent = 192
		cs = 6
	case 21:
		ent = 224
		cs = 7
	case 24:
		ent = 256
		cs = 8
	default:
		return nil
	}

	return &Bip39{
		SentenceLength: length,
		Entropy:        ent,
		ChecksumBits:   cs,
		separator:      separator,
		mnemonic:       nil,
		entropy:        nil,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer by returning the last-generated mnemonic
// sentence as a string
func (b *Bip39) String() string {
	return strings.Join(b.mnemonic, string(b.separator))
}

// Generate a mnemonic sentence of the selected word quantity
// specified in the constructor. The entropy is calculated internally.
func (b *Bip39) GenerateMnemonic() ([]string, error) {
	// Generate random entropy targeting a specific sentence length
	entropy, err := b.generateEntropy(int(b.Entropy))
	if err != nil {
		return nil, err
	}

	return b.generateMnemonic(entropy)
}

// Generate a mnemonic sentence of the selected word quantity
// specified in the constructor but using the provided entropy.
func (b *Bip39) GenerateMnemonicFromEntropy(entropy []byte) ([]string, error) {
	if err := b.validateEntropy(entropy); err != nil {
		return nil, err
	}

	return b.generateMnemonic(entropy)
}

// find out the entropy value used to generate the provided mnemonic sentence
func (b *Bip39) EntropyFromMnemonic(sentence []string) ([]byte, error) {
	// validate sentence length
	vbl := []int{12, 15, 18, 21, 24}
	csl := []int{4, 5, 6, 7, 8}
	size := len(sentence)
	if !slices.Contains(vbl, size) {
		return nil, fmt.Errorf("invalid mnemonic sentence length: %d", size)
	}

	// official English word list
	const WORD_SEP string = " "
	wordList := strings.Split(BIP39_WORDS, WORD_SEP)

	// for generating the ENT+CS bit stream
	var bitStream string = ""
	// generate indices
	for _, word := range sentence {
		// find the word index in the BIP39 English word list
		index := slices.Index(wordList, word)
		// update the binary representation of ENT+CS
		bitStream = bitStream + fmt.Sprintf("%.11b", index)
	}
	// remove CS
	csIdx := slices.Index(vbl, size)
	csLen := csl[csIdx]
	// ENT
	entBits := bitStream[0 : len(bitStream)-csLen]
	// CS
	csBits := bitStream[len(bitStream)-csLen:]
	entropy, err := binaryStringToBytes(entBits)
	if err != nil {
		return nil, err
	}

	// generate the hash of the recovered ENTropy
	hash := b.generateHash(entropy)
	// get the significant bits of that hash for that BIP39 sentence length
	hashBits := bytesToBinaryString(hash)[0:csLen]
	// compare
	if hashBits != csBits {
		return nil, fmt.Errorf("the entropy checksums do not match %s != %s", hashBits, csBits)
	}

	return entropy, nil
}

// the last used entropy slice
func (b *Bip39) GetEntropy() []byte {
	return b.entropy
}

// Generate a cryptographic seed based on the Mnemonic Sentence
// and Passphrase strings.
func (b *Bip39) ToSeed(mnemonic, passphrase string) []byte {
	// normalize mnemonic
	mnemonic = norm.NFKD.String(mnemonic)
	// normalize passphrase
	passphrase = norm.NFKD.String(passphrase)
	// compose passphrase
	passphrase = "mnemonic" + passphrase
	mnemonicBytes := []byte(mnemonic)
	passphraseBytes := []byte(passphrase)
	// PBKDF2
	const KEY_LEN_BYTES = 64
	seed := DeriveKey(mnemonicBytes, passphraseBytes, 2048, KEY_LEN_BYTES)
	return seed
}

// Generate a cryptographic seed based on the Mnemonic Sentence
// and Passphrase strings and return the seed as a Hexadecimal string.
func (b *Bip39) ToSeedHex(mnemonic, passphrase string) string {
	binSeed := b.ToSeed(mnemonic, passphrase)
	return hex.EncodeToString(binSeed)
}

/* ----------------------------------------------------------------
 *				P r i v a t e	M e t h o d s
 *-----------------------------------------------------------------*/

// Generate a mnemonic sentence of the selected word quantity
// specified in the constructor.
func (b *Bip39) generateMnemonic(entropy []byte) ([]string, error) {
	var bitStream string = ""
	// ENT
	bitStream = bytesToBinaryString(entropy)

	// generate the entropy's SHA256
	sha256 := b.generateHash(entropy)
	// ENT+CS
	csBits := bytesToBinaryString(sha256)[0:b.ChecksumBits]
	bitStream = bitStream + csBits
	// Separate ENT+CS into 11-bit groups (each group represents 0..2047)
	const GROUP_SEP rune = '*'
	const WORD_SEP string = " "
	bitStream = GroupBySize(bitStream, 11, GROUP_SEP)
	// official English word list
	wordList := strings.Split(BIP39_WORDS, WORD_SEP)
	// generate indexes of 0..2047
	groupsSlice := strings.Split(bitStream, string(GROUP_SEP))
	// prepare the wordlist
	b.mnemonic = make([]string, b.SentenceLength)
	for i, groupStr := range groupsSlice {
		// convert the group (in binary number base) to integer
		if num, err := strconv.ParseInt(groupStr, 2, 0); err != nil {
			return nil, err
		} else {
			b.mnemonic[i] = wordList[num]
		}
	}

	b.entropy = entropy
	return b.mnemonic, nil
}

// validates an entropy's bit size
func (b *Bip39) validateEntropy(entropy []byte) error {
	ebl := []int{16, 20, 24, 28, 32}
	if !slices.Contains(ebl, len(entropy)) {
		return fmt.Errorf("not a valid entropy size: %d-bits", len(entropy)*8)
	}
	return nil
}

// generate an entropy of the selected bit length
func (b *Bip39) generateEntropy(bitLen int) ([]byte, error) {
	vbl := []int{128, 160, 192, 224, 256}
	if !slices.Contains(vbl, bitLen) {
		return nil, fmt.Errorf("not a valid entropy bit length: %d", bitLen)
	}
	// Generate random entropy
	entropy := make([]byte, bitLen/8)
	_, err := rand.Read(entropy)
	return entropy, err
}

// generates a SHA256 hash of an entropy
func (b *Bip39) generateHash(entropy []byte) []byte {
	hash := sha256.New()
	hash.Write(entropy)
	return hash.Sum(nil)
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Convert the specified byte slice to a binary string.
func bytesToBinaryString(slice []byte) string {
	// Convert each byte to its bits representation as string
	var strBuff bytes.Buffer
	for _, b := range slice {
		strBuff.WriteString(fmt.Sprintf("%.8b", b))
	}

	return strBuff.String()
}

// Convert the specified binary string to a byte slice.
func binaryStringToBytes(binStr string) ([]byte, error) {
	// Length of the binary string shall be multiple of 8
	if (len(binStr) % 8) != 0 {
		return nil, errors.New("the specified binary string is not valid")
	}

	// Create slice
	slice := make([]byte, 0, len(binStr)/8)

	// Split the string into groups of 8-bit and convert each of them to byte
	for i := 0; i < len(binStr); i += 8 {
		// Convert current byte
		byteStrBin := binStr[i : i+8]
		byteVal, err := strconv.ParseInt(byteStrBin, 2, 16)
		// Stop if conversion error
		if err != nil {
			return nil, err
		}
		// Append new byte
		slice = append(slice, byte(byteVal))
	}

	return slice, nil
}

// Take a string s and divide it in N groups separated by char.
func GroupBySize(s string, n uint, char rune) string {
	if n == 0 {
		return s
	}

	var buffer bytes.Buffer
	var n1 = int(n - 1)
	var l1 = len(s) - 1 // we are dealing only with ASCII chars
	letters := []rune(s)
	for i, rune := range letters {
		buffer.WriteRune(rune)
		if i%int(n) == n1 && i != l1 {
			buffer.WriteRune(char)
		}
	}

	return buffer.String()
}

// derive a cryptographic key using PBKDF2 with HMAC-SHA512.
// For Bip39 use 2048 iterations (PBKDF2_ROUNDS) and 64-byte
// key length.
func DeriveKey(password, salt []byte, iter, keyLen int) []byte {
	return pbkdf2.Key(password, salt, iter, keyLen, sha512.New)
}

/*
func DemoBip39() {
	// generate the 24-word mnemonic sentence
	bip39 := NewBip39(24, '·')
	list, err := bip39.GenerateMnemonic()
	if err != nil {
		fmt.Println("Generate Error", err)
	}
	fmt.Println("Length", len(list))
	fmt.Println("Mnemonic", bip39.String())

	// retrieve the entropy used to generate that list
	_, err = bip39.EntropyFromMnemonic(list)
	if err != nil {
		fmt.Println("Verify Error", err)
	}

	seed := bip39.ToSeed(bip39.String(), "hello")
	fmt.Println("Seed", seed)
}
*/
