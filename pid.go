package tinyfile

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrAnotherProcessRunning = errors.New("pid file is aleady exist")
	ErrAleadySetPid          = errors.New("aleady seted pid on this process")

	PidFileFlg                    = os.O_CREATE | os.O_WRONLY
	PidFilePermission os.FileMode = 0664
)

var pidFilepath string = ""

// PidSet create pid file to argument path
// require run PidClean() on process end
func PidSet(path string) error {
	if pidFilepath != "" {
		return ErrAleadySetPid
	}

	if FileExist(path) {
		return ErrAnotherProcessRunning
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

// PidClean delete pid file of created by PidSet()
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

// FileExist check path is file and exist
// true is only file exist and type is file
func FileExist(path string) bool {
	f, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return !f.IsDir()
}
