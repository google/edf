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

import "time"

type dataSignal struct {
	Signal

	e *edfSignal
}

func newDataSignal(edfSignal *edfSignal) *dataSignal {
	return &dataSignal{
		Signal: edfSignal,
		e:      edfSignal,
	}
}

// Returns the time between two recording samples of this signal.
func (s *dataSignal) SamplingRate() time.Duration {
	return time.Duration(s.e.edf.Header.DurationDataRecords/float32(s.Definition().SamplesRecord)) * time.Second
}

// Returns the recording data, in physical units.
func (s *dataSignal) Recording(start, end time.Time) ([]float64, error) {
	r, err := getSignalData(s.e.edf, s.e.signalIndex, start, end)
	if err != nil {
		return nil, err
	}
	result := make([]float64, len(r))
	for i, dataPoint := range r {
		result[i] = s.e.a*float64(dataPoint) + s.e.b
	}
	return result, nil
}
