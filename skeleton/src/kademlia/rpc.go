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
    node_state := make(map[ID]string)

for {
    //Check updated closestContact boolean or   
    rpc_nodes := NodesToRPC(node_state, shortlist)
    main_chan := make(chan IDandContacts)
    timer_chan := make(chan bool)

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

    Loop:
    for {
        select {
            case <- timer_chan:
                //remove inactive contacts
                shortlist = FindAndRemoveInactiveContacts(shortlist, node_state)
                //Update shortlist
                response := <-  main_chan
                alpha_contacts := NodesToRPC(node_state, response.Contacts)
                shortlist = UpdateShortlist(shortlist, alpha_contacts, req.NodeID, node_state)
                break Loop
            case response := <- main_chan:
                node_state[response.NodeID] = "active"
                alpha_contacts := NodesToRPC(node_state, response.Contacts)
                shortlist = UpdateShortlist(shortlist, alpha_contacts, req.NodeID, node_state)
                if closestNode.NodeID == shortlist[0].NodeID {
                    break
                }
                closestNode = shortlist[0]    
                //Update closestNode
                    //check for exit conditions
                    //exit outside loop
                    //If none of the RPCed nodes updates closestNode : exit
            }
        }
    //make RPC calls again
    }

    res.MsgID = ds.NodeID
    res.Nodes = ContactsToFoundNodes(shortlist)
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

func UpdateShortlist(shortlist []Contact, alpha_contacts[]Contact, dest_id ID, node_state map[ID]string) (new_shortlist []Contact) {
    //add new alpha contacts to shortlist
    //make sure they aren't duplicated though or "inactive"
        for _, alpha_contact := range alpha_contacts {
            for _, cur_contact := range new_shortlist {
                if alpha_contact.NodeID == cur_contact.NodeID ||
                node_state[alpha_contact.NodeID] == "inactive" {
                    continue
                } else {
                    new_shortlist = append(new_shortlist, alpha_contact)
                    break
                }
            }
        }
        //sort new_shortlist
        ds := new(IDandContacts)
        ds.Contacts = new_shortlist
        ds.NodeID = dest_id
        sort.Sort(ds)
        new_shortlist = ds.Contacts
        return new_shortlist
}

func NodesToRPC(node_state map[ID]string, nodes []Contact)(nodes_to_call_rpc_on []Contact) {
//Takes a map of the "active/inactive" contacts and a list of contacts
//returns a list of the contacts we didn't query before that we should make RPC calls to
return nil
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
