package main

import (
	"github.com/DyadyaRodya/gofermart/internal/app"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultServerAddress        = `:8080`
	defaultAccrualServerAddress = "http://localhost:8100"
	defaultLogLevel             = "debug"
	defaultSaltSize             = 16
)

func main() {
	server := app.InitApp(defaultServerAddress, defaultAccrualServerAddress, defaultLogLevel, defaultSaltSize)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		s := <-c
		err := server.Shutdown(s)
		if err != nil {
			panic(err)
		}
	}()
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
