# SFS

## Description

 SFS(Smarter Function Scheduler) is an entirely user-space scheduler which carefully orchestrates existing Linux FIFO and CFS schedulers to approximate shortest remaining Time First(SRTF).

## Getting Started

### Dependencies

* Operating System: Ubuntu 20.04
* go version == 1.17.2
* schedtool version == v1.3.0

### Installing

#### SFS-Standalone-image

```
docker run --privileged --name test --mount type=bind,source="$(pwd)"/result,target=/result fuyuqi1995/sfs
```

#### SFS-Standalone

```go build && ./test.sh```

#### SFS-port OpenLambda

* openlambda
```
make imgs/lambda && make install #install openlambda
./create.sh && python cp_default.py && python replace.py #set up 64 workers
```
* SFS scheduler
```
go build & go run main.go #start SFS scheduler
```
* Http client
```
go build & go run run.go #submit requests
```
