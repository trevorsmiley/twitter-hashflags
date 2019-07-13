package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"twitter-hashflags/hashflag"
	"twitter-hashflags/twitter"
	"twitter-hashflags/utils"
)

const (
	hashflagDIR = "downloaded_hashflags"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Invalid arguments")
	}

	op := os.Args[1]

	hashflags, err := twitter.GetHashflagsFromTwitter()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d active hashflags\n", len(hashflags))

	switch op {
	case "get":
		for _, name := range alphaList(hashflags) {
			fmt.Printf("%s\n", name)
		}
	case "get-full":
		for _, hf := range hashflags {
			fmt.Printf("%s | %s\n%s\n", hf.GetName(), hf.URL.String(), strings.Join(hf.Hashtags, ", "))
		}
	case "download":
		utils.CreateDirIfMissing(hashflagDIR)
		missingHashflags := filterMissingHashflags(hashflags)
		fmt.Printf("Syncing %d hashflags to /%s\n", len(missingHashflags), hashflagDIR)

		for _, hf := range missingHashflags {
			fmt.Printf("Downloading '%s' from %s\n", hf.GetFileName(), hf.URL.String())
			hf.Download(hashflagDIR)
		}
	case "force-download":
		fmt.Printf("Downloading all hashflags to /%s\n", hashflagDIR)
		forceDownloadAll(hashflags)

	case "diff":
		missingHashflags := filterMissingHashflags(hashflags)
		fmt.Printf("%d Missing hashflags\n", len(missingHashflags))
		for _, hf := range missingHashflags {
			fmt.Printf("%s\n", hf.GetFileName())
		}
	}

}

func alphaList(hashflags map[string]hashflag.Hashflag) []string {
	list := make([]string, 0)
	for _, hf := range hashflags {
		list = append(list, hf.GetName())
	}

	//case insensitive sort
	sort.Slice(list, func(i, j int) bool { return strings.ToLower(list[i]) < strings.ToLower(list[j]) })
	return list
}

func forceDownloadAll(hashflags map[string]hashflag.Hashflag) {
	err := utils.CreateOrClearDir(hashflagDIR)
	if err != nil {
		log.Fatal(err)
	}
	for _, hf := range hashflags {
		hf.Download(hashflagDIR)
	}
}

func filterMissingHashflags(hashflags map[string]hashflag.Hashflag) []hashflag.Hashflag {
	filtered := make([]hashflag.Hashflag, 0)
	filenames, err := utils.GetFileNames(hashflagDIR)
	if err != nil {
		log.Fatal(err)
	}
	for _, hf := range hashflags {
		if utils.ContainsString(filenames, hf.GetFileName()) {
			continue
		}
		filtered = append(filtered, hf)
	}
	hashflag.SortHashflags(filtered)
	return filtered
}
