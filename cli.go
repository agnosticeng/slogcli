package slogcli

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
	slogctx "github.com/veqryn/slog-context"
)

func SlogFlags() []cli.Flag {
	return []cli.Flag{
		&cli.IntFlag{
			Name:    "log-level",
			Value:   0,
			EnvVars: []string{"LOG_LEVEL"},
		},
		&cli.StringFlag{
			Name:    "log-path",
			EnvVars: []string{"LOG_PATH"},
		},
		&cli.StringFlag{
			Name:    "log-format",
			Value:   "TEXT",
			EnvVars: []string{"LOG_FORMAT"},
		},
	}
}

func SlogBefore(ctx *cli.Context) error {
	var (
		path   = ctx.String("log-path")
		format = ctx.String("log-format")
		w      io.WriteCloser
		err    error
		lvl    slog.LevelVar
	)

	lvl.Set(slog.Level(ctx.Int("log-level")))

	if len(path) == 0 {
		w = os.Stderr
	} else {
		w, err = os.Create(path)
	}

	if err != nil {
		return err
	}

	var handler slog.Handler

	switch format {
	case "text", "TEXT":
		handler = slog.NewTextHandler(w, &slog.HandlerOptions{AddSource: true, Level: &lvl})
	case "json", "JSON":
		handler = slog.NewJSONHandler(w, &slog.HandlerOptions{AddSource: true, Level: &lvl})
	default:
		return fmt.Errorf("unknown log format: %s", format)
	}

	var (
		slogCtxHandler = slogctx.NewHandler(handler, nil)
		logger         = slog.New(slogCtxHandler)
	)

	slog.SetDefault(logger)
	ctx.Context = slogctx.NewCtx(ctx.Context, logger)
	return nil
}
