# 示例1 - Hello World #

## 关于 ##
本示例将给您一个Dragonboat的基本功能的概览，包括：

* 如何配置并开始一个新的NodeHost实例
* 如何开始一个Raft组
* 如何发起一个proposal以改变Raft状态
* 如何做强一致读
* 如何触发一个snapshot快照
* 如何完成Raft组的成员变更

## 编译 ##
执行下面的命令以编译helloworld程序：
```
cd $HOME/src/dragonboat-example
make helloworld
```

## 首次运行 ##
在同一主机上三个不同的终端terminals上启动三个helloworld进程：

```
./example-helloworld -nodeid 1
```
```
./example-helloworld -nodeid 2
```
```
./example-helloworld -nodeid 3
```
这将组建一个三节点的Raft集群，每个节点都由上述命令行命令中-nodeid所指定的NodeID值来标示。为求简易，本示例被设定为需用三个节点且NodeID值必须为1, 2, 3。

在任意一个终端terminal中输入一个字符串并按下键盘回车，这样的外部输入通常被称为一个提议（proposal），它会被复制到其它节点中。在不同的终端内重复几次这样的操作，此时请注意各个terminal上显示的计数信息和回显的消息的次序，在三个不同的terminals上它们应该是完全一致对应的。这演示了Dragonboat所实现的Raft的最核心功能：在分布的节点上达成共识。

```
2017-09-09 00:02:08.873755 I | transport: raft gRPC stream from [00128:00002] to [00128:00001] established
from Update() : msg: hi there, count:1
```
In the log message above, it mentions "[00128:00002]" and "[00128:00001]". They are used to identify two Raft nodes in log messages, the first one means it is a node from Raft Cluster with ClusterID 128, its NodeID is 2.

在上述log消息中提到了"[00128:00002]"和"[00128:00001]"。它们是用来在log消息中指代Raft节点的，比如第一个的意义是一个来自ClusterID为128的Raft组的NodeID为2的节点。

每个helloworld进程都有一个后台的goroutine以每10秒一次的频率做强一致读（Linearizable Read）。它查询已经应用的提议个数并将此结果显示到terminal上。

## 多数派 Quorum ##
只要任意多数的节点在正常工作，那么该Raft集群被认为具备多数派Quorum。对于一个三节点集群，只要至少两个任意节点在正常工作，那么系统便具备多数派Quorum。

现在选择任意一个terminal中按Ctrl+C以中止它的运行。在剩余的两个terminal中的任意一个输入消息，所输入的消息应该依旧可以被复制到另外一个节点上并被回显到terminal上。此时，使用Ctrl+C再次中止一个节点后再次输入一些消息，此时便不能继续看到所输入的消息被回显到terminal上，并且应该很快看到一个如下的超时错误提示消息：

```
failed to get a proposal session, timeout
2017-09-09 00:04:09.039257 W | transport: breaker localhost:63003 failed, connect and process failed: context deadline exceeded
```

## 重启节点 ##
我们选定一个已经被停止的节点，然后用首次启动它的时候一样的命令来重启它。比如，那个首次启动时候用了含有*-nodeid 2*字样命令的节点，可以用下面的命令在其原来的terminal中重启：
```
./example-helloworld -nodeid 2
```

此时，之前被复制的消息再次被应用并回显，并且它们的次序和之前完全一致。这是因为Dragonboat内部记录所有已提交commit成功的Proposals，在重启的时候它们将再次被应用到重启以后的进程中，以确保重启后的节点能恢复到之前失效时候的状态。

## 状态快照 Snapshot ##
在提交一定数量提议后，终端上会提示"snapshot captured"字样的消息。快照snapshot这一行为由raft.Config中的SnapshotEntries参数来控制。快照使得一个Raft程序的状态可以被快速整体捕获保存，从而可以被用于整体恢复程序的状态，而不需要逐一、大量地应用所有已经应用的提议。

## 成员变更 ##
在本例程中，Raft cluster的成员可以通过输入特殊消息来完成变更。下列输入到helloworld例程中的特殊消息将请求使得位于localhost:63100地址的具有node ID 4的一个新节点被加入到Raft cluster中。
```
add localhost:63100 4
```
一旦成员变更完成，便可以在一个新的终端中使用下列shell命令启动这个节点
```
./example-helloworld -nodeid 4 -addr localhost:63100 -join
```
命令行中的-join使得在调用nodehost.StartCluster时的join参数被设为true，以表示这是一个新的节点加入。

新节点将首先追赶上其余节点的进度，再开始接收新的复制过来的消息。

下列特殊消息将使得node ID为4的节点被从Raft cluster被移除掉。
```
remove 4
```
一旦被移除，该节点将无法继续接收复制过来的消息。

请注意，将一个已经移除的节点再次加入回到Raft cluster中是不被允许的。

## 重新开始 ##
所有保存的数据均位于example-data的子目录内，可以手工删除这个example-data目录从而重新开始本例程。

## 代码 ##
在[statemachine.go](statemachine.go)中，ExampleStateMachine结构体用来实现statemachine.IStateMachine接口。这个data store结构体用来实现应用程序自己的状态机逻辑。具体的内容将在[下一示例](README.DS.CHS.md)中展开。

[main.go](main.go)含有main()入口函数，在里面我们实例化了一个NodeHost实例，把所创建的Raft集群节点加入了这个实例。它同时使用多个goroutine来做用户输入消息和Ctrl+C信号的处理。同时请留意MakeProposal()函数的错误处理部分代码和注释。

makeMembershipChange()函数完成成员的变更。
