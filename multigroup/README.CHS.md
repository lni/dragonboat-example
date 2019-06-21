# 示例3 - 多个Raft组 #

## 关于 ##
本示例展示如何使用多个Raft组。

## 编译 ##
执行下面的命令以编译示例程序：
```
cd $HOME/src/dragonboat-example
make multigroup
```

## 运行 ##
在同一主机上三个不同的终端terminals上启动三个示例程序的进程：

```
./example-multigroup -nodeid 1
```
```
./example-multigroup -nodeid 2
```
```
./example-multigroup -nodeid 3
```
这将组建两个三节点的Raft集群，每个节点都由上述命令行命令中-nodeid所指定的NodeID值来标示。为求简易，本示例被设定为需用三个节点且NodeID值必须为1, 2, 3。

与之前的helloworld示例一样，你可以在一个终端上输入一条消息，它将会被复制到别的节点上。在本示例中，如果你输入一个以问好"?"结尾的消息，那么它会被提议到第二个raft组中，其余不以"?"结尾的消息均提议至第一个raft组中。

本例中，我们仅使用两个Raft组，这是为了使得例子程序尽可能简单。在一个实际的应用中，用户可以轻易使用大量的Raft组。

## 重新开始 ##
所有保存的数据均位于example-data的子目录内，可以手工删除这个example-data目录从而重新开始本例程。

## 代码 ##
[main.go](main.go)含有main()入口函数，在里面我们实例化了一个NodeHost实例，把所创建的两个Raft集群节点加入了这个实例。它同时实现了根据用户输入不同而对不同Raft组来提交提议的逻辑。注意代码中各阶段是如何指定ClusterID值的。
