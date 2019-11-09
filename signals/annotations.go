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
	timestamp   timestamp
	annotations []string
}

func (tsa *timeStampedAnnotation) Time() time.Time {
	return tsa.timestamp.base.Add(time.Duration(tsa.timestamp.Onset * float64(time.Second)))
}

func (tsa *timeStampedAnnotation) End() time.Time {
	return tsa.Time().Add(time.Duration(tsa.timestamp.Duration * float64(time.Second)))
}

func (tsa *timeStampedAnnotation) Annotations() []string {
	return tsa.annotations
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
	for annotationIndex := range as.annotations {
		annotation := &(as.annotations[annotationIndex])
		at := annotation.End()
		if (at == start || at.After(start)) && (at == end || at.Before(end)) {
			result = append(result, annotation)
		}
	}
	return result, nil
}

func newAnnotationSignal(baseSignal *edfSignal) (AnnotationSignal, error) {
	records := baseSignal.edf.Records
	aS := annotationSignal{baseSignal, []timeStampedAnnotation{}}
	for _, record := range records {
		tsa := timeStampedAnnotation{timestamp{baseSignal.StartTime(), 0, 0}, []string{}}
		signal := record.Signals[baseSignal.signalIndex]
		// Extract bytes from 16-bit integers.
		buffer := new(bytes.Buffer)
		for _, biChar := range signal.Samples {
			buffer.WriteByte(byte(biChar & 0xFF))
			buffer.WriteByte(byte(biChar >> 8))
		}
		// Zero bytes don't count.
		rawAnnotations := bytes.Split(bytes.Replace(buffer.Bytes(), []byte{'\x00'}, []byte{}, -1), []byte{'\x14'})
		realAnnotations := [][]byte{}
		for _, annotation := range rawAnnotations {
			if len(annotation) != 0 {
				realAnnotations = append(realAnnotations, annotation)
			}
		}
		for annotationIndex, realAnnotation := range realAnnotations {
			if len(realAnnotation) == 0 {
				continue
			}
			// I have found in the wild timestamps both as the first or the second annotation.
			if annotationIndex == 0 || annotationIndex == 1 {
				timestamp, err := parseTimestamp(baseSignal.StartTime(), realAnnotation)
				if err == nil {
					tsa.timestamp = timestamp
					continue
				}
			}
			tsa.annotations = append(tsa.annotations, string(realAnnotation))
		}
		aS.annotations = append(aS.annotations, tsa)
	}
	return &aS, nil
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
