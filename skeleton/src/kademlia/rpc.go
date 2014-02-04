package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
    "net"
    "fmt"
)


// Host identification.
type Contact struct {
    NodeID ID
    Host net.IP
    Port uint16
}


// PING
type Ping struct {
    Sender Contact
    MsgID ID
}

type Pong struct {
    MsgID ID
    Sender Contact
}

func (k *Kademlia) Ping(ping Ping, pong *Pong) error {
    // This one's a freebie.
    pong.MsgID = CopyID(ping.MsgID)
    fmt.Printf("ping.MsgID from RPC call: %v\n", ping.MsgID.AsString())
    fmt.Printf("pong.MsgID from RPC call: %v\n", pong.MsgID.AsString())
    pong.Sender.NodeID = k.NodeID
    pong.Sender.Host = k.Host
    pong.Sender.Port = k.Port
    return nil
}


// STORE
type StoreRequest struct {
    Sender Contact
    MsgID ID
    Key ID
    Value []byte
}

type StoreResult struct {
    MsgID ID
    Err error
}

func (k *Kademlia) Store(req StoreRequest, res *StoreResult) error {
    k.Data[req.Key] = req.Value
    res.MsgID = CopyID(req.MsgID)
    fmt.Printf("I now have: %v\n", k.Data[req.Key]) 
    return nil
}


// FIND_NODE
type FindNodeRequest struct {
    Sender Contact
    MsgID ID
    NodeID ID
}

type FoundNode struct {
    IPAddr string
    Port uint16
    NodeID ID
}

type FindNodeResult struct {
    MsgID ID
    Nodes []FoundNode
    Err error
}

/*func (k *Kademlia) FindNode(req FindNodeRequest, res *FindNodeResult) error {
    //check if we're the node in question
    start_from_one := false //hacky way to know if we should fill the array from 0 or 1
    for key, _ := range k.Data {
        if key == req.NodeID {
            start_from_one = true
            res.Nodes[0].NodeID = k.NodeID
            res.Nodes[0].IPAddr = k.Host.String()
            res.Nodes[0].Port = k.Port
        }
    }
    if start_from_one == true {
    //loop through each contact in each bucket to calculate distance
        for index, bucket := range k.Buckets {
            for num, contact := range k.Buckets[index].Contacts {
                //distance(x, y) = x ^ y                 
                //there is an Xor function in id.go!
                //XOR the distance between yourself and your contacts.
                //find the .. 2 or 3 closest contacts?
            }
        }
    }

    return nil
}
*/

// FIND_VALUE
type FindValueRequest struct {
    Sender Contact
    MsgID ID
    Key ID
}

// If Value is nil, it should be ignored, and Nodes means the same as in a
// FindNodeResult.
type FindValueResult struct {
    MsgID ID
    Value []byte
    Nodes []FoundNode
    Err error
}

func (k *Kademlia) FindValue(req FindValueRequest, res *FindValueResult) error {
    res.MsgID = CopyID(req.MsgID)
    if val,ok := k.Data[req.Key]; ok {
        res.Value = val
    }else{
        res.Value = nil
        //call find node
    }
    return nil
}

