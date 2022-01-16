// Copyright 2022 Alexey Nosov.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "C"
import "fmt"

//export println
func println(v *C.char) {
	s := C.GoString(v)

	fmt.Println(s)
}

func main() {}
