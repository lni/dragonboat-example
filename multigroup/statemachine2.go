// Copyright 2017,2018 Lei Ni (nilei81@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/lni/dragonboat/statemachine"
)

// SecondStateMachine is the IStateMachine implementation used in the
// multigroup example for handling all inputs ends with "?".
// See https://github.com/lni/dragonboat/blob/master/statemachine/rsm.go for
// more details of the IStateMachine interface.
// The behaviour of SecondStateMachine is similar to the ExampleStateMachine.
// The biggest difference is that its Update() method has different print
// out messages. See Update() for details.
type SecondStateMachine struct {
	ClusterID uint64
	NodeID    uint64
	Count     uint64
}

// NewSecondStateMachine creates and return a new SecondStateMachine object.
func NewSecondStateMachine(clusterID uint64,
	nodeID uint64) statemachine.IStateMachine {
	return &SecondStateMachine{
		ClusterID: clusterID,
		NodeID:    nodeID,
		Count:     0,
	}
}

// Lookup performs local lookup on the SecondStateMachine instance. In this example,
// we always return the Count value as a little endian binary encoded byte
// slice.
func (s *SecondStateMachine) Lookup(query []byte) []byte {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, s.Count)
	return result
}

// Update updates the object using the specified committed raft entry.
func (s *SecondStateMachine) Update(data []byte) uint64 {
	// in this example, we regard the input as a question.
	s.Count++
	fmt.Printf("got a question from user: %s, count:%d\n", string(data), s.Count)
	return uint64(len(data))
}

// SaveSnapshot saves the current IStateMachine state into a snapshot using the
// specified io.Writer object.
func (s *SecondStateMachine) SaveSnapshot(w io.Writer,
	fc statemachine.ISnapshotFileCollection,
	done <-chan struct{}) (uint64, error) {
	// as shown above, the only state that can be saved is the Count variable
	// there is no external file in this IStateMachine example, we thus leave
	// the fc untouched
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, s.Count)
	_, err := w.Write(data)
	if err != nil {
		return 0, err
	}
	return uint64(len(data)), nil
}

// RecoverFromSnapshot recovers the state using the provided snapshot.
func (s *SecondStateMachine) RecoverFromSnapshot(r io.Reader,
	files []statemachine.SnapshotFile,
	done <-chan struct{}) error {
	// restore the Count variable, that is the only state we maintain in this
	// example, the input files is expected to be empty
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	v := binary.LittleEndian.Uint64(data)
	s.Count = v
	return nil
}

// Close closes the IStateMachine instance. There is nothing for us to cleanup
// or release as this is a pure in memory data store. Note that the Close
// method is not guaranteed to be called as node can crash at any time.
func (s *SecondStateMachine) Close() {}

// GetHash returns a uint64 representing the current object state.
func (s *SecondStateMachine) GetHash() uint64 {
	// the only state we have is that Count variable. that uint64 value pretty much
	// represents the state of this IStateMachine
	return s.Count
}
