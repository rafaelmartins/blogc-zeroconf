package main

import (
	"os"
	"path/filepath"
	"regexp"
	_ "sort"

	"github.com/blogc/go-blogc"
	"github.com/sirupsen/logrus"
)

var (
	reIndex = regexp.MustCompile(`^(README|readme|INDEX|index|_INDEX|_index)\.txt$`)
	rePosts = regexp.MustCompile(`^([A-Za-z0-9_-]+)\.txt$`)
)

type source struct {
	path   blogc.FilePath
	slug   string
	logCtx *logrus.Entry
}

func (s *source) getVariable(variable string) (string, bool, error) {
	entry := &blogc.BuildContext{
		InputFiles: []blogc.File{s.path},
	}

	return entry.GetEvaluatedVariable(variable)
}

func getSources(dir string) (*source, []*source, string) {
	logrus.WithField("directory", dir).Info("discovering sources")

	posts := []*source{}
	template := ""
	root := false

	var index *source = nil

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		logCtx := logrus.WithField("file", path)

		if err != nil {
			logCtx.Error(err)
			return nil
		}

		if info.IsDir() {
			if root {
				logCtx.Info("skipping directory")
				return filepath.SkipDir // we only want to walk the root directory
			}
			root = true
			return nil
		}

		basePath := filepath.Base(path)

		if basePath == "main.tmpl" {
			logCtx.Info("found template")
			template = path
			return nil
		}

		if matches := reIndex.FindStringSubmatch(basePath); index == nil && matches != nil && len(matches) == 2 {
			logCtx.Info("found index")
			index = &source{
				path: blogc.FilePath(path),
				slug: matches[1],
				logCtx: logrus.WithFields(logrus.Fields{
					"path": path,
					"slug": matches[1],
				}),
			}
			return nil
		}

		if matches := rePosts.FindStringSubmatch(basePath); matches != nil && len(matches) == 2 {
			logCtx.Info("found post")
			posts = append(posts, &source{
				path: blogc.FilePath(path),
				slug: matches[1],
				logCtx: logrus.WithFields(logrus.Fields{
					"path": path,
					"slug": matches[1],
				}),
			})
			return nil
		}

		logCtx.Info("skipping file")
		return nil
	})

	//sort.Slice(posts, func(i, j int) bool {
	// TODO: allow changing order of posts
	//	return posts[i].index > posts[j].index
	//})

	return index, posts, template
}
