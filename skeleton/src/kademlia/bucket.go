package kademlia
//Bucket definition and modules

import (
    "errors"
)

type Bucket struct {
    Contacts []Contact
}

func NewBucket() *Bucket{
    bucket_ptr := new(Bucket)
    (*bucket_ptr).Contacts = make([]Contact,0,NUMCONTACTS)
    return new(Bucket)
}

func Update(contact Contact, bucket Bucket) error {
    /*
    in_bucket := InBucket(contact, bucket)
    is_full := IsFull(bucket)
    switch {
    case in_bucket:
        //move to end of bucket
        //remove first
    case !in_bucket && !is_full:
        bucket.Contacts[bucket.Last+1] = contact
        bucket.Last++
        bucket.Size++
    case !in_bucket && is_full:
        //ping head of bucket
        //if head fails to respond
            //drop head
            //append contact to end of list
        //else
            //Move head to tail

    }
    */
    return errors.New("function not implemented")
}

func InBucket(contact Contact, bucket Bucket) bool {
    /*
    for  i:=0; i<bucket.Last; i++{
        if bucket.Contacts[i].NodeID == contact.NodeID{
            return true
        }
    }
    */
    return false
}

func IsFull(bucket Bucket) bool {
    /*
    if NUMBUCKETS == bucket.Last{
        return true
    }
    */
    return false
}

