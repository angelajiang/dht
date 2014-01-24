package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

// Core Kademlia type. You can put whatever state you want in this.

const NUMBUCKETS int =  160
const NUMCONTACTS int = 20

type Kademlia struct {
    NodeID ID
    Buckets []Bucket
}

func NewKademlia() *Kademlia {
    // TODO: Assign yourself a random ID and prepare other state here.
    kptr := new(Kademlia)
    k := *kptr
    k.Buckets = make([]Bucket, NUMBUCKETS)
    k.Buckets[0] = *(NewBucket())

    return kptr
}


