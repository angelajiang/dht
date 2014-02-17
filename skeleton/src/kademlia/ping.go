package kademlia
// Contains definitions mirroring the Kademlia spec. You will need to stick
// strictly to these to be compatible with the reference implementation and
// other groups' code.

import (
	"net/rpc"
	"log"
    "net"
    "fmt"
    "errors"
)

type Ping struct {
    Sender Contact
    MsgID ID
}

type Pong struct {
    MsgID ID
    Sender Contact
}

//RPC Call
func (k *Kademlia) Ping(ping Ping, pong *Pong) error {
    // This one's a freebie.
    pong.MsgID = CopyID(ping.MsgID)
    fmt.Printf("ping.MsgID from RPC call: %v\n", ping.MsgID.AsString())
    fmt.Printf("pong.MsgID from RPC call: %v\n", pong.MsgID.AsString())
    fmt.Printf("Ping Recepient NodeID: %v\n", k.NodeID)
    pong.Sender.NodeID = k.NodeID
    pong.Sender.Host = k.Host
    pong.Sender.Port = k.Port
    Update(k, &ping.Sender)
    return nil
}

func CallPing(k *Kademlia, remote_host net.IP, port uint16) (Pong, error){
    /* CallPing should probably take a Kademlia object here */
    //TODO: run the Update function?
    peer_str := HostPortToPeerStr(remote_host, port)
    fmt.Printf("peer_str for ping: %v\n", peer_str)
    client, err := rpc.DialHTTP("tcp", peer_str)
    if err != nil {
          log.Fatal("Call: ", err)
    }

    fmt.Printf("Making Ping struct\n")
    ping := new(Ping)
    ping.MsgID = NewRandomID()
    ping.Sender = k.GetContact()

    var pong Pong
    err = client.Call("Kademlia.Ping", ping, &pong)
    if err != nil {
        err = errors.New("Call: No resonse from ping")
          //log.Fatal("Call: ", err)
    } else {
        fmt.Printf("Calling Update From Ping!\n")
        Update(k, &pong.Sender)
    }

    client.Close()
    return pong, err
}


