## About / [中文版](README.CHS.md) ##
This repo contains examples for [dragonboat](http://github.com/lni/dragonboat).

The master branch and the release-3.3 branch of this repo target Dragonboat's master and v3.3.x releases.

Go 1.17 or later releases with [Go module](https://github.com/golang/go/wiki/Modules) support is required.

## Notice ##

Programs provided here in this repo are examples - they are intentionally created in a more straight forward way to help users to understand the basics of the dragonboat library. They are not benchmark programs.

## Install ##

To download the example code to say $HOME/src/dragonboat-example:
```
$ cd $HOME/src
$ git clone https://github.com/lni/dragonboat-example
```
Build all examples:
```
$ cd $HOME/src/dragonboat-example
$ make
```

## Examples ##

Click links below for more details.

* [Example 1](helloworld) - Hello World
* [Example 2](helloworld/README.DS.md) - State Machine
* [Example 3](multigroup) - Multiple Raft Groups
* [Example 4](ondisk) - On Disk State Machine
* [Example 5](optimistic-write-lock) - Optimistic Write Lock

## Next Step ##
* [godoc](https://godoc.org/github.com/lni/dragonboat)
* Contribute code or report bugs for [dragonboat](http://github.com/lni/dragonboat)
