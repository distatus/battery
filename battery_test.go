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
	"fmt"
	"reflect"
	"testing"
)

func TestNewState(t *testing.T) {
	cases := []struct {
		in       string
		stateOut State
		errorOut error
	}{
		{"Charging", Charging, nil},
		{"charging", Charging, nil},
		{"Not Charging", NotCharging, nil},
		{"strange", Unknown, fmt.Errorf("Invalid state `strange`")},
	}

	for i, c := range cases {
		state, err := newState(c.in)

		if state != c.stateOut {
			t.Errorf("%d: %v != %v", i, state, c.stateOut)
		}
		if !reflect.DeepEqual(err, c.errorOut) {
			t.Errorf("%d: %v != %v", i, err, c.errorOut)
		}
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		batteryIn  *Battery
		errorIn    error
		batteryOut *Battery
		errorOut   error
	}{{
		&Battery{Full: 1}, nil,
		&Battery{Full: 1}, nil,
	}, {
		nil, fmt.Errorf("t1"),
		nil, ErrFatal{fmt.Errorf("t1")},
	}, {
		&Battery{Full: 2}, ErrPartial{Current: fmt.Errorf("t2")},
		&Battery{Full: 2}, ErrPartial{Current: fmt.Errorf("t2")},
	}, {
		&Battery{Full: 3}, ErrPartial{},
		&Battery{Full: 3}, nil,
	}, {
		nil, ErrPartial{State: fmt.Errorf("t3"), Current: fmt.Errorf("t4"), Full: fmt.Errorf("t5"), Design: fmt.Errorf("t6"), ChargeRate: fmt.Errorf("t7"), Voltage: fmt.Errorf("t8"), DesignVoltage: fmt.Errorf("t9")},
		nil, ErrFatal{ErrAllNotNil},
	}}

	for i, c := range cases {
		f := func(idx int) (*Battery, error) {
			return c.batteryIn, c.errorIn
		}

		battery, err := get(f, 0)

		if !reflect.DeepEqual(battery, c.batteryOut) {
			t.Errorf("%d: %v != %v", i, battery, c.batteryOut)
		}
		if !reflect.DeepEqual(err, c.errorOut) {
			t.Errorf("%d: %v != %v", i, err, c.errorOut)
		}
	}
}

func TestGetAll(t *testing.T) {
	cases := []struct {
		batteriesIn  []*Battery
		errorsIn     error
		batteriesOut []*Battery
		errorsOut    error
	}{{
		[]*Battery{{Full: 1}}, nil,
		[]*Battery{{Full: 1}}, nil,
	}, {
		[]*Battery{}, fmt.Errorf("t1"),
		[]*Battery{}, ErrFatal{fmt.Errorf("t1")},
	}, {
		[]*Battery{{Full: 2}, {Full: 3}}, Errors{ErrPartial{}, ErrPartial{}},
		[]*Battery{{Full: 2}, {Full: 3}}, nil,
	}, {
		[]*Battery{{Full: 4}, {Full: 5}}, Errors{ErrPartial{State: fmt.Errorf("t2"), Current: fmt.Errorf("t3"), Full: fmt.Errorf("t4"), Design: fmt.Errorf("t5"), ChargeRate: fmt.Errorf("t6"), Voltage: fmt.Errorf("t101"), DesignVoltage: fmt.Errorf("t102")}, ErrPartial{State: fmt.Errorf("t7"), Current: fmt.Errorf("t8"), Full: fmt.Errorf("t9"), Design: fmt.Errorf("t10"), ChargeRate: fmt.Errorf("t11"), Voltage: fmt.Errorf("t103"), DesignVoltage: fmt.Errorf("t104")}},
		nil, ErrFatal{ErrAllNotNil},
	}, {
		[]*Battery{{Full: 6}, {Full: 7}}, Errors{ErrPartial{State: fmt.Errorf("t12")}, fmt.Errorf("t13")},
		[]*Battery{{Full: 6}, {Full: 7}}, Errors{ErrPartial{State: fmt.Errorf("t12")}, ErrFatal{fmt.Errorf("t13")}},
	}, {
		[]*Battery{{}, {Full: 8}}, Errors{ErrPartial{State: fmt.Errorf("t14"), Current: fmt.Errorf("t15"), Full: fmt.Errorf("t16"), Design: fmt.Errorf("t17"), ChargeRate: fmt.Errorf("t18"), Voltage: fmt.Errorf("t105"), DesignVoltage: fmt.Errorf("t106")}, nil},
		[]*Battery{{}, {Full: 8}}, Errors{ErrFatal{ErrAllNotNil}, nil},
	}, {
		[]*Battery{{Full: 9}, {Full: 10}}, Errors{ErrPartial{}, fmt.Errorf("t19")},
		[]*Battery{{Full: 9}, {Full: 10}}, Errors{nil, ErrFatal{fmt.Errorf("t19")}},
	}}

	for i, c := range cases {
		f := func() ([]*Battery, error) {
			return c.batteriesIn, c.errorsIn
		}
		batteries, err := getAll(f)

		if !reflect.DeepEqual(batteries, c.batteriesOut) {
			t.Errorf("%d: %v != %v", i, batteries, c.batteriesOut)
		}
		if !reflect.DeepEqual(err, c.errorsOut) {
			t.Errorf("%d: %v != %v", i, err, c.errorsOut)
		}
	}
}

func ExampleGetAll() {
	batteries, err := GetAll()
	if err != nil {
		fmt.Println("Could not get batteries info")
		return
	}

	for i, battery := range batteries {
		fmt.Printf("Bat%d: ", i)
		fmt.Printf("state: %s, ", battery.State)
		fmt.Printf("current capacity: %f mWh, ", battery.Current)
		fmt.Printf("last full capacity: %f mWh, ", battery.Full)
		fmt.Printf("design capacity: %f mWh, ", battery.Design)
		fmt.Printf("charge rate: %f mW\n", battery.ChargeRate)
	}
}

func ExampleGetAll_errors() {
	_, err := GetAll()
	if err == nil {
		fmt.Println("Got batteries info")
		return
	}

	switch perr := err.(type) {
	case ErrFatal:
		fmt.Println("Fatal error! No info retrieved")
	case Errors:
		for i, err := range perr {
			if err != nil {
				fmt.Printf("Could not get battery info for `%d`\n", i)
				continue
			}

			fmt.Printf("Got battery info for `%d`\n", i)
		}
	}
}

func ExampleGet() {
	battery, err := Get(0)
	if err != nil {
		fmt.Println("Could not get battery info")
		return
	}

	fmt.Printf("Bat%d: ", 0)
	fmt.Printf("state: %s, ", battery.State)
	fmt.Printf("current capacity: %f mWh, ", battery.Current)
	fmt.Printf("last full capacity: %f mWh, ", battery.Full)
	fmt.Printf("design capacity: %f mWh, ", battery.Design)
	fmt.Printf("charge rate: %f mW\n", battery.ChargeRate)
}

func ExampleGet_errors() {
	_, err := Get(0)
	if err == nil {
		fmt.Println("Got battery info")
		return
	}

	switch perr := err.(type) {
	case ErrFatal:
		fmt.Println("Fatal error! No info retrieved")
	case ErrPartial:
		if perr.Current != nil {
			fmt.Println("Could not get current battery capacity")
			return
		}

		fmt.Println("Got current battery capacity")
	}
}
