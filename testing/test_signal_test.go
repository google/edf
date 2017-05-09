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

package testing

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestBiLevelBasics(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Duration(10) * time.Second)
	records := []float64{rand.NormFloat64(), rand.NormFloat64(),
		rand.NormFloat64(), rand.NormFloat64(),
		rand.NormFloat64(), rand.NormFloat64(),
		rand.NormFloat64(), rand.NormFloat64(),
		rand.NormFloat64(), rand.NormFloat64(),
		rand.NormFloat64(), rand.NormFloat64(),
		rand.NormFloat64(), rand.NormFloat64(),
		rand.NormFloat64(), rand.NormFloat64(),
		rand.NormFloat64(), rand.NormFloat64(),
		rand.NormFloat64(), rand.NormFloat64()}
	baseSignal := NewTestingSignal(t1, t2, records)
	if baseSignal.StartTime() != t1 {
		t.Errorf("%v should be equal to %v", baseSignal.StartTime(), t1)
	}
	if baseSignal.EndTime() != t2 {
		t.Errorf("%v should be equal to %v", baseSignal.EndTime(), t2)
	}
	if baseSignal.SamplingRate() != time.Duration(500)*time.Millisecond {
		t.Errorf("Wrong sampling rate: %v should be %v", baseSignal.SamplingRate(), time.Duration(500)*time.Millisecond)
	}
	actualRecords, err := baseSignal.Recording(
		t1.Add(time.Duration(500)*time.Millisecond*3),
		t1.Add(time.Duration(500)*time.Millisecond*15))
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(actualRecords, records[3:15]) {
		t.Errorf("%v should be equal to %v", actualRecords, records[3:15])
	}
}
