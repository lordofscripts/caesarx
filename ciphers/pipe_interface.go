/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package ciphers

import (
	"fmt"
	"lordofscripts/caesarx/app/mlog"
	"lordofscripts/caesarx/cmn"
)

/* ----------------------------------------------------------------
 *							G l o b a l s
 *-----------------------------------------------------------------*/
const (
	PipeNoOp   PipeAction = iota // the piped command does nothing
	PipeEncode                   // the piped command encodes the input
	PipeDecode                   // the piped command decodes the input
)

/* ----------------------------------------------------------------
 *						I n t e r f a c e s
 *-----------------------------------------------------------------*/

var _ IPipe = (*Pipe)(nil)

type IPipe interface {
	WithPipe(command any) error
}

type PipeAction uint8

/* ----------------------------------------------------------------
 *							T y p e s
 *-----------------------------------------------------------------*/

/**
 * The Pipe object does not do anything by itself other than validating
 * the piped command. This object must be embedded into a command that
 * would become 'pipeable'.
 * Examples:
 *	· chain/pipe the output of a Caesar cipher to a Base64 encoder
 *	· pipe the output of a Caesar cipher to a formatter like Trigram
 */
type Pipe struct {
	pipeCmd any
}

/* ----------------------------------------------------------------
 *							C o n s t r u c t o r s
 *-----------------------------------------------------------------*/

/**
 * (Ctor) create an empty pipe that will later be initialized
 * through the WithPipe() method.
 */
func NewEmptyPipe() Pipe {
	return Pipe{pipeCmd: nil}
}

/* ----------------------------------------------------------------
 *							M e t h o d s
 *-----------------------------------------------------------------*/

func (e PipeAction) String() string {
	labels := []string{"PipeNoOp", "PipeEncode", "PipeDecode"}
	return labels[e]
}

/**
 * Fit a pipe at the output of the command so that our output
 * becomes its input.
 */
func (p *Pipe) WithPipe(command any) error {
	var err error = nil
	if p.isFittingPipe(command) {
		p.pipeCmd = command
	} else {
		err = fmt.Errorf("%T is not a fitting pipe, sorry", command)
	}

	return err
}

/**
 * Call this prior to PipeOutput() to avoid errors on submitting
 * data to a closed pipe.
 */
func (p *Pipe) IsPipeOpen() bool {
	return p.pipeCmd != nil
}

/**
 * Perform the named action through the pipe.
 */
func (p *Pipe) PipeOutput(action PipeAction, data string) (string, error) {
	if p.pipeCmd == nil {
		mlog.ErrorT("action on nil pipe passing unchanged", mlog.String("Action", action.String()))
		return data, nil
	}

	if action == PipeNoOp { // pass-through unchanged
		return data, nil
	}

	switch v := (p.pipeCmd).(type) {
	case ICipherCommand:
		switch action {
		case PipeEncode:
			return v.Encode(data)
		case PipeDecode:
			return v.Decode(data)
		}
		panic("invalid PipeAction")

	case cmn.ICommand:
		return v.Execute(data)

	default:
		return data, fmt.Errorf("type %T is not pipeable. How did we get here?", v)
	}
}

// Common factor for checking whether a Pipe candidate fulfills
// either ICommand (post-processing commands) or ICipherCommand.
func (p *Pipe) isFittingPipe(command any) bool {
	accept := false
	if _, ok := command.(cmn.ICommand); ok {
		accept = true
	}
	if _, ok := command.(ICipherCommand); ok {
		accept = true
	}

	return accept
}

/* ----------------------------------------------------------------
 *							F u n c t i o n s
 *-----------------------------------------------------------------*/

// Common factor for checking whether a Pipe candidate fulfills
// either ICommand (post-processing commands) or ICipherCommand.
func CheckCipherPipe(command any) bool { // @audit needed? DEPRECATE
	accept := false
	if _, ok := command.(cmn.ICommand); ok {
		accept = true
	}
	if _, ok := command.(ICipherCommand); ok {
		accept = true
	}
	return accept
}
