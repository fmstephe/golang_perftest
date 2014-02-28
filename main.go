package main

import (
	"log"
	"flag"
	"runtime"
	"runtime/pprof"
	"os"
)

var msgCount = flag.Int("m", 1, "The number of messages (in millions) to send")
var bufferSize = flag.Int("b", 100, "The size of the buffer for the communicating channel")
var notLockThread = flag.Bool("nl", false, "Indicates that we don't want to lock messaging goroutines to an OS thread")

const countMult = 1000 * 1000

const arrSize = 10
const msgSize = 64

type msgArr struct {
	a [arrSize]msg
}

type msg struct {
	b [msgSize]byte
}

func main() {
	runtime.GOMAXPROCS(4)
	flag.Parse()
	startProfile()
	prodcon()
	endProfile()
}

func startProfile() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
}

func endProfile() {
	pprof.StopCPUProfile()
}

func lockToThread() {
	if *notLockThread {
		println("Not locking threads")
		return
	}
	runtime.LockOSThread()
}
