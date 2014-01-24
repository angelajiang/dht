package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

// Core Kademlia type. You can put whatever state you want in this.

const NUMBUCKETS int =  160
const NUMCONTACTS int = 20

type Kademlia struct {
    NodeID ID
    Buckets [NUMBUCKETS]Bucket
}

func NewKademlia() *Kademlia {
    // TODO: Assign yourself a random ID and prepare other state here.
    return new(Kademlia)
}


