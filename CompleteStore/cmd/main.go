package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	tcpserver "pdfparser/pkg"
	"syscall"
)

func main() {

	TcpPort := flag.Int("port", 8080, "The port to listen on; default is 8080.")
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM)

	tcp := tcpserver.NewTcp(*TcpPort)
	fmt.Println("starting ")
	go func() {
		tcp.Start()
	}()
	<-quitChan
	tcp.Shutdown()

}
