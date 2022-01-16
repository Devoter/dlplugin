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
import (
	"flag"
	"fmt"
	"os"
	"unsafe"

	"github.com/Devoter/dlplugin"
)

type PluginAPI struct {
	Println func(s string)
}

// Init initializes the plugin API.
func (papi *PluginAPI) Init(lookup func(symName string) (uintptr, error)) error {
	printlnPtr, err := lookup("println")
	if err != nil {
		return err
	}

	papi.Println = func(s string) {
		cs := C.CString(s)
		defer C.free(unsafe.Pointer(cs))

		C.println(C.uintptr_t(printlnPtr), cs)
	}

	return nil
}

func main() {
	pluginFilename := flag.String("plugin", "", "plugin filename")
	help := flag.Bool("help", false, "show this text")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	var papi PluginAPI

	plug, err := dlplugin.Open(*pluginFilename, &papi)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not initialize a plugin by the reason: %v\n", err)
		os.Exit(1)
	}

	defer plug.Close()

	papi.Println("Hello, world!")
}
