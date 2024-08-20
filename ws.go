package ws

import (
	"context"
	"net"
	"net/http"
)

type Status struct {
	Err    <-chan error
	Done   <-chan struct{}
	Listen net.Addr
}

func (s Status) Wait() (err error) {
	select {
	case <-s.Done:
	case err = <-s.Err:
	}
	return
}

func Serve(ctx context.Context, addr string, mux http.Handler) (ss Status) {
	done, listen, err := make(chan struct{}), make(chan net.Addr, 1), make(chan error, 1)
	ss.Done, ss.Err = done, err

	srv := &http.Server{Handler: mux, Addr: addr}

	srv.BaseContext = func(ln net.Listener) context.Context {
		listen <- ln.Addr()
		return ctx
	}

	go func() {
		select {
		case <-done:
		case <-ctx.Done():
			srv.Shutdown(context.TODO())
		}
	}()

	go func() {
		e := srv.ListenAndServe()
		if e != nil && e != http.ErrServerClosed {
			err <- e
		}
		close(done)
		close(err)
		close(listen)
	}()

	ss.Listen = <-listen
	return
}
