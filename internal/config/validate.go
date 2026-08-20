package config

import (
	"fmt"
	"strings"
)

func (c Config) Validate() error {
	if strings.TrimSpace(c.HTTPAddr) == "" {
		return fmt.Errorf("http address required")
	}
	if c.WorkerInterval <= 0 {
		return fmt.Errorf("worker interval must be positive")
	}
	if c.MaxBufferBytes < 1<<20 {
		return fmt.Errorf("buffer too small")
	}
	if c.CommandTimeout <= 0 {
		return fmt.Errorf("command timeout must be positive")
	}
	return nil
}
func (c Config) SafeSummary() map[string]any {
	return map[string]any{"http_addr": c.HTTPAddr, "data_dir": c.DataDir, "worker_interval": c.WorkerInterval.String(), "max_buffer_bytes": c.MaxBufferBytes, "command_timeout": c.CommandTimeout.String()}
}

func (c Config) DataPath() string { return c.DataDir }
