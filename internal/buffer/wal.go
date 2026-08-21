package buffer

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"industrial-edge-protocol/internal/domain"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Record struct {
	Sequence  uint64         `json:"sequence"`
	Reading   domain.Reading `json:"reading"`
	Digest    string         `json:"digest"`
	CreatedAt time.Time      `json:"created_at"`
}
type WAL struct {
	mu       sync.Mutex
	dir      string
	maxBytes int64
	sequence uint64
}

var openWALFile = func(path string) (io.ReadCloser, error) { return os.Open(path) }

func New(dir string, maxBytes int64) *WAL {
	_ = os.MkdirAll(dir, 0750)
	return &WAL{dir: dir, maxBytes: maxBytes}
}
func (w *WAL) Append(readings []domain.Reading) ([]Record, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	path := filepath.Join(w.dir, "segment-current.jsonl")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	st, _ := f.Stat()
	if st.Size() > w.maxBytes {
		return nil, fmt.Errorf("buffer quota exceeded")
	}
	records := make([]Record, 0, len(readings))
	enc := json.NewEncoder(f)
	for _, r := range readings {
		w.sequence++
		raw, _ := json.Marshal(r)
		sum := sha256.Sum256(raw)
		rec := Record{Sequence: w.sequence, Reading: r, Digest: fmt.Sprintf("%x", sum), CreatedAt: time.Now().UTC()}
		if err := enc.Encode(rec); err != nil {
			return records, err
		}
		records = append(records, rec)
	}
	return records, nil
}
func (w *WAL) Replay(fn func(Record) error) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	files, err := filepath.Glob(filepath.Join(w.dir, "segment-*.jsonl"))
	if err != nil {
		return err
	}
	for _, name := range files {
		f, err := openWALFile(name)
		if err != nil {
			return err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			var rec Record
			if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
				f.Close()
				return err
			}
			if err := fn(rec); err != nil {
				f.Close()
				return err
			}
		}
		if err := scanner.Err(); err != nil {
			f.Close()
			return err
		}
	}
	return nil
}
func (w *WAL) Compact(upto uint64) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	path := filepath.Join(w.dir, "segment-current.jsonl")
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		in.Close()
		return err
	}
	scanner := bufio.NewScanner(in)
	enc := json.NewEncoder(out)
	for scanner.Scan() {
		var rec Record
		if json.Unmarshal(scanner.Bytes(), &rec) == nil && rec.Sequence > upto {
			_ = enc.Encode(rec)
		}
	}
	in.Close()
	out.Close()
	return os.Rename(tmp, path)
}
