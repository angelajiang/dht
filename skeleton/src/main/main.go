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
    first_peer_str := args[1]

    fmt.Printf("kademlia starting up!\n")
    host, port := kademlia.PeerStrToHostPort(listen_str)
    kadem := kademlia.NewKademlia(host, port)
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
    if err != nil {
        log.Fatal("DialHTTP: ", err)
    }
    client, err := rpc.DialHTTP("tcp", first_peer_str)
    ping := new(kademlia.Ping)
    ping.MsgID = kademlia.NewRandomID()
    var pong kademlia.Pong
    err = client.Call("Kademlia.Ping", ping, &pong)
    if err != nil {
        log.Fatal("Call: ", err)
    }

    log.Printf("ping msgID: %s\n", ping.MsgID.AsString())
    log.Printf("pong msgID: %s\n", pong.MsgID.AsString())
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
                var pong_from_host kademlia.Pong
                host := cmdline_args[1]
                if strings.Contains(host, ":") {
                    listen_netip, peer_uint16 := kademlia.PeerStrToHostPort(host)
                    pong_from_host, err = kademlia.DoPing(listen_netip, peer_uint16)
                    log.Printf("pong MsgID: %v\n", pong_from_host.MsgID.AsString())
                }
            case "whoami":
                fmt.Printf("whoami")
            case "local_find_value":
                fmt.Printf("local_find_value")

        }
    }
    var pong2 kademlia.Pong
    listen_netip, peer_uint16 := kademlia.PeerStrToHostPort(first_peer_str)
    pong2, err = kademlia.DoPing(listen_netip, peer_uint16)
    log.Printf("pong msg from doping %v\n", pong2.MsgID.AsString())

    //Making new contacts and calling Update
    tmp_id := kadem.NodeID
    tmp_ip := net.ParseIP("127.0.0.1")
    tmp_contact := kademlia.Contact{tmp_id, tmp_ip, 7890}
    kademlia.Update(tmp_contact, &kadem.Buckets[0])

    tmp_id = kademlia.NewRandomID()
    tmp_ip = net.ParseIP("123.123.123.123")
    tmp_contact = kademlia.Contact{tmp_id, tmp_ip, 7890}
    kademlia.Update(tmp_contact, &kadem.Buckets[0])

    //Putting value into tmp_contact for testing DoFindValue
    s := make([]byte, 5)
    tmp_data := s
    fmt.Printf("tmpdata: %v\n", tmp_data)
}

