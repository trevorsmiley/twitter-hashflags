package main

import (
	"flag"
	"fmt"
	"github.com/gookit/color"
	"log"
	"os"
	"path/filepath"
	"twitter-hashflags/hashflag"
	"twitter-hashflags/twitter"
)

const (
	hashflagDIR    = "downloaded_hashflags"
	deactivatedDIR = "deactivated"
	detailsFile    = "hashflag-list.txt"
)

var cyan = color.FgCyan.Render
var green = color.FgGreen.Render
var red = color.FgRed.Render
var yellow = color.FgYellow.Render

func main() {

	//list
	listCommand := flag.NewFlagSet("info", flag.ExitOnError)
	fullDetails := listCommand.Bool("out", false, "write full details to file")
	fileDiff := listCommand.Bool("diff", false, "diff")
	deactivated := listCommand.Bool("deactivated", false, "deactivated")

	//sync
	syncCommand := flag.NewFlagSet("sync", flag.ExitOnError)
	moveInactive := syncCommand.Bool("m", false, "move deactivated")
	force := syncCommand.Bool("force", false, "force download")

	switch os.Args[1] {
	case "info":
		err := listCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	case "sync":
		err := syncCommand.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
	default:
		printHelp()
	}

	if listCommand.Parsed() {
		hashflag.Info(getHashflags(), hashflagDIR, detailsFile, *fullDetails, *deactivated, *fileDiff)
	}
	if syncCommand.Parsed() {
		hashflag.Sync(getHashflags(), hashflagDIR, deactivatedDIR, *moveInactive, *force)
	}
}

func printHelp() {
	opList := []struct {
		op   string
		desc string
		flag bool
	}{
		{"info", "Print a list of all active hashflags", false},
		{"-out", color.Sprintf("Write a list of all active hashflags with hashtags to %s", green(detailsFile)), true},
		{"-diff", color.Sprintf("List only hashflags missing in %s", green(hashflagDirPath())), true},
		{"-deactivated", color.Sprintf("List only deactivated files that are in %s", green(hashflagDirPath())), true},
		{"sync", color.Sprintf("Download all missing hashflags to %s", green(hashflagDirPath())), false},
		{"-m", color.Sprintf(
			"Move all deactivated files in %s to %s and download all missing hashflags to %s",
			red(hashflagDirPath()),
			red(dirPath(filepath.Join(hashflagDIR, deactivatedDIR))),
			green(hashflagDirPath()),
		), true},
		{"-force", color.Sprintf("<red>Clear directory</> %s and download all active hashflags", red(hashflagDirPath())), true},
	}
	fmt.Println("Available commands:")
	for _, op := range opList {
		opText := cyan(op.op)
		if op.flag {
			opText = color.Sprintf("\t%s", yellow(op.op))
		}
		color.Printf("\t%s: %s\n", opText, op.desc)
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

func hashflagDirPath() string {
	return dirPath(hashflagDIR)
}

func dirPath(dir string) string {
	return fmt.Sprintf("%s%s", string(os.PathSeparator), dir)
}
