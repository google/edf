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

// Package processing contains a processing utilities for EDF signal.
package processing

import (
	"errors"
	"math"
	"time"

	"github.com/google/edf"
	"github.com/google/edf/signals"
)

type Level int

const (
	TRANSITION Level = iota
	LOW
	HIGH
)

// BiLevelSignal is an EDF signal coerced into two values (levels) only.
type BiLevelSignal interface {
	signals.Signal
	Low() float64
	High() float64
	BiLevelRecording(start, end time.Time) ([]Level, error)
}

// biLevelSignal is the internal representation of BiLevelSignal
type biLevelSignal struct {
	s         signals.Signal
	low       float64
	high      float64
	tolerance float64
}

func (s *biLevelSignal) Low() float64 {
	return s.low
}

func (s *biLevelSignal) High() float64 {
	return s.high
}

func (s *biLevelSignal) Label() string {
	return s.s.Label() + " (bilevel)"
}

func (s *biLevelSignal) StartTime() time.Time {
	return s.s.StartTime()
}

func (s *biLevelSignal) EndTime() time.Time {
	return s.s.EndTime()
}

func (s *biLevelSignal) Definition() *edf.SignalDefinition {
	return nil
}

func (s *biLevelSignal) SamplingRate() time.Duration {
	ds, ok := s.s.(signals.DataSignal)
	if !ok {
		return 0
	}
	return ds.SamplingRate()
}

func (s *biLevelSignal) BiLevelRecording(start, end time.Time) ([]Level, error) {
	ds, ok := s.s.(signals.DataSignal)
	if !ok {
		return nil, errors.New("BiLevelRecording can only be created for data signals")
	}
	r, err := ds.Recording(start, end)
	if err != nil {
		return nil, err
	}
	data := make([]Level, len(r))
	for i, dataPoint := range r {
		if math.Abs(dataPoint-s.low) < s.tolerance {
			data[i] = LOW
		} else if math.Abs(dataPoint-s.high) < s.tolerance {
			data[i] = HIGH
		} else {
			data[i] = TRANSITION
		}
	}
	return data, nil
}

func (s *biLevelSignal) Recording(start, end time.Time) ([]float64, error) {
	ds, ok := s.s.(signals.DataSignal)
	if !ok {
		return nil, errors.New("BiLevelRecording can only be created for data signals")
	}
	r, err := ds.Recording(start, end)
	if err != nil {
		return nil, err
	}
	for i, dataPoint := range r {
		if math.Abs(dataPoint-s.low) < s.tolerance {
			r[i] = s.low
		} else if math.Abs(dataPoint-s.high) < s.tolerance {
			r[i] = s.high
		} else {
			r[i] = dataPoint
		}
	}
	return r, nil
}

// NewBiLevelSignal transforms a signal into a bi-level signal.
func NewBiLevelSignal(s signals.Signal, lowlevel, highlevel, tolerance float64) BiLevelSignal {
	return &biLevelSignal{s, lowlevel, highlevel, tolerance}
}
