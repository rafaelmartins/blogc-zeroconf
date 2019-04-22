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
		if !c.blogcCtx.NeedsBuild() {
			c.logCtx.Trace("up to date")
			continue
		}
		c.logCtx.Info("building")
		if err := c.blogcCtx.Build(); err != nil {
			return err
		}
	}

	for basePath, source := range ctx.copy {
		dst := filepath.Join(out, basePath)
		logCtx := logrus.WithFields(logrus.Fields{
			"from": source,
			"copy": dst,
		})

		needsCopy := func() bool {
			st, err := os.Stat(source)
			if err != nil {
				return false // source not found?
			}
			smtime := st.ModTime()

			st, err = os.Stat(dst)
			if err != nil {
				return true
			}
			dmtime := st.ModTime()

			return dmtime.Before(smtime)
		}()

		if !needsCopy {
			logCtx.Trace("up to date")
			continue
		}

		from, err := os.Open(source)
		if err != nil {
			return err
		}
		defer from.Close()

		to, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		defer to.Close()

		logCtx.Info("copying")
		_, err = io.Copy(to, from)
		if err != nil {
			return err
		}
	}

	return nil
}
