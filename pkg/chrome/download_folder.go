package chrome

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

type DownloadFolder struct {
	Directory string
}

func (df DownloadFolder) LatestFileWithPrefix(prefix string) (string, error) {
	fileList, err := ioutil.ReadDir(df.Directory)
	if err != nil {
		return "", err
	}

	var files []os.FileInfo
	for _, f := range fileList {
		if strings.HasPrefix(f.Name(), prefix) {
			files = append(files, f)
		}
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Before(files[j].ModTime())
	})

	return fmt.Sprintf("%s/%s", df.Directory, files[len(files)-1].Name()), nil
}
