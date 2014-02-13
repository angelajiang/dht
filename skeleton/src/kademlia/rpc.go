package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
    "net"
    //"net/rpc"
    "fmt"
    "sort"
    //"log"
    "time"
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
    if req.NodeID == k.NodeID {
        res.Nodes[0].NodeID = k.NodeID
        res.Nodes[0].IPAddr = k.Host.String()
        res.Nodes[0].Port = k.Port
    } else {
        closestContacts := FindClosestContacts(k, req.NodeID)
        res.Nodes = ContactsToFoundNodes(closestContacts)
    }
    res.MsgID = CopyID(req.MsgID)
    return nil
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

func  (k *Kademlia) IterativeFindNode(req FindNodeRequest, res *FindNodeResult) error {
    //1. FindClosestContacts -> this returns 3 closest nodes.
    closestContacts := FindClosestContacts(k, req.NodeID)
    //2. Make a sorted shortlist, add initial closest contacts to it. Set initial value of closestNode = closest contact in shortlist.
    ds := new(DistanceSorter)
    ds.Contacts = closestContacts
    ds.DestID = req.NodeID
    sort.Sort(ds)
    shortlist := ds.Contacts
    closestNode := shortlist[0]

    //3. NodesToRPC function that takes a shortlist and returns up to alpha
    //nodes that we need to contact (checks if node is marked active doesn't add it to list)
    
    //4. Send parallel FindNode RPC calls to contacts returned from NodesToRPC. *** CHANNELS GO HERE
    rpc1 := make(chan []Contact)
    rpc2 := make(chan []Contact)
    rpc3 := make(chan []Contact)
    //the third argument here would be whatever NodesToRPC returns
    go FindNodeWithChannel(k, rpc1, &shortlist[0], req.NodeID)
    go FindNodeWithChannel(k, rpc2, &shortlist[1], req.NodeID)
    go FindNodeWithChannel(k, rpc3, &shortlist[2], req.NodeID)
    
    response1 := false
    response2 := false
    response3 := false
    for {
        select {
        case res1 := <- rpc1:
            response1 = true
            //mark contact1 as active
            //Update shortlist
            //Update closestNode
            //check for exit conditions
        case res2 := <- rpc2:
            response2 = true
            //mark contact2 as active
            //Update shortlist
            //Update closestNode
            //check for exit conditions
        case res3 := <- rpc3:
            response3 = true
            //mark contact3 as active
            //Update shortlist 
            //Update closestNode
            //check for exit conditions
        case <- time.After(10 * 1e9): //timeout after 10 seconds
            if !response1 {
                //mark contact1 as inactive
                //remove it from shortlist
                //check for exit conditions
            }
            if !response2 {
                //mark contact2 as inactive
                //remove it from shortlist
                //check for exit conditions
            }
            if !response3 {
                //mark contact3 as inactive
                //remove it from shortlist
                //check for exit conditions
            }
        }
    
        //Make new RPC calls
    }
/*
Rula wrote this stuff:  
for {
    
rpcs_returned := 3
select {
    case <- chan:
        //add to short list
        //check for exit conditions
    }

    go MakeFindNodeCall()

*/

    //if contact responds: mark it as "active" -> map from contact node id to "active", "inactive"
    //5. UpdateShortlist: Using responses from FindNode calls: update shortlist -> if contact is in active/inactive map, don't add it to shortlist
    //6. Update closestNode
    //7. Call General Update Function That We Haven't Done Before
    //8. Send FindNode RPCs again until: -- none of the new contacts are closer (i.e. closestNode doesn't change)  -- there are k active "already been queried" contacts in the shortlist
    return nil;
}

func FindNodeWithChannel(k *Kademlia, c chan []Contact, remoteContact *Contact, search_id ID) error {
    FoundNodes, err := CallFindNode(k, remoteContact, search_id)
    FoundContacts := FoundNodesToContacts(FoundNodes)
    c <- FoundContacts
    return nil
}
