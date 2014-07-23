package profiler

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type Config struct {
	CPU       bool   `json:"cpu"`
	Memory    bool   `json:"memory"`
	Block     bool   `json:"block"`
	Goroutine bool   `json:"goroutine"`
	Prefix    string `json:"prefix"`
	Interval  int    `json:"interval"`
}

type profiler struct {
	conf        Config
	terminateCh chan struct{}
	closers     []func()
}

func NewProfiler(conf Config) *profiler {
	return &profiler{
		conf:        conf,
		terminateCh: make(chan struct{}, 0),
		closers:     make([]func(), 0),
	}
}

func (p profiler) Run() {

	if p.conf.CPU {
		p.startProfilingCPU()
	}
	if p.conf.Block {
		runtime.SetBlockProfileRate(1)
	}

	if p.conf.Interval > 0 {
		timer := time.NewTimer(p.conf.Interval * time.Millisecond)
		select {
		case <-timer.C:
			p.TakeSnapshot()
			// start again
			p.Run()
		case <-p.terminateCh:
			p.TakeSnapshot()
		}
	}
}

func (p profiler) TakeSnapshot() {
	if p.conf.CPU {
		p.takeCPUSnapshot()
	}
	if p.conf.Memory {
		p.takeMemorySnapshot()
	}
	if p.conf.Block {
		p.takeBlockSnapshot()
	}
	if p.conf.Goroutine {
		p.takeGoroutineSnapshot()
	}

	for _, c := range p.closers {
		c()
	}
	p.closers = nil
}

func (p profiler) Stop() {
	close(p.terminateCh)
}

func (p profiler) startProfilingCPU() {
	pprofFile := fmt.Sprintf("%scpu.%d.pprof", p.conf.Prefix, time.Now().Unix())
	fmt.Println("Starting new CPU Profiler:", pprofFile)
	f, err := os.Create(pprofFile)
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	p.closers = append(p.closers, func() {
		f.Close()
	})
}

// collect CPU profiling information
func (p profiler) takeCPUSnapshot() {
	fmt.Println("Stopping CPU Profiler")
	pprof.StopCPUProfile()
}

// collect Memory profiling information
func (p profiler) takeMemorySnapshot() {
	pprofFile := fmt.Sprintf("%smem.%d.pprof", p.conf.Prefix, time.Now().Unix())
	fmt.Println("Taking Memory Profile Snapshot:", pprofFile)
	f, err := os.Create(pprofFile)
	if err != nil {
		fmt.Println(err)
	}
	pprof.WriteHeapProfile(f)
	f.Close()
}

// collect Block profiling information
func (p profiler) takeBlockSnapshot() {
	pprofFile := fmt.Sprintf("%sblock.%d.pprof", p.conf.Prefix, time.Now().Unix())
	fmt.Println("Taking Block Profile Snapshot:", pprofFile)
	f, err := os.Create(pprofFile)
	if err != nil {
		fmt.Println(err)
	}
	profile := pprof.Lookup("block")
	profile.WriteTo(f, 2)
	f.Close()
	runtime.SetBlockProfileRate(0)
}

// collect Goroutine profiling information
func (p profiler) takeGoroutineSnapshot() {
	pprofFile := fmt.Sprintf("%sgoroutine.%d.pprof", p.conf.Prefix, time.Now().Unix())
	fmt.Println("Taking Goroutine Profile Snapshot:", pprofFile)
	f, err := os.Create(pprofFile)
	if err != nil {
		fmt.Println(err)
	}
	profile := pprof.Lookup("goroutine")
	profile.WriteTo(f, 2)
	f.Close()
}
