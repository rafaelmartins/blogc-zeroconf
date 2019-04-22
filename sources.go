package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/blogc/go-blogc"
	"github.com/sirupsen/logrus"
)

var (
	reIndex = regexp.MustCompile(`^(README|readme|INDEX|index|_INDEX|_index)\.txt$`)
	rePosts = regexp.MustCompile(`^([A-Za-z0-9_-]+)\.txt$`)
)

type source struct {
	path      blogc.FilePath
	slug      string
	logCtx    *logrus.Entry
	timestamp int64
}

func (s *source) getVariable(variable string) (string, bool, error) {
	entry := &blogc.BuildContext{
		Listing:    false,
		InputFiles: []blogc.File{s.path},
	}

	return entry.GetEvaluatedVariable(variable)
}

func (s *source) setTimestamp() {
	ctx := &blogc.BuildContext{
		Listing:         false,
		InputFiles:      []blogc.File{s.path},
		GlobalVariables: []string{"DATE_FORMAT=%s"},
	}

	v, found, err := ctx.GetEvaluatedVariable("DATE_FORMATTED")
	if err != nil || !found {
		s.logCtx.Warning("failed to get post timestamp")
		return
	}

	t, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		s.logCtx.Warning("failed to parse post timestamp")
		return
	}

	s.timestamp = t
}

func getSources(dir string) (*source, []*source, map[string]string, string) {
	logrus.WithField("path", dir).Trace("discovering sources")

	posts := []*source{}
	copy := map[string]string{}
	template := ""
	root := false

	var index *source = nil

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		logCtx := logrus.WithField("path", path)

		if err != nil {
			logCtx.Error(err)
			return nil
		}

		if info.IsDir() {
			if root {
				logCtx.Trace("skipping directory")
				return filepath.SkipDir // we only want to walk the root directory
			}
			root = true
			return nil
		}

		basePath := filepath.Base(path)

		if basePath == "main.tmpl" {
			logCtx.Trace("found template")
			template = path
			return nil
		}

		if matches := reIndex.FindStringSubmatch(basePath); index == nil && matches != nil && len(matches) == 2 {
			index = &source{
				path: blogc.FilePath(path),
				slug: matches[1],
				logCtx: logrus.WithFields(logrus.Fields{
					"source": path,
					"slug":   matches[1],
				}),
				timestamp: -1,
			}
			index.logCtx.Trace("found index")
			return nil
		}

		if matches := rePosts.FindStringSubmatch(basePath); matches != nil && len(matches) == 2 {
			entry := &source{
				path: blogc.FilePath(path),
				slug: matches[1],
				logCtx: logrus.WithFields(logrus.Fields{
					"source": path,
					"slug":   matches[1],
				}),
				timestamp: -1,
			}
			entry.logCtx.Trace("found post")
			entry.setTimestamp()
			posts = append(posts, entry)
			return nil
		}

		if basePath[0:1] != "." && basePath[0:1] != "_" {
			copy[basePath] = path
			logrus.WithField("source", path).Trace("found copy")
			return nil
		}

		logCtx.Trace("skipping file")
		return nil
	})

	return index, posts, copy, template
}
