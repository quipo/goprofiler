package profiler

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

// Config wraps the profiler settings
type Config struct {
	CPU                  bool   `json:"cpu"`                    // enable CPU profiling
	Memory               bool   `json:"memory"`                 // enable memory profiling
	Block                bool   `json:"block"`                  // enable block profiling
	Goroutine            bool   `json:"goroutine"`              // enable goroutine profiling
	Mutex                bool   `json:"mutex"`                  // enable contended mutex profiling
	Prefix               string `json:"prefix"`                 // prefix for the name of the file storing the profiling snapshots
	Interval             string `json:"interval"`               // interval between 2 subsequent snapshots
	MemoryProfileRate    int    `json:"memory_profile_rate"`    // set to 1 to include every allocated block in the profile, 0 to disable
	CPUProfileRate       int    `json:"cpu_profile_rate"`       // set to a value above zero to enable collection (hz samples per second)
	MutexProfileFraction int    `json:"mutex_profile_fraction"` // set to a value above zero to enable collection
}

// profiler is unexported to force initialisation via constructor
type profiler struct {
	conf        Config
	terminateCh chan struct{}
	closers     []func()
	logger      *log.Logger
}

// NewProfiler initialises a new instance of a profiler
func NewProfiler(conf Config) *profiler {
	return &profiler{
		conf:        conf,
		terminateCh: make(chan struct{}),
		closers:     make([]func(), 0),
		logger:      log.New(os.Stdout, "[profiler] ", log.Ldate|log.Ltime),
	}
}

func (c Config) isOn() bool {
	return c.CPU || c.Memory || c.Goroutine || c.Block || c.Mutex
}

// Run starts the profiler
func (p *profiler) Run() {
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
	if p.conf.Mutex {
		runtime.SetMutexProfileFraction(p.conf.MutexProfileFraction)
	}

	if p.conf.isOn() && "" != p.conf.Interval {
		interval, err := time.ParseDuration(p.conf.Interval)
		if nil != err {
			log.Println("Error parsing interval parameter:", err)
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
func (p *profiler) TakeSnapshot() {
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
	if p.conf.Mutex {
		p.takeMutexSnapshot()
	}

	for _, c := range p.closers {
		c()
	}
	p.closers = p.closers[:0]
}

// Stop terminates the active profiler(s)
func (p profiler) Stop() {
	close(p.terminateCh)
}

// opens a new output file to collect CPU profiling information
func (p *profiler) startProfilingCPU() {
	pprofFile := fmt.Sprintf("%scpu.%d.pprof", p.conf.Prefix, time.Now().Unix())
	p.logger.Println("Starting new CPU Profiler:", pprofFile)
	f, err := os.Create(pprofFile)
	if err != nil {
		panic(err)
	}
	if err = pprof.StartCPUProfile(f); err != nil {
		p.logger.Println("could not start CPU profile: ", err)
	}
	p.closers = append(p.closers, func() {
		if err = f.Close(); err != nil {
			p.logger.Println(err)
		}
	})
}

// collect CPU profiling information
func (p profiler) takeCPUSnapshot() {
	p.logger.Println("Stopping CPU Profiler")
	pprof.StopCPUProfile()
}

// collect Memory profiling information
func (p profiler) takeMemorySnapshot() {
	pprofFile := fmt.Sprintf("%smem.%d.pprof", p.conf.Prefix, time.Now().Unix())
	p.logger.Println("Taking Memory Profile Snapshot:", pprofFile)
	f, err := os.Create(pprofFile)
	if err != nil {
		p.logger.Println(err)
	}

	if err = pprof.WriteHeapProfile(f); err != nil {
		p.logger.Println(err)
	}
	if err = f.Close(); err != nil {
		p.logger.Println(err)
	}
}

// collect Block profiling information
func (p profiler) takeBlockSnapshot() {
	pprofFile := fmt.Sprintf("%sblock.%d.pprof", p.conf.Prefix, time.Now().Unix())
	p.logger.Println("Taking Block Profile Snapshot:", pprofFile)
	f, err := os.Create(pprofFile)
	if err != nil {
		p.logger.Println(err)
	}
	profile := pprof.Lookup("block")
	if err = profile.WriteTo(f, 2); err != nil {
		p.logger.Println(err)
	}

	if err = f.Close(); err != nil {
		p.logger.Println(err)
	}
	runtime.SetBlockProfileRate(0)
}

// collect Goroutine profiling information
func (p profiler) takeGoroutineSnapshot() {
	pprofFile := fmt.Sprintf("%sgoroutine.%d.pprof", p.conf.Prefix, time.Now().Unix())
	p.logger.Println("Taking Goroutine Profile Snapshot:", pprofFile)
	f, err := os.Create(pprofFile)
	if err != nil {
		p.logger.Println(err)
	}
	profile := pprof.Lookup("goroutine")
	if err = profile.WriteTo(f, 2); err != nil {
		p.logger.Println(err)
	}
	if err = f.Close(); err != nil {
		p.logger.Println(err)
	}
}

// collect Mutex profiling information
func (p profiler) takeMutexSnapshot() {
	pprofFile := fmt.Sprintf("%smutex.%d.pprof", p.conf.Prefix, time.Now().Unix())
	p.logger.Println("Taking Mutex Profile Snapshot:", pprofFile)
	f, err := os.Create(pprofFile)
	if err != nil {
		p.logger.Println(err)
	}
	profile := pprof.Lookup("mutex")
	if err = profile.WriteTo(f, 2); err != nil {
		p.logger.Println(err)
	}
	if err = f.Close(); err != nil {
		p.logger.Println(err)
	}
}
