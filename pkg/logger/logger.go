package logger

import (
	"github.com/RCSE2025/backend-go/pkg/logger/handlers/slogpretty"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func NewLogger(prod bool) *slog.Logger {
	var log *slog.Logger

	switch prod {
	case false:
		log = setupPrettySlog()
	case true:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

//func NewLogger(env string) *slog.Logger {
//	var log *slog.Logger
//
//
//
//	switch env {
//	case envLocal:
//		log = setupPrettySlog()
//	case envDev:
//		log = slog.New(
//			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
//		)
//	case envProd:
//		log = slog.New(
//			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
//		)
//	default: // If env config is invalid, set prod settings by default due to security
//		log = slog.New(
//			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
//		)
//	}
//
//	return log
//}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
