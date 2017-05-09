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

package processing

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	edf_testing "github.com/google/edf/testing"
)

func randFloat32(mean float64) float64 {
	return rand.Float64() + mean - 0.5
}

func TestBiLevelBasics(t *testing.T) {
	t1 := time.Now()
	t2 := t1.Add(time.Duration(10) * time.Second)
	records := []float64{randFloat32(10), randFloat32(10),
		randFloat32(5), randFloat32(5),
		randFloat32(10), randFloat32(10),
		randFloat32(7), randFloat32(5),
		randFloat32(10), randFloat32(10),
		randFloat32(5), randFloat32(5),
		randFloat32(10), randFloat32(10),
		randFloat32(5), randFloat32(5),
		randFloat32(10), randFloat32(10),
		randFloat32(5), randFloat32(5)}
	baseSignal := edf_testing.NewTestingSignal(t1, t2, records)
	biLevel := NewBiLevelSignal(baseSignal, 5, 10, 1)
	if biLevel.StartTime() != t1 {
		t.Errorf("%v should be equal to %v", biLevel.StartTime(), t1)
	}
	if biLevel.EndTime() != t2 {
		t.Errorf("%v should be equal to %v", biLevel.EndTime(), t2)
	}
	expectedSignal := []Level{HIGH, HIGH,
		LOW, LOW,
		HIGH, HIGH,
		TRANSITION, LOW,
		HIGH, HIGH,
		LOW, LOW,
		HIGH, HIGH,
		LOW, LOW,
		HIGH, HIGH,
		LOW, LOW}
	actualSignal, err := biLevel.BiLevelRecording(t1, t2)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(actualSignal, expectedSignal) {
		t.Errorf("%v should be equal to %v", actualSignal, expectedSignal)
	}
}
