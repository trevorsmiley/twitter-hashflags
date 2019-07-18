package main

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/gookit/color"
	"github.com/trevorsmiley/fileutils"
	"log"
	"os"
	"twitter-hashflags/hashflag"
	"twitter-hashflags/twitter"
)

const (
	hashflagDIR = "downloaded_hashflags"
	detailsFile = "hashflag-list.txt"
)

var cyan = color.FgCyan.Render
var green = color.FgGreen.Render
var red = color.FgRed.Render

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Invalid arguments")
		printHelp()
		return
	}

	op := os.Args[1]

	switch op {
	case "list":
		list(getHashflags())
	case "list-fulldetails":
		listFullDetails(getHashflags())
	case "sync":
		sync(getHashflags(), hashflagDIR)
	case "force-download":
		forceDownload(getHashflags(), hashflagDIR)
	case "diff":
		diff(getHashflags(), hashflagDIR)
	case "help":
		printHelp()
	default:
		fmt.Println("Unknown command")
		printHelp()
	}
}

func getHashflags() []hashflag.Hashflag {
	hashflags, err := twitter.GetHashflagsFromTwitter()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %s active hashflags\n\n", cyan(len(hashflags)))
	return hashflags
}

func printHelp() {
	oplist := []struct {
		op   string
		desc string
	}{
		{"list", "Print a list of all active hashflags"},
		{"list-fulldetails", color.Sprintf("Write a list of all active hashflags with hashtags to %s", green(detailsFile))},
		{"diff", color.Sprintf("List all hashflags missing in %s", green(hashflagDirPath()))},
		{"sync", color.Sprintf("Download all missing hashflags to %s", green(hashflagDirPath()))},
		{"force-download", color.Sprintf("<red>Clear directory</> %s and download all active hashflags", red(hashflagDirPath()))},
	}
	fmt.Println("Available commands:")
	for _, op := range oplist {
		color.Printf("\t%s: %s\n", cyan(op.op), op.desc)
	}
}

func hashflagDirPath() string {
	return fmt.Sprintf("%s%s", string(os.PathSeparator), hashflagDIR)
}

func forceDownload(hashflags []hashflag.Hashflag, dir string) {
	fmt.Printf("Downloading all hashflags to /%s\n", dir)
	downloadAll(hashflags, true, dir)
}

func diff(hashflags []hashflag.Hashflag, dir string) {
	missingHashflags := hashflag.FilterMissingHashflags(hashflags, dir)
	numMissing := red(len(missingHashflags))
	if len(missingHashflags) == 0 {
		numMissing = green(len(missingHashflags))
	}
	fmt.Printf("%s missing hashflags\n\n", numMissing)
	for _, hf := range missingHashflags {
		fmt.Printf("%s\n", hf.GetFileName())
	}
}

func sync(hashflags []hashflag.Hashflag, dir string) {
	missingHashflags := hashflag.FilterMissingHashflags(hashflags, dir)
	if len(missingHashflags) > 0 {
		fmt.Printf("Syncing %s hashflags to /%s\n", green(len(missingHashflags)), cyan(dir))
		downloadAll(missingHashflags, false, dir)
		color.FgGreen.Println("Complete")
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
	tmpl, err := hashflag.GetTemplate()
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
	fmt.Printf("Writing details to %s\n", cyan(detailsFile))
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
