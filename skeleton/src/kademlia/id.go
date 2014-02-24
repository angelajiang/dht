package kademlia
// Contains definitions for the 160-bit identifiers used throughout kademlia.

import (
    "encoding/hex"
    "math/rand"
    "time"
    "fmt"
)


// IDs are 160-bit ints. We're going to use byte arrays with a number of
// methods.
const IDBytes = 20
type ID [IDBytes]byte

func (id ID) AsString() string {
    return hex.EncodeToString(id[0:IDBytes])
}

func (id ID) Xor(other ID) (ret ID) {
    for i := 0; i < IDBytes; i++ {
        ret[i] = id[i] ^ other[i]
    }
    return
}

// Return -1, 0, or 1, with the same meaning as strcmp, etc.
func (id ID) Compare(other ID) int {
    for i := 0; i < IDBytes; i++ {
        difference := int(id[i]) - int(other[i])
        switch {
        case difference == 0:
            continue
        case difference < 0:
            return -1
        case difference > 0:
            return 1
        }
    }
    return 0
}

func (id ID) Equals(other ID) bool {
    return id.Compare(other) == 0
}

func (id ID) Less(other ID) bool {
    return id.Compare(other) < 0
}

// Return the number of consecutive zeroes, starting from the low-order bit, in
// a ID.
func (id ID) PrefixLen() int {
    for i:= 0; i < IDBytes; i++ {
        for j := 0; j < 8; j++ {
            if (id[i] >> uint8(j)) & 0x1 != 0 {
                return (8 * i) + j
            }
        }
    }
    return IDBytes * 8
}

func GetSetBits(distance ID)(ones []int){
    //Xor'ed result distance -> slice of bits in distance that are one from MSB->LSB 
    //Returns indices of set bits in distance
    //Backwards bit traversal for compatibility with PrefixLen
    //ex) 1011 0110-> [1, 3, 4, 6, 7]
    indices := make([]int, 0)
    for i:= IDBytes-1; i >= 0; i-- {
        for j := 7; j >= 0; j-- {
            if (distance[i] >> uint8(j)) & 0x1 != 0 {
                indices = append(indices, (8*IDBytes) - (8*i+j)-1)
            }
        }
    }
    return indices
}

// Generate a new ID from nothing.
func NewRandomID() (ret ID) {
    rand.Seed(time.Now().UnixNano())
    for i := 0; i < IDBytes; i++ {
        ret[i] = uint8(rand.Intn(256))
    }
    return
}

// Generate an ID identical to another.
func CopyID(id ID) (ret ID) {
    for i := 0; i < IDBytes; i++ {
        ret[i] = id[i]
    }
    return
}

// Generate a ID matching a given string.
func FromString(idstr string) (ret ID, err error) {
    bytes, err := hex.DecodeString(idstr)

    fmt.Printf("bytes: %v\n", bytes)
    fmt.Printf("len bytes: %v\n", len(bytes))
    if err != nil {
        return
    }

    for i := 0; i < IDBytes; i++ {
        ret[i] = bytes[i]
    }
    return
}

func FromByteArray(byte_array []byte) (ret ID, err error) {
    for i := 0; i < IDBytes; i++ {
        ret[i] = byte_array[i]
    }
    if err != nil {
        return
    }
    return
}
