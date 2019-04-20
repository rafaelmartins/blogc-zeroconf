package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func clean(ctx *context, out string) error {
	bctxs, err := ctx.getBuildCtxs(out, false)
	if err != nil {
		return err
	}

	for _, c := range bctxs {
		c.logCtx.Info("removing")
		if err := os.Remove(c.blogcCtx.OutputFile.Path()); err != nil {
			return err
		}
	}

	dirs := []string{}
	filepath.Walk(out, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}

		if err != nil {
			logrus.WithField("path", path).Error(err)
			return nil
		}

		// prepend to slice, because we want subdirectories first
		dirs = append([]string{path}, dirs...)
		return nil
	})

	for _, dir := range dirs {
		logCtx := logrus.WithField("path", dir)

		f, err := os.Open(dir)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err = f.Readdirnames(1); err != io.EOF {
			logCtx.Warning("directory not empty")
			continue
		}

		logCtx.Info("removing")
		if err := os.Remove(dir); err != nil {
			return err
		}
	}

	return nil
}
