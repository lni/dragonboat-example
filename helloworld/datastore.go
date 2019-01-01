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

	"github.com/lni/dragonboat/datastore"
)

// ExampleStore is the IDataStore implementation used in the helloworld example
type ExampleStore struct {
	ClusterID uint64
	NodeID    uint64
	Count     uint64
}

// NewExampleStore creates and return a new ExampleStore object.
func NewExampleStore(clusterID uint64, nodeID uint64) datastore.IDataStore {
	return &ExampleStore{
		ClusterID: clusterID,
		NodeID:    nodeID,
		Count:     0,
	}
}

// Lookup performs local lookup on the ExampleStore instance. In this example,
// we always return the Count value as a little endian binary encoded byte
// slice.
func (s *ExampleStore) Lookup(query []byte) []byte {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, s.Count)
	return result
}

// Update updates the object using the specified committed raft entry.
func (s *ExampleStore) Update(data []byte) uint64 {
	// in this example, we print out the following hello world message for each
	// incoming update request. we also increase the counter by one to remember
	// how many updates we have applied
	s.Count++
	fmt.Printf("from ExampleStore.Update(), msg: %s, count:%d\n",
		string(data), s.Count)
	return uint64(len(data))
}

// SaveSnapshot saves the current object state into a snapshot using the
// specified io.Writer object.
func (s *ExampleStore) SaveSnapshot(w io.Writer,
	fc datastore.ISnapshotFileCollection,
	done <-chan struct{}) (uint64, error) {
	// as shown above, the only state that can be saved is the Count variable
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, s.Count)
	_, err := w.Write(data)
	if err != nil {
		return 0, err
	}
	return uint64(len(data)), nil
}

// RecoverFromSnapshot recovers the state using the provided snapshot.
func (s *ExampleStore) RecoverFromSnapshot(r io.Reader,
	files []datastore.SnapshotFile,
	done <-chan struct{}) error {
	// restore the Count variable
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	v := binary.LittleEndian.Uint64(data)
	s.Count = v
	return nil
}

// Close closes the IDataStore instance. There is nothing for us to cleanup or
// release as this is a pure in memory data store.
func (s *ExampleStore) Close() {}

// GetHash returns a uint64 representing the current object state.
func (s *ExampleStore) GetHash() uint64 {
	// again - the only state we have is that Count variable. that
	// uint64 pretty much represents the state of this IDataStore
	return s.Count
}
