package loadbalance

import (
	"google.golang.org/grpc/balancer"
)

type Aperture interface {
	// Next returns next selected item.
	Next() (interface{}, func(balancer.DoneInfo))
	// Set logical width
	SetLogicalWidth(width int)
	// Set local peer id
	SetLocalPeerID(string)
	// Set local peers.
	SetLocalPeers([]string)
	// Set remote peers.
	SetRemotePeers([]interface{})
}

type P2C interface {
	// Next returns next selected item.
	Next() (interface{}, func(balancer.DoneInfo))
	// Add a item.
	Add(interface{}, float64)
}