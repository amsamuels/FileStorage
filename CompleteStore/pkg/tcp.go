package pkg

import (
	"fmt"
	"net"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	// The underlying connection of the peer. Which in this case
	// is a TCP connection.
	conn net.Conn
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
	ln     net.Listener
	rpcch  chan RPC
	ONPeer func(Peer) error
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close implements the peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func NewTcp(opts TCPTransportOpts) *TCPServer {
	return &TCPServer{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

// Consume implements the Transport interface,
// returns a channel which you can read the incomming messages from
func (t *TCPServer) Consume() <-chan RPC {

	return t.rpcch
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

func (t *TCPServer) handelConn(conn net.Conn) {

	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s \n ", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.ONPeer != nil {
		if err = t.ONPeer(peer); err != nil {
			return
		}
	}
	// Read Loop
	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)
		if err != nil {
			return
		}

		rpc.From = conn.RemoteAddr()

		t.rpcch <- rpc
	}
}

func (t *TCPServer) Shutdown() {
	<-t.shutdown
}
