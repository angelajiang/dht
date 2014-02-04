package kademlia

import (
	"log"
	"net/rpc"
	"net"
	"fmt"

)

func TestPingFirstPeer(k *Kademlia, first_peer_str string){
    // Confirm our server is up with a PING request and then exit.
    client, err := rpc.DialHTTP("tcp", first_peer_str)
    if err != nil {
        log.Fatal("DialHTTP: ", err)
    }
    ping := new(Ping)
    ping.MsgID = NewRandomID()
    var pong Pong
    err = client.Call("Kademlia.Ping", ping, &pong)
    if err != nil {
        log.Fatal("Call: ", err)
    }

    log.Printf("ping msgID: %s\n", ping.MsgID.AsString())
    log.Printf("pong msgID: %s\n", pong.MsgID.AsString())
}

func TestUpdate(k *Kademlia, c Contact){
    //Making new contacts and calling Update
    Update(c, &k.Buckets[0])
}
//func TestFindValue(k *Kademlia, key ID, contact_with_value Contact, contact_without_value Contact){
func TestFindValue(k *Kademlia, contact_with_value Contact, contact_without_value Contact){
    //Putting value into tmp_contact for testing DoFindValue
    s := make([]byte, 5)
    data_key := NewRandomID()
    tmp_data := s
    fmt.Printf("tmpdata: %v\n", tmp_data)
    k.Data[data_key] = tmp_data
    res, _ := CallFindValue(k, &contact_with_value, data_key)
    fmt.Printf("DoFindValue: contact with value response: %v\n",res)
    res, _ = CallFindValue(k, &contact_without_value, data_key)
    fmt.Printf("DoFindValue: contact without value response: %v\n",res)
}

func TestBasicRPCs(k *Kademlia, first_peer_str string){
	TestPingFirstPeer(k, first_peer_str)
    tmp_id := NewRandomID()
    tmp_ip := net.ParseIP("127.0.0.1")
    tmp_port := k.Port
    contact1 := Contact{tmp_id, tmp_ip, tmp_port}
	TestUpdate(k, contact1)

    tmp_id = NewRandomID()
    tmp_ip = net.ParseIP("127.0.0.1")
    tmp_port = k.Port
    contact2 := Contact{tmp_id, tmp_ip, tmp_port}
	TestFindValue(k, contact1, contact2)

}