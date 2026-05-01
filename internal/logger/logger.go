package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Logger struct {
	*slog.Logger
	file *os.File
}

func New(cfg Config) (*Logger, error) {
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	filepath := filepath.Join(
		cfg.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("openning file: %w", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	var level slog.Level
	if err := level.UnmarshalText([]byte(strings.ToUpper(cfg.Level))); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	opts := slog.HandlerOptions{
		Level: level,
	}

	logger := slog.New(slog.NewJSONHandler(multiWriter, &opts))

	return &Logger{
		logger,
		file,
	}, nil
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close app logger:", err)
	}
}
