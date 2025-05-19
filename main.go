package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nyudlts/bytemath"
)

var root string

const version = "0.1.0-alpha"

func init() {
	flag.StringVar(&root, "root", "", "Root directory for the project")
}

func main() {
	fmt.Printf("ddp version %s\n", version)
	flag.Parse()

	outFile, err := os.Create("output.tsv")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	writer.Comma = '\t'

	writer.Write([]string{"DIRECTORY", "FILE_COUNT", "FILE_SIZE", "FILE_SIZE_HUMAN"})

	if _, err := os.Stat(root); err != nil {
		panic(err.Error())
	}

	var totalFiles int = 0
	var totalFilesSize int64 = 0
	var totalDirectories int = 0

	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			totalDirectories++
			filesCount, filesSize, err := getCountFiles(path)
			if err != nil {
				return (err)
			}
			totalFiles += filesCount
			totalFilesSize += filesSize
			writer.Write([]string{path, fmt.Sprintf("%d", filesCount), fmt.Sprintf("%d", filesSize), bytemath.ConvertBytesToHumanReadable(filesSize)})
			writer.Flush()
		}
		return nil
	}); err != nil {
		panic(err.Error())
	}

	writer.Write([]string{"Totals", fmt.Sprintf("%d", totalFiles), fmt.Sprintf("%d", totalFilesSize), bytemath.ConvertBytesToHumanReadable(totalFilesSize)})
	writer.Flush()

	fmt.Printf("Found %d files (%s) in %d directories", totalFiles, bytemath.ConvertBytesToHumanReadable(totalFilesSize), totalDirectories)
}

func getCountFiles(path string) (int, int64, error) {
	childFileCount := 0
	var childFileSize int64 = 0
	children, err := os.ReadDir(path)
	if err != nil {
		return childFileCount, childFileSize, err
	}
	for _, child := range children {
		if !child.IsDir() {
			childFileCount++
			fi, err := os.Stat(filepath.Join(path, child.Name()))
			if err != nil {
				return childFileCount, childFileSize, err
			}
			childFileSize += fi.Size()
		}
	}
	return childFileCount, childFileSize, nil
}
