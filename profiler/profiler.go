package profiler

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// Config wraps the profiler settings
type Config struct {
	CPU               bool   `json:"cpu"`
	Memory            bool   `json:"memory"`
	Block             bool   `json:"block"`
	Goroutine         bool   `json:"goroutine"`
	Prefix            string `json:"prefix"`
	Interval          string `json:"interval"`
	MemoryProfileRate int    `json:"memory_profile_rate"`
	CPUProfileRate    int    `json:"cpu_profile_rate"`
}

type profiler struct {
	conf        Config
	terminateCh chan struct{}
	closers     []func()
}

// NewProfiler initialises a new instance of a profiler
func NewProfiler(conf Config) *profiler {
	return &profiler{
		conf:        conf,
		terminateCh: make(chan struct{}, 0),
		closers:     make([]func(), 0),
	}
}

func (c Config) isOn() bool {
	return c.CPU || c.Memory || c.Goroutine || c.Block
}

// Run starts the profiler
func (p profiler) Run() {
	if p.conf.CPU {
		if p.conf.CPUProfileRate > 0 {
			runtime.SetCPUProfileRate(p.conf.CPUProfileRate)
		}
		p.startProfilingCPU()
	}
	if p.conf.Memory {
		runtime.MemProfileRate = p.conf.MemoryProfileRate
	}
	if p.conf.Block {
		runtime.SetBlockProfileRate(1)
	}

	if p.conf.isOn() && "" != p.conf.Interval {
		interval, err := time.ParseDuration(p.conf.Interval)
		if nil != err {
			fmt.Println(err)
			return
		}
		timer := time.NewTimer(interval)
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

// TakeSnapshot takes a profiling data snapshot for the enabled resources
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

// Stop terminates the active profiler(s)
func (p profiler) Stop() {
	close(p.terminateCh)
}

// opens a new output file to collect CPU profiling information
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
