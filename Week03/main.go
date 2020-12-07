package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"

	"golang.org/x/sync/errgroup"
)

func serveApp(ctx context.Context) error {
	mux := http.NewServeMux()
	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: mux}
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(resp, "Hello,QCon!")
	})
	go func() {
		<-ctx.Done()
		shutdownctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		srv.Shutdown(shutdownctx)
		log.Printf("Shutdown App")
	}()
	log.Printf("Start App")
	return srv.ListenAndServe()
}

func serveDebug(ctx context.Context) error {
	srv := &http.Server{Addr: "127.0.0.1:8081", Handler: http.DefaultServeMux}
	go func() {
		<-ctx.Done()
		shutdownctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		srv.Shutdown(shutdownctx)
		log.Printf("Shutdown Debug")
	}()
	log.Printf("Start Debug")
	return srv.ListenAndServe()
}

func processSignal(ctx context.Context, c chan os.Signal) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case s := <-c:
			log.Printf("get a signal %s", s.String())
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				return errors.New("Close by signal " + s.String())
			case syscall.SIGHUP:
			default:
				return errors.New("Undefined signal")
			}
		}
	}
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return serveApp(ctx)
	})
	g.Go(func() error {
		return serveDebug(ctx)
	})
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	g.Go(func() error {
		return processSignal(ctx, c)
	})
	if err := g.Wait(); err != nil {
		log.Printf("Server Error:%v\n", err)
	}
}
