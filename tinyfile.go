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
	file, err := os.OpenFile(s.fp.Name(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
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

func NewSync(path string) (*Sync, error) {
	return DefaultRepo.New(path)
}

func Watch(ctx context.Context) {
	DefaultRepo.Watch(ctx)
}

func Close() {
	DefaultRepo.Close()
}

type SyncRepo struct {
	mu    sync.Mutex
	files []*Sync
}

func (sr *SyncRepo) New(path string) (*Sync, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0664)
	if err != nil {
		return nil, err
	}

	snk := &Sync{
		fp: file,
	}

	sr.files = append(sr.files, snk)

	return snk, nil
}

func (sr *SyncRepo) Watch(ctx context.Context) {
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)

		for {
			select {
			case <-sig:
				sr.ReOpen()
			case <-ctx.Done():
				goto end
			}
		}
	end:
	}()
}

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

func (sr *SyncRepo) Close() {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	for _, snk := range sr.files {
		snk.Close()
	}
}
