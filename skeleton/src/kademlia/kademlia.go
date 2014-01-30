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
)

const NUMBUCKETS int =  160
const NUMCONTACTS int = 1

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
    c := new(Contact)
    c.NodeID = kptr.NodeID
    c.Host = host
    c.Port = port
    kptr.KContact = *c
    return kptr
}

func DoPing(remote_host net.IP, port uint16) (Pong, error){
    peer_str := HostPortToPeerStr(remote_host, port)
    client, err := rpc.DialHTTP("tcp", peer_str)
    if err != nil {
          log.Fatal("Call: ", err)
    }
    ping := new(Ping)
    ping.MsgID = NewRandomID()
    var pong Pong
    err = client.Call("Kademlia.Ping", ping, &pong)
    if err != nil {
          log.Fatal("Call: ", err)
    }

    return pong, nil
}

func DoFindValue(k *Kademlia, remoteContact *Contact, Key ID)(*FindValueResult, error){
    //Set up client.
    peer_str := HostPortToPeerStr(remoteContact.Host, remoteContact.Port)
    client, err := rpc.DialHTTP("tcp", peer_str)
    if err != nil {
        //maybe get rid of contact?
        log.Fatal("DialHttp: ", err)
    }
    fmt.Printf("Client in DoFindValue: %v\n", client)
    //Create FindValueRequest
    req := new(FindValueRequest)
    req.Sender = k.KContact         //Sender of the request?
    req.MsgID = NewRandomID()
    req.Key = Key

    fmt.Printf("req in DoFindValue: %v\n", req)

    //Call Kademlia.FindValue
    result := new(FindValueResult)
    err = client.Call("Kademlia.FindValue", req, &result)
    if err != nil {
          log.Fatal("Call: ", err)
    }
    //If value is there, return data
    //else, call DoFindNode
    return result, nil
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



