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

/*
linearizable is an example program for building a linearizable state machine using dragonboat.
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/config"
)

var (
	datadir = "/tmp/dragonboat-example-linearizable"
	members = map[uint64]string{
		1: "localhost:61001",
		2: "localhost:61002",
		3: "localhost:61003",
	}
	httpAddr = []string{
		":8001",
		":8002",
		":8003",
	}
	clusterID uint64 = 128
)

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	for i, nodeAddr := range members {
		dir := fmt.Sprintf("%s/%d", datadir, i)
		if err := os.MkdirAll(dir, 0777); err != nil {
			panic(err)
		}
		log.Printf("Starting node %s", nodeAddr)
		nh, err := dragonboat.NewNodeHost(config.NodeHostConfig{
			RaftAddress:    nodeAddr,
			NodeHostDir:    dir,
			RTTMillisecond: 100,
		})
		if err != nil {
			panic(err)
		}
		fsm := NewLinearizableFSM()
		err = nh.StartConcurrentCluster(members, false, fsm, config.Config{
			NodeID:             uint64(i),
			ClusterID:          clusterID,
			ElectionRTT:        10,
			HeartbeatRTT:       1,
			CheckQuorum:        true,
			SnapshotEntries:    10,
			CompactionOverhead: 5,
		})
		if err != nil {
			panic(err)
		}
		go func(s *http.Server) {
			log.Fatal(s.ListenAndServe())
		}(&http.Server{
			Addr:    httpAddr[i-1],
			Handler: &handler{nh},
		})
	}
	<-stop
}
