// Copyright 2022 Alexey Nosov.
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && cgo) || (darwin && cgo) || (freebsd && cgo)
// +build linux,cgo darwin,cgo freebsd,cgo

package dlplugin

/*
#cgo linux LDFLAGS: -ldl
#include <dlfcn.h>
#include <limits.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>

static uintptr_t dlplugin_open(const char* path, char** err)
{
	void *h = dlopen(path, RTLD_NOW|RTLD_GLOBAL);

	if (!h) *err = (char *)dlerror();

	return (uintptr_t)h;
}

static uintptr_t dlplugin_lookup(uintptr_t h, const char* name, char** err)
{
	void *r = dlsym((void *)h, name);

	if (!r) *err = (char *)dlerror();

	return (uintptr_t)r;
}

static void dlplugin_close(uintptr_t h)
{
	dlclose((void *)h);
}
*/
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

func libOpen(name string, initializer PluginInitializer) (*Plugin, error) {
	cPath := make([]byte, C.PATH_MAX+1)
	cRelName := make([]byte, len(name)+1)
	copy(cRelName, name)

	if C.realpath(
		(*C.char)(unsafe.Pointer(&cRelName[0])),
		(*C.char)(unsafe.Pointer(&cPath[0]))) == nil {
		return nil, errors.New(`dlplugin.Open("` + name + `"): realpath failed`)
	}

	filepath := C.GoString((*C.char)(unsafe.Pointer(&cPath[0])))

	pluginsMu.Lock()

	var p *Plugin

	if p = plugins[filepath]; p != nil {
		pluginsMu.Unlock()

		return p, nil
	}

	var cErr *C.char
	h := C.dlplugin_open((*C.char)(unsafe.Pointer(&cPath[0])), &cErr)

	if h == 0 {
		pluginsMu.Unlock()

		return nil, errors.New(`dlplugin.Open("` + name + `"): ` + C.GoString(cErr))
	}

	if len(name) > 3 && name[len(name)-3:] == ".so" {
		name = name[:len(name)-3]
	}

	if plugins == nil {
		plugins = make(map[string]*Plugin)
	}

	p = &Plugin{handler: uintptr(h), filepath: filepath}

	lookupFunc := func(symName string) (uintptr, error) {
		cname := make([]byte, len(symName)+1)
		copy(cname, symName)

		p := C.dlplugin_lookup(h, (*C.char)(unsafe.Pointer(&cname[0])), &cErr)
		if p == 0 {
			return 0, errors.New(`dlplugin.Open("` + name + `"): could not find symbol ` + symName + `: ` + C.GoString(cErr))
		}

		return uintptr(p), nil
	}

	if err := initializer.Init(lookupFunc); err != nil {
		C.dlplugin_close(C.uintptr_t(p.handler))
		pluginsMu.Unlock()

		return nil, err
	}

	plugins[filepath] = p
	pluginsMu.Unlock()

	return p, nil
}

func libClose(p *Plugin) error {
	pluginsMu.Lock()

	C.dlplugin_close(C.uintptr_t(p.handler))
	delete(plugins, p.filepath)

	pluginsMu.Unlock()

	return nil
}

var (
	pluginsMu sync.Mutex
	plugins   map[string]*Plugin
)
