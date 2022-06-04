package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	dbsm "github.com/lni/dragonboat/v4/statemachine"
)

const (
	ResultCodeFailure = iota
	ResultCodeSuccess
	ResultCodeVersionMismatch
)

type Query struct {
	Key string
}

type Entry struct {
	Key string `json:"key"`
	Ver uint64 `json:"ver"`
	Val string `json:"val"`
}

func NewLinearizableFSM() dbsm.CreateConcurrentStateMachineFunc {
	return dbsm.CreateConcurrentStateMachineFunc(func(shardID, replicaID uint64) dbsm.IConcurrentStateMachine {
		return &linearizableFSM{
			shardID:   shardID,
			replicaID: replicaID,
			data:      map[string]interface{}{},
		}
	})
}

type linearizableFSM struct {
	shardID   uint64
	replicaID uint64
	data      map[string]interface{}
}

func (fsm *linearizableFSM) Update(entries []dbsm.Entry) ([]dbsm.Entry, error) {
	for i, ent := range entries {
		var entry Entry
		if err := json.Unmarshal(ent.Cmd, &entry); err != nil {
			return entries, fmt.Errorf("Invalid entry %#v, %w", ent, err)
		}
		if v, ok := fsm.data[entry.Key]; ok {
			// Reject entries with mismatched versions
			if v.(Entry).Ver != entry.Ver {
				data, _ := json.Marshal(v)
				entries[i].Result = dbsm.Result{
					Value: ResultCodeVersionMismatch,
					Data:  data,
				}
				continue
			}
		}
		entry.Ver = ent.Index
		fsm.data[entry.Key] = entry
		b, _ := json.Marshal(entry)
		entries[i].Result = dbsm.Result{
			Value: ResultCodeSuccess,
			Data:  b,
		}
	}

	return entries, nil
}

func (fsm *linearizableFSM) Lookup(e interface{}) (val interface{}, err error) {
	query, ok := e.(Query)
	if !ok {
		return nil, fmt.Errorf("Invalid query %#v", e)
	}
	val, _ = fsm.data[query.Key]

	return
}

func (fsm *linearizableFSM) PrepareSnapshot() (ctx interface{}, err error) {
	return
}

func (fsm *linearizableFSM) SaveSnapshot(ctx interface{}, w io.Writer, sfc dbsm.ISnapshotFileCollection, stopc <-chan struct{}) (err error) {
	b, err := json.Marshal(fsm.data)
	if err == nil {
		_, err = io.Copy(w, bytes.NewReader(b))
	}

	return
}

func (fsm *linearizableFSM) RecoverFromSnapshot(r io.Reader, sfc []dbsm.SnapshotFile, stopc <-chan struct{}) (err error) {
	return json.NewDecoder(r).Decode(fsm.data)
}

func (fsm *linearizableFSM) Close() (err error) {
	return
}
