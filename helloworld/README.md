# Example 1 - Hello World #

## About ##
This example aims to give you an overview of Dragonboat's basic features, such as:

* how to config and start a NodeHost instance
* how to start a Raft group (also known as Raft shard)
* how to make proposals to update Raft group state
* how to perform linearizable read from Raft group
* how to trigger snapshot to be generated
* how to perform Raft group membership change

## Build ##
To build the helloworld executable -
```
cd $HOME/src/dragonboat-example
make helloworld
```

## First Run ##
Start three instances of the helloworld program on the same machine in three different terminals:

```
./example-helloworld -replicaid 1
```
```
./example-helloworld -replicaid 2
```
```
./example-helloworld -replicaid 3
```
This forms a Raft group with 3 nodes, each of them is identified by the NodeID specified on the command line. For simplicity, this example program has been hardcoded to use 3 nodes and their node id values are 1, 2 and 3.

Type in a message (also known as a proposal) and press enter in any one of those terminals, the message will be replicated to other nodes. You can try repeating this a few times, probably in different terminals. Note the count numbers and the order of messages, they should be identical across all three terminals. See example log messages below. This demonstrates the core feature of Raft implemented by Dragonboat - reaching consensus in a distributed environment.

```
2017-09-09 00:02:08.873755 I | transport: raft gRPC stream from [00128:00002] to [00128:00001] established
from Update() : msg: hi there, count:1
```
In the log message above, it mentions "[00128:00002]" and "[00128:00001]". They are used to identify two Raft nodes in log messages, the first one means it is a node from Raft group with ShardID 128, its NodeID is 2. 

Each helloworld process has a background goroutine performing linearizable read every 10 seconds. It queries the number of applied proposals and print out the result to the terminal. 

## Quorum ##
As long as the majority of nodes in the Raft group are available, the group is said to has the quorum. For such a 3-nodes Raft group, any two nodes need to be available to have the quorum.

Now press CTRL+C in any one of the terminal to stop the running exmple-helloworld program. Type in some messsages again in any of the remaining two terminals, you should still be able to have the message replicated across to the other node as the Raft group still has the quorum. Let's press CTRL+C to stop another instance, you should notice that further input messages will no longer be printed back on the remaining terminal and you are expected to see a timeout message. 

```
failed to get a proposal session, timeout
2017-09-09 00:04:09.039257 W | transport: breaker localhost:63003 failed, connect and process failed: context deadline exceeded
```

## Restart a node ##
Let's pick a stopped instance and restart it using the exact same command, e.g. for the one which we previously started with the *-replicaid 2* command line option, it can be restarted using the command below - 
```
./example-helloworld -replicaid 2
```

Previously replicated messages are printed back onto the terminal again in the same order as they were initially replicated across. Dragonboat internally records the state of the node and all updates to make sure it can be correctly restored after restart. 

## Snapshotting ##
After proposing several more messages, there will be logs mentioning that "snapshot captured". This is controled by the SnapshotEntries parameter specified in the raft.Config object. Snapshots can be used to restore the state of the program without requiring every single proposed messages to be applied one by one.

## Membership Change ##
In this example program, Raft group membership can be changed by inputing some messages with special format. The following special message causes a new node with node ID 4 running at localhost:63100 to be added to the Raft group.

```
add localhost:63100 4
```
Once the membership change is completed, you can start this recently added node in a new terminal using the following shell command - 
```
./example-helloworld -replicaid 4 -addr localhost:63100 -join
```
The -join option tells the progress to set the join parameter to be true when calling the nodehost.StartReplica() function. 

The new node will catch up with other nodes and start to receive replicated messages.

The following special message removes the node with node ID 4 from the cluster.
```
remove 4
```
Once removed, node will stop receiving further replicated messages.

Note that adding a previously removed node back to the cluster is not allowed.

## Start Over ##
All saved data is saved into the example-data folder, you can delete this example-data folder and restart all processes to start over again.

## Code! ##
In [statemachine.go](statemachine.go), we have this ExampleStateMachine struct which implements the statemachine.IStateMachine interface. This is the data store struct for implementing the application state machine. We will leave the details to the [next example](README.DS.md). 

[main.go](main.go) contains the main() function, it is the place where we instantiated the NodeHost instance, added the created example Raft group to it. It uses multiple goroutines for input and signal handling. Updates to the IStateMachine instance is achieved by making proposals.

makeMembershipChange() shows how to make membership changes, including add or remove nodes.
