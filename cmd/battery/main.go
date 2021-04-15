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

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gmelchett/battery"
)

func printBattery(idx int, bat *battery.Battery) {
	fmt.Printf(
		"BAT%d: %s, %.2f%%",
		idx,
		bat.State,
		bat.Current/bat.Full*100,
	)
	defer fmt.Printf(" [Voltage: %.2fV (design: %.2fV)]\n", bat.Voltage, bat.DesignVoltage)

	var str string
	var timeNum float64
	switch bat.State {
	case battery.Discharging:
		if bat.ChargeRate == 0 {
			fmt.Print(", discharging at zero rate - will never fully discharge")
			return
		}
		str = "remaining"
		timeNum = bat.Current / bat.ChargeRate
	case battery.Charging:
		if bat.ChargeRate == 0 {
			fmt.Print(", charging at zero rate - will never fully charge")
			return
		}
		str = "until charged"
		timeNum = (bat.Full - bat.Current) / bat.ChargeRate
	default:
		return
	}
	duration, _ := time.ParseDuration(fmt.Sprintf("%fh", timeNum))
	fmt.Printf(", %s %s", duration, str)
}

func main() {
	batteries, err := battery.GetAll()
	if err, isFatal := err.(battery.ErrFatal); isFatal {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(batteries) == 0 {
		fmt.Fprintln(os.Stderr, "No batteries")
		os.Exit(1)
	}
	errs, partialErrs := err.(battery.Errors)
	for i, bat := range batteries {
		if partialErrs && errs[i] != nil {
			fmt.Fprintf(os.Stderr, "Error getting info for BAT%d: %s\n", i, errs[i])
			continue
		}
		printBattery(i, bat)
	}
}
