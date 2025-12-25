package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Wait(ctx context.Context, cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 10*time.Second)
	defer cancelTimeout()

	cancel()
	<-timeoutCtx.Done()
}
