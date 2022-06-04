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

	sm "github.com/lni/dragonboat/v4/statemachine"
)

// SecondStateMachine is the IStateMachine implementation used in the
// multigroup example for handling all inputs ends with "?".
// See https://github.com/lni/dragonboat/blob/master/statemachine/rsm.go for
// more details of the IStateMachine interface.
// The behaviour of SecondStateMachine is similar to the ExampleStateMachine.
// The biggest difference is that its Update() method has different print
// out messages. See Update() for details.
type SecondStateMachine struct {
	ShardID   uint64
	ReplicaID uint64
	Count     uint64
}

// NewSecondStateMachine creates and return a new SecondStateMachine object.
func NewSecondStateMachine(shardID uint64, replicaID uint64) sm.IStateMachine {
	return &SecondStateMachine{
		ShardID:   shardID,
		ReplicaID: replicaID,
		Count:     0,
	}
}

// Lookup performs local lookup on the SecondStateMachine instance. In this example,
// we always return the Count value as a little endian binary encoded byte
// slice.
func (s *SecondStateMachine) Lookup(query interface{}) (interface{}, error) {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, s.Count)
	return result, nil
}

// Update updates the object using the specified committed raft entry.
func (s *SecondStateMachine) Update(e sm.Entry) (sm.Result, error) {
	// in this example, we regard the input as a question.
	s.Count++
	fmt.Printf("got a question from user: %s, count:%d\n", string(e.Cmd), s.Count)
	return sm.Result{Value: uint64(len(e.Cmd))}, nil
}

// SaveSnapshot saves the current IStateMachine state into a snapshot using the
// specified io.Writer object.
func (s *SecondStateMachine) SaveSnapshot(w io.Writer,
	fc sm.ISnapshotFileCollection, done <-chan struct{}) error {
	// as shown above, the only state that can be saved is the Count variable
	// there is no external file in this IStateMachine example, we thus leave
	// the fc untouched
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, s.Count)
	_, err := w.Write(data)
	return err
}

// RecoverFromSnapshot recovers the state using the provided snapshot.
func (s *SecondStateMachine) RecoverFromSnapshot(r io.Reader,
	files []sm.SnapshotFile, done <-chan struct{}) error {
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
func (s *SecondStateMachine) Close() error { return nil }
