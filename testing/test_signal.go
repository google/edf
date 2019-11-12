// Copyright 2017 Google Inc. All Rights Reserved.
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

// Package testing contains a testing utilities for the edf package and sub-packages.
package testing

import (
	"fmt"
	"time"

	"github.com/google/edf"
	"github.com/google/edf/signals"
)

type testingSignal struct {
	start   time.Time
	end     time.Time
	records []float64
}

func (ts *testingSignal) Label() string {
	return "Testing signal"
}

func (ts *testingSignal) StartTime() time.Time {
	return ts.start
}

func (ts *testingSignal) EndTime() time.Time { return ts.end }

func (ts *testingSignal) Definition() *edf.SignalDefinition { return nil }

func (ts *testingSignal) SamplingRate() time.Duration {
	return ts.end.Sub(ts.start) / time.Duration(len(ts.records))
}

func (ts *testingSignal) Recording(start, end time.Time) ([]float64, error) {
	if start.Before(ts.start) {
		return nil, fmt.Errorf("%v is before %v", start, ts.start)
	}
	if end.After(ts.end) {
		return nil, fmt.Errorf("%v is after %v", end, ts.end)
	}
	begin := start.Sub(ts.start) / ts.SamplingRate()
	finish := end.Sub(ts.start) / ts.SamplingRate()
	result := make([]float64, int64(finish-begin))
	copy(result, ts.records[begin:finish])
	return result, nil
}

func NewTestingSignal(start, end time.Time, records []float64) signals.DataSignal {
	return &testingSignal{start: start, end: end, records: records}
}
