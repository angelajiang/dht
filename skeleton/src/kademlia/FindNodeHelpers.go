package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"sort"
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
    ds := new(DistanceSorter)
    ds.Contacts = k.Buckets[index].Contacts
    ds.DestID = requestID
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


func PrefixLength(id ID, other ID) (dist int) {
    dist = id.PrefixLen() - other.PrefixLen()
    if dist < 0 {
        dist = dist * -1
    }
    return
}

func ContactsToFoundNodes(contacts []Contact)(foundNodes []FoundNode){
    //Takes a splice of contacts and transforms it into a splice of foundNodes
    //Output can be stored in a FindNodeResult
    foundNodes = make([]FoundNode, 0)
    for _, contact := range contacts{
        f := new(FoundNode)
        f.NodeID = contact.NodeID
        f.Port = contact.Port
        f.IPAddr = contact.Host.String()
        foundNodes = append(foundNodes, *f)
    }
    return
}

