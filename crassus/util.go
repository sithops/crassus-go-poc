package crassus

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func waitForSigterm() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	log.Println("Press Ctrl-C to exit...")
	for {
		time.Sleep(10 * time.Second) // or runtime.Gosched() or similar per @misterbee
	}
}

func cleanup() {
	log.Println("Exiting!")
}
