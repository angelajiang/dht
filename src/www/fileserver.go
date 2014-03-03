package main

import (
    "flag"
    "net/http"
    "log"
)

var port = flag.String("port", "8080", "Define what TCP port to bind to")
var root = flag.String("root", "static", "Define the root filesystem path")

func main() {
    flag.Parse()
    log.Fatal(http.ListenAndServe(":"+*port, http.FileServer(http.Dir(*root))))
}
