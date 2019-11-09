// Copyright 2019 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/// This package is an example tool showing how to use the edf library.
package main

import (
	"flag"
	"fmt"

	"github.com/google/edf"
	"github.com/google/edf/signals"
)

var (
	input           = flag.String("input", "", "input")
	signalLabel     = flag.String("signal", "", "signal")
	annotationLabel = flag.Bool("annotations", false, "annotations")
)

func main() {
	flag.Parse()
	edfFile, err := edf.ReadEDF(*input)
	if err != nil {
		panic(err)
	}

	if *signalLabel == "" && !*annotationLabel {
		for _, signal := range edfFile.Header.Signals {
			fmt.Printf("Signal: '%s'\n", signal.Label)
		}
		return
	}

	edfSignals, err := signals.GetSignals(edfFile)
	if err != nil {
		panic(err)
	}
	for _, signal := range edfSignals {
		if signal.Label() != *signalLabel && (!*annotationLabel || signal.Label() != "EDF Annotations") {
			continue
		}
		fmt.Printf("Signal: %s\n", signal.Label())
		if !(*annotationLabel) {
			values, err := signal.(signals.DataSignal).Recording(signal.StartTime(), signal.EndTime())
			if err != nil {
				panic(err)
			}
			fmt.Println(values)
		} else {
			values, err := signal.(signals.AnnotationSignal).Annotations(signal.StartTime(), signal.EndTime())
			if err != nil {
				panic(err)
			}
			for _, value := range values {
				fmt.Println(value.Time(), value.Annotations())
			}
		}
	}
}
