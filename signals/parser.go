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

package signals

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/edf"
)

// GetSignals return the signals from an EDF dataset.
func GetSignals(e *edf.Edf) ([]Signal, error) {
	signals := make([]Signal, e.Header.NumSignals)
	for i := range e.Header.Signals {
		signal, err := newEdfSignal(e, i)
		if err != nil {
			return nil, err
		}
		if e.Header.Signals[i].Label == "EDF Annotations" {
			signals[i], err = newAnnotationSignal(signal)
			if err != nil {
				return nil, err
			}
		} else {
			signals[i] = newDataSignal(signal)
		}
	}
	return signals, nil
}

type edfSignal struct {
	edf         *edf.Edf
	startTime   time.Time
	endTime     time.Time
	signalIndex int

	// digital to physical conversion parameters
	a float64
	b float64
}

func newEdfSignal(e *edf.Edf, signalIndex int) (*edfSignal, error) {
	s := new(edfSignal)
	s.edf = e
	s.signalIndex = signalIndex
	start, err := getStartTime(e.Header)
	if err != nil {
		return nil, err
	}
	s.startTime = start
	end, err := getEndTime(e.Header)
	if err != nil {
		return nil, err
	}
	s.endTime = end

	def := &s.edf.Header.Signals[signalIndex]
	physMin, err := strconv.ParseFloat(def.PhysicalMinimum, 32)
	if err != nil {
		return nil, err
	}
	physMax, err := strconv.ParseFloat(def.PhysicalMaximum, 32)
	if err != nil {
		return nil, err
	}
	digiMin, err := strconv.ParseFloat(def.DigitalMinimum, 32)
	if err != nil {
		return nil, err
	}
	digiMax, err := strconv.ParseFloat(def.DigitalMaximum, 32)
	if err != nil {
		return nil, err
	}

	s.a = (physMax - physMin) / (digiMax - digiMin)
	s.b = physMin - s.a*digiMin

	return s, nil
}

func (s *edfSignal) Label() string {
	return s.Definition().Label
}

// Start date and time of the recording.
func (s *edfSignal) StartTime() time.Time {
	return s.startTime
}

// End date and time of the recording.
func (s *edfSignal) EndTime() time.Time {
	return s.endTime
}

// Signal definition. This may be nil for composite/generated signals.
func (s *edfSignal) Definition() *edf.SignalDefinition {
	return &s.edf.Header.Signals[s.signalIndex]
}

// getStartTime returns the starting date and time of the recording
func getStartTime(h *edf.Header) (time.Time, error) {
	return time.Parse("02.01.06 15.04.05", h.StartDate+" "+h.StartTime)
}

// getEndTime returns the end date and time of the recording
func getEndTime(h *edf.Header) (time.Time, error) {
	start, err := getStartTime(h)
	if err != nil {
		return start, err
	}
	end := start.Add(
		time.Duration(
			float32(h.NumDataRecords)*h.DurationDataRecords) * time.Second)
	return end, nil
}

// getSignalData returns the signal samples between the specified times.
func getSignalData(e *edf.Edf, signalIndex int, start, end time.Time) ([]int16, error) {
	recordingStart, err := getStartTime(e.Header)
	if err != nil {
		return nil, err
	}
	if recordingStart.After(start) {
		return nil, fmt.Errorf("Requesting data before the recording")
	}

	recordingEnd, err := getEndTime(e.Header)
	if err != nil {
		return nil, err
	}
	if recordingEnd.Before(end) {
		return nil, fmt.Errorf("Requesting data after the recording")
	}

	durationSample := float64(e.Header.DurationDataRecords) / float64(e.Header.Signals[signalIndex].SamplesRecord)

	skipStart := start.Sub(recordingStart)
	startRecord := uint32(
		skipStart.Seconds() / float64(e.Header.DurationDataRecords))
	startSample := uint32(
		(skipStart.Seconds() - (float64(startRecord) * float64(e.Header.DurationDataRecords))) / (durationSample))

	endRecord := uint32(
		end.Sub(recordingStart).Seconds() / float64(e.Header.DurationDataRecords))
	endSample := uint32((end.Sub(recordingStart).Seconds() - (float64(endRecord) * float64(e.Header.DurationDataRecords))) / (durationSample))

	numSamples := e.Header.Signals[signalIndex].SamplesRecord*(endRecord-startRecord-1) + endSample + (e.Header.Signals[signalIndex].SamplesRecord - startSample)

	result := make([]int16, numSamples)
	s := 0
	for i := startRecord; i <= endRecord; i++ {
		for j := uint32(0); j < e.Header.Signals[signalIndex].SamplesRecord; j++ {
			if i == startRecord && j == 0 {
				j = startSample
			}
			if i == endRecord && j >= endSample {
				break
			}
			result[s] = e.Records[i].Signals[signalIndex].Samples[j]
			s++
		}
	}
	return result, nil
}
