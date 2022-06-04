## 关于 ##
本repo含[dragonboat](http://github.com/lni/dragonboat)项目的示例程序

本repo的master branch和release-3.3 branch针对dragonboat repo的master branch和各v3.3.x发布版。

需Go 1.17或更新的带[Go module](https://github.com/golang/go/wiki/Modules)支持的Go版本。

## 注意事项 ##
本repo中的程序均为示例，为了便于向用户展现dragonboat的基本用途，它们被刻意以最简单的方式实现而忽略了基本所有性能考虑。这些示例程序不能用于跑分用途。

## 安装 ##
假设计划下载例程代码到$HOME/src/dragonboat-example：
```
$ cd $HOME/src
$ git clone https://github.com/lni/dragonboat-example
```
编译所有例程：
```
$ cd $HOME/src/dragonboat-example
$ make
```

## 示例 ##

点选下列链接以获取具体例程信息。

* [示例 1](helloworld) - Hello World
* [示例 2](helloworld/README.DS.md) - State Machine 状态机
* [示例 3](multigroup/README.CHS.md) - 多个Raft组
* [示例 4](ondisk/README.CHS.md) - 基于磁盘的State Machine 状态机

## 下一步 ##
* [godoc](https://godoc.org/github.com/lni/dragonboat)
* 为[dragonboat](http://github.com/lni/dragonboat)项目贡献代码或报告bug！
