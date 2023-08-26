// battery
// Copyright (C) 2016-2017 Karol 'Kenji Takahashi' Woźniak
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

package battery

import (
	"errors"
	"testing"
)

func TestErrPartial(t *testing.T) {
	cases := []struct {
		in    ErrPartial
		str   string
		isnil bool
		nonil bool
	}{
		{ErrPartial{}, "{}", true, false},
		{ErrPartial{Full: errors.New("t1")}, "{Full:t1}", false, false},
		{ErrPartial{State: errors.New("t2"), Full: errors.New("t3")}, "{State:t2 Full:t3}", false, false},
		{ErrPartial{State: errors.New("t4"), Current: errors.New("t5"), Full: errors.New("t6"), Design: errors.New("t7"), ChargeRate: errors.New("t8"), Voltage: errors.New("t9"), DesignVoltage: errors.New("t10")}, "{State:t4 Current:t5 Full:t6 Design:t7 ChargeRate:t8 Voltage:t9 DesignVoltage:t10}", false, true},
	}

	for i, c := range cases {
		str := c.in.Error()
		isnil := c.in.isNil()
		nonil := c.in.noNil()

		if str != c.str {
			t.Errorf("%d: %v != %v", i, str, c.str)
		}
		if isnil != c.isnil {
			t.Errorf("%d: %v != %v", i, isnil, c.isnil)
		}
		if nonil != c.nonil {
			t.Errorf("%d: %v != %v", i, nonil, c.nonil)
		}
	}
}

func TestErrors(t *testing.T) {
	cases := []struct {
		in  Errors
		str string
	}{
		{Errors{nil}, "[]"},
		{Errors{ErrPartial{}}, "[{}]"},
		{Errors{ErrFatal{errors.New("t1")}}, "[Could not retrieve battery info: `t1`]"},
		{Errors{ErrPartial{Full: errors.New("t2")}, ErrFatal{errors.New("t3")}}, "[{Full:t2} Could not retrieve battery info: `t3`]"},
		{Errors{ErrPartial{Full: errors.New("t4")}, ErrPartial{Current: errors.New("t5")}}, "[{Full:t4} {Current:t5}]"},
	}

	for i, c := range cases {
		str := c.in.Error()

		if str != c.str {
			t.Errorf("%d: %v != %v", i, str, c.str)
		}
	}
}

func BenchmarkErrors_Error(b *testing.B) {
	e := getErrors()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Error()
	}
}

func BenchmarkErrors_ErrorBuilder(b *testing.B) {
	e := getErrors()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.ErrorBuilder()
	}
}

func getErrors() Errors {
	return Errors([]error{
		errors.New("1"),
		errors.New("2"),
		errors.New("3"),
		errors.New("4"),
		errors.New("5"),
	})
}
