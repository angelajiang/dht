package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
    "net"
    "fmt"
    "net/rpc"
    "log"
)

// Host identification.
type Contact struct {
    NodeID ID
    Host net.IP
    Port uint16
}

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
    if val, ok := k.Data[req.Key]; ok {
        res.Value = val

    } else {
        res.Value = nil
        closestContacts := FindClosestContacts(k, req.Key)
        res.Nodes = ContactsToFoundNodes(closestContacts)
    }
    return nil
}

func FindValueLocally(k *Kademlia, Key ID) error {
    //1. Hash key
    hashed_key := HashKey(Key)
    hashed_id, err := FromByteArray(hashed_key)
    if err != nil {
        fmt.Printf("error hashing key\n")
    }
    //2. Find data corresponding to hashed key
    Val := k.Data[hashed_id]
    if Val == nil {
        fmt.Printf("ERR")
    } else {
        fmt.Printf("Val: %v\n", Val)
    }
    return nil
}

func CallFindValue(k *Kademlia, remoteContact *Contact, Key ID)(*FindValueResult, error){
    //Set up client.
    peer_str := HostPortToPeerStr(remoteContact.Host, remoteContact.Port)
    client, err := rpc.DialHTTP("tcp", peer_str)
    if err != nil {
        //maybe get rid of contact?
        log.Fatal("DialHttp: ", err)
    }
    //Create FindValueRequest
    hashed_key := HashKey(Key)
    hashed_id, err := FromByteArray(hashed_key)
    req := new(FindValueRequest)
    req.Sender = k.KContact
    req.MsgID = NewRandomID()
    req.Key = hashed_id

    //Call Kademlia.FindValue
    result := new(FindValueResult)
    err = client.Call("Kademlia.FindValue", req, &result)
    if err != nil {
          log.Fatal("Call: ", err)
    }
    //result either has value or Nodes
    return result, nil
}


