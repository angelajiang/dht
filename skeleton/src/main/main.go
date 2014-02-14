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
    "bufio"
    "os"
    "strings"
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
    //first_peer_str := args[1]

    fmt.Printf("kademlia starting up!\n")
    host, port := kademlia.PeerStrToHostPort(listen_str)
    //fmt.Printf("host: %v\n", host)
    //fmt.Printf("port: %v\n", port)
    kadem := kademlia.NewKademlia(host, port)
    //fmt.Printf("kadem NodeID: %v\n", kadem.NodeID)
    rpc.Register(kadem)
    rpc.HandleHTTP()
    l, err := net.Listen("tcp", listen_str)
    if err != nil {
        log.Fatal("Listen: ", err)
    }

    // Serve forever.
    go http.Serve(l, nil)

    //Testing
    //kademlia.TestPingFirstPeer(kadem, first_peer_str)
    //kademlia.TestBasicRPCs(kadem, first_peer_str)

    /* looping forever, reading from stdin */
    for {
        bio := bufio.NewReader(os.Stdin)
        line, _, cmd_err := bio.ReadLine()
        if cmd_err != nil {
            log.Fatal("ReadLine failed: ", cmd_err)
        }
        // convert line from an array of bytes to a string array 
        cmdline := string(line)
        cmdline_args := strings.Split(cmdline, " ")
        command := cmdline_args[0]
        switch command {
            case "ping":
                ping := new(kademlia.Ping)
                ping.MsgID = kademlia.NewRandomID()
                ping.Sender.NodeID = kadem.NodeID
                fmt.Printf("sender NodeID: %v\n", ping.Sender.NodeID)
                ping.Sender.Host = host
                ping.Sender.Port = port
                var pong_from_host kademlia.Pong
                host_to_ping := cmdline_args[1]
                if strings.Contains(host_to_ping, ":") {
                    listen_netip, peer_uint16 := kademlia.PeerStrToHostPort(host_to_ping)
                    pong_from_host, err = kademlia.CallPing(listen_netip, peer_uint16)
                    log.Printf("pong MsgID: %v\n", pong_from_host.MsgID.AsString())
                } else {
                }
            case "whoami":
                fmt.Printf("%v\n", kadem.NodeID)
            case "local_find_value":
                fmt.Printf("local_find_value")
                key_array := []byte(cmdline_args[1])
                key_id, err := kademlia.FromByteArray(key_array)
                if err != nil {
                    fmt.Printf("error converting from byte array to ID\n")
                } else {
                    kademlia.FindValueLocally(kadem, key_id)
                }
            case "get_contact":
                id_array := []byte(cmdline_args[1])
                id, err := kademlia.FromByteArray(id_array)
                if err != nil {
                    fmt.Printf("error converting from byte array to ID\n")
                } else {
                    kademlia.FindContactLocally(kadem, id)
                }
            case "store":
                /*
                node_id := []byte(cmdline_args[1])
                key := []byte(cmdline_args[2])
                val := []byte(cmdline_args[2])
                
                TODO: write a function that returns a contact from the node_id
                contact := kademlia.FindNodeFromNodeID(node_id)
                contact := new(Contact)
                contact.NodeID = node_id
                //contact.Port =
                //contact.Host = 
                err := kademlia.CallStore(contact, key, val)
                if err != nil {
                    fmt.Printf("error storing value\n")
                }
                */
            case "find_node":
                /*
                node_id := []byte(cmdline_args[1])
                key := []byte(cmdline_args[2])
                */
            case "find_value":
                /*
                node_id := []byte(cmdline_args[1])
                key := []byte(cmdline_args[2])
                */
            case "iterativeStore":
                /*
                key := []byte(cmdline_args[1])
                val := []byte(cmdline_args[2])
                */
            case "iterativeFindNode":
                /*
                node_id := []byte(cmdline_args[1])
                */
            case "iterativeFindValue":
                /*
                key := []byte(cmdline_args[1])
                */
        }
    }
}

