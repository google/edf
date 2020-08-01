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
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// ReadEDF reads an EDF file.
func ReadEDF(filename string) (*Edf, error) {
	var err error

	fileInput, err := os.Open(filename)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}

	defer fileInput.Close()
	input := bufio.NewReader(fileInput)
	edf := &Edf{}
	edf.Header, err = readHeader(input)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	if err := readRecords(input, edf); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}

	return edf, nil
}

func readNextBytes(input io.Reader, size uint) ([]byte, error) {
	data := make([]byte, size)
	_, err := io.ReadFull(input, data)
	return data, err
}

// Reads the header of the EDF+ file.
func readHeader(input io.Reader) (*Header, error) {
	header := Header{}
	var data []byte
	var iData uint64
	var fData float64
	var sData string
	var err error

	if data, err = readNextBytes(input, 8); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	header.Version = strings.TrimSpace(string(data))

	if data, err = readNextBytes(input, 80); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	header.PatiendID = strings.TrimSpace(string(data))

	if data, err = readNextBytes(input, 80); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	header.RecordingID = strings.TrimSpace(string(data))

	if data, err = readNextBytes(input, 8); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	header.StartDate = strings.TrimSpace(string(data))

	if data, err = readNextBytes(input, 8); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	header.StartTime = strings.TrimSpace(string(data))

	if data, err = readNextBytes(input, 8); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	iData, err = strconv.ParseUint(strings.TrimSpace(string(data)), 10, 32)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	header.HeaderSize = uint32(iData)

	if data, err = readNextBytes(input, 44); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	header.Reserved = strings.TrimSpace(string(data))

	if data, err = readNextBytes(input, 8); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	iData, err = strconv.ParseUint(strings.TrimSpace(string(data)), 10, 32)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	header.NumDataRecords = uint32(iData)

	if data, err = readNextBytes(input, 8); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	fData, err = strconv.ParseFloat(strings.TrimSpace(string(data)), 32)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	header.DurationDataRecords = float32(fData)

	if data, err = readNextBytes(input, 4); err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	iData, err = strconv.ParseUint(strings.TrimSpace(string(data)), 10, 32)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	header.NumSignals = uint32(iData)

	header.Signals = make([]SignalDefinition, header.NumSignals)

	for signalIndex := uint32(0); signalIndex < header.NumSignals; signalIndex++ {
		signal := &header.Signals[signalIndex]
		if data, err = readNextBytes(input, 16); err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}
		signal.Label = strings.TrimSpace(string(data))
	}

	for signalIndex := uint32(0); signalIndex < header.NumSignals; signalIndex++ {
		signal := &header.Signals[signalIndex]
		if data, err = readNextBytes(input, 80); err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}
		signal.TransducerType = strings.TrimSpace(string(data))
	}

	for signalIndex := uint32(0); signalIndex < header.NumSignals; signalIndex++ {
		signal := &header.Signals[signalIndex]
		if data, err = readNextBytes(input, 8); err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}
		signal.PhysicalDimension = strings.TrimSpace(string(data))
	}

	for signalIndex := uint32(0); signalIndex < header.NumSignals; signalIndex++ {
		signal := &header.Signals[signalIndex]
		if data, err = readNextBytes(input, 8); err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}
		signal.PhysicalMinimum = strings.TrimSpace(string(data))
	}

	for signalIndex := uint32(0); signalIndex < header.NumSignals; signalIndex++ {
		signal := &header.Signals[signalIndex]
		if data, err = readNextBytes(input, 8); err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}
		signal.PhysicalMaximum = strings.TrimSpace(string(data))
	}

	for signalIndex := uint32(0); signalIndex < header.NumSignals; signalIndex++ {
		signal := &header.Signals[signalIndex]
		if data, err = readNextBytes(input, 8); err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}
		signal.DigitalMinimum = strings.TrimSpace(string(data))
	}

	for signalIndex := uint32(0); signalIndex < header.NumSignals; signalIndex++ {
		signal := &header.Signals[signalIndex]
		if data, err = readNextBytes(input, 8); err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}
		signal.DigitalMaximum = strings.TrimSpace(string(data))
	}

	for signalIndex := uint32(0); signalIndex < header.NumSignals; signalIndex++ {
		signal := &header.Signals[signalIndex]
		if data, err = readNextBytes(input, 80); err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}
		signal.Prefiltering = strings.TrimSpace(string(data))
	}

	for signalIndex := uint32(0); signalIndex < header.NumSignals; signalIndex++ {
		signal := &header.Signals[signalIndex]
		if data, err = readNextBytes(input, 8); err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}
		sData = strings.TrimSpace(string(data))
		if iData, err = strconv.ParseUint(sData, 10, 32); err != nil {
			log.Printf("Error at signal %d: %v\n", signalIndex, err)
			return nil, err
		}
		signal.SamplesRecord = uint32(iData)
	}

	for signalIndex := uint32(0); signalIndex < header.NumSignals; signalIndex++ {
		signal := &header.Signals[signalIndex]
		if data, err = readNextBytes(input, 32); err != nil {
			log.Printf("Error: %v\n", err)
			return nil, err
		}
		signal.Reserved = strings.TrimSpace(string(data))
	}

	return &header, nil
}

// Reads the data records from the EDF+ file. The header of the edf must be
// parsed and filled.
func readRecords(input io.Reader, edf *Edf) error {
	edf.Records = make([]Record, edf.Header.NumDataRecords)
	for i := uint32(0); i < edf.Header.NumDataRecords; i++ {
		record := &edf.Records[i]
		record.Signals = make([]SignalRecord, edf.Header.NumSignals)
		for s := uint32(0); s < edf.Header.NumSignals; s++ {
			signal := &record.Signals[s]
			signal.Samples = make([]int16, edf.Header.Signals[s].SamplesRecord)
			for d := uint32(0); d < edf.Header.Signals[s].SamplesRecord; d++ {
				err := binary.Read(input, binary.LittleEndian, &signal.Samples[d])
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
