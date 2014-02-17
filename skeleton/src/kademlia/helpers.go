package kademlia
// Contains the core kademlia type. In addition to core state, this type serves
// as a receiver for the RPC methods, which is required by that package.

import (
    "net"
    "strconv"
    "strings"
    "crypto/sha1"
    "math/rand"
    "time"
    "sort"
)

//SORTING//

//Implemented sort for contacts
type IDandContacts struct {
	NodeID ID
	Contacts []Contact
}
func (ds *IDandContacts) Len() int {
    return len(ds.Contacts)
}
func (ds *IDandContacts) Swap(i, j int) {
    ds.Contacts[i], ds.Contacts[j] = ds.Contacts[j], ds.Contacts[i]
}
func (ds *IDandContacts) Less(i, j int) bool {
    i_dist := ds.Contacts[i].NodeID.Xor(ds.NodeID)
    j_dist := ds.Contacts[j].NodeID.Xor(ds.NodeID)
    return i_dist.PrefixLen() > j_dist.PrefixLen()
}

func SortContacts(contacts []Contact, destID ID)([]Contact){
    //Wrapper for using sorting function
    idc := new(IDandContacts)
    idc.Contacts = contacts
    idc.NodeID = destID
    sort.Sort(idc)
    return idc.Contacts
}

//Implemented sort for IDs
type IDSorter struct{
    NodeID ID
    IDs []ID
}
func (idsort *IDSorter) Len() int {
    return len(idsort.IDs)
}
func (idsort *IDSorter) Swap(i, j int) {
    idsort.IDs[i], idsort.IDs[j] = idsort.IDs[j], idsort.IDs[i]
}
func (idsort *IDSorter) Less(i, j int) bool {
    i_dist := idsort.IDs[i].Xor(idsort.NodeID)
    j_dist := idsort.IDs[j].Xor(idsort.NodeID)
    return i_dist.PrefixLen() > j_dist.PrefixLen()
}


func PeerStrToHostPort(listen_str string) (net.IP, uint16){
    /*Parsing*/
    input_arr := strings.Split(listen_str, ":")
    host_str := input_arr[0]
    port_str := input_arr[1]
    //Check if localhost
    if host_str == "localhost"{
        host_str = "127.0.0.1"
    }
    listen_netip := net.ParseIP(host_str)
    peer_uint64, _ := strconv.ParseUint(port_str, 10, 16)
    peer_uint16 := uint16(peer_uint64)

    return listen_netip, peer_uint16
}

func MoveToEnd(list []Contact, index int) (ret []Contact) {
    //list = a, b, c -> b, c, a
    contact := list[index]
    ret = append(list[:index], list[(index+1):]...)
    ret = append(ret, contact)
    return
}

func HostPortToPeerStr(remote_host net.IP, port uint16) (peer_str string){
    remote_host_str := remote_host.String()
    port_uint64 := uint64(port)
    port_str :=  strconv.FormatUint(port_uint64, 10)
    peer_str = remote_host_str + ":" + port_str
    return peer_str
}

func HashKey(key ID) []byte {
    //fmt.Printf("size of key: %v\n", len(key))
    h := sha1.New()
    h.Write(key[:])
    bs := h.Sum([]byte{})
    //fmt.Printf("bs is :%v\n", bs)
    return bs
}

func GetBucketIndex(distance ID)(index int){
	/*Given distance, returns first set bit counting from MSB
	 	ex) 0011 0101 -> 3	*/
	index = 0
    for i:= IDBytes-1; i >= 0; i-- {
        for j := 7; j >= 0; j-- {
            if (distance[i] >> uint8(j)) & 0x1 != 0 {
                index = (8*IDBytes) - (8*i+j)
                return
            }
        }
    }
	return 
}
    
//TESTING HELPERS

func HexDigitToID(hex_digit string, n int) (id ID){
    //Takes hex number of all digits as "digit"
    //returns number as n-byte id
    //HexDigitToID(f,160) -> ffff....ffff to byte array of size 160
    id, _ = FromString(strings.Repeat(hex_digit, n*2))
    return
}


func Random(min, max int) int {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(max - min) + min
}

func NewIterativeContacts(n int)(contacts []Contact) {
    starting_port := 4000
    for i := 0; i < n; i++ {
        nodeid := NewRandomID()
        ip := net.ParseIP("127.0.0.1")
        port := uint16(starting_port)
        new_contact := Contact{nodeid, ip, port}
        contacts = append(contacts, new_contact)
    }
    return 
}

func NewRandomContact()(*Contact) {
    port := uint16(Random(4000,4003))
    ip := net.ParseIP("127.0.0.1")
    nodeid := NewRandomID()
    return &Contact{nodeid, ip, port}
}


func FirstBytesOfIDs(ids []ID) (bytes []byte){
    bytes = make([]byte, 0)
    for _, id := range ids{
        bytes = append(bytes, id[0])
    }
    return
}
func FirstBytesOfContactIDs(contacts []Contact) (bytes []byte){
    bytes = make([]byte, 0)
    for _, c := range contacts{
        bytes = append(bytes, c.NodeID[0])
    }
    return
}


