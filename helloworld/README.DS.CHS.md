# 示例2 - IDataStore #

## 关于 ##
本示例介绍如何实现一个[datastore.IDataStore](https://godoc.org/github.com/lni/dragonboat/datastore#IDataStore)。

## 代码 ##
[datastore.go](datastore.go)实现了Dragonboat应用所需要的用于管理用户数据的[datastore.IDataStore](https://godoc.org/github.com/lni/dragonboat/datastore#IDataStore)接口。在本例子中，我们介绍上一个Helloworld示例中所使用的IDataStore实例是如何实现的。

首先需要实现Update()与Lookup()方法，用于处理收到的更新与查询请求。在本示例IDataStore中，它只有一个名为Count的整数计数器，每当Update()被调用时，该计数器会被递增。在Update()方法中，我们同时打印出所收到的输入参数用以演示目的。在一个实际应用里，用户可根据这样的输入参数相应更新IDataStore的状态。

Lookup()是一个只读的用于查询IDataStore的方法，在本示例的实现中，仅简单的把计数器的值放入一个byte slice中并返回。名为query的byte slice参数通常是应用提供的用于表明需查询内容的输入参数。

SaveSnapshot()和RecoverFromSnapshot() 用来实现快照的保存与读取。快照需要能包含IDataStore的完整状态。IDataStore所维护的内存内的数据可以通过所提供的以磁盘文件为后台存储的io.Writer和io.Reader来保存与恢复。请注意，SaveSnapshot()也是一个只读的方法，它不应该改变IDataStore的状态。

Close()的可以被认为是可选的。因为系统并不保证Close回被最终调用，因此IDataStore的数据完整性不能依赖于Close()方法。
