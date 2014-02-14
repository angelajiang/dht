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
                //Update shortlist
                //remove inactive contacts
                break Loop
            case result := <- main_chan:
                node_state[result.NodeID] = "active"
                    //contacts_to_rpc := NodesToRPC(node_state, result.Contacts)
                    //Update shortlist 
                    //Update closestNode
                    //check for exit conditions
                    //exit outside loop
                    //If none of the RPCed nodes updates closestNode : exit
            }
        }
    //make RPC calls again
    }

    //UpdateShortlist: Using responses from FindNode calls: update shortlist -> if contact is in active/inactive map, don't add it to shortlist
    //Update closestNode
    //Call General Update Function That We Haven't Done Before
    //Send FindNode RPCs again until: -- none of the new contacts are closer (i.e. closestNode doesn't change)  -- there are k active "already been queried" contacts in the shortlist
    return nil;
}

func UpdateShortlist(shortlist []Contact, rpc_contact Contact, alpha_contacts[]Contact, node_state map[ID]string) error {
    //remove inactive contact from shortlist
    if node_state[rpc_contact.NodeID]=="inactive" {
        //remove from shortlist
    }
    
    //add new alpha contacts to shortlist
    //make sure they aren't duplicated though or "inactive"
    return nil
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
