package main

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/trevorsmiley/fileutils"
	"html/template"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"twitter-hashflags/hashflag"
	"twitter-hashflags/twitter"
)

const (
	hashflagDIR = "downloaded_hashflags"
	detailsFile = "hashflag-list.txt"
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
	case "list":
		list(hashflags)
	case "list-fulldetails":
		listFullDetails(hashflags)
	case "sync":
		sync(hashflags, hashflagDIR)
	case "force-download":
		forceDownload(hashflags, hashflagDIR)
	case "diff":
		diff(hashflags, hashflagDIR)
	}

}

func forceDownload(hashflags []hashflag.Hashflag, dir string) {
	fmt.Printf("Downloading all hashflags to /%s\n", dir)
	downloadAll(hashflags, true, dir)
}

func diff(hashflags []hashflag.Hashflag, dir string) {
	missingHashflags := hashflag.FilterMissingHashflags(hashflags, dir)
	fmt.Printf("%d Missing hashflags\n", len(missingHashflags))
	for _, hf := range missingHashflags {
		fmt.Printf("%s\n", hf.GetFileName())
	}
}

func sync(hashflags []hashflag.Hashflag, dir string) {
	missingHashflags := hashflag.FilterMissingHashflags(hashflags, dir)
	if len(missingHashflags) > 0 {
		fmt.Printf("Syncing %d hashflags to /%s\n", len(missingHashflags), dir)
		downloadAll(missingHashflags, false, dir)
	} else {
		fmt.Println("No new hashflags to download")
	}
}

func list(hashflags []hashflag.Hashflag) {
	for _, hf := range hashflags {
		fmt.Printf("%s\n", hf.GetFileName())
	}
}

func listFullDetails(hashflags []hashflag.Hashflag) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("No caller information")
	}
	tmpl, err := template.New("hashflags.tmpl").Funcs(template.FuncMap{"StringsJoin": strings.Join}).ParseFiles(path.Dir(filename) + "/hashflag/hashflags.tmpl")
	if err != nil {
		log.Fatal("Error with template", err)
	}
	f, err := os.Create(detailsFile)
	if err != nil {
		log.Fatalf("Couldn't create file %s\n%v", detailsFile, err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatal("Couldn't close file", err)
		}
	}()
	fmt.Printf("Writing details to %s\n", detailsFile)
	err = tmpl.Execute(f, hashflags)
	if err != nil {
		log.Fatal("Error executing template", err)
	}
}

func downloadAll(hashflags []hashflag.Hashflag, clearDIR bool, dir string) {
	if clearDIR {
		err := fileutils.CreateOrClearDir(dir)
		if err != nil {
			log.Fatalf("Error with directory %s\n%v", dir, err)
		}
	} else {
		err := fileutils.CreateDirIfMissing(dir)
		if err != nil {
			log.Fatalf("Error with directory %s\n%v", dir, err)
		}
	}
	bar := pb.StartNew(len(hashflags))
	for _, hf := range hashflags {
		bar.Increment()
		hf.Download(dir)
	}
	bar.Finish()
}
