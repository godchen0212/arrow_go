// Copyright (C) 2019-2020 Zilliz. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions and limitations under the License.

package storage

import (
	"bytes"
)

// EventReader is used to parse the events contained in the Binlog file.
type EventReader struct {
	PayloadReaderInterface
	buffer   *bytes.Buffer
	isClosed bool
}



// Close closes EventReader object. It mainly calls the Close method of inner PayloadReaderInterface and
// mark itself as closed.
func (reader *EventReader) Close() error {
	if !reader.isClosed {
		reader.isClosed = true
		return reader.PayloadReaderInterface.Close()
	}
	return nil
}

func newEventReader(datatype DataType, buffer *bytes.Buffer) (*EventReader, error) {
	reader := &EventReader{
		buffer:   buffer,
		isClosed: false,
	}


	next := 0
	payloadBuffer := buffer.Next(next)
	payloadReader, err := NewPayloadReader(datatype, payloadBuffer)
	if err != nil {
		return nil, err
	}
	reader.PayloadReaderInterface = payloadReader
	return reader, nil
}
