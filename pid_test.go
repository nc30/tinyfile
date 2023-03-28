package tinyfile_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/nc30/tinyfile"
)

func TestPidSet(t *testing.T) {
	tmp, err := os.MkdirTemp("", "gotest")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)
	defer tinyfile.PidClean()

	filename := tmp + "/test.pid"
	err = tinyfile.PidSet(filename)
	if err != nil {
		t.Error(err)
	}

	b, err := os.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}

	if string(b) != strconv.Itoa(os.Getpid()) {
		t.Errorf("pid is not equal %s != %d", string(b), os.Getpid())
	}

	err = tinyfile.PidSet(filename)
	if err != tinyfile.ErrAleadySetPid {
		t.Errorf("invalid error")
	}
}

func TestPidSetAleadyRunning(t *testing.T) {
	tmp := os.TempDir()
	tmp, err := os.MkdirTemp("", "gotest")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)
	defer tinyfile.PidClean()

	filename := tmp + "/test.pid"
	os.WriteFile(filename, []byte(""), 0664)

	err = tinyfile.PidSet(filename)
	if err != tinyfile.ErrAnotherProcessRunning {
		t.Errorf("invalid error of %v", err)
	}
}

func TestPidSetTwiceSet(t *testing.T) {
	tmp, err := os.MkdirTemp("", "gotest")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)
	defer tinyfile.PidClean()

	filename := tmp + "/test.pid"
	filename2 := tmp + "/test2.pid"

	err = tinyfile.PidSet(filename)
	if err != nil {
		t.Error(err)
	}

	err = tinyfile.PidSet(filename2)
	if err != tinyfile.ErrAleadySetPid {
		t.Errorf("invalid error of %s", err)
	}
}

func TestPidClean(t *testing.T) {
	err := tinyfile.PidClean()
	if err != nil {
		t.Error(err)
	}

	tmp, err := os.MkdirTemp("", "gotest")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	filename := tmp + "/test.pid"
	err = tinyfile.PidSet(filename)
	if err != nil {
		t.Error(err)
	}

	_, err = os.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}

	err = tinyfile.PidClean()
	if err != nil {
		t.Error(err)
	}

	_, err = os.ReadFile(filename)
	if err == nil {
		t.Error(err)
	}
}

func TestFileExist(t *testing.T) {
	f, err := os.CreateTemp(os.TempDir(), "gotest")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())

	if !tinyfile.FileExist(f.Name()) {
		t.Errorf("file is exist")
	}

	os.Remove(f.Name())

	if tinyfile.FileExist(f.Name()) {
		t.Errorf("file is not exist")
	}
}

func TestFileExistDirectory(t *testing.T) {
	tmp, err := os.MkdirTemp("", "gotest")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	if tinyfile.FileExist(tmp) {
		t.Errorf("this is directory")
	}
}
