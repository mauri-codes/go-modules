package sqs_worker

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type SetShutDownConditionsInput struct {
	Configuration            *Configuration
	ShouldKeepAliveOnTimeOut func() bool
	ShutDownAction           func()
}

func SetShutDownConditions(input SetShutDownConditionsInput) (*time.Timer, chan struct{}, context.Context, context.CancelFunc) {
	timeoutCtx, cancel := context.WithCancel(context.Background())
	idleTimeoutDuration := time.Duration(input.Configuration.IdleTimeout) * time.Minute
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	idleTimer := time.NewTimer(idleTimeoutDuration)
	go func() {
		x := <-sigs
		log.Println(x)
		log.Println(x.String())
		log.Println("Shutdown signal received. Stopping message polling...")
		cancel()
	}()
	resetChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-idleTimer.C:
				if input.ShouldKeepAliveOnTimeOut() {
					log.Println("Reseting Timeout")
					idleTimer.Reset(idleTimeoutDuration)
				} else {
					log.Println("Idle timeout reached, shutting down")
					input.ShutDownAction()
					cancel()
					return
				}
			case <-resetChan:
				if !idleTimer.Stop() {
					select {
					case <-idleTimer.C:
					default:
					}
				}
				idleTimer.Reset(idleTimeoutDuration)
				log.Println("Idle timer reset")
			case <-timeoutCtx.Done():
				log.Println("Context canceled, shutting down timer goroutine")
				return
			}
		}
	}()
	go monitorSpotTermination(cancel)

	return idleTimer, resetChan, timeoutCtx, cancel
}

func monitorSpotTermination(cancel context.CancelFunc) {
	for {
		resp, err := http.Get("http://169.254.169.254/latest/meta-data/spot/termination-time")
		if err == nil && resp.StatusCode == 200 {
			log.Println("Spot instance termination scheduled!")
			cancel()
			return
		}
		time.Sleep(10 * time.Second)
	}
}
