// Copyright 2022 Alexey Nosov.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Devoter/dlplugin"
)

func main() {
	pluginFilename := flag.String("plugin", "", "plugin filename")
	timeout := flag.Int64("timeout", 100, "generator timeout (ms)")
	help := flag.Bool("help", false, "show this text")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return
	}

	var plugAPI PluginAPI

	plug, err := dlplugin.Open(*pluginFilename, &plugAPI)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not initialize a plugin by the reason: %v\n", err)
		os.Exit(1)
	}

	defer plug.Close()

	plugAPI.Start(time.Duration(*timeout) * time.Millisecond)
	defer plugAPI.Stop()

	for i := 0; i < 10; i++ {
		fmt.Printf("value: %d\n", plugAPI.Read())
		<-time.After(300 * time.Millisecond)
	}
}
