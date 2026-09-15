package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/arthurvicttor/screenshot-guard/internal/capture"
)

const version = "0.1"

func main() {
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	interval := flag.Duration("interval", 500*time.Millisecond, "clipboard polling interval")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ltime)
	logger.Printf("Screenshot Guard v%s", version)
	if *verbose {
		logger.Println("Verbose mode enabled")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	watcher := capture.NewClipboardWatcher(*interval)

	go func() {
		if err := watcher.Run(ctx); err != nil {
			logger.Printf("clipboard watcher error: %v", err)
		}
	}()

	logger.Println("Guard started (Ctrl+C to stop)")

	for ev := range watcher.Events {
		msg, ok := capture.DecideLogMessage(ev, *verbose)
		if !ok {
			continue
		}
		logger.Println(msg)
		if ev.IsImage {
			logger.Println("Action: LOGGED")
		}
	}

	logger.Println("Guard stopped")
}