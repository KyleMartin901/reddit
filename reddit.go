// Package reddit implements basic client for the Reddit API.
package reddit

import (
  "net/http"
  "encoding/json"
  "fmt"
  "errors"
  "log"
)

// Item describes a Reddit item.
type Item struct {
  Title string
  URL string
  Comments int `json:"num_comments"`
}

type response struct {
  Data struct{
    Children []struct {
      Data Item
    }
  }
}

// Get fetches the most recent Items posted to the specified subreddit.
func Get(reddit string) ([]Item, error) {
  url := fmt.Sprintf("http://reddit.com/r/%s.json", reddit)

  client := &http.Client{}

  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    log.Fatal(err)
  }

  req.Header.Set("User-Agent", "Golang_Spider_Bot")

  resp, err := client.Do(req)
  if err != nil {
    log.Fatal(err)
  }

  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK {
    return nil, errors.New(resp.Status)
  }

  r := new(response)
  err = json.NewDecoder(resp.Body).Decode(r)
  if err != nil {
    log.Fatal(err)
  }

  items := make([]Item, len(r.Data.Children))

  for i, child := range r.Data.Children  {
    items[i] = child.Data
  }

  return items, nil

}

func (i Item) String() string  {
  comm := ""

  switch i.Comments {
  case 0:
    // nothing 
  case 1:
    comm = " (1 comment)"
  default:
    comm = fmt.Sprintf(" (%d comments)", i.Comments)
  }
  return fmt.Sprintf("%s%s\n%s", i.Title, comm, i.URL)
}
