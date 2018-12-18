# Example 3 - C++ Hello World #

## About ##
This example aims to give you an overview of Dragonboat's C++11 binding, which allows you to use Dragonboat in your C++ projects. You can safely ignore the C++ binding if your project is in Golang.

## Build ##
To build the C++ hello-world executable, first build the C++ binding in Dragonboat's top level directory:
```
cd $GOPATH/src/github.com/lni/dragonboat
make clean
make binding
```
then build the C++ hello-world executable:
```
cd examples/cpphelloworld
make
```

## Run ##
Start three instances of the C++ helloworld program on the same machine in three different terminals:

```
./cpphelloworld 1
```
```
./cpphelloworld 2
```
```
./cpphelloworld 3
```

This forms a Raft cluster with 3 nodes. The behaviour of this example is very similar to the Go [helloworld example](../helloworld). See the example description of the Go helloworld example for more details.

## Code ##
[datastore.h](datastore.h) and [datastore.cpp](datastore.cpp) implements the Data Store used by the application. Together with [plugin.h](plugin.h) and [plugin.cpp](plugin.cpp), a .so plugin named dragonboat-cpp-plugin-cpphelloworld.so can be built.

[helloworld.cpp](helloworld.cpp) contains the main function that starts a Raft based node and implements the features of this example. Please see comments in [helloworld.cpp](helloworld.cpp) for details.  
