package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

// Core Kademlia type. You can put whatever state you want in this.

import (
    "net"
    "net/rpc"
    "strconv"
    "strings"
    "log"
    "fmt"
    "crypto/sha1"
)

const NUMBUCKETS int =  160
const NUMCONTACTS int = 1
const VALUESIZE int = 160

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

func CallFindNode(k *Kademlia, remoteContact *Contact, search_key ID) (close_contacts []FoundNode, err error){
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
    req.NodeID = search_key
    err = client.Call("Kademlia.FindNode", req, &res) 
    if err != nil {
          log.Fatal("Call FindNode: ", err)
    }

   return res.Nodes, nil
}
/*HELPERS*/

func PeerStrToHostPort(listen_str string) (net.IP, uint16){
    /*Parsing*/
    input_arr := strings.Split(listen_str, ":")
    host_str := input_arr[0]
    port_str := input_arr[1]
    //Check if localhost
    if host_str == "localhost"{
        host_str = "127.0.0.1"
    }
    listen_netip := net.ParseIP(host_str)
    peer_uint64, _ := strconv.ParseUint(port_str, 10, 16)
    peer_uint16 := uint16(peer_uint64)

    return listen_netip, peer_uint16
}

func HostPortToPeerStr(remote_host net.IP, port uint16) (peer_str string){
    remote_host_str := remote_host.String()
    port_uint64 := uint64(port)
    port_str :=  strconv.FormatUint(port_uint64, 10)
    peer_str = remote_host_str + ":" + port_str
    return peer_str
}

func HashKey(key ID) []byte {
    //fmt.Printf("size of key: %v\n", len(key))
    h := sha1.New()
    h.Write(key[:])
    bs := h.Sum([]byte{})
    //fmt.Printf("bs is :%v\n", bs)
    return bs
}


