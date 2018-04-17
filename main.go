package main

import (
	"io"
	"os"
	"path/filepath"

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

	tmpl, err := ctx.getTemplate()
	if err != nil {
		logrus.Fatal(err)
	}
	defer tmpl.Close()

	posts := []*buildCtx{}
	postsFiles := []blogc.File{}
	vars := ctx.globalVariables()

	appendEntryCtx := func(src blogc.File, dst blogc.File) {
		posts = append(posts, &buildCtx{
			blogcCtx: &blogc.BuildContext{
				Listing:         false,
				InputFiles:      []blogc.File{src},
				TemplateFile:    tmpl,
				OutputFile:      dst,
				GlobalVariables: vars,
			},
			logCtx: logrus.WithFields(logrus.Fields{
				"file":  src.Path(),
				"entry": dst.Path(),
			}),
		})
	}

	for _, p := range ctx.posts {
		appendEntryCtx(
			p.path,
			blogc.FilePath(filepath.Join(out, "post", p.slug, "index.html")),
		)
		postsFiles = append(postsFiles, p.path)
	}

	dst := blogc.FilePath(filepath.Join(out, "index.html"))

	if len(posts) > 0 {
		listing := &blogc.BuildContext{
			Listing:         true,
			InputFiles:      postsFiles,
			TemplateFile:    tmpl,
			OutputFile:      dst,
			GlobalVariables: vars,
		}

		logCtx := logrus.WithField("index", dst.Path())

		if ctx.index != nil {
			listing.ListingEntryFile = ctx.index.path
			logCtx = logCtx.WithField("file", ctx.index.path.Path())
		}

		posts = append(posts, &buildCtx{
			blogcCtx: listing,
			logCtx:   logCtx,
		})

	} else if ctx.index != nil {
		appendEntryCtx(ctx.index.path, dst)
	}

	if len(os.Args) > 1 && os.Args[1] == "clean" {
		for _, c := range posts {
			c.logCtx.Info("removing")
			if err := os.Remove(c.blogcCtx.OutputFile.Path()); err != nil {
				c.logCtx.Fatal(err)
			}
		}

		dirs := []string{}
		filepath.Walk(out, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				return nil
			}

			if err != nil {
				logrus.WithField("dir", path).Error(err)
				return nil
			}

			// prepend to slice, because we want subdirectories first
			dirs = append([]string{path}, dirs...)
			return nil
		})

		for _, dir := range dirs {
			logCtx := logrus.WithField("dir", dir)

			f, err := os.Open(dir)
			if err != nil {
				logCtx.Fatal(err)
			}
			defer f.Close()

			if _, err = f.Readdirnames(1); err != io.EOF {
				logCtx.Warning("directory not empty")
				continue
			}

			logCtx.Info("removing")
			if err := os.Remove(dir); err != nil {
				logCtx.Fatal(err)
			}
		}

		return
	}

	for _, c := range posts {
		c.logCtx.Info("building")
		if err := c.blogcCtx.Build(); err != nil {
			c.logCtx.Fatal(err)
		}
	}
}
