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

func Update(contact *Contact, bucket_addr *Bucket) error {
    fmt.Printf("bucket len in Update is: %v\n", len(bucket_addr.Contacts))
    bucket := *bucket_addr
    in_bucket, index:= InBucket(contact, bucket)
    is_full := IsFull(bucket)
    switch {
    case in_bucket:
        /*Move contact to end of bucket's contact list*/
        fmt.Printf("Case: in_bucket\n")
        bucket.Contacts = append(bucket.Contacts[:index-1],bucket.Contacts[(index+1):]...)
        bucket.Contacts = append(bucket.Contacts, *contact)
    case !in_bucket && !is_full:
        if len(bucket_addr.Contacts) == 0{
            fmt.Printf("Case: !in_bucket, !is_full, empty\n")
            bucket_addr.Contacts = append(bucket_addr.Contacts, *contact)
        } else {
            fmt.Printf("Case: !in_bucket, !is_full, !empty\n")
            pong, err := CallPing(bucket_addr.Contacts[0].Host, bucket_addr.Contacts[0].Port)
            fmt.Printf("%+v\n", pong)
            if err != nil{
                bucket_addr.Contacts = append(bucket_addr.Contacts[1:], *contact)
            }
            bucket_addr.Contacts = append(bucket_addr.Contacts, *contact)
        }
    case !in_bucket && is_full:
        fmt.Printf("Case: !in_bucket and is_full\n")
        /*Replace head of list if head doesn't respond. Otherwise, ignore*/
        pong, err := CallPing(bucket_addr.Contacts[0].Host, bucket_addr.Contacts[0].Port)
        fmt.Printf("%+v\n", pong)
        if err != nil{
            //drop head append contact to end of list
            bucket_addr.Contacts = append(bucket_addr.Contacts[1:], *contact)
        } else {
            //Move head to tail
            bucket_addr.Contacts = append(bucket_addr.Contacts[1:],bucket_addr.Contacts[0])
        }
    }
    return errors.New("function not implemented")
}

func InBucket(contact *Contact, bucket Bucket) (in_bucket bool, index int) {
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

