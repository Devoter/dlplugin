// Copyright 2022 Alexey Nosov.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

/*
#include <stdint.h>
#include <stdlib.h>

extern void ConcatStringsCallback(uintptr_t h, char *result);

static void concat_strings(uintptr_t r, uintptr_t cb_id, char *s1, char *s2)
{
	typedef void (*concat_strings_t)(char *, char *, uintptr_t, uintptr_t);
	typedef void (*callback_t)(uintptr_t, char *);

	((concat_strings_t)r)(s1, s2, cb_id, (uintptr_t)ConcatStringsCallback);
}
*/
import "C"
import (
	"flag"
	"fmt"
	"os"
	"runtime/cgo"
	"unsafe"

	"github.com/Devoter/dlplugin"
)

type Plugin struct {
	concatStrings func(s1 string, s2 string) string
}

func (p *Plugin) ConcatStrings(s1 string, s2 string) string {
	return p.concatStrings(s1, s2)
}

func (p *Plugin) Init(lookup func(symName string) (uintptr, error)) error {
	concatStringsPtr, err := lookup("concat_strings")
	if err != nil {
		return err
	}

	p.concatStrings = func(s1, s2 string) string {
		cs1 := C.CString(s1)
		cs2 := C.CString(s2)

		var result string

		cb := func(s string) { result = s }
		h := cgo.NewHandle(cb)

		C.concat_strings(C.uintptr_t(concatStringsPtr), C.uintptr_t(h), cs1, cs2)

		h.Delete()
		C.free(unsafe.Pointer(cs1))
		C.free(unsafe.Pointer(cs2))

		return result
	}

	return nil
}

//export ConcatStringsCallback
func ConcatStringsCallback(h C.uintptr_t, result *C.char) {
	cb := cgo.Handle(h).Value().(func(s string))

	cb(C.GoString(result))
}

func main() {
	pluginFilename := flag.String("plugin", "", "plugin filename")
	help := flag.Bool("help", false, "show this text")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	var plugInst Plugin

	// loading the plugin an initializing its inteface.
	plug, err := dlplugin.Open(*pluginFilename, &plugInst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not initialize a plugin by the reason: %v\n", err)
		os.Exit(1)
	}

	defer plug.Close() // release a plugin library

	result := plugInst.ConcatStrings("Hello, ", "world!") // call plugin function

	fmt.Println(result)
}
