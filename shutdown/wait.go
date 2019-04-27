package shutdown

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// WaitShutdown monitors SIGTERM for a signal. When it is recieved, WaitShutdown calls the server's graceful shutdown method. This function is blocking, so it should be run on a seperate goroutine
func WaitShutdown(server *http.Server, wg *sync.WaitGroup) {
	defer wg.Done()

	c := make(chan os.Signal, 2)
	defer close(c)

	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	s := <-c
	log.Printf("OS signal recieved: %s\n", s.String())

	const timeout = time.Minute/2 - 5*time.Second // 25 seconds
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
	log.Println("Server shutdown successful")
}
