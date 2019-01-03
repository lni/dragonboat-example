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
helloworld is an example program for dragonboat.
*/
package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/lni/dragonboat"
	"github.com/lni/dragonboat-example/utils"
	"github.com/lni/dragonboat/config"
	"github.com/lni/dragonboat/logger"
)

const (
	exampleClusterID uint64 = 128
)

var (
	// initial nodes count is fixed to three, their addresses are also fixed
	addresses = []string{
		"localhost:63001",
		"localhost:63002",
		"localhost:63003",
	}
	errNotMembershipChange = errors.New("not a membership change request")
)

// makeMembershipChange makes membership change request.
func makeMembershipChange(nh *dragonboat.NodeHost,
	cmd string, addr string, nodeID uint64) {
	var rs *dragonboat.RequestState
	var err error
	if cmd == "add" {
		// orderID is ignored in standalone mode
		rs, err = nh.RequestAddNode(exampleClusterID, nodeID, addr, 0, 5000)
	} else if cmd == "remove" {
		rs, err = nh.RequestDeleteNode(exampleClusterID, nodeID, 0, 5000)
	} else {
		panic("unknown cmd")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "membership change failed, %v\n", err)
	}
	select {
	case r := <-rs.CompletedC:
		if r.Completed() {
			fmt.Fprintf(os.Stdout, "membership change completed successfully\n")
		} else {
			fmt.Fprintf(os.Stderr, "membership change failed\n")
		}
	}
}

// splitMembershipChangeCmd tries to parse the input string as membership change
// request. ADD node request has the following expected format -
// add localhost:63100 4
// REMOVE node request has the following expected format -
// remove 4
func splitMembershipChangeCmd(v string) (string, string, uint64, error) {
	parts := strings.Split(v, " ")
	if len(parts) == 2 || len(parts) == 3 {
		cmd := strings.ToLower(strings.TrimSpace(parts[0]))
		if cmd != "add" && cmd != "remove" {
			return "", "", 0, errNotMembershipChange
		}
		addr := ""
		var nodeIDStr string
		var nodeID uint64
		var err error
		if cmd == "add" {
			addr = strings.TrimSpace(parts[1])
			nodeIDStr = strings.TrimSpace(parts[2])
		} else {
			nodeIDStr = strings.TrimSpace(parts[1])
		}
		if nodeID, err = strconv.ParseUint(nodeIDStr, 10, 64); err != nil {
			return "", "", 0, errNotMembershipChange
		}
		return cmd, addr, nodeID, nil
	}
	return "", "", 0, errNotMembershipChange
}

func main() {
	nodeID := flag.Int("nodeid", 1, "NodeID to use")
	addr := flag.String("addr", "", "Nodehost address")
	join := flag.Bool("join", false, "Joining a new node")
	flag.Parse()
	if len(*addr) == 0 && *nodeID != 1 && *nodeID != 2 && *nodeID != 3 {
		fmt.Fprintf(os.Stderr, "node id must be 1, 2 or 3 when address is not specified\n")
		os.Exit(1)
	}
	// https://github.com/golang/go/issues/17393
	if runtime.GOOS == "darwin" {
		signal.Ignore(syscall.Signal(0xd))
	}
	peers := make(map[uint64]string)
	for idx, v := range addresses {
		// key is the NodeID, NodeID is not allowed to be 0
		// value is the raft address
		peers[uint64(idx+1)] = v
	}
	var nodeAddr string
	if len(*addr) != 0 {
		nodeAddr = *addr
	} else {
		nodeAddr = peers[uint64(*nodeID)]
	}
	fmt.Fprintf(os.Stdout, "node address: %s\n", nodeAddr)
	// change the log verbosity
	logger.GetLogger("raft").SetLevel(logger.ERROR)
	logger.GetLogger("rsm").SetLevel(logger.WARNING)
	logger.GetLogger("transport").SetLevel(logger.WARNING)
	logger.GetLogger("grpc").SetLevel(logger.WARNING)
	// config for raft
	rc := config.Config{
		NodeID:             uint64(*nodeID),
		ClusterID:          exampleClusterID,
		ElectionRTT:        5,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}
	datadir := filepath.Join(
		"example-data",
		"helloworld-data",
		fmt.Sprintf("node%d", *nodeID))
	// config for the nodehost
	// by default, insecure transport is used, you can choose to use Mutual TLS
	// Authentication to authenticate both servers and clients. To use Mutual
	// TLS Authentication, set the MutualTLS field in NodeHostConfig to true, set
	// the CAFile, CertFile and KeyFile fields to point to the path of your CA
	// file, certificate and key files.
	// by default, TCP based RPC module is used, set the RaftRPCFactory field in
	// NodeHostConfig to rpc.NewRaftGRPC (github.com/lni/dragonboat/plugin/rpc) to
	// use gRPC based transport. To use gRPC based RPC module, you need to install
	// the gRPC library first -
	//
	// $ go get -u google.golang.org/grpc
	//
	nhc := config.NodeHostConfig{
		WALDir:         datadir,
		NodeHostDir:    datadir,
		RTTMillisecond: 200,
		RaftAddress:    nodeAddr,
		// RaftRPCFactory: rpc.NewRaftGRPC,
	}
	nh := dragonboat.NewNodeHost(nhc)
	// peers is always populated here, but its content is only used when you are
	// starting a brand new cluster for the first time. Its content is ignored
	// when you are restarting a stopped node or joining a new node
	if err := nh.StartCluster(peers, *join, NewExampleStateMachine, rc); err != nil {
		fmt.Fprintf(os.Stderr, "failed to add cluster, %v\n", err)
		os.Exit(1)
	}
	raftStopper := utils.NewStopper()
	consoleStopper := utils.NewStopper()
	ch := make(chan string, 16)
	consoleStopper.RunWorker(func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				close(ch)
				return
			}
			if s == "exit\n" {
				raftStopper.Stop()
				// no data will be lost/corrupted if nodehost.Stop() is not called
				nh.Stop()
				return
			}
			ch <- s
		}
	})
	raftStopper.RunWorker(func() {
		// this goroutine makes a linearizable read every 10 second. it returns the
		// Count value maintained in IStateMachine. see datastore.go for details.
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				result, err := nh.SyncRead(ctx, exampleClusterID, []byte{})
				cancel()
				if err == nil {
					var count uint64
					count = binary.LittleEndian.Uint64(result)
					fmt.Fprintf(os.Stdout, "count: %d\n", count)
				}
			case <-raftStopper.ShouldStop():
				return
			}
		}
	})
	raftStopper.RunWorker(func() {
		// use a NO-OP client session here
		// check the example in godoc to see how to use a regular client session
		cs := nh.GetNoOPSession(exampleClusterID)
		for {
			select {
			case v, ok := <-ch:
				if !ok {
					return
				}
				// remove the \n char
				msg := strings.Replace(v, "\n", "", 1)
				if cmd, addr, nodeID, err := splitMembershipChangeCmd(msg); err == nil {
					// input is a membership change request
					makeMembershipChange(nh, cmd, addr, nodeID)
				} else {
					// input is a regular message need to be proposed
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					// make a proposal to update the IStateMachine instance
					_, err := nh.SyncPropose(ctx, cs, []byte(msg))
					cancel()
					if err != nil {
						fmt.Fprintf(os.Stderr, "SyncPropose returned error %v\n", err)
					}
				}
			case <-raftStopper.ShouldStop():
				return
			}
		}
	})
	raftStopper.Wait()
}
