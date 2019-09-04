package main

import (
	"fmt"
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

func dirTree(out *os.File, path string, printFiles bool) error {
	var deepnessNode string
	return drawTree(out, path, printFiles, deepnessNode)
}

func drawTree(out *os.File, path string, printFiles bool, deepnessNode string) error {
	filesInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	lastElement := filesInfo[len(filesInfo)-1]

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
				fileSize = string(fileInfo.Size())
			}
			outputFormat = fmt.Sprintf("%s%s──────%s (%sb)\n", deepnessNode, nodeType, fileInfo.Name(), fileSize)
		} else {
			outputFormat = fmt.Sprintf("%s%s──────%s\n", deepnessNode, nodeType, fileInfo.Name())
		}

		out.WriteString(outputFormat)

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
	//Если папка пуста
	if len(filesInfo) > 0 {
		return true
	} else {
		return false
	}
}
