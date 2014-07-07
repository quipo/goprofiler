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
		Interval: time.Second * 15
	}
	prof := profiler.NewProfiler(pprofConf)
	go prof.Run()
	
	// your app's code here...

	prof.Stop()
}
```

## Author

Lorenzo Alberton

* Web: [http://alberton.info](http://alberton.info)
* Twitter: [@lorenzoalberton](https://twitter.com/lorenzoalberton)
* Linkedin: [/in/lorenzoalberton](https://www.linkedin.com/in/lorenzoalberton)


## Copyright

See [LICENSE](LICENSE) document
