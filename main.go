package main

import (
	"os"

	"github.com/blogc/go-blogc"
	"github.com/sirupsen/logrus"
)

func main() {
	out, found := os.LookupEnv("OUTPUT_DIR")
	if !found {
		out = "_build"
	}

	level, found := os.LookupEnv("LOG_LEVEL")
	if !found {
		level = logrus.InfoLevel.String()
	}

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetLevel(lvl)

	command := "build"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	logCtx := logrus.WithFields(logrus.Fields{
		"blogc_version": blogc.Version,
		"command":       command,
	})

	logCtx.Info("blogc-zeroconf")

	ctx, err := newContext()
	if err != nil {
		logrus.Fatal(err)
	}
	defer ctx.close()

	switch command {
	case "build":
		err = build(ctx, out)
	case "clean":
		err = clean(ctx, out)
	case "runserver":
		err = runserver(ctx, out)
	default:
		logCtx.Fatal("command not found")
	}

	if err != nil {
		logCtx.Fatal(err)
	}
}
