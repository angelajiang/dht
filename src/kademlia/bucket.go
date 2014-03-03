package kademlia
//Bucket definition and modules

import (
)

type Bucket struct {
    Contacts []Contact
}

func NewBucket() *Bucket{
    bucket_ptr := new(Bucket)
    bucket_ptr.Contacts = make([]Contact,0,NUMCONTACTS)
    return bucket_ptr
}

func (b *Bucket) InBucket(contact *Contact) (in_bucket bool, index int) {
    /*Returns true if contact is in contact list of bucket*/
    in_bucket = false
    for i, cur_contact := range b.Contacts {
        index = i
        if contact.NodeID == cur_contact.NodeID{
            in_bucket = true
            return
        }
    }
    return
}

func (b *Bucket) IsFull() bool {
    /*Returns true if bucket is full*/
    if len(b.Contacts) == cap(b.Contacts){
        return true
    }
    return false
}

