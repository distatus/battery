// battery
// Copyright (C) 2016-2017,2023 Karol 'Kenji Takahashi' Woźniak
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
// TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package battery provides cross-platform, normalized battery information.
//
// Gives access to a system independent, typed battery state, capacity, charge and voltage values recalculated as necessary to be returned in mW, mWh or V units.
//
// Currently supported systems:
//
//	Linux 2.6.39+
//	OS X 10.10+
//	Windows XP+
//	FreeBSD
//	DragonFlyBSD
//	NetBSD
//	OpenBSD
//	Solaris
package battery

import (
	"fmt"
)

// AgnosticState type enumerates possible battery states, using platform agnostic naming.
type AgnosticState int8

const (
	// Undefined specifies a state that was returned by the controller, but there is no
	// platform agnostic mapping for it.
	// This generally shouldn't happen, if it does consider opening a report for the library
	// (ideally with the contents of `State.Explain()` call as well).
	Undefined AgnosticState = -1
	// Unknown specifies a literal unknown state returned by the controller.
	// This state is also considered the "default", therefore it will be set when an Error
	// was returned from the Get/GetAll call.
	Unknown AgnosticState = iota - 1
	Empty
	Full
	Charging
	Discharging
	// Idle specifies a state where battery is in "capacity saving" mode.
	// It usually means that it sits idle at around 80% charge while power source is plugged in.
	Idle
)

var states = map[AgnosticState]string{
	Undefined:   "Undefined",
	Unknown:     "Unknown",
	Empty:       "Empty",
	Full:        "Full",
	Charging:    "Charging",
	Discharging: "Discharging",
	Idle:        "Idle",
}

func (s AgnosticState) String() string {
	return states[s]
}

type State struct {
	Raw      AgnosticState
	specific string
}

func (s State) Explain() string {
	return s.specific
}

func (s State) String() string {
	return s.Raw.String()
}

func (s State) GoString() string {
	return fmt.Sprintf("%s (%s)", s.Raw, s.specific)
}

// Battery type represents a single battery entry information.
type Battery struct {
	// Current battery state.
	State State
	// Current (momentary) capacity (in mWh).
	Current float64
	// Last known full capacity (in mWh).
	Full float64
	// Reported design capacity (in mWh).
	Design float64
	// Current (momentary) charge rate (in mW).
	// It is always non-negative, consult .State field to check
	// whether it means charging or discharging.
	ChargeRate float64
	// Current voltage (in V).
	Voltage float64
	// Design voltage (in V).
	// Some systems (e.g. macOS) do not provide a separate
	// value for this. In such cases, or if getting this fails,
	// but getting `Voltage` succeeds, this field will have
	// the same value as `Voltage`, for convenience.
	DesignVoltage float64
	// CycleCount charge times.
	CycleCount int
}

func (b *Battery) String() string {
	return fmt.Sprintf("%+v", *b)
}

func get(sg func(idx int) (*Battery, error), idx int) (*Battery, error) {
	b, err := sg(idx)
	return b, wrapError(err)
}

// Get returns battery information for given index.
//
// Note that index taken here is normalized, such that GetAll()[idx] == Get(idx).
// It does not necessarily represent the "name" or "position" a battery was given
// by the underlying system.
//
// If error != nil, it will be either ErrFatal or ErrPartial.
func Get(idx int) (*Battery, error) {
	return get(systemGet, idx)
}

func getAll(sg func() ([]*Battery, error)) ([]*Battery, error) {
	bs, err := sg()
	if errors, ok := err.(Errors); ok {
		nils := 0
		partials := 0
		for i, err := range errors {
			err = wrapError(err)
			if err == nil {
				nils++
			}
			if _, ok := err.(ErrPartial); ok {
				partials++
			}
			errors[i] = err
		}
		if nils == len(errors) {
			return bs, nil
		}
		if nils > 0 || partials > 0 {
			return bs, errors
		}
		return nil, ErrFatal{ErrAllNotNil}
	}
	if err != nil {
		return bs, ErrFatal{err}
	}
	return bs, nil
}

// GetAll returns information about all batteries in the system.
//
// If error != nil, it will be either ErrFatal or Errors.
// If error is of type Errors, it is guaranteed that length of both returned slices is the same and that i-th error coresponds with i-th battery structure.
func GetAll() ([]*Battery, error) {
	return getAll(systemGetAll)
}
