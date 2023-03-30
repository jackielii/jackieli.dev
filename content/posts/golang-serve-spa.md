+++
draft = true
date = 2022-09-29T15:10:15+01:00
title = "Go serve SPA front-end"
description = "Go serve SPA front-end"
slug = "go-serve-spa"
authors = []
tags = ["go", "golang", "spa"]
categories = []
externalLink = ""
series = []
+++

```go
package runtime

import (
	"bytes"
	"errors"
	"fmt"
	"hash/maphash"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"tespkg.in/df-workflow/ui"
)

var (
	vars     bytes.Buffer
	varsHash string
)

func init() {
	godotenv.Load()
	h := maphash.Hash{}
	w := io.MultiWriter(&vars, &h)
	envs := os.Environ()
	for _, v := range envs {
		if strings.HasPrefix(v, "REACT_APP_") {
			i := strings.Index(v, "=")
			if i < 0 {
				log.Fatalf("env not correct: %v", v)
			}
			key := v[:i]
			value := v[i+1:]
			value = strings.ReplaceAll(value, "\"", "\\\"")
			fmt.Fprintf(w, `window.%s="%s";`, key, value)
		}
	}
	varsHash = fmt.Sprintf("%x", h.Sum(nil))
}

func serveEnv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/javascript")
	// sensible cache expires: 10 min
	w.Header().Set("expires", time.Now().Add(time.Minute*10).Format(http.TimeFormat))
	w.Header().Set("etag", varsHash)
	if r.Header.Get("if-none-match") == varsHash {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Write(vars.Bytes())
}

func notfound(w http.ResponseWriter, _ *http.Request, err string) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(err))
}

func internal(w http.ResponseWriter, _ *http.Request, err string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err))
}

func denied(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("not authorised"))
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/_env.js" {
		serveEnv(w, r)
		return
	}
	if path == "/index.html" || path == "/" {
		// good idea, but always need a https server to work: http/2 needs https
		if pusher, ok := w.(http.Pusher); ok {
			// Push is supported.
			options := &http.PushOptions{
				Header: http.Header{
					"Accept-Encoding": r.Header["Accept-Encoding"],
				},
			}
			if err := pusher.Push("/_env.js", options); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
	}
	allow := false
	if strings.HasPrefix(path, "/static/js/__p_") {
		if !allow {
			// TODO: auth
			foo, err := r.Cookie("foo")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					denied(w, r)
					return
				}
				internal(w, r, err.Error())
				return
			}
			if foo.Value != "bar" {
				denied(w, r)
				return
			}
		}
	}
	if path == "/service-worker.js" {
		w.Header().Set("Cache-Control", "no-cache")
	}

	f, err := ui.StaticAssets.Open(filepath.Join(".", path))
	if errors.Is(err, os.ErrNotExist) {
		f, err = ui.StaticAssets.Open("index.html")
		if err != nil {
			internal(w, r, fmt.Errorf("failed to open index.html: %w", err).Error())
			return
		}
	} else if err != nil {
		internal(w, r, "open iii"+err.Error())
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		internal(w, r, err.Error())
		return
	}

	http.ServeContent(w, r, fi.Name(), fi.ModTime(), ioFile{f})
}

type ioFile struct {
	file fs.File
}

func (f ioFile) Close() error               { return f.file.Close() }
func (f ioFile) Read(b []byte) (int, error) { return f.file.Read(b) }
func (f ioFile) Stat() (fs.FileInfo, error) { return f.file.Stat() }
func (f ioFile) Seek(offset int64, whence int) (int64, error) {
	s, ok := f.file.(io.Seeker)
	if !ok {
		return 0, errors.New("not seeker")
	}
	return s.Seek(offset, whence)
}
```
