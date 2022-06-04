# Example 4 - On Disk State Machine #

## About ##
This example uses a [Pebble](https://github.com/cockroachdb/pebble) based distributed key-value store to demonstrate the on disk state machine support in Dragonboat.

## Build ##
To build the executable -
```
cd $HOME/src/dragonboat-example
make ondisk
```

## Run This Example ##
Start three instances of the example program on the same machine in three different terminals:

```
./example-ondisk -replicaid 1
```
```
./example-ondisk -replicaid 2
```
```
./example-ondisk -replicaid 3
```
This forms a Raft group with 3 replicas, each of them is identified by the ReplicaID specified on the command line. For simplicity, this example program has been hardcoded to use 3 nodes and their node id values are 1, 2 and 3.

You can type in a message in one of the two following formats - 
```
put key value
```
or 
```
get key
```

The first command above sets the specified input value to key, the second command queries the underlying on disk state machine and returns the value associated with key. 

## Start Over ##
All saved data is saved into the example-data folder, you can delete this example-data folder and restart all processes to start over again.

## Code ##
In [diskkv.go](diskkv.go), the DiskKV struct implements the statemachine.IOnDiskStateMachine interface. It employs Pebble as its on disk storage engine to store all state machine managed data, it thus doesn't need to be restored from snapshot or saved Raft logs after each reboot. This also ensures that the total amount of data that can be managed by the state machine is limited by available disk capacity rather than memory size. 

The Open method of the statemachine.IOnDiskStateMachine interface opens existing on disk state machine and returns the index of the last updated Raft log entry. It is important for all statemachine.IOnDiskStateMachine implmentations to atomically persist the index of the last updated Raft log entry together with the outcome of the update operation when updating such on disk state machines. In-core state should also be synchronized with disk, e.g. using fsync(). In this example, we always use Pebble's WriteBatch type to atomically write incoming records, including the the index of the last updated Raft log entry, to the underlying Pebble database. fsync() is invoked by Pebble at the end of each write.

Compared with statemachine.IStateMachine based state machine, another major difference is that concurrent read and write are supported by statemachine.IOnDiskStateMachine based on disk state machines. The Lookup and the SaveSnapshot method can be concurrently invoked when the state machine is being updated by the Update method. The Lookup method can also be invoked when the state machine is being resotred by the RecoverFromSnapshot method. 

To support the above described concurrent access to statemachine.IOnDiskStateMachine types, the way how state machine snapshot is saved is also different from previously described statemachine.IStateMachine types. The PrepareSnapshot method will be first invoked to capture and return a so called state identifier object that can describe the point in time state of the state machine. In this example, we take a Pebble snapshot and return it as a part of the generated diskKVCtx instance. Update is not allowed by the system when PrepareSnapshot is being invoked. SaveSnapshot is then invoked concurrent to the Update method to actually save the point in time state of the state machine identified by the provided state identifier. In this example, we iterate over all key-value pairs covered in the Pebble snapshot and write all of them to the provided io.Writer.

See godoc in [diskkv.go](diskkv.go) for more detials.

The main function can be found in [main.go](main.go), it is the place where we instantiated the NodeHost instance, added the created example Ragroupp to it. User inputs are also handled here to allow users to put or get key-value pairs.
