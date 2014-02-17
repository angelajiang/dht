package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"sort"
    "net"
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


