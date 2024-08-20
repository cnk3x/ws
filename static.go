package ws

import (
	"io/fs"
	"net/http"
	"sync"
)

func Static(fsys fs.FS) func(next http.Handler) http.Handler {
	staticServe := useStaticServe()
	hs := http.FileServer(http.FS(fsys))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && staticServe(w, r, hs) {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

type staticResponse struct {
	status int
	http.ResponseWriter
	written bool
}

func (resp *staticResponse) Write(data []byte) (int, error) {
	if resp.status != 404 && resp.status != 405 {
		resp.written = true
		return resp.ResponseWriter.Write(data)
	}
	return 0, nil
}

func (resp *staticResponse) WriteHeader(statusCode int) {
	if !resp.written {
		if resp.status = statusCode; resp.status != 404 && resp.status != 405 {
			resp.ResponseWriter.WriteHeader(resp.status)
		}
	}
}

func useStaticServe() func(w http.ResponseWriter, r *http.Request, exec http.Handler) (written bool) {
	p := &sync.Pool{New: func() any { return &staticResponse{} }}
	return func(w http.ResponseWriter, r *http.Request, exec http.Handler) (written bool) {
		resp := p.Get().(*staticResponse)
		resp.ResponseWriter = w
		exec.ServeHTTP(resp, r)
		written = resp.written
		resp.ResponseWriter = nil
		resp.status = 0
		p.Put(resp)
		return
	}
}
