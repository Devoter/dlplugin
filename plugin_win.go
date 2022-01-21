// Copyright 2022 Alexey Nosov.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows && cgo
// +build windows,cgo

package dlplugin

/*
#include <stdint.h>
#include <windows.h>

static uintptr_t dlplugin_open(const char *path, int32_t *errcode)
{
	HMODULE h = LoadLibrary(path);

	if (!h) *errcode = GetLastError();

	return (uintptr_t)h;
}

static uintptr_t dlplugin_lookup(uintptr_t h, const char *name, int32_t *errcode)
{
	void *r = GetProcAddress((HMODULE)h, name);

	if (!r) *errcode = GetLastError();

	return (uintptr_t)r;
}

static int32_t dlplugin_close(uintptr_t h)
{
	return FreeLibrary((HMODULE)h) ? GetLastError() : 0;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

func formatWinErr(errCode C.int32_t) error {
	return fmt.Errorf("windows error code: %d", int32(errCode))
}

func libOpen(name string, initializer PluginInitializer) (*Plugin, error) {
	cPath := make([]byte, C.PATH_MAX+1)
	cRelName := make([]byte, len(name)+1)
	copy(cRelName, name)

	cPathLen := C.GetFullPathName(C.LPCSTR(unsafe.Pointer(&cRelName[0])), C.PATH_MAX, C.LPSTR(unsafe.Pointer(&cPath[0])),
		(*C.LPSTR)(unsafe.Pointer(uintptr(0))))

	if cPathLen == 0 {
		return nil, formatWinErr(C.int32_t(C.GetLastError()))
	} else if cPathLen > C.PATH_MAX {
		return nil, fmt.Errorf("the buffer could not contain a full file path, current size: %d, required: %d",
			int32(C.PATH_MAX), int32(cPathLen))
	}

	filepath := C.GoString((*C.char)(unsafe.Pointer(&cPath[0])))

	pluginsMu.Lock()

	var p *Plugin

	if p = plugins[filepath]; p != nil {
		pluginsMu.Unlock()

		return p, nil
	}

	var cErr C.int32_t

	h := C.dlplugin_open((*C.char)(unsafe.Pointer(&cPath[0])), &cErr)

	if h == 0 {
		pluginsMu.Unlock()

		return nil, errors.New(`dlplugin.Open("` + name + `"): ` + formatWinErr(cErr).Error())
	}

	if len(name) > 3 && name[len(name)-3:] == ".so" {
		name = name[:len(name)-3]
	}

	if plugins == nil {
		plugins = make(map[string]*Plugin)
	}

	p = &Plugin{handler: uintptr(h), filepath: filepath}

	if initializer != nil {
		lookupFunc := produceLookup(h, name)

		if err := initializer.Init(lookupFunc); err != nil {
			C.dlplugin_close(C.uintptr_t(p.handler))
			pluginsMu.Unlock()

			return nil, err
		}
	}

	plugins[filepath] = p
	pluginsMu.Unlock()

	return p, nil
}

func libClose(p *Plugin) error {
	pluginsMu.Lock()

	cErr := C.dlplugin_close(C.uintptr_t(p.handler))
	delete(plugins, p.filepath)

	pluginsMu.Unlock()

	if cErr != 0 {
		return formatWinErr(cErr)
	}

	return nil
}

func produceLookup(h C.uintptr_t, name string) func(symName string) (uintptr, error) {
	return func(symName string) (uintptr, error) {
		var cErr *C.char
		cname := make([]byte, len(symName)+1)
		copy(cname, symName)

		p := C.dlplugin_lookup(h, (*C.char)(unsafe.Pointer(&cname[0])), &cErr)
		if p == 0 {
			return 0, errors.New(`dlplugin.Open("` + name + `"): could not find symbol ` + symName + `: ` + formatWinErr(cErr).Error())
		}

		return uintptr(p), nil
	}
}

func libInit(p *Plugin, initializer PluginInitializer) error {
	pluginsMu.RLock()

	lookup := produceLookup(C.uintptr_t(p.handler), p.name)
	err := initializer.Init(lookup)

	pluginsMu.RUnlock()

	return err
}

var (
	pluginsMu sync.RWMutex
	plugins   map[string]*Plugin
)
