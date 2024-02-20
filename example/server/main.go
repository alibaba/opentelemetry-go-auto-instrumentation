package main

import (
	"os"
	"os/signal"
	"server/pkgs"
	"syscall"
)

func main() {
	go func() {
		pkgs.SetupHttp()
	}()

	go func() {
		pkgs.SetupGin()
	}()

	go func() {
		pkgs.SetupGRPC()
	}()

	go func() {
		pkgs.SetMux()
	}()

	go func() {
		pkgs.SetupEcho()
	}()

	signalCh := make(chan os.Signal, 1)

	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	<-signalCh

	os.Exit(0)
}
