package twitter

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"twitter-hashflags/hashflag"

	"golang.org/x/net/html"
)

const (
	twitterHashflagTag = "init-data"
	twitterURL         = "https://www.twitter.com"
)

type InitData struct {
	ActiveHashflags map[string]string `json:"activeHashflags"`
}

func GetHashflagsFromTwitter() ([]hashflag.Hashflag, error) {
	resp, err := http.Get(twitterURL)
	if err != nil {
		log.Fatalf("Couldn't fetch %s\n%v\n", twitterURL, err)
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal("Error closing body", err)
		}
	}()

	root, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal("Error parsing body", err)
	}
	element, ok := getElementById(twitterHashflagTag, root)
	if !ok {
		log.Fatal("element not found")
	}
	for _, a := range element.Attr {
		if a.Key == "value" {
			hashflags := getActiveHashflags(a.Val)
			return groupHashflags(hashflags), nil
		}
	}
	log.Fatal("element missing value")
	return nil, err
}

func getElementById(id string, n *html.Node) (element *html.Node, ok bool) {
	for _, a := range n.Attr {
		if a.Key == "id" && a.Val == id {
			return n, true
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if element, ok = getElementById(id, c); ok {
			return
		}
	}
	return
}

func getActiveHashflags(s string) InitData {
	hf := InitData{}
	err := json.Unmarshal([]byte(s), &hf)
	if err != nil {
		log.Fatal("Error unmarshaling active hashflags", err)
	}
	return hf
}

func groupHashflags(hashflags InitData) []hashflag.Hashflag {
	grouped := make(map[string]hashflag.Hashflag)
	for hashtag, uri := range hashflags.ActiveHashflags {
		filename := hashflag.GetFileName(uri)
		if hf, ok := grouped[filename]; ok {
			hf.Hashtags = append(hf.Hashtags, hashtag)
			grouped[filename] = hf
		} else {
			u, _ := url.Parse(uri)
			grouped[filename] = hashflag.Hashflag{
				URL:      *u,
				Hashtags: []string{hashtag},
			}
		}
	}

	for filename, hf := range grouped {
		sort.Strings(hf.Hashtags)
		grouped[filename] = hf
		//fmt.Printf("%s - %s - %v\n", hf.GetName(), hf.GetFileExtension(), hf.URL.String())
	}

	list := make([]hashflag.Hashflag, 0)
	for _, hf := range grouped {
		list = append(list, hf)
	}

	sort.Slice(list, func(i, j int) bool {
		return strings.ToLower(list[i].GetFileName()) < strings.ToLower(list[j].GetFileName())
	})

	return list

}
