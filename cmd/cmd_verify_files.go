/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 * Concrete implementation of ICommander to verify two files for equality.
 *-----------------------------------------------------------------*/
package cmd

import (
	"fmt"
	"lordofscripts/caesarx/cmn"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/

const (
	HashCRC64 HashType = iota
	HashMD5
	HashSHA256
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ ICommander = (*VerifyFilesCommand)(nil)

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

type HashType uint

type VerifyFilesCommand struct {
	CommanderBase
	filenameA string
	filenameB string
	hash      HashType
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

// creates an instance of VerifyFilesCommand to verify fileA with fileB
// using the with file hash/digest.
func NewVerifyFilesCommand(fileA, fileB string, with HashType) *VerifyFilesCommand {
	return &VerifyFilesCommand{
		CommanderBase{},
		fileA,
		fileB,
		with,
	}
}

// creates an instance of VerifyFilesCommand to verify fileA
// using the with file hash/digest.
func NewVerifyFileCommand(fileA string, with HashType) *VerifyFilesCommand {
	return &VerifyFilesCommand{
		CommanderBase{},
		fileA,
		"",
		with,
	}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

// implements fmt.Stringer giving general info about the command
func (vc *VerifyFilesCommand) String() string {
	var hashName string
	switch vc.hash {
	case HashCRC64:
		hashName = "CRC64"

	case HashMD5:
		hashName = "MD5"

	case HashSHA256:
		hashName = "SHA256"

	default:
		hashName = ""
	}

	n := 1
	if vc.hasSecondFile() {
		n = 2
	}

	return fmt.Sprintf("Verify %d file(s) using %s", n, hashName)
}

// execute the command with the input provided in the constructor.
// returns nil on success. Check command output with GetOutput()
func (vc *VerifyFilesCommand) Execute() error {
	var err error = nil

	funcVerify := func(val1, val2 any) bool {
		var areEqual bool
		switch val1.(type) {
		case uint64:
			valU1 := val1.(uint64)
			valU2 := val2.(uint64)
			areEqual = valU1 == valU2

		case string:
			valU1 := val1.(string)
			valU2 := val2.(string)
			areEqual = valU1 == valU2

		}

		if areEqual {
			vc.Output("· OK both files are equal\n")
		} else {
			vc.Output("· FAIL files are different!\n")
		}

		return areEqual
	}

	switch vc.hash {
	case HashCRC64:
		if hash1, errA := cmn.CalculateFileCRC64(vc.filenameA); errA == nil {
			vc.Output("·%16X %s\n", hash1, vc.filenameA)
			if vc.hasSecondFile() {
				if hash2, errB := cmn.CalculateFileCRC64(vc.filenameB); errB == nil {
					vc.Output("·%16X %s\n", hash2, vc.filenameB)
					funcVerify(hash1, hash2)
				} else {
					err = errB
				}
			}
		} else {
			err = errA
		}

	case HashMD5:
		if hash1, errA := cmn.CalculateFileMD5(vc.filenameA); errA == nil {
			vc.Output("·%32s %s\n", hash1, vc.filenameA)
			if vc.hasSecondFile() {
				if hash2, errB := cmn.CalculateFileMD5(vc.filenameB); errB == nil {
					vc.Output("·%32s %s\n", hash2, vc.filenameB)
					funcVerify(hash1, hash2)
				} else {
					err = errB
				}
			}
		} else {
			err = errA
		}

	case HashSHA256:
		if hash1, errA := cmn.CalculateFileSHA256(vc.filenameA); errA == nil {
			vc.Output("·%64s %s\n", hash1, vc.filenameA)
			if vc.hasSecondFile() {
				if hash2, errB := cmn.CalculateFileSHA256(vc.filenameB); errB == nil {
					vc.Output("·%64s %s\n", hash2, vc.filenameB)
					funcVerify(hash1, hash2)
				} else {
					err = errB
				}

			}
		} else {
			err = errA
		}

	}

	return err
}

// check whether the command will be comparing the hash of two files
// or just computing the hash of one for reference.
func (vc *VerifyFilesCommand) hasSecondFile() bool {
	return len(vc.filenameB) > 0
}
