package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

// Core Kademlia type. You can put whatever state you want in this.

import (
    "net"
)

const NUMBUCKETS int =  160
const NUMCONTACTS int = 1

type Kademlia struct {
    NodeID ID
    Buckets []Bucket
}

func NewKademlia() *Kademlia {
    // TODO: Assign yourself a random ID and prepare other state here.
    kptr := new(Kademlia)
    kptr.NodeID = NewRandomID()
    kptr.Buckets = make([]Bucket, NUMBUCKETS)
    kptr.Buckets[0] = *(NewBucket())

    return kptr
}

func DoPing(remote_host net.IP, port uint16) (Pong, error){
    tmp_id := NewRandomID()
    tmp_ip := net.ParseIP("127.0.0.1")
    tmp_contact := Contact{tmp_id, tmp_ip, 4000}
    tmp_pong := Pong{tmp_id, tmp_contact}
    return tmp_pong, nil
}


