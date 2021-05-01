## Primitives for Non-blocking Concurrent Access

# Problem

What are some primitives that we could use for non-blocking concurrent access to update statistics, counters, etc. inline in a high-speed data path?

# Primitives

We consider three primitives here, all supported by Golang libraries. More to be added soon.

1. Mutex locks and synchronization
2. Atomic integer operations
3. Map-based packet size updates with the assumption that certain packet sizes are more common than others (e.g. 64, 128, etc.)
4. Concurrent and atomic integer operations

# Usage

To build the program:

```
make all
```

To execute the program, simply type:

```
make run
```

-or-

```
$ ./build/server -h
Usage of ./build/server:
  -mode int
    	mutex=0 atomic=1 map=2 concurrent=3
  -packets int
    	#packets to send (default 10000)
  -threads int
    	#concurrent threads (default 10)
```

Here is a sample run.

```
$ ./build/server -mode 0 -threads 10 -packets 1000
2021/05/01 00:30:11 Initializing datapath..
2021/05/01 00:30:11 Populating datapath queues [10]
2021/05/01 00:30:11 datapath pkt queue populated: thread 0, pkts 100
2021/05/01 00:30:11 datapath pkt queue populated: thread 1, pkts 100
2021/05/01 00:30:11 datapath pkt queue populated: thread 2, pkts 100
2021/05/01 00:30:11 datapath pkt queue populated: thread 3, pkts 100
2021/05/01 00:30:11 datapath pkt queue populated: thread 4, pkts 100
2021/05/01 00:30:11 datapath pkt queue populated: thread 5, pkts 100
2021/05/01 00:30:11 datapath pkt queue populated: thread 6, pkts 100
2021/05/01 00:30:11 datapath pkt queue populated: thread 7, pkts 100
2021/05/01 00:30:11 datapath pkt queue populated: thread 8, pkts 100
2021/05/01 00:30:11 datapath pkt queue populated: thread 9, pkts 100
2021/05/01 00:30:11 Init datapath: type mutex, threads 10, total pkts 1000
2021/05/01 00:30:11 Running datapath threads 10, pkts 1000
2021/05/01 00:30:11 RunDatapath took 898.917µs
2021/05/01 00:30:11 average packet size 33227
```

# Tests

Testcases can be run using the following:

```
make test
```

# Results

One such test result is presented below. With 10 threads, and 10,000 packets evenly spread across, following is the result of Golang benachmark on a Darwin OS.

Mutex Context

```
./build/server -mode 0 -threads 10 -packets 10000
2021/05/01 00:31:25 Initializing datapath..
2021/05/01 00:31:25 Populating datapath queues [10]
2021/05/01 00:31:25 datapath pkt queue populated: thread 0, pkts 1000
2021/05/01 00:31:25 datapath pkt queue populated: thread 1, pkts 1000
2021/05/01 00:31:25 datapath pkt queue populated: thread 2, pkts 1000
2021/05/01 00:31:25 datapath pkt queue populated: thread 3, pkts 1000
2021/05/01 00:31:25 datapath pkt queue populated: thread 4, pkts 1000
2021/05/01 00:31:25 datapath pkt queue populated: thread 5, pkts 1000
2021/05/01 00:31:25 datapath pkt queue populated: thread 6, pkts 1000
2021/05/01 00:31:25 datapath pkt queue populated: thread 7, pkts 1000
2021/05/01 00:31:25 datapath pkt queue populated: thread 8, pkts 1000
2021/05/01 00:31:25 datapath pkt queue populated: thread 9, pkts 1000
2021/05/01 00:31:25 Init datapath: type mutex, threads 10, total pkts 10000
2021/05/01 00:31:25 Running datapath threads 10, pkts 10000
2021/05/01 00:31:25 RunDatapath took 1.7415ms
2021/05/01 00:31:25 average packet size 30147
```

Atomic Context

```
./build/server -mode 1 -threads 10 -packets 10000
2021/05/01 00:31:39 Initializing datapath..
2021/05/01 00:31:39 Populating datapath queues [10]
2021/05/01 00:31:39 datapath pkt queue populated: thread 0, pkts 1000
2021/05/01 00:31:39 datapath pkt queue populated: thread 1, pkts 1000
2021/05/01 00:31:39 datapath pkt queue populated: thread 2, pkts 1000
2021/05/01 00:31:39 datapath pkt queue populated: thread 3, pkts 1000
2021/05/01 00:31:39 datapath pkt queue populated: thread 4, pkts 1000
2021/05/01 00:31:39 datapath pkt queue populated: thread 5, pkts 1000
2021/05/01 00:31:39 datapath pkt queue populated: thread 6, pkts 1000
2021/05/01 00:31:39 datapath pkt queue populated: thread 7, pkts 1000
2021/05/01 00:31:39 datapath pkt queue populated: thread 8, pkts 1000
2021/05/01 00:31:39 datapath pkt queue populated: thread 9, pkts 1000
2021/05/01 00:31:39 Init datapath: type atomic, threads 10, total pkts 10000
2021/05/01 00:31:39 Running datapath threads 10, pkts 10000
2021/05/01 00:31:39 RunDatapath took 1.139958ms
2021/05/01 00:31:39 average packet size 30810
```

Map Context

```
./build/server -mode 2 -threads 10 -packets 10000
2021/05/01 00:32:01 Initializing datapath..
2021/05/01 00:32:01 Populating datapath queues [10]
2021/05/01 00:32:01 datapath pkt queue populated: thread 0, pkts 1000
2021/05/01 00:32:01 datapath pkt queue populated: thread 1, pkts 1000
2021/05/01 00:32:01 datapath pkt queue populated: thread 2, pkts 1000
2021/05/01 00:32:01 datapath pkt queue populated: thread 3, pkts 1000
2021/05/01 00:32:01 datapath pkt queue populated: thread 4, pkts 1000
2021/05/01 00:32:01 datapath pkt queue populated: thread 5, pkts 1000
2021/05/01 00:32:01 datapath pkt queue populated: thread 6, pkts 1000
2021/05/01 00:32:01 datapath pkt queue populated: thread 7, pkts 1000
2021/05/01 00:32:01 datapath pkt queue populated: thread 8, pkts 1000
2021/05/01 00:32:01 datapath pkt queue populated: thread 9, pkts 1000
2021/05/01 00:32:01 Init datapath: type map, threads 10, total pkts 10000
2021/05/01 00:32:01 Running datapath threads 10, pkts 10000
2021/05/01 00:32:01 RunDatapath took 64.638459ms
2021/05/01 00:32:01 average packet size 33387
```

Concurrent Atomic Context

```
./build/server -mode 3 -threads 10 -packets 10000
2021/05/01 00:32:17 Initializing datapath..
2021/05/01 00:32:17 Populating datapath queues [10]
2021/05/01 00:32:17 datapath pkt queue populated: thread 0, pkts 1000
2021/05/01 00:32:17 datapath pkt queue populated: thread 1, pkts 1000
2021/05/01 00:32:17 datapath pkt queue populated: thread 2, pkts 1000
2021/05/01 00:32:17 datapath pkt queue populated: thread 3, pkts 1000
2021/05/01 00:32:17 datapath pkt queue populated: thread 4, pkts 1000
2021/05/01 00:32:17 datapath pkt queue populated: thread 5, pkts 1000
2021/05/01 00:32:17 datapath pkt queue populated: thread 6, pkts 1000
2021/05/01 00:32:17 datapath pkt queue populated: thread 7, pkts 1000
2021/05/01 00:32:17 datapath pkt queue populated: thread 8, pkts 1000
2021/05/01 00:32:17 datapath pkt queue populated: thread 9, pkts 1000
2021/05/01 00:32:17 Init datapath: type concurrent, threads 10, total pkts 10000
2021/05/01 00:32:17 Running datapath threads 10, pkts 10000
2021/05/01 00:32:17 RunDatapath took 352.208µs
2021/05/01 00:32:17 average packet size 32539
```

# Benchmarks

Benchmark tests can be run using the following:

```
make benchmark
```

# Intrepretation

- `BenchmarkWithMap` seems to perform the worst in terms of space and time complexity owing to the costly map update and read operations
- `BenchmarkMutex` seems to perform average with about 874 us/op
- `BenchmarkAtomic` seems to outperform mutex solution by a good 20% almost at all times
- `BenchmarkConcurrent` seems to outperform the atomic by a phenomenal 3x degree.

More results and references will be added upon further experiments.
