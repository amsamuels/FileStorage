package pkg

import (
	"fmt"
	"net"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	// The underlying connection of the peer. Which in this case
	// is a TCP connection.
	conn chan net.Conn
	// if we dial and retrieve a conn => outbound == true
	// if we accept and retrieve a conn => outbound == false
	outbound bool
}

type TCPServer struct {
	listenAddr    string
	ln            net.Listener
	HandshakeFunc HandshakeFunc

	decoder Decoder

	messageChan chan Message
	peers       map[net.Addr]Peer
	shutdown    chan int
}

type Message struct {
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     make(chan net.Conn),
		outbound: outbound,
	}
}

func NewTcp(listenAddr int) *TCPServer {
	return &TCPServer{
		HandshakeFunc: NOPHandshakeFunc,
		listenAddr:    fmt.Sprintf(":%d", listenAddr),
		shutdown:      make(chan int),
	}
}

//var cancel = make(chan string)

func (t *TCPServer) Start() error {
	var err error

	t.ln, err = net.Listen("tcp", t.listenAddr)
	if err != nil {
		return fmt.Errorf("hehe")
	}

	go t.acceptLoop()

	return nil
}

func (t *TCPServer) acceptLoop() {
	for {
		conn, err := t.ln.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		go t.handelConn(conn)
	}
}

type Temp struct{}

func (t *TCPServer) handelConn(conn net.Conn) {

	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(conn); err != nil {
		fmt.Printf("TCP error: %s \n ", err)
	}

	fmt.Printf("%v", peer)

	// Read Loop
	msg := &Temp{}
	//buf := make([]byte, 2048)
	for {

		if err := t.decoder.Decoder(conn, msg); err != nil {
			fmt.Printf("TCP error: %s \n ", err)
		}

		// data, err := conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("read err")
		// }
		// msg := buf[:data]
		fmt.Println(msg)
	}
}

func (t *TCPServer) Shutdown() {
	<-t.shutdown
}
