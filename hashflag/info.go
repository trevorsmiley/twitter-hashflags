package hashflag

import (
	"fmt"
	"log"
	"os"
)

func Info(hashflags []Hashflag, hashflagDIR, detailsFile string, fullDetails, deactivated, fileDiff bool) {
	if fileDiff {
		diff(hashflags, hashflagDIR)
	} else if fullDetails {
		listFullDetails(hashflags, detailsFile)

	} else if deactivated {
		listDeactivated(hashflags, hashflagDIR)
	} else {
		for _, hf := range hashflags {
			fmt.Printf("%s\n", hf.GetFileName())
		}
	}
}

func listFullDetails(hashflags []Hashflag, detailsFile string) {
	tmpl, err := GetTemplate()
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

func diff(hashflags []Hashflag, dir string) {
	missingHashflags := FilterMissingHashflags(hashflags, dir)
	numMissing := red(len(missingHashflags))
	if len(missingHashflags) == 0 {
		numMissing = green(len(missingHashflags))
	}
	fmt.Printf("%s missing hashflags\n\n", numMissing)
	for _, hf := range missingHashflags {
		fmt.Printf("%s\n", hf.GetFileName())
	}
}

func listDeactivated(hashflags []Hashflag, hashflagDIR string) {
	deactivated := FindDeactivatedHashflags(hashflags, hashflagDIR)
	fmt.Printf("%s deactivated Hashflags\n", red(len(deactivated)))
	for _, filename := range deactivated {
		fmt.Printf("%s\n", filename)
	}
}
