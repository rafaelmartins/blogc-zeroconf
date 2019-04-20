package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func clean(ctx *context, out string) error {
	bctxs, err := ctx.getBuildContexts(out, false)
	if err != nil {
		return err
	}

	for _, c := range bctxs {
		if err := os.Remove(c.blogcCtx.OutputFile.Path()); err != nil {
			if os.IsNotExist(err) {
				c.logCtx.Warning("not found, skipping")
				continue
			}
			return err
		}
		c.logCtx.Info("removed")
	}

	dirs := []string{}
	filepath.Walk(out, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logCtx := logrus.WithField("path", path)
			if os.IsNotExist(err) {
				logCtx.Warning("not found, skipping")
				return nil
			}
			logCtx.Error(err)
			return nil
		}

		if !info.IsDir() {
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

		if err := os.Remove(dir); err != nil {
			if os.IsNotExist(err) {
				logCtx.Warning("not found, skipping")
				continue
			}
			return err
		}
		logCtx.Info("removed")
	}

	return nil
}
