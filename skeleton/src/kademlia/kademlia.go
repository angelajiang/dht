package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

// Core Kademlia type. You can put whatever state you want in this.

import (
    "net"
    "net/rpc"
    "log"
    "fmt"
    "errors"
)

const NUMBUCKETS int =  160
const NUMCONTACTS int = 100
const VALUESIZE int = 160
const ALPHA int = 3

type Kademlia struct {
    NodeID ID
    Buckets []Bucket
    Host net.IP
    Port uint16
    Data map[ID][]byte
    KContact Contact
}

func NewKademlia(host net.IP, port uint16) *Kademlia {
    // TODO: Assign yourself a random ID and prepare other state here.
    kptr := new(Kademlia)
    kptr.NodeID = NewRandomID()
    kptr.Buckets = make([]Bucket, NUMBUCKETS)
    for i,_ := range kptr.Buckets{
        kptr.Buckets[i] = *(NewBucket())
    }
    kptr.Host = host
    kptr.Port = port
    kptr.Data = make(map[ID][]byte, VALUESIZE)
    c := new(Contact)
    c.NodeID = kptr.NodeID
    c.Host = host
    c.Port = port
    kptr.KContact = *c
    return kptr
}

func FindContactLocally(k *Kademlia, contact_id ID) error {
    dist := k.NodeID.Xor(contact_id)
    bucket_index := GetBucketIndex(dist) 
    for _, contact := range k.Buckets[bucket_index].Contacts {
        if contact.NodeID == contact_id {
            fmt.Printf("%v  %v\n", contact.Host, contact.Port)
            break
        }
    }
    fmt.Printf("ERR\n")
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

func CallPing(remote_host net.IP, port uint16) (Pong, error){
    /* DoPing should probably take a Kademlia object here */
    //TODO: run the Update function?
    peer_str := HostPortToPeerStr(remote_host, port)
    client, err := rpc.DialHTTP("tcp", peer_str)
    if err != nil {
          log.Fatal("Call: ", err)
    }
    ping := new(Ping)
    ping.MsgID = NewRandomID()
    //ping.Sender.NodeID = 

    var pong Pong
    err = client.Call("Kademlia.Ping", ping, &pong)
    if err != nil {
          log.Fatal("Call: ", err)
    }

    return pong, nil
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

func CallFindNode(k *Kademlia, remoteContact *Contact, search_id ID) (close_contacts []FoundNode, err error){
   //set up client 
    peer_str := HostPortToPeerStr(remoteContact.Host, remoteContact.Port)
    client, err := rpc.DialHTTP("tcp", peer_str)
    if err != nil {
          log.Fatal("DialHTTP in FindNode: ", err)
    }
    fmt.Printf("Client in CallFindNode: %v\n", client)

    req := new(FindNodeRequest)
    var res FindNodeResult
    req.Sender.NodeID = k.NodeID
    req.Sender.Host = k.Host
    req.Sender.Port = k.Port
    req.MsgID = NewRandomID()
    req.NodeID = search_id
    err = client.Call("Kademlia.FindNode", req, &res) 
    if err != nil {
          log.Fatal("Call FindNode: ", err)
    }

   return res.Nodes, nil
}


func Update(k *Kademlia, contact *Contact) error {
    //Choose correct bucket to put contact
    distance := k.NodeID.Xor(contact.NodeID)
    bucket_index := GetBucketIndex(distance)
    bucket_addr := &k.Buckets[bucket_index]

    bucket := *bucket_addr
    fmt.Printf("Bucket %v length before update: %v\n", bucket_index, len(k.Buckets[bucket_index].Contacts))
    in_bucket, index:= InBucket(contact, *bucket_addr)
    is_full := IsFull(bucket)
    switch {
    case in_bucket:
        /*Move contact to end of bucket's contact list*/
        fmt.Printf("Case: in_bucket\n")
        bucket.Contacts = append(bucket.Contacts[:index-1],bucket.Contacts[(index+1):]...)
        bucket.Contacts = append(bucket.Contacts, *contact)
    case !in_bucket && !is_full:
        if len(bucket_addr.Contacts) == 0{
            fmt.Printf("Case: !in_bucket, !is_full, empty\n")
            bucket_addr.Contacts = append(bucket_addr.Contacts, *contact)
        } else {
            fmt.Printf("Case: !in_bucket, !is_full, !empty\n")
            pong, err := CallPing(bucket_addr.Contacts[0].Host, k.Port)//bucket_addr.Contacts[0].Port)
            fmt.Printf("%+v\n", pong)
            if err != nil{
                bucket_addr.Contacts = append(bucket_addr.Contacts[1:], *contact)
            }
            bucket_addr.Contacts = append(bucket_addr.Contacts, *contact)
        }
    case !in_bucket && is_full:
        fmt.Printf("Case: !in_bucket and is_full\n")
        /*Replace head of list if head doesn't respond. Otherwise, ignore*/
        pong, err := CallPing(bucket_addr.Contacts[0].Host, k.Port)//bucket_addr.Contacts[0].Port)
        fmt.Printf("%+v\n", pong)
        if err != nil{
            //drop head append contact to end of list
            bucket_addr.Contacts = append(bucket_addr.Contacts[1:], *contact)
        } else {
            //Move head to tail
            bucket_addr.Contacts = append(bucket_addr.Contacts[1:],bucket_addr.Contacts[0])
        }
    }
    fmt.Printf("After: %v\n", len(k.Buckets[bucket_index].Contacts))
    return errors.New("function not implemented")
}

/*
func FindNodeFromNodeId(k *Kademlia, node_id ID) {
   //GetBucketIndex?    

}
*/
