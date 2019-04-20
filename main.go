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

	appendEntryCtx := func(src *source, dst blogc.File) {
		posts = append(posts, &buildCtx{
			blogcCtx: &blogc.BuildContext{
				Listing:         false,
				InputFiles:      []blogc.File{src.path},
				TemplateFile:    tmpl,
				OutputFile:      dst,
				GlobalVariables: vars,
			},
			logCtx: src.logCtx.WithField("entry", dst.Path()),
		})
	}

	for _, p := range ctx.posts {
		appendEntryCtx(
			p,
			blogc.FilePath(filepath.Join(out, ctx.postsPrefix, p.slug, "index.html")),
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
			logCtx = logCtx.WithField("source", ctx.index.path.Path())
		}

		posts = append(posts, &buildCtx{
			blogcCtx: listing,
			logCtx:   logCtx,
		})

		atomDst := blogc.FilePath(filepath.Join(out, "atom.xml"))
		atomLogCtx := logrus.WithField("atom", atomDst.Path())

		if ctx.baseDomain != "" && ctx.authorName != "" && ctx.authorEmail != "" {
			atomTmpl, err := ctx.getAtomTemplate()
			if err != nil {
				atomLogCtx.Fatal(err)
			}
			defer atomTmpl.Close()

			posts = append(posts, &buildCtx{
				blogcCtx: &blogc.BuildContext{
					Listing:         true,
					InputFiles:      postsFiles,
					TemplateFile:    atomTmpl,
					OutputFile:      atomDst,
					GlobalVariables: append(vars, "DATE_FORMAT=%Y-%m-%dT%H:%M:%SZ"),
				},
				logCtx: atomLogCtx,
			})
		} else {
			atomLogCtx.Warning("atom feed disabled. to generate, add BASE_DOMAIN, AUTHOR_NAME and AUTHOR_EMAIL to index file")
		}

	} else if ctx.index != nil {
		appendEntryCtx(ctx.index, dst)
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
