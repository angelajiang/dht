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

func (k *Kademlia) FindNode(req FindNodeRequest, res *FindNodeResult) error {
    //check if we're the node in question
    if req.NodeID == k.NodeID{
        res.Nodes[0].NodeID = k.NodeID
        res.Nodes[0].IPAddr = k.Host.String()
        res.Nodes[0].Port = k.Port
    }else{
        //closestContacts = FindClosestContacts
        //res.Nodes = ContactsToFoundNodes
    }
    res.MsgID = CopyID(req.MsgID)
    return nil
}

func FindClosestContacts(k *Kademlia, requestID ID) (closestContacts []Contact){
    //Basically a wrapper to find closest contacts
    closestContacts = make([]Contact, 0) //Need to make specific to alpha parameter
    distance := k.NodeID.Xor(requestID)
    //First 1 from MSB is index to closest bucket
    indices := DistanceToOnes(distance)
    for index := range indices {
        //Add contacts from buckets[index] until closestContacts is full
        should_continue := AddNodesFromBucket(k, index, requestID, closestContacts)
        if !should_continue {
            //closestContacts is full
            break
        }
        //Otherwise, move on to next closest bucket
    }
    return
}

func AddNodesFromBucket(k *Kademlia, index int, requestID ID, closestContacts []Contact)(IsFull bool){
    //Sort contacts in bucket
    IsFull = false
    //make a sorted contacts slice
    sorted_contacts := make([]Contact, 3)
    //add an inital contact in there
    sorted_contacts = append(sorted_contacts, k.Buckets[index].Contacts[0])
    //loop through every contact in the bucket
    for _, contact := range k.Buckets[index].Contacts {
        //compare against distances of sorted contacts 
        for n, sorted_contact := range sorted_contacts {
           if PrefixDistance(contact.NodeID, k.NodeID) < PrefixDistance(sorted_contact.NodeID, k.NodeID) {
               if n == 0 {
                   //If it's closer than the first contact in sorted_contacts:
                   //Prepend to list (must do this in a stupid way because
                   //adding a contact to the front of a list is apparently not a
                   //straightforward process!
                   closestContacts = append(closestContacts, contact)
                   closestContacts = append(closestContacts, sorted_contacts[(n+1):]...)
                   //check if closestContacts is full
                   if len(closestContacts) == 3 {
                        IsFull = true
                        return
                   }
               } else {
                   //If it's closer than one of the other elements in
                   //sorted_contacts:
                   closestContacts = append(sorted_contacts[:(n-1)], contact, sorted_contacts[(n)])
                   //check again for length of closestContact
                   if len(closestContacts) == 3 {
                        IsFull = true
                        return
                   }
               }
           }
        }
    }
    //we check here in case we've looped through all the contacts and the loop
    //terminated
    if len(closestContacts) == 3 {
        IsFull = true
    }
    return
}


func PrefixDistance(id ID, other ID) (dist int) {
    return
}
func DistanceToOnes(distance ID)(ones []int){
    //Xor'ed result distance -> slice of bits in distance that are one from MSB->LSB 
    //ex) 1011 -> [1, 3, 4]
    ones = make([]int, 160)

    return
}

func ContactsToFoundNodes(contacts []Contact)(foundNodes []FoundNode){
    //Takes a splice of contacts and transforms it into a splice of foundNodes
    //Output can be stored in a FindNodeResult
    foundNodes = make([]FoundNode, 0)
    return
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
    if val,ok := k.Data[req.Key]; ok {
        res.Value = val

    }else{
        res.Value = nil
        //call find node
    }
    return nil
}

