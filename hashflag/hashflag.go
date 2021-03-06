package hashflag

import (
	"html/template"
	"log"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"twitter-hashflags/utils"

	"github.com/trevorsmiley/fileutils"
)

type Hashflag struct {
	URL      url.URL
	Hashtags []string
	Name     string
	Filename string
}

func GetTextOutputTemplate() (*template.Template, error) {
	t := `{{ $hashflags := . }}
Active Hashflags
-----------------
{{- range $index, $hashflag := $hashflags }}
File: {{ .GetFileName }}
URL: {{ .URL.String }}
Hashtags: {{ len .Hashtags }}
	{{ StringsJoin .Hashtags ", " }}
{{end}}`
	return getTemplate(t, "hashflags.tmpl")
}

func GetHtmlTemplate() (*template.Template, error) {
	t := `{{ $hashflags := . }}
<html>
<body>
{{- range $index, $hashflag := $hashflags }}
<img src="{{ .URL.String }}" title="{{ .GetFileName }}: {{ StringsJoin .Hashtags ", " }}"/>
{{end}}
</body>
</html>`
	return getTemplate(t, "html.tmpl")
}

func getTemplate(tmpl, name string) (*template.Template, error) {
	return template.New(name).Funcs(template.FuncMap{"StringsJoin": strings.Join}).Parse(tmpl)
}

func (h *Hashflag) GetFileName() string {
	if h.Filename == "" {
		return GetFileNameFromURL(h.URL)
	}
	return h.Filename
}

func GetFileName(uri string) string {
	u, _ := url.Parse(uri)
	return GetFileNameFromURL(*u)
}

func GetFileNameFromURL(u url.URL) string {
	path := u.Path
	splitPath := strings.Split(path, "/")
	return splitPath[len(splitPath)-1]
}

func (h *Hashflag) GetName() string {
	if h.Name == "" {
		filename := h.GetFileName()
		s := strings.Split(filename, ".")
		h.Name = s[0]
	}
	return h.Name
}

func (h *Hashflag) GetFileExtension() string {
	filename := h.GetFileName()
	return filepath.Ext(filename)
}

func (h *Hashflag) Download(dir string) {

	err := fileutils.DownloadFile(filepath.Join(dir, h.GetFileName()), h.URL.String())
	if err != nil {
		log.Fatalf("Error downloading file %s\n%v", h.URL.String(), err)
	}
}

func SortHashflags(hashflags []Hashflag) {
	sort.Slice(hashflags, func(i, j int) bool {
		return strings.ToLower(hashflags[i].GetName()) < strings.ToLower(hashflags[j].GetName())
	})
}

func FilterMissingHashflags(hashflags []Hashflag, hashflagDIR string) []Hashflag {
	filtered := make([]Hashflag, 0)
	filenames, err := fileutils.GetFileNames(hashflagDIR)
	if err != nil {
		log.Fatalf("Error getting filenames from %s\n%v", hashflagDIR, err)
	}
	for _, hf := range hashflags {
		if utils.ContainsString(filenames, hf.GetFileName()) {
			continue
		}
		filtered = append(filtered, hf)
	}
	SortHashflags(filtered)
	return filtered
}

func FindDeactivatedHashflags(activeHashflags []Hashflag, hashflagDIR string) []string {
	deactivated := make([]string, 0)
	filenames, err := fileutils.GetFileNames(hashflagDIR)
	imageExt := []string{".png"}
	if err != nil {
		log.Fatalf("Error getting filenames from %s\n%v", hashflagDIR, err)
	}
	for _, filename := range filenames {
		if !utils.ContainsString(imageExt, filepath.Ext(filename)) {
			continue
		}
		if !HashflagsContains(activeHashflags, filename) {
			deactivated = append(deactivated, filename)
		}
	}
	return deactivated
}

func HashflagsContains(hashflags []Hashflag, s string) bool {
	for _, hf := range hashflags {
		if hf.GetFileName() == s || hf.GetName() == s {
			return true
		}
	}
	return false
}
