package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

// Core Kademlia type. You can put whatever state you want in this.

import (
    "net"
    "fmt"
    "errors"
)

const NUMBUCKETS int =  160
const NUMCONTACTS int = 20
const VALUESIZE int = 160
const ALPHA int = 3

type Kademlia struct {
    NodeID ID
    Buckets []Bucket
    Host net.IP
    Port uint16
    Data map[ID][]byte
}

func NewKademlia(host net.IP, port uint16) *Kademlia {
    // TODO: Assign yourself a random ID and prepare other state here.
    kptr := new(Kademlia)
    kptr.NodeID = NewRandomID()
    kptr.Buckets = make([]Bucket, NUMBUCKETS)
    for i,_ := range kptr.Buckets{
        kptr.Buckets[i] = *(NewBucket())
    }
    kptr.Host = host
    kptr.Port = port
    kptr.Data = make(map[ID][]byte, VALUESIZE)
    return kptr
}

func (k *Kademlia) GetContact() (Contact){
    c := new(Contact)
    c.NodeID = k.NodeID
    c.Host = k.Host
    c.Port = k.Port
    return *c
}

func Update(k *Kademlia, contact *Contact) error {
    //Choose correct bucket to put contact
    fmt.Printf("\nContact used in update: %v\n", contact.NodeID)
    fmt.Printf("New Contact Host: %v, Port: %v\n", contact.Host, contact.Port)
    distance := k.NodeID.Xor(contact.NodeID)
    bucket_index := GetBucketIndex(distance)
    fmt.Printf("Adding to Bucket %v\n", bucket_index)
    bucket_addr := &k.Buckets[bucket_index]
    bucket := *bucket_addr
    
    fmt.Printf("len(Bucket %v) before update: %v\n", bucket_index, len(k.Buckets[bucket_index].Contacts))
    in_bucket, index := bucket_addr.InBucket(contact)
    is_full := bucket_addr.IsFull()
    switch {
    case in_bucket:
        /*Move contact to end of bucket's contact list*/
        //fmt.Printf("Case: in_bucket\n")
        //FIXED: GIVES OUT OF BOUNDS ERROR
        if len(bucket.Contacts) > 1 {
            bucket.Contacts = MoveToEnd(bucket.Contacts, index)
        }
    case !in_bucket && !is_full:
        //fmt.Printf("Case: !in_bucket, !is_full\n")
        bucket_addr.Contacts = append(bucket_addr.Contacts, *contact)
    case !in_bucket && is_full:
        //fmt.Printf("Case: !in_bucket and is_full\n")
        /*Replace head of list if head doesn't respond. Otherwise, ignore*/
        fmt.Printf("Ping'd Contact. Host: %v, Port: %v\n",
        bucket_addr.Contacts[0].Host, bucket_addr.Contacts[0].Port)
        pong, err := CallPing(k, bucket_addr.Contacts[0].Host,
        bucket_addr.Contacts[0].Port)//bucket_addr.Contacts[0].Port)
        fmt.Printf("%+v\n", pong)
        if err != nil {
            //drop head and append contact to end of list
            bucket_addr.Contacts = append(bucket_addr.Contacts[1:], *contact)
        } else {
            //Move head to tail
            bucket_addr.Contacts = append(bucket_addr.Contacts[1:],bucket_addr.Contacts[0])
        }
    }
    fmt.Printf("len(Bucket %v) After update: %v\n",
    bucket_index, len(k.Buckets[bucket_index].Contacts))
    return errors.New("function not implemented")
}
