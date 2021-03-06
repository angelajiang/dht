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
    "strconv"
)

import (
    "kademlia"
)

func main() {

    var err error

    // By default, Go seeds its RNG with 1. This would cause every program to
    // generate the same sequence of IDs.
    rand.Seed(time.Now().UnixNano())

    // Get the bind and connect connection strings from command-line arguments.
    flag.Parse()
    args := flag.Args()
    if len(args) != 3 {
        log.Fatal("<host>:<port> <host>:<port> <debugmode 0/1>")
    }
    var listenStr string = args[0]
    var firstPeerStr string = args[1]
    var debug bool = false
    var paramDebug int
    paramDebug, _ = strconv.Atoi(args[2])
    if (paramDebug == 1){
        debug = true
    }

    host, port := kademlia.PeerStrToHostPort(listenStr)

    if (debug){
        log.Printf("DEBUG MODE: %v\n", debug)
        port = uint16(kademlia.Random(4000,5000))
        fmt.Printf("port: %v\n", port)
        listenStr = kademlia.HostPortToPeerStr(host, port)
    }
    
    kadem := kademlia.NewKademlia(host, port)
    rpc.Register(kadem)
    rpc.HandleHTTP()
    l, err := net.Listen("tcp", listenStr)
    if err != nil {
        log.Fatal("Listen: ", err)
    }

    // Serve forever.
    go http.Serve(l, nil)

    /* looping forever, reading from stdin */
    for {
        bio := bufio.NewReader(os.Stdin)
        line, _, cmd_err := bio.ReadLine()
        if cmd_err != nil {
            log.Fatal("ReadLine failed: ", cmd_err)
        }
        // convert line from an array of bytes to a string array 
        cmdline := string(line)
        cmdlineArgs := strings.Split(cmdline, " ")
        command := cmdlineArgs[0]
        switch command {
            case "ping":
                //ping 123.12.12.0 1231
                //ping 1111 --> will ping localhost:1111
                ping := new(kademlia.Ping)
                ping.MsgID = kademlia.NewRandomID()
                ping.Sender.NodeID = kadem.NodeID
                fmt.Printf("sender NodeID: %v\n", ping.Sender.NodeID)
                ping.Sender.Host = host
                ping.Sender.Port = port
                var pong_from_host kademlia.Pong
                host_to_ping := cmdlineArgs[1]
                
                if strings.Contains(host_to_ping, ":") {
                    remoteIP, remotePort := kademlia.PeerStrToHostPort(host_to_ping)
                    fmt.Printf("IP to ping: %v\n peer: %v\n", remoteIP, remotePort)
                    pong_from_host, err = kademlia.CallPing(kadem, remoteIP, remotePort)
                    if err != nil {
                        log.Fatal("ReadLine failed: ", cmd_err)
                    }
                    log.Printf("pong MsgID: %v\n", pong_from_host.MsgID.AsString())
                } else { 
                    node_id, err := kademlia.FromString(cmdlineArgs[1])
                    if err != nil {
                        fmt.Printf("error converting from byte array to ID\n")
                    } else {
                        contact, err := kademlia.FindContactLocally(kadem, node_id)
                        if err != nil {
                            closestContacts, err := kademlia.IterativeFindNode(kadem, node_id)
                            if err != nil{
                                fmt.Printf("%v\n", err)
                            } else {
                                //ping closest contact
                                pong_from_host, err = kademlia.CallPing(kadem, closestContacts[0].Host, closestContacts[0].Port)
                                if err != nil {
                                    log.Fatal("ReadLine failed: ", cmd_err)
                                }
                                log.Printf("pong MsgID: %v\n", pong_from_host.MsgID.AsString())
                            }
                        } else {
                            //ping local contact
                            pong_from_host, err = kademlia.CallPing(kadem, contact.Host, contact.Port)
                            if err != nil {
                                log.Fatal("ReadLine failed: ", cmd_err)
                            }
                            log.Printf("pong MsgID: %v\n", pong_from_host.MsgID.AsString())
                        }
                    } 
                }
            case "whoami":
                fmt.Printf("%v\n", kadem.NodeID)
            case "local_find_value":
                fmt.Printf("local_find_value")
                key, err := kademlia.FromString(cmdlineArgs[1])
                if err != nil {
                    fmt.Printf("error converting from byte array to ID\n")
                } else {
                    v, err := kademlia.FindValueLocally(kadem, key)
                    if err != nil {
                        fmt.Printf("ERR in find_value")
                    } else {
                       fmt.Printf("id %v val %v\n", kadem.NodeID, v)  
                    }
                }
            case "get_contact":
                node_id, err := kademlia.FromString(cmdlineArgs[1])
                if err != nil {
                    fmt.Printf("error converting from byte array to ID\n")
                } else {
                    kademlia.FindContactLocally(kadem, node_id)
                }
            case "store":
                if len(cmdlineArgs) != 4 {
                    log.Printf("Error: Wrong number of arguments calling store. Expected 4, got %v\n", len(cmdlineArgs))
                    break
                }
                node_id, err := kademlia.FromString(cmdlineArgs[1])
                key, err := kademlia.FromString(cmdlineArgs[2])
                val := []byte(cmdlineArgs[3])
                contact, err := kademlia.FindContactLocally(kadem, node_id)
                if err != nil {
                    //IterativeStore
                    storedIn, err := kademlia.IterativeStore(kadem, key, val)
                    if err != nil{
                        log.Printf("%v\n", err)
                    }
                    fmt.Printf("%v stored in %v\n", key, storedIn[len(storedIn)-1].NodeID)
                } else {
                    err := kademlia.CallStore(&contact, key, val)
                    if err != nil {
                        fmt.Printf("error storing value\n")
                    }
                }
            case "find_node":
                //find_node nodeID
                node_id, err := kademlia.FromString(cmdlineArgs[1])
                if err != nil {
                    fmt.Printf("error converting from byte array to ID\n")
                } else {
                    contact, err := kademlia.FindContactLocally(kadem, node_id)
                    if err != nil {
                        closestContacts, err := kademlia.IterativeFindNode(kadem, node_id)
                        if err != nil{
                            fmt.Printf("%v\n", err)
                        }else{
                            fmt.Printf("Closest IDs: %v\n", kademlia.ContactsToIDs(closestContacts))
                        }
                    }else{
                        fmt.Printf("Found ID: %v\n", contact.NodeID)
                    }
                }
            case "find_value":
                key, err := kademlia.FromString(cmdlineArgs[2])
                if err != nil {
                    fmt.Printf("error converting from byte array to ID\n")
                } else {
                    v, local_err := kademlia.FindValueLocally(kadem, key)
                    if local_err != nil {
                        retID, foundValue, err := kademlia.IterativeFindValue(kadem, key)
                        if err != nil{
                            fmt.Printf("%v\n", err)
                            break
                        }
                        fmt.Printf("%v %v\n", retID, foundValue)
                    } else {
                       fmt.Printf("%v %v\n", kadem.NodeID, v)  
                    }
                }
            case "iterativeStore":
                if len(cmdlineArgs) != 4 {
                    log.Printf("Error: Wrong number of arguments calling iterativeStore. Expected 4, got %v\n", len(cmdlineArgs))
                    break
                }
                key, err := kademlia.FromString(cmdlineArgs[1])
                val := []byte(cmdlineArgs[2])
                storedIn, err := kademlia.IterativeStore(kadem, key, val)
                if err != nil{
                    log.Printf("%v\n", err)
                    break
                }
                fmt.Printf("%v stored in %v\n", key, storedIn[len(storedIn)-1].NodeID)
            case "iterativeFindNode":
                if len(cmdlineArgs) != 2 {
                    fmt.Printf("Error in command \"iterativeFindNode\": command must be of the form \"iterativeFindNode nodeID\"\n")
                    break
                }
                node_id, err := kademlia.FromString(cmdlineArgs[1])
                closestContacts, err := kademlia.IterativeFindNode(kadem, node_id)
                if err != nil {
                    log.Fatal("IterativeFindNode failed:\n", err)
                }
                var closestIDs []kademlia.ID = make([]kademlia.ID, 0)
                for _, contact := range closestContacts {
                    closestIDs = append(closestIDs, contact.NodeID)
                }
            case "iterativeFindValue":
                //iterativeFindValue key(ID)
                key, err := kademlia.FromString(cmdlineArgs[1])
                retID, foundValue, err := kademlia.IterativeFindValue(kadem, key)
                if err != nil{
                    fmt.Printf("%v\n", err)
                    break
                }
                fmt.Printf("%v %v\n", retID, foundValue)
            //Cases for debugging
            case "setid":
                //setid f
                d := string(cmdlineArgs[1])
                kadem.NodeID = kademlia.HexDigitToID(d, 20)
            case "test":
                //test
                kademlia.TestBasicRPCs(kadem, firstPeerStr)
            case "is":
                if len(cmdlineArgs) != 3 {
                    log.Printf("Error: Wrong number of arguments calling ifn. Expected 3, got %v\nex) is f message\n", len(cmdlineArgs))
                    break
                }
                d := string(cmdlineArgs[1])
                key := kademlia.HexDigitToID(d, 20)
                value := []byte(cmdlineArgs[2]) 
                storedIn, err := kademlia.IterativeStore(kadem, key, value)
                if err != nil{
                    log.Printf("%v\n", err)
                    break
                }
                fmt.Printf("%v stored in %v\n", key, storedIn[len(storedIn)-1])

            case "ifn":
                //ifn f
                if len(cmdlineArgs) != 2 {
                    log.Printf("Error: Wrong number of arguments calling ifn. Expected 2, got %v\n", len(cmdlineArgs))
                    break
                }
                d := string(cmdlineArgs[1])
                destID := kademlia.HexDigitToID(d, 20)

                closestContacts, err := kademlia.IterativeFindNode(kadem, destID)
                if err != nil {
                    log.Fatal("IterativeFindNode failed:\n", err)
                }
                fmt.Printf("Closest alpha contacts: %v\n", kademlia.FirstBytesOfContactIDs(closestContacts))

            case "ifv":
                //ifv f f
                if len(cmdlineArgs) != 3 {
                    log.Printf("Error: Wrong number of arguments calling ifn. Expected 3, got %v\n", len(cmdlineArgs))
                    break
                }
                d := string(cmdlineArgs[1])
                key := kademlia.HexDigitToID(d, 20)
                retID, foundValue, err := kademlia.IterativeFindValue(kadem, key)
                if err != nil{
                    fmt.Printf("%v\n", err)
                    break
                }
                fmt.Printf("%v %v\n", retID, foundValue)

            case "test_update":
                kademlia.TestUpdate(kadem, 3)
            case "test_store_local":
                contact := new(kademlia.Contact)
                contact.NodeID = kadem.NodeID
                contact.Host = kadem.Host
                contact.Port = kadem.Port
                key := kademlia.NewRandomID()
                val := make([]byte, 5)
                x := 2
                for i := 0; i < len(val); i++ {
                    val[i] = byte(x)
                    x++
                }
                kademlia.CallStore(contact, key, val)

                res := new(kademlia.FindValueResult)
                res, err := kademlia.CallFindValue(kadem, contact, key)
                if err != nil {
                    fmt.Printf("err")
                }
                fmt.Printf("value found is: %v\n", res.Value)
            case "test_store_remote":
              //make random key for storing
              data_key := kademlia.NewRandomID()
              //make random byte array for the value
              dummy_data := make([]byte, 5)
              x := 1
              //turn it into a [1 2 3 4 5] array instead of [0 0 0 0 0]
              for i := 0; i < len(dummy_data); i++ {
                dummy_data[i] = byte(x)
                x++
              }
              fmt.Printf("data to be stored: %v at key: %v\n", dummy_data,
              data_key)
              //call iterativestore
              storedIn, err := kademlia.IterativeStore(kadem, data_key,
              dummy_data)
              if err != nil{
                log.Printf("%v\n", err)
                break
              }
              fmt.Printf("%v stored in %v\n", data_key, storedIn[len(storedIn)-1])
              //call iterativefindvalue for that value
              retID, foundValue, err := kademlia.IterativeFindValue(kadem, data_key)
              if err != nil{
                fmt.Printf("%v\n", err)
                break
              }
              fmt.Printf("%v %v\n", retID, foundValue)
            
            case "get_stored_values":
              fmt.Printf("stored values in Node %v are %v\n", kadem.NodeID,
              kadem.Data)
          }
    }
}

