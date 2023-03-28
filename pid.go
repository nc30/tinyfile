package tinyfile

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrAleadyRunning = errors.New("pid file is aleady exist")
	ErrAleadySetPid  = errors.New("aleady seted pid on this process")

	PidFileFlg                    = os.O_CREATE | os.O_WRONLY
	PidFilePermission os.FileMode = 0664
)

var pidFilepath string = ""

func PidSet(path string) error {
	if pidFilepath != "" {
		return ErrAleadySetPid
	}

	_, e := os.Stat(path)
	if e == nil {
		return ErrAleadyRunning
	}

	pid := os.Getpid()
	f, err := os.OpenFile(path, PidFileFlg, PidFilePermission)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "%d", pid)

	pidFilepath = path

	return nil
}

func PidClean() error {
	if pidFilepath != "" {
		err := os.Remove(pidFilepath)
		if err == nil {
			pidFilepath = ""
		}

		return err
	}

	return nil
}
