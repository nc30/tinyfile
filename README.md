[![Go Reference](https://pkg.go.dev/badge/github.com/nc30/tinyfile.svg)](https://pkg.go.dev/github.com/nc30/tinyfile) [![GoTest](https://github.com/nc30/tinyfile/actions/workflows/gotest.yml/badge.svg)](https://github.com/nc30/tinyfile/actions/workflows/gotest.yml)

tinyfile is tiny service library

## install

`go get github.com/nc30/tinyfile`


## PidFile controll

`tinyfile` support simple pid control and logrotatable file object


### example

```go
package main

import (
	"fmt"
	"os"

	"github.com/nc30/tinyfile"
)

var PidFilePath = "/tmp/test.pid"

func main(){
	err := tinyfile.PidSet(PidFilePath)
	if err != nil {
		fmt.FPrintln(os.Stderr, err)
		os.Exit(1)
	}
	defer tinyfile.PidClean()


	/*****

		some process

	*****/
}

```

## Rotatable file

create `io.WriteCloser` object at reopen of SIGHUP

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/nc30/tinyfile"
)

var logPath = "/tmp/test.log"

func main() {
	// get rotatable io.WriteCloser object
	f, err := NewSync(logPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// close file handler
	defer tinyfile.Close()

	// set outpu to log
	log.SetOutput(f)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// start watch signal of SIGHUP
	tinyfile.Watch(ctx)

	/*****

		some process

	*****/
}
```
