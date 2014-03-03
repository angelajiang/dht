package main

import (
    "flag"
    "net/http"
    "fmt"
    "encoding/json"
    "io"
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
/*
	type Request struct {
		Method string // GET, POST, PUT, etc.
		URL *url.URL
		Proto      string // "HTTP/1.0"
		ProtoMajor int    // 1
		ProtoMinor int    // 0
    	// A header maps request lines to their values.
   		Header Header
   		Body io.ReadCloser
   		ContentLength int64
   		TransferEncoding []string
   		Close bool
   		Host string
   		Form url.Values
   		PostForm url.Values
		...
   	}
 */
func handlerProcessTags(w http.ResponseWriter, r *http.Request){
	fmt.Printf("client requested %v\n", r.URL.Path[1:])
	var buf = make([]byte, 200)
  //var bytesRead, _ = io.ReadFull(r.Body, buf)
	var _, _ = io.ReadFull(r.Body, buf)
	//fmt.Fprintf(w, "# of bytes read: %v!", bytesRead)
  w.Header().Set("Content-Type", "application/json")
  //fmt.Fprint(w, Response{"resp":"hi!"})
  fmt.Fprint(w, "alert('hi');")
}

func main() {
    flag.Parse()
    http.HandleFunc("/processTags", handlerProcessTags)
    http.ListenAndServe(":"+*port, nil)
    //log.Fatal(http.ListenAndServe(":"+*port, nil))
    //panic(http.ListenAndServe(":"+*port, http.FileServer(http.Dir(*root))))
}
