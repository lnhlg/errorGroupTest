package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello world!")
	})

	server := http.Server{
		Handler: mux,
		Addr:    ":8800",
	}

	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		fmt.Println("server start")
		return server.ListenAndServe()
	})

	g.Go(func() error {
		<-ctx.Done()
		return server.Shutdown(ctx)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-quit:
			cancel()
		}
		return nil
	})

	fmt.Println("group exiting:", g.Wait().Error())
}
