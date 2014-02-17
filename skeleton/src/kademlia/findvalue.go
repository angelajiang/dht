package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
    "net"
    "fmt"
    "net/rpc"
    "log"
    "errors"
    "time"
)

// Host identification.
type Contact struct {
    NodeID ID
    Host net.IP
    Port uint16
}

// FIND_VALUE
type FindValueRequest struct {
    Sender Contact
    MsgID ID
    Key ID
}

type IterativeFindValueResult struct{
    NodeID ID
    Contacts []Contact
    Value []byte
    Err error
}

// If Value is nil, it should be ignored, and Nodes means the same as in a
// FindNodeResult.
type FindValueResult struct {
    MsgID ID
    Value []byte
    Nodes []FoundNode
    Err error
}

func FindValueWithChannel(k *Kademlia, remoteContact *Contact, search_id ID) (*IterativeFindValueResult) {
    res , err := CallFindValue(k, remoteContact, search_id)
    FoundContacts := FoundNodesToContacts(res.Nodes)
    ret := new(IterativeFindValueResult)
    ret.NodeID = remoteContact.NodeID
    ret.Contacts = FoundContacts
    ret.Value = res.Value
    ret.Err = err
    return ret
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

func FindValueLocally(k *Kademlia, Key ID) error {
    //1. Hash key
    hashed_key := HashKey(Key)
    hashed_id, err := FromByteArray(hashed_key)
    if err != nil {
        fmt.Printf("error hashing key\n")
    }
    //2. Find data corresponding to hashed key
    Val := k.Data[hashed_id]
    if Val == nil {
        fmt.Printf("Error in FindValueLocally\n")
    } else {
        fmt.Printf("Val: %v\n", Val)
    }
    return nil
}

func CallFindValue(k *Kademlia, remoteContact *Contact, Key ID)(*FindValueResult, error){
    //Set up client.
    peer_str := HostPortToPeerStr(remoteContact.Host, remoteContact.Port)
    client, err := rpc.DialHTTP("tcp", peer_str)
    if err != nil {
        //maybe get rid of contact?
        log.Fatal("DialHttp: ", err)
    }
    //Create FindValueRequest
    hashed_key := HashKey(Key)
    hashed_id, err := FromByteArray(hashed_key)
    req := new(FindValueRequest)
    req.Sender = k.GetContact()
    req.MsgID = NewRandomID()
    req.Key = hashed_id

    //Call Kademlia.FindValue
    result := new(FindValueResult)
    err = client.Call("Kademlia.FindValue", req, &result)
    if err != nil {
          log.Fatal("Call: ", err)
    }
    //result either has value or Nodes
    return result, nil
}

func IterativeFindValue(k *Kademlia, key ID) (retID ID, foundValue []byte, e error){
    
    e = errors.New("ERR")

    closestContacts := FindClosestContacts(k, key)
    if len(closestContacts) == 0{
        return
    }
    shortlist := SortContacts(closestContacts, key)
    closestNode := shortlist[0]
    updatedClosestContact := true

    node_state := make(map[ID]string)

    main_chan := make(chan *IterativeFindValueResult)
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
                main_chan <- FindValueWithChannel(k, &cur, key)
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
                    shortlist = SortContacts(shortlist, key)
                    break Alpha_Loop
                case response := <- main_chan:
                    //UPDATE
                    //TODO: do something with response.Err?
                    for _,c := range response.Contacts{
                        cur := c
                        Update(k, &cur)
                    }
                    if response.Value != nil{
                        fmt.Printf("Found value %v from node %v\n", response.Value, response.NodeID)
                        foundValue = response.Value
                        retID = response.NodeID
                        e = nil
                        return
                    }
                    fmt.Printf("%v responds to RPC with %v\n", response.NodeID[0], FirstBytesOfContactIDs(response.Contacts))
                    node_state[response.NodeID] = "active"
                    shortlist = UpdateShortlist(shortlist, response.Contacts, key, node_state)
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
    
    return;
}



