package main

import (
	"os"

	"github.com/blogc/go-blogc"
	"github.com/sirupsen/logrus"
)

type buildCtx struct {
	blogcCtx *blogc.BuildContext
	logCtx   *logrus.Entry
}

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

	ctx, err := newCtx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer ctx.close()

	if len(os.Args) <= 1 {
		build(ctx, out)
		return
	}

	switch os.Args[1] {
	case "build":
		err = build(ctx, out)
	case "clean":
		err = clean(ctx, out)
	default:
		logrus.Fatalf("command not found: %s", os.Args[1])
	}

	if err != nil {
		logrus.Fatal(err)
	}
}
