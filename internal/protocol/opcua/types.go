package opcua

import (
	"encoding/binary"
	"fmt"
)

type NodeID struct {
	Namespace  uint16
	Identifier uint32
}

func (n NodeID) String() string { return fmt.Sprintf("ns=%d;i=%d", n.Namespace, n.Identifier) }
func ParseNodeID(value string) (NodeID, error) {
	var n NodeID
	var ns, id uint32
	if _, err := fmt.Sscanf(value, "ns=%d;i=%d", &ns, &id); err != nil || ns > 65535 {
		return n, fmt.Errorf("invalid node id")
	}
	n.Namespace = uint16(ns)
	n.Identifier = id
	return n, nil
}
func EncodeNodeID(n NodeID) []byte {
	out := make([]byte, 6)
	binary.LittleEndian.PutUint16(out, n.Namespace)
	binary.LittleEndian.PutUint32(out[2:], n.Identifier)
	return out
}
func DecodeNodeID(data []byte) (NodeID, error) {
	if len(data) < 6 {
		return NodeID{}, fmt.Errorf("short node id")
	}
	return NodeID{Namespace: binary.LittleEndian.Uint16(data), Identifier: binary.LittleEndian.Uint32(data[2:])}, nil
}
