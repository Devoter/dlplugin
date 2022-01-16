// Copyright 2022 Alexey Nosov.
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (!linux && !freebsd && !darwin && !windows) || !cgo
// +build !linux,!freebsd,!darwin,!windows !cgo

package dlplugin

import "errors"

func libOpen(name string, initializer PluginInitializer) (*Plugin, error) {
	return nil, errors.New("dlplugin: not implemented")
}

func libClose(p *Plugin) error {
	return errors.New("dlplugin: not implemented")
}
