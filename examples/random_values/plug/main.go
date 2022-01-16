// Copyright 2022 Alexey Nosov.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// #include <stdint.h>
import "C"
import "time"

var generator RandomGenerator = *NewRandomGenerator()

//export start
func start(timeout C.int64_t) {
	generator.Start(time.Duration(timeout))
}

//export stop
func stop() {
	generator.Stop()
}

//export read
func read() C.int64_t {
	return C.int64_t(generator.Read())
}

func main() {}
