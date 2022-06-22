package test_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// signal.Notify(signalChan, syscall.SIGINT)
	// signal.Notify(signalChan, syscall.SIGTERM)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				doLog("The test is still running! Don't kill me!")
			}
			time.Sleep(10 * time.Second)
		}
	}()
	// go func() {
	// 	<-signalChan
	// 	doLog("Interrupt signal received. Quitting...")
	// 	os.Exit(1)
	// }()
	exitVal := m.Run()
	cancel()
	os.Exit(exitVal)
}

// DoLog logs the given arguments to the given writer, along with a timestamp.
func doLog(args ...interface{}) {
	date := time.Now()
	prefix := fmt.Sprintf("%s:", date.Format(time.RFC3339))
	allArgs := append([]interface{}{prefix}, args...)
	fmt.Println(allArgs...) //nolint:forbidigo
}
