package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"sort"
	"strconv"
)

const prefix string = "├───"
const lastPrefix string = "└───"

type File struct {
	Name  string
	Size  int64
	IsDir bool
}

func rec(curPath string, curIden string, full bool, outBuf *bytes.Buffer) error {
	fl, err := ioutil.ReadDir(curPath)
	if err != nil {
		return err
	}
	files := make([]File, 0)
	for _, f := range fl {
		if f.IsDir() {
			files = append(files, File{
				Name:  f.Name(),
				IsDir: true,
			})
		} else {
			if full {
				files = append(files, File{
					Name:  f.Name(),
					Size:  f.Size(),
					IsDir: false,
				})
			}
		}
	}
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})
	for ind, curFile := range files {
		curPrefix := ""
		if ind == len(files)-1 {
			curPrefix = lastPrefix
		} else {
			curPrefix = prefix
		}
		if curFile.IsDir {
			outBuf.WriteString(curIden + curPrefix + curFile.Name + "\n")
			newIden := ""
			if ind == len(files)-1 {
				newIden = curIden + "\t"
			} else {
				newIden = curIden + "│\t"
			}
			newPath := curPath + "/" + curFile.Name
			rec(newPath, newIden, full, outBuf)
		} else {
			if full {
				outBuf.WriteString(curIden + curPrefix + curFile.Name)
				if curFile.Size == 0 {
					outBuf.WriteString(" (empty)")
				} else {
					outBuf.WriteString(" (" + strconv.Itoa(int(curFile.Size)) + "b)")
				}
				outBuf.WriteString("\n")
			}
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, full bool) error {
	outBuf := bytes.NewBuffer(nil)
	err := rec(path, "", full, outBuf)
	if err != nil {
		return err
	}
	_, err = out.Write(outBuf.Bytes())
	if err != nil {
		return err
	}
	return nil
}
