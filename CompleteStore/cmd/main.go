package main

import (
	"CompleteStore/pkg"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	tcpOpts := pkg.TCPTransportOpts{
		ListenAddr:    ":8080",
		HandshakeFunc: pkg.NOPHandshakeFunc,
		Decoder:       pkg.DefaultDecoder{},
	}
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM)

	tcp := pkg.NewTcp(tcpOpts)
	fmt.Println("starting ")
	go func() {
		tcp.Start()
	}()

	<-quitChan
	tcp.Shutdown()

}
