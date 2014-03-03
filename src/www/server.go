package main

import (
    "flag"
    "net/http"
    "fmt"
    "encoding/json"
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

type test_struct struct {
    Test string
}

func handlerProcessTags(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Access-Control-Allow-Origin", "*")
  	fmt.Printf("client requested %v\n", r.URL.Path[1:])
    err := r.ParseForm()
    if err != nil{
      log.Print(err)
    }
    values := r.Form
    fmt.Printf("values: %v\n", values)
    jsonResp := Response{"tags":values["tags"]}
    fmt.Printf("json: %v\n", jsonResp)
    fmt.Fprint(w, jsonResp)
}

func main() {
    flag.Parse()
    http.HandleFunc("/processTags", handlerProcessTags)
    http.ListenAndServe(":"+*port, nil)
}
