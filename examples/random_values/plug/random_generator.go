// Copyright 2022 Alexey Nosov.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	if err := binary.Read(crand.Reader, binary.BigEndian, &v); err != nil {
		panic(err)
	}

	return v
}

// RandomGenerator provides an asynchronous random numbers generator.
type RandomGenerator struct {
	mx      sync.RWMutex
	wg      sync.WaitGroup
	stopCh  chan struct{}
	started bool
	random  *rand.Rand
	value   int64
}

func NewRandomGenerator() *RandomGenerator {
	var src cryptoSource

	rnd := rand.New(src)
	rnd.Seed(time.Now().UnixNano())
	rnd.Uint64() // may be panic
	msg := recover()

	if msg != "" && msg != nil {
		log.Println("could not use /dev/urandom", msg)
		psrc := rand.NewSource(time.Now().UnixNano())
		rnd = rand.New(psrc)
		rnd.Seed(time.Now().UnixNano())
	}

	return &RandomGenerator{
		stopCh: make(chan struct{}, 1),
		random: rnd,
	}
}

func (rg *RandomGenerator) Start(timeout time.Duration) {
	rg.mx.RLock()

	if rg.started {
		rg.mx.RUnlock()

		return
	}

	rg.mx.RUnlock()

	started := make(chan struct{}, 1)

	rg.wg.Add(1)
	go rg.start(started, timeout)
	<-started
}

func (rg *RandomGenerator) Stop() {
	rg.mx.Lock()
	defer rg.mx.Unlock()

	if !rg.started {
		return
	}

	rg.stopCh <- struct{}{}
	rg.wg.Wait()
	rg.started = false
}

func (rg *RandomGenerator) Read() int64 {
	return atomic.LoadInt64(&rg.value)
}

func (rg *RandomGenerator) start(started chan<- struct{}, timeout time.Duration) {
	defer rg.wg.Done()

	rg.mx.Lock()
	rg.started = true
	atomic.StoreInt64(&rg.value, rg.random.Int63())
	rg.mx.Unlock()

	started <- struct{}{}

	for {
		select {
		case <-rg.stopCh:
			return
		case <-time.After(timeout):
			atomic.StoreInt64(&rg.value, rg.random.Int63())
		}
	}
}
