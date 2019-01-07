package dht



type Node struct {
	Info NodeInfo
	Routing *Routing
	krpc *KRPC
	nw *Network
}

func (node *Node) ID() Identifier {
	return node.Info.ID
}
