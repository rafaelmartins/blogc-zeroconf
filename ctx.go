package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/blogc/go-blogc"
	"github.com/sirupsen/logrus"
)

type buildContext struct {
	blogcCtx *blogc.BuildContext
	logCtx   *logrus.Entry
}

type context struct {

	// guessed from index
	title       string
	subtitle    string
	postsPrefix string
	baseUrl     string
	baseDomain  string
	authorName  string
	authorEmail string
	withDate    bool

	// read from directory
	index          *source
	posts          []*source
	postsFiles     []blogc.File
	postsAtomFiles []blogc.File
	mainTemplate   string

	// not filled by newCtx
	mainTemplateFile blogc.File
	atomTemplateFile blogc.File
}

func newCtx() (*context, error) {
	dir, found := os.LookupEnv("SOURCE_DIR")
	if !found {
		dir = "."

		for _, d := range []string{"doc", "docs"} {
			if st, err := os.Stat(d); err == nil && st.IsDir() {
				dir = d
				break
			}
		}
	}

	ctx := context{
		title:       "Untitled",
		postsPrefix: "post",
		authorName:  "Unknown Author",
		withDate:    true,
	}

	ctx.index, ctx.posts, ctx.mainTemplate = getSources(dir)

	if ctx.index == nil {
		if len(ctx.posts) == 0 {
			return nil, fmt.Errorf("no sources found")
		}
	} else {
		var err error

		title, found, err := ctx.index.getVariable("TITLE")
		if err != nil {
			return nil, err
		}
		if found {
			ctx.title = title
		}

		subtitle, found, err := ctx.index.getVariable("SUBTITLE")
		if err != nil {
			return nil, err
		}
		if found {
			ctx.subtitle = subtitle
		}

		postsPrefix, found, err := ctx.index.getVariable("POSTS_PREFIX")
		if err != nil {
			return nil, err
		}
		if found {
			ctx.postsPrefix = postsPrefix
		}

		baseUrl, found, err := ctx.index.getVariable("BASE_URL")
		if err != nil {
			return nil, err
		}
		if found {
			ctx.baseUrl = baseUrl
		}

		baseDomain, found, err := ctx.index.getVariable("BASE_DOMAIN")
		if err != nil {
			return nil, err
		}
		if found {
			ctx.baseDomain = baseDomain
		}

		authorName, found, err := ctx.index.getVariable("AUTHOR_NAME")
		if err != nil {
			return nil, err
		}
		if found {
			ctx.authorName = authorName
		}

		authorEmail, found, err := ctx.index.getVariable("AUTHOR_EMAIL")
		if err != nil {
			return nil, err
		}
		if found {
			ctx.authorEmail = authorEmail
		}
	}

	for _, v := range ctx.posts {
		if v.timestamp == -1 {
			ctx.withDate = false
		}

		ctx.postsFiles = append(ctx.postsFiles, &v.path)
		ctx.postsAtomFiles = append(ctx.postsAtomFiles, &v.path)
	}

	sort.Slice(ctx.postsFiles, func(i int, j int) bool {
		rv := func(i int, j int) bool {
			// we assume that posts and postsFiles have the same size and order
			if ctx.posts[i].timestamp != ctx.posts[j].timestamp {
				return ctx.posts[i].timestamp > ctx.posts[j].timestamp
			}
			return ctx.posts[i].slug > ctx.posts[j].slug
		}(i, j)

		if ctx.index != nil {
			// FIXME: check value?
			_, asc, err := ctx.index.getVariable("POSTS_ASC")
			if err == nil && asc {
				return !rv
			}
		}

		return rv
	})

	return &ctx, nil
}

func (c *context) globalVariables() []string {
	rv := []string{}

	if c.title != "" {
		rv = append(rv, fmt.Sprintf("SITE_TITLE=%s", c.title))
	}

	if c.subtitle != "" {
		rv = append(rv, fmt.Sprintf("SITE_SUBTITLE=%s", c.subtitle))
	}

	if c.postsPrefix != "" {
		rv = append(rv, fmt.Sprintf("POSTS_PREFIX=%s", c.postsPrefix))
	}

	if c.baseUrl != "" {
		rv = append(rv, fmt.Sprintf("BASE_URL=%s", c.baseUrl))
	}

	if c.baseDomain != "" {
		rv = append(rv, fmt.Sprintf("BASE_DOMAIN=%s", c.baseDomain))
	}

	if c.authorName != "" {
		rv = append(rv, fmt.Sprintf("AUTHOR_NAME=%s", c.authorName))
	}

	if c.authorEmail != "" {
		rv = append(rv, fmt.Sprintf("AUTHOR_EMAIL=%s", c.authorEmail))
	}

	return rv
}

func (c *context) getBuildCtxs(out string, withTemplates bool) ([]*buildContext, error) {
	rv := []*buildContext{}

	withAtom := len(c.posts) > 0 && c.baseDomain != "" && c.withDate

	if withTemplates {
		var err error

		if c.mainTemplate != "" {
			c.mainTemplateFile = blogc.FilePath(c.mainTemplate)
		}

		if c.mainTemplateFile, err = blogc.NewFileBytes([]byte(mainTemplate)); err != nil {
			return nil, err
		}

		if withAtom {
			if c.atomTemplateFile, err = blogc.NewFileBytes([]byte(atomTemplate)); err != nil {
				return nil, err
			}
		}
	}

	vars := c.globalVariables()

	appendEntryCtx := func(src *source, dst blogc.File) {
		rv = append(rv, &buildContext{
			blogcCtx: &blogc.BuildContext{
				Listing:         false,
				InputFiles:      []blogc.File{src.path},
				TemplateFile:    c.mainTemplateFile,
				OutputFile:      dst,
				GlobalVariables: vars,
			},
			logCtx: src.logCtx.WithField("entry", dst.Path()),
		})
	}

	for _, p := range c.posts {
		appendEntryCtx(
			p,
			blogc.FilePath(filepath.Join(out, c.postsPrefix, p.slug, "index.html")),
		)
	}

	dst := blogc.FilePath(filepath.Join(out, "index.html"))

	if len(c.posts) > 0 {
		listing := &blogc.BuildContext{
			Listing:         true,
			InputFiles:      c.postsFiles,
			TemplateFile:    c.mainTemplateFile,
			OutputFile:      dst,
			GlobalVariables: vars,
		}

		logCtx := logrus.WithField("index", dst.Path())

		if c.index != nil {
			listing.ListingEntryFile = c.index.path
			logCtx = logCtx.WithField("source", c.index.path.Path())
		}

		rv = append(rv, &buildContext{
			blogcCtx: listing,
			logCtx:   logCtx,
		})

		atomDst := blogc.FilePath(filepath.Join(out, "atom.xml"))
		atomLogCtx := logrus.WithField("atom", atomDst.Path())

		if withAtom {
			rv = append(rv, &buildContext{
				blogcCtx: &blogc.BuildContext{
					Listing:         true,
					InputFiles:      c.postsAtomFiles,
					TemplateFile:    c.atomTemplateFile,
					OutputFile:      atomDst,
					GlobalVariables: append(vars, "DATE_FORMAT=%Y-%m-%dT%H:%M:%SZ"),
				},
				logCtx: atomLogCtx,
			})
		} else {
			errs := []string{}
			if c.baseDomain == "" {
				errs = append(errs, "index source BASE_DOMAIN variable (e.g. 'http://foo.com')")
			}
			if !c.withDate {
				errs = append(errs, "posts timestamp (DATE variable, e.g '2019-01-01 12:00:00')")
			}

			atomLogCtx.WithField("missing", strings.Join(errs, ", ")).Warning("atom support disabled")
		}

	} else if c.index != nil {
		appendEntryCtx(c.index, dst)
	}

	return rv, nil
}

func (c *context) close() {
	if c.mainTemplateFile != nil {
		c.mainTemplateFile.Close()
	}

	if c.atomTemplateFile != nil {
		c.atomTemplateFile.Close()
	}
}
