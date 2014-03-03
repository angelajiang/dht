// httpserver.go
package main

import (
    "flag"
    "net/http"
    "fmt"
    "log"
    "io"
)

var port = flag.String("port", "8080", "Define what TCP port to bind to")
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
	var bytesRead, _ = io.ReadFull(r.Body, buf)
	fmt.Fprintf(w, "# of bytes read: %v!", bytesRead)
}

func main() {
    flag.Parse()
    http.HandleFunc("/processTags", handlerProcessTags)
    log.Fatal(http.ListenAndServe(":"+*port, nil))
    //panic(http.ListenAndServe(":"+*port, http.FileServer(http.Dir(*root))))
}
