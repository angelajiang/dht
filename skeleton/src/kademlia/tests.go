package kademlia

import (
	"log"
	"net/rpc"
	"net"
	"fmt"

)

func TestPingFirstPeer(k *Kademlia, first_peer_str string){
    // Confirm our server is up with a PING request and then exit.
	fmt.Printf("\nTESTING: Pinging first peer\n")
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

func TestUpdate(k *Kademlia){
    //Making new contacts and calling Update
	fmt.Printf("\nTESTING: Update\n")
    tmp_id := NewRandomID()
    tmp_ip := net.ParseIP("127.0.0.1")
    tmp_port := k.Port
    contact1 := Contact{tmp_id, tmp_ip, tmp_port}
    Update(&contact1, &k.Buckets[0])
}

func TestStoreAndFindValue(k *Kademlia){
	fmt.Printf("\nTESTING: Store and Find Value\n")
    data_key := NewRandomID()
    s := make([]byte, 5)
    tmp_data := s
    fmt.Printf("Data to store: %v\n", tmp_data)

   	//Try to find value before store 
	dest_contact := k.KContact
    res, err := CallFindValue(k, &k.KContact, data_key)
    if err != nil{
    	log.Fatal("CallFindValue: ", err)
    }
    fmt.Printf("CallFindValue before store: %v\n",res.Value)

    //Call store
    err = CallStore(&dest_contact, data_key, tmp_data)
    if err != nil{
    	log.Fatal("CallStore: ", err)
    }

    //Try to find value after store
    res, err = CallFindValue(k, &dest_contact, data_key)
    if err != nil{
    	log.Fatal("CallFindValue: ", err)
    }
    fmt.Printf("CallFindValue after store: %v\n",res.Value)

}

func TestGetSetBits(){
	id1 := NewRandomID()
	id2 := NewRandomID()
	distance := id1.Xor(id2) 
	fmt.Printf("Distance: %v\n", distance)
	indices := GetSetBits(distance)
	fmt.Printf("Indices: %v\n", indices)

}

func TestContactsToFoundNodes(k *Kademlia){
	closestContacts := make([]Contact, 0)
	FillTestContactSlice(&closestContacts)
    fmt.Printf("closestContacts: %v\n", closestContacts)
    foundNodes := ContactsToFoundNodes(closestContacts)
    fmt.Printf("foundNodes: %v\n", foundNodes)
}

func FillTestContactSlice(contacts *[]Contact){
    contact1 := NewRandomContact()
    contact2 := NewRandomContact()
    *contacts = append(*contacts, contact1)
    *contacts = append(*contacts, contact2)
}
/*
func TestContactsToFoundNodes(k *Kademlia){
	fmt.Println("TestContactsToFoundNodes:\n")
	contacts := make([]Contact, 0)
	FillTestContactSlice(&contacts)
    fmt.Printf("ContactsToSort: %v\n", contacts)
}
*/




func TestBasicRPCs(k *Kademlia, first_peer_str string){
	TestUpdate(k)
	TestStoreAndFindValue(k)
	TestGetSetBits()
	TestContactsToFoundNodes(k)

}