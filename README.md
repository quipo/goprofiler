# Golang profiling utility

[![GoDoc](https://godoc.org/github.com/quipo/goprofiler/profiler?status.png)](http://godoc.org/github.com/quipo/goprofiler/profiler)

## Introduction

Simple profiling support package for Go. 

Inspired by [Dave Cheney's library](https://github.com/davecheney/profile), with extra functionality and automation.

The library can take CPU, Memory, Block and Goroutine profiles. 
It can either take a single snapshot, or can take snapshots at regular intervals.

## Installation

    go get github.com/quipo/goprofiler/profiler

## Sample usage

```go
package main

import (
	"time"

	"github.com/quipo/goprofiler/profiler"
)

func main() {
	pprofConf := profiler.Config{
		Prefix: "/tmp/myapp.",
		CPU: true,
		Memory: true,
		Block: false,
		Goroutine: false,
		Interval: "15s",      // one snapshot every 15 seconds,
		MemoryProfileRate: 1 // collection information about all allocations
	}
	prof := profiler.NewProfiler(pprofConf)
	go prof.Run()
	
	// your app's code here...

	// you can take a snapshot at any point in the code:
	prof.TakeSnapshot()

	// take one last snapshot and clean resources
	prof.Stop()
}
```

## Useful links

* [Debugging performance issues in Go programs (intel)](https://software.intel.com/en-us/blogs/2014/05/10/debugging-performance-issues-in-go-programs)
* [Intel® Performance Counter Monitor](https://software.intel.com/en-us/articles/intel-performance-counter-monitor-a-better-way-to-measure-cpu-utilization)
* [runtime/pprof WriteTo](http://golang.org/pkg/runtime/pprof/#Profile.WriteTo)
* [Description of testing flags](http://golang.org/cmd/go/#hdr-Description_of_testing_flags)
* [Monitoring A Production Golang Server With Memstats](http://pythonic.zoomquiet.io/data/20131112090955/index.html)
* [Go performance tales](https://www.datadoghq.com/2014/04/go-performance-tales/)
* [Cache coherency primer](http://fgiesen.wordpress.com/2014/07/07/cache-coherency/)
* [CPU Cache flushing fallacy](http://mechanical-sympathy.blogspot.dk/2013/02/cpu-cache-flushing-fallacy.html)
* [cmd/gc: make liveness ~10x faster](https://codereview.appspot.com/125720043)
* [Go 1.4+ garbage collector YC thread](https://news.ycombinator.com/item?id=8148666)
* [Plans for pauseless GC algorithm](https://groups.google.com/forum/#!msg/golang-dev/GvA0DaCI2BU/1EpYa8HbxdIJ)
* [Arena allocator in Go](http://blog.tuxychandru.com/2014/07/arena-allocation-in-go.html)
* [Google search for Go arena allocator](https://www.google.co.uk/search?q=golang+gc+arena+allocator&oq=golang+gc+arena+allocator&aqs=chrome..69i57j69i64.4822j0j7&sourceid=chrome&es_sm=91&ie=UTF-8)
* [Google search for Go GC flags](https://www.google.co.uk/search?q=golang+gc+flags&oq=golang+gc+&aqs=chrome.4.69i57j0l5.8358j0j7&sourceid=chrome&es_sm=91&ie=UTF-8#q=golang+gc+flags&start=10&tbs=qdr:m)
* [How To Determine Web Application Thread Pool Size](http://venkateshcm.com/2014/05/How-To-Determine-Web-Applications-Thread-Poll-Size/)
* [Go's unsafe.Pointer Pointer Type](http://learngowith.me/gos-pointer-pointer-type/)
* [perf Counting](http://www.brendangregg.com/blog/2014-07-03/perf-counting.html)


FlameGraph links:
* [Brendan Gregg's link page](http://www.brendangregg.com/flamegraphs.html)
* [LISA'13](http://www.brendangregg.com/Slides/LISA13_Flame_Graphs.pdf)
* [FlameGraph github repo, by Brendan Gregg](https://github.com/brendangregg/FlameGraph)
* [Java Flame Graphs](http://www.brendangregg.com/blog/2014-06-12/java-flame-graphs.html)
* [CPU Flame Graphs](http://www.brendangregg.com/FlameGraphs/cpuflamegraphs.html)
* [off-CPU Flame Graphs](http://agentzh.org/misc/slides/off-cpu-flame-graphs.pdf)
* [StackGraph command](http://godoc.org/code.google.com/p/rog-go/cmd/stackgraph) [see also](https://plus.google.com/+rogerpeppe/posts/XfK6UR57xNK)
* [Go Flame Graphs](https://github.com/kisielk/goflamegraph)


## Author

Lorenzo Alberton

* Web: [http://alberton.info](http://alberton.info)
* Twitter: [@lorenzoalberton](https://twitter.com/lorenzoalberton)
* Linkedin: [/in/lorenzoalberton](https://www.linkedin.com/in/lorenzoalberton)


## Copyright

See [LICENSE](LICENSE) document
