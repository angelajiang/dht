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

