package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"sort"
    "net"
    "fmt"
)

func FindClosestContacts(k *Kademlia, requestID ID) (closestContacts []Contact){
    //Basically a wrapper to find closest contacts
    closestContacts = make([]Contact, 0) //Need to make specific to alpha parameter
    distance := k.NodeID.Xor(requestID)
    //First 1 from MSB is index to closest bucket
    indices := GetSetBits(distance)
    for index := range indices {
        //Add contacts from buckets[index] until closestContacts is full
        is_full := AddNodesFromBucket(k, index, requestID, &closestContacts)
        if is_full {
            //closestContacts is full
            return closestContacts
        }
        //Otherwise, move on to next closest bucket
    }
    return
}

func AddNodesFromBucket(k *Kademlia, index int, requestID ID, closestContactsPtr *[]Contact)(IsFull bool){
    IsFull = false
    if (len(k.Buckets[index].Contacts) == 0){
        //No contacts in this bucket
        return
    }
    ds := new(IDandContacts)
    ds.Contacts = k.Buckets[index].Contacts
    ds.NodeID = requestID
    sort.Sort(ds)
    sorted_contacts := make([]Contact, 0)
    sorted_contacts = ds.Contacts       //Holds sorted contacts of bucket
    closestContacts := *closestContactsPtr
    for _, contact := range sorted_contacts {
	//Add contact from sorted list to closestContacts until full
        *closestContactsPtr = append(closestContacts, contact)
        if len(closestContacts) == ALPHA {
            IsFull = true
        	return
        }
    }
    //Added all contacts in bucket and still not full
    return
}

func ContactToFoundNode(c Contact) (f *FoundNode){
    f = new(FoundNode)
    f.NodeID = c.NodeID
    f.Port = c.Port
    f.IPAddr = c.Host.String()
    return
}

func ContactsToFoundNodes(contacts []Contact)(foundNodes []FoundNode){
    //Takes a slice of contacts and transforms it into a slice of foundNodes
    //Output can be stored in a FindNodeResult
    foundNodes = make([]FoundNode, 0)
    for _, contact := range contacts{
        f := ContactToFoundNode(contact)
        foundNodes = append(foundNodes, *f)
    }
    return
}

func FoundNodesToContacts(foundNodes []FoundNode) (contacts []Contact) {
    //Takes a slice of foundNodes and returns it as a slice of contacts
    contacts = make([]Contact, 0)
    for _, found_node := range foundNodes {
        c := new(Contact)
        c.NodeID = found_node.NodeID
        c.Port = found_node.Port
        c.Host = net.ParseIP(found_node.IPAddr)
        contacts = append(contacts, *c)
    }
    return contacts
}

//ITERATIVE HELPERS//

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

func GetFirstAlphaContacts(contacts []Contact)([]Contact){
    //Gets the first alpha contacts of a slice of contacts
    alphaClosest := make([]Contact, 0, ALPHA)
    for _,c := range contacts{
        if (len(alphaClosest) < cap(alphaClosest)){
            alphaClosest = append(alphaClosest, c)
        }else{
            return alphaClosest
        }
    }
    return alphaClosest
}


func RemoveInactiveContacts(shortlist []Contact, node_state map[ID]string) (new_shortlist []Contact) {
    for _, c := range shortlist {
        if node_state[c.NodeID] == "inactive"{
            continue
        }else{
            new_shortlist = append(new_shortlist, c)
        }
   } 
   return
}

func UpdateShortlist(shortlist []Contact, alpha_contacts[]Contact, destID ID, node_state map[ID]string) []Contact {
    //Adds alpha_contacts to short list.
    //Remove duplicates or inactive contacts from shortlist.
    //Then sorts shortlist
    for _, alpha_contact := range alpha_contacts {
        should_add := true
        for _, cur_contact := range shortlist {
            if alpha_contact.NodeID == cur_contact.NodeID ||
            node_state[alpha_contact.NodeID] == "inactive" {
                should_add = false
                break
            } 
        }
        if should_add {
            shortlist = append(shortlist, alpha_contact)
        }
    }
    //sort new_shortlist
    shortlist = SortContacts(shortlist, destID)
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
        } else {
            alpha_contacts_to_rpc = append(alpha_contacts_to_rpc, c)
            if len(alpha_contacts_to_rpc) == cap(alpha_contacts_to_rpc) {
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
