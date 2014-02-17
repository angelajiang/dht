package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
    "fmt"
    "time"
    "errors"
    "log"
    "net/rpc"
)

// Datatypes
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

//RPC
func (k *Kademlia) FindNode(req FindNodeRequest, res *FindNodeResult) error {
    //check if we're the node in question
    if req.NodeID == k.NodeID {
        foundContact := ContactToFoundNode(k.KContact)
        res.Nodes = append(res.Nodes, *foundContact)
    } else {
        closestContacts := FindClosestContacts(k, req.NodeID)
        res.Nodes = ContactsToFoundNodes(closestContacts)
    }
    res.MsgID = CopyID(req.MsgID)
    return nil
}

func CallFindNode(k *Kademlia, remoteContact *Contact, search_id ID) (close_contacts []FoundNode, err error){
   //set up client 
    peer_str := HostPortToPeerStr(remoteContact.Host, remoteContact.Port)
    client, err := rpc.DialHTTP("tcp", peer_str)
    if err != nil {
          log.Fatal("DialHTTP in FindNode: ", err)
    }
    fmt.Printf("Client in CallFindNode\n")

    req := new(FindNodeRequest)
    var res FindNodeResult
    req.Sender.NodeID = k.NodeID
    req.Sender.Host = k.Host
    req.Sender.Port = k.Port
    req.MsgID = NewRandomID()
    req.NodeID = search_id
    err = client.Call("Kademlia.FindNode", req, &res) 
    if err != nil {
          log.Fatal("Error in CallFindNode: ", err)
    }

   return res.Nodes, nil
}


//Iterative
func IterativeFindNode(k *Kademlia, destID ID) (closestContacts []Contact, err error) {
    
    //1. FindClosestContacts -> this returns 3 closest nodes.
    closestContacts = FindClosestContacts(k, destID)
    //2. Make a sorted shortlist, add initial closest contacts to it. Set initial value of closestNode = closest contact in shortlist.
    if len(closestContacts) == 0{
        err = errors.New("Error in IterativeFindNode: No contacts to send initial RPCs to.")
        return
    }
    shortlist := SortContacts(closestContacts, destID)
    closestNode := shortlist[0]
    updatedClosestContact := true

    node_state := make(map[ID]string)

    main_chan := make(chan IDandContacts)
    timer_chan := make(chan bool)


    Shortlist_Loop:
    for {
        //Exit conditions
        if updatedClosestContact == false{
            fmt.Printf("IFN: exit condition 1\n")
            break Shortlist_Loop
        }
        rpc_nodes := GetAlphaNodesToRPC(shortlist, node_state)
        if len(rpc_nodes) == 0{
            fmt.Printf("IFN: exit condition 2\n")
            //No more nodes in shortlist to query
            break Shortlist_Loop
        }

        //Did not exit. Start another iteration of sending alpha RPCs
        updatedClosestContact = false
        for _, c := range rpc_nodes {
            fmt.Printf("Sending RPC to %v\n", c.NodeID)
            cur := c
            go func() {
                node_state[cur.NodeID] = "inactive"
                main_chan <- FindNodeWithChannel(k, &cur, destID)
            }()
        }
        go func() {
            time.Sleep(1000 * time.Millisecond)
            timer_chan <- true
        }()

        Alpha_Loop:
        for {
            //Wait for reponses from RPC calls and timer
            select {
                case <- timer_chan:
                    //remove contacts that did not respond to RPC from shortlist
                    shortlist = RemoveInactiveContacts(shortlist, node_state)
                    shortlist = SortContacts(shortlist, destID)

                    break Alpha_Loop
                case response := <- main_chan:
                    //UPDATE
                    fmt.Printf("%v responds to RPC with %v\n", response.NodeID, FirstBytesOfContactIDs(response.Contacts))
                    node_state[response.NodeID] = "active"
                    shortlist = UpdateShortlist(shortlist, response.Contacts, destID, node_state)
                    fmt.Printf("shortlist after adding: %v\n", FirstBytesOfContactIDs(shortlist))
                    if len(shortlist) > 0 {
                        fmt.Printf("closest: %v\nShortlist[0]%v\n", closestNode.NodeID[0], shortlist[0].NodeID[0])
                        if closestNode.NodeID != shortlist[0].NodeID {
                            closestNode = shortlist[0]    
                            updatedClosestContact = true
                        }
                    }
                }
            }
    }

    shortlist = RemoveInactiveContacts(shortlist, node_state)
    shortlist = RemoveNodesToRPC(shortlist, node_state)
    fmt.Printf("shortlist returned: %v\n", FirstBytesOfContactIDs(shortlist))
    closestContacts = GetFirstAlphaContacts(shortlist)

    //UpdateShortlist: Using responses from FindNode calls: update shortlist -> if contact is in active/inactive map, don't add it to shortlist
    //Update closestNode
    //Call General Update Function That We Haven't Done Before
    //Send FindNode RPCs again until: -- none of the new contacts are closer (i.e. closestNode doesn't change)  -- there are k active "already been queried" contacts in the shortlist
    
    return;
}



