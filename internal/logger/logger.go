package logger

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"time"
)

func NewLogger() *slog.Logger {
	return slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
			AddSource:  true,
		}),
	)
}
