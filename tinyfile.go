package tinyfile

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type OpenFileFunc func(string, int, os.FileMode) (*os.File, error)

var OpenFile OpenFileFunc = os.OpenFile

var (
	FileFlg              = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	FilePerm os.FileMode = 0664
)

// Sync is Reopenable io.ReadCloser
type Sync struct {
	mu sync.Mutex
	fp *os.File
}

func (s *Sync) ReOpen() error {
	if s.fp == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.fp.Close()
	file, err := OpenFile(s.fp.Name(), FileFlg, FilePerm)
	if err != nil {
		return err
	}

	s.fp = file
	return nil
}

func (s *Sync) Close() error {
	if s.fp == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.fp.Close()
}

func (s *Sync) Write(p []byte) (n int, err error) {
	if s.fp == nil {
		return 0, errors.New("file not found")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.fp.Write(p)
}

func (s *Sync) Sync() error {
	log.Println("synced")
	return nil
}

var DefaultRepo = NewRepo()

func NewRepo() *SyncRepo {
	return &SyncRepo{
		files: []*Sync{},
	}
}

// create sync object
func NewSync(path string) (*Sync, error) {
	return DefaultRepo.New(path)
}

// watch SIGHUP
func Watch(ctx context.Context) {
	DefaultRepo.WatchWithSignal(ctx, syscall.SIGHUP)
}

// watch signal
func WatchWithSignal(ctx context.Context, sig ...os.Signal) {
	DefaultRepo.WatchWithSignal(ctx, sig...)
}

// Close is close all files
func Close() {
	DefaultRepo.Close()
}

type SyncRepo struct {
	mu    sync.Mutex
	files []*Sync
}

// create SyncRepository
func (sr *SyncRepo) New(path string) (*Sync, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	file, err := OpenFile(path, FileFlg, FilePerm)
	if err != nil {
		return nil, err
	}

	snk := &Sync{
		fp: file,
	}

	sr.files = append(sr.files, snk)

	return snk, nil
}

// Watch start SIGHUP watchng process
func (sr *SyncRepo) Watch(ctx context.Context) {
	sr.WatchWithSignal(ctx, syscall.SIGHUP)
}

// WatchWithSignal start signal watching function
func (sr *SyncRepo) WatchWithSignal(ctx context.Context, sig ...os.Signal) {
	go func() {
		sigchan := make(chan os.Signal, 1)
		defer func() {
			close(sigchan)
		}()
		signal.Notify(sigchan, sig...)

		for {
			select {
			case <-sigchan:
				sr.ReOpen()
			case <-ctx.Done():
				goto end
			}
		}
	end:
	}()
}

// ReOpen refresh i-node of repository files
func (sr *SyncRepo) ReOpen() {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	for _, snk := range sr.files {
		err := snk.ReOpen()
		if err != nil {
			log.Println(err)
		}
	}
}

// Close is close all repository objects
func (sr *SyncRepo) Close() {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	for _, snk := range sr.files {
		snk.Close()
	}
}
