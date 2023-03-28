[![GoTest](https://github.com/nc30/tinyfile/actions/workflows/gotest.yml/badge.svg)](https://github.com/nc30/tinyfile/actions/workflows/gotest.yml)


tinyfile is tiny service library

## install

`go get github.com/nc30/tinyfile`


## PidFile controll

tinyfile support simple pid control.


### example

```go
package main

import (
	"fmt"
	"os"

	"github.com/nc30/tinyfile"
)

var PidFilePath = "/tmp/pid"

func main(){
	err := tinyfile.PidSet(PidFilePath)
	if err != nil {
		fmt.FPrintln(err)
		os.Exit(1)
	}
	defer tinyfile.PidClean()


	/// some process
}

```
