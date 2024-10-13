package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	_root, _ := os.Getwd()
	rootPath, err := filepath.Abs(_root)

	if err != nil {
		panic(err)
	}

	files, err := filepath.Glob(path.Join(rootPath, "pl/*.go"))
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		var fileName = strings.Split(file, "/")[len(strings.Split(file, "/"))-1]
		fileNameWithoutExt := strings.Split(fileName, ".")[0]
		// fileExtension := strings.Split(fileName, ".")[1]
		outputFile := fmt.Sprintf("pl/%s.so", fileNameWithoutExt)
		inputFile := fmt.Sprintf("pl/%s.go", fileNameWithoutExt)

		cmd := exec.Command("go", "build", "-buildmode=plugin", "--trimpath", "-o", outputFile, inputFile)

		err := cmd.Start()

		if err != nil {
			panic(err)
		}
		fmt.Printf("Built module: %s \n", fileNameWithoutExt)

	}
}
