package slogcli

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
	slogctx "github.com/veqryn/slog-context"
)

type fileContextKey struct{}

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
		lvl    slog.LevelVar
	)

	lvl.Set(slog.Level(ctx.Int("log-level")))

	if len(path) == 0 {
		w = os.Stderr
	} else {
		f, err := os.Create(path)

		if err != nil {
			return err
		}

		ctx.Context = context.WithValue(ctx.Context, fileContextKey{}, f)
		w = f
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

func SlogAfter(ctx *cli.Context) error {
	f, found := ctx.Context.Value(fileContextKey{}).(*os.File)

	if found {
		return f.Close()
	}

	return nil
}
