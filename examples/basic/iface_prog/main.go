// Copyright 2022 Alexey Nosov.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Devoter/dlplugin"
)

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

	var papi PluginAPI = &plugInst

	papi.Println("Hello, world!") // call plugin function
}
