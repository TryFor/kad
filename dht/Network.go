package dht

import "net"

type Network struct {
	ownNode *Node
	Conn    *net.UDPConn
	broker  *Broker
}

type Broker struct {

}
