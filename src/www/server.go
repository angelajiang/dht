package main

import (
    "flag"
    "net/http"
    "log"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "strconv"
)

type Response map[string][]RedditPost

/*
func (r Response) String() (s string) {
        b, err := json.Marshal(r)
        if err != nil {
                s = ""
                return
        }
        s = string(b)
        return
}
*/

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

func perror(err error, who string, why string) {
    if err != nil {
      log.Printf("Error in %v when %v\n", who, why)
      panic(err)
    }
}

func getRedditListing(subreddit string) RedditListing{
  url := "http://www.reddit.com/r/"+subreddit+".json"
  resp, err := http.Get(url)
  perror(err, "GetRedditListing", "getting URL "+url)
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  perror(err, "GetRedditListing", "reading body of "+url)
  var rl RedditListing
  err = json.Unmarshal(body, &rl)
  return rl
}

// Gets first n RedditPosts from a RedditListing
// TODO: Deal with replicated links
func GetNPosts(listing RedditListing, n int64)([]RedditPost){
  var posts = make([]RedditPost, 0, n) 
  var children []RedditT3 = listing.Data.Children 
  for _, child := range(children){
    var post RedditPost = child.Data
    if (cap(posts)!=len(posts)){
      posts = append(posts, post)
    }
  }
  return posts
}

func handlerProcessTags(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Access-Control-Allow-Origin", "*")
  	//fmt.Printf("client requested %v\n", r.URL.Path[1:])
    err := r.ParseForm()
    perror(err, "handlerProcessTags", "parsing client request form")
    values := r.Form
    tags := values["tags"]
    numLinks, err := strconv.ParseInt(values["numLinks"][0], 10, 0)
    perror(err, "handlerProcessTags", "Parsing value[numLinks] as int")
    var resp map[string][]RedditPost = Response{}
    for _, tag := range(tags){
      var rl RedditListing = getRedditListing(tag)
      var posts []RedditPost = GetNPosts(rl, numLinks)
      resp[tag] = posts
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
