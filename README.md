# SFS

## Description

 SFS(Smarter Function Scheduler) is an entirely user-space scheduler which carefully orchestrates existing Linux FIFO and CFS schedulers to approximate shortest remaining Time First(SRTF).

## Getting Started

### Dependencies

* Operating System: Ubuntu 20.04
* go version == 1.17.2
* schedtool version == v1.3.0

### Installing

#### SFS-Standalone

```go build && ./test.sh```

#### SFS-port OpenLambda

##### openlambda

make imgs/lambda && make install
* How/where to download your program
* Any modifications needed to be made to files/folders
