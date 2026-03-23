package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	logFileName = "unraid-tui.log"
	maxLogSize  = 5 * 1024 * 1024 // 5 MB
	fileMode    = 0600
)

// Init sets up file-based logging in the given directory.
// Returns a close function to flush and close the log file.
// If logging cannot be initialized, it falls back to discarding logs.
func Init(dir string) func() {
	if err := os.MkdirAll(dir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "logging: cannot create dir %s: %v\n", dir, err)
		slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
		return func() {}
	}

	logPath := filepath.Join(dir, logFileName)
	rotate(logPath)

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, fileMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logging: cannot open %s: %v\n", logPath, err)
		slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
		return func() {}
	}

	handler := slog.NewTextHandler(f, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))

	return func() { f.Close() }
}

// rotate renames the log file to .log.1 if it exceeds maxLogSize.
func rotate(path string) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}
	if info.Size() < maxLogSize {
		return
	}
	os.Rename(path, path+".1")
}
