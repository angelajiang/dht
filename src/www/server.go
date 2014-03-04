package main

import (
    "flag"
    "net/http"
    "fmt"
    "encoding/json"
    "io/ioutil"
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

type RedditListing struct {
    Kind string
    Data RedditListingData
}

type RedditListingData struct {
    Children []RedditT3
}

type RedditT3 struct{
    Data RedditPost
}

type RedditPost struct {
    Url string
    Title string
    Ups int64
    Num_comments int64
}

func perror(err error) {
    if err != nil {
        panic(err)
    }
}

func getRedditLinks(subreddit string){
  url := "http://www.reddit.com/r/"+subreddit+".json"
  resp, err := http.Get(url)
  fmt.Printf("Getting URL: %v\n", url)
  perror(err)
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  perror(err)
  //var p interface{}
  var p RedditListing
  err = json.Unmarshal(body, &p)
  //var rpost RedditPost 
  fmt.Printf("message: %v\n", p)
}

func handlerProcessTags(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Access-Control-Allow-Origin", "*")
  	fmt.Printf("client requested %v\n", r.URL.Path[1:])
    err := r.ParseForm()
    perror(err)
    values := r.Form
    resp := Response{}
    tags := values["tags"]
    for i, tag := range(tags){
      fmt.Printf("tag %v: %v\n", i, tag)
      getRedditLinks(tag)
      //resp[string(i)] = tag
    }
    fmt.Printf("resp: %v\n", resp)
    fmt.Fprint(w, resp)
}

var port = flag.String("port", "5555", "Define what TCP port to bind to")
var root = flag.String("root", "static", "Define the root filesystem path")

func main() {
    flag.Parse()
    http.HandleFunc("/processTags", handlerProcessTags)
    http.ListenAndServe(":"+*port, nil)
}
