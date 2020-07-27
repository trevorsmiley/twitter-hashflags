package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
	"twitter-hashflags/hashflag"
)

type HashFlagTwitter struct {
	CampaignName        string `json:"campaignName"`
	Hashtag             string `json:"hashtag"`
	AssetUrl            string `json:"assetUrl"`
	StartingTimestampMs string `json:"startingTimestampMs"`
	EndingTimestampMs   string `json:"endingTimestampMs"`
}

func GetHashflagsFromTwitter() ([]hashflag.Hashflag, error) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	hourFlags, err := GetHashFlagsUnmarshalled(now.Format("2006-01-02-15"))
	dayFlags, err := GetHashFlagsUnmarshalled(now.Format("2006-01-02"))

	allFlags := append(dayFlags, hourFlags...)
	if len(allFlags) == 0 {
		return nil, err
	}
	return groupHashflags(allFlags), nil
}

func GetHashFlagsUnmarshalled(timeString string) ([]HashFlagTwitter, error) {
	uri := fmt.Sprintf("https://pbs.twimg.com/hashflag/config-%s.json", timeString)
	r, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Timestring used: %s\n", timeString)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatal("Error closing body", err)
		}
	}()

	var hfs []HashFlagTwitter
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &hfs)
	if err != nil {
		return nil, err
	}
	return hfs, nil
}

func groupHashflags(hashflags []HashFlagTwitter) []hashflag.Hashflag {
	grouped := make(map[string]hashflag.Hashflag)
	for _, hashf := range hashflags {
		filename := hashflag.GetFileName(hashf.AssetUrl)
		if hf, ok := grouped[filename]; ok {
			hf.Hashtags = append(hf.Hashtags, hashf.Hashtag)
			grouped[filename] = hf
		} else {
			u, _ := url.Parse(hashf.AssetUrl)
			grouped[filename] = hashflag.Hashflag{
				URL:      *u,
				Hashtags: []string{hashf.Hashtag},
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
