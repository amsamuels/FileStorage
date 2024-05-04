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

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	shutdown      chan int
}

type TCPServer struct {
	TCPTransportOpts
	ln net.Listener

	messageChan chan ProcessMessage
	peers       map[net.Addr]Peer
}

type ProcessMessage struct {
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     make(chan net.Conn),
		outbound: outbound,
	}
}

func NewTcp(opts TCPTransportOpts) *TCPServer {
	return &TCPServer{
		TCPTransportOpts: opts,
	}
}

//var cancel = make(chan string)

func (t *TCPServer) Start() error {
	var err error

	t.ln, err = net.Listen("tcp", t.ListenAddr)
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

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP error: %s \n ", err)
	}

	fmt.Printf("%v", peer)

	// Read Loop
	msg := &RPC{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s \n ", err)
			continue
		}
		fmt.Println(msg)
	}
}

func (t *TCPServer) Shutdown() {
	<-t.shutdown
}
