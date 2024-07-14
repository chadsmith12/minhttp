package tcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type TCPServer struct {
    Listener net.Listener
    HandlerFunc func(context.Context, net.Conn)
    wg sync.WaitGroup
}

// Listen has the TCPServer run and listen for new connections.
// All new connections get a cancellable context and are processed on a new goroutine.
func (server *TCPServer) Listen() error {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    for {
        conn, err := server.Listener.Accept()
        if err != nil {
            if errors.Is(err, net.ErrClosed) {
                return nil
            }

            return err
        }

        server.wg.Add(1)
        go func() {
            defer server.wg.Done()
            server.HandlerFunc(ctx, conn)
        }()
    }
}

// Stop will close the listener on the tcp server.
func (server *TCPServer) Stop() error {
    return server.Listener.Close()
}

// Wait will wait for the specific timeout duration for all connections to finish.
// If not all connections are finished before the timeout period, then connections will be forciably closed.
func (server *TCPServer) Wait(timeout time.Duration) {
    timeoutCh := time.After(timeout)
    doneCh := make(chan struct{})

    go func ()  {
      server.wg.Wait()
        close(doneCh)
    }()

    select {
    case <-timeoutCh:
    case <-doneCh:
    }
}

// Registers the server to start listening to terminal signal from the OS to attempt to gracefully shutdown.
func (server *TCPServer) RegisterCloseSignal() {
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

    go func() {
        <-sig
        fmt.Println("\nClosing server...")
        server.Stop()
        server.Wait(5 * time.Second)
    }()
}
