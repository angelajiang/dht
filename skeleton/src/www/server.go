// httpserver.go
package main

import (
    "flag"
    "net/http"
)

var port = flag.String("port", "8080", "Define what TCP port to bind to")
var root = flag.String("root", "static", "Define the root filesystem path")

func handlerProcessTags(http.ResponseWriter, *http.Request){
	//handler logic
}

func main() {
    flag.Parse()
    http.HandleFunc("/processTags", handlerProcessTags)
    panic(http.ListenAndServe(":"+*port, http.FileServer(http.Dir(*root))))
}
