package main

import (
	"time"
)

func prodcon() {
	ct := make(chan msgArr, *bufferSize)
	cb := make(chan bool)
	f := &producer{consumer: ct}
	t := &consumer{producer: ct, fin: cb}
	go f.run()
	go t.run()
	<-cb
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
