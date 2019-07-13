package utils

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

func GetFileNameWithoutExtension(filename string) string {
	extension := filepath.Ext(filename)
	return filename[0 : len(filename)-len(extension)]
}

func GetFileNames(dir string) ([]string, error) {
	filenames := make([]string, 0)
	files, err := ioutil.ReadDir("./" + dir)
	if err != nil {
		return filenames, err
	}

	for _, f := range files {
		filenames = append(filenames, f.Name())
	}
	sort.Strings(filenames)
	return filenames, nil
}

func RemoveFolderContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer func() {
		err := d.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadFile(filepath, uri string) error {
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		err := out.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	_, err = io.Copy(out, resp.Body)
	return err
}

func CreateDirIfMissing(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.Mkdir(dir, os.ModePerm)
	}
	return nil
}

func CreateOrClearDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.Mkdir(dir, os.ModePerm)
	} else {
		return RemoveFolderContents(dir)
	}
}
