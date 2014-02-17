package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"net/rpc"
	"log"
)

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
    return nil
}


func CallStore(remote_contact *Contact, Key ID, Value []byte) error {
    //initialize request and result structs
    request := new(StoreRequest)
    var store_result StoreResult

    //set up rpc dial and all that jazz 
    peer_str := HostPortToPeerStr(remote_contact.Host, remote_contact.Port)
    client, err := rpc.DialHTTP("tcp", peer_str)
    if err != nil {
        log.Fatal("DialHttp: ", err)
    }

    hashed_key := HashKey(Key)
    hashed_id, err := FromByteArray(hashed_key)

    //set up request struct
    request.Sender = *(remote_contact)
    request.MsgID = NewRandomID()
//    request.Key = Key
    request.Key = hashed_id
    request.Value = Value

    //make rpc call 
    err = client.Call("Kademlia.Store", request, &store_result)
    if err != nil {
          log.Fatal("Call: ", err)
    }

    return nil
}

/*
func IterativeStore(k *Kademlia, Key ID, Value []byte){
	nodes, err := IterativeFindNode(k, Key)


}
*/
