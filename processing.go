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

// Package edf contains a parser for EDF+ files.
package edf

import (
	"math"
	"time"
)

// BiLevelSignal is an EDF signal coerced into two values (levels) only.
type BiLevelSignal interface {
	Signal
	Low() float32
	High() float32
}

// biLevelSignal is the internal representation of BiLevelSignal
type biLevelSignal struct {
	s    Signal
	low  float32
	high float32
}

func (s *biLevelSignal) Low() float32 {
	return s.low
}

func (s *biLevelSignal) High() float32 {
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

func (s *biLevelSignal) Definition() *SignalDefinition {
	return nil
}

func (s *biLevelSignal) SamplingRate() time.Duration {
	return s.s.SamplingRate()
}

func (s *biLevelSignal) Recording(start, end time.Time) ([]float32, error) {
	r, err := s.s.Recording(start, end)
	if err != nil {
		return nil, err
	}
	for i, dataPoint := range r {
		if math.Abs(float64(dataPoint-s.low)) < math.Abs(float64(dataPoint-s.high)) {
			r[i] = s.low
		} else {
			r[i] = s.high
		}
	}
	return r, nil
}

// NewBiLevelSignal transforms a signal into a bi-level signal.
func NewBiLevelSignal(s Signal, lowlevel, highlevel float32) BiLevelSignal {
	return &biLevelSignal{s, lowlevel, highlevel}
}
