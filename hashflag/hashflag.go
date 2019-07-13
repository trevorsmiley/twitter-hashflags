package hashflag

import (
	"log"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/trevorsmiley/fileutils"
)

type Hashflag struct {
	URL      url.URL
	Hashtags []string
	Name     string
	Filename string
}

func (h *Hashflag) GetFileName() string {
	if h.Filename == "" {
		u := h.URL
		path := u.Path
		splitPath := strings.Split(path, "/")
		h.Filename = splitPath[len(splitPath)-1]
	}
	return h.Filename
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
		log.Fatal(err)
	}
}

func SortHashflags(hashflags []Hashflag) {
	sort.Slice(hashflags, func(i, j int) bool {
		return strings.ToLower(hashflags[i].GetName()) < strings.ToLower(hashflags[j].GetName())
	})
}
