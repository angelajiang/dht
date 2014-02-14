package kademlia

import (
	"log"
	"net/rpc"
	"fmt"
	"sort"

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

func TestUpdate(k *Kademlia, n int){
    //Making new contacts and calling Update
	fmt.Printf("\nTESTING: Update\n")
	for i := 0; i < n; i++ {
		c := NewRandomContact()
	    Update(k, c)
	}
	fmt.Printf("After adding %v contacts:\n%v\n", n, k.Buckets)
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
    fmt.Printf("Nodes before store: %v\n",res.Nodes)

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
    fmt.Printf("Nodes after store: %v\n",res.Nodes)

}

func TestGetSetBits(){
	fmt.Printf("\nTESTING: GetSetBits\n")
	id1 := NewRandomID()
	id2 := NewRandomID()
	distance := id1.Xor(id2) 
	fmt.Printf("Distance: %v\n", distance)
	indices := GetSetBits(distance)
	fmt.Printf("Indices: %v\n", indices)

}

func TestContactsToFoundNodes(k *Kademlia){
	fmt.Printf("\nTESTING: ContactsToFoundNodes\n")
	closestContacts := make([]Contact, 0)
	FillTestContactSlice(&closestContacts, 3)
    fmt.Printf("closestContacts: %v\n", closestContacts)
    foundNodes := ContactsToFoundNodes(closestContacts)
    fmt.Printf("foundNodes: %v\n", foundNodes)
}

func FillTestContactSlice(contacts *[]Contact, n int){
	for i := 0; i < n; i++ {
		c := NewRandomContact()
	    *contacts = append(*contacts, *c)
	}
}

func TestSortByDistance(){
	fmt.Println("\nTESTING: TestSortByDistance\n")
	contacts := make([]Contact, 0)
	FillTestContactSlice(&contacts, 3)
    fmt.Printf("ContactsToSort: %v\n\n", contacts)
    ds := new(IDandContacts)
    ds.Contacts = contacts
    dest_id := NewRandomID()
    ds.NodeID = dest_id
    sort.Sort(ds)
    fmt.Printf("Sorted: %v\n\n", ds.Contacts)
}

func TestFindClosestContacts(k *Kademlia){
	fmt.Println("\nTESTING: TestFindClosestContacts\n")
	TestUpdate(k, 100)		
	requestID := NewRandomID()
	//FindClosestContacts(k, requestID)
	closestContacts := FindClosestContacts(k, requestID)
	fmt.Printf("ClosestContacts: %v\n", closestContacts)
}

func TestFindNode(k *Kademlia){

	fmt.Println("\nTESTING: TestFindNode\n")
	c := NewRandomContact()
	c.Port = 7777
	closestContacts, _ := CallFindNode(k, c, NewRandomID())
	fmt.Printf("ClosestContacts: %v\n", closestContacts)




}


func TestBasicRPCs(k *Kademlia, first_peer_str string){
	TestUpdate(k, 100)
	//TestStoreAndFindValue(k)
	//TestGetSetBits()
	//TestContactsToFoundNodes(k)
	//TestSortByDistance()
	//TestFindClosestContacts(k)
	TestFindNode(k)
	fmt.Printf("\n\n")

}
