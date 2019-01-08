package dht

import (
	"time"
	"math/big"
	"math/rand"
	"fmt"
	"log"
	"io"
	"bytes"
	"bufio"
	"encoding/binary"
)

type Bucket struct {
	Min, Max   *big.Int
	Nodes      []*NodeInfo
	LastUpdate time.Time
}

func NewBucket(min, max *big.Int) *Bucket {
	b := new(Bucket)
	b.Min = min
	b.Max = max
	b.Nodes = nil
	b.LastUpdate = time.Now()
	return b
}

func (bucket *Bucket) Touch() {
	bucket.LastUpdate = time.Now()
}

func (bucket *Bucket) Len() int {
	return len(bucket.Nodes)
}

func (bucket *Bucket) Add(n *NodeInfo) {
	bucket.Nodes = append(bucket.Nodes, n)
	bucket.Touch()
}

func (bucket *Bucket) updateIfExists(node *NodeInfo) bool {
	for i,n := range bucket.Nodes {
		if n.ID.CompareTo(node.ID) == 0 {
			bucket.Nodes[i] = node
			bucket.Touch()
			return true
		}
	}
	return  false
}

func (bucket *Bucket) Copy(result *[]*NodeInfo, maxsize int) int {
	nw := 0
	for _,n := range bucket.Nodes {
		if n.Status == GOOD {
			*result = append(*result, n)
			nw++
			if len(*result) == maxsize {
				break
			}
		}
	}
	return nw
}

func (bucket *Bucket) RandID() Identifier{
	d := bisub(bucket.Max, bucket.Min)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	z := biadd(bucket.Min, d.Rand(random,d))
	ret := make([]byte,20)
	for idx, b := range z.Bytes() {
		ret[idx] = b
	}
	return ret
}

func (bucket *Bucket) Print() {
	for _, ni := range bucket.Nodes {
		fmt.Printf("\t%s\n", ni)
	}
}

type Routing struct {
	ownNode *Node
	table []*Bucket
	log *log.Logger
}


func NewRouting(ownNode *Node) *Routing {
	routing := new(Routing)
	routing.ownNode = ownNode
	b := NewBucket(binew(0), bilsh(1, 160))
	routing.table = make([]*Bucket, 1)
	routing.table[0] = b
	data, err := GetPersist().LoadNodeInfo(ownNode.ID())
	if err == nil && len(data) > 0 {
		err := routing.LoadRouting(bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}
	}
	return routing

}

func (routing *Routing) Len() int{
	var length int
	for _,v := range routing.table{
		length += v.Len()
	}
	return length
}


func (routing *Routing) LoadRouting(reader io.Reader) error {
	buf := bufio.NewReader(reader)
	var data []byte = make([]byte, 24)
	_, err := buf.Read(data)
	if err != nil {
		return err
	}
	var length uint32 = 0
	err = binary.Read(bytes.NewReader(data[20:24]), binary.LittleEndian, &length)
	if err != nil {
		return err
	}
	var stream []byte = make([]byte, length)
	_, err = buf.Read(stream)
	if err != nil {
		return err
	}
	nodes := ParseBytesStream(stream)
	for _, v := range nodes {
		routing.InsertNode(v)
	}
	return nil
}


func (routing *Routing) InsertNode(other *NodeInfo) {
	if routing.isMe(other) {
		return
	}
	bucket, idx := routing.findBucket(other.ID)
	if bucket.updateIfExists(other) {
		return
	}
	if bucket.Len() < K{
		bucket.Add(other)
		return
	}
	if idx == len(routing.table)-1 {

	}

}


func (routing *Routing) isMe(other *NodeInfo) bool{
	return routing.ownNode.Info.ID.CompareTo(other.ID) == 0
}

func (routing *Routing) findBucket(dst Identifier)(*Bucket, int) {
	idx := routing.bucketIndex(dst)
	length := len(routing.table)
	if length > idx {
		return routing.table[idx],idx
	}
	return routing.table[length-1],length-1
}

func (routing *Routing) bucketIndex(dst Identifier) int {
	return BucketIndex(routing.ownNode.ID(), dst)
}





