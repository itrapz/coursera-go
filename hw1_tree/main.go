package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"

	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var deepnessNode string
	return drawTree(out, path, printFiles, deepnessNode)
}

func drawTree(out io.Writer, path string, printFiles bool, deepnessNode string) error {
	filesInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	lastElement := getLastElement(filesInfo, printFiles)
	for _, fileInfo := range filesInfo {
		var outputFormat string
		var nodeType string
		var deepnessPattern string

		if lastElement == fileInfo {
			nodeType = "└"
			deepnessPattern = " \t"
		} else {
			nodeType = "├"
			deepnessPattern = "│\t"
		}
		if printFiles && !fileInfo.IsDir() {
			var fileSize string
			if fileInfo.Size() == 0 {
				fileSize = "empty"
			} else {
				fileSize = fmt.Sprintf("%db", fileInfo.Size())
			}
			outputFormat = fmt.Sprintf("%s%s───%s (%s)\n", deepnessNode, nodeType, fileInfo.Name(), fileSize)
		} else if fileInfo.IsDir() {
			outputFormat = fmt.Sprintf("%s%s───%s\n", deepnessNode, nodeType, fileInfo.Name())
		} else {
			continue
		}
		out.Write([]byte(outputFormat))
		if fileInfo.IsDir() && dirHasFiles(path+"/"+fileInfo.Name()) {
			drawTree(out, path+"/"+fileInfo.Name(), printFiles, deepnessNode+deepnessPattern)
		}
	}
	return nil
}

func dirHasFiles(path string) bool {
	filesInfo, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	if len(filesInfo) > 0 { //Если папка пуста
		return true
	} else {
		return false
	}
}

func getLastElement(filesInfo []os.FileInfo, printFiles bool) os.FileInfo {
	var lastElement os.FileInfo
	if printFiles {
		lastElement = filesInfo[len(filesInfo)-1]
	} else {
		for _, fileInfo := range filesInfo {
			if fileInfo.IsDir() {
				lastElement = fileInfo
			}
		}
	}
	return lastElement
}
