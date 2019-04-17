package main

import (
	"fmt"
	"os"

	"github.com/blogc/go-blogc"
)

type context struct {

	// guessed from index
	title    string
	subtitle string

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

	title := "Untitled"
	subtitle := ""

	if index == nil {
		if len(posts) == 0 {
			return nil, fmt.Errorf("no sources found")
		}
	} else {
		var err error

		title, _, err = index.getVariable("TITLE")
		if err != nil {
			return nil, err
		}

		subtitle, _, err = index.getVariable("SUBTITLE")
		if err != nil {
			return nil, err
		}
	}

	rv := context{
		title:    title,
		subtitle: subtitle,
		index:    index,
		posts:    posts,
		template: template,
	}

	return &rv, nil
}

func (c *context) globalVariables(withSiteTitle bool) []string {
	rv := []string{}

	if withSiteTitle && c.title != "" {
		rv = append(rv, fmt.Sprintf("SITE_TITLE=%s", c.title))
	}

	if c.subtitle != "" {
		rv = append(rv, fmt.Sprintf("SITE_SUBTITLE=%s", c.subtitle))
	}

	return rv
}

func (c *context) getTemplate() (blogc.File, error) {
	if c.template != "" {
		return blogc.FilePath(c.template), nil
	}

	return blogc.NewFileBytes([]byte(mainTemplate))
}
