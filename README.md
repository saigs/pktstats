## Primitives for Non-blocking Concurrent Access

# Problem

What are some primitives that we could use for non-blocking concurrent access to update statistics, counters, etc. inline in a high-speed data path?

# Primitives

We consider three primitives here, all supported by Golang libraries. More to be added soon.

1. Mutex locks and synchronization
2. Atomic integer operations
3. Map-based packet size updates with the assumption that certain packet sizes are more common than others (e.g. 64, 128, etc.)

# Results

One such test result is presented below. With 10 threads, and 10,000 packets evenly spread across, following is the result of Golang benachmark on a Darwin OS.

```
$ go clean -testcache
$ go test -bench=. -benchtime=10s
DP processed 10000 pkts
goos: darwin
goarch: amd64
BenchmarkMutex-8     	DP processed 1010000 pkts
DP processed 101010000 pkts
DP processed 238310000 pkts
   13730	    874447 ns/op
DP processed 238320000 pkts
BenchmarkAtomic-8    	DP processed 239320000 pkts
DP processed 339320000 pkts
DP processed 509320000 pkts
   17000	    705469 ns/op
DP processed 509320000 pkts
BenchmarkWithMap-8   	DP processed 509320000 pkts
     100	 333719894 ns/op
PASS
ok  	_/Users/sais/bench	73.644s
```

# Intrepretation

- `BenchmarkWithMap` seems to perform the worst in terms of space and time complexity owing to the costly map update and read operations
- `BenchmarkMutex` seems to perform average with about 874 us/op
- `BenchmarkAtomic` seems to outperform mutex solution by a good 20% almost at all times

More results and references will be added upon further study.
