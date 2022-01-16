// Copyright 2022 Alexey Nosov.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
#include <stdint.h>

static void start(uintptr_t r, int64_t timeout)
{
	((void (*)(int64_t))r)(timeout);
}

static void stop(uintptr_t r)
{
	((void (*)())r)();
}

static int64_t read(uintptr_t r)
{
	return ((int64_t (*)())r)();
}
*/
import "C"
import "time"

type PluginAPI struct {
	start func(timeout time.Duration)
	stop  func()
	read  func() int64
}

func (papi *PluginAPI) Start(timeout time.Duration) {
	papi.start(timeout)
}

func (papi *PluginAPI) Stop() {
	papi.stop()
}

func (papi *PluginAPI) Read() int64 {
	return papi.read()
}

func (papi *PluginAPI) Init(lookup func(symName string) (uintptr, error)) error {
	startPtr, err := lookup("start")
	if err != nil {
		return err
	}

	stopPtr, err := lookup("stop")
	if err != nil {
		return err
	}

	readPtr, err := lookup("read")
	if err != nil {
		return err
	}

	papi.start = func(timeout time.Duration) {
		C.start(C.uintptr_t(startPtr), C.int64_t(timeout))
	}

	papi.stop = func() {
		C.stop(C.uintptr_t(stopPtr))
	}

	papi.read = func() int64 {
		return int64(C.read(C.uintptr_t(readPtr)))
	}

	return nil
}
