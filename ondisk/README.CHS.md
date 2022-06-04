# 示例4 基于磁盘的状态机 #

## 关于 ##
本例用一个基于(Pebble)[https://github.com/cockroachdb/pebble]的分布式Key-Value数据库来展示Dragonboat的基于磁盘的状态机支持。

## 编译 ##
用下列命令编译本例的可执行文件 - 
```
cd $HOME/src/dragonboat-example
make ondisk
```

## 运行 ##
使用下列命令在同一台计算机的三个终端上启动三个本例程的实例。

```
./example-ondisk -replicaid 1
```
```
./example-ondisk -replicaid 2
```
```
./example-ondisk -replicaid 3
```

这将组建一个三节点的Raft集群，每个节点都由上述命令行命令中-replicaid所指定的ReplicaID值来标示。为求简易，本示例被
设定为需用三个节点且ReplicaID值必须为1, 2, 3。

您可以以下列两种格式输入一个命令以使用本例程 -

```
put key value
```
或 
```
get key
```

第一个命令将所指定的输入值Value设入所指定的键值Key，第二个命令通过查询底层的基于磁盘的状态机以返回键值Key所指向的值。

## 重新开始 ##
所有保存的数据均位于example-data的子目录内，可以手工删除这个example-data目录从而重新开始本例程。

## 代码 ##
在[diskkv.go](diskkv.go)中，DiskKV类型实现了statemachine.IOnDiskStateMachine这一接口。它使用Pebble作为存储引擎以在磁盘上存储所有状态机内容，这使它无需在每次重启后用快照Snapshot和已保存的Raft Log来恢复状态。这同时使得状态机管理下的数据量受磁盘大小制约，而不再受限于内存大小。

statemachine.IOnDiskStateMachine接口中的Open方法用以打开一个在磁盘上已存在的状态机并返回其最后一个已处理的Raft Log的index值。所有的实现了statemachine.IOnDiskStateMachine接口的类型均必须在保存状态机状态的同时，原子的同时保存其最后一个已处理的Raft Log的index值。内存那缓存的状态与磁盘同步，比如通过使用fsync()。本例中，我们始终使用Pebble的WriteBatch来原子地写入多个记录到底层的Pebble数据库中，包括最后一个已处理的Raft Log的index值。fsync()也始终在每次写以后被调用。

与基于statemachine.IStateMachine的状态机相比较，本例另一主要区别在于基于statemachine.IOnDiskStateMachine的状态机支持并发的读写。在状态机正在被Update更新时，Lookup和SaveSnapshot方法可以被同时并发的调用。而Lookup方法也可以在状态机正在被RecoverFromSnapshot方法恢复的时候被并发调用。

为了支持这样的并发访问，状态机快照的产生方法也与基于statemachine.IStateMachine的状态机有所不同。PrepareSnapshot方法首先被执行，它保存一个称为状态ID的对象，它用来标示状态机在某一具体时间点的状态。本例中，我们使用Pebble的快照功能来创建这样一个状态ID，并将其作为所产生的diskKVCtx对象的一部分返回。Update方法不会与PrepareSnapshot方法并发执行。接着，SaveSnapshot方法便可以与Update方法并发的执行了，SaveSnapshot会根据所提供的状态ID来产生状态机快照。本例中，我们遍历Pebble的快照中所涵盖的所有key-value对并将它们写入所提供的io.Writer中。

请参考[diskkv.go](diskkv.go)代码了解更详细的实现。

main函数在[main.go](main.go)中，它是我们创建NodeHost对象并启动所用的Raft组的地方。用户输入也在其中被处理，以允许用户通过put和get命令操作key-value数据。
