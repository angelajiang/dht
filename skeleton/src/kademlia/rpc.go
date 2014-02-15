package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
    "net"
    "fmt"
    "sort"
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
    ds := new(IDandContacts)
    ds.Contacts = closestContacts
    ds.NodeID = req.NodeID
    sort.Sort(ds)
    shortlist := ds.Contacts
    closestNode := shortlist[0]
    updatedClosestContact := true

    node_state := make(map[ID]string)

    main_chan := make(chan IDandContacts)
    timer_chan := make(chan bool)


Shortlist_Loop:
for {
    //Exit conditions
    if updatedClosestContact == false{
        //RPC calls did not return any contacts closer than current closest node
        break Shortlist_Loop
    }
    rpc_nodes := GetAlphaNodesToRPC(shortlist, node_state)
    if len(rpc_nodes) == 0{
        //No more nodes in shortlist to query
        break Shortlist_Loop
    }

    //Did not exit. Start another iteration of sending alpha RPCs
    updatedClosestContact = false
    for _, c := range rpc_nodes {
        go func() {
            node_state[c.NodeID] = "inactive"
            main_chan <- FindNodeWithChannel(k, &c, req.NodeID)
        }()
    }
    go func() {
        time.Sleep(300 * time.Millisecond)
        timer_chan <- true
    }()

    Alpha_Loop:
    for {
        //Wait for reponses from RPC calls and timer
        select {
            case <- timer_chan:
                //remove contacts that did not respond to RPC from shortlist
                shortlist = FindAndRemoveInactiveContacts(shortlist, node_state)

                //Shouldn't care about results in main_chan after the timer runs out

                //Update shortlist
                //response := <-  main_chan
                //alpha_contacts := NodesToRPC(node_state, response.Contacts)
                //shortlist = UpdateShortlist(shortlist, alpha_contacts, req.NodeID, node_state)
                break Alpha_Loop
            case response := <- main_chan:
                node_state[response.NodeID] = "active"
                //alpha_contacts := NodesToRPC(node_state, response.Contacts)
                shortlist = UpdateShortlist(shortlist, response.Contacts, req.NodeID, node_state)
                if closestNode.NodeID != shortlist[0].NodeID {
                    closestNode = shortlist[0]    
                    updatedClosestContact = true
                }
            }
        }
    //make RPC calls again
    }
    shortlist = FindAndRemoveInactiveContacts(shortlist, node_state)
    shortlist = RemoveNodesToRPC(shortlist, node_state)
    res.Nodes = ContactsToFoundNodes(shortlist)
    res.MsgID = ds.NodeID
    res.Err = nil

    //UpdateShortlist: Using responses from FindNode calls: update shortlist -> if contact is in active/inactive map, don't add it to shortlist
    //Update closestNode
    //Call General Update Function That We Haven't Done Before
    //Send FindNode RPCs again until: -- none of the new contacts are closer (i.e. closestNode doesn't change)  -- there are k active "already been queried" contacts in the shortlist
    
    return nil;
    
}

func FindAndRemoveInactiveContacts(shortlist []Contact, node_state map[ID]string) (new_shortlist []Contact) {
    for _, contact := range shortlist {
        if node_state[contact.NodeID] == "inactive" {
            new_shortlist = RemoveInactiveContact(shortlist, contact, node_state)
        }
    }
    return
}

func RemoveInactiveContact(shortlist []Contact, contact Contact, node_state map[ID]string) (new_shortlist []Contact) {
    //remove inactive rpc contact from shortlist
    if node_state[contact.NodeID]=="inactive" {
        //remove from shortlist
        var index int
        for i, cur_contact := range shortlist {
            if cur_contact.NodeID == contact.NodeID {
                index = i
                break
            }
        }
        new_shortlist = append(new_shortlist, shortlist[0:index]...)
        new_shortlist = append(new_shortlist, shortlist[index+1:]...)
    }
    return new_shortlist
}

func UpdateShortlist(shortlist []Contact, alpha_contacts[]Contact, dest_id ID, node_state map[ID]string) []Contact {
    //add new alpha contacts to shortlist
    //make sure they aren't duplicated though or "inactive"
    for _, alpha_contact := range alpha_contacts {
        for _, cur_contact := range shortlist {
            if alpha_contact.NodeID == cur_contact.NodeID ||
            node_state[alpha_contact.NodeID] == "inactive" {
                continue
            } else {
                shortlist = append(shortlist, alpha_contact)
                break
            }
        }
    }
    //sort new_shortlist
    ds := new(IDandContacts)
    ds.Contacts = shortlist
    ds.NodeID = dest_id
    sort.Sort(ds)
    shortlist = ds.Contacts
    return shortlist
}

func GetAlphaNodesToRPC(nodes []Contact, node_state map[ID]string) (alpha_contacts_to_rpc []Contact) {
//Takes a map of the "active/inactive" contacts and a list of contacts
//returns a list of alpha contacts we didn't query before that we should make RPC calls to
    alpha_contacts_to_rpc = make([]Contact, 0, ALPHA)
    for _, c := range nodes{
        if _,ok := node_state[c.NodeID]; ok {
            //If it's in node_state, then we've already sent an rpc
            continue
        } else{
            alpha_contacts_to_rpc = append(alpha_contacts_to_rpc, c)
            if len(alpha_contacts_to_rpc) == cap(alpha_contacts_to_rpc){
                return
            }
        }
    }
    return
}

func RemoveNodesToRPC(shortlist []Contact, node_state map[ID]string) []Contact{
    //Update shortlist to only include nodes that we've sent RPCs for
    new_shortlist := make([]Contact,0)
    for _, c := range shortlist{
        if _,ok := node_state[c.NodeID]; ok {
            //If it's in node_state, then we've already sent an rpc
            new_shortlist = append(new_shortlist, c)
        }
    }
    return new_shortlist
}

func FindNodeWithChannel(k *Kademlia, remoteContact *Contact, search_id ID) (ret IDandContacts) {
    FoundNodes, err := CallFindNode(k, remoteContact, search_id)
    if err != nil {
        fmt.Printf("Error calling FindNode RPC")
    }
    FoundContacts := FoundNodesToContacts(FoundNodes)
    ret.NodeID = remoteContact.NodeID
    ret.Contacts = FoundContacts
    return
}
