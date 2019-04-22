package main

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func runserver(ctx *context, out string) error {
	fs := http.FileServer(http.Dir(out))
	http.Handle("/", fs)

	listenAddr, found := os.LookupEnv("LISTEN_ADDR")
	if !found {
		listenAddr = ":8080"
	}

	if err := build(ctx, out); err != nil {
		return err
	}

	logrus.WithField("listen_addr", listenAddr).Info("listening HTTP")
	return http.ListenAndServe(listenAddr, nil)
}
