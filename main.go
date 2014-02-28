package main

import (
	"log"
	"time"
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
	ct := make(chan msgArr, *bufferSize)
	cb := make(chan bool)
	f := &producer{consumer: ct}
	t := &consumer{producer: ct, fin: cb}
	startProfile()
	go f.run()
	go t.run()
	<-cb
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

type producer struct {
	consumer chan msgArr
}

func (p *producer) run() {
	lockToThread()
	lim := *msgCount * countMult
	var buf msgArr
	for i := 1; i < lim; i++ {
		m := &buf.a[i%arrSize]
		for j := 0; j < msgSize; j++ {
			m.b[j] = (byte(i) / 2) + 1
		}
		if i%arrSize == arrSize-1 {
			p.consumer <- buf
		}
	}
	close(p.consumer)
}

type consumer struct {
	producer chan msgArr
	fin chan bool
}

func (c *consumer) run() {
	lockToThread()
	var comp msgArr
	s := time.Now()
	count := 0
	for {
		im :=<-c.producer
		if im == comp {
			break
		}
		count++
	}
	total := time.Now().UnixNano() - s.UnixNano()
	println("count = ", count)
	println("Nano ", total)
	println("Micro ", total/1000)
	println("Milli ", total/(1000*1000))
	println("Seconds ", total/(1000*1000*1000))
	c.fin <- true
}
