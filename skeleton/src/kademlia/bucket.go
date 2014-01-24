package kademlia
//Bucket definition and modules

import (
    "errors"
    "fmt"
)

type Bucket struct {
    Contacts []Contact
}

func NewBucket() *Bucket{
    bucket_ptr := new(Bucket)
    bucket_ptr.Contacts = make([]Contact,0,NUMCONTACTS)
    return bucket_ptr
}

func Update(contact Contact, bucket Bucket) error {
    in_bucket, index:= InBucket(contact, bucket)
    is_full := IsFull(bucket)
    switch {
    case in_bucket:
        /*Move contact to end of bucket's contact list*/
        bucket.Contacts = append(bucket.Contacts[:index-1],bucket.Contacts[(index+1):]...)
        bucket.Contacts = append(bucket.Contacts, contact)
    case !in_bucket && !is_full:
        //TODO: if empty, just add to list
        if len(bucket.Contacts) == 0{
            bucket.Contacts = append(bucket.Contacts, contact)
        }else{
            fmt.Printf("Bucket capacity: %v\n", cap(bucket.Contacts))
            pong, err := DoPing(bucket.Contacts[0].Host, bucket.Contacts[0].Port)
            fmt.Printf("%+v\n", pong)
            fmt.Printf("%s\n", err)
            if err != nil{
                bucket.Contacts = append(bucket.Contacts[1:], contact)
            }
            bucket.Contacts = append(bucket.Contacts, contact)
        }
    case !in_bucket && is_full:
        /*Replace head of list if head doesn't respond. Otherwise, ignore*/
        pong, err := DoPing(bucket.Contacts[0].Host, bucket.Contacts[0].Port)
        fmt.Printf("%+v\n", pong)
        fmt.Printf("%s\n", err)
        if err != nil{
            bucket.Contacts = append(bucket.Contacts[1:], contact)
        }
        
        //if head fails to respond
            //drop head
            //append contact to end of list
        //else
            //Move head to tail

    }
    return errors.New("function not implemented")
}

func InBucket(contact Contact, bucket Bucket) (in_bucket bool, index int) {
    /*Returns true if contact is in contact list of bucket*/
    in_bucket = false
    for i,cur_contact := range bucket.Contacts {
        index = i
        if contact.NodeID == cur_contact.NodeID{
            in_bucket = true
            return
        }
    }
    return
}

func IsFull(bucket Bucket) bool {
    /*Returns true if bucket is full*/
    if len(bucket.Contacts) == cap(bucket.Contacts){
        return true
    }
    return false
}

