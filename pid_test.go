package tinyfile

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestPidSet(t *testing.T) {
	tmp, err := os.MkdirTemp("", "gotest")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)
	defer PidClean()

	filename := filepath.Join(tmp, "test.pid")
	err = PidSet(filename)
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

	err = PidSet(filename)
	if err != ErrAleadySetPid {
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
	defer PidClean()

	filename := filepath.Join(tmp, "test.pid")
	os.WriteFile(filename, []byte(""), 0664)

	err = PidSet(filename)
	if err != ErrAnotherProcessRunning {
		t.Errorf("invalid error of %v", err)
	}
}

func TestPidSetTwiceSet(t *testing.T) {
	tmp, err := os.MkdirTemp("", "gotest")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)
	defer PidClean()

	filename := filepath.Join(tmp, "test.pid")
	filename2 := filepath.Join(tmp, "test2.pid")

	err = PidSet(filename)
	if err != nil {
		t.Error(err)
	}

	err = PidSet(filename2)
	if err != ErrAleadySetPid {
		t.Errorf("invalid error of %s", err)
	}
}

func TestPidClean(t *testing.T) {
	err := PidClean()
	if err != nil {
		t.Error(err)
	}

	tmp, err := os.MkdirTemp("", "gotest")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	filename := filepath.Join(tmp, "test.pid")
	err = PidSet(filename)
	if err != nil {
		t.Error(err)
	}

	_, err = os.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}

	err = PidClean()
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

	if !FileExist(f.Name()) {
		t.Errorf("file is exist")
	}

	os.Remove(f.Name())

	if FileExist(f.Name()) {
		t.Errorf("file is not exist")
	}
}

func TestFileExistDirectory(t *testing.T) {
	tmp, err := os.MkdirTemp("", "gotest")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmp)

	if FileExist(tmp) {
		t.Errorf("this is directory")
	}
}
