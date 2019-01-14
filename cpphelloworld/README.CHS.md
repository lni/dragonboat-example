# 示例4 - C++ Hello World #

## 关于 ##
如果您计划在Go项目中使用Dragonboat，您可以跳过本示例并忽略Dragonboat的C++ binding功能。

本示例将给您一个C++使用Dragonboat库的概览。

## 编译 ##
为了编译helloworld示例可执行文件，首先编译Dragonboat C++支持库。在Dragonboat的顶层目录：
```
cd $GOPATH/src/github.com/lni/dragonboat
make clean
make binding
sudo make install-binding
```
然后编译helloworld程序：
```
cd $GOPATH/src/github.com/lni/dragonboat-example/cpphelloworld
make
```

## 运行 ##
在同一主机上三个不同的终端terminals上启动三个cpphelloworld进程：

```
./cpphelloworld 1
```
```
./cpphelloworld 2
```
```
./cpphelloworld 3
```
这将组建了一个三节点的Raft集群。本示例程序的行为和Go版本的[示例1 - Hello World](../helloworld/README.CHS.md)十分相似。

## 代码 ##
[statemachine.h](statemachine.h)和[statemachine.cpp](statemachine.cpp)实现本应用所使用的StateMachine. 和[plugin.cpp](plugin.cpp)一起, 它们可以用来编译产生一个名为dragonboat-cpp-plugin-cpphelloworld.so的plugin。

[helloworld.cpp](helloworld.cpp)含main函数，它将启动一个Raft节点并实现本程序的功能与行为。[helloworld.cpp](helloworld.cpp)含详细注释，以解读各细节。
