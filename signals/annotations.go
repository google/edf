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
	"bytes"
	"errors"
	"io"
	"strconv"
	"time"
)

type timestamp struct {
	base     time.Time
	Onset    float64
	Duration float64
}

type timeStampedAnnotation struct {
	timestamp  timestamp
	annotation string
}

func (tsa *timeStampedAnnotation) Time() time.Time {
	return tsa.timestamp.base.Add(time.Duration(tsa.timestamp.Onset * float64(time.Second)))
}

func (tsa *timeStampedAnnotation) End() time.Time {
	return tsa.Time().Add(time.Duration(tsa.timestamp.Duration * float64(time.Second)))
}

func (tsa *timeStampedAnnotation) Annotation() string {
	return tsa.annotation
}

type annotationSignal struct {
	Signal
	annotations []timeStampedAnnotation
}

func (as *annotationSignal) Annotations(start, end time.Time) ([]Annotation, error) {
	result := make([]Annotation, 0)
	if start.Before(as.StartTime()) || end.After(as.EndTime()) || start.After(end) {
		return nil, errors.New("Invalid start or end time")
	}
	for _, annotation := range as.annotations {
		at := annotation.End()
		if (at == start || at.After(start)) && (at == end || at.Before(end)) {
			result = append(result, &annotation)
		}
	}
	return result, nil
}

func newAnnotationSignal(baseSignal *edfSignal) (AnnotationSignal, error) {
	records := baseSignal.edf.Records
	buffer := new(bytes.Buffer)
	aS := &annotationSignal{baseSignal, make([]timeStampedAnnotation, 0)}
	durationSample := time.Duration(
		float64(baseSignal.edf.Header.DurationDataRecords)/float64(baseSignal.Definition().SamplesRecord)) * time.Second
	for i, record := range records {
		tsa := timeStampedAnnotation{}

		signal := record.Signals[baseSignal.signalIndex]
		// Read signal into string
		for _, biChar := range signal.Samples {
			buffer.WriteByte(byte(biChar & 0xFF))
			buffer.WriteByte(byte(biChar >> 8))
		}

		for {
			sampleBytes, err := buffer.ReadBytes('\x00')
			if err == io.EOF {
				break
			} else if err != nil {
				return nil, err
			}
			sampleBuffer := bytes.NewBuffer(sampleBytes[0 : len(sampleBytes)-1])
			// Parse string
			tsBytes, err := sampleBuffer.ReadBytes('\x14')
			if err == io.EOF {
				// No timestamp
				tsa.annotation = string(tsBytes)
				tsa.timestamp.base = baseSignal.StartTime().Add(time.Duration(i) * durationSample)
				aS.annotations = append(aS.annotations, tsa)
			} else if err != nil {
				return nil, err
			} else {
				tsa.annotation = string(sampleBuffer.Bytes())
				tsa.timestamp, err = parseTimestamp(baseSignal.StartTime(), tsBytes[0:len(tsBytes)-1])
				if err != nil {
					return nil, err
				}
			}
			if len(tsa.annotation) == 0 {
				continue
			}
			aS.annotations = append(aS.annotations, tsa)
		}
	}
	return aS, nil
}

func parseTimestamp(base time.Time, tsBytes []byte) (timestamp, error) {
	buffer := bytes.NewBuffer(tsBytes)
	onsetBytes, err := buffer.ReadString('\x15')
	if err == io.EOF {
		// No duration
		onset, err := strconv.ParseFloat(onsetBytes, 64)
		if err != nil {
			return timestamp{}, err
		}
		return timestamp{base, onset, 0}, nil
	} else if err != nil {
		return timestamp{}, err
	}
	onset, err := strconv.ParseFloat(onsetBytes[0:len(onsetBytes)-1], 64)
	if err != nil {
		return timestamp{}, err
	}
	durationBytes := buffer.Bytes()
	duration, err := strconv.ParseFloat(string(durationBytes), 64)
	if err != nil {
		return timestamp{}, err
	}
	return timestamp{base, onset, duration}, nil
}
