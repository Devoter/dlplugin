# dlplugin

This package is based on the official Go `plugin` package, but modified to use any dynamic C libraries (Only Linux, FreeBSD, and macOS).

It provides a thread-safe interface for loading/unloading dynamic libraries, but the library symbols should be loaded manually using `PluginInitializer`.


## Installation

`go get github.com/Devoter/dlplugin`

or use `go mod` tool.


## Usage

**WARNING:** Windows implementation was not tested and should not be used.

This package uses `cgo`, it is highly recommended to read the [official CGO documentation](https://pkg.go.dev/cmd/cgo).

Open and prepare a plugin via `dlplugin.Open()` function:

```go
Open(path string, initializer PluginInitializer) (*Plugin, error)
```

It accepts a library filename and an initializer. The `Init()` method denotes a function type which initializes a plugin API.

```go
type PluginInitializer interface {
	Init(lookup func(symName string) (uintptr, error)) error
}
```

An opened library may be closed using the `Close()` method of the `Plugin` or the `Close()` function:

```go
func (p *Plugin) Close() error

func Close(p *Plugin) error
```

## Examples

All examples have Makefiles, therefore you can build each example with the `make` command.

### Basic

This example is a program that prints "Hello, world!" via dynamic library call. The example contains two implementations of program: naive and with an interface.

<details>
  <summary>Plugin code</summary>

```go
package main

import "C"
import "fmt"

//export println
func println(v *C.char) {
	s := C.GoString(v)

	fmt.Println(s)
}

func main() {}

```
</details>

<details>
  <summary>Program code</summary>

```go
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
```
</details>

[See](./examples/basic).


### Random values

This example starts a random values generator from the library and reads generated values.

[See](./examples/random_values).


### Callback

This example concatenates two string with a dynamic library and returns the result via a callback function.

[See](./examples/callback).


### Multilib

This example loads two libs with the single interface. The program instanciates remote objects and works with them.

[See](https://github.com/Devoter/dlplugin_multilib_example)

## License

[LICENSE](./LICENSE)
