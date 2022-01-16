// Copyright 2022 Alexey Nosov.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
#include <stdint.h>
#include <stdlib.h>

static void run_callback(uintptr_t cb, uintptr_t id, char *result)
{
	typedef void (*callback_t)(uintptr_t r, char *result);

	((callback_t)cb)(id, result);
}
*/
import "C"
import "unsafe"

//export concat_strings
func concat_strings(s1 *C.char, s2 *C.char, cbID C.uintptr_t, callback C.uintptr_t) {
	res := C.CString(C.GoString(s1) + C.GoString(s2))
	C.run_callback(callback, cbID, res)
	C.free(unsafe.Pointer(res))
}

func main() {}
