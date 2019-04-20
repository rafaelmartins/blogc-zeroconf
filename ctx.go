package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/blogc/go-blogc"
)

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
	index    *source
	posts    []*source
	template string
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

	index, posts, template := getSources(dir)

	ctx := context{
		title:       "Untitled",
		postsPrefix: "post",
		authorName:  "Unknown Author",
		index:       index,
		posts:       posts,
		template:    template,
		withDate:    true,
	}

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
			break
		}
	}

	sort.Slice(ctx.posts, func(i int, j int) bool {
		rv := func(i int, j int) bool {
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

func (c *context) getTemplate() (blogc.File, error) {
	if c.template != "" {
		return blogc.FilePath(c.template), nil
	}

	return blogc.NewFileBytes([]byte(mainTemplate))
}

func (c *context) getAtomTemplate() (blogc.File, error) {
	return blogc.NewFileBytes([]byte(atomTemplate))
}
