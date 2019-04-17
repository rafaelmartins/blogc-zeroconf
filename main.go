package main

import (
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
	vars := ctx.globalVariables(true)

	for _, p := range ctx.posts {
		dst := blogc.FilePath(filepath.Join(out, "post", p.slug, "index.html"))

		entry := &blogc.BuildContext{
			Listing:         false,
			InputFiles:      []blogc.File{p.path},
			TemplateFile:    tmpl,
			OutputFile:      dst,
			GlobalVariables: vars,
		}

		logCtx := logrus.WithFields(logrus.Fields{
			"file": p.path.Path(),
			"post": dst.Path(),
		})

		posts = append(posts, &buildCtx{blogcCtx: entry, logCtx: logCtx})
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
		entry := &blogc.BuildContext{
			Listing:         false,
			InputFiles:      []blogc.File{ctx.index.path},
			TemplateFile:    tmpl,
			OutputFile:      dst,
			GlobalVariables: ctx.globalVariables(false),
		}

		posts = append(posts, &buildCtx{
			blogcCtx: entry,
			logCtx: logrus.WithFields(logrus.Fields{
				"file": ctx.index.path.Path(),
				"post": dst.Path(),
			}),
		})
	}

	for _, c := range posts {
		c.logCtx.Info("building")
		if err := c.blogcCtx.Build(); err != nil {
			c.logCtx.Fatal(err)
		}
	}
}
