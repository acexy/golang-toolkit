package sys

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	toolkitError "github.com/acexy/golang-toolkit/error"
)

var defaultShutdownSig = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}

var shutdownSignalMu sync.Mutex
var shutdownSignals = defaultShutdownSig
var shutdownSignalsInitialized bool
var exitSignChan = make(chan struct{})
var exitOnce sync.Once

// SetShutdownSignals sets nonempty process shutdown signals before the first subscription.
func SetShutdownSignals(sig ...os.Signal) error {
	if len(sig) == 0 {
		return toolkitError.ErrInvalidShutdownSignals
	}
	for _, item := range sig {
		if item == nil {
			return toolkitError.ErrInvalidShutdownSignals
		}
	}

	shutdownSignalMu.Lock()
	defer shutdownSignalMu.Unlock()
	if shutdownSignalsInitialized {
		return toolkitError.ErrShutdownSignalsInitialized
	}
	shutdownSignals = append([]os.Signal(nil), sig...)
	return nil
}

// ShutdownHolding blocks until the shared process shutdown signal is received.
func ShutdownHolding() {
	<-ShutdownSignal()
}

// ShutdownCallback waits for shutdown and runs a nonnil callback synchronously.
func ShutdownCallback(f func()) {
	<-ShutdownSignal()
	if f != nil {
		f()
	}
}

// ShutdownSignal returns the shared shutdown channel after synchronously registering the OS signals.
// The first received signal closes the channel for all waiters; registration happens only once.
func ShutdownSignal() <-chan struct{} {
	exitOnce.Do(func() {
		shutdownSignalMu.Lock()
		shutdownSignalsInitialized = true
		sig := append([]os.Signal(nil), shutdownSignals...)
		shutdownSignalMu.Unlock()
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, sig...)
		go func() {
			defer signal.Stop(ch)
			<-ch
			close(exitSignChan)
		}()
	})
	return exitSignChan
}

// ShutdownContext returns a context canceled by shutdown, its parent, or the returned cancel function.
// The shared signal subscription is initialized before this function returns.
func ShutdownContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	ctxNew, cancel := context.WithCancel(ctx)
	shutdown := ShutdownSignal()
	go func() {
		select {
		case <-shutdown:
			cancel()
		case <-ctxNew.Done():
		}
	}()
	return ctxNew, cancel
}

// ShutdownCallbackDeadline waits for shutdown, then waits for the callback up to deadline.
// A timeout does not terminate a callback that is still running.
func ShutdownCallbackDeadline(f func(), deadline time.Duration) {
	<-ShutdownSignal()
	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	done := make(chan struct{})
	go func() {
		defer close(done)
		if f != nil {
			f()
		}
	}()
	select {
	case <-ctx.Done():
	case <-done:
	}
}
