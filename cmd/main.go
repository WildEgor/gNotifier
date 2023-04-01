package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	server "github.com/WildEgor/gNotifier/internal"
	log "github.com/sirupsen/logrus"
)

func main() {
	server, _ := server.NewServer()
	log.Fatal(server.Listen(fmt.Sprintf(":%v", "8888")))

	// block main thread - wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println()
		log.Println(sig)
		done <- true
	}()

	log.Println("[Main] Awaiting signal")
	<-done
	log.Println("[Main] Stopping consumer")
}
