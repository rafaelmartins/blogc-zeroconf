package main

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func build(ctx *context, out string) error {
	bctxs, err := ctx.getBuildContexts(out, true)
	if err != nil {
		return err
	}

	for _, c := range bctxs {
		c.logCtx.Info("building")
		if err := c.blogcCtx.Build(); err != nil {
			return err
		}
	}

	for basePath, source := range ctx.copy {
		from, err := os.Open(source)
		if err != nil {
			return err
		}
		defer from.Close()

		dst := filepath.Join(out, basePath)
		to, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer to.Close()

		logrus.WithFields(logrus.Fields{
			"from": source,
			"to":   dst,
		}).Info("copying")

		_, err = io.Copy(to, from)
		if err != nil {
			return err
		}
	}

	return nil
}
