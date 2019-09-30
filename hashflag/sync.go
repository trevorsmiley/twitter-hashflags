package hashflag

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/gookit/color"
	"github.com/trevorsmiley/fileutils"
	"log"
	"os"
	"path/filepath"
)

var cyan = color.FgCyan.Render
var green = color.FgGreen.Render
var red = color.FgRed.Render

func Sync(hashflags []Hashflag, hashflagDIR, deactivatedDIR string, moveInactive, force bool) {
	if force {
		forceDownload(hashflags, hashflagDIR)
		return
	}

	if moveInactive {
		moveDeactivated(hashflags, hashflagDIR, deactivatedDIR)
	}

	missingHashflags := FilterMissingHashflags(hashflags, hashflagDIR)
	if len(missingHashflags) > 0 {
		fmt.Printf("Syncing %s hashflags to /%s\n", green(len(missingHashflags)), cyan(hashflagDIR))
		downloadAll(missingHashflags, false, hashflagDIR)
		color.FgGreen.Println("Complete")
	} else {
		fmt.Println("No new hashflags to download")
	}
}

func forceDownload(hashflags []Hashflag, dir string) {
	fmt.Printf("Downloading all hashflags to /%s\n", dir)
	downloadAll(hashflags, true, dir)
}

func moveDeactivated(hashflags []Hashflag, hashflagDIR, deactivatedDIR string) {
	deactivated := FindDeactivatedHashflags(hashflags, hashflagDIR)
	fmt.Printf("%s deactivated Hashflags\n", red(len(deactivated)))
	fmt.Printf("Moving deactivated Hashflags into %s\n", green(filepath.Join(hashflagDIR, deactivatedDIR)))
	err := fileutils.CreateDirIfMissing(filepath.Join(hashflagDIR, deactivatedDIR))
	if err != nil {
		log.Fatal(err)
	}
	err = moveFiles(deactivated, filepath.Join(hashflagDIR), filepath.Join(hashflagDIR, deactivatedDIR))
	if err != nil {
		log.Fatal(err)
	}
	color.Green.Println("Complete")
}

func downloadAll(hashflags []Hashflag, clearDIR bool, dir string) {
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
	for _, hf := range hashflags{
		color.Printf("%s\n", green(hf.GetFileName()))
	}
}

func moveFiles(filenames []string, oldDIR, newDIR string) error {
	fmt.Println("Moving...")
	for _, filename := range filenames {
		color.Red.Println(filename)
		err := os.Rename(filepath.Join(oldDIR, filename), filepath.Join(newDIR, filename))
		if err != nil {
			return err
		}
	}
	return nil
}
