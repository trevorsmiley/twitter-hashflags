package hashflag

import (
	"fmt"
	"html/template"
	"log"
	"os"
)

func Info(hashflags []Hashflag, hashflagDIR, detailsFile, htmlFile string, fullDetails, deactivated, fileDiff, htmlPage bool) {
	if fileDiff {
		diff(hashflags, hashflagDIR)
	} else if fullDetails {
		listFullDetails(hashflags, detailsFile)
	} else if deactivated {
		listDeactivated(hashflags, hashflagDIR)
	}else if htmlPage {
		buildLocalPage(hashflags, htmlFile)
	} else {
		for _, hf := range hashflags {
			fmt.Printf("%s\n", hf.GetFileName())
		}
	}
}

func buildLocalPage(hashflags []Hashflag, htmlFile string){
	tmpl, err := GetHtmlTemplate()
	if err != nil{
		log.Fatal("HTML template error", err)
	}
	writeTemplateToFile(htmlFile, tmpl, hashflags)
}

func writeTemplateToFile(fileName string, tmpl *template.Template, hashflags []Hashflag) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Couldn't create file %s\n%v", fileName, err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatal("Couldn't close file", err)
		}
	}()
	fmt.Printf("Writing details to %s\n", cyan(fileName))
	err = tmpl.Execute(f, hashflags)
	if err != nil {
		log.Fatal("Error executing template", err)
	}
}

func listFullDetails(hashflags []Hashflag, detailsFile string) {
	tmpl, err := GetTextOutputTemplate()
	if err != nil {
		log.Fatal("Error with template", err)
	}
	writeTemplateToFile(detailsFile, tmpl, hashflags)
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
