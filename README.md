# go-wavm

Experimental Golang wrapper for [WAVM](https://github.com/WAVM/WAVM).


## Benchmark results

These are some preliminary results, useful to compare against [go-wasm-benchmark](https://github.com/matiasinsaurralde/go-wasm-benchmark).

[WAVM](https://github.com/WAVM/WAVM) seems to be a good candidate for a Go wrapper, specially when considering [this](https://medium.com/fluence-network/a-standalone-webassembly-vm-benchmark-5300d534a04d) and its [proposed extensions](https://github.com/WAVM/WAVM#webassembly-10).

In the future I'm expecting more progress around Go based VMs. I will try to extend this small package as a learning process (`cgo` based for now).

```
goos: darwin
goarch: amd64
BenchmarkWAVMSum-8            	     500	   3143148 ns/op	     728 B/op	      35 allocs/op
BenchmarkWAVMSumReentrant-8   	 3000000	       507 ns/op	     144 B/op	       6 allocs/op
PASS
ok  	_/Users/matias/dev/wasm/go-wavm	4.032s
```