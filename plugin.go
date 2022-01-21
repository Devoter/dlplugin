// Copyright 2022 Alexey Nosov.
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dlplugin implements loading and symbol resolution of C plugins.
//
// A plugin is a shared object (dynamic library) with exported functions and variables.
//
// Currently plugins are only supported on Linux, FreeBSD, and macOS.
//
// See examples in examples directory.
package dlplugin

// PluginInitializer provides an interface which is used to initialize a plugin interface.
type PluginInitializer interface {
	// Init initializes a plugin API.
	Init(lookup func(symName string) (uintptr, error)) error
}

// Plugin is a loaded plugin.
type Plugin struct {
	handler  uintptr
	filepath string
	name     string
}

// Close closes a plugin.
func (p *Plugin) Close() error {
	return libClose(p)
}

// Init provides a lookup function and calls an initializer's Init method.
func (p *Plugin) Init(initializer PluginInitializer) error {
	return libInit(p, initializer)
}

// Open opens a plugin.
// If a path has already been opened, then the existing *Plugin is returned.
// It is safe for concurrent use by multiple goroutines.
// If the initializer is nil the initialization process will be skipped.
func Open(path string, initializer PluginInitializer) (*Plugin, error) {
	return libOpen(path, initializer)
}

// Close closes a plugin.
func Close(p *Plugin) error {
	return libClose(p)
}
