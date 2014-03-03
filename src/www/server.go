package main

import (
    "flag"
    "net/http"
    "fmt"
    "encoding/json"
    "io"
    "log"
)

type Response map[string]interface{}

func (r Response) String() (s string) {
        b, err := json.Marshal(r)
        if err != nil {
                s = ""
                return
        }
        s = string(b)
        return
}

var port = flag.String("port", "5555", "Define what TCP port to bind to")
var root = flag.String("root", "static", "Define the root filesystem path")

func handlerProcessTags(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Printf("client requested %v\n", r.URL.Path[1:])
	var buf = make([]byte, 200)
  var _, err = io.ReadFull(r.Body, buf)
  if err != nil{
    //TODO: make sure to read all bytes
    //EOF = no bytes read, ErrUnexpectedEOF = some, not all bytes read
    log.Print(err)
  }
  fmt.Fprint(w, Response{"myresp":1})
}

func main() {
    flag.Parse()
    http.HandleFunc("/processTags", handlerProcessTags)
    http.ListenAndServe(":"+*port, nil)
}
