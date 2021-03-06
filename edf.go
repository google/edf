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

// Edf represents an EDF+ file.
type Edf struct {
	Header *Header

	// Records
	Records []Record
}

// Header represents an EDF+ header.
type Header struct {
	Version             string
	PatiendID           string
	RecordingID         string
	StartDate           string
	StartTime           string
	HeaderSize          uint32
	Reserved            string
	NumDataRecords      uint32
	DurationDataRecords float32
	NumSignals          uint32
	Signals             []SignalDefinition
}

// SignalDefinition holds the definition of an EDF signal.
type SignalDefinition struct {
	Label             string
	TransducerType    string
	PhysicalDimension string
	PhysicalMinimum   string
	PhysicalMaximum   string
	DigitalMinimum    string
	DigitalMaximum    string
	Prefiltering      string
	SamplesRecord     uint32
	Reserved          string
}

// Record holds a single record entry from the EDF file.
type Record struct {
	Signals []SignalRecord
}

// SignalRecord holds the samples for a single signal inside a data record.
type SignalRecord struct {
	Samples []int16
}
