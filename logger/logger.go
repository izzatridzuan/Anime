package logger

import (
	"io"
	"log/slog"
	"os"
)

func Init() error {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	writer := io.MultiWriter(os.Stdout, file)
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))
	return nil
}
