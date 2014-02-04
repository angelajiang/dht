package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import (
    "net"
    "strconv"
    "strings"
    "crypto/sha1"
)

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

func InBucket(contact *Contact, bucket Bucket) (in_bucket bool, index int) {
    /*Returns true if contact is in contact list of bucket*/
    in_bucket = false
    for i,cur_contact := range bucket.Contacts {
        index = i
        if contact.NodeID == cur_contact.NodeID{
            in_bucket = true
            return
        }
    }
    return
}

func IsFull(bucket Bucket) bool {
    /*Returns true if bucket is full*/
    if len(bucket.Contacts) == cap(bucket.Contacts){
        return true
    }
    return false
}
