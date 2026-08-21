package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"industrial-edge-protocol/internal/domain"
)

func TestReplayClosesEachSegmentPromptly(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 4; i++ {
		path := filepath.Join(dir, fmt.Sprintf("telemetry-%06d.jsonl", i))
		if err := os.WriteFile(path, []byte(`{"point_id":"p1","value":1}`+"\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	originalOpen := openSegmentFile
	active := 0
	overlapped := false
	openSegmentFile = func(path string) (io.ReadCloser, error) {
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		if active != 0 {
			overlapped = true
		}
		active++
		return &trackedSegmentFile{ReadCloser: file, active: &active}, nil
	}
	defer func() { openSegmentFile = originalOpen }()
	segment := NewSegment(dir, 1024)
	if err := segment.Scan(func(domain.Reading) error { return nil }); err != nil {
		t.Fatal(err)
	}
	if overlapped || active != 0 {
		t.Fatalf("segment files overlapped during scan: overlap=%v active=%d", overlapped, active)
	}
}

type trackedSegmentFile struct {
	io.ReadCloser
	active *int
}

func (f *trackedSegmentFile) Close() error {
	*f.active--
	return f.ReadCloser.Close()
}
