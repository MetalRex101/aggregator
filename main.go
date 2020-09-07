package main

import (
	"context"
	"fmt"
	"github.com/MetalRex101/affise/src/api"
	"github.com/MetalRex101/affise/src/network"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const addr = ":8090"
const timeout = time.Second * 10

func run() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-signalChan
		cancel()
	}()

	mux := http.NewServeMux()

	handler := api.NewHandler(api.NewScript(network.NewHttpClient()))
	mux.HandleFunc("/", handler.Handle)

	server := &http.Server{
		Addr:         addr,
		Handler:      network.RateLimit(mux),
		WriteTimeout: timeout,
	}

	serve(ctx, server)
}

func serve(ctx context.Context, server *http.Server) {
	fmt.Println(fmt.Sprintf("[*] http started on %s", addr))

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(fmt.Sprintf("failed to serve http: %s", err))
		}
	}()

	<-ctx.Done()

	// could replace with timeout context if jobs on http endpoint could be huge
	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Println(fmt.Sprintf("[*] failed to shutdown server properly: %s", err))
	}
}

func main() {
	run()
}
