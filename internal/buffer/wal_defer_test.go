package buffer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestReplayPreservesUploadError(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 4; i++ {
		path := filepath.Join(dir, fmt.Sprintf("segment-%06d.jsonl", i))
		if err := os.WriteFile(path, []byte(`{"sequence":1,"reading":{"point_id":"p1"}}`+"\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	originalOpen := openWALFile
	active := 0
	overlapped := false
	openWALFile = func(path string) (io.ReadCloser, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		if active != 0 {
			overlapped = true
		}
		active++
		return &trackedWALFile{ReadCloser: file, active: &active}, nil
	}
	defer func() { openWALFile = originalOpen }()
	wal := New(dir, 1<<20)
	if err := wal.Replay(func(Record) error { return nil }); err != nil {
		t.Fatal(err)
	}
	if overlapped || active != 0 {
		t.Fatalf("wal files stayed open across replay: overlap=%v active=%d", overlapped, active)
	}
}

type trackedWALFile struct {
	io.ReadCloser
	active *int
}

func (f *trackedWALFile) Close() error {
	*f.active--
	return f.ReadCloser.Close()
}
