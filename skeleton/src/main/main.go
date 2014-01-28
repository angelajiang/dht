package main 
import (
    "flag"
    "fmt"
    "log"
    "math/rand"
    "net"
    "net/http"
    "net/rpc"
    "time"
)

import (
    "kademlia"
)

func main() {


    // By default, Go seeds its RNG with 1. This would cause every program to
    // generate the same sequence of IDs.
    rand.Seed(time.Now().UnixNano())

    // Get the bind and connect connection strings from command-line arguments.
    flag.Parse()
    args := flag.Args()
    if len(args) != 2 {
        log.Fatal("Must be invoked with exactly two arguments!\n")
    }
    listen_str := args[0]
    firstPeerStr := args[1]

    fmt.Printf("kademlia starting up!\n")
    kadem := kademlia.NewKademlia()
    /*
    tmp_id := kadem.NodeID
    tmp_ip := net.ParseIP("127.0.0.1")
    tmp_contact := kademlia.Contact{tmp_id, tmp_ip, 4000}
    kademlia.Update(tmp_contact, &kadem.Buckets[0])
    fmt.Printf("contact list length (1) %v\n", len(kadem.Buckets[0].Contacts))

    tmp_id = kademlia.NewRandomID()
    tmp_ip = net.ParseIP("123.123.123.123")
    tmp_contact = kademlia.Contact{tmp_id, tmp_ip, 5000}
    kademlia.Update(tmp_contact, &kadem.Buckets[0])
    fmt.Printf("contact list length (2) %v\n", len(kadem.Buckets[0].Contacts))
    */
    rpc.Register(kadem)
    rpc.HandleHTTP()
    l, err := net.Listen("tcp", listen_str)
    if err != nil {
        log.Fatal("Listen: ", err)
    }

    // Serve forever.
    go http.Serve(l, nil)

    // Confirm our server is up with a PING request and then exit.
    // Your code should loop forever, reading instructions from stdin and
    // printing their results to stdout. See README.txt for more details.
    client, err := rpc.DialHTTP("tcp", firstPeerStr)
    if err != nil {
        log.Fatal("DialHTTP: ", err)
    }
    ping := new(kademlia.Ping)
    ping.MsgID = kademlia.NewRandomID()
    var pong kademlia.Pong
    err = client.Call("Kademlia.Ping", ping, &pong)
    if err != nil {
        log.Fatal("Call: ", err)
    }

    log.Printf("ping msgID: %s\n", ping.MsgID.AsString())
    log.Printf("pong msgID: %s\n", pong.MsgID.AsString())

    var pong2 kademlia.Pong
    listen_netip, peer_uint16 := kademlia.PeerStrToHostPort(listen_str)
    pong2, err = kademlia.DoPing(listen_netip, peer_uint16)
    log.Printf("pong msg from doping %v\n", pong2.MsgID.AsString())
}

