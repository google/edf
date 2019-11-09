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

// Package signals interprets a raw EDF file into data or annotation signals.
package signals

import (
	"time"

	"github.com/google/edf"
)

// Signal wraps all the data recorded on for a signal.
type Signal interface {
	// Label of the signal.
	Label() string

	// StartTime returns the date and time of the recording.
	StartTime() time.Time

	// EndTime returns the end date and time of the recording.
	EndTime() time.Time

	// Definition returns the sgnal definition. This may be nil for
	// composite/generated signals.
	Definition() *edf.SignalDefinition
}

// DataSignal is a signal representing measures of a physical quantity sampled at regular intervals.
type DataSignal interface {
	Signal

	// SamplingRate returns the time between two recording samples of this signal.
	SamplingRate() time.Duration

	// Recording returns the recording data, in physical units.
	Recording(start, end time.Time) ([]float64, error)
}

// Annotation is a single annotation from an annotation signal.
type Annotation interface {
	// Time of the annotation.
	Time() time.Time

	// End time of the annotation.
	End() time.Time

	// Annotation contents.
	Annotations() []string
}

// AnnotationSignal is a signal containing text annotations (timestamped or at regular intervals) per the
// EDF+ specification.
type AnnotationSignal interface {
	Signal

	// Annotations returns the annotations.
	Annotations(start, end time.Time) ([]Annotation, error)
}
