package pkg

import "net"

// RPC holds any arbitrary data that is being sent over the
// each transport between two nodes in the network.
type RPC struct {
	From   net.Addr
	Paylad []byte
}