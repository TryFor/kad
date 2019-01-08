package dht

import "time"

type KRPC struct {
	ownNode *Node
	tid uint32
}

func  ParseBytesStream(data []byte) []*NodeInfo {
	var nodes []*NodeInfo = nil
	for j:=0; j<len(data); j = j+6 {
		if j+26 > len(data) {
			break
		}
		kn := data[j:j+26]
		node := new(NodeInfo)
		node.ID = Identifier(kn[0:20])
		node.IP = kn[20:24]
		port := kn[24:26]
		node.Port = int(port[0])<<8 + int(port[1])
		node.Status = GOOD
		node.LastSeen = time.Now()
		nodes = append(nodes, node)
	}
	return nodes
}
