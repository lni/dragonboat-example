# Example 2 - IDataStore #

## About ##
This example demonstrates how to implement your own [datastore.IDataStore](https://godoc.org/github.com/lni/dragonboat/datastore#IDataStore). 

## Code ##
[datastore.go](datastore.go) implements the [datastore.IDataStore](https://godoc.org/github.com/lni/dragonboat/datastore#IDataStore) interface required for managing application data in Dragonboat based applications. In this example, we show how the IDataStore instance used in the previous Helloworld example is implemented.

We first implement the Update() and Lookup() methods for handling incoming updates and queries. In this example IDataStore, there is a single integer named Count for representing the state of the IDataStore, the Count integer is increased every time when Update() is invoked. In the Update method, we also print out the input payload for demonstration purposes. In a real application, users are free to interpret the input data slice and update the state of their IDataStore accordingly. 

Lookup() is a read-only method for querying the state of the IDataStore, in the implementation of this example, we just put the Count value into a byte slice and return it. The input byte slice named key is usually used by applications to specify what to query. 

SaveSnapshot() and RecoverFromSnapshot() are used to implement snapshot save and load operations. Those in memory data maintained by your IDataStore can be saved to or read from the disk file backed io.Writer and io.Reader. The SaveSnapshot() method is also read-only, which means it is not suppose to change the state of the IDataStore. 

The Close() method should be considered as optional, as there is no guarantee that the Close() method will be called, data integrity of the IDataStore instance must not rely on the Close() method.  
