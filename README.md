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

disable Adaptive time slice by adding comment on line 351 in SFS-standaline/sfs.go and we could configure a fix time slice by change initial value on line 337.

#### Overhead handling

disable hyhird strategy by adding comment on line 348 in SFS-standaline/sfs.go

#### IO handling

disable IO handling by adding comment on line 352 in SFS-standaline/sfs.go

#### IAT rate

could configure IAT rate on line 45 in SFS-standaline/readTrace.go
