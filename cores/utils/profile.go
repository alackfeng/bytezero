package utils

import (
	"os"

	//	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"
)

// go build -o bin/bytezero.exe -v -gcflags '-N -l' main.go
// export BZGO_PROF=1; ./bin/bytezero
// export BZGO_PROF=1; GODEBUG=gctrace=1 ./bin/bytezero
// go tool pprof ./bin/bytezero bzgo.cpuprof
// go tool pprof ./bin/bytezero bzgo.memprof
// help: web, top100, peek, list
// go tool pprof -alloc_space /home/vagrant/dcs/dcs/bin/bytezero bzgo.memprof

// web
// http://127.0.0.1:17700/debug/pprof
// http://127.0.0.1:17700/debug/pprof/profile
// http://127.0.0.1:17700/debug/pprof/heap

// go tool pprof -alloc_space -cum -svg http://127.0.0.1:17700/debug/pprof/heap > heap.svg
// go tool pprof -inuse_space -cum -svg http://127.0.0.1:17700/debug/pprof/heap > heap.svg

// EnvEnableProfiling
const (
	EnvEnableProfiling = "BZGO_PROF"
	cpuProfile         = "bzgo.cpuprof"
	heapProfile        = "bzgo.memprof"
)

var stat runtime.MemStats
func printMemStat() {
	// memory.
	runtime.ReadMemStats(&stat)
	log.Printf("Alloc: %v, TotalAlloc: %v, HeapAlloc: %v, NumGC: %v.", stat.Alloc, stat.TotalAlloc, stat.HeapAlloc, stat.NumGC)
}

// startProfiling begins CPU profiling and returns a `stop` function to be
// executed as late as possible. The stop function captures the memprofile.
func startProfiling() (func(), error) {
	// start CPU profiling as early as possible
	ofi, err := os.Create(cpuProfile)
	if err != nil {
		return nil, err
	}
	pprof.StartCPUProfile(ofi)
	go func() {
		for range time.NewTicker(time.Second * 30).C {
			err := writeHeapProfileToFile()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	// web.
   	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:17700", nil))
	}()

	stopProfiling := func() {
		pprof.StopCPUProfile()
		ofi.Close() // captured by the closure
	}
	return stopProfiling, nil
}

func writeHeapProfileToFile() error {
	mprof, err := os.Create(heapProfile)
	if err != nil {
		return err
	}
	// fmt.Println("writeHeapProfileToFile - 30s.")
	MakeRecycle()

	defer mprof.Close() // _after_ writing the heap profile
	return pprof.WriteHeapProfile(mprof)
}

func InitGC(c int) {
	debug.SetGCPercent(c)
}

func MakeRecycle() {
	// runtime.GC()
	printMemStat()
	// debug.FreeOSMemory()
}


// ProfileIfEnabled -
func ProfileIfEnabled() (func(), error) {
	// FIXME this is a temporary hack so profiling of asynchronous operations
	// works as intended.
	printMemStat()

	if os.Getenv(EnvEnableProfiling) != "" {
		log.Println("startProfiling true.")
		InitGC(400)
		stopProfilingFunc, err := startProfiling() // TODO maybe change this to its own option... profiling makes it slower.
		if err != nil {
			return nil, err
		}
		return stopProfilingFunc, nil
	}
	return func() {}, nil
}
