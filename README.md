# SFS: Smart OS Scheduling for Serverless Functions
[![DOI](https://zenodo.org/badge/480597485.svg)](https://zenodo.org/badge/latestdoi/480597485)

## Description

SFS (Smart Function Scheduler) is a user-space OS scheduler that carefully orchestrates existing Linux kernel-space scheduling policies, FIFO and CFS, to approximate the optimal offline scheduling policy Shortest Remaining Time First (SC'22).

To access the paper from ACM DL: https://dl.acm.org/doi/abs/10.5555/3571885.3571940

To access our arXiv preprint: https://arxiv.org/abs/2209.01709 

## Getting Started

### Dependencies

* Operating System: Ubuntu 20.04
* go version == 1.17.2
* schedtool version == v1.3.0

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

### To cite SFS
```
@inproceedings{10.5555/3571885.3571940,
author = {Fu, Yuqi and Liu, Li and Wang, Haoliang and Cheng, Yue and Chen, Songqing},
title = {SFS: Smart OS Scheduling for Serverless Functions},
year = {2022},
isbn = {9784665454445},
publisher = {IEEE Press},
abstract = {Serverless computing enables a new way of building and scaling cloud applications by allowing developers to write fine-grained serverless or cloud functions. The execution duration of a cloud function is typically short---ranging from a few milliseconds to hundreds of seconds. However, due to resource contentions caused by public clouds' deep consolidation, the function execution duration may get significantly prolonged and fail to accurately account for the function's true resource usage. We observe that the function duration can be highly unpredictable with huge amplification of more than 50\texttimes{} for an open-source FaaS platform (OpenLambda). Our experiments show that the OS scheduling policy of cloud functions' host server can have a crucial impact on performance. The default Linux scheduler, CFS (Completely Fair Scheduler), being oblivious to workloads, frequently context-switches short functions, causing a turnaround time that is much longer than their service time.We propose SFS (Smart Function Scheduler), which works entirely in the user space and carefully orchestrates existing Linux FIFO and CFS schedulers to approximate Shortest Remaining Time First (SRTF). SFS uses two-level scheduling that seamlessly combines a new FILTER policy with Linux CFS, to trade off increased duration of long functions for significant performance improvement for short functions. We implement SFS in the Linux user space and port it to OpenLambda. Evaluation results show that SFS significantly improves short functions' duration with a small impact on relatively longer functions, compared to CFS.},
booktitle = {Proceedings of the International Conference on High Performance Computing, Networking, Storage and Analysis},
articleno = {42},
numpages = {16},
keywords = {performance evaluation, operating systems, cloud computing},
location = {Dallas, Texas},
series = {SC '22}
}
```
