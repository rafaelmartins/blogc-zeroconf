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

	appendEntryCtx := func(src blogc.File, dst blogc.File, vars []string) {
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
			vars,
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
		appendEntryCtx(ctx.index.path, dst, ctx.globalVariables(false))
	}

	for _, c := range posts {
		c.logCtx.Info("building")
		if err := c.blogcCtx.Build(); err != nil {
			c.logCtx.Fatal(err)
		}
	}
}
