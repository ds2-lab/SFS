# SFS
[![DOI](https://zenodo.org/badge/480597485.svg)](https://zenodo.org/badge/latestdoi/480597485)
## Description

 SFS(Smarter Function Scheduler) is an entirely user-space scheduler which carefully orchestrates existing Linux FIFO and CFS schedulers to approximate shortest remaining Time First(SRTF).

## Getting Started

### Dependencies

* Operating System: Ubuntu 20.04
* go version == 1.17.2
* schedtool version == v1.3.0
* python2.7

### System Configuration

* ulimit -n 1024000

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

### Evaluation

We have lots of sentivity analysis on evaluation section, here we provide how we make configurations.

#### CPU, Policy and trace

```
Usage of ./main:
  -n int
    	# of cpu cores (default 16)
  -o string
    	STCF optimal values (default "optimal.txt")
  -p string
    	scheduling policys: m:SFS; c:CFS, s: SRTF (default "m")
  -t string
    	trace
```
#### Adaptive time slice

To disable Adaptive time slice you can comment out line 351 in SFS-standalone/sfs.go and you could configure a fix time slice by change initial value on line 337.

#### Overhead handling

To disable hybrid strategy, you can comment out line 348 of SFS-standalone/sfs.go.

#### IO handling

To disable IO handling, you can comment out line 352 in SFS-standalone/sfs.go

#### IAT rate

You could configure IAT rate on line 45 in SFS-standalone/readTrace.go
