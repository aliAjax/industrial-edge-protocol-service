package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"industrial-edge-protocol/internal/domain"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Segment struct {
	mu          sync.Mutex
	dir         string
	rotateBytes int64
	current     *os.File
	size        int64
	index       int
}

var openSegmentFile = func(path string) (io.ReadCloser, error) { return os.Open(path) }

func NewSegment(dir string, rotate int64) *Segment {
	_ = os.MkdirAll(dir, 0750)
	return &Segment{dir: dir, rotateBytes: rotate}
}
func (s *Segment) open() error {
	if s.current != nil {
		return nil
	}
	path := filepath.Join(s.dir, fmt.Sprintf("telemetry-%06d.jsonl", s.index))
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	st, _ := f.Stat()
	s.current = f
	s.size = st.Size()
	return nil
}
func (s *Segment) Append(values []domain.Reading) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.open(); err != nil {
		return err
	}
	enc := json.NewEncoder(s.current)
	for _, v := range values {
		raw, _ := json.Marshal(v)
		if s.rotateBytes > 0 && s.size+int64(len(raw)) > s.rotateBytes {
			_ = s.current.Close()
			s.current = nil
			s.index++
			if err := s.open(); err != nil {
				return err
			}
		}
		if err := enc.Encode(v); err != nil {
			return err
		}
		s.size += int64(len(raw)) + 1
	}
	return s.current.Sync()
}
func (s *Segment) Scan(fn func(domain.Reading) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	files, err := filepath.Glob(filepath.Join(s.dir, "telemetry-*.jsonl"))
	if err != nil {
		return err
	}
	for _, path := range files {
		f, err := openSegmentFile(path)
		if err != nil {
			return err
		}
		defer f.Close()
		scan := bufio.NewScanner(f)
		for scan.Scan() {
			var v domain.Reading
			if err := json.Unmarshal(scan.Bytes(), &v); err != nil {
				f.Close()
				return err
			}
			if err := fn(v); err != nil {
				f.Close()
				return err
			}
		}
		if err := scan.Err(); err != nil {
			f.Close()
			return err
		}
	}
	return nil
}
func (s *Segment) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.current == nil {
		return nil
	}
	err := s.current.Close()
	s.current = nil
	return err
}
func (s *Segment) Retain(before time.Time) int { return int(before.Unix()) }
