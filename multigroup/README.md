# Example 3 - Multiple Raft Groups #

## About ##
This example aims to give you an overview on how to use multiple raft groups in Dragonboat.

## Build ##
To build the executable -
```
cd $HOME/src/dragonboat-example
make multigroup
```

## Let's try it ##
Start three instances of the example program on the same machine in three different terminals:

```
./example-multigroup -nodeid 1
```
```
./example-multigroup -nodeid 2
```
```
./example-multigroup -nodeid 3
```
This forms two Raft clusters each with 3 nodes, nodes are identified by the NodeID values specified on the command line while two Raft groups are identified by the ClusterID values harded coded in main.go. For simplicity, this example program has been hardcoded to use 3 nodes for each raft group and their node id values are 1, 2 and 3.

Similar to the previous helloworld example, you can type in a message in one of the terminals and your input message will be replicated to other nodes. In this example, if you type something ends with a question mark "?", then that particular message is going to be proposed to the second raft group, while messages without the question mark at the end are proposed to the first raft group. 

Note that we use two Raft groups here for simplicity, in a real world application, users can scale to much larger number of Raft groups.  

## Start Over ##
All saved data is saved into the example-data folder, you can delete this example-data folder and restart all processes to start over again.

## Code! ##
[main.go](main.go) contains the main() function, it is the place where we instantiated the NodeHost instance, added the created two raft groups to it. It also implements the logic of making proposals to one of the two available raft groups based on user input. Note how ClusterID value is specified during different stages of this demo. 
