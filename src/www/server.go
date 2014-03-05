package main

import (
    "flag"
    "net/http"
    "net/url"
    "log"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "strconv"
    "strings"
)

type PostsResponse map[string][]RedditPost

func (r PostsResponse) String() (s string) {
    b, err := json.Marshal(r)
    perror(err, "handlerProcessTags", "marshalling response into JSON")
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
    Downs int64
    Num_comments int64
}

func perror(err error, who string, why string) {
    if err != nil {
      log.Printf("Error in %v when %v\n", who, why)
      panic(err)
    }
}

func getRedditListing(subreddit string) RedditListing{
  var url string = "http://www.reddit.com/r/"+subreddit+".json"
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
    var err error = r.ParseForm()
    perror(err, "handlerProcessTags", "parsing client request form")
    var values url.Values = r.Form
    //var tags []string = values["tags"]  //use if format: ["tag1", "tag2"]
    var strTags string = values["tags"][0] 
    var tags []string = strings.Split(strTags, " ")
    numLinks, err := strconv.ParseInt(values["numLinks"][0], 10, 0)
    perror(err, "handlerProcessTags", "Parsing value[numLinks] as int")
    var resp PostsResponse = PostsResponse{}
    for _, tag := range(tags){
      var rl RedditListing = getRedditListing(tag)
      var posts []RedditPost = GetNPosts(rl, numLinks)
      resp[tag] = posts
    }
    fmt.Fprint(w, resp.String())
}

var port = flag.String("port", "5555", "Define what TCP port to bind to")
var root = flag.String("root", "static", "Define the root filesystem path")

func main() {
    flag.Parse()
    http.HandleFunc("/processTags", handlerProcessTags)
    http.ListenAndServe(":"+*port, nil)
}
