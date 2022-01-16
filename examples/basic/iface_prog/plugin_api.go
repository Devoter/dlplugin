// Copyright 2022 Alexey Nosov.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
#include <stdint.h>
#include <stdlib.h>

static void println(uintptr_t r, char *s)
{
	((void (*)(char *))r)(s);
}
*/
import "C"
import "unsafe"

// PluginAPI declares a plugin interface.
type PluginAPI interface {
	Println(s string)
}

// Plugin struct declares a plugin interface implementation.
type Plugin struct {
	println func(s string)
}

// Println calls a plugin `println` function transparently.
func (p *Plugin) Println(s string) {
	p.println(s)
}

// InitPluginAPIFactory returns an interface instance and a function that initializes the plugin API.
func (p *Plugin) Init(lookup func(symName string) (uintptr, error)) error {
	printlnPtr, err := lookup("println") // search for the library symbol
	if err != nil {
		return err
	}

	// defining a Go wrapper.
	p.println = func(s string) {
		cs := C.CString(s)
		defer C.free(unsafe.Pointer(cs))

		C.println(C.uintptr_t(printlnPtr), cs)
	}

	return nil
}
