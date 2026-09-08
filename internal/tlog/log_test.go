package tlog

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
)

func BenchmarkSlogHandler(b *testing.B) {
	slog.SetDefault(slog.New(get()))

	for b.Loop() {
		slog.LogAttrs(context.Background(), slog.LevelInfo,
			"INFOINFOINFOINFOINFOINFOINFOINFOINFOINFOINFOINFO",
			slog.Any("error", errors.New("ERROR")),
			slog.Bool("val", true),
			slog.Int("INT", 10),
			slog.String("STR", "INTINTINTINTINTINTINTINTINTINTINTINTINTINTINTINTINT"),
			slog.Duration("DURATION", 10*time.Second))

		slog.LogAttrs(context.Background(), slog.LevelInfo,
			"INFOINFOINFOINFOINFOINFOINFOINFOINFOINFOINFOINFO",
			slog.Any("error", errors.New("ERROR")),
			slog.Bool("val", true),
			slog.Int("INT", 10),
			slog.Int("INT", 10),
			slog.Int("INT", 10),
			slog.Int("INT", 10),
			slog.Int("INT", 10),
			slog.Int("INT", 10),
			slog.Int("INT", 10),
			slog.Int("INT", 10),
			slog.Int("INT", 10),
			slog.Int("INT", 10),
			slog.Int("INT", 10),
			slog.Int("INT", 10),
			slog.String("STR", "INTINTINTINTINTINTINTINTINTINTINTINTINTINTINTINTINT"),
			slog.String("STR", "INTINTINTINTINTINTINTINTINTINTINTINTINTINTINTINTINT"),
			slog.Duration("DURATION", 10*time.Second))

		slog.LogAttrs(context.Background(), slog.LevelInfo,
			"INFOINFOINFOINFOINFOINFOINFOINFOINFOINFOINFOINFO",
			slog.String("STR", "INTINTINTINTINTINTINTINTINTINTINTINTINTINTINTINTINT"),
			slog.Duration("DURATION", 10*time.Second))

		slog.LogAttrs(context.Background(), slog.LevelInfo,
			"INFOINFOINFOINFOINFOINFOINFOINFOINFOINFOINFOINFO",
			slog.Any("error", errors.New("ERROR")),
			slog.Any("error", errors.New("ERROR")),
			slog.Any("error", errors.New("ERROR")),
			slog.Any("error", errors.New("ERROR")),
			slog.Any("error", errors.New("ERROR")),
			slog.Bool("val", true),
			slog.Bool("val", true),
			slog.Bool("val", true),
			slog.Bool("val", true),
			slog.Bool("val", true),
			slog.Bool("val", true),
			slog.Bool("val", true),
			slog.Bool("val", true),
			slog.Bool("val", true),
			slog.Bool("val", true),
			slog.Int("INT", 10),
			slog.String("STR", "INTINTINTINTINTINTINTINTINTINTINTINTINTINTINTINTINT"),
			slog.Duration("DURATION", 10*time.Second))
	}
}
