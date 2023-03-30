package tinyfile

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type OpenFileFunc func(string, int, os.FileMode) (*os.File, error)

var OpenFile OpenFileFunc = os.OpenFile

var (
	FileFlg                    = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	FilePermission os.FileMode = 0664
)

// TinyWriter is zapcore.WriteSyncer
type TinyWriter struct {
	mu sync.Mutex
	fp *os.File
}

func (s *TinyWriter) ReOpen() error {
	if s.fp == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.fp.Close()
	file, err := OpenFile(s.fp.Name(), FileFlg, FilePermission)
	if err != nil {
		return err
	}

	s.fp = file
	return nil
}

func (s *TinyWriter) Close() error {
	if s.fp == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.fp.Close()
}

// write to file
func (s *TinyWriter) Write(p []byte) (n int, err error) {
	if s.fp == nil {
		return 0, errors.New("file not found")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.fp.Write(p)
}

func (s *TinyWriter) Sync() error {
	if s.fp == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.fp.Sync()
}

var defaultRepo = NewRepo()

func NewRepo() *SyncRepo {
	return &SyncRepo{
		files: []*TinyWriter{},
	}
}

// create sync object
func NewSync(path string) (*TinyWriter, error) {
	return defaultRepo.New(path)
}

// watch SIGHUP
func Watch(ctx context.Context) {
	defaultRepo.WatchWithSignal(ctx, syscall.SIGHUP)
}

// watch signal
func WatchWithSignal(ctx context.Context, sig ...os.Signal) {
	defaultRepo.WatchWithSignal(ctx, sig...)
}

// Close is close all files
func Close() {
	defaultRepo.Close()
}

// SyncRepo is TinyFile struct controller.
// that wach some signal and ma
type SyncRepo struct {
	mu    sync.Mutex
	files []*TinyWriter
}

// create SyncRepository
func (sr *SyncRepo) New(path string) (*TinyWriter, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	file, err := OpenFile(path, FileFlg, FilePermission)
	if err != nil {
		return nil, err
	}

	snk := &TinyWriter{
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

// ReOpen all management under TinyWriter
func (sr *SyncRepo) ReOpen() []error {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	ers := []error{}

	for _, snk := range sr.files {
		err := snk.ReOpen()
		if err != nil {
			ers = append(ers, err)
		}
	}

	if len(ers) == 0 {
		return nil
	}

	return ers
}

// Close is close all management under TinyWriter
func (sr *SyncRepo) Close() {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	for _, snk := range sr.files {
		snk.Close()
	}
}
